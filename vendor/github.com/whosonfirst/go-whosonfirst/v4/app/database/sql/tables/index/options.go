package index

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

type RunOptions struct {
	MaxProcesses      int
	Verbose           bool
	SpatialTables     bool
	SpelunkerTables   bool
	AllTables         bool
	RTreeTable        bool
	GeoJSONTable      bool
	PropertiesTable   bool
	SPRTable          bool
	SpelunkerTable    bool
	ConcordancesTable bool
	AncestorsTable    bool
	SearchTable       bool
	NamesTable        bool
	SupersedesTable   bool
	IndexAlt          []string
	StrictAltFiles    bool
	DatabaseURI       string
	IteratorURI       string
	IteratorSources   []string
	Optimize          bool
	IndexRelations    bool
	RelationsURI      string
}

func RunOptionsFromFlagSet(fs *flag.FlagSet) (*RunOptions, error) {

	flagset.Parse(fs)
	args := fs.Args()

	return RunOptionsFromParsedFlags(args...)
}

func RunOptionsFromParsedFlags(args ...string) (*RunOptions, error) {

	opts := &RunOptions{
		Verbose:           verbose,
		MaxProcesses:      procs,
		SpatialTables:     spatial_tables,
		SpelunkerTables:   spelunker_tables,
		AllTables:         all,
		RTreeTable:        rtree,
		GeoJSONTable:      geojson,
		PropertiesTable:   properties,
		SPRTable:          spr,
		SpelunkerTable:    spelunker,
		ConcordancesTable: concordances,
		AncestorsTable:    ancestors,
		SearchTable:       search,
		NamesTable:        names,
		SupersedesTable:   supersedes,
		IndexAlt:          index_alt,
		StrictAltFiles:    strict_alt_files,
		DatabaseURI:       db_uri,
		IteratorURI:       iterator_uri,
		IteratorSources:   args,
		Optimize:          optimize,
		IndexRelations:    index_relations,
		RelationsURI:      relations_uri,
	}

	return opts, nil
}
