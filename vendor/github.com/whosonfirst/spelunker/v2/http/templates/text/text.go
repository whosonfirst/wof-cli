// Package text provides functions for working with ".txt" templates.
package text

import (
	"context"
	"embed"
	"text/template"

	sfomuseum_text "github.com/sfomuseum/go-template/text"
)

//go:embed *.txt
var FS embed.FS

// LoadTemplates instantiates ".txt" templates.
func LoadTemplates(ctx context.Context) (*template.Template, error) {

	return sfomuseum_text.LoadTemplates(ctx, FS)
}
