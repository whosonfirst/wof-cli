package edit

import (
	"bytes"
	"encoding/json"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/whosonfirst/go-whosonfirst-export/v3"
	"github.com/whosonfirst/go-whosonfirst-validate"
	"github.com/whosonfirst/go-writer/v3"
	"github.com/whosonfirst/wof/edit/static"
)

func staticHandler() http.Handler {
	http_fs := http.FS(static.FS)
	return http.FileServer(http_fs)
}

func dataHandler(data_root *os.Root) http.Handler {
	http_fs := http.FS(data_root.FS())
	return http.FileServer(http_fs)
}

func apiListHandler(data_root *os.Root) http.Handler {

	once := sync.OnceValues(func() ([]string, error) {

		paths := make([]string, 0)

		err := fs.WalkDir(data_root.FS(), ".", func(path string, d fs.DirEntry, err error) error {

			if err != nil {
				return err
			}

			if strings.HasPrefix(path, ".") {
				return nil
			}

			paths = append(paths, path)
			return nil
		})

		return paths, err
	})

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		logger := slog.Default()
		logger = logger.With("uri", req.URL.Path)

		paths, err := once()

		if err != nil {
			logger.Error("Failed to derive paths", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		rsp.Header().Set("Content-type", "application/json")

		enc := json.NewEncoder(rsp)
		err = enc.Encode(paths)

		if err != nil {
			logger.Error("Failed to encode paths", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

	}

	return http.HandlerFunc(fn)
}

func apiSaveHandler(data_root *os.Root, wr writer.Writer) http.Handler {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		logger := slog.Default()
		logger = logger.With("uri", req.URL.Path)

		if req.Method != http.MethodPost {
			http.Error(rsp, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		logger.Info("Save record")

		fname := filepath.Base(req.URL.Path)

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

		if has_changes {

			logger.Info("Record has changed")

			_, new_body, err := export.Export(ctx, body)

			if err != nil {
				logger.Error("Failed to export body", "error", err)
				http.Error(rsp, "Internal server error", http.StatusInternalServerError)
				return
			}

			// Note how this is NOT using whosonfirst/go-whosonfirst-writer that
			// will explicitly write files as 123/456/7/1234567.geojson which is
			// not necessarily the desired effect.
			
			br := bytes.NewReader(new_body)
			_, err = wr.Write(ctx, fname, br)

			if err != nil {
				logger.Error("Failed to write body", "error", err)
				http.Error(rsp, "Internal server error", http.StatusInternalServerError)
				return
			}

			err = data_root.WriteFile(fname, new_body, 0644)

			if err != nil {
				logger.Error("Failed to write body to data root", "error", err)
				http.Error(rsp, "Internal server error", http.StatusInternalServerError)
				return
			}

			body = new_body
		}

		rsp.Header().Set("Content-type", "application/json")
		rsp.Write([]byte(body))
	}

	return http.HandlerFunc(fn)
}
