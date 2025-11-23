package sql

import (
	"context"

	"github.com/aaronland/go-pagination"
	wof_spr "github.com/whosonfirst/go-whosonfirst-spr/v2"
	"github.com/whosonfirst/spelunker/v2"
)

// GetAlternatePlacetypes retrieves the list of alternate placetype ("wof:placetype_alt") in a SQLSpelunker database.
func (s *SQLSpelunker) GetAlternatePlacetypes(ctx context.Context) (*spelunker.Faceting, error) {
	return nil, spelunker.ErrNotImplemented
}

// HasAlternatePlacetypes retrieves the list of Who's On First records with a given alternate placetype ("wof:placetype_alt") in a SQLSpelunker database.
func (s *SQLSpelunker) HasAlternatePlacetype(ctx context.Context, pg_opts pagination.Options, pt string, filters []spelunker.Filter) (wof_spr.StandardPlacesResults, pagination.Results, error) {
	return nil, nil, spelunker.ErrNotImplemented
}

// HasAlternatePlacetypeFaceted retrieves faceted properties for records with a given alternate placetype ("wof:placetype_alt") in a SQLSpelunker database.
func (s *SQLSpelunker) HasAlternatePlacetypeFaceted(ctx context.Context, pt string, filters []spelunker.Filter, facets []*spelunker.Facet) ([]*spelunker.Faceting, error) {
	return nil, spelunker.ErrNotImplemented
}
