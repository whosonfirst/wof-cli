package www

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/sfomuseum/go-http-auth"
	"github.com/whosonfirst/go-whosonfirst-spelunker-httpd"
)

type TemplateHandlerOptions struct {
	Authenticator auth.Authenticator
	Templates     *template.Template
	TemplateName  string
	PageTitle     string
	URIs          *httpd.URIs
}

type TemplateHandlerVars struct {
	Id         int64
	PageTitle  string
	URIs       *httpd.URIs
	Properties string
	OpenGraph  *OpenGraph
}

func TemplateHandler(opts *TemplateHandlerOptions) (http.Handler, error) {

	t := opts.Templates.Lookup(opts.TemplateName)

	if t == nil {
		return nil, fmt.Errorf("Failed to locate ihelp' template")
	}

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		logger := httpd.LoggerWithRequest(req, nil)

		vars := TemplateHandlerVars{
			PageTitle: opts.PageTitle,
			URIs:      opts.URIs,
		}

		vars.OpenGraph = &OpenGraph{
			Type:        "Article",
			SiteName:    "Who's On First Spelunker",
			Title:       fmt.Sprintf("Who's On First Spelunker – %s", opts.PageTitle),
			Description: "",
			Image:       "",
		}

		rsp.Header().Set("Content-Type", "text/html")

		err := t.Execute(rsp, vars)

		if err != nil {
			logger.Error("Failed to render template ", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
		}

	}

	h := http.HandlerFunc(fn)
	return h, nil
}
