package spelunker

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// IS_CURRENT_FILTER_SCHEME defines the URI scheme for `IsCurrentFilter` implementation of the `Filter` interface.
const IS_CURRENT_FILTER_SCHEME string = "iscurrent"

// IsCurrentFilter implements the `Filter` interface for filtering results by their `mz:is_current` property.
type IsCurrentFilter struct {
	Filter
	is_current int
}

// NewIsCurrentFilterFromString derives a new `Filter` implementation for filtering results whose `mz:is_current` value matches the numeric value of 'name'.
func NewIsCurrentFilterFromString(ctx context.Context, name string) (Filter, error) {
	uri := fmt.Sprintf("%s://?flag=%s", IS_CURRENT_FILTER_SCHEME, name)
	return NewIsCurrentFilter(ctx, uri)
}

// NewIsCurrentFilter derives a new `Filter` implementation for filtering results whose `mz:is_current` value is encoded in 'uri'
// which is expected to take the form of:
//
//	iscurrent://?flag={VALUE}
//
// Where {VALUE} may be 1 (true), 0 (false) or -1 (unknown)
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

	switch fl {
	case -1, 0, 1:
		// pass
	default:
		return nil, fmt.Errorf("Invalid is current value")
	}

	f := &IsCurrentFilter{
		is_current: fl,
	}

	return f, nil
}

// Scheme returns the value of `IS_CURRENT_FILTER_SCHEME`.
func (f *IsCurrentFilter) Scheme() string {
	return IS_CURRENT_FILTER_SCHEME
}

// Value returns the `mz:is_current` value (int) that results should be filtered by.
func (f *IsCurrentFilter) Value() any {
	return f.is_current
}
