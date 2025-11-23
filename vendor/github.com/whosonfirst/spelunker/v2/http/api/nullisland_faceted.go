package api

import (
	"encoding/json"
	"net/http"

	// TBD
	// "github.com/aaronland/go-http/v4/auth"

	"github.com/aaronland/go-http/v4/slog"
	"github.com/whosonfirst/spelunker/v2"
	sp_http "github.com/whosonfirst/spelunker/v2/http"
)

// NullIslandFacetedHandlerOptions defines options for invoking the `NullIslandFacetedHandler` method.
type NullIslandFacetedHandlerOptions struct {
	// An instance implemeting the `spelunker.Spelunker` interface.
	Spelunker spelunker.Spelunker
	// Authenticator auth.Authenticator
}

// NullIslandFacetedHandler returns an `http.Handler` for returning faceted results for Who's On First records "visiting" Null Island.
func NullIslandFacetedHandler(opts *NullIslandFacetedHandlerOptions) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()
		logger := slog.LoggerWithRequest(req, nil)

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

		facets_rsp, err := opts.Spelunker.VisitingNullIslandFaceted(ctx, filters, facets)

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
