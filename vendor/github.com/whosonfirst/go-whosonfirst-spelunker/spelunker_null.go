package spelunker

import (
	"context"
	"time"

	"github.com/aaronland/go-pagination"
	"github.com/whosonfirst/go-whosonfirst-placetypes"
	"github.com/whosonfirst/go-whosonfirst-spr/v2"
	"github.com/whosonfirst/go-whosonfirst-uri"
)

// NullSpelunker implements the [Spelunker] interface but returns an `ErrNotImplemented` error for every method.
// The easiest way to think about NullSpelunker is that its a template for implementing the Spelunker interface
// for an actual working database.
type NullSpelunker struct {
	Spelunker
}

func init() {
	ctx := context.Background()
	RegisterSpelunker(ctx, "null", NewNullSpelunker)
}

func NewNullSpelunker(ctx context.Context, uri string) (Spelunker, error) {
	s := &NullSpelunker{}
	return s, nil
}

func (s *NullSpelunker) GetRecordForId(ctx context.Context, id int64, uri_args *uri.URIArgs) ([]byte, error) {
	return nil, ErrNotImplemented
}

func (s *NullSpelunker) GetSPRForId(ctx context.Context, id int64, uri_args *uri.URIArgs) (spr.StandardPlacesResult, error) {
	return nil, ErrNotImplemented
}

func (s *NullSpelunker) GetFeatureForId(ctx context.Context, id int64, uri_args *uri.URIArgs) ([]byte, error) {
	return nil, ErrNotImplemented
}

func (s *NullSpelunker) GetDescendants(ctx context.Context, pg_opts pagination.Options, id int64, filters []Filter) (spr.StandardPlacesResults, pagination.Results, error) {
	return nil, nil, ErrNotImplemented
}

func (s *NullSpelunker) GetDescendantsFaceted(ctx context.Context, id int64, filters []Filter, facets []*Facet) ([]*Faceting, error) {
	return nil, ErrNotImplemented
}

func (s *NullSpelunker) CountDescendants(ctx context.Context, id int64) (int64, error) {
	return 0, ErrNotImplemented
}

func (s *NullSpelunker) Search(ctx context.Context, pg_opts pagination.Options, q *SearchOptions, filters []Filter) (spr.StandardPlacesResults, pagination.Results, error) {
	return nil, nil, ErrNotImplemented
}

func (s *NullSpelunker) SearchFaceted(ctx context.Context, q *SearchOptions, filters []Filter, facets []*Facet) ([]*Faceting, error) {
	return nil, ErrNotImplemented
}

func (s *NullSpelunker) GetRecent(ctx context.Context, pg_opts pagination.Options, d time.Duration, filters []Filter) (spr.StandardPlacesResults, pagination.Results, error) {
	return nil, nil, ErrNotImplemented
}

func (s *NullSpelunker) GetRecentFaceted(ctx context.Context, d time.Duration, filters []Filter, facets []*Facet) ([]*Faceting, error) {
	return nil, ErrNotImplemented
}

func (s *NullSpelunker) GetPlacetypes(ctx context.Context) (*Faceting, error) {
	return nil, ErrNotImplemented
}

func (s *NullSpelunker) HasPlacetype(ctx context.Context, pg_opts pagination.Options, pt *placetypes.WOFPlacetype, filters []Filter) (spr.StandardPlacesResults, pagination.Results, error) {
	return nil, nil, ErrNotImplemented
}

func (s *NullSpelunker) HasPlacetypeFaceted(ctx context.Context, pt *placetypes.WOFPlacetype, filters []Filter, facets []*Facet) ([]*Faceting, error) {
	return nil, ErrNotImplemented
}

func (s *NullSpelunker) GetConcordances(ctx context.Context) (*Faceting, error) {
	return nil, ErrNotImplemented
}

func (s *NullSpelunker) HasConcordance(ctx context.Context, pg_opts pagination.Options, namespace string, predicate string, value any, filters []Filter) (spr.StandardPlacesResults, pagination.Results, error) {
	return nil, nil, ErrNotImplemented
}

func (s *NullSpelunker) HasConcordanceFaceted(ctx context.Context, namespace string, predicate string, value any, filters []Filter, facets []*Facet) ([]*Faceting, error) {
	return nil, ErrNotImplemented
}

func (s *NullSpelunker) GetTags(ctx context.Context) (*Faceting, error) {
	return nil, ErrNotImplemented
}

func (s *NullSpelunker) HasTag(ctx context.Context, pg_opts pagination.Options, tag string, filters []Filter) (spr.StandardPlacesResults, pagination.Results, error) {
	return nil, nil, ErrNotImplemented
}

func (s *NullSpelunker) HasTagFaceted(ctx context.Context, tag string, filters []Filter, facets []*Facet) ([]*Faceting, error) {
	return nil, ErrNotImplemented
}

func (s *NullSpelunker) VisitingNullIsland(ctx context.Context, pg_opts pagination.Options, filters []Filter) (spr.StandardPlacesResults, pagination.Results, error) {
	return nil, nil, ErrNotImplemented
}

func (s *NullSpelunker) VisitingNullIslandFaceted(ctx context.Context, filters []Filter, facets []*Facet) ([]*Faceting, error) {
	return nil, ErrNotImplemented
}
