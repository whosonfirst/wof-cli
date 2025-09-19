package edit

import (
	"encoding/json"
	"io/fs"
	"log/slog"
	"net/http"
	"strings"
	"sync"

	"github.com/whosonfirst/wof/edit/static"
)

func staticHandler() http.Handler {
	http_fs := http.FS(static.FS)
	return http.FileServer(http_fs)
}

func dataHandler(data_fs fs.FS) http.Handler {
	http_fs := http.FS(data_fs)
	return http.FileServer(http_fs)
}

func apiListHandler(data_fs fs.FS) http.Handler {

	once := sync.OnceValues(func() ([]string, error) {

		paths := make([]string, 0)

		err := fs.WalkDir(data_fs, ".", func(path string, d fs.DirEntry, err error) error {

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
