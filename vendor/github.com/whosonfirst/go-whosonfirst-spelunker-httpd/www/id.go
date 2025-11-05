package www

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/sfomuseum/go-http-auth"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-whosonfirst-placetypes"
	"github.com/whosonfirst/go-whosonfirst-spelunker"
	"github.com/whosonfirst/go-whosonfirst-spelunker-httpd"
	wof_funcs "github.com/whosonfirst/go-whosonfirst-spelunker-httpd/templates/funcs"
	"github.com/whosonfirst/go-whosonfirst-uri"
)

type IdHandlerOptions struct {
	Spelunker     spelunker.Spelunker
	Authenticator auth.Authenticator
	Templates     *template.Template
	URIs          *httpd.URIs
}

type IdHandlerAncestor struct {
	Placetype string
	Id        int64
}

type IdHandlerVars struct {
	Id               int64
	RequestId        string
	URIArgs          *uri.URIArgs
	PageTitle        string
	URIs             *httpd.URIs
	Properties       string
	CountDescendants int64
	Hierarchies      [][]*IdHandlerAncestor
	RelPath          string
	GitHubURL        string
	WriteFieldURL    string
	OpenGraph        *OpenGraph
}

func IdHandler(opts *IdHandlerOptions) (http.Handler, error) {

	t := opts.Templates.Lookup("id")

	if t == nil {
		return nil, fmt.Errorf("Failed to locate 'id' template")
	}

	alt_t := opts.Templates.Lookup("alt")

	if alt_t == nil {
		return nil, fmt.Errorf("Missing alt template")
	}

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()
		logger := httpd.LoggerWithRequest(req, nil)

		req_uri, err, status := httpd.ParseURIFromRequest(req, nil)

		if err != nil {
			logger.Error("Failed to parse URI from request", "error", err)
			http.Error(rsp, spelunker.ErrNotFound.Error(), status)
			return
		}

		wof_id := req_uri.Id

		req_id, err := uri.Id2Fname(req_uri.Id, req_uri.URIArgs)

		if err != nil {
			logger.Error("Failed to derive request ID", "error", err)
			http.Error(rsp, spelunker.ErrNotFound.Error(), http.StatusNotFound)
			return
		}

		req_id = strings.Replace(req_id, filepath.Ext(req_id), "", 1)

		logger = logger.With("wof id", wof_id)

		uri_args := new(uri.URIArgs)

		f, err := opts.Spelunker.GetRecordForId(ctx, wof_id, uri_args)

		if err != nil {
			logger.Error("Failed to get by ID", "error", err)
			http.Error(rsp, spelunker.ErrNotFound.Error(), http.StatusNotFound)
			return
		}

		name_rsp := gjson.GetBytes(f, "wof:name")
		wof_name := name_rsp.String()

		country_rsp := gjson.GetBytes(f, "wof:country")
		wof_country := country_rsp.String()

		country_name, country_exists := httpd.CountryCodeLookup[wof_country]

		if !country_exists {
			country_name = wof_country
		}

		rel_path, err := uri.Id2RelPath(wof_id, req_uri.URIArgs)

		if err != nil {
			logger.Error("Failed to derive relative path for record", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		repo_name := gjson.GetBytes(f, "wof:repo")

		github_url := fmt.Sprintf("https://github.com/whosonfirst-data/%s/blob/master/data/%s", repo_name, rel_path)

		vars := IdHandlerVars{
			Id:         wof_id,
			RequestId:  req_id,
			URIArgs:    req_uri.URIArgs,
			Properties: string(f),
			PageTitle:  wof_name,
			GitHubURL:  github_url,
			URIs:       opts.URIs,
			RelPath:    rel_path,
		}

		if req_uri.IsAlternate {

			rsp.Header().Set("Content-Type", "text/html")

			err = alt_t.Execute(rsp, vars)

			if err != nil {
				logger.Error("Failed to return ", "error", err)
				http.Error(rsp, "womp womp", http.StatusInternalServerError)
			}

			return
		}

		count_descendants, err := opts.Spelunker.CountDescendants(ctx, wof_id)

		if err != nil {
			logger.Error("Failed to count descendants", "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		// START OF there's got to be a better way to do this...

		str_pt := gjson.GetBytes(f, "wof:placetype")

		pt, err := placetypes.GetPlacetypeByName(str_pt.String())

		if err != nil {
			logger.Error("Failed to load placetype", "placetype", str_pt, "error", err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		roles := []string{
			"common",
			"optional",
			"common_optional",
		}

		// If custom placetype assume the most granular placetype for constructing
		// ordered list from hierarchy.

		if pt.Name == "custom" {

			v, err := placetypes.GetPlacetypeByName("installation")

			if err != nil {
				logger.Warn("Failed to instantiate 'installation' placetype, %w", err)
			} else {
				pt = v
			}
		}

		ancestors := placetypes.AncestorsForRoles(pt, roles)
		count_ancestors := len(ancestors)

		sorted := make([]string, 0)

		for i := count_ancestors - 1; i >= 0; i-- {
			n := ancestors[i]
			sorted = append(sorted, n.String())
		}

		hierarchies := make([]map[string]int64, 0)

		h_rsp := gjson.GetBytes(f, "wof:hierarchy")

		if h_rsp.Exists() {

			for _, h := range h_rsp.Array() {

				dict := make(map[string]int64)

				for k, v := range h.Map() {
					dict[k] = v.Int()
				}

				hierarchies = append(hierarchies, dict)
			}
		}

		handler_hierarchies := make([][]*IdHandlerAncestor, len(hierarchies))

		for idx, hier := range hierarchies {

			handler_ancestors := make([]*IdHandlerAncestor, 0)

			for _, n := range sorted {

				k := fmt.Sprintf("%s_id", n)
				v, ok := hier[k]

				if !ok {
					continue
				}

				a := &IdHandlerAncestor{
					Placetype: n,
					Id:        v,
				}

				handler_ancestors = append(handler_ancestors, a)
			}

			handler_hierarchies[idx] = handler_ancestors
		}

		// END OF there's got to be a better way to do this...

		writefield_url := fmt.Sprintf("https://raw.githubusercontent.com/whosonfirst-data/%s/master/data/%s", repo_name, rel_path)

		vars.CountDescendants = count_descendants
		vars.Hierarchies = handler_hierarchies
		vars.WriteFieldURL = writefield_url

		// START OF put me in a function or something...

		is_pt := wof_funcs.IsAPlacetype(str_pt.String())

		var og_desc string

		switch str_pt.String() {
		case "continent", "planet", "empire", "ocean":
			og_desc = fmt.Sprintf("%s (%d) is %s", wof_name, wof_id, is_pt)
		case "country":
			og_desc = fmt.Sprintf("%s (%d) is %s :flag-%s:", wof_name, wof_id, is_pt, strings.ToLower(wof_country))
		default:

			var og_country string

			switch wof_country {
			case "US":
				og_country = fmt.Sprintf("the %s", country_name)
			default:
				og_country = country_name
			}

			og_desc = fmt.Sprintf("%s (%d) is %s in %s :flag-%s:", wof_name, wof_id, is_pt, og_country, strings.ToLower(wof_country))
		}

		// END OF put me in a function or something...

		svg_url := httpd.URIForIdSimple(opts.URIs.SVG, wof_id)

		og_image, err := opts.URIs.Abs(svg_url)

		if err != nil {
			logger.Error("Failed to derive absolute URL for SVG image", "url", svg_url, "error", err)
		}

		vars.OpenGraph = &OpenGraph{
			Type:        "Article",
			SiteName:    "Who's On First Spelunker",
			Title:       wof_name,
			Description: og_desc,
			Image:       og_image,
		}

		rsp.Header().Set("Content-Type", "text/html")

		err = t.Execute(rsp, vars)

		if err != nil {
			logger.Error("Failed to return ", "error", err)
			http.Error(rsp, "womp womp", http.StatusInternalServerError)
		}

	}

	h := http.HandlerFunc(fn)
	return h, nil
}
