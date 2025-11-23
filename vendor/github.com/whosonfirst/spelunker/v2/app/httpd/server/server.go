package server

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	_ "github.com/whosonfirst/spelunker/v2/app"

	"github.com/aaronland/go-http/v4/route"
	"github.com/aaronland/go-http/v4/server"
)

func Run(ctx context.Context) error {
	fs := DefaultFlagSet()
	return RunWithFlagSet(ctx, fs)
}

func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet) error {

	opts, err := RunOptionsFromFlagSet(ctx, fs)

	if err != nil {
		return fmt.Errorf("Failed to derive run options from flagset, %w", err)
	}

	return RunWithOptions(ctx, opts)
}

func RunWithOptions(ctx context.Context, opts *RunOptions) error {

	// First create a local copy of RunOptions that can't be
	// modified after the fact. 'run_options' is defined in vars.go

	v, err := opts.Clone()

	if err != nil {
		return fmt.Errorf("Failed to create local run options, %w", err)
	}

	run_options = v

	if run_options.Verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("Verbose (debug) logging enabled")
	}

	// To do: Add/consult "is enabled" flags

	// START OF defer loading handlers (and all their dependencies) until they are actually routed to
	// in case we are running in a "serverless" environment like AWS Lambda

	path_urisjs, err := url.JoinPath(run_options.URIs.Static, "/javascript/whosonfirst.spelunker.uris.js")

	if err != nil {
		return fmt.Errorf("Failed to construct path for whosonfirst.spelunker.uris.js, %w", err)
	}

	mux_handlers := map[string]route.RouteHandlerFunc{

		// Common handler things
		"/robots.txt": robotsTxtHandlerFunc,

		// Map tile handler/config stuff is dealth with below

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
		run_options.URIs.Index:             indexHandlerFunc,
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
		run_options.URIs.WKT:                      wktHandlerFunc,
	}

	map_cfg_handler, map_tile_handler, map_tile_url, err := mapConfigHandlers(ctx)

	if err != nil {
		return fmt.Errorf("Failed to derive map config handlers, %w", err)
	}

	mux_handlers["/maps.json"] = func(ctx context.Context) (http.Handler, error) {
		return map_cfg_handler, nil
	}

	if map_tile_handler != nil {
		mux_handlers[map_tile_url] = func(ctx context.Context) (http.Handler, error) {
			return map_tile_handler, nil
		}
	}

	assign_handlers := func(handler_map map[string]route.RouteHandlerFunc, paths []string, handler_func route.RouteHandlerFunc) {

		for _, p := range paths {
			handler_map[p] = handler_func
		}
	}

	assign_handlers(mux_handlers, run_options.URIs.IdAlt, idHandlerFunc)
	assign_handlers(mux_handlers, run_options.URIs.DescendantsAlt, descendantsHandlerFunc)

	// API/machine-readable
	assign_handlers(mux_handlers, run_options.URIs.GeoJSONAlt, geoJSONHandlerFunc)
	assign_handlers(mux_handlers, run_options.URIs.GeoJSONLDAlt, geoJSONLDHandlerFunc)
	assign_handlers(mux_handlers, run_options.URIs.NavPlaceAlt, navPlaceHandlerFunc)
	assign_handlers(mux_handlers, run_options.URIs.SelectAlt, selectHandlerFunc)
	assign_handlers(mux_handlers, run_options.URIs.RecentAlt, recentHandlerFunc)
	assign_handlers(mux_handlers, run_options.URIs.SPRAlt, sprHandlerFunc)
	assign_handlers(mux_handlers, run_options.URIs.SVGAlt, svgHandlerFunc)
	assign_handlers(mux_handlers, run_options.URIs.WKTAlt, wktHandlerFunc)

	route_handler_opts := &route.RouteHandlerOptions{
		Handlers: mux_handlers,
	}

	route_handler, err := route.RouteHandlerWithOptions(route_handler_opts)

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
		for uri, h := range mux_handlers {
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
