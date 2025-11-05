package www

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/whosonfirst/go-whosonfirst-spelunker-httpd"
)

type RobotsTxtHandlerOptions struct {
	Templates *template.Template
}

func RobotsTxtHandler(opts *RobotsTxtHandlerOptions) (http.Handler, error) {

	t_name := "robots_txt"
	t := opts.Templates.Lookup(t_name)

	if t == nil {
		return nil, fmt.Errorf("Failed to locate '%s' template", t_name)
	}

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		logger := httpd.LoggerWithRequest(req, nil)

		rsp.Header().Set("Content-Type", "text/plain")

		err := t.Execute(rsp, nil)

		if err != nil {
			logger.Error("Failed to render template ", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
		}

	}

	h := http.HandlerFunc(fn)
	return h, nil
}
