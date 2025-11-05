package sql

import (
	"context"
	"fmt"
	"strings"

	"github.com/aaronland/go-pagination"
	"github.com/whosonfirst/go-whosonfirst-placetypes"
	"github.com/whosonfirst/go-whosonfirst-spelunker"
	wof_spr "github.com/whosonfirst/go-whosonfirst-spr/v2"
	"github.com/whosonfirst/go-whosonfirst-sql/tables"
)

func (s *SQLSpelunker) GetPlacetypes(ctx context.Context) (*spelunker.Faceting, error) {

	facet_counts := make([]*spelunker.FacetCount, 0)

	// TBD alt files...
	q := fmt.Sprintf("SELECT placetype, COUNT(id) AS count FROM %s WHERE is_alt=0 GROUP BY placetype ORDER BY count DESC", tables.SPR_TABLE_NAME)

	rows, err := s.db.QueryContext(ctx, q)

	if err != nil {
		return nil, fmt.Errorf("Failed to execute query, %w", err)
	}

	for rows.Next() {

		var pt string
		var count int64

		err := rows.Scan(&pt, &count)

		if err != nil {
			return nil, fmt.Errorf("Failed to scan row, %w", err)
		}

		f := &spelunker.FacetCount{
			Key:   pt,
			Count: count,
		}

		facet_counts = append(facet_counts, f)
	}

	err = rows.Close()

	if err != nil {
		return nil, fmt.Errorf("Failed to close results rows, %w", err)
	}

	f := spelunker.NewFacet("placetype")

	faceting := &spelunker.Faceting{
		Facet:   f,
		Results: facet_counts,
	}

	return faceting, nil
}

func (s *SQLSpelunker) HasPlacetype(ctx context.Context, pg_opts pagination.Options, pt *placetypes.WOFPlacetype, filters []spelunker.Filter) (wof_spr.StandardPlacesResults, pagination.Results, error) {

	where, args, err := s.hasPlacetypeQueryWhere(pt, filters)

	if err != nil {
		return nil, nil, fmt.Errorf("Failed to derive placetype query, %w", err)
	}

	str_where := strings.Join(where, " AND ")
	return s.querySPR(ctx, pg_opts, str_where, args...)
}

func (s *SQLSpelunker) HasPlacetypeFaceted(ctx context.Context, pt *placetypes.WOFPlacetype, filters []spelunker.Filter, facets []*spelunker.Facet) ([]*spelunker.Faceting, error) {

	q_where, q_args, err := s.hasPlacetypeQueryWhere(pt, filters)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive query where statement, %w", err)
	}

	results := make([]*spelunker.Faceting, len(facets))

	// START OF do this in go routines

	for idx, f := range facets {

		q := s.hasPlacetypeQueryFacetStatement(ctx, f, q_where)
		// slog.Info("FACET", "q", q, "args", q_args)

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

func (s *SQLSpelunker) hasPlacetypeQueryWhere(pt *placetypes.WOFPlacetype, filters []spelunker.Filter) ([]string, []interface{}, error) {

	where := []string{
		"placetype = ?",
	}

	args := []interface{}{
		pt.Name,
	}

	where, args, err := s.assignFilters(where, args, filters)

	if err != nil {
		return nil, nil, err
	}

	return where, args, nil
}

func (s *SQLSpelunker) hasPlacetypeQueryFacetStatement(ctx context.Context, facet *spelunker.Facet, where []string) string {

	facet_label := s.facetLabel(facet)

	cols := []string{
		fmt.Sprintf("%s.%s AS %s", tables.SPR_TABLE_NAME, facet_label, facet),
		fmt.Sprintf("COUNT(%s.id) AS count", tables.SPR_TABLE_NAME),
	}

	q := s.hasPlacetypeQueryStatement(ctx, cols, where)
	return fmt.Sprintf("%s GROUP BY %s.%s ORDER BY count DESC", q, tables.SPR_TABLE_NAME, facet_label)
}

func (s *SQLSpelunker) hasPlacetypeQueryStatement(ctx context.Context, cols []string, where []string) string {

	str_cols := strings.Join(cols, ",")
	str_where := strings.Join(where, " AND ")

	return fmt.Sprintf("SELECT %s FROM %s WHERE %s", str_cols, tables.SPR_TABLE_NAME, str_where)

}
