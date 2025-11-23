package spelunker

import (
	"context"
	"fmt"
	"net/url"
)

// TAG_FILTER_SCHEME defines the URI scheme for `TagFilter` implementation of the `Filter` interface.
const tag_FILTER_SCHEME string = "tag"

// TagFilter implements the `Filter` interface for filtering results by tag.
type TagFilter struct {
	Filter
	tag string
}

// NewTagFilterFromString derives a new `Filter` implementation for filtering results whose tag value 't'.
func NewTagFilterFromString(ctx context.Context, t string) (Filter, error) {
	uri := fmt.Sprintf("%s://%s", tag_FILTER_SCHEME, t)
	return NewTagFilter(ctx, uri)
}

// NewTagFilter derives a new `Filter` implementation for filtering results whose ISO tag code is encoded in 'uri'
// which is expected to take the form of:
//
//	tag://{TAG}
func NewTagFilter(ctx context.Context, uri string) (Filter, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	t := u.Host

	f := &TagFilter{
		tag: t,
	}

	return f, nil
}

// Scheme returns the value of `TAG_FILTER_SCHEME`.
func (f *TagFilter) Scheme() string {
	return tag_FILTER_SCHEME
}

// Value returns the tag value (string) that results should be filtered by.
func (f *TagFilter) Value() any {
	return f.tag
}
