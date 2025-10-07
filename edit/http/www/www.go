package www

import (
	"io/fs"
	"net/http"
	"os"
)

func StaticHandler(static_fs fs.FS) http.Handler {
	http_fs := http.FS(static_fs)
	return http.FileServer(http_fs)
}

func DataHandler(data_root *os.Root) http.Handler {
	http_fs := http.FS(data_root.FS())
	return http.FileServer(http_fs)
}
