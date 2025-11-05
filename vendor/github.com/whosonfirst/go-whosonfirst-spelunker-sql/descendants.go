package sql

// Dealing with descendants means querying the `ancestors` table or joining on
// the `ancestors` table and the `spr` table. Originally querying for descendants
// was done using SQLite's "instr" function but a) that probably wouldn't have
// worked for MySQL, etc. and b) it (the function) triggered timeouts and errors
// when querying using a remote VFS-enabled SQLite database, or at least that is
// the current theory.

import (
	"context"
	db_sql "database/sql"
	"fmt"
	// "log/slog"
	"strings"

	"github.com/aaronland/go-pagination"
	"github.com/aaronland/go-pagination/countable"
	"github.com/whosonfirst/go-whosonfirst-spelunker"
	wof_spr "github.com/whosonfirst/go-whosonfirst-spr/v2"
	"github.com/whosonfirst/go-whosonfirst-sql/tables"
	"github.com/whosonfirst/go-whosonfirst-sqlite-spr"
)

func (s *SQLSpelunker) GetDescendants(ctx context.Context, pg_opts pagination.Options, id int64, filters []spelunker.Filter) (wof_spr.StandardPlacesResults, pagination.Results, error) {

	q_where, q_args, err := s.descendantsQueryWhere(ctx, id, filters)

	if err != nil {
		return nil, nil, fmt.Errorf("Failed to derive query where statement, %w", err)
	}

	q_cols := s.descendantsQueryColumnsAll(ctx)

	q := s.descendantsQueryStatement(ctx, q_cols, q_where)

	if pg_opts != nil {
		limit, offset := s.deriveLimitOffset(pg_opts)
		q = fmt.Sprintf("%s LIMIT %d OFFSET %d", q, limit, offset)
	}

	pg_ch := make(chan pagination.Results)
	results_ch := make(chan wof_spr.StandardPlacesResults)

	done_ch := make(chan bool)
	err_ch := make(chan error)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {

		defer func() {
			done_ch <- true
		}()

		count_q := s.descendantsQueryCountStatement(ctx, q_where)

		count, err := s.queryCount(ctx, fmt.Sprintf("%s.id", tables.SPR_TABLE_NAME), count_q, q_args...)

		if err != nil {
			err_ch <- fmt.Errorf("Failed to derive query count, %w", err)
			return
		}

		var pg_results pagination.Results
		var pg_err error

		if pg_opts != nil {
			pg_results, pg_err = countable.NewResultsFromCountWithOptions(pg_opts, count)
		} else {
			pg_results, pg_err = countable.NewResultsFromCount(count)
		}

		if pg_err != nil {
			err_ch <- fmt.Errorf("Failed to derive pagination results, %w", pg_err)
			return
		}

		pg_ch <- pg_results
	}()

	go func() {

		defer func() {
			done_ch <- true
		}()

		rows, err := s.db.QueryContext(ctx, q, q_args...)

		if err != nil {
			err_ch <- fmt.Errorf("Failed to query where '%s', %w", q, err)
			return
		}

		results := make([]wof_spr.StandardPlacesResult, 0)

		for rows.Next() {

			select {
			case <-ctx.Done():
				break
			default:
				// pass
			}

			spr_row, err := spr.RetrieveSPRWithRows(ctx, rows)

			if err != nil {
				err_ch <- fmt.Errorf("Failed to derive SPR from row, %w", err)
				return
			}

			results = append(results, spr_row)
		}

		err = rows.Close()

		if err != nil {
			err_ch <- fmt.Errorf("Failed to close results rows for descendants, %w", err)
			return
		}

		spr_results := &spr.SQLiteResults{
			Places: results,
		}

		results_ch <- spr_results
	}()

	var pg_results pagination.Results
	var spr_results wof_spr.StandardPlacesResults

	remaining := 2

	for remaining > 0 {
		select {
		case <-done_ch:
			remaining -= 1
		case r := <-pg_ch:
			pg_results = r
		case r := <-results_ch:
			spr_results = r
		case err := <-err_ch:
			return nil, nil, err
		}
	}

	return spr_results, pg_results, nil
}

func (s *SQLSpelunker) GetDescendantsFaceted(ctx context.Context, id int64, filters []spelunker.Filter, facets []*spelunker.Facet) ([]*spelunker.Faceting, error) {

	q_where, q_args, err := s.descendantsQueryWhere(ctx, id, filters)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive query where statement, %w", err)
	}

	results := make([]*spelunker.Faceting, len(facets))

	// START OF do this in go routines

	for idx, f := range facets {

		q := s.descendantsQueryFacetStatement(ctx, f, q_where)

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

func (s *SQLSpelunker) CountDescendants(ctx context.Context, id int64) (int64, error) {

	var count int64

	q := fmt.Sprintf("SELECT COUNT(id) FROM %s WHERE ancestor_id = ? AND id != ?", tables.ANCESTORS_TABLE_NAME)
	row := s.db.QueryRowContext(ctx, q, id, id)

	err := row.Scan(&count)

	switch {
	case err == db_sql.ErrNoRows:
		return 0, spelunker.ErrNotFound
	case err != nil:
		return 0, fmt.Errorf("Failed to execute count descendants query for %d, %w", id, err)
	default:
		return count, nil
	}
}

func (s *SQLSpelunker) descendantsQueryWhere(ctx context.Context, id int64, filters []spelunker.Filter) ([]string, []interface{}, error) {

	where := []string{
		fmt.Sprintf("%s.ancestor_id = ?", tables.ANCESTORS_TABLE_NAME),
	}

	args := []interface{}{
		id,
	}

	where, args, err := s.assignFilters(where, args, filters)

	if err != nil {
		return nil, nil, err
	}

	return where, args, nil
}

func (s *SQLSpelunker) descendantsQueryColumnsAll(ctx context.Context) []string {

	// START OF put me in a function
	str_cols := `id, parent_id, name, placetype,
		inception, cessation,
		country, repo,
		latitude, longitude,
		min_latitude, min_longitude,
		max_latitude, max_longitude,
		is_current, is_deprecated, is_ceased,is_superseded, is_superseding,
		supersedes, superseded_by, belongsto,
		is_alt, alt_label,
		lastmodified`

	cols := strings.Split(str_cols, ",")
	// END OF put me in a function

	count_cols := len(cols)

	fq_cols := make([]string, count_cols)

	for idx, c := range cols {
		c = strings.TrimSpace(c)
		fq_cols[idx] = fmt.Sprintf("%s.%s AS %s", tables.SPR_TABLE_NAME, c, c)
	}

	return fq_cols
}

func (s *SQLSpelunker) descendantsQueryStatement(ctx context.Context, cols []string, where []string) string {

	str_cols := strings.Join(cols, ",")
	str_where := strings.Join(where, " AND ")

	return fmt.Sprintf("SELECT %s FROM %s JOIN %s ON %s.id = %s.id AND %s", str_cols, tables.SPR_TABLE_NAME, tables.ANCESTORS_TABLE_NAME, tables.SPR_TABLE_NAME, tables.ANCESTORS_TABLE_NAME, str_where)

}

func (s *SQLSpelunker) descendantsQueryCountStatement(ctx context.Context, where []string) string {

	cols := []string{
		fmt.Sprintf("%s.id AS id", tables.SPR_TABLE_NAME),
	}

	return s.descendantsQueryStatement(ctx, cols, where)
}

func (s *SQLSpelunker) descendantsQueryFacetStatement(ctx context.Context, facet *spelunker.Facet, where []string) string {

	facet_label := s.facetLabel(facet)

	cols := []string{
		fmt.Sprintf("%s.%s AS %s", tables.SPR_TABLE_NAME, facet_label, facet),
		fmt.Sprintf("COUNT(%s.id) AS count", tables.SPR_TABLE_NAME),
	}

	q := s.descendantsQueryStatement(ctx, cols, where)
	return fmt.Sprintf("%s GROUP BY %s.%s ORDER BY count DESC", q, tables.SPR_TABLE_NAME, facet_label)
}
