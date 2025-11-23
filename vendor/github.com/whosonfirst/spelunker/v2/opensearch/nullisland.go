package opensearch

import (
	"context"
	"strings"

	"github.com/aaronland/go-pagination"
	opensearchapi "github.com/opensearch-project/opensearch-go/v4/opensearchapi"
	wof_spr "github.com/whosonfirst/go-whosonfirst-spr/v2"
	"github.com/whosonfirst/spelunker/v2"
)

// VisitingNullIsland retrieves the list of records that are "visiting Null Island" (have a latitude, longitude value of "0.0, 0.0" in an OpenSearchSpelunker index.
func (s *OpenSearchSpelunker) VisitingNullIsland(ctx context.Context, pg_opts pagination.Options, filters []spelunker.Filter) (wof_spr.StandardPlacesResults, pagination.Results, error) {

	q := s.visitingNullIslandQuery(filters)
	return s.searchPaginated(ctx, pg_opts, q)
}

// VisitingNullIslandFaceted retrieves faceted properties for records that are "visiting Null Island" (have a latitude, longitude value of "0.0, 0.0" in an OpenSearchSpelunker index.
func (s *OpenSearchSpelunker) VisitingNullIslandFaceted(ctx context.Context, filters []spelunker.Filter, facets []*spelunker.Facet) ([]*spelunker.Faceting, error) {

	q := s.visitingNullIslandFacetedQuery(filters, facets)
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
