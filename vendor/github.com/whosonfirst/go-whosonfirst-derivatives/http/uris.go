package http

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

type URIs struct {
	GeoJSON      string   `json:"geojson"`
	GeoJSONAlt   []string `json:"geojson_alt"`
	GeoJSONLD    string   `json:"geojsonld"`
	GeoJSONLDAlt []string `json:"geojsonld_alt"`
	NavPlace     string   `json:"navplace"`
	NavPlaceAlt  []string `json:"navplace_alt"`
	Select       string   `json:"select"`
	SelectAlt    []string `json:"select_alt"`
	SPR          string   `json:"spr"`
	SPRAlt       []string `json:"spr_alt"`
	SVG          string   `json:"svg"`
	SVGAlt       []string `json:"svg_alt"`
	WKT          string   `json:"wkt"`
	WKTAlt       []string `json:"wkt_alt"`
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

		// API/machine-readable

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
		Select: "/id/{id}/select",
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
