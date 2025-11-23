package api

import (
	"net/http"
	"regexp"

	"encoding/json"
	"github.com/aaronland/go-http/v4/sanitize"
	"github.com/aaronland/go-http/v4/slog"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-whosonfirst-derivatives"
	derivatives_http "github.com/whosonfirst/go-whosonfirst-derivatives/http"			
	wof_http "github.com/whosonfirst/go-whosonfirst/http"
)

type SelectHandlerOptions struct {
	Pattern  *regexp.Regexp
	Provider derivatives.Provider
}

func SelectHandler(opts *SelectHandlerOptions) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()
		logger := slog.LoggerWithRequest(req, nil)

		query, err := sanitize.GetString(req, "select")

		if err != nil {
			http.Error(rsp, "Invalid parameters", http.StatusBadRequest)
			return
		}

		if query == "" {
			http.Error(rsp, "Missing select", http.StatusBadRequest)
			return
		}

		if !opts.Pattern.MatchString(query) {
			http.Error(rsp, "Invalid select", http.StatusBadRequest)
			return
		}

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

		query_rsp := gjson.GetBytes(r, query)

		var rsp_body []byte

		if query_rsp.Exists() {

			enc, err := json.Marshal(query_rsp.Value())

			if err != nil {
				logger.Error("Failed to marshal response", "error", err)
				http.Error(rsp, err.Error(), http.StatusInternalServerError)
				return
			}

			rsp_body = enc
		}

		rsp.Header().Set("Content-Type", "application/json")
		rsp.Write(rsp_body)
	}

	h := http.HandlerFunc(fn)
	return h, nil
}
