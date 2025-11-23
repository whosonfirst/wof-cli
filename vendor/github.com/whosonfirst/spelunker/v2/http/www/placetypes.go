package www

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/aaronland/go-http/v4/auth"
	"github.com/aaronland/go-http/v4/slog"
	"github.com/whosonfirst/spelunker/v2"
	wof_http "github.com/whosonfirst/spelunker/v2/http"
)

type placetypesHandlerVars struct {
	PageTitle string
	URIs      *wof_http.URIs
	Facets    []*spelunker.FacetCount
	OpenGraph *OpenGraph
}

// PlacetypesHandlerOptions  defines configuration options for the `PlacetypesHandler` method.
type PlacetypesHandlerOptions struct {
	// An instance implemeting the `spelunker.Spelunker` interface.
	Spelunker spelunker.Spelunker
	// An instance implementing the `aaronland/go-http/v4/auth.Authenticator` interface.
	Authenticator auth.Authenticator
	// An `html/template.Template` instance containing the named template "placetypes".
	Templates *template.Template
	// URIs are the `wof_http.URIs` details for this Spelunker instance.
	URIs *wof_http.URIs
}

// PlacetypesHandler returns an `http.Handler` instance to display webpage listing all the placetypes in Spelunker index.
func PlacetypesHandler(opts *PlacetypesHandlerOptions) (http.Handler, error) {

	t := opts.Templates.Lookup("placetypes")

	if t == nil {
		return nil, fmt.Errorf("Failed to locate 'placetypes' template")
	}

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()
		logger := slog.LoggerWithRequest(req, nil)

		faceting, err := opts.Spelunker.GetPlacetypes(ctx)

		if err != nil {
			logger.Error("Failed to get placetypes", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		vars := placetypesHandlerVars{
			PageTitle: "Placetypes",
			URIs:      opts.URIs,
			Facets:    faceting.Results,
		}

		vars.OpenGraph = &OpenGraph{
			Type:        "Article",
			SiteName:    "Who's On First Spelunker",
			Title:       "Who's On First Placetypes",
			Description: "Who's On First records grouped by their place types",
			Image:       "",
		}

		rsp.Header().Set("Content-Type", "text/html")

		err = t.Execute(rsp, vars)

		if err != nil {
			logger.Error("Failed to render template", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
		}

	}

	h := http.HandlerFunc(fn)
	return h, nil
}
