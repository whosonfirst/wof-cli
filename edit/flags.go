package edit

import (
	"flag"
	"fmt"
	"os"

	"github.com/aaronland/go-http-maps/v2"
	"github.com/sfomuseum/go-flags/flagset"
)

var reader_uri string
var writer_uri string

var access_token_uri string

var ensure_rel_path bool

var verbose bool

var map_provider string
var map_tile_uri string
var protomaps_theme string

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("edit")

	fs.StringVar(&reader_uri, "reader-uri", "", "An optional whosonfirst/go-reader/v2.Reader URI used to read WOF records from alternate sources. If defined then the -writer-uri flag must also be populated.")
	fs.StringVar(&writer_uri, "writer-uri", "", "An optional whosonfirst/go-writer/v3.Writer URI used to write records to alternate sources. If defined then the -reader-uri flag must also be populated.")

	fs.BoolVar(&ensure_rel_path, "ensure-relative-path", false, "Boolean flag signaling that each URI should be expanded to its fully-quality WOF-style relative path. This flag is only processed if the -reader-uri flag is not-empty.")

	fs.StringVar(&map_provider, "map-provider", "leaflet", "Valid options are: leaflet, protomaps")
	fs.StringVar(&map_tile_uri, "map-tile-uri", maps.LEAFLET_OSM_TILE_URL, "A valid Leaflet tile layer URI. See documentation for special-case (interpolated tile) URIs.")
	fs.StringVar(&protomaps_theme, "protomaps-theme", "light", "A valid Protomaps theme label.")

	fs.StringVar(&access_token_uri, "gh-access-token-uri", "", "A valid GitHub API access token. This is only necessary if -writer-uri is \"wof-pr://\".")
	fs.BoolVar(&verbose, "verbose", false, "Enable verbose (debug) logging.")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Launch a local web server running on a random port hosting a web application for editing one or more Who's On First records.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s path(N) path(N)\n", os.Args[0])
		fs.PrintDefaults()
	}

	return fs
}
