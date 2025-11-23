package www

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/aaronland/go-http/v4/auth"
	"github.com/aaronland/go-http/v4/slog"
	"github.com/aaronland/go-pagination"
	"github.com/whosonfirst/go-whosonfirst-spr/v2"
	"github.com/whosonfirst/spelunker/v2"
	wof_http "github.com/whosonfirst/spelunker/v2/http"
)

type nullIslandHandlerVars struct {
	PageTitle        string
	URIs             *wof_http.URIs
	Places           []spr.StandardPlacesResult
	Pagination       pagination.Results
	PaginationURL    string
	FacetsURL        string
	FacetsContextURL string
	OpenGraph        *OpenGraph
}

// NullIslandHandlerOptions  defines configuration options for the `NullIslandHandler` method.
type NullIslandHandlerOptions struct {
	// An instance implemeting the `spelunker.Spelunker` interface.
	Spelunker spelunker.Spelunker
	// An instance implementing the `aaronland/go-http/v4/auth.Authenticator` interface.
	Authenticator auth.Authenticator
	// An `html/template.Template` instance containing the named template "nullisland".
	Templates *template.Template
	// URIs are the `wof_http.URIs` details for this Spelunker instance.
	URIs *wof_http.URIs
}

// NullIslandHandler returns an `http.Handler` instance to display webpage listing Who's On First records "visiting" Null Island (with lat,lon coordinates of "0.0,0.0").
func NullIslandHandler(opts *NullIslandHandlerOptions) (http.Handler, error) {

	t := opts.Templates.Lookup("nullisland")

	if t == nil {
		return nil, fmt.Errorf("Failed to locate 'recent' template")
	}

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()
		logger := slog.LoggerWithRequest(req, nil)

		pg_opts, err := wof_http.PaginationOptionsFromRequest(req)

		if err != nil {
			logger.Error("Failed to create pagination options", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		filter_params := wof_http.DefaultFilterParams()

		filters, err := wof_http.FiltersFromRequest(ctx, req, filter_params)

		if err != nil {
			logger.Error("Failed to derive filters from request", "error", err)
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		r, pg_r, err := opts.Spelunker.VisitingNullIsland(ctx, pg_opts, filters)

		if err != nil {
			logger.Error("Failed to get recent", "error", err)
			http.Error(rsp, "InternalServerError", http.StatusInternalServerError)
			return
		}

		// This is not ideal but I am not sure what is better yet...
		pagination_url := wof_http.URIForNullIsland(opts.URIs.NullIsland, filters, nil)

		// This is not ideal but I am not sure what is better yet...
		facets_url := wof_http.URIForNullIsland(opts.URIs.NullIslandFaceted, filters, nil)
		facets_context_url := pagination_url

		vars := nullIslandHandlerVars{
			Places:           r.Results(),
			Pagination:       pg_r,
			URIs:             opts.URIs,
			PaginationURL:    pagination_url,
			FacetsURL:        facets_url,
			FacetsContextURL: facets_context_url,
		}

		svg_url := wof_http.URIForIdSimple(opts.URIs.SVG, 0)

		og_image, err := opts.URIs.Abs(svg_url)

		if err != nil {
			logger.Error("Failed to derive absolute URL for SVG image", "url", svg_url, "error", err)
		}

		vars.OpenGraph = &OpenGraph{
			Type:        "Article",
			SiteName:    "Who's On First Spelunker",
			Title:       `Who's On First records that are "visiting" Null Island`,
			Description: "Who's On First records with missing or undetermined geographies",
			Image:       og_image,
		}

		rsp.Header().Set("Content-Type", "text/html")

		err = t.Execute(rsp, vars)

		if err != nil {
			logger.Error("Failed to return ", "error", err)
			http.Error(rsp, "InternalServerError", http.StatusInternalServerError)
		}

	}

	h := http.HandlerFunc(fn)
	return h, nil
}
