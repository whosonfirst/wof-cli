package api

import (
	"encoding/json"
	"net/http"

	"github.com/aaronland/go-http/v4/slog"
	"github.com/paulmach/orb/encoding/wkt"
	"github.com/paulmach/orb/geojson"
	"github.com/whosonfirst/go-whosonfirst-derivatives"
	derivatives_http "github.com/whosonfirst/go-whosonfirst-derivatives/http"			
	wof_http "github.com/whosonfirst/go-whosonfirst/http"
)

type WKTHandlerOptions struct {
	Provider derivatives.Provider
}

func WKTHandler(opts *WKTHandlerOptions) (http.Handler, error) {

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

		var f geojson.Feature
		err = json.Unmarshal(r, &f)

		if err != nil {
			logger.Error("Failed to unmarshal feature", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		geom := f.Geometry
		wkt := wkt.Marshal(geom)

		rsp.Header().Set("Content-Type", "text/plain")

		rsp.Write([]byte(wkt))
		return
	}

	h := http.HandlerFunc(fn)
	return h, nil
}
