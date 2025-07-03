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

	"github.com/sfomuseum/go-sfomuseum-mapshaper"
	"github.com/tidwall/sjson"
	"github.com/whosonfirst/go-whosonfirst-export/v3"
	"github.com/whosonfirst/go-whosonfirst-spatial/database"
	"github.com/whosonfirst/go-whosonfirst-spatial/filter"
	"github.com/whosonfirst/go-whosonfirst-spatial/hierarchy"
	hierarchy_filter "github.com/whosonfirst/go-whosonfirst-spatial/hierarchy/filter"
	"github.com/whosonfirst/wof"
	"github.com/whosonfirst/wof/reader"
	"github.com/whosonfirst/wof/writer"
)

type RunOptions struct {
	// A list of URIs to perform point-in-polygon operations on.
	URIs []string
	// Write output to STDOUT
	Stdout bool
	// ...
	Export bool
	// ...
	Exporter export.Exporter
	// ...
	Placetype string
	// Resolver is a `PointInPolygonHierarchyResolver` instance which is what manages all the point-in-polygon operations and related decision making.
	Resolver *hierarchy.PointInPolygonHierarchyResolver
	// ResultsFunc is function used to derive a single `spr.StandardPlacesResult` from a list of `spr.StandardPlacesResult` instances (point-in-polygon candidates).
	ResultsFunc hierarchy_filter.FilterSPRResultsFunc
	// InputFilters are filterd used to limits the point-in-polygon candidates to consider.
	InputFilters *filter.SPRInputs
	// UpdateCallback is a function definition for a custom callback to convert point-in-polygon candidates in to a dictionary of properties containining hierarchy information
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
		Stdout:         stdout,
		Placetype:      placetype,
		Resolver:       resolver,
		ResultsFunc:    results_func,
		InputFilters:   input_filters,
		UpdateCallback: update_cb,
	}

	if exportify {

		ex, err := export.NewExporter(ctx, "whosonfirst://")

		if err != nil {
			return fmt.Errorf("Failed to create new exporter, %w", err)
		}

		opts.Export = true
		opts.Exporter = ex
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

		if opts.Placetype != "" {

			body, err = sjson.SetBytes(body, "properties.wof:placetype", opts.Placetype)

			if err != nil {
				return fmt.Errorf("Failed to assign wof:placetype property for '%s', %w", uri, err)
			}
		}

		has_changed, new_body, err := opts.Resolver.PointInPolygonAndUpdate(ctx, opts.InputFilters, opts.ResultsFunc, opts.UpdateCallback, body)

		if err != nil {
			return fmt.Errorf("Failed to perform point in polygon operation, %w", err)
		}

		if has_changed {

			wr_uri := uri

			if opts.Stdout {
				wr_uri = writer.STDOUT
			}

			if opts.Export && opts.Exporter != nil {

				_, new_body, err = opts.Exporter.Export(ctx, new_body)

				if err != nil {
					return fmt.Errorf("Failed to export new record for '%s' ('%s'), %w", wr_uri, uri, err)
				}
			}

			err = writer.Write(ctx, wr_uri, new_body)

			if err != nil {
				return fmt.Errorf("Failed to write body for '%s' ('%s'), %w", wr_uri, uri, err)
			}
		}
	}

	return nil
}
