package api

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/aaronland/go-http/v3/slog"
	"github.com/google/uuid"
	"github.com/whosonfirst/go-whosonfirst-export/v3"
	"github.com/whosonfirst/go-whosonfirst-feature/properties"
	"github.com/whosonfirst/go-whosonfirst-validate"
	gh_writer "github.com/whosonfirst/go-writer-github/v3"
	"github.com/whosonfirst/go-writer/v3"
)

func SaveHandler(data_root *os.Root, uri_map *sync.Map, writer_uri string) http.Handler {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		logger := slog.LoggerWithRequest(req, nil)

		if req.Method != http.MethodPost {
			http.Error(rsp, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		logger.Info("Save record")

		fname := filepath.Base(req.URL.Path)

		v, exists := uri_map.Load(fname)

		if !exists {
			logger.Error("Failed to find URI in lookup", "fname", fname)
			http.Error(rsp, "Not found", http.StatusNotFound)
			return
		}

		uri := v.(string)
		logger = logger.With("uri", uri)

		old_body, err := data_root.ReadFile(fname)

		if err != nil {
			logger.Error("Failed to stat URI", "fname", fname, "error", err)
			http.Error(rsp, "Not found", http.StatusNotFound)
			return
		}

		body, err := validate.EnsureValidGeoJSON(req.Body)

		if err != nil {
			logger.Error("Failed to ensure GeoJSON", "error", err)
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		err = validate.Validate(body)

		if err != nil {
			logger.Error("Failed to validate body", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		has_changes, err := export.HasChanges(ctx, old_body, body)

		if err != nil {
			logger.Error("Failed to determine if body has changes", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Return 304 if !has_changes?

		if !has_changes {
			rsp.Header().Set("Content-type", "application/json")
			rsp.Write([]byte(body))
			return
		}

		logger.Info("Record has changed")

		_, new_body, err := export.Export(ctx, body)

		if err != nil {
			logger.Error("Failed to export body", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		if strings.HasPrefix(writer_uri, gh_writer.GITHUBAPI_PR_SCHEME) {

			id, err := properties.Id(new_body)

			if err != nil {
				logger.Error("Failed to derive wof:id from update", "error", err)
				http.Error(rsp, "Internal server error", http.StatusInternalServerError)
				return
			}

			name, err := properties.Name(new_body)

			if err != nil {
				logger.Error("Failed to derive wof:name from update", "error", err)
				http.Error(rsp, "Internal server error", http.StatusInternalServerError)
				return
			}

			pt, err := properties.Placetype(new_body)

			if err != nil {
				logger.Error("Failed to derive wof:placetype from update", "error", err)
				http.Error(rsp, "Internal server error", http.StatusInternalServerError)
				return
			}

			repo, err := properties.Repo(new_body)

			if err != nil {
				logger.Error("Failed to derive wof:repo from update", "error", err)
				http.Error(rsp, "Internal server error", http.StatusInternalServerError)
				return
			}

			u, err := url.Parse(writer_uri)

			if err != nil {
				logger.Error("Failed to parse githubapi-pr writer URI", "error", err)
				http.Error(rsp, "Internal server error", http.StatusInternalServerError)
				return
			}

			q := u.Query()

			if q.Has("pr-branch") {
				q.Del("pr-branch")
			}

			uid := uuid.New()

			q.Set("pr-branch", fmt.Sprintf("wof-cli-edit-%d-%s", id, uid.String()))

			if q.Has("pr-title") {
				q.Del("pr-title")
			}

			q.Set("pr-title", fmt.Sprintf("Update %s %s", pt, name))

			if q.Has("pr-description") {
				q.Del("pr-description")
			}

			q.Set("pr-description", fmt.Sprintf("Updating %s using wof-cli edit", fname))

			u.Path = repo
			u.RawQuery = q.Encode()

			writer_uri = u.String()
		}

		wr, err := writer.NewWriter(ctx, writer_uri)

		if err != nil {
			logger.Error("Failed to create new writer", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		logger.Debug("Write record", "uri", uri)

		_, err = wr.Write(ctx, uri, bytes.NewReader(new_body))

		if err != nil {
			logger.Error("Failed to write body", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		err = wr.Close(ctx)

		if err != nil {
			logger.Error("Failed to close writer after writing", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Now write back to the temporary directory with the files being edited

		err = data_root.WriteFile(fname, new_body, 0644)

		if err != nil {
			logger.Error("Failed to write body to data root", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		body = new_body

		rsp.Header().Set("Content-type", "application/json")
		rsp.Write([]byte(body))
	}

	return http.HandlerFunc(fn)
}
