package text

import (
	"context"
	"embed"

	sfomuseum_text "github.com/sfomuseum/go-template/text"
	"text/template"
)

//go:embed *.txt
var FS embed.FS

func LoadTemplates(ctx context.Context) (*template.Template, error) {

	return sfomuseum_text.LoadTemplates(ctx, FS)
}
