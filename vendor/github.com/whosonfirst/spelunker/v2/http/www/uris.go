package www

import (
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"

	"github.com/aaronland/go-http/v4/slog"
	wof_http "github.com/whosonfirst/spelunker/v2/http"
)

// URIsJSHandlerOptions defines configuration options for invoking the `URIsJSHandler` method.
type URIsJSHandlerOptions struct {
	// Templates are a `text/template` instance containing the "whosonfirst_spelunker_uris" template.
	Templates *template.Template
	// URIs are the `wof_http.URIs` details for this Spelunker instance.
	URIs *wof_http.URIs
}

type urisJSVars struct {
	Table string
}

// URIsJSHandler returns JSON-encoded information about Spelunker-specific URIs.
func URIsJSHandler(opts *URIsJSHandlerOptions) (http.Handler, error) {

	t := opts.Templates.Lookup("whosonfirst_spelunker_uris")

	if t == nil {
		return nil, fmt.Errorf("Failed to locate 'whosonfirst_spelunker_uris' template")
	}

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		logger := slog.LoggerWithRequest(req, nil)

		enc_table, err := json.Marshal(opts.URIs)

		if err != nil {
			logger.Error("Failed to marshal URIs table", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		vars := urisJSVars{
			Table: string(enc_table),
		}

		rsp.Header().Set("Content-type", "text/javascript")
		err = t.Execute(rsp, vars)

		if err != nil {
			logger.Error("Failed to execute template", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		return
	}

	return http.HandlerFunc(fn), nil
}
