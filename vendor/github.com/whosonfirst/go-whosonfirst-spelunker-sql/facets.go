package sql

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/whosonfirst/go-whosonfirst-spelunker"
)

func (s *SQLSpelunker) facetLabel(f *spelunker.Facet) string {

	var facet_label string

	switch f.Property {
	case "iscurrent":
		facet_label = "is_current"
	case "isdeprecated":
		facet_label = "is_deprecated"
	default:
		facet_label = f.Property
	}

	return facet_label
}

func (s *SQLSpelunker) facetWithQuery(ctx context.Context, q string, args ...interface{}) ([]*spelunker.FacetCount, error) {

	rows, err := s.db.QueryContext(ctx, q, args...)

	if err != nil {
		slog.Error("Failed to query facets", "query", q, "args", args, "error", err)
		return nil, fmt.Errorf("Failed to query facets, %w", err)
	}

	counts := make([]*spelunker.FacetCount, 0)

	for rows.Next() {

		var facet string
		var count int64

		err := rows.Scan(&facet, &count)

		if err != nil {
			slog.Error("Failed to scan facet columns", "query", q, "args", args, "error", err)
			return nil, fmt.Errorf("Failed to scan facet columns, %w", err)
		}

		f := &spelunker.FacetCount{
			Key:   facet,
			Count: count,
		}

		counts = append(counts, f)
	}

	err = rows.Close()

	if err != nil {
		return nil, fmt.Errorf("Failed to close results rows for facets, %w", err)
	}

	return counts, nil
}
