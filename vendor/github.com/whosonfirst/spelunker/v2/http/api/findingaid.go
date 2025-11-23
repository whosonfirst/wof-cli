package api

import (
	"net/http"

	"github.com/aaronland/go-http/v4/slog"
	wof_http "github.com/whosonfirst/go-whosonfirst/http"
	"github.com/whosonfirst/spelunker/v2"
	sp_http "github.com/whosonfirst/spelunker/v2/http"
)

// FindingAidHandlerOptions defines options for invoking the `FindingAidHandler` method.
type FindingAidHandlerOptions struct {
	// An instance implemeting the `spelunker.Spelunker` interface.
	Spelunker spelunker.Spelunker
}

// FindingAidHandler returns an `http.Handler` for deriving the Who's On First repository for a given record.
func FindingAidHandler(opts *FindingAidHandlerOptions) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()
		logger := slog.LoggerWithRequest(req, nil)

		req_uri, err, status := wof_http.ParseURIFromRequest(req)

		if err != nil {
			logger.Error("Failed to parse URI from request", "error", err)
			http.Error(rsp, spelunker.ErrNotFound.Error(), status)
			return
		}

		spr, err := sp_http.SPRFromRequestURI(ctx, opts.Spelunker, req_uri)

		if err != nil {
			logger.Error("Failed to get by ID", "id", req_uri.Id, "error", err)
			http.Error(rsp, spelunker.ErrNotFound.Error(), http.StatusNotFound)
			return
		}

		repo := spr.Repo()

		rsp.Header().Set("Content-Type", "text/plain")
		rsp.Write([]byte(repo))
	}

	h := http.HandlerFunc(fn)
	return h, nil
}
