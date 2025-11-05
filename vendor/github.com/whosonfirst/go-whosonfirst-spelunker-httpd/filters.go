package httpd

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aaronland/go-http-sanitize"
	"github.com/whosonfirst/go-whosonfirst-spelunker"
)

func DefaultFilterParams() []string {

	return []string{
		"placetype",
		"country",
		"tag",
		"iscurrent",
		"isdeprecated",
	}
}

func FiltersFromRequest(ctx context.Context, req *http.Request, params []string) ([]spelunker.Filter, error) {

	filters := make([]spelunker.Filter, 0)

	for _, p := range params {

		switch p {
		case "iscurrent":

			str_fl, err := sanitize.GetString(req, "iscurrent")

			if err != nil {
				return nil, fmt.Errorf("Failed to derive ?is_current= query parameter, %w", err)
			}

			if str_fl != "" {

				switch str_fl {
				case "-1", "0", "1":
					// ok
				default:
					return nil, fmt.Errorf("Invalid ?iscurrent= query parameter")
				}

				is_current_f, err := spelunker.NewIsCurrentFilterFromString(ctx, str_fl)

				if err != nil {
					return nil, fmt.Errorf("Failed to create new is current filter, %w", err)
				}

				filters = append(filters, is_current_f)
			}

		case "isdeprecated":

			str_fl, err := sanitize.GetString(req, "isdeprecated")

			if err != nil {
				return nil, fmt.Errorf("Failed to derive ?isdeprecated= query parameter, %w", err)
			}

			if str_fl != "" {

				switch str_fl {
				case "-1", "0", "1":
					// ok
				default:
					return nil, fmt.Errorf("Invalid ?isdeprecated query parameter")
				}

				is_deprecated_f, err := spelunker.NewIsDeprecatedFilterFromString(ctx, str_fl)

				if err != nil {
					return nil, fmt.Errorf("Failed to create new is deprecated filter, %w", err)
				}

				filters = append(filters, is_deprecated_f)
			}

		case "country":

			country, err := sanitize.GetString(req, "country")

			if err != nil {
				return nil, fmt.Errorf("Failed to derive ?placetype= query parameter, %w", err)
			}

			if country != "" {

				country_f, err := spelunker.NewCountryFilterFromString(ctx, country)

				if err != nil {
					return nil, fmt.Errorf("Failed to create country filter from string '%s', %w", country, err)
				}

				filters = append(filters, country_f)
			}

		case "tag":

			tag, err := sanitize.GetString(req, "tag")

			if err != nil {
				return nil, fmt.Errorf("Failed to derive ?placetype= query parameter, %w", err)
			}

			if tag != "" {

				tag_f, err := spelunker.NewTagFilterFromString(ctx, tag)

				if err != nil {
					return nil, fmt.Errorf("Failed to create tag filter from string '%s', %w", tag, err)
				}

				filters = append(filters, tag_f)
			}

		case "placetype":

			placetype, err := sanitize.GetString(req, "placetype")

			if err != nil {
				return nil, fmt.Errorf("Failed to derive ?placetype= query parameter, %w", err)
			}

			if placetype != "" {

				placetype_f, err := spelunker.NewPlacetypeFilterFromString(ctx, placetype)

				if err != nil {
					return nil, fmt.Errorf("Failed to create placetype filter from string '%s', %w", placetype, err)
				}

				filters = append(filters, placetype_f)
			}

		default:
			return nil, fmt.Errorf("Invalid or unsupported parameter, %s", p)
		}
	}

	return filters, nil
}
