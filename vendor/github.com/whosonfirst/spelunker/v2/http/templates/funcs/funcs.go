// Package funcs defines Spelunker-specific HTML template functions.
package funcs

import (
	"fmt"
	"log/slog"
	"net/url"
	"slices"
	"strings"

	"github.com/whosonfirst/go-whosonfirst-sources"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// Name returns the full name for a source identified by its prefix. If not full name can be found returns the prefix itself.
func NameForSource(source string) string {

	nspred := strings.Split(source, ":")
	prefix := nspred[0]

	src, err := sources.GetSourceByPrefix(prefix)

	if err != nil {
		return prefix
	}

	return src.Fullname
}

// FormatNumbers returns the string-formatted value of 'i' using the `language.English` printer.
func FormatNumber(i int64) string {
	p := message.NewPrinter(language.English)
	return p.Sprintf("%d", i)
}

// AppendPagination appends pagination query parameters (k=v) to 'uri'.
func AppendPagination(uri string, k string, v any) string {

	u, err := url.Parse(uri)

	if err != nil {
		slog.Error("Failed to parse URI to append pagination", "uri", uri, "error", err)
		return "#"
	}

	q := u.Query()
	q.Set(k, fmt.Sprintf("%v", v))

	u.RawQuery = q.Encode()
	return u.String()
}

// IsAPlacetype returns 'pt' prefixed with 'a' or 'an'.
func IsAPlacetype(pt string) string {

	if pt == "custom" {
		return "a custom placetype"
	}

	// https://github.com/whosonfirst/spelunker/v2-httpd/issues/46

	vowels := []string{
		"a", "e", "i", "o", "u",
	}

	first := pt[0]

	if slices.Contains(vowels, string(first)) {
		return fmt.Sprintf("an %s", pt)
	} else {
		return fmt.Sprintf("a %s", pt)
	}

	return pt
}
