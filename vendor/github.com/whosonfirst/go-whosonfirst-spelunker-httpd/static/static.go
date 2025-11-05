package static

import (
	"embed"
)

//go:embed css/*.css javascript/*.js fonts/*
var FS embed.FS
