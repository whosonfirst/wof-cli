package www

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/sfomuseum/go-http-auth"
	"github.com/whosonfirst/go-whosonfirst-spelunker"
	"github.com/whosonfirst/go-whosonfirst-spelunker-httpd"
)

type ConcordancesHandlerOptions struct {
	Spelunker     spelunker.Spelunker
	Authenticator auth.Authenticator
	Templates     *template.Template
	URIs          *httpd.URIs
}

type ConcordancesHandlerVars struct {
	PageTitle string
	URIs      *httpd.URIs
	Facets    []*spelunker.FacetCount
	OpenGraph *OpenGraph
}

func ConcordancesHandler(opts *ConcordancesHandlerOptions) (http.Handler, error) {

	t := opts.Templates.Lookup("concordances")

	if t == nil {
		return nil, fmt.Errorf("Failed to locate 'concordances' template")
	}

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()
		logger := httpd.LoggerWithRequest(req, nil)

		faceting, err := opts.Spelunker.GetConcordances(ctx)

		if err != nil {
			logger.Error("Failed to get concordances", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		vars := ConcordancesHandlerVars{
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
