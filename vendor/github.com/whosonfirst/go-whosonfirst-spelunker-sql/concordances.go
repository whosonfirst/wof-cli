package sql

import (
	"context"
	"fmt"
	_ "log/slog"
	"strings"

	"github.com/aaronland/go-pagination"
	"github.com/aaronland/go-pagination/countable"
	"github.com/whosonfirst/go-whosonfirst-spelunker"
	wof_spr "github.com/whosonfirst/go-whosonfirst-spr/v2"
	"github.com/whosonfirst/go-whosonfirst-sql/tables"
	"github.com/whosonfirst/go-whosonfirst-sqlite-spr"
)

func (s *SQLSpelunker) GetConcordances(ctx context.Context) (*spelunker.Faceting, error) {

	facet_counts := make([]*spelunker.FacetCount, 0)

	q := fmt.Sprintf("SELECT other_source, COUNT(other_id) AS count FROM %s GROUP BY other_source ORDER BY count DESC", tables.CONCORDANCES_TABLE_NAME)

	rows, err := s.db.QueryContext(ctx, q)

	if err != nil {
		return nil, fmt.Errorf("Failed to execute query, %w", err)
	}

	for rows.Next() {

		var source string
		var count int64

		err := rows.Scan(&source, &count)

		if err != nil {
			return nil, fmt.Errorf("Failed to scan row, %w", err)
		}

		nspred := strings.Split(source, ":")
		ns := nspred[0]

		f := &spelunker.FacetCount{
			Key:   ns,
			Count: count,
		}

		facet_counts = append(facet_counts, f)
	}

	err = rows.Close()

	if err != nil {
		return nil, fmt.Errorf("Failed to close results rows, %w", err)
	}

	f := spelunker.NewFacet("concordance")

	faceting := &spelunker.Faceting{
		Facet:   f,
		Results: facet_counts,
	}

	return faceting, nil
}

func (s *SQLSpelunker) HasConcordance(ctx context.Context, pg_opts pagination.Options, namespace string, predicate string, value any, filters []spelunker.Filter) (wof_spr.StandardPlacesResults, pagination.Results, error) {

	var q string

	where := make([]string, 0)
	args := make([]interface{}, 0)

	switch {
	case namespace != "" && predicate != "":
		where = append(where, fmt.Sprintf("%s.other_source = ?", tables.CONCORDANCES_TABLE_NAME))
		args = append(args, fmt.Sprintf("%s:%s", namespace, predicate))
	case namespace != "":
		where = append(where, fmt.Sprintf("%s.other_source LIKE ?", tables.CONCORDANCES_TABLE_NAME))
		args = append(args, namespace+":%")
	case predicate != "":
		where = append(where, fmt.Sprintf("%s.other_source LIKE ?", tables.CONCORDANCES_TABLE_NAME))
		args = append(args, "%:"+predicate)
	default:
		return nil, nil, fmt.Errorf("Missing namespace and predicate")
	}

	if value != "" {
		where = append(where, fmt.Sprintf("%s.other_id = ?", tables.CONCORDANCES_TABLE_NAME))
		args = append(args, value)
	}

	var err error

	where, args, err = s.assignFilters(where, args, filters)

	if err != nil {
		return nil, nil, err
	}

	if len(filters) == 0 {

		str_where := strings.Join(where, " AND ")

		q = fmt.Sprintf("SELECT id FROM %s WHERE %s", tables.CONCORDANCES_TABLE_NAME, str_where)

	} else {

		// slog.Info("WHERE", "where", where, "args", args)
		str_where := strings.Join(where, " AND ")

		q = fmt.Sprintf("SELECT %s.id AS id FROM %s LEFT JOIN %s ON %s.id = %s.id WHERE %s",
			tables.SPR_TABLE_NAME,
			tables.SPR_TABLE_NAME,
			tables.CONCORDANCES_TABLE_NAME,
			tables.SPR_TABLE_NAME,
			tables.CONCORDANCES_TABLE_NAME,
			str_where,
		)

	}

	// slog.Info("QUERY", "q", q, "args", args)

	rows, err := s.db.QueryContext(ctx, q, args...)

	if err != nil {
		return nil, nil, fmt.Errorf("Failed to execute query, %w", err)
	}

	ids := make([]interface{}, 0)
	qms := make([]string, 0)

	for rows.Next() {

		var str_id int64

		err := rows.Scan(&str_id)

		if err != nil {
			return nil, nil, fmt.Errorf("Failed to scan row, %w", err)
		}

		ids = append(ids, str_id)
		qms = append(qms, "?")
	}

	err = rows.Close()

	if err != nil {
		return nil, nil, fmt.Errorf("Failed to close results rows, %w", err)
	}

	if len(ids) == 0 {

		var pg_results pagination.Results
		var pg_err error

		if pg_opts != nil {
			pg_results, pg_err = countable.NewResultsFromCountWithOptions(pg_opts, 0)
		} else {
			pg_results, pg_err = countable.NewResultsFromCount(0)
		}

		if pg_err != nil {
			return nil, nil, fmt.Errorf("Failed to create pagination results, %w", err)
		}

		results := make([]wof_spr.StandardPlacesResult, 0)

		spr_results := &spr.SQLiteResults{
			Places: results,
		}

		return spr_results, pg_results, nil
	}

	spr_where := []string{
		fmt.Sprintf("id IN (%s)", strings.Join(qms, ",")),
	}

	str_spr_where := strings.Join(spr_where, " AND ")
	return s.querySPR(ctx, pg_opts, str_spr_where, ids...)
}

func (s *SQLSpelunker) HasConcordanceFaceted(ctx context.Context, namespace string, predicate string, value any, filters []spelunker.Filter, facets []*spelunker.Facet) ([]*spelunker.Faceting, error) {

	where := make([]string, 0)
	args := make([]interface{}, 0)

	switch {
	case namespace != "" && predicate != "":
		where = append(where, fmt.Sprintf("%s.other_source = ?", tables.CONCORDANCES_TABLE_NAME))
		args = append(args, fmt.Sprintf("%s:%s", namespace, predicate))
	case namespace != "":
		where = append(where, fmt.Sprintf("%s.other_source LIKE ?", tables.CONCORDANCES_TABLE_NAME))
		args = append(args, namespace+":%")
	case predicate != "":
		where = append(where, fmt.Sprintf("%s.other_source LIKE ?", tables.CONCORDANCES_TABLE_NAME))
		args = append(args, "%:"+predicate)
	default:
		return nil, fmt.Errorf("Missing namespace and predicate")
	}

	if value != "" {
		where = append(where, fmt.Sprintf("%s.other_id = ?", tables.CONCORDANCES_TABLE_NAME))
		args = append(args, value)
	}

	where, args, err := s.assignFilters(where, args, filters)

	if err != nil {
		return nil, err
	}

	// slog.Info("WHERE", "where", where, "args", args)

	str_where := strings.Join(where, " AND ")

	facetings := make([]*spelunker.Faceting, len(facets))

	// START OF do this in go routines

	for idx, f := range facets {

		facet_label := s.facetLabel(f)

		q := fmt.Sprintf("SELECT %s.%s AS %s, COUNT(%s.id) AS count FROM %s LEFT JOIN %s ON %s.id = %s.id WHERE %s GROUP BY %s ORDER BY count DESC",
			tables.SPR_TABLE_NAME,
			facet_label,
			facet_label,
			tables.SPR_TABLE_NAME,
			tables.SPR_TABLE_NAME,
			tables.CONCORDANCES_TABLE_NAME,
			tables.SPR_TABLE_NAME,
			tables.CONCORDANCES_TABLE_NAME,
			str_where,
			facet_label,
		)

		// slog.Info("QUERY", "q", q, "args", args)

		rows, err := s.db.QueryContext(ctx, q, args...)

		if err != nil {
			return nil, err
		}

		facet_counts := make([]*spelunker.FacetCount, 0)

		for rows.Next() {

			var key string
			var count int64

			err := rows.Scan(&key, &count)

			if err != nil {
				return nil, fmt.Errorf("Failed to scan row, %w", err)
			}

			fc := &spelunker.FacetCount{
				Key:   key,
				Count: count,
			}

			facet_counts = append(facet_counts, fc)
		}

		faceting := &spelunker.Faceting{
			Facet:   f,
			Results: facet_counts,
		}

		facetings[idx] = faceting
	}

	// END OF do this in go routines

	return facetings, nil
}
