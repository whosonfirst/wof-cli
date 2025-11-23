package api

import (
	"encoding/json"
	"net/http"

	// TBD...
	// "github.com/aaronland/go-http/v4/auth"

	"github.com/aaronland/go-http/v4/slog"
	wof_http "github.com/whosonfirst/go-whosonfirst/http"
	"github.com/whosonfirst/spelunker/v2"
	sp_http "github.com/whosonfirst/spelunker/v2/http"
)

// DescendantsFacetedHandlerOptions defines options for invoking the `DescendantsFacetedHandler` method.
type DescendantsFacetedHandlerOptions struct {
	// An instance implemeting the `spelunker.Spelunker` interface.
	Spelunker spelunker.Spelunker
	// TBD...
	// Authenticator auth.Authenticator
}

// HasConcordanceFacetedHandler returns an `http.Handler` for returning faceted results for a descendants of a given Who's On First record.
func DescendantsFacetedHandler(opts *DescendantsFacetedHandlerOptions) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()
		logger := slog.LoggerWithRequest(req, nil)

		uri, err, status := wof_http.ParseURIFromRequest(req)

		if err != nil {
			logger.Error("Failed to parse URI from request", "error", err)
			http.Error(rsp, spelunker.ErrNotFound.Error(), status)
			return
		}

		logger = logger.With("wofid", uri.Id)

		filter_params := sp_http.DefaultFilterParams()

		filters, err := sp_http.FiltersFromRequest(ctx, req, filter_params)

		if err != nil {
			logger.Error("Failed to derive filters from request", "error", err)
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		facets, err := sp_http.FacetsFromRequest(ctx, req, filter_params)

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

		facets_rsp, err := opts.Spelunker.GetDescendantsFaceted(ctx, uri.Id, filters, facets)

		if err != nil {
			logger.Error("Failed to get facets for descendants", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
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
