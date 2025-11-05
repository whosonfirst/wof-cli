package spelunker

import (
	"context"
	"fmt"
	"net/url"
)

const COUNTRY_FILTER_SCHEME string = "country"

type CountryFilter struct {
	Filter
	code string
}

func NewCountryFilterFromString(ctx context.Context, code string) (Filter, error) {
	uri := fmt.Sprintf("%s://%s", COUNTRY_FILTER_SCHEME, code)
	return NewCountryFilter(ctx, uri)
}

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

func (f *CountryFilter) Scheme() string {
	return COUNTRY_FILTER_SCHEME
}

func (f *CountryFilter) Value() any {
	return f.code
}
