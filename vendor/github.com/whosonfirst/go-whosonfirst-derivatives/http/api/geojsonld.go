package api

import (
	"net/http"
	
	"github.com/aaronland/go-http/v4/slog"
	"github.com/sfomuseum/go-geojsonld"
	"github.com/whosonfirst/go-whosonfirst-derivatives"
	wof_http "github.com/whosonfirst/go-whosonfirst/http"
)

type GeoJSONLDHandlerOptions struct {
	Provider derivatives.Provider
}

func GeoJSONLDHandler(opts *GeoJSONLDHandlerOptions) (http.Handler, error) {

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

		r, err := opts.Provider.GetFeature(ctx, req_uri.Id, req_uri.URIArgs)

		if err != nil {
			logger.Error("Failed to get by ID", "error", err)
			http.Error(rsp, derivatives.ErrNotFound.Error(), http.StatusNotFound)
			return
		}

		body, err := geojsonld.AsGeoJSONLDWithReader(ctx, r)

		if err != nil {
			logger.Error("Failed to render geojson", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		rsp.Header().Set("Content-Type", "application/geo+json")
		rsp.Write([]byte(body))
	}

	h := http.HandlerFunc(fn)
	return h, nil
}
