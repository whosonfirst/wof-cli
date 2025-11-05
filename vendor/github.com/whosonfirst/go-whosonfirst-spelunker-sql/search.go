package sql

import (
	"context"
	"fmt"
	"strings"

	"github.com/aaronland/go-pagination"
	"github.com/whosonfirst/go-whosonfirst-spelunker"
	wof_spr "github.com/whosonfirst/go-whosonfirst-spr/v2"
	"github.com/whosonfirst/go-whosonfirst-sql/tables"
)

func (s *SQLSpelunker) Search(ctx context.Context, pg_opts pagination.Options, search_opts *spelunker.SearchOptions, filters []spelunker.Filter) (wof_spr.StandardPlacesResults, pagination.Results, error) {

	where, args, err := s.searchQueryWhere(search_opts, filters)

	if err != nil {
		return nil, nil, err
	}

	str_where := strings.Join(where, " AND ")

	if len(filters) == 0 {
		str_where := strings.Join(where, " AND ")
		return s.querySearch(ctx, pg_opts, str_where, args...)
	}

	return s.querySearchWithFilters(ctx, pg_opts, str_where, args...)
}

func (s *SQLSpelunker) SearchFaceted(ctx context.Context, search_opts *spelunker.SearchOptions, filters []spelunker.Filter, facets []*spelunker.Facet) ([]*spelunker.Faceting, error) {

	q_where, q_args, err := s.searchQueryWhere(search_opts, filters)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive query where statement, %w", err)
	}

	where := strings.Join(q_where, " AND ")

	results := make([]*spelunker.Faceting, len(facets))

	// START OF do this in go routines

	for idx, f := range facets {

		facet_label := s.facetLabel(f)

		q := fmt.Sprintf("SELECT %s.%s AS %s, COUNT(%s.id) AS count FROM %s JOIN %s ON %s.id = CAST(%s.id AS INTEGER) WHERE %s GROUP BY %s.%s ORDER BY count DESC",
			tables.SPR_TABLE_NAME,
			facet_label,
			facet_label,
			tables.SPR_TABLE_NAME,
			tables.SEARCH_TABLE_NAME,
			tables.SPR_TABLE_NAME,
			tables.SEARCH_TABLE_NAME,
			tables.SPR_TABLE_NAME,
			where,
			tables.SPR_TABLE_NAME,
			facet_label,
		)

		counts, err := s.facetWithQuery(ctx, q, q_args...)

		if err != nil {
			return nil, fmt.Errorf("Failed to facet columns, %w", err)
		}

		fc := &spelunker.Faceting{
			Facet:   f,
			Results: counts,
		}

		results[idx] = fc
	}

	// END OF do this in go routines

	return results, nil
}

func (s *SQLSpelunker) searchQueryWhere(search_opts *spelunker.SearchOptions, filters []spelunker.Filter) ([]string, []interface{}, error) {

	if len(filters) == 0 {

		where := []string{
			"names_all MATCH ?",
		}

		args := []interface{}{
			search_opts.Query,
		}

		return s.assignFilters(where, args, filters)
	}

	// join search on spr table

	where := []string{
		fmt.Sprintf("%s.names_all MATCH ?", tables.SEARCH_TABLE_NAME),
	}

	args := []interface{}{
		search_opts.Query,
	}

	return s.assignFilters(where, args, filters)
}
