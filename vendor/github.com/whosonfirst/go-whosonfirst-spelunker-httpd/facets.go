package httpd

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aaronland/go-http-sanitize"
	"github.com/whosonfirst/go-whosonfirst-spelunker"
)

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
