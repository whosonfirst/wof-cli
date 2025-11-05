package index

import (
	"flag"
	"fmt"
	"runtime"
	"strings"

	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/multi"
	"github.com/whosonfirst/go-whosonfirst-iterate/v3"
)

var iterator_uri string

var db_uri string

var all bool
var ancestors bool
var concordances bool
var geojson bool
var spelunker bool
var geometries bool
var names bool
var rtree bool
var properties bool
var search bool
var spr bool
var supersedes bool

var spatial_tables bool
var spelunker_tables bool

var optimize bool

var strict_alt_files bool
var index_alt multi.MultiString

var index_relations bool
var relations_uri string

var procs int
var verbose bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("index")

	valid_schemes := strings.Join(iterate.IteratorSchemes(), ",")
	iterator_desc := fmt.Sprintf("A valid whosonfirst/go-whosonfirst-iterate/v3.Iterator URI. Supported iterator URI schemes are: %s", valid_schemes)

	fs.StringVar(&iterator_uri, "iterator-uri", "repo://", iterator_desc)

	fs.StringVar(&db_uri, "database-uri", "", "A URI in the form of 'sql://{DATABASE_SQL_ENGINE}?dsn={DATABASE_SQL_DSN}'. For example: sql://sqlite3?dsn=test.db")

	fs.BoolVar(&all, "all", false, "Index all tables (except the 'search' and 'geometries' tables which you need to specify explicitly)")
	fs.BoolVar(&ancestors, "ancestors", false, "Index the 'ancestors' tables")
	fs.BoolVar(&concordances, "concordances", false, "Index the 'concordances' tables")
	fs.BoolVar(&geojson, "geojson", false, "Index the 'geojson' table")
	fs.BoolVar(&spelunker, "spelunker", false, "Index the 'spelunker' table")
	fs.BoolVar(&geometries, "geometries", false, "Index the 'geometries' table (requires that libspatialite already be installed if using SQLite)")
	fs.BoolVar(&names, "names", false, "Index the 'names' table")
	fs.BoolVar(&rtree, "rtree", false, "Index the 'rtree' table")
	fs.BoolVar(&properties, "properties", false, "Index the 'properties' table")
	fs.BoolVar(&search, "search", false, "Index the 'search' table. If using the SQLite FTS5 full-text indexer requires the `fts5` build tag.")
	fs.BoolVar(&spr, "spr", false, "Index the 'spr' table")
	fs.BoolVar(&supersedes, "supersedes", false, "Index the 'supersedes' table")

	fs.BoolVar(&spatial_tables, "spatial-tables", false, "If true then index the necessary tables for use with the whosonfirst/go-whosonfirst-spatial-sqlite package.")
	fs.BoolVar(&spelunker_tables, "spelunker-tables", false, "If true then index the necessary tables for use with the whosonfirst/go-whosonfirst-spelunker packages")

	fs.BoolVar(&optimize, "optimize", true, "Attempt to optimize the database before closing connection")

	fs.Var(&index_alt, "index-alt", "Zero or more table names where alt geometry files should be indexed.")
	fs.BoolVar(&strict_alt_files, "strict-alt-files", true, "Be strict when indexing alt geometries")

	fs.BoolVar(&index_relations, "index-relations", false, "Index the records related to a feature, specifically wof:belongsto, wof:depicts and wof:involves. Alt files for relations are not indexed at this time.")
	fs.StringVar(&relations_uri, "index-relations-reader-uri", "", "A valid go-reader.Reader URI from which to read data for a relations candidate.")

	fs.IntVar(&procs, "processes", (runtime.NumCPU() * 2), "The number of concurrent processes to index data with")

	fs.BoolVar(&verbose, "verbose", false, "Enable verbose (debug) logging")
	return fs
}
