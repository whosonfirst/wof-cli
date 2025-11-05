package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	// "github.com/sfomuseum/go-http-auth"
	"github.com/sfomuseum/iso8601duration"
	"github.com/whosonfirst/go-whosonfirst-spelunker"
	"github.com/whosonfirst/go-whosonfirst-spelunker-httpd"
)

type RecentFacetedHandlerOptions struct {
	Spelunker spelunker.Spelunker
	// Authenticator auth.Authenticator
}

func RecentFacetedHandler(opts *RecentFacetedHandlerOptions) (http.Handler, error) {

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

		filter_params := httpd.DefaultFilterParams()

		filters, err := httpd.FiltersFromRequest(ctx, req, filter_params)

		if err != nil {
			logger.Error("Failed to derive filters from request", "error", err)
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		facets, err := httpd.FacetsFromRequest(ctx, req, filter_params)

		if err != nil {
			logger.Error("Failed to derive facets from requrst", "error", err)
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		if len(facets) == 0 {
			logger.Error("No facets from requrst")
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		facets_rsp, err := opts.Spelunker.GetRecentFaceted(ctx, d.ToDuration(), filters, facets)

		if err != nil {
			logger.Error("Failed to get recent", "error", err)
			http.Error(rsp, "womp womp", http.StatusInternalServerError)
			return
		}

		rsp.Header().Set("Content-Type", "application/json")

		enc := json.NewEncoder(rsp)
		err = enc.Encode(facets_rsp)

		if err != nil {
			logger.Error("Failed to encode facets response", "error", err)
			http.Error(rsp, "womp womp", http.StatusInternalServerError)
			return
		}
	}

	h := http.HandlerFunc(fn)
	return h, nil
}
