package www

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/aaronland/go-http/v4/auth"
	"github.com/aaronland/go-http/v4/slog"
	wof_http "github.com/whosonfirst/spelunker/v2/http"
)

type templateHandlerVars struct {
	Id         int64
	PageTitle  string
	URIs       *wof_http.URIs
	Properties string
	OpenGraph  *OpenGraph
}

// TemplateHandlerOptions  defines configuration options for the `TemplateHandler` method.
type TemplateHandlerOptions struct {
	// An instance implementing the `aaronland/go-http/v4/auth.Authenticator` interface.
	Authenticator auth.Authenticator
	// An `html/template.Template` instance containing the named template defined in the `TemplateName` property.
	Templates *template.Template
	// The name of the template to display.
	TemplateName string
	// The title of the page to display.
	PageTitle string
	// URIs are the `wof_http.URIs` details for this Spelunker instance.
	URIs *wof_http.URIs
}

// TemplateHandler returns an `http.Handler` instance to display webpage defined by a named template.
func TemplateHandler(opts *TemplateHandlerOptions) (http.Handler, error) {

	t := opts.Templates.Lookup(opts.TemplateName)

	if t == nil {
		return nil, fmt.Errorf("Failed to locate ihelp' template")
	}

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		logger := slog.LoggerWithRequest(req, nil)

		vars := templateHandlerVars{
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
