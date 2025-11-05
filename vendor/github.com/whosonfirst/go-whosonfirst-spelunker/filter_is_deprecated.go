package spelunker

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

const IS_DEPRECATED_FILTER_SCHEME string = "isdeprecated"

type IsDeprecatedFilter struct {
	Filter
	is_deprecated int
}

func NewIsDeprecatedFilterFromString(ctx context.Context, name string) (Filter, error) {
	uri := fmt.Sprintf("%s://?flag=%s", IS_DEPRECATED_FILTER_SCHEME, name)
	return NewIsDeprecatedFilter(ctx, uri)
}

func NewIsDeprecatedFilter(ctx context.Context, uri string) (Filter, error) {

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

	f := &IsDeprecatedFilter{
		is_deprecated: fl,
	}

	return f, nil
}

func (f *IsDeprecatedFilter) Scheme() string {
	return IS_DEPRECATED_FILTER_SCHEME
}

func (f *IsDeprecatedFilter) Value() any {
	return f.is_deprecated
}
