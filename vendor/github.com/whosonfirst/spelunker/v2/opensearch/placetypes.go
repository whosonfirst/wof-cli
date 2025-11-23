package opensearch

import (
	"context"
	"fmt"
	"strings"

	"github.com/aaronland/go-pagination"
	opensearchapi "github.com/opensearch-project/opensearch-go/v4/opensearchapi"
	"github.com/whosonfirst/go-whosonfirst-placetypes"
	wof_spr "github.com/whosonfirst/go-whosonfirst-spr/v2"
	"github.com/whosonfirst/spelunker/v2"
)

// GetPlacetypes retrieves the list of unique placetypes in a Spleunker index in an OpenSearchSpelunker index.
func (s *OpenSearchSpelunker) GetPlacetypes(ctx context.Context) (*spelunker.Faceting, error) {

	pt_facet := spelunker.NewFacet("placetype")

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
		return nil, fmt.Errorf("Failed to facet placetypes, %w", err)
	}

	return f[0], nil
}

// HasPlacetype retrieves the list of records with a given placetype in an OpenSearchSpelunker index.
func (s *OpenSearchSpelunker) HasPlacetype(ctx context.Context, pg_opts pagination.Options, pt *placetypes.WOFPlacetype, filters []spelunker.Filter) (wof_spr.StandardPlacesResults, pagination.Results, error) {

	q := s.hasPlacetypeQuery(pt.Name, filters)
	return s.searchPaginated(ctx, pg_opts, q)
}

// HasPlacetypeFaceted retrieves faceted properties for records with a given placetype in an OpenSearchSpelunker index.
func (s *OpenSearchSpelunker) HasPlacetypeFaceted(ctx context.Context, pt *placetypes.WOFPlacetype, filters []spelunker.Filter, facets []*spelunker.Facet) ([]*spelunker.Faceting, error) {

	q := s.hasPlacetypeFacetedQuery(pt.Name, filters, facets)
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
