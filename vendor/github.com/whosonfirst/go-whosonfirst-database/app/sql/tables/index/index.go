package index

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"runtime"
	"slices"

	_ "github.com/whosonfirst/go-whosonfirst-database/sql"

	database_sql "github.com/sfomuseum/go-database/sql"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/whosonfirst/go-reader/v2"
	"github.com/whosonfirst/go-whosonfirst-database/sql/indexer"
	"github.com/whosonfirst/go-whosonfirst-database/sql/tables"
)

const index_alt_all string = "*"

func Run(ctx context.Context) error {
	fs := DefaultFlagSet()
	return RunWithFlagSet(ctx, fs)
}

// To do: Add RunWithOptions...

func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet) error {

	flagset.Parse(fs)

	if verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("Verbose logging enabled")
	}

	runtime.GOMAXPROCS(procs)

	if spatial_tables {
		rtree = true
		geojson = true
		properties = true
		spr = true
	}

	if spelunker_tables {
		// rtree = true
		spr = true
		spelunker = true
		geojson = true
		concordances = true
		ancestors = true
		search = true

		to_index_alt := []string{
			tables.GEOJSON_TABLE_NAME,
		}

		for _, table_name := range to_index_alt {

			if !slices.Contains(index_alt, table_name) {
				index_alt = append(index_alt, table_name)
			}
		}

	}

	logger := slog.Default()

	db, err := database_sql.OpenWithURI(ctx, db_uri)

	if err != nil {
		return err
	}

	defer func() {

		err := db.Close()

		if err != nil {
			logger.Error("Failed to close database connection", "error", err)
		}
	}()

	db_driver := database_sql.Driver(db)

	switch db_driver {
	case database_sql.POSTGRES_DRIVER:

	case database_sql.SQLITE_DRIVER:

		// optimize query performance
		// https://www.sqlite.org/pragma.html#pragma_optimize
		if optimize {

			defer func() {

				_, err = db.Exec("PRAGMA optimize")

				if err != nil {
					logger.Error("Failed to optimize", "error", err)
					return
				}
			}()

		}

	}

	init_opts := &tables.InitTablesOptions{
		RTree:           rtree,
		GeoJSON:         geojson,
		Properties:      properties,
		SPR:             spr,
		Spelunker:       spelunker,
		Concordances:    concordances,
		Ancestors:       ancestors,
		Search:          search,
		Names:           names,
		Supersedes:      supersedes,
		SpatialTables:   spatial_tables,
		SpelunkerTables: spelunker_tables,
		All:             all,
		IndexAlt:        index_alt,
		StrictAltFiles:  strict_alt_files,
	}

	to_index, err := tables.InitTables(ctx, db, init_opts)

	if err != nil {
		return err
	}

	if len(to_index) == 0 {
		return fmt.Errorf("You forgot to specify which (any) tables to index")
	}

	record_opts := &indexer.LoadRecordFuncOptions{
		StrictAltFiles: strict_alt_files,
	}

	record_func := indexer.LoadRecordFunc(record_opts)

	idx_opts := &indexer.IndexerOptions{
		DB:             db,
		Tables:         to_index,
		LoadRecordFunc: record_func,
	}

	if index_relations {

		r, err := reader.NewReader(ctx, relations_uri)

		if err != nil {
			return fmt.Errorf("Failed to load reader (%s), %v", relations_uri, err)
		}

		belongsto_func := indexer.IndexRelationsFunc(r)
		idx_opts.PostIndexFunc = belongsto_func
	}

	idx, err := indexer.NewIndexer(idx_opts)

	if err != nil {
		return fmt.Errorf("failed to create sqlite indexer because %v", err)
	}

	idx.Timings = timings

	uris := fs.Args()

	err = idx.IndexURIs(ctx, iterator_uri, uris...)

	if err != nil {
		return fmt.Errorf("Failed to index paths in %s mode because: %s", iterator_uri, err)
	}

	return nil
}
