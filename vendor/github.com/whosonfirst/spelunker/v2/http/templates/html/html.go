// Package html provides functions for working with ".html" templates.
package html

import (
	"context"
	"embed"

	sfomuseum_html "github.com/sfomuseum/go-template/html"
	"html/template"
)

//go:embed *.html
var FS embed.FS

// LoadTemplates instantiates ".html" templates.
func LoadTemplates(ctx context.Context) (*template.Template, error) {

	return sfomuseum_html.LoadTemplates(ctx, FS)
}
