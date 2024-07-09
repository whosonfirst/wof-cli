package emit

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/aaronland/go-json-query"
	"github.com/sfomuseum/go-flags/flagset"
)

var writer_uri string
var iterator_uri string

var forgiving bool
var include_alt_geoms bool

var as_spr bool
var as_spr_geojson bool

var queries query.QueryFlags
var query_mode string

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("emit")

	fs.StringVar(&writer_uri, "writer-uri", "jsonl://?writer=stdout://", "A valid whosonfirst/go-writer.Writer URI.")
	fs.StringVar(&iterator_uri, "iterator-uri", "repo://", "A valid whosonfirst/go-whosonfirst-iterate/v2/emitter URI. If URI is \"-\" then this flag will be assigned a value of \"file://\" whose input will be the expanded URIs derived from additional arguments.")

	fs.BoolVar(&as_spr, "as-spr", false, "Emit Who's On First records formatted as Standard Place Response (SPR) records.")
	fs.BoolVar(&as_spr_geojson, "as-spr-geojson", false, "Emit Who's On First records as GeoJSON records where the 'properties' element is replaced by a Standard Place Response (SPR) representation of the record.")

	fs.BoolVar(&include_alt_geoms, "include-alt-geoms", true, "Emit alternate geometry records.")

	valid_modes := strings.Join([]string{query.QUERYSET_MODE_ALL, query.QUERYSET_MODE_ANY}, ", ")
	desc_modes := fmt.Sprintf("Specify how query filtering should be evaluated. Valid modes are: %s", valid_modes)

	fs.Var(&queries, "query", "One or more {PATH}={REGEXP} parameters for filtering records.")
	fs.StringVar(&query_mode, "query-mode", query.QUERYSET_MODE_ALL, desc_modes)

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Emit one or more Who's On First records.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s emit [options] path(N) path(N)\n", os.Args[0])
		fs.PrintDefaults()
	}

	return fs
}
