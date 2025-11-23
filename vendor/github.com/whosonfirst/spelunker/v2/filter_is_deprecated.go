package spelunker

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// IS_DEPRECATED_FILTER_SCHEME defines the URI scheme for `IsDeprecatedFilter` implementation of the `Filter` interface.
const IS_DEPRECATED_FILTER_SCHEME string = "isdeprecated"

// IsDeprecatedFilter implements the `Filter` interface for filtering results whose deprecation status	is true, false or unknown.
type IsDeprecatedFilter struct {
	Filter
	is_deprecated int
}

// NewIsDeprecatedFilterFromString derives a new `Filter` implementation for filtering results whose deprecation status is true, false or unknown.
func NewIsDeprecatedFilterFromString(ctx context.Context, name string) (Filter, error) {
	uri := fmt.Sprintf("%s://?flag=%s", IS_DEPRECATED_FILTER_SCHEME, name)
	return NewIsDeprecatedFilter(ctx, uri)
}

// NewIsDeprecatedFilter derives a new `Filter` implementation for filtering results whose deprecation status is encoded in 'uri'
// which is expected to take the form of:
//
//	isdeprecated://?flag={VALUE}
//
// Where {VALUE} may be 1 (true), 0 (false) or -1 (unknown)
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

// Scheme returns the value of `IS_DEPRECATED_FILTER_SCHEME`.
func (f *IsDeprecatedFilter) Scheme() string {
	return IS_DEPRECATED_FILTER_SCHEME
}

// Value returns the deprecation status (int) that results should be filtered by.
func (f *IsDeprecatedFilter) Value() any {
	return f.is_deprecated
}
