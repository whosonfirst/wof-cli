package emit

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/aaronland/go-json-query"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/multi"
	"github.com/whosonfirst/go-whosonfirst-iterate/v3"
	"github.com/whosonfirst/go-writer/v3"
)

var writer_uri string
var iterator_uri string

var forgiving bool
var include_alt_geoms bool

var as_spr bool
var as_spr_geojson bool

var format string

var queries query.QueryFlags
var query_mode string

var csv_append_properties multi.KeyValueString

func DefaultFlagSet() *flag.FlagSet {

	writer_schemes := strings.Join(writer.Schemes(), ", ")
	iter_schemes := strings.Join(iterate.IteratorSchemes(), ", ")

	fs := flagset.NewFlagSet("emit")

	fs.StringVar(&writer_uri, "writer-uri", "jsonl://?writer=stdout://", "A valid whosonfirst/go-writer.Writer URI. Available options are: "+writer_schemes)

	fs.StringVar(&iterator_uri, "iterator-uri", "repo://", "A valid whosonfirst/go-whosonfirst-iterate/v3.Iterator URI. If URI is \"-\" then this flag will be assigned a value of \"file://\" whose input will be the expanded URIs derived from additional arguments. Available options are: "+iter_schemes)

	fs.BoolVar(&as_spr, "as-spr", false, "Emit Who's On First records formatted as Standard Place Response (SPR) records. This flag is DEPRECATED. Please use '-format spr' instead.")
	fs.BoolVar(&as_spr_geojson, "as-spr-geojson", false, "Emit Who's On First records as GeoJSON records where the 'properties' element is replaced by a Standard Place Response (SPR) representation of the record. This flag is DEPRECATED. Please use '-format geojson' instead.")

	fs.StringVar(&format, "format", "", "Valid options are: csv, spr, spr-geojson or [none]. If none then the raw GeoJSON for each matching record will be emitted.")

	fs.BoolVar(&forgiving, "forgiving", false, "Do not stop processing when errors are encountered.")

	fs.BoolVar(&include_alt_geoms, "include-alt-geoms", true, "Emit alternate geometry records.")

	valid_modes := strings.Join([]string{query.QUERYSET_MODE_ALL, query.QUERYSET_MODE_ANY}, ", ")
	desc_modes := fmt.Sprintf("Specify how query filtering should be evaluated. Valid modes are: %s", valid_modes)

	fs.Var(&csv_append_properties, "csv-append-property", "Zero or more additional properties to append to each CSV row. Properties should be in the format of {COLUMN_NAME}={PATH}. This flag is only honoured if the -format flag has a value of \"csv\".")

	fs.Var(&queries, "query", "One or more {PATH}={REGEXP} parameters for filtering records.")
	fs.StringVar(&query_mode, "query-mode", query.QUERYSET_MODE_ALL, desc_modes)

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Emit one or more Who's On First records.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s emit [options] path(N) path(N)\n", os.Args[0])
		fs.PrintDefaults()
	}

	return fs
}
