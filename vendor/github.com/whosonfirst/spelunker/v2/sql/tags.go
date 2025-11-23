package sql

// Tags are not currently indexed in any SQL tables
// https://github.com/whosonfirst/go-whosonfirst-sql/tree/main/tables
// https://github.com/whosonfirst/go-whosonfirst-sqlite-features/tree/main/tables

import (
	"context"

	"github.com/aaronland/go-pagination"
	wof_spr "github.com/whosonfirst/go-whosonfirst-spr/v2"
	"github.com/whosonfirst/spelunker/v2"
)

// GetTags retrieves the list of unique tags in a Spelunker index in a SQLSpelunker database.
func (s *SQLSpelunker) GetTags(ctx context.Context) (*spelunker.Faceting, error) {
	return nil, spelunker.ErrNotImplemented
}

// HasTag retrieves the list of records that have a given tag in a SQLSpelunker database.
func (s *SQLSpelunker) HasTag(ctx context.Context, pg_opts pagination.Options, tag string, filters []spelunker.Filter) (wof_spr.StandardPlacesResults, pagination.Results, error) {
	return nil, nil, spelunker.ErrNotImplemented
}

// HasTagFaceted retrieves faceted properties for records that have a given tag in a SQLSpelunker database.
func (s *SQLSpelunker) HasTagFaceted(ctx context.Context, tag string, filters []spelunker.Filter, facets []*spelunker.Facet) ([]*spelunker.Faceting, error) {
	return nil, spelunker.ErrNotImplemented
}
