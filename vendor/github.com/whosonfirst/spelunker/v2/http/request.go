package http

import (
	"context"
	"fmt"
	go_http "net/http"

	"github.com/aaronland/go-http/v4/sanitize"
	"github.com/aaronland/go-pagination"
	"github.com/aaronland/go-pagination/countable"
	"github.com/aaronland/go-pagination/cursor"
	"github.com/whosonfirst/go-whosonfirst-spr/v2"
	wof_http "github.com/whosonfirst/go-whosonfirst/http"
	"github.com/whosonfirst/spelunker/v2"
)

// PaginationOptionsFromRequests derives a new `pagination.Options` instance from query parameters present in 'req'.
func PaginationOptionsFromRequest(req *go_http.Request) (pagination.Options, error) {

	q_cursor, err := sanitize.GetString(req, "cursor")

	if err != nil {
		return nil, fmt.Errorf("Failed to derive ?cursor= parameter, %w", err)
	}

	if q_cursor != "" {

		pg_opts, err := cursor.NewCursorOptions()

		if err != nil {
			return nil, fmt.Errorf("Failed to create cursor options, %w", err)
		}

		pg_opts.Pointer(q_cursor)
		return pg_opts, nil
	}

	page, err := sanitize.GetInt64(req, "page")

	if err != nil {
		return nil, fmt.Errorf("Failed to derive ?page= parameter, %w", err)
	}

	if page == 0 {
		page = 1
	}

	pg_opts, err := countable.NewCountableOptions()

	if err != nil {
		return nil, fmt.Errorf("Failed to create countable options, %w", err)
	}

	pg_opts.Pointer(page)
	return pg_opts, nil
}

// ParsePageNumberFromRequest derives a pagination page number from the 'page' query parameter in 'req'.
func ParsePageNumberFromRequest(req *go_http.Request) (int64, error) {

	page, err := sanitize.GetInt64(req, "page")

	if err != nil {
		return 0, fmt.Errorf("Failed to derive ?page= parameter, %w", err)
	}

	if page == 0 {
		page = 1
	}

	return page, nil
}

// FeatureFromRequestURI returns the GeoJSON Feature for the Who's On First ID derived from 'req_uri'.
func FeatureFromRequestURI(ctx context.Context, sp spelunker.Spelunker, req_uri *wof_http.URI) ([]byte, error) {

	wof_id := req_uri.Id

	f, err := sp.GetFeatureForId(ctx, wof_id, req_uri.URIArgs)

	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve feature for %d, %w", wof_id, err)
	}

	return f, nil
}

// RecordFromRequestURI returns the internal Spelunker record (representation) for the Who's On First ID derived from 'req_uri'.
// Remember: This might be a GeoJSON Feature but it might not be.
func RecordFromRequestURI(ctx context.Context, sp spelunker.Spelunker, req_uri *wof_http.URI) ([]byte, error) {

	wof_id := req_uri.Id

	f, err := sp.GetRecordForId(ctx, wof_id, req_uri.URIArgs)

	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve record for %d, %w", wof_id, err)
	}

	return f, nil
}

// SPRFromRequestURI returns the Standard Places Response (SPR) for the Who's On First ID derived from 'req_uri'.
func SPRFromRequestURI(ctx context.Context, sp spelunker.Spelunker, req_uri *wof_http.URI) (spr.StandardPlacesResult, error) {

	wof_id := req_uri.Id

	f, err := sp.GetSPRForId(ctx, wof_id, req_uri.URIArgs)

	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve SPR for %d, %w", wof_id, err)
	}

	return f, nil
}
