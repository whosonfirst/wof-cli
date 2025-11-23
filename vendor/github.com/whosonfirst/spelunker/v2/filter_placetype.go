package spelunker

import (
	"context"
	"fmt"
	"net/url"

	"github.com/whosonfirst/go-whosonfirst-placetypes"
)

// PLACETYPE_FILTER_SCHEME defines the URI scheme for `PlacetypeFilter` implementation of the `Filter` interface.
const PLACETYPE_FILTER_SCHEME string = "placetype"

// PlacetypeFilter implements the `Filter` interface for filtering results by Who's On First placetype.
type PlacetypeFilter struct {
	Filter
	placetype string
}

// NewPlacetypeFilter derives a new `Filter` implementation for filtering results whose Who's On First placetype code is encoded in 'uri'.
func NewPlacetypeFilterFromString(ctx context.Context, name string) (Filter, error) {
	uri := fmt.Sprintf("%s://%s", PLACETYPE_FILTER_SCHEME, name)
	return NewPlacetypeFilter(ctx, uri)
}

// NewPlacetypeFilter derives a new `Filter` implementation for filtering results whose Who's On First placetype code is encoded in 'uri'
// which is expected to take the form of:
//
//	placetype://{WHOSONFIRST_PLACETYPE}
func NewPlacetypeFilter(ctx context.Context, uri string) (Filter, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	pt := u.Host

	if !placetypes.IsValidPlacetype(pt) {
		return nil, fmt.Errorf("Invalid placetype")
	}

	f := &PlacetypeFilter{
		placetype: pt,
	}

	return f, nil
}

// Scheme returns the value of `PLACETYPE_FILTER_SCHEME`.
func (f *PlacetypeFilter) Scheme() string {
	return PLACETYPE_FILTER_SCHEME
}

// Value returns the Who's On First placetype code (string) that results should be filtered by.
func (f *PlacetypeFilter) Value() any {
	return f.placetype
}
