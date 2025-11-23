// Package javascript provides functions for working with ".js" templates.
package javascript

import (
	"context"
	"embed"

	sfomuseum_text "github.com/sfomuseum/go-template/text"
	"text/template"
)

//go:embed *.js
var FS embed.FS

// LoadTemplates instantiates ".js" templates.
func LoadTemplates(ctx context.Context) (*template.Template, error) {

	return sfomuseum_text.LoadTemplatesMatching(ctx, "*.js", FS)
}
