package api

import (
	"encoding/json"
	"net/http"

	"github.com/aaronland/go-http/v4/slog"
	"github.com/whosonfirst/go-whosonfirst-derivatives"
	derivatives_http "github.com/whosonfirst/go-whosonfirst-derivatives/http"			
	wof_http "github.com/whosonfirst/go-whosonfirst/http"
	"github.com/whosonfirst/go-whosonfirst-spr/v2"
)

type SPRHandlerOptions struct {
	Provider derivatives.Provider
}

func SPRHandler(opts *SPRHandlerOptions) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()
		logger := slog.LoggerWithRequest(req, nil)

		req_uri, err, status := wof_http.ParseURIFromRequest(req)

		if err != nil {
			logger.Error("Failed to parse URI from request", "error", err)
			http.Error(rsp, derivatives.ErrNotFound.Error(), status)
			return
		}

		if req_uri.Id <= -1 {
			http.Error(rsp, "Not found", http.StatusNotFound)
			return
		}

		logger = logger.With("id", req_uri.Id)

		r, err := derivatives_http.FeatureFromRequestURI(ctx, opts.Provider, req_uri)

		if err != nil {
			logger.Error("Failed to get by ID", "error", err)
			http.Error(rsp, derivatives.ErrNotFound.Error(), http.StatusNotFound)
			return
		}

		s, err := spr.WhosOnFirstSPR(r)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		rsp.Header().Set("Content-Type", "application/json")

		enc := json.NewEncoder(rsp)
		err = enc.Encode(s)

		if err != nil {
			logger.Error("Failed to marshal response", "error", err)
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		return
	}

	h := http.HandlerFunc(fn)
	return h, nil
}
