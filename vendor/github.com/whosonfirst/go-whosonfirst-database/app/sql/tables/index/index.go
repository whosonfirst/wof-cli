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
	_ "github.com/sfomuseum/go-flags/flagset"
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

	opts, err := RunOptionsFromFlagSet(fs)

	if err != nil {
		return err
	}

	return RunWithOptions(ctx, opts)
}

func RunWithOptions(ctx context.Context, opts *RunOptions) error {

	if opts.Verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("Verbose logging enabled")
	}

	runtime.GOMAXPROCS(opts.MaxProcesses)

	if opts.SpatialTables {
		opts.RTreeTable = true
		opts.GeoJSONTable = true
		opts.PropertiesTable = true
		opts.SPRTable = true
	}

	if opts.SpelunkerTables {
		// rtree = true
		opts.SPRTable = true
		opts.SpelunkerTable = true
		opts.GeoJSONTable = true
		opts.ConcordancesTable = true
		opts.AncestorsTable = true
		opts.SearchTable = true

		to_index_alt := []string{
			tables.GEOJSON_TABLE_NAME,
		}

		for _, table_name := range to_index_alt {

			if !slices.Contains(index_alt, table_name) {
				opts.IndexAlt = append(opts.IndexAlt, table_name)
			}
		}

	}

	logger := slog.Default()

	db, err := database_sql.OpenWithURI(ctx, opts.DatabaseURI)

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
		if opts.Optimize {

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
		RTree:           opts.RTreeTable,
		GeoJSON:         opts.GeoJSONTable,
		Properties:      opts.PropertiesTable,
		SPR:             opts.SPRTable,
		Spelunker:       opts.SpelunkerTable,
		Concordances:    opts.ConcordancesTable,
		Ancestors:       opts.AncestorsTable,
		Search:          opts.SearchTable,
		Names:           opts.NamesTable,
		Supersedes:      opts.SupersedesTable,
		SpatialTables:   opts.SpatialTables,
		SpelunkerTables: opts.SpelunkerTables,
		All:             opts.AllTables,
		IndexAlt:        opts.IndexAlt,
		StrictAltFiles:  opts.StrictAltFiles,
	}

	to_index, err := tables.InitTables(ctx, db, init_opts)

	if err != nil {
		return err
	}

	if len(to_index) == 0 {
		return fmt.Errorf("You forgot to specify which (any) tables to index")
	}

	record_opts := &indexer.LoadRecordFuncOptions{
		StrictAltFiles: opts.StrictAltFiles,
	}

	record_func := indexer.LoadRecordFunc(record_opts)

	idx_opts := &indexer.IndexerOptions{
		DB:             db,
		Tables:         to_index,
		LoadRecordFunc: record_func,
		Workers:        opts.MaxProcesses,
	}

	if opts.IndexRelations {

		r, err := reader.NewReader(ctx, opts.RelationsURI)

		if err != nil {
			return fmt.Errorf("Failed to load reader (%s), %v", opts.RelationsURI, err)
		}

		belongsto_func := indexer.IndexRelationsFunc(r)
		idx_opts.PostIndexFunc = belongsto_func
	}

	idx, err := indexer.NewIndexer(idx_opts)

	if err != nil {
		return fmt.Errorf("failed to create sqlite indexer because %v", err)
	}

	err = idx.IndexURIs(ctx, opts.IteratorURI, opts.IteratorSources...)

	if err != nil {
		return fmt.Errorf("Failed to index paths in %s mode because: %s", iterator_uri, err)
	}

	return nil
}
