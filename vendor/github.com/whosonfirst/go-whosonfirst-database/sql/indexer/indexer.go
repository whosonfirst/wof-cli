package indexer

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log/slog"
	"sync"
	"time"

	database_sql "github.com/sfomuseum/go-database/sql"
	"github.com/whosonfirst/go-whosonfirst-iterate/v3"
)

// IndexerPostIndexFunc is a custom function to invoke after a record has been indexed.
type IndexerPostIndexFunc func(context.Context, *sql.DB, []database_sql.Table, interface{}) error

// IndexerLoadRecordFunc is a custom function to be invoked for each record processed by the `IndexURIs` method.
type IndexerLoadRecordFunc func(context.Context, string, io.ReadSeeker, ...interface{}) (interface{}, error)

// Indexer is a struct that provides methods for indexing records in one or more SQLite database_sql.tables
type Indexer struct {
	table_timings map[string]time.Duration
	mu            *sync.RWMutex
	options       *IndexerOptions
	// Timings is a boolean flag indicating whether timings (time to index records) should be recorded)
	Timings bool
}

// IndexerOptions
type IndexerOptions struct {
	// DB is the `database_sql.sql.DB` instance that records will be indexed in.
	DB *sql.DB
	// Tables is the list of `sfomuseum/go-database_sql.Table` instances that records will be indexed in.
	Tables []database_sql.Table
	// LoadRecordFunc is a custom `whosonfirst/go-whosonfirst-iterate/v2` callback function to be invoked
	// for each record processed by	the `IndexURIs`	method.
	LoadRecordFunc IndexerLoadRecordFunc
	// PostIndexFunc is an optional custom function to invoke after a record has been indexed.
	PostIndexFunc IndexerPostIndexFunc
}

// NewSQLiteInder returns a `Indexer` configured with 'opts'.
func NewIndexer(opts *IndexerOptions) (*Indexer, error) {

	table_timings := make(map[string]time.Duration)
	mu := new(sync.RWMutex)

	i := Indexer{
		table_timings: table_timings,
		mu:            mu,
		options:       opts,
		Timings:       false,
	}

	return &i, nil
}

// IndexURIs will index records returned by the `whosonfirst/go-whosonfirst-iterate` instance for 'uris',
func (idx *Indexer) IndexURIs(ctx context.Context, iterator_uri string, uris ...string) error {

	iter, err := iterate.NewIterator(ctx, iterator_uri)

	if err != nil {
		return fmt.Errorf("Failed to create new iterator, %w", err)
	}

	defer iter.Close()

	done_ch := make(chan bool)
	t1 := time.Now()

	// ideally this could be a proper stand-along package method but then
	// we have to set up a whole bunch of scaffolding just to pass 'indexer'
	// around so... we're not doing that (20180205/thisisaaronland)

	show_timings := func() {

		t2 := time.Since(t1)
		i := iter.Seen()

		slog.Info("Time to index all", "count", i, "time", t2)
	}

	if idx.Timings {

		go func() {

			for {

				select {
				case <-done_ch:
					return
				case <-time.After(1 * time.Minute):
					show_timings()
				}
			}
		}()

		defer func() {
			done_ch <- true
		}()
	}

	for rec, err := range iter.Iterate(ctx, uris...) {

		if err != nil {
			return err
		}

		logger := slog.Default()
		logger = logger.With("path", rec.Path)

		defer rec.Body.Close()

		err = idx.IndexIteratorRecord(ctx, rec)

		if err != nil {
			logger.Error("Failed to index record", "error", err)
			return err
		}
	}

	return nil
}

// IndexIterateRecord will index 'rec' in the underlying database.
func (idx *Indexer) IndexIteratorRecord(ctx context.Context, rec *iterate.Record) error {

	logger := slog.Default()
	logger = logger.With("path", rec.Path)

	record, err := idx.options.LoadRecordFunc(ctx, rec.Path, rec.Body)

	if err != nil {
		return err
	}

	if record == nil {
		return nil
	}

	idx.mu.Lock()
	idx.mu.Unlock()

	err = database_sql.IndexRecord(ctx, idx.options.DB, record, idx.options.Tables...)

	if err != nil {
		return fmt.Errorf("Failed to index record, %w", err)
	}

	if idx.options.PostIndexFunc != nil {

		err := idx.options.PostIndexFunc(ctx, idx.options.DB, idx.options.Tables, record)

		if err != nil {
			return err
		}
	}

	logger.Debug("Indexed database record")
	return nil
}
