package spelunker

import (
	"context"
	"fmt"
	"net/url"

	"github.com/whosonfirst/go-whosonfirst-placetypes"
)

const PLACETYPE_FILTER_SCHEME string = "placetype"

type PlacetypeFilter struct {
	Filter
	placetype string
}

func NewPlacetypeFilterFromString(ctx context.Context, name string) (Filter, error) {
	uri := fmt.Sprintf("%s://%s", PLACETYPE_FILTER_SCHEME, name)
	return NewPlacetypeFilter(ctx, uri)
}

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

func (f *PlacetypeFilter) Scheme() string {
	return PLACETYPE_FILTER_SCHEME
}

func (f *PlacetypeFilter) Value() any {
	return f.placetype
}
