package api

// https://preview.iiif.io/api/navplace_extension/api/extension/navplace/

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/whosonfirst/go-whosonfirst-spelunker"
	"github.com/whosonfirst/go-whosonfirst-spelunker-httpd"
)

type NavPlaceHandlerOptions struct {
	Spelunker   spelunker.Spelunker
	MaxFeatures int
}

// NavPlaceHandler will return a given record as a FeatureCollection for use by the IIIF navPlace extension,
// specifically as navPlace "reference" objects.
func NavPlaceHandler(opts *NavPlaceHandlerOptions) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()
		logger := httpd.LoggerWithRequest(req, nil)

		q := req.URL.Query()
		base := q.Get("id")

		if base == "" {
			path := req.URL.Path
			base = filepath.Base(path)

			base = strings.TrimLeft(base, "/")
			base = strings.TrimRight(base, "/")
		}

		ids := strings.Split(base, ",")

		uris := make([]*httpd.URI, len(ids))

		for idx, str_id := range ids {

			req_uri, err, status := httpd.ParseURIFromPath(ctx, str_id, nil)

			if err != nil {
				logger.Error("Failed to parse URI from request", "id", str_id, "error", err)
				http.Error(rsp, err.Error(), status)
				return
			}

			uris[idx] = req_uri
		}

		count := len(uris)

		if count == 0 {
			http.Error(rsp, "No IDs to include", http.StatusBadRequest)
			return
		}

		if count > opts.MaxFeatures {
			logger.Error("Exceed maximum number of features")
			http.Error(rsp, "Maximum number of IDs exceeded", http.StatusBadRequest)
			return
		}

		rsp.Header().Set("Content-Type", "application/geo+json")

		rsp.Write([]byte(`{"type":"FeatureCollection", "features":[`))

		for i, req_uri := range uris {

			r, err := httpd.FeatureFromRequestURI(ctx, opts.Spelunker, req_uri)

			if err != nil {
				logger.Error("Failed to retrieve record", "id", req_uri.Id, "error", err)
				http.Error(rsp, "Failed to retrieve ID", http.StatusInternalServerError)
				return
			}

			rsp.Write(r)

			if i+1 < count {
				rsp.Write([]byte(`,`))
			}
		}

		rsp.Write([]byte(`]}`))
	}

	h := http.HandlerFunc(fn)
	return h, nil
}
