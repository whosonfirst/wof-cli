package spelunker

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

const IS_CURRENT_FILTER_SCHEME string = "iscurrent"

type IsCurrentFilter struct {
	Filter
	is_current int
}

func NewIsCurrentFilterFromString(ctx context.Context, name string) (Filter, error) {
	uri := fmt.Sprintf("%s://?flag=%s", IS_CURRENT_FILTER_SCHEME, name)
	return NewIsCurrentFilter(ctx, uri)
}

func NewIsCurrentFilter(ctx context.Context, uri string) (Filter, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	q := u.Query()

	if !q.Has("flag") {
		return nil, fmt.Errorf("Missing ?flag= parameter")
	}

	str_fl := q.Get("flag")

	fl, err := strconv.Atoi(str_fl)

	if err != nil {
		return nil, fmt.Errorf("Invalid ?flag= parameter, %w", err)
	}

	f := &IsCurrentFilter{
		is_current: fl,
	}

	return f, nil
}

func (f *IsCurrentFilter) Scheme() string {
	return IS_CURRENT_FILTER_SCHEME
}

func (f *IsCurrentFilter) Value() any {
	return f.is_current
}
