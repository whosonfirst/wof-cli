package www

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/aaronland/go-pagination"
	"github.com/sfomuseum/go-http-auth"
	"github.com/whosonfirst/go-whosonfirst-sources"
	"github.com/whosonfirst/go-whosonfirst-spelunker"
	"github.com/whosonfirst/go-whosonfirst-spelunker-httpd"
	"github.com/whosonfirst/go-whosonfirst-spr/v2"
)

type HasConcordanceHandlerOptions struct {
	Spelunker     spelunker.Spelunker
	Authenticator auth.Authenticator
	Templates     *template.Template
	URIs          *httpd.URIs
}

type HasConcordanceHandlerVars struct {
	PageTitle        string
	URIs             *httpd.URIs
	Concordance      *spelunker.Concordance
	Places           []spr.StandardPlacesResult
	Pagination       pagination.Results
	PaginationURL    string
	FacetsURL        string
	FacetsContextURL string
	Source           *sources.WOFSource
	OpenGraph        *OpenGraph
}

func HasConcordanceHandler(opts *HasConcordanceHandlerOptions) (http.Handler, error) {

	t := opts.Templates.Lookup("concordance")

	if t == nil {
		return nil, fmt.Errorf("Failed to locate 'concordance' template")
	}

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()
		logger := httpd.LoggerWithRequest(req, nil)

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

		pg_opts, err := httpd.PaginationOptionsFromRequest(req)

		if err != nil {
			logger.Error("Failed to create pagination options", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		filter_params := httpd.DefaultFilterParams()

		filters, err := httpd.FiltersFromRequest(ctx, req, filter_params)

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
			pagination_url = httpd.URIForConcordanceTriple(opts.URIs.ConcordanceTriple, ns, pred, value, filters, nil)
			facets_url = httpd.URIForConcordanceTriple(opts.URIs.ConcordanceTripleFaceted, ns, pred, value, filters, nil)
		case ns != "" && pred != "":
			pagination_url = httpd.URIForConcordanceNSPred(opts.URIs.ConcordanceNSPred, ns, pred, filters, nil)
			facets_url = httpd.URIForConcordanceNSPred(opts.URIs.ConcordanceNSPredFaceted, ns, pred, filters, nil)
		case pred != "" && value != "":
			pagination_url = httpd.URIForConcordanceTriple(opts.URIs.ConcordanceTriple, "*", pred, value, filters, nil)
			facets_url = httpd.URIForConcordanceTriple(opts.URIs.ConcordanceTripleFaceted, "*", pred, value, filters, nil)
		case ns != "" && value != "":
			pagination_url = httpd.URIForConcordanceTriple(opts.URIs.ConcordanceTriple, ns, "*", value, filters, nil)
			facets_url = httpd.URIForConcordanceTriple(opts.URIs.ConcordanceTripleFaceted, ns, "*", value, filters, nil)
		case ns != "":
			pagination_url = httpd.URIForConcordanceNS(opts.URIs.ConcordanceNS, ns, filters, nil)
			facets_url = httpd.URIForConcordanceNS(opts.URIs.ConcordanceNSFaceted, ns, filters, nil)
		case pred != "":
			pagination_url = httpd.URIForConcordanceTriple(opts.URIs.ConcordanceTriple, "*", pred, "*", filters, nil)
			facets_url = httpd.URIForConcordanceTriple(opts.URIs.ConcordanceTripleFaceted, "*", pred, "*", filters, nil)
		case value != "":
			pagination_url = httpd.URIForConcordanceTriple(opts.URIs.ConcordanceTriple, "*", "*", value, filters, nil)
			facets_url = httpd.URIForConcordanceTriple(opts.URIs.ConcordanceTripleFaceted, "*", "*", value, filters, nil)
		default:
			logger.Info("WUT")
		}

		facets_context_url = req.URL.Path

		vars := HasConcordanceHandlerVars{
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

		vars.OpenGraph = &OpenGraph{
			Type:        "Article",
			SiteName:    "Who's On First Spelunker",
			Title:       fmt.Sprintf(`Who's On First concordances for \"%s\"`, c),
			Description: fmt.Sprintf(`Who's On First records that "hold hands" with records from %s`, src.Fullname),
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
