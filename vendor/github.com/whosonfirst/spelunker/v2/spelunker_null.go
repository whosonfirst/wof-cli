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

// NewNullSpelunker returns an implementation of the [Spelunker] interface that returns an `ErrNotImplemented` error
// for every method configured by 'uri' which is expected to take the form of:
//
//	null://
func NewNullSpelunker(ctx context.Context, uri string) (Spelunker, error) {
	s := &NullSpelunker{}
	return s, nil
}

// GetRecordForId retrieves properties (or more specifically the "document") for a given ID in a NullSpelunker database.
func (s *NullSpelunker) GetRecordForId(ctx context.Context, id int64, uri_args *uri.URIArgs) ([]byte, error) {
	return nil, ErrNotImplemented
}

// GetSPRForId retrieves the `spr.StandardPlaceResult` instance for a given ID in a NullSpelunker database.
func (s *NullSpelunker) GetSPRForId(ctx context.Context, id int64, uri_args *uri.URIArgs) (spr.StandardPlacesResult, error) {
	return nil, ErrNotImplemented
}

// GetFeatureForId retrieves the GeoJSON Feature record for a given ID in a NullSpelunker database.
func (s *NullSpelunker) GetFeatureForId(ctx context.Context, id int64, uri_args *uri.URIArgs) ([]byte, error) {
	return nil, ErrNotImplemented
}

// GetDescendants retrieves all the Who's On First record that are a descendant of a specific Who's On First ID in a NullSpelunker database.
func (s *NullSpelunker) GetDescendants(ctx context.Context, pg_opts pagination.Options, id int64, filters []Filter) (spr.StandardPlacesResults, pagination.Results, error) {
	return nil, nil, ErrNotImplemented
}

// GetDescendantsFaceted retrieves faceted properties for records that are a descendant of a specific Who's On First ID in a NullSpelunker database.
func (s *NullSpelunker) GetDescendantsFaceted(ctx context.Context, id int64, filters []Filter, facets []*Facet) ([]*Faceting, error) {
	return nil, ErrNotImplemented
}

// CountDescendants returns the total number of Who's On First records that are a descendant of a specific Who's On First ID in a NullSpelunker database.
func (s *NullSpelunker) CountDescendants(ctx context.Context, id int64) (int64, error) {
	return 0, ErrNotImplemented
}

// Search retrieves all the Who's On First records that match a search criteria in a NullSpelunker database.
func (s *NullSpelunker) Search(ctx context.Context, pg_opts pagination.Options, q *SearchOptions, filters []Filter) (spr.StandardPlacesResults, pagination.Results, error) {
	return nil, nil, ErrNotImplemented
}

// SearchFaceted retrieves faceted properties for records match a search criteria in a NullSpelunker database.
func (s *NullSpelunker) SearchFaceted(ctx context.Context, q *SearchOptions, filters []Filter, facets []*Facet) ([]*Faceting, error) {
	return nil, ErrNotImplemented
}

// GetRecent retrieves all the Who's On First records that have been modified with a window of time in a NullSpelunker database.
func (s *NullSpelunker) GetRecent(ctx context.Context, pg_opts pagination.Options, d time.Duration, filters []Filter) (spr.StandardPlacesResults, pagination.Results, error) {
	return nil, nil, ErrNotImplemented
}

// GetRecentFaceted retrieves faceted properties for records that have been modified with a window of time in a NullSpelunker database.
func (s *NullSpelunker) GetRecentFaceted(ctx context.Context, d time.Duration, filters []Filter, facets []*Facet) ([]*Faceting, error) {
	return nil, ErrNotImplemented
}

// GetPlacetypes retrieves the list of unique placetypes in a Spleunker index in a NullSpelunker database.
func (s *NullSpelunker) GetPlacetypes(ctx context.Context) (*Faceting, error) {
	return nil, ErrNotImplemented
}

// HasPlacetype retrieves the list of records with a given placetype in a NullSpelunker database.
func (s *NullSpelunker) HasPlacetype(ctx context.Context, pg_opts pagination.Options, pt *placetypes.WOFPlacetype, filters []Filter) (spr.StandardPlacesResults, pagination.Results, error) {
	return nil, nil, ErrNotImplemented
}

// HasPlacetypeFaceted retrieves faceted properties for records with a given placetype in a NullSpelunker database.
func (s *NullSpelunker) HasPlacetypeFaceted(ctx context.Context, pt *placetypes.WOFPlacetype, filters []Filter, facets []*Facet) ([]*Faceting, error) {
	return nil, ErrNotImplemented
}

// GetAlternatePlacetypes retrieves the list of alternate placetype ("wof:placetype_alt") in a NullSpelunker database.
func (s *NullSpelunker) GetAlternatePlacetypes(ctx context.Context) (*Faceting, error) {
	return nil, ErrNotImplemented
}

// HasAlternatePlacetypes retrieves the list of Who's On First records with a given alternate placetype ("wof:placetype_alt") in a NullSpelunker database.
func (s *NullSpelunker) HasAlternatePlacetype(ctx context.Context, pg_opts pagination.Options, pt string, filters []Filter) (spr.StandardPlacesResults, pagination.Results, error) {
	return nil, nil, ErrNotImplemented
}

// HasAlternatePlacetypeFaceted retrieves faceted properties for records with a given alternate placetype ("wof:placetype_alt") in a NullSpelunker database.
func (s *NullSpelunker) HasAlternatePlacetypeFaceted(ctx context.Context, pt string, filters []Filter, facets []*Facet) ([]*Faceting, error) {
	return nil, ErrNotImplemented
}

// GetConcordances retrieves the list of unique concordances in a NullSpelunker database.
func (s *NullSpelunker) GetConcordances(ctx context.Context) (*Faceting, error) {
	return nil, ErrNotImplemented
}

// HasConcordance retrieve the list of records with a given concordance in a NullSpelunker database.
func (s *NullSpelunker) HasConcordance(ctx context.Context, pg_opts pagination.Options, namespace string, predicate string, value any, filters []Filter) (spr.StandardPlacesResults, pagination.Results, error) {
	return nil, nil, ErrNotImplemented
}

// HasConcordanceFaceted retrieves faceted properties for records with a given concordance in a NullSpelunker database.
func (s *NullSpelunker) HasConcordanceFaceted(ctx context.Context, namespace string, predicate string, value any, filters []Filter, facets []*Facet) ([]*Faceting, error) {
	return nil, ErrNotImplemented
}

// GetTags retrieves the list of unique tags in a Spelunker index in a NullSpelunker database.
func (s *NullSpelunker) GetTags(ctx context.Context) (*Faceting, error) {
	return nil, ErrNotImplemented
}

// HasTag retrieves the list of records that have a given tag in a NullSpelunker database.
func (s *NullSpelunker) HasTag(ctx context.Context, pg_opts pagination.Options, tag string, filters []Filter) (spr.StandardPlacesResults, pagination.Results, error) {
	return nil, nil, ErrNotImplemented
}

// HasTagFaceted retrieves faceted properties for records that have a given tag in a NullSpelunker database.
func (s *NullSpelunker) HasTagFaceted(ctx context.Context, tag string, filters []Filter, facets []*Facet) ([]*Faceting, error) {
	return nil, ErrNotImplemented
}

// VisitingNullIsland retrieves the list of records that are "visiting Null Island" (have a latitude, longitude value of "0.0, 0.0" in a NullSpelunker database.
func (s *NullSpelunker) VisitingNullIsland(ctx context.Context, pg_opts pagination.Options, filters []Filter) (spr.StandardPlacesResults, pagination.Results, error) {
	return nil, nil, ErrNotImplemented
}

// VisitingNullIslandFaceted retrieves faceted properties for records that are "visiting Null Island" (have a latitude, longitude value of "0.0, 0.0" in a NullSpelunker database.
func (s *NullSpelunker) VisitingNullIslandFaceted(ctx context.Context, filters []Filter, facets []*Facet) ([]*Faceting, error) {
	return nil, ErrNotImplemented
}
