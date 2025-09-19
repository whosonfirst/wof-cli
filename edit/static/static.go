package static

import (
	"embed"
)

//go:embed css/* javascript/* *.html
var FS embed.FS
