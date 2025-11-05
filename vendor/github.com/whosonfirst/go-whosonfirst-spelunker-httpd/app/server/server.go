package server

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/aaronland/go-http-server"
	"github.com/aaronland/go-http-server/handler"
)

func Run(ctx context.Context, logger *slog.Logger) error {
	fs := DefaultFlagSet()
	return RunWithFlagSet(ctx, fs, logger)
}

func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet, logger *slog.Logger) error {

	opts, err := RunOptionsFromFlagSet(ctx, fs)

	if err != nil {
		return fmt.Errorf("Failed to derive run options from flagset, %w", err)
	}

	return RunWithOptions(ctx, opts, logger)
}

func RunWithOptions(ctx context.Context, opts *RunOptions, logger *slog.Logger) error {

	slog.SetDefault(logger)

	// First create a local copy of RunOptions that can't be
	// modified after the fact. 'run_options' is defined in vars.go

	v, err := opts.Clone()

	if err != nil {
		return fmt.Errorf("Failed to create local run options, %w", err)
	}

	run_options = v

	// To do: Add/consult "is enabled" flags

	// START OF defer loading handlers (and all their dependencies) until they are actually routed to
	// in case we are running in a "serverless" environment like AWS Lambda

	path_urisjs, err := url.JoinPath(run_options.URIs.Static, "/javascript/whosonfirst.spelunker.uris.js")

	if err != nil {
		return fmt.Errorf("Failed to construct path for whosonfirst.spelunker.uris.js, %w", err)
	}

	handlers := map[string]handler.RouteHandlerFunc{

		// Common handler things
		"/robots.txt": robotsTxtHandlerFunc,

		// WWW/human-readable
		run_options.URIs.Placetypes:        placetypesHandlerFunc,
		run_options.URIs.Placetype:         hasPlacetypeHandlerFunc,
		run_options.URIs.Concordances:      concordancesHandlerFunc,
		run_options.URIs.ConcordanceNS:     hasConcordanceHandlerFunc,
		run_options.URIs.ConcordanceNSPred: hasConcordanceHandlerFunc,
		run_options.URIs.ConcordanceTriple: hasConcordanceHandlerFunc,
		run_options.URIs.Recent:            recentHandlerFunc,
		run_options.URIs.NullIsland:        nullIslandHandlerFunc,
		run_options.URIs.Descendants:       descendantsHandlerFunc,
		run_options.URIs.Id:                idHandlerFunc,
		run_options.URIs.Search:            searchHandlerFunc,
		run_options.URIs.About:             aboutHandlerFunc,
		run_options.URIs.Code:              codeHandlerFunc,
		run_options.URIs.HowTo:             howtoHandlerFunc,
		run_options.URIs.Index:             indexHandlerFunc,
		run_options.URIs.Tiles:             tilesHandlerFunc,
		run_options.URIs.OpenSearch:        openSearchHandlerFunc,

		// Static assets
		run_options.URIs.Static: staticHandlerFunc,
		// Run-time static assets
		path_urisjs: urisJSHandlerFunc,

		// API/machine-readable
		run_options.URIs.ConcordanceNSFaceted:     hasConcordanceFacetedHandlerFunc,
		run_options.URIs.ConcordanceNSPredFaceted: hasConcordanceFacetedHandlerFunc,
		run_options.URIs.ConcordanceTripleFaceted: hasConcordanceFacetedHandlerFunc,
		run_options.URIs.DescendantsFaceted:       descendantsFacetedHandlerFunc,
		run_options.URIs.FindingAid:               findingAidHandlerFunc,
		run_options.URIs.GeoJSON:                  geoJSONHandlerFunc,
		run_options.URIs.GeoJSONLD:                geoJSONLDHandlerFunc,
		run_options.URIs.NavPlace:                 navPlaceHandlerFunc,
		run_options.URIs.NullIslandFaceted:        nullIslandFacetedHandlerFunc,
		run_options.URIs.PlacetypeFaceted:         placetypeFacetedHandlerFunc,
		run_options.URIs.RecentFaceted:            recentFacetedHandlerFunc,
		run_options.URIs.SearchFaceted:            searchFacetedHandlerFunc,
		run_options.URIs.Select:                   selectHandlerFunc,
		run_options.URIs.SPR:                      sprHandlerFunc,
		run_options.URIs.SVG:                      svgHandlerFunc,
	}

	assign_handlers := func(handler_map map[string]handler.RouteHandlerFunc, paths []string, handler_func handler.RouteHandlerFunc) {

		for _, p := range paths {
			handler_map[p] = handler_func
		}
	}

	assign_handlers(handlers, run_options.URIs.IdAlt, idHandlerFunc)
	assign_handlers(handlers, run_options.URIs.DescendantsAlt, descendantsHandlerFunc)

	// API/machine-readable
	assign_handlers(handlers, run_options.URIs.GeoJSONAlt, geoJSONHandlerFunc)
	assign_handlers(handlers, run_options.URIs.GeoJSONLDAlt, geoJSONLDHandlerFunc)
	assign_handlers(handlers, run_options.URIs.NavPlaceAlt, navPlaceHandlerFunc)
	assign_handlers(handlers, run_options.URIs.SelectAlt, selectHandlerFunc)
	assign_handlers(handlers, run_options.URIs.RecentAlt, recentHandlerFunc)
	assign_handlers(handlers, run_options.URIs.SPRAlt, sprHandlerFunc)
	assign_handlers(handlers, run_options.URIs.SVGAlt, svgHandlerFunc)

	log_logger := slog.NewLogLogger(logger.Handler(), slog.LevelInfo)

	route_handler_opts := &handler.RouteHandlerOptions{
		Handlers: handlers,
		Logger:   log_logger,
	}

	route_handler, err := handler.RouteHandlerWithOptions(route_handler_opts)

	if err != nil {
		return fmt.Errorf("Failed to configure route handler, %w", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", route_handler)

	// END OF defer loading handlers (and all their dependencies) until they are actually routed to

	s, err := server.NewServer(ctx, run_options.ServerURI)

	if err != nil {
		return fmt.Errorf("Failed to create new server, %w", err)
	}

	go func() {
		for uri, h := range handlers {
			slog.Debug("Enable handler", "uri", uri, "handler", fmt.Sprintf("%T", h))
		}
	}()

	slog.Info("Listening for requests", "address", s.Address())
	err = s.ListenAndServe(ctx, mux)

	if err != nil {
		return fmt.Errorf("Failed to start server, %w", err)
	}

	return nil
}
