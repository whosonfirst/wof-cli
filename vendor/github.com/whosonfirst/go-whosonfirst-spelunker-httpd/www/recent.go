package www

import (
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"time"

	"github.com/aaronland/go-pagination"
	"github.com/dustin/go-humanize"
	"github.com/sfomuseum/go-http-auth"
	"github.com/sfomuseum/iso8601duration"
	"github.com/whosonfirst/go-whosonfirst-spelunker"
	"github.com/whosonfirst/go-whosonfirst-spelunker-httpd"
	"github.com/whosonfirst/go-whosonfirst-spr/v2"
)

type RecentHandlerOptions struct {
	Spelunker     spelunker.Spelunker
	Authenticator auth.Authenticator
	Templates     *template.Template
	URIs          *httpd.URIs
}

type RecentHandlerVars struct {
	PageTitle     string
	URIs          *httpd.URIs
	Places        []spr.StandardPlacesResult
	Pagination    pagination.Results
	PaginationURL string
	// Duration         time.Duration
	Duration         *duration.Duration
	Since            string
	FacetsURL        string
	FacetsContextURL string
	OpenGraph        *OpenGraph
}

func RecentHandler(opts *RecentHandlerOptions) (http.Handler, error) {

	t := opts.Templates.Lookup("recent")

	if t == nil {
		return nil, fmt.Errorf("Failed to locate 'recent' template")
	}

	re_full, err := regexp.Compile(`P((?P<year>\d+)Y)?((?P<month>\d+)M)?((?P<day>\d+)D)?(T((?P<hour>\d+)H)?((?P<minute>\d+)M)?((?P<second>\d+)S)?)?`)

	if err != nil {
		return nil, fmt.Errorf("Failed to compile ISO8601 duration pattern, %w", err)
	}

	re_week, err := regexp.Compile(`P((?P<week>\d+)W)`)

	if err != nil {
		return nil, fmt.Errorf("Failed to compile ISO8601 duration week pattern, %w", err)
	}

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()
		logger := httpd.LoggerWithRequest(req, nil)

		str_d := req.PathValue("duration")

		switch {
		case re_week.MatchString(str_d):
			// ok
		case re_full.MatchString(str_d):
			// ok
		default:
			str_d = "P30D"
		}

		logger = logger.With("duration", str_d)

		d, err := duration.FromString(str_d)

		if err != nil {
			logger.Error("Failed to parse duration", "error", err)
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

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

		r, pg_r, err := opts.Spelunker.GetRecent(ctx, pg_opts, d.ToDuration(), filters)

		if err != nil {
			logger.Error("Failed to get recent", "error", err)
			http.Error(rsp, "womp womp", http.StatusInternalServerError)
			return
		}

		// This is not ideal but I am not sure what is better yet...
		pagination_url := httpd.URIForRecent(opts.URIs.Recent, str_d, filters, nil)

		// This is not ideal but I am not sure what is better yet...
		facets_url := httpd.URIForRecent(opts.URIs.RecentFaceted, str_d, filters, nil)
		facets_context_url := pagination_url

		now := time.Now()
		now_ts := now.Unix()

		then_ts := now_ts - int64(d.ToDuration().Seconds())
		then := time.Unix(then_ts, 0)

		since := humanize.RelTime(now, then, "", "")

		vars := RecentHandlerVars{
			Places:        r.Results(),
			Pagination:    pg_r,
			URIs:          opts.URIs,
			PaginationURL: pagination_url,
			// Duration:         d.ToDuration(),
			Duration:         d,
			Since:            since,
			FacetsURL:        facets_url,
			FacetsContextURL: facets_context_url,
		}

		vars.OpenGraph = &OpenGraph{
			Type:        "Article",
			SiteName:    "Who's On First Spelunker",
			Title:       "Who's On First recently updated records",
			Description: fmt.Sprintf("Who's On First records that have been updated since %s", since),
			Image:       "",
		}

		rsp.Header().Set("Content-Type", "text/html")

		err = t.Execute(rsp, vars)

		if err != nil {
			logger.Error("Failed to return ", "error", err)
			http.Error(rsp, "womp womp", http.StatusInternalServerError)
		}

	}

	h := http.HandlerFunc(fn)
	return h, nil
}
