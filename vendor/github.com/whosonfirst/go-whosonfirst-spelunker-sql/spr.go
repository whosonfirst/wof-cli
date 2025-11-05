package sql

// Common code for querying the `spr` table.

import (
	"context"
	"fmt"
	"strings"

	"github.com/whosonfirst/go-whosonfirst-spelunker"
	"github.com/whosonfirst/go-whosonfirst-sql/tables"
)

func (s *SQLSpelunker) facetSPR(ctx context.Context, facet *spelunker.Facet, where string, args ...interface{}) ([]*spelunker.FacetCount, error) {

	q := fmt.Sprintf("SELECT %s, COUNT(id) AS count FROM %s WHERE %s GROUP BY %s ORDER BY count DESC", facet, tables.SPR_TABLE_NAME, where, facet)

	return s.facetWithQuery(ctx, q, args...)
}

func (s *SQLSpelunker) sprQueryColumnsAll(ctx context.Context) []string {

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
