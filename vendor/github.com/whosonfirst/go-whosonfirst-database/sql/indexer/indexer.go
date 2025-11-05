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
	// The maxiumum number of Go routines (workers) to use when indexing records.
	Workers int
}

// NewSQLiteInder returns a `Indexer` configured with 'opts'.
func NewIndexer(opts *IndexerOptions) (*Indexer, error) {

	table_timings := make(map[string]time.Duration)
	mu := new(sync.RWMutex)

	i := Indexer{
		table_timings: table_timings,
		mu:            mu,
		options:       opts,
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

	iter_workers := idx.options.Workers
	iter_throttle := make(chan bool, iter_workers)

	for i := 0; i < iter_workers; i++ {
		iter_throttle <- true
	}

	iter_ctx, iter_cancel := context.WithCancel(ctx)
	iter_wg := new(sync.WaitGroup)

	iter_err_ch := make(chan error)
	var iter_err error

	go func() {

		for {
			select {
			case err := <-iter_err_ch:
				iter_err = err
				iter_cancel()
				return
			}
		}
	}()

	for rec, err := range iter.Iterate(iter_ctx, uris...) {

		if err != nil {
			return err
		}

		select {
		case <-iter_ctx.Done():
			break
		default:
			<-iter_throttle
		}

		select {
		case <-iter_ctx.Done():
			break
		default:
			// carry on
		}

		iter_wg.Go(func() {

			logger := slog.Default()
			logger = logger.With("path", rec.Path)

			defer func() {
				rec.Body.Close()
				iter_throttle <- true
			}()

			err = idx.IndexIteratorRecord(ctx, rec)

			if err != nil {
				logger.Error("Failed to index record", "error", err)
				iter_err_ch <- err
			}
		})
	}

	iter_wg.Wait()

	if iter_err != nil {
		return iter_err
	}

	return nil
}

// IndexIterateRecord will index 'rec' in the underlying database.
func (idx *Indexer) IndexIteratorRecord(ctx context.Context, rec *iterate.Record) error {

	logger := slog.Default()
	logger = logger.With("path", rec.Path)

	record, err := idx.options.LoadRecordFunc(ctx, rec.Path, rec.Body)

	if err != nil {
		return fmt.Errorf("Failed to load record func, %w", err)
	}

	if record == nil {
		return nil
	}

	idx.mu.Lock()
	defer idx.mu.Unlock()

	err = database_sql.IndexRecord(ctx, idx.options.DB, record, idx.options.Tables...)

	if err != nil {
		return fmt.Errorf("Failed to index record, %w", err)
	}

	if idx.options.PostIndexFunc != nil {

		err := idx.options.PostIndexFunc(ctx, idx.options.DB, idx.options.Tables, record)

		if err != nil {
			return fmt.Errorf("Post index func failed, %w", err)
		}
	}

	logger.Debug("Indexed database record")
	return nil
}
