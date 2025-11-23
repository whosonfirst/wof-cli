package server

import (
	"flag"
	"fmt"
	"os"

	"github.com/aaronland/go-http-maps/v2"
	"github.com/sfomuseum/go-flags/flagset"
)

var server_uri string
var spelunker_uri string
var authenticator_uri string

var map_provider string
var map_tile_uri string
var protomaps_theme string
var protomaps_max_data_zoom int

var root_url string

var verbose bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("spelunker")

	fs.StringVar(&server_uri, "server-uri", "http://localhost:8080", "A valid `aaronland/go-http/v3/server.Server URI.")
	fs.StringVar(&spelunker_uri, "spelunker-uri", "null://", "A URI in the form of '{SPELUNKER_SCHEME}://{IMPLEMENTATION_DETAILS}' referencing the underlying Spelunker database. For example: sql://sqlite3?dsn=spelunker.db")
	fs.StringVar(&authenticator_uri, "authenticator-uri", "null://", "A valid aaronland/go-http/v3/auth.Authenticator URI. This is future-facing work and can be ignored for now.")

	fs.StringVar(&map_provider, "map-provider", "leaflet", "Valid options are: leaflet, protomaps")
	fs.StringVar(&map_tile_uri, "map-tile-uri", maps.LEAFLET_OSM_TILE_URL, "A valid Leaflet tile layer URI. See documentation for special-case (interpolated tile) URIs.")
	fs.StringVar(&protomaps_theme, "protomaps-theme", "white", "A valid Protomaps theme label.")
	fs.IntVar(&protomaps_max_data_zoom, "protomaps-max-data-zoom", 15, "The maximum zoom (tile) level for data in a PMTiles database")

	fs.StringVar(&root_url, "root-url", "", "The root URL for all public-facing URLs and links. If empty then the value of the -server-uri flag will be used.")

	fs.BoolVar(&verbose, "verbose", false, "Enable verbose (debug) logging.")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Start the Spelunker web application.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t%s [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Valid options are:\n")
		fs.PrintDefaults()
	}

	return fs
}
