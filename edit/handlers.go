package edit

import (
	"net/http"

	"github.com/whosonfirst/wof/edit/static"	
)

func staticHandler() http.Handler {

	http_fs := http.FS(static.FS)
	return http.FileServer(http_fs)
}

