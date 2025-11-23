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

type concordancesHandlerVars struct {
	PageTitle string
	URIs      *wof_http.URIs
	Facets    []*spelunker.FacetCount
	OpenGraph *OpenGraph
}

// ConcordanceHandlerOptions defines configuration options for the `ConcordancesHandler` method.
type ConcordancesHandlerOptions struct {
	// An instance implemeting the `spelunker.Spelunker` interface.
	Spelunker spelunker.Spelunker
	// An instance implementing the `aaronland/go-http/v4/auth.Authenticator` interface.
	Authenticator auth.Authenticator
	// An `html/template.Template` instance containing the named template "concordances".
	Templates *template.Template
	// URIs are the `wof_http.URIs` details for this Spelunker instance.
	URIs *wof_http.URIs
}

// ConcordancesHandler returns an `http.Handler` instance to display a webpage listing all the concordances in a Spelunker index.
func ConcordancesHandler(opts *ConcordancesHandlerOptions) (http.Handler, error) {

	t := opts.Templates.Lookup("concordances")

	if t == nil {
		return nil, fmt.Errorf("Failed to locate 'concordances' template")
	}

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()
		logger := slog.LoggerWithRequest(req, nil)

		faceting, err := opts.Spelunker.GetConcordances(ctx)

		if err != nil {
			logger.Error("Failed to get concordances", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		vars := concordancesHandlerVars{
			PageTitle: "Concordances",
			URIs:      opts.URIs,
			Facets:    faceting.Results,
		}

		vars.OpenGraph = &OpenGraph{
			Type:        "Article",
			SiteName:    "Who's On First Spelunker",
			Title:       "Concordances with Who's On First",
			Description: `Other data sources that Who's On First "holds hands" with`,
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
