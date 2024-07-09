package geometry

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/paulmach/orb/geojson"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-whosonfirst-export/v2"
	"github.com/whosonfirst/wof"
	"github.com/whosonfirst/wof/reader"
	"github.com/whosonfirst/wof/uris"
	"github.com/whosonfirst/wof/writer"
)

type GeometryCommand struct {
	wof.Command
}

func init() {
	ctx := context.Background()
	wof.RegisterCommand(ctx, "geometry", NewGeometryCommand)
}

func NewGeometryCommand(ctx context.Context, cmd string) (wof.Command, error) {
	c := &GeometryCommand{}
	return c, nil
}

func (c *GeometryCommand) Run(ctx context.Context, args []string) error {

	fs := DefaultFlagSet()
	fs.Parse(args)

	var exporter export.Exporter
	var source_geom *geojson.Geometry
	var geom_cb func(context.Context, string, []byte) error

	switch action {
	case "update":
		ex, err := export.NewExporter(ctx, "whosonfirst://")

		if err != nil {
			return fmt.Errorf("Failed to create new exporter, %w", err)
		}

		exporter = ex

		geom, err := c.deriveSourceGeometry(ctx, source, format)

		if err != nil {
			return fmt.Errorf("Failed to derive source geometry, %w", err)
		}

		source_geom = geom

		geom_cb = func(ctx context.Context, uri string, body []byte) error {

			updates := map[string]interface{}{
				"geometry": source_geom,
			}

			has_changes, new_body, err := export.AssignPropertiesIfChanged(ctx, body, updates)

			if err != nil {
				return fmt.Errorf("Failed to assign properties for '%s', %w", uri, err)
			}

			if !has_changes {
				slog.Info("No difference between geometries, not updating")
				return nil
			}

			new_body, err = exporter.Export(ctx, new_body)

			if err != nil {
				return fmt.Errorf("Failed to export body for '%s', %w", uri, err)
			}

			wr_uri := uri

			if stdout {
				wr_uri = writer.STDOUT
			}

			err = writer.Write(ctx, wr_uri, new_body)

			if err != nil {
				return fmt.Errorf("Failed to write body for '%s', %w", uri, err)
			}

			return nil
		}

	case "show":

		geom_cb = func(ctx context.Context, uri string, body []byte) error {
			geom_rsp := gjson.GetBytes(body, "geometry")
			fmt.Println(geom_rsp.String())
			return nil
		}

	default:
		return fmt.Errorf("Invalid or unsupported action")
	}

	uris_cb := func(ctx context.Context, uri string) error {

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

		return geom_cb(ctx, uri, body)
	}

	input_uris := fs.Args()

	err := uris.ExpandURIsWithCallback(ctx, uris_cb, input_uris...)

	if err != nil {
		return fmt.Errorf("Failed to run, %w", err)
	}

	return nil
}

func (c *GeometryCommand) deriveSourceGeometry(ctx context.Context, source string, format string) (*geojson.Geometry, error) {

	// To do: Support alternate sources (STDIN, etc) and formats (WKB, etc)

	source_r, err := os.Open(source)

	if err != nil {
		return nil, fmt.Errorf("Failed to open source, %w", err)
	}

	defer source_r.Close()

	source_body, err := io.ReadAll(source_r)

	if err != nil {
		return nil, fmt.Errorf("Failed to read source, %w", err)
	}

	source_f, err := geojson.UnmarshalFeature(source_body)

	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal source feature, %w", err)
	}

	source_geom := source_f.Geometry
	geojson_geom := geojson.NewGeometry(source_geom)

	return geojson_geom, nil
}
