package www

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/aaronland/go-pagination"
	"github.com/sfomuseum/go-http-auth"
	"github.com/whosonfirst/go-whosonfirst-spelunker"
	"github.com/whosonfirst/go-whosonfirst-spelunker-httpd"
	"github.com/whosonfirst/go-whosonfirst-spr/v2"
)

type NullIslandHandlerOptions struct {
	Spelunker     spelunker.Spelunker
	Authenticator auth.Authenticator
	Templates     *template.Template
	URIs          *httpd.URIs
}

type NullIslandHandlerVars struct {
	PageTitle        string
	URIs             *httpd.URIs
	Places           []spr.StandardPlacesResult
	Pagination       pagination.Results
	PaginationURL    string
	FacetsURL        string
	FacetsContextURL string
	OpenGraph        *OpenGraph
}

func NullIslandHandler(opts *NullIslandHandlerOptions) (http.Handler, error) {

	t := opts.Templates.Lookup("nullisland")

	if t == nil {
		return nil, fmt.Errorf("Failed to locate 'recent' template")
	}

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()
		logger := httpd.LoggerWithRequest(req, nil)

		pg_opts, err := httpd.PaginationOptionsFromRequest(req)

		if err != nil {
			logger.Error("Failed to create pagination options", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		filter_params := httpd.DefaultFilterParams()

		filters, err := httpd.FiltersFromRequest(ctx, req, filter_params)

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
		pagination_url := httpd.URIForNullIsland(opts.URIs.NullIsland, filters, nil)

		// This is not ideal but I am not sure what is better yet...
		facets_url := httpd.URIForNullIsland(opts.URIs.NullIslandFaceted, filters, nil)
		facets_context_url := pagination_url

		vars := NullIslandHandlerVars{
			Places:           r.Results(),
			Pagination:       pg_r,
			URIs:             opts.URIs,
			PaginationURL:    pagination_url,
			FacetsURL:        facets_url,
			FacetsContextURL: facets_context_url,
		}

		svg_url := httpd.URIForIdSimple(opts.URIs.SVG, 0)

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
