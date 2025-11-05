package spelunker

import (
	"context"
	"fmt"
	"net/url"
)

const tag_FILTER_SCHEME string = "tag"

type TagFilter struct {
	Filter
	tag string
}

func NewTagFilterFromString(ctx context.Context, t string) (Filter, error) {
	uri := fmt.Sprintf("%s://%s", tag_FILTER_SCHEME, t)
	return NewTagFilter(ctx, uri)
}

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

func (f *TagFilter) Scheme() string {
	return tag_FILTER_SCHEME
}

func (f *TagFilter) Value() any {
	return f.tag
}
