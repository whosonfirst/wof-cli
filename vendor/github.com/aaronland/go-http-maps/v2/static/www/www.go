package www

import (
	"embed"
)

//go:embed css/* javascript/* index.html
var FS embed.FS
