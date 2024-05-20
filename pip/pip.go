package pip

/*

#!/bin/sh

TILES_URI=file:///usr/local/data
TILES_DATABASE=architecture
TILES_ZOOM=12
TILES_LAYER=architecture

ENC_TILES_URI=`urlescape ${TILES_URI}`

./bin/wof pip \
    -spatial-database-uri "pmtiles://?tiles=${ENC_TILES_URI}&database=${TILES_DATABASE}&zoom=${TILES_ZOOM}&layer=${TILES_LAYER}&enable-cache=true" \
    -mapshaper-client-uri https://localhost:8080 \
    ~/Desktop/sfo.geojson


*/

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/sfomuseum/go-sfomuseum-mapshaper"
	_ "github.com/whosonfirst/go-whosonfirst-spatial-pmtiles"
	_ "github.com/whosonfirst/go-whosonfirst-spatial-sqlite"
	"github.com/whosonfirst/go-whosonfirst-spatial/database"
	"github.com/whosonfirst/go-whosonfirst-spatial/filter"
	"github.com/whosonfirst/go-whosonfirst-spatial/hierarchy"
	hierarchy_filter "github.com/whosonfirst/go-whosonfirst-spatial/hierarchy/filter"
	"github.com/whosonfirst/wof"
	"github.com/whosonfirst/wof/reader"
	"github.com/whosonfirst/wof/writer"
)

type RunOptions struct {
	URIs           []string
	Overwrite      bool
	Resolver       *hierarchy.PointInPolygonHierarchyResolver
	ResultsFunc    hierarchy_filter.FilterSPRResultsFunc
	InputFilters   *filter.SPRInputs
	UpdateCallback hierarchy.PointInPolygonHierarchyResolverUpdateCallback
}

type PointInPolygonCommand struct {
	wof.Command
}

func init() {
	ctx := context.Background()
	wof.RegisterCommand(ctx, "pip", NewPointInPolygonCommand)
}

func NewPointInPolygonCommand(ctx context.Context, cmd string) (wof.Command, error) {

	c := &PointInPolygonCommand{}
	return c, nil
}

func (c *PointInPolygonCommand) Run(ctx context.Context, args []string) error {

	fs := DefaultFlagSet()
	fs.Parse(args)

	uris := fs.Args()

	spatial_db, err := database.NewSpatialDatabase(ctx, spatial_database_uri)

	if err != nil {
		return fmt.Errorf("Failed to create new spatial database, %w", err)
	}

	resolver_opts := &hierarchy.PointInPolygonHierarchyResolverOptions{
		Database: spatial_db,
	}

	if mapshaper_client_uri != "" {

		mapshaper_client, err := mapshaper.NewClient(ctx, mapshaper_client_uri)

		if err != nil {
			return fmt.Errorf("Failed to create new mapshaper client, %w", err)
		}

		resolver_opts.Mapshaper = mapshaper_client
	}

	resolver, err := hierarchy.NewPointInPolygonHierarchyResolver(ctx, resolver_opts)

	if err != nil {
		return fmt.Errorf("Failed to create hierarchy resolver, %w", err)
	}

	// To do: update this from CLI flags

	input_filters := &filter.SPRInputs{}

	// To do: Implement Who's On First hierarchy rules

	results_func := hierarchy_filter.FirstButForgivingSPRResultsFunc

	update_cb := hierarchy.DefaultPointInPolygonHierarchyResolverUpdateCallback()

	opts := &RunOptions{
		URIs:           uris,
		Overwrite:      overwrite,
		Resolver:       resolver,
		ResultsFunc:    results_func,
		InputFilters:   input_filters,
		UpdateCallback: update_cb,
	}

	return RunWithOptions(ctx, opts)
}

func RunWithOptions(ctx context.Context, opts *RunOptions) error {

	for _, uri := range opts.URIs {

		r, is_stdin, err := reader.ReadCloserFromURI(ctx, uri)

		if err != nil {
			return fmt.Errorf("Failed to open '%s' for reading, %w", uri, err)
		}

		if !is_stdin {
			defer r.Close()
		}

		body, err := io.ReadAll(r)

		if err != nil {
			return fmt.Errorf("Failed to read '%s', %w", uri, err)
		}

		has_changed, new_body, err := opts.Resolver.PointInPolygonAndUpdate(ctx, opts.InputFilters, opts.ResultsFunc, opts.UpdateCallback, body)

		if err != nil {
			return fmt.Errorf("Failed to perform point in polygon operation, %w", err)
		}

		if !has_changed {
			slog.Info("No changes")
			continue
		}

		// To do: Still need to "export" record here

		err = writer.Write(ctx, uri, new_body, opts.Overwrite)

		if err != nil {
			return fmt.Errorf("Failed to write body for '%s', %w", uri, err)
		}
	}

	return nil
}
