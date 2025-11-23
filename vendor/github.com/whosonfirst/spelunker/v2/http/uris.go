package http

import (
	"fmt"
	"log/slog"
	"net/url"
	"reflect"
	"strings"

	"github.com/whosonfirst/spelunker/v2"
)

// URIs is a struct for defining paths (URIs) for specific Spelunker web application endpoints.
type URIs struct {
	// Id defines the URI for individual records.
	Id string `json:"id"`
	// Id defines the URI for individual alternate geometry records.
	IdAlt []string `json:"id_alt"`
	// Concordances defines the URI for all the concordances (across all records).
	Concordances string `json:"concordances"`
	// ConcordancesNS defines the URI for all the concordances with a given namespace.
	ConcordanceNS string `json:"concordance_ns"`
	// ConcordancesNSPred defines the URI for all the concordances with a given namespace and predicate pair.
	ConcordanceNSPred string `json:"concordance_ns_pred"`
	// ConcordancesTriple defines the URI for all the concordances with a fully-qualified concordance (ns:pred=value).
	ConcordanceTriple string `json:"concordance_triple"`
	// Descendants defines the URI for all the descendants of a given record.
	Descendants string `json:"descendants"`
	// DescendantsAlt defines zero or more alternate URIs to display the descendants of a given record.
	DescendantsAlt []string `json:"descendants_alt"`
	// Index defines the URI for the initial landing (or index) page for the Spelunker.
	Index string `json:"index"`
	// Placetypes defines the URI for all the placetypes.
	Placetypes string `json:"placetypes"`
	// Placetypes defines the URI for all the record with a given placetype.
	Placetype string `json:"placetype"`
	// Placetypes defines the URI for all the records "visiting" Null Island (have a lat,lon of "0.0, 0.0").
	NullIsland string `json:"nullisland"`
	// Recent defined the URI for all the records that have been updated within a given time period.
	Recent string `json:"recent"`
	// RecentAlt defines zero or more alternate URIs to display records that have been updated within a given time period.
	RecentAlt []string `json:"recent_alt"`
	// Search defines the URI for searching the Spelunker.
	Search string `json:"search"`
	// About defined the URI for an "about the Spelunker" page.
	About string `json:"about"`
	// OpenSearch defines the URI for the OpenSearch browser plugin (search) definition..
	OpenSearch string `json:"opensearch"`

	// ConcordanceNSFaceted defines the URI for the API endpoint to return faceted results for a given namespace.
	ConcordanceNSFaceted string `json:"concordance_ns"`
	// ConcordanceNSPredFaceted defines the URI for the API endpoint to return faceted results for a namespace and predicate pair.
	ConcordanceNSPredFaceted string `json:"concordance_ns_pred"`
	// ConcordanceTripleFaceted defines the URI for the API endpoint to return faceted results for a concordance (ns:pred=value).
	ConcordanceTripleFaceted string `json:"concordance_triple_faceted"`
	// DescendantsFaceted defines the URI for the API endpoint to return faceted results for the descendants of a given record.
	DescendantsFaceted string `json:"descendants_faceted"`
	// FindingAid defines the URI for the API endpoint to return the repository (as defined by the "wof:repo" property) for a given ID.
	FindingAid string `json:"finding_aid"`
	// GeoJSON defines the URI for the API endpoint to render a Who's On First record as a GeoJSON Feature.
	GeoJSON string `json:"geojson"`
	// GeoJSON defines zero or more URIs for alternate API endpoints to render a Who's On First record as a GeoJSON Feature.
	GeoJSONAlt []string `json:"geojson_alt"`
	// GeoJSON defines the URI for the API endpoint to render a Who's On First record as a GeoJSON-LD Feature.
	GeoJSONLD string `json:"geojsonld"`
	// GeoJSON defines zero or more URIs for alternate API endpoints to render a Who's On First record as a GeoJSON-LD Feature.
	GeoJSONLDAlt []string `json:"geojsonld_alt"`
	// NavPlace defines the URI to render a Who's On First record as a IIIF NavPlace document.
	NavPlace string `json:"navplace"`
	// GeoJSON defines zero or more URIs for alternate API endpoints to render a Who's On First record as a IIIF NavPlace Feature.
	NavPlaceAlt []string `json:"navplace_alt"`
	// NullIslandFaceted defines the URI for the API endpoint to return faceted results for Who's Of First records "visiting" Null Island (have lat,lon coordinates of "0.0,0.0").
	NullIslandFaceted string `json:"nullisland_faceted"`
	// PlacetypeFaceted defines the URI for the API endpoint to return faceted results for records with a specific placetype.
	PlacetypeFaceted string `json:"placetype_faceted"`
	// RecentFaceted defines the URI for the API endpoint to return faceted results for records which have been updated within a given time period.
	RecentFaceted string `json:"recent_faceted"`
	// SearchFaceted defines the URI for the API endpoint to return faceted results for a search query.
	SearchFaceted string `json:"search_faceted"`
	// Select defines the URIs to emit specific properties in a a Who's On First record.
	Select string `json:"select"`
	// SelectAlt defines zero or more URIs for alternate API endpoints to emit specific properties in a a Who's On First record.
	SelectAlt []string `json:"select_alt"`
	// SPR defines the URI to render a Who's On First record as a Standard Places Response (SPR) document.
	SPR string `json:"spr"`
	// SPRAlt defines zero or more URIs for alternate API endpoints to a Who's On First record as a Standard Places Response (SPR) document.
	SPRAlt []string `json:"spr_alt"`
	// SPR defines the URI to render a Who's On First record as an SVG document.
	SVG string `json:"svg"`
	// SVGAlt defines zero or more URIs for alternate API endpoints to a Who's On First record as an SVG document.
	SVGAlt []string `json:"svg_alt"`
	// WKT defines the URI to render a Who's On First record's geometry property as "well-known text" (WKT).
	WKT string `json:"wkt"`
	// WKTAlt defines zero or more URIs for alternate API endpoints to a Who's On First record's geometry property as "well-known text" (WKT).
	WKTAlt []string `json:"wkt_alt"`

	// RootURL defines the root URL (inclusive of scheme, host and port details) for the Spelunker.
	RootURL string `json:"root_url"`
	// Static defines the URI for static assets (JavaScript, CSS, etc.).
	Static string `json:"static"`
}

// DefaultURIs returns a `URIs` struct with default values for Spelunker web application endpoint paths (URIs).
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
		WKT: "/id/{id}/wkt",
		WKTAlt: []string{
			"/wkt/",
		},
	}

	return uris_table
}

// Abs returns the fully-qualified URI for 'path'.
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
	id_uri := replaceAll(uri, "{id}", id)
	return uriWithFilters(id_uri, nil, nil)
}

func URIForId(uri string, id int64, filters []spelunker.Filter, facets []spelunker.Facet) string {

	id_uri := replaceAll(uri, "{id}", id)
	return uriWithFilters(id_uri, filters, facets)
}

func URIForPlacetype(uri string, pt string, filters []spelunker.Filter, facets []spelunker.Facet) string {

	pt_uri := replaceAll(uri, "{placetype}", pt)
	return uriWithFilters(pt_uri, filters, facets)
}

func URIForRecentSimple(uri string, d string) string {
	r_uri := replaceAll(uri, "{duration}", d)
	return uriWithFilters(r_uri, nil, nil)
}

func URIForRecent(uri string, d string, filters []spelunker.Filter, facets []spelunker.Facet) string {

	r_uri := replaceAll(uri, "{duration}", d)
	return uriWithFilters(r_uri, filters, facets)
}

func URIForConcordanceNS(uri string, ns string, filters []spelunker.Filter, facets []spelunker.Facet) string {

	c_uri := replaceAll(uri, "{namespace}", ns)
	return uriWithFilters(c_uri, filters, facets)
}

func URIForConcordanceNSPred(uri string, ns string, pred string, filters []spelunker.Filter, facets []spelunker.Facet) string {

	c_uri := uri

	c_uri = replaceAll(c_uri, "{namespace}", ns)
	c_uri = replaceAll(c_uri, "{predicate}", pred)
	return uriWithFilters(c_uri, filters, facets)
}

func URIForConcordanceTriple(uri string, ns string, pred string, value any, filters []spelunker.Filter, facets []spelunker.Facet) string {

	c_uri := uri

	c_uri = replaceAll(c_uri, "{namespace}", ns)
	c_uri = replaceAll(c_uri, "{predicate}", pred)
	c_uri = replaceAll(c_uri, "{value}", value)
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

func replaceAll(input string, pattern string, value any) string {
	str_value := fmt.Sprintf("%v", value)
	return strings.Replace(input, pattern, str_value, -1)
}

func (u *URIs) applyPrefix(prefix string) error {

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
