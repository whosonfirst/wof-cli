package api

import (
	"net/http"

	"github.com/sfomuseum/go-geojsonld"
	"github.com/whosonfirst/go-whosonfirst-spelunker"
	"github.com/whosonfirst/go-whosonfirst-spelunker-httpd"
)

type GeoJSONLDHandlerOptions struct {
	Spelunker spelunker.Spelunker
}

func GeoJSONLDHandler(opts *GeoJSONLDHandlerOptions) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()
		logger := httpd.LoggerWithRequest(req, nil)

		req_uri, err, status := httpd.ParseURIFromRequest(req, nil)

		if err != nil {
			logger.Error("Failed to parse URI from request", "error", err)
			http.Error(rsp, spelunker.ErrNotFound.Error(), status)
			return
		}

		if req_uri.Id <= -1 {
			http.Error(rsp, "Not found", http.StatusNotFound)
			return
		}

		logger = logger.With("id", req_uri.Id)

		r, err := httpd.FeatureFromRequestURI(ctx, opts.Spelunker, req_uri)

		if err != nil {
			logger.Error("Failed to get by ID", "error", err)
			http.Error(rsp, spelunker.ErrNotFound.Error(), http.StatusNotFound)
			return
		}

		body, err := geojsonld.AsGeoJSONLD(ctx, r)

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
