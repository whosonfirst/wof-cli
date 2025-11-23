package www

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/aaronland/go-http/v4/auth"
	"github.com/aaronland/go-http/v4/slog"
	"github.com/aaronland/go-pagination"
	"github.com/whosonfirst/go-whosonfirst-sources"
	"github.com/whosonfirst/go-whosonfirst-spr/v2"
	"github.com/whosonfirst/spelunker/v2"
	wof_http "github.com/whosonfirst/spelunker/v2/http"
)

type hasConcordanceHandlerVars struct {
	PageTitle        string
	URIs             *wof_http.URIs
	Concordance      *spelunker.Concordance
	Places           []spr.StandardPlacesResult
	Pagination       pagination.Results
	PaginationURL    string
	FacetsURL        string
	FacetsContextURL string
	Source           *sources.WOFSource
	OpenGraph        *OpenGraph
}

// HasConcordancesHandlerOptions defines configuration options for the `HasConcordanceHandler` method.
type HasConcordanceHandlerOptions struct {
	// An instance implemeting the `spelunker.Spelunker` interface.
	Spelunker spelunker.Spelunker
	// An instance implementing the `aaronland/go-http/v4/auth.Authenticator` interface.
	Authenticator auth.Authenticator
	// An `html/template.Template` instance containing the named template "concordance".
	Templates *template.Template
	// URIs are the `wof_http.URIs` details for this Spelunker instance.
	URIs *wof_http.URIs
}

// HasConcordanceHandler returns an `http.Handler` instance to display a webpage for a specific concordance.
func HasConcordanceHandler(opts *HasConcordanceHandlerOptions) (http.Handler, error) {

	t := opts.Templates.Lookup("concordance")

	if t == nil {
		return nil, fmt.Errorf("Failed to locate 'concordance' template")
	}

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()
		logger := slog.LoggerWithRequest(req, nil)

		ns := req.PathValue("namespace")
		pred := req.PathValue("predicate")
		value := req.PathValue("value")

		ns = strings.TrimRight(ns, ":")
		pred = strings.TrimLeft(pred, ":")
		pred = strings.TrimRight(pred, "=")

		logger = logger.With("namespace", ns)
		logger = logger.With("predicate", pred)
		logger = logger.With("value", value)

		if ns == "*" {
			ns = ""
		}

		if pred == "*" {
			pred = ""
		}

		if value == "*" {
			value = ""
		}

		c := spelunker.NewConcordanceFromTriple(ns, pred, value)

		pg_opts, err := wof_http.PaginationOptionsFromRequest(req)

		if err != nil {
			logger.Error("Failed to create pagination options", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		filter_params := wof_http.DefaultFilterParams()

		filters, err := wof_http.FiltersFromRequest(ctx, req, filter_params)

		if err != nil {
			logger.Error("Failed to derive filters from request", "error", err)
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		r, pg_r, err := opts.Spelunker.HasConcordance(ctx, pg_opts, ns, pred, value, filters)

		if err != nil {
			logger.Error("Failed to get records having concordance", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		page_title := fmt.Sprintf("Concordances for %s", c)

		var src *sources.WOFSource

		if ns != "" {

			v, err := sources.GetSourceByPrefix(ns)

			if err != nil {
				logger.Warn("Failed to derive source from namespace", "error", err)
			} else {
				src = v
			}
		}

		var pagination_url string
		var facets_url string
		var facets_context_url string

		switch {
		case ns != "" && pred != "" && value != "":
			pagination_url = wof_http.URIForConcordanceTriple(opts.URIs.ConcordanceTriple, ns, pred, value, filters, nil)
			facets_url = wof_http.URIForConcordanceTriple(opts.URIs.ConcordanceTripleFaceted, ns, pred, value, filters, nil)
		case ns != "" && pred != "":
			pagination_url = wof_http.URIForConcordanceNSPred(opts.URIs.ConcordanceNSPred, ns, pred, filters, nil)
			facets_url = wof_http.URIForConcordanceNSPred(opts.URIs.ConcordanceNSPredFaceted, ns, pred, filters, nil)
		case pred != "" && value != "":
			pagination_url = wof_http.URIForConcordanceTriple(opts.URIs.ConcordanceTriple, "*", pred, value, filters, nil)
			facets_url = wof_http.URIForConcordanceTriple(opts.URIs.ConcordanceTripleFaceted, "*", pred, value, filters, nil)
		case ns != "" && value != "":
			pagination_url = wof_http.URIForConcordanceTriple(opts.URIs.ConcordanceTriple, ns, "*", value, filters, nil)
			facets_url = wof_http.URIForConcordanceTriple(opts.URIs.ConcordanceTripleFaceted, ns, "*", value, filters, nil)
		case ns != "":
			pagination_url = wof_http.URIForConcordanceNS(opts.URIs.ConcordanceNS, ns, filters, nil)
			facets_url = wof_http.URIForConcordanceNS(opts.URIs.ConcordanceNSFaceted, ns, filters, nil)
		case pred != "":
			pagination_url = wof_http.URIForConcordanceTriple(opts.URIs.ConcordanceTriple, "*", pred, "*", filters, nil)
			facets_url = wof_http.URIForConcordanceTriple(opts.URIs.ConcordanceTripleFaceted, "*", pred, "*", filters, nil)
		case value != "":
			pagination_url = wof_http.URIForConcordanceTriple(opts.URIs.ConcordanceTriple, "*", "*", value, filters, nil)
			facets_url = wof_http.URIForConcordanceTriple(opts.URIs.ConcordanceTripleFaceted, "*", "*", value, filters, nil)
		default:
			logger.Info("WUT")
		}

		facets_context_url = req.URL.Path

		vars := hasConcordanceHandlerVars{
			PageTitle:        page_title,
			URIs:             opts.URIs,
			Concordance:      c,
			Places:           r.Results(),
			Pagination:       pg_r,
			Source:           src,
			PaginationURL:    pagination_url,
			FacetsURL:        facets_url,
			FacetsContextURL: facets_context_url,
		}

		var desc string

		if src != nil {
			desc = fmt.Sprintf(`Who's On First records that "hold hands" with records from %s`, src.Fullname)
		} else {
			desc = fmt.Sprintf(`Who's On First records that "hold hands" with records matching this concordance`)
		}

		vars.OpenGraph = &OpenGraph{
			Type:        "Article",
			SiteName:    "Who's On First Spelunker",
			Title:       fmt.Sprintf(`Who's On First concordances for \"%s\"`, c),
			Description: desc,
			Image:       "",
		}

		rsp.Header().Set("Content-Type", "text/html")

		err = t.Execute(rsp, vars)

		if err != nil {
			logger.Error("Failed to render template", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
		}

	}

	h := http.HandlerFunc(fn)
	return h, nil
}
