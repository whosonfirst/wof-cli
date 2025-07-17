package reader

import (
	"context"
	"fmt"
	"iter"
	"log/slog"
	"net/url"
	"sync/atomic"

	"github.com/whosonfirst/go-reader/v2"
	"github.com/whosonfirst/go-whosonfirst-iterate/v3"
	"github.com/whosonfirst/go-whosonfirst-iterate/v3/filters"
	wof_uri "github.com/whosonfirst/go-whosonfirst-uri"
)

func init() {
	ctx := context.Background()
	iterate.RegisterIterator(ctx, "reader", NewReaderIterator)
}

// GitIterator implements the `Iterator` interface for crawling records with a `whosonfirst/go-reader.Reader` instance.
type ReaderIterator struct {
	iterate.Iterator
	// reader is the `whosonfirst/go-reader.Reader` instance used to reade documents
	reader reader.Reader
	// filters is a `filters.Filters` instance used to include or exclude specific records from being crawled.
	filters filters.Filters
	// seen is the count of documents that have been processed.
	seen int64
	// iterating is a boolean value indicating whether records are still being iterated.
	iterating *atomic.Bool
}

// NewReaderIterator() returns a new `GitIterator` instance configured by 'uri' in the form of:
//
//	reader://?{PARAMETERS}
//
// Where {PATH} is an optional path on disk where a repository will be clone to (default is to clone repository in memory) and {PARAMETERS} may be:
// * `?include=` Zero or more `aaronland/go-json-query` query strings containing rules that must match for a document to be considered for further processing.
// * `?exclude=` Zero or more `aaronland/go-json-query`	query strings containing rules that if matched will prevent a document from being considered for further processing.
// * `?include_mode=` A valid `aaronland/go-json-query` query mode string for testing inclusion rules.
// * `?exclude_mode=` A valid `aaronland/go-json-query` query mode string for testing exclusion rules.
// * `?reader=` A valid `whosonfirst/go-reader` URI used to create the underlying reader instance.
func NewReaderIterator(ctx context.Context, uri string) (iterate.Iterator, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	q := u.Query()

	reader_uri := q.Get("reader")

	r, err := reader.NewReader(ctx, reader_uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new reader, %w", err)
	}

	f, err := filters.NewQueryFiltersFromURI(ctx, uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive query filters from URI, %w", err)
	}

	idx := &ReaderIterator{
		reader:    r,
		filters:   f,
		seen:      int64(0),
		iterating: new(atomic.Bool),
	}

	return idx, nil
}

// Iterate will return an `iter.Seq2[*Record, error]` for each record encountered in 'uris'.
func (it *ReaderIterator) Iterate(ctx context.Context, uris ...string) iter.Seq2[*iterate.Record, error] {

	return func(yield func(rec *iterate.Record, err error) bool) {

		it.iterating.Swap(true)
		defer it.iterating.Swap(false)

		for _, uri := range uris {

			logger := slog.Default()
			logger = logger.With("uri", uri)

			id, uri_args, err := wof_uri.ParseURI(uri)

			if err != nil {

				if !yield(nil, fmt.Errorf("Failed to parse '%s', %w", uri, err)) {
					return
				}

				continue
			}

			rel_path, err := wof_uri.Id2RelPath(id, uri_args)

			if err != nil {

				if !yield(nil, fmt.Errorf("Failed to derived relative path for '%s', %w", uri, err)) {
					return
				}

				continue
			}

			atomic.AddInt64(&it.seen, 1)

			logger = logger.With("path", rel_path)

			r, err := it.reader.Read(ctx, rel_path)

			if err != nil {

				if !yield(nil, fmt.Errorf("Failed to read path (%s) for '%s', %w", rel_path, uri, err)) {
					return
				}

				continue
			}

			if it.filters != nil {

				ok, err := iterate.ApplyFilters(ctx, r, it.filters)

				if err != nil {
					r.Close()
					if !yield(nil, fmt.Errorf("Failed to apply filters to %s, %w", rel_path, err)) {
						return
					}

					continue
				}

				if !ok {
					r.Close()
					continue
				}

			}

			rec := iterate.NewRecord(rel_path, r)
			yield(rec, nil)
		}
	}
}

// Seen() returns the total number of records processed so far.
func (it *ReaderIterator) Seen() int64 {
	return atomic.LoadInt64(&it.seen)
}

// IsIterating() returns a boolean value indicating whether 'it' is still processing documents.
func (it *ReaderIterator) IsIterating() bool {
	return it.iterating.Load()
}

// Close performs any implementation specific tasks before terminating the iterator.
func (it *ReaderIterator) Close() error {
	return nil
}
