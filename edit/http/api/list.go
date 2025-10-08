package api

import (
	"encoding/json"
	"io/fs"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/aaronland/go-http/v3/slog"
)

func ListHandler(data_root *os.Root) http.Handler {

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

		logger := slog.LoggerWithRequest(req, nil)
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
