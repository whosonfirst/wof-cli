package api

import (
	"bytes"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/aaronland/go-http/v3/slog"
	"github.com/whosonfirst/go-whosonfirst-export/v3"
	"github.com/whosonfirst/go-whosonfirst-validate"
	"github.com/whosonfirst/go-writer/v3"
)

func SaveHandler(data_root *os.Root, uri_map *sync.Map, wr writer.Writer) http.Handler {

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

		if has_changes {

			logger.Info("Record has changed")

			_, new_body, err := export.Export(ctx, body)

			if err != nil {
				logger.Error("Failed to export body", "error", err)
				http.Error(rsp, "Internal server error", http.StatusInternalServerError)
				return
			}

			_, err = wr.Write(ctx, uri, bytes.NewReader(new_body))

			// Note how this is NOT using whosonfirst/go-whosonfirst-writer that
			// will explicitly write files as 123/456/7/1234567.geojson which is
			// not necessarily the desired effect.

			if err != nil {
				logger.Error("Failed to write body", "error", err)
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
		}

		rsp.Header().Set("Content-type", "application/json")
		rsp.Write([]byte(body))
	}

	return http.HandlerFunc(fn)
}
