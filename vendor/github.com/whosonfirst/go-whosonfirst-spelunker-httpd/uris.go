package httpd

import (
	"fmt"
	"log/slog"
	"net/url"
	"reflect"
	"strings"

	"github.com/whosonfirst/go-whosonfirst-spelunker"
)

type URIs struct {
	// WWW/human-readable
	Id                string   `json:"id"`
	IdAlt             []string `json:"id_alt"`
	Concordances      string   `json:"concordances"`
	ConcordanceNS     string   `json:"concordance_ns"`
	ConcordanceNSPred string   `json:"concordance_ns_pred"`
	ConcordanceTriple string   `json:"concordance_triple"`
	Descendants       string   `json:"descendants"`
	DescendantsAlt    []string `json:"descendants_alt"`
	Index             string   `json:"index"`
	Placetypes        string   `json:"placetypes"`
	Placetype         string   `json:"placetype"`
	NullIsland        string   `json:"nullisland"`
	Recent            string   `json:"recent"`
	RecentAlt         []string `json:"recent_alt"`
	Search            string   `json:"search"`
	About             string   `json:"about"`
	Code              string   `json:"code"`
	HowTo             string   `json:"how_to"`
	Tiles             string   `json:"tiles"`
	OpenSearch        string   `json:"opensearch"`

	// Static assets
	Static string `json:"static"`

	// API/machine-readable
	ConcordanceNSFaceted     string   `json:"concordance_ns"`
	ConcordanceNSPredFaceted string   `json:"concordance_ns_pred"`
	ConcordanceTripleFaceted string   `json:"concordance_triple_faceted"`
	DescendantsFaceted       string   `json:"descendants_faceted"`
	FindingAid               string   `json:"finding_aid"`
	GeoJSON                  string   `json:"geojson"`
	GeoJSONAlt               []string `json:"geojson_alt"`
	GeoJSONLD                string   `json:"geojsonld"`
	GeoJSONLDAlt             []string `json:"geojsonld_alt"`
	NavPlace                 string   `json:"navplace"`
	NavPlaceAlt              []string `json:"navplace_alt"`
	NullIslandFaceted        string   `json:"nullisland_faceted"`
	PlacetypeFaceted         string   `json:"placetype_faceted"`
	RecentFaceted            string   `json:"recent_faceted"`
	SearchFaceted            string   `json:"search_faceted"`
	Select                   string   `json:"select"`
	SelectAlt                []string `json:"select_alt"`
	SPR                      string   `json:"spr"`
	SPRAlt                   []string `json:"spr_alt"`
	SVG                      string   `json:"svg"`
	SVGAlt                   []string `json:"svg_alt"`

	RootURL string `json:"root_url"`
}

func (u *URIs) ApplyPrefix(prefix string) error {

	val := reflect.ValueOf(*u)

	for i := 0; i < val.NumField(); i++ {

		field := val.Field(i)
		v := field.String()

		if v == "" {
			continue
		}

		if strings.HasPrefix(v, prefix) {
			continue
		}

		new_v, err := url.JoinPath(prefix, v)

		if err != nil {
			return fmt.Errorf("Failed to assign prefix to %s, %w", v, err)
		}

		reflect.ValueOf(u).Elem().Field(i).SetString(new_v)
	}

	return nil
}

func DefaultURIs() *URIs {

	// Note that the default path for ID-related URIs is "/id/{id}/foo"
	// mostly so that the URIForId template function will work. More generic
	// catch-all paths are stored in {NAME}Alt URI definitions. For example:
	// GeoJSON: "/id/{id}/geojson" handles: "http://localhost:8080/id/1327010993/geojson"
	// GeoJSONAlt: []string{ "/geojson", } handles: "http://localhost:8080/geojson/132/701/099/3/1327010993.geojson"

	uris_table := &URIs{

		// WWW/human-readable

		Index:             "/",
		Search:            "/search",
		About:             "/about",
		Code:              "/code",
		HowTo:             "/howto",
		NullIsland:        "/nullisland",
		Placetypes:        "/placetypes",
		Placetype:         "/placetypes/{placetype}",
		Concordances:      "/concordances",
		ConcordanceNS:     "/concordances/{namespace}",
		ConcordanceNSPred: "/concordances/{namespace}:{predicate}",
		ConcordanceTriple: "/concordances/{namespace}:{predicate}={value}",
		Recent:            "/recent/{duration}",
		RecentAlt: []string{
			"/recent",
		},
		Id:          "/id/{id}",
		Descendants: "/id/{id}/descendants",
		Tiles:       "/tiles/{z}/{x}/{y}",
		OpenSearch:  "/opensearch",

		// Static Assets
		Static: "/static/",

		// API/machine-readable
		ConcordanceNSFaceted:     "/concordances/{namespace}/facets",
		ConcordanceNSPredFaceted: "/concordances/{namespace}:{predicate}/facets",
		ConcordanceTripleFaceted: "/concordances/{namespace}:{predicate}={value}/facets",
		DescendantsFaceted:       "/id/{id}/descendants/facets",

		FindingAid: "/findingaid/",

		GeoJSON: "/id/{id}/geojson",
		GeoJSONAlt: []string{
			"/geojson/",
		},
		GeoJSONLD: "/id/{id}/geojsonld",
		GeoJSONLDAlt: []string{
			"/geojsonld/",
		},
		NavPlace: "/id/{id}/navplace",
		NavPlaceAlt: []string{
			"/navplace/",
		},
		NullIslandFaceted: "/nullisland/facets",
		PlacetypeFaceted:  "/placetypes/{placetype}/facets",
		RecentFaceted:     "/recent/{duration}/facets",
		SearchFaceted:     "/search/facets",
		Select:            "/id/{id}/select",
		SelectAlt: []string{
			"/select/",
		},
		SPR: "/id/{id}/spr",
		SPRAlt: []string{
			"/spr/",
		},
		SVG: "/id/{id}/svg",
		SVGAlt: []string{
			"/svg/",
		},
	}

	return uris_table
}

func (uris_table *URIs) Abs(path string) (string, error) {

	root_u, err := url.Parse(uris_table.RootURL)

	if err != nil {
		return "", fmt.Errorf("Failed to parse root URL, %w", err)
	}

	this_u := url.URL{}
	this_u.Host = root_u.Host
	this_u.Scheme = root_u.Scheme
	this_u.Path = path

	return this_u.String(), nil
}

func URIForIdSimple(uri string, id int64) string {
	id_uri := ReplaceAll(uri, "{id}", id)
	return uriWithFilters(id_uri, nil, nil)
}

func URIForId(uri string, id int64, filters []spelunker.Filter, facets []spelunker.Facet) string {

	id_uri := ReplaceAll(uri, "{id}", id)
	return uriWithFilters(id_uri, filters, facets)
}

func URIForPlacetype(uri string, pt string, filters []spelunker.Filter, facets []spelunker.Facet) string {

	pt_uri := ReplaceAll(uri, "{placetype}", pt)
	return uriWithFilters(pt_uri, filters, facets)
}

func URIForRecentSimple(uri string, d string) string {
	r_uri := ReplaceAll(uri, "{duration}", d)
	return uriWithFilters(r_uri, nil, nil)
}

func URIForRecent(uri string, d string, filters []spelunker.Filter, facets []spelunker.Facet) string {

	r_uri := ReplaceAll(uri, "{duration}", d)
	return uriWithFilters(r_uri, filters, facets)
}

func URIForConcordanceNS(uri string, ns string, filters []spelunker.Filter, facets []spelunker.Facet) string {

	c_uri := ReplaceAll(uri, "{namespace}", ns)
	return uriWithFilters(c_uri, filters, facets)
}

func URIForConcordanceNSPred(uri string, ns string, pred string, filters []spelunker.Filter, facets []spelunker.Facet) string {

	c_uri := uri

	c_uri = ReplaceAll(c_uri, "{namespace}", ns)
	c_uri = ReplaceAll(c_uri, "{predicate}", pred)
	return uriWithFilters(c_uri, filters, facets)
}

func URIForConcordanceTriple(uri string, ns string, pred string, value any, filters []spelunker.Filter, facets []spelunker.Facet) string {

	c_uri := uri

	c_uri = ReplaceAll(c_uri, "{namespace}", ns)
	c_uri = ReplaceAll(c_uri, "{predicate}", pred)
	c_uri = ReplaceAll(c_uri, "{value}", value)
	return uriWithFilters(c_uri, filters, facets)
}

func URIForSearch(uri string, query string, filters []spelunker.Filter, facets []spelunker.Facet) string {

	u, _ := url.Parse(uri)
	q := u.Query()

	q.Set("q", query)
	u.RawQuery = q.Encode()

	return uriWithFilters(u.String(), filters, facets)
}

func URIForNullIsland(uri string, filters []spelunker.Filter, facets []spelunker.Facet) string {

	return uriWithFilters(uri, filters, facets)
}

func uriWithFilters(uri string, filters []spelunker.Filter, facets []spelunker.Facet) string {

	u, _ := url.Parse(uri)
	q := u.Query()

	for _, f := range filters {
		q.Set(f.Scheme(), fmt.Sprintf("%v", f.Value()))
	}

	for _, f := range facets {
		q.Set("facet", f.String())
	}

	u.RawQuery = q.Encode()

	slog.Debug("URI", "with filters and facets", u.String())
	return u.String()
}

func ReplaceAll(input string, pattern string, value any) string {
	str_value := fmt.Sprintf("%v", value)
	return strings.Replace(input, pattern, str_value, -1)
}
