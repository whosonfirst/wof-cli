package opensearch

import (
	"context"
	"fmt"
	"strings"

	"github.com/aaronland/go-pagination"
	opensearchapi "github.com/opensearch-project/opensearch-go/v4/opensearchapi"
	wof_spr "github.com/whosonfirst/go-whosonfirst-spr/v2"
	"github.com/whosonfirst/spelunker/v2"
)

// GetConcordances retrieves the list of unique concordances in a OpenSearchSpelunker index.
func (s *OpenSearchSpelunker) GetConcordances(ctx context.Context) (*spelunker.Faceting, error) {

	c_facet := spelunker.NewFacet("concordances_sources.keyword")

	facets := []*spelunker.Facet{
		c_facet,
	}

	q := s.matchAllFacetedQuery(facets)
	sz := 0

	req := &opensearchapi.SearchReq{
		Indices: []string{
			s.index,
		},
		Body: strings.NewReader(q),
		Params: opensearchapi.SearchParams{
			Size: &sz,
		},
	}

	f, err := s.facet(ctx, req, facets)

	if err != nil {
		return nil, fmt.Errorf("Failed to facet concordances, %w", err)
	}

	return f[0], nil
}

// HasConcordance retrieve the list of records with a given concordance in an OpenSearchSpelunker index.
func (s *OpenSearchSpelunker) HasConcordance(ctx context.Context, pg_opts pagination.Options, namespace string, predicate string, value any, filters []spelunker.Filter) (wof_spr.StandardPlacesResults, pagination.Results, error) {

	q := s.hasConcordanceQuery(namespace, predicate, value, filters)
	return s.searchPaginated(ctx, pg_opts, q)

	return nil, nil, spelunker.ErrNotImplemented
}

// HasConcordanceFaceted retrieves faceted properties for records with a given concordance in an OpenSearchSpelunker index.
func (s *OpenSearchSpelunker) HasConcordanceFaceted(ctx context.Context, namespace string, predicate string, value any, filters []spelunker.Filter, facets []*spelunker.Facet) ([]*spelunker.Faceting, error) {

	q := s.hasConcordanceFacetedQuery(namespace, predicate, value, filters, facets)
	sz := 0

	req := &opensearchapi.SearchReq{
		Indices: []string{
			s.index,
		},
		Body: strings.NewReader(q),
		Params: opensearchapi.SearchParams{
			Size: &sz,
		},
	}

	return s.facet(ctx, req, facets)
}
