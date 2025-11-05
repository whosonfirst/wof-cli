package sql

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"strconv"
	"strings"

	"github.com/aaronland/go-pagination"
	"github.com/aaronland/go-pagination/countable"
	"github.com/whosonfirst/go-whosonfirst-spelunker"
	wof_spr "github.com/whosonfirst/go-whosonfirst-spr/v2"
	"github.com/whosonfirst/go-whosonfirst-sql/tables"
	"github.com/whosonfirst/go-whosonfirst-sqlite-spr"
)

func (s *SQLSpelunker) queryCount(ctx context.Context, col string, q string, args ...interface{}) (int64, error) {

	parts := strings.Split(q, " FROM ")
	parts = strings.Split(parts[1], " LIMIT ")
	parts = strings.Split(parts[0], " ORDER ")

	conditions := parts[0]

	count_query := fmt.Sprintf("SELECT COUNT(%s) FROM %s", col, conditions)

	row := s.db.QueryRowContext(ctx, count_query, args...)

	var count int64
	err := row.Scan(&count)

	if err != nil {
		return 0, fmt.Errorf("Failed to execute count query '%s', %w", count_query, err)
	}

	return count, nil
}

func (s *SQLSpelunker) deriveLimitOffset(pg_opts pagination.Options) (int, int) {

	page_num := countable.PageFromOptions(pg_opts)
	page := int(math.Max(1.0, float64(page_num)))

	per_page := int(math.Max(1.0, float64(pg_opts.PerPage())))
	spill := int(math.Max(1.0, float64(pg_opts.Spill())))

	if spill >= per_page {
		spill = per_page - 1
	}

	offset := 0
	limit := per_page

	offset = (page - 1) * per_page

	return limit, offset
}

func (s *SQLSpelunker) selectSPR(ctx context.Context, where string) string {
	return fmt.Sprintf(`SELECT 
		id, parent_id, name, placetype,
		inception, cessation,
		country, repo,
		latitude, longitude,
		min_latitude, min_longitude,
		max_latitude, max_longitude,
		is_current, is_deprecated, is_ceased,is_superseded, is_superseding,
		supersedes, superseded_by, belongsto,
		is_alt, alt_label,
		lastmodified
	FROM %s WHERE %s`, tables.SPR_TABLE_NAME, where)
}

func (s *SQLSpelunker) querySPR(ctx context.Context, pg_opts pagination.Options, where string, args ...interface{}) (wof_spr.StandardPlacesResults, pagination.Results, error) {

	if pg_opts != nil {
		limit, offset := s.deriveLimitOffset(pg_opts)
		where = fmt.Sprintf("%s LIMIT %d OFFSET %d", where, limit, offset)
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

		count_q := fmt.Sprintf("SELECT %s.id AS id FROM %s WHERE %s", tables.SPR_TABLE_NAME, tables.SPR_TABLE_NAME, where)
		count, err := s.queryCount(ctx, "id", count_q, args...)

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

		results_q := s.selectSPR(ctx, where)

		rows, err := s.db.QueryContext(ctx, results_q, args...)

		if err != nil {
			err_ch <- fmt.Errorf("Failed to query where '%s', %w", results_q, err)
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

func (s *SQLSpelunker) querySearch(ctx context.Context, pg_opts pagination.Options, where string, args ...interface{}) (wof_spr.StandardPlacesResults, pagination.Results, error) {

	q := fmt.Sprintf("SELECT id FROM %s WHERE %s", tables.SEARCH_TABLE_NAME, where)

	return s.querySearchDo(ctx, pg_opts, q, args...)
}

func (s *SQLSpelunker) querySearchWithFilters(ctx context.Context, pg_opts pagination.Options, where string, args ...interface{}) (wof_spr.StandardPlacesResults, pagination.Results, error) {

	// Note the SQLite specific-iness of this.
	// TO DO: Add code to "do the right thing" depending on a.engine (MySQL, etc...)
	// SELECT search.id AS id FROM search JOIN spr ON (search.id = CAST(spr.id AS INTEGER) AND search.names_all MATCH 'montreal') LIMIT 10 OFFSET 0;

	q := fmt.Sprintf("SELECT %s.id AS id FROM %s JOIN %s ON %s.id = CAST(%s.id AS INTEGER) WHERE %s",
		tables.SEARCH_TABLE_NAME,
		tables.SEARCH_TABLE_NAME,
		tables.SPR_TABLE_NAME,
		tables.SEARCH_TABLE_NAME,
		tables.SPR_TABLE_NAME,
		where,
	)

	return s.querySearchDo(ctx, pg_opts, q, args...)
}

func (s *SQLSpelunker) querySearchDo(ctx context.Context, pg_opts pagination.Options, q string, args ...interface{}) (wof_spr.StandardPlacesResults, pagination.Results, error) {

	slog.Info("SEARCH", "q", q)

	// https://www.sqlite.org/fts5.html

	if pg_opts != nil {
		limit, offset := s.deriveLimitOffset(pg_opts)
		q = fmt.Sprintf("%s LIMIT %d OFFSET %d", q, limit, offset)
	}

	pg_ch := make(chan pagination.Results)
	id_ch := make(chan int64)

	done_ch := make(chan bool)
	err_ch := make(chan error)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {

		defer func() {
			done_ch <- true
		}()

		count, err := s.queryCount(ctx, fmt.Sprintf("%s.id", tables.SEARCH_TABLE_NAME), q, args...)

		if err != nil {
			err_ch <- fmt.Errorf("Failed to derive query count, %w", err)
			return
		}

		slog.Info("SEARCH", "query", q, "count", count)

		pg_results, err := countable.NewResultsFromCountWithOptions(pg_opts, count)

		if err != nil {
			err_ch <- fmt.Errorf("Failed to derive pagination results, %w", err)
			return
		}

		pg_ch <- pg_results
	}()

	go func() {

		defer func() {
			done_ch <- true
		}()

		slog.Info("SEARCH", "do", q, "args", args)

		rows, err := s.db.QueryContext(ctx, q, args...)

		if err != nil {
			err_ch <- fmt.Errorf("Failed to query where '%s', %w", q, err)
			return
		}

		for rows.Next() {

			select {
			case <-ctx.Done():
				break
			default:
				// pass
			}

			var id int64
			err := rows.Scan(&id)

			if err != nil {
				err_ch <- fmt.Errorf("Failed to scan ID, %w", err)
				return
			}

			id_ch <- id
		}

		err = rows.Close()

		if err != nil {
			err_ch <- fmt.Errorf("Failed to close results rows for descendants, %w", err)
			return
		}
	}()

	var pg_results pagination.Results
	str_ids := make([]string, 0)

	remaining := 2

	for remaining > 0 {
		select {
		case <-done_ch:
			remaining -= 1
		case r := <-pg_ch:
			pg_results = r
		case id := <-id_ch:
			str_id := strconv.FormatInt(id, 10)
			str_ids = append(str_ids, str_id)
		case err := <-err_ch:
			return nil, nil, err
		}
	}

	spr_where := fmt.Sprintf("id IN (%s)", strings.Join(str_ids, ","))

	spr_results, _, err := s.querySPR(ctx, nil, spr_where)

	if err != nil {
		return nil, nil, fmt.Errorf("Failed to retrieve SPR records, %w", err)
	}

	return spr_results, pg_results, nil
}

func (s *SQLSpelunker) assignFilters(where []string, args []interface{}, filters []spelunker.Filter) ([]string, []interface{}, error) {

	for _, f := range filters {

		switch f.Scheme() {
		case spelunker.COUNTRY_FILTER_SCHEME:
			where = append(where, fmt.Sprintf("%s.country = ?", tables.SPR_TABLE_NAME))
			args = append(args, f.Value())
		case spelunker.PLACETYPE_FILTER_SCHEME:
			where = append(where, fmt.Sprintf("%s.placetype = ?", tables.SPR_TABLE_NAME))
			args = append(args, f.Value())
		case spelunker.IS_CURRENT_FILTER_SCHEME:
			where = append(where, fmt.Sprintf("%s.is_current = ?", tables.SPR_TABLE_NAME))
			args = append(args, f.Value())
		case spelunker.IS_DEPRECATED_FILTER_SCHEME:
			switch f.Value().(int) {
			case 0:
				where = append(where, fmt.Sprintf("%s.is_deprecated != 1", tables.SPR_TABLE_NAME))
			default:
				where = append(where, fmt.Sprintf("%s.is_deprecated = 1", tables.SPR_TABLE_NAME))
			}
		default:
			return nil, nil, fmt.Errorf("Invalid or unsupported filter scheme, %s", f.Scheme())
		}
	}

	return where, args, nil
}
