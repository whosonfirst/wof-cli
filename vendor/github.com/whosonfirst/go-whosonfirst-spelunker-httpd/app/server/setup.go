package server

import (
	"context"
	"fmt"
	html_template "html/template"
	"log/slog"

	"github.com/rs/cors"
	"github.com/sfomuseum/go-http-auth"
	"github.com/whosonfirst/go-whosonfirst-spelunker"
)

func setupCommon() {

	ctx := context.Background()
	var err error

	// defined in vars.go
	sp, err = spelunker.NewSpelunker(ctx, run_options.SpelunkerURI)

	if err != nil {
		setupCommonError = fmt.Errorf("Failed to set up network, %w", err)
		return
	}
}

func setupAPI() {

	setupCommonOnce.Do(setupCommon)

	if setupCommonError != nil {
		slog.Error("Failed to set up common configuration", "error", setupCommonError)
		setupAPIError = fmt.Errorf("Failed to set up common configuration, %w", setupCommonError)
		return
	}

	// Please finesse me...
	cors_origins := []string{
		"*",
	}

	cors_wrapper = cors.New(cors.Options{
		AllowedOrigins:   cors_origins,
		AllowCredentials: false,
		Debug:            false,
	})
}

func setupWWW() {

	ctx := context.Background()
	var err error

	setupCommonOnce.Do(setupCommon)

	if setupCommonError != nil {
		slog.Error("Failed to set up common configuration", "error", setupCommonError)
		setupWWWError = fmt.Errorf("Common setup failed, %w", setupCommonError)
		return
	}

	// defined in vars.go
	authenticator, err = auth.NewAuthenticator(ctx, run_options.AuthenticatorURI)

	if err != nil {
		setupWWWError = fmt.Errorf("Failed to create new authenticator, %w", err)
		return
	}

	// defined in vars.go
	html_templates = html_template.New("html").Funcs(run_options.HTMLTemplateFuncs)

	for idx, f := range run_options.HTMLTemplates {

		html_templates, err = html_templates.ParseFS(f, "*.html")

		if err != nil {
			setupWWWError = fmt.Errorf("Failed to load templates from FS at offset %d, %w", idx, err)
			return
		}
	}

}
