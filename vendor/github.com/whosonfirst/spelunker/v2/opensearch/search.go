package opensearch

import (
	"context"
	_ "fmt"
	"strings"

	"github.com/aaronland/go-pagination"
	opensearchapi "github.com/opensearch-project/opensearch-go/v4/opensearchapi"
	wof_spr "github.com/whosonfirst/go-whosonfirst-spr/v2"
	"github.com/whosonfirst/spelunker/v2"
)

// Search retrieves all the Who's On First records that match a search criteria in an OpenSearchSpelunker index.
func (s *OpenSearchSpelunker) Search(ctx context.Context, pg_opts pagination.Options, search_opts *spelunker.SearchOptions, filters []spelunker.Filter) (wof_spr.StandardPlacesResults, pagination.Results, error) {

	q := s.searchQuery(search_opts, filters)
	return s.searchPaginated(ctx, pg_opts, q)
}

// SearchFaceted retrieves faceted properties for records match a search criteria in an OpenSearchSpelunker index.
func (s *OpenSearchSpelunker) SearchFaceted(ctx context.Context, search_opts *spelunker.SearchOptions, filters []spelunker.Filter, facets []*spelunker.Facet) ([]*spelunker.Faceting, error) {

	q := s.searchFacetedQuery(search_opts, filters, facets)
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
