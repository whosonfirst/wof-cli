package www

import (
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"strconv"

	"github.com/aaronland/go-http-sanitize"
	"github.com/aaronland/go-pagination"
	"github.com/sfomuseum/go-http-auth"
	"github.com/whosonfirst/go-whosonfirst-spelunker"
	"github.com/whosonfirst/go-whosonfirst-spelunker-httpd"
	"github.com/whosonfirst/go-whosonfirst-spr/v2"
	"github.com/whosonfirst/go-whosonfirst-uri"
)

type SearchHandlerOptions struct {
	Spelunker     spelunker.Spelunker
	Authenticator auth.Authenticator
	Templates     *template.Template
	URIs          *httpd.URIs
}

type SearchHandlerVars struct {
	PageTitle        string
	URIs             *httpd.URIs
	Places           []spr.StandardPlacesResult
	Pagination       pagination.Results
	PaginationURL    string
	FacetsURL        string
	FacetsContextURL string
	Feature          spr.StandardPlacesResult
	SearchOptions    *spelunker.SearchOptions
	OpenGraph        *OpenGraph
}

func SearchHandler(opts *SearchHandlerOptions) (http.Handler, error) {

	form_t := opts.Templates.Lookup("search")

	if form_t == nil {
		return nil, fmt.Errorf("Failed to locate 'search' template")
	}

	results_t := opts.Templates.Lookup("search_results")

	if results_t == nil {
		return nil, fmt.Errorf("Failed to locate 'search_results' template")
	}

	re_wofid, err := regexp.Compile(`^\d+$$`)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse WOF ID regular expression, %w", err)
	}

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()
		logger := httpd.LoggerWithRequest(req, nil)

		vars := SearchHandlerVars{
			URIs:      opts.URIs,
			PageTitle: "Search",
		}

		vars.OpenGraph = &OpenGraph{
			Type:        "Article",
			SiteName:    "Who's On First Spelunker",
			Title:       "Who's On First Search",
			Description: "Search for Who's On First records in the Spelunker",
			Image:       "",
		}

		q, err := sanitize.GetString(req, "q")

		if err != nil {
			logger.Error("Failed to determine query string", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		if q == "" {

			rsp.Header().Set("Content-Type", "text/html")

			err = form_t.Execute(rsp, vars)

			if err != nil {
				logger.Error("Failed to return ", "error", err)
				http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			}

			return
		}

		pg_opts, err := httpd.PaginationOptionsFromRequest(req)

		if err != nil {
			logger.Error("Failed to create pagination options", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		search_opts := &spelunker.SearchOptions{
			Query: q,
		}

		filter_params := httpd.DefaultFilterParams()

		filters, err := httpd.FiltersFromRequest(ctx, req, filter_params)

		if err != nil {
			logger.Error("Failed to derive filters from request", "error", err)
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		// TBD - Do this concurrently in Go routines? It kind of feels like yak-shaving
		// at this stage...

		r, pg_r, err := opts.Spelunker.Search(ctx, pg_opts, search_opts, filters)

		if err != nil {
			logger.Error("Failed to get search", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		vars.OpenGraph.Title = fmt.Sprintf(`Search results for \"%s\"`, q)
		vars.OpenGraph.Description = fmt.Sprintf(`Who's On First records matching the query term \"%s\"`, q)

		vars.Places = r.Results()
		vars.Pagination = pg_r

		pagination_url := httpd.URIForSearch(opts.URIs.Search, q, filters, nil)
		facets_url := httpd.URIForSearch(opts.URIs.SearchFaceted, q, filters, nil)
		facets_context_url := httpd.URIForSearch(opts.URIs.Search, q, filters, nil)

		vars.PaginationURL = pagination_url
		vars.FacetsURL = facets_url
		vars.FacetsContextURL = facets_context_url
		vars.SearchOptions = search_opts

		//

		if re_wofid.MatchString(q) {

			var wof_id int64
			var wof_f []byte
			var wof_s spr.StandardPlacesResult

			wofid_ok := true

			if wofid_ok {

				v, err := strconv.ParseInt(q, 10, 64)

				if err != nil {
					logger.Error("Failed to parse ID", "q", q, "error", err)
					wofid_ok = false
				}

				wof_id = v
			}

			// Check min/max here...

			// To do: Replace this with opts.Spelunker.GetSPRForId
			// once the kinks have been worked out

			if wofid_ok {

				uri_args := new(uri.URIArgs)

				f, err := opts.Spelunker.GetFeatureForId(ctx, wof_id, uri_args)

				if err != nil {
					logger.Error("Failed to get by ID", "error", err)
					wofid_ok = false
				}

				wof_f = f
			}

			if wofid_ok {

				v, err := spr.WhosOnFirstSPR(wof_f)

				if err != nil {
					logger.Error("Failed to derive SPR for feature", "id", wof_id, "error", err)
					wofid_ok = false
				}

				wof_s = v
			}

			if wofid_ok {
				vars.Feature = wof_s
			}
		}

		//

		rsp.Header().Set("Content-Type", "text/html")

		err = results_t.Execute(rsp, vars)

		if err != nil {
			logger.Error("Failed to return ", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
		}

	}

	h := http.HandlerFunc(fn)
	return h, nil
}
