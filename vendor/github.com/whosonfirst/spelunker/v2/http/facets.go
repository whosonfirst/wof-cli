package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aaronland/go-http/v4/sanitize"
	"github.com/whosonfirst/spelunker/v2"
)

// FiltersFromRequest derives faceting criteria from 'req' for query parameters matching 'params'.
func FacetsFromRequest(ctx context.Context, req *http.Request, params []string) ([]*spelunker.Facet, error) {

	// TBD...
	facets := make([]*spelunker.Facet, 0)

	v, err := sanitize.GetString(req, "facet")

	if err != nil {
		return nil, fmt.Errorf("Failed to derive ?facet= query  parameter, %w", err)
	}

	if v == "" {
		return nil, fmt.Errorf("Empty facet paramter")
	}

	facets = append(facets, spelunker.NewFacet(v))
	return facets, nil
}
