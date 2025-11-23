package api

import (
	"encoding/json"
	"net/http"
	"strings"

	// TBD...
	// "github.com/aaronland/go-http/v4/auth"

	"github.com/aaronland/go-http/v4/slog"
	"github.com/whosonfirst/spelunker/v2"
	sp_http "github.com/whosonfirst/spelunker/v2/http"
)

// HasConcordanceFacetedHandlerOptions defines options for invoking the `HasConcordanceFacetedHandler` method.
type HasConcordanceFacetedHandlerOptions struct {
	// An instance implemeting the `spelunker.Spelunker` interface.
	Spelunker spelunker.Spelunker
	// Authenticator auth.Authenticator
}

// HasConcordanceFacetedHandler returns an `http.Handler` for returning faceted results for a partial or complete concordance.
func HasConcordanceFacetedHandler(opts *HasConcordanceFacetedHandlerOptions) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()
		logger := slog.LoggerWithRequest(req, nil)

		ns := req.PathValue("namespace")
		pred := req.PathValue("predicate")
		value := req.PathValue("value")

		ns = strings.TrimRight(ns, ":")
		pred = strings.TrimLeft(pred, ":")
		pred = strings.TrimRight(pred, "=")

		if ns == "*" {
			ns = ""
		}

		if pred == "*" {
			pred = ""
		}

		if value == "*" {
			value = ""
		}

		// c := spelunker.NewConcordanceFromTriple(ns, pred, value)

		logger = logger.With("namespace", ns)
		logger = logger.With("predicate", pred)
		logger = logger.With("value", value)

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

		facets_rsp, err := opts.Spelunker.HasConcordanceFaceted(ctx, ns, pred, value, filters, facets)

		if err != nil {
			logger.Error("Failed to get facets for concordance", "error", err)
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
