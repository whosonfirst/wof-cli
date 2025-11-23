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

// GetAlternatePlacetypes retrieves the list of alternate placetype ("wof:placetype_alt") in an OpenSearchSpelunker index.
func (s *OpenSearchSpelunker) GetAlternatePlacetypes(ctx context.Context) (*spelunker.Faceting, error) {

	pt_facet := spelunker.NewFacet("placetypealt")

	facets := []*spelunker.Facet{
		pt_facet,
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
		return nil, fmt.Errorf("Failed to facet alternate placetypes, %w", err)
	}

	return f[0], nil
}

// HasAlternatePlacetypes retrieves the list of Who's On First records with a given alternate placetype ("wof:placetype_alt") in an OpenSearchSpelunker index.
func (s *OpenSearchSpelunker) HasAlternatePlacetype(ctx context.Context, pg_opts pagination.Options, pt string, filters []spelunker.Filter) (wof_spr.StandardPlacesResults, pagination.Results, error) {

	q := s.hasAlternatePlacetypeQuery(pt, filters)
	return s.searchPaginated(ctx, pg_opts, q)
}

// HasAlternatePlacetypeFaceted retrieves faceted properties for records with a given alternate placetype ("wof:placetype_alt") in an OpenSearchSpelunker index.
func (s *OpenSearchSpelunker) HasAlternatePlacetypeFaceted(ctx context.Context, pt string, filters []spelunker.Filter, facets []*spelunker.Facet) ([]*spelunker.Faceting, error) {

	q := s.hasAlternatePlacetypeFacetedQuery(pt, filters, facets)
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
