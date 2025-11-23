package spelunker

import (
	"context"
	"fmt"
	"net/url"
)

// COUNTRY_FILTER_SCHEME defines the URI scheme for `CountryFilter` implementation of the `Filter` interface.
const COUNTRY_FILTER_SCHEME string = "country"

// CountryFilter implements the `Filter` interface for filtering results by country.
type CountryFilter struct {
	Filter
	code string
}

// NewCountryFilterFromString derives a new `Filter` implementation for filtering results whose ISO country code is 'code'.
func NewCountryFilterFromString(ctx context.Context, code string) (Filter, error) {
	uri := fmt.Sprintf("%s://%s", COUNTRY_FILTER_SCHEME, code)
	return NewCountryFilter(ctx, uri)
}

// NewCountryFilter derives a new `Filter` implementation for filtering results whose ISO country code is encoded in 'uri'
// which is expected to take the form of:
//
//	country://{ISO_COUNTRY_CODE}
func NewCountryFilter(ctx context.Context, uri string) (Filter, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	code := u.Host

	// Validate code here...

	f := &CountryFilter{
		code: code,
	}

	return f, nil
}

// Scheme returns the value of `COUNTRY_FILTER_SCHEME`.
func (f *CountryFilter) Scheme() string {
	return COUNTRY_FILTER_SCHEME
}

// Value returns the ISO country code (string) that results should be filtered by.
func (f *CountryFilter) Value() any {
	return f.code
}
