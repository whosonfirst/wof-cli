package api

import (
	"net/http"

	"github.com/aaronland/go-http/v4/slog"
	"github.com/whosonfirst/go-whosonfirst-derivatives"
	derivatives_http "github.com/whosonfirst/go-whosonfirst-derivatives/http"	
	wof_http "github.com/whosonfirst/go-whosonfirst/http"
)

type GeoJSONHandlerOptions struct {
	Provider derivatives.Provider
}

func GeoJSONHandler(opts *GeoJSONHandlerOptions) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()
		logger := slog.LoggerWithRequest(req, nil)

		req_uri, err, status := wof_http.ParseURIFromRequest(req)

		if err != nil {
			logger.Error("Failed to parse URI from request", "error", err)
			http.Error(rsp, derivatives.ErrNotFound.Error(), status)
			return
		}

		wof_id := req_uri.Id

		if wof_id <= -1 {
			http.Error(rsp, "Not found", http.StatusNotFound)
			return
		}

		logger = logger.With("wof id", wof_id)

		r, err := derivatives_http.FeatureFromRequestURI(ctx, opts.Provider, req_uri)

		if err != nil {
			logger.Error("Failed to get by ID", "error", err)
			http.Error(rsp, derivatives.ErrNotFound.Error(), http.StatusNotFound)
			return
		}

		rsp.Header().Set("Content-Type", "application/json")
		rsp.Write(r)
	}

	h := http.HandlerFunc(fn)
	return h, nil
}
