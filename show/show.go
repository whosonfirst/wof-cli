package show

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"slices"

	"github.com/paulmach/orb/geojson"
	sfom_show "github.com/sfomuseum/go-geojson-show"
	"github.com/whosonfirst/wof"
	"github.com/whosonfirst/wof/reader"
	"github.com/whosonfirst/wof/uris"
)

//go:embed style.json
var style_json string

//go:embed point_style.json
var point_style_json string

type ShowCommand struct {
	wof.Command
}

func init() {
	ctx := context.Background()
	wof.RegisterCommand(ctx, "show", NewShowCommand)
}

func NewShowCommand(ctx context.Context, cmd string) (wof.Command, error) {

	c := &ShowCommand{}
	return c, nil
}

func (c *ShowCommand) Run(ctx context.Context, args []string) error {

	var is_featurecollection bool

	fs := sfom_show.DefaultFlagSet()
	fs.BoolVar(&is_featurecollection, "featurecollection", false, "Boolean flag indicating input data is a GeoJSON FeatureCollection.")

	fs.Parse(args)

	fs_uris := fs.Args()

	run_opts, err := sfom_show.RunOptionsFromFlagSet(ctx, fs)

	if err != nil {
		return fmt.Errorf("Failed to derive run options, %w", err)
	}

	// Ensure custom label properties

	label_props := []string{
		"wof:id",
		"wof:name",
		"wof:placetype",
		"mz:is_current",
		"src:geom",
	}

	for _, prop := range run_opts.LabelProperties {

		if !slices.Contains(label_props, prop) {
			label_props = append(label_props, prop)
		}
	}

	run_opts.LabelProperties = label_props

	// Ensure custom styles

	if run_opts.Style == nil {

		style, err := sfom_show.UnmarshalStyle(style_json)

		if err != nil {
			return fmt.Errorf("Failed to unmarshal style, %w", err)
		}

		run_opts.Style = style
	}

	if run_opts.PointStyle == nil {

		point_style, err := sfom_show.UnmarshalStyle(point_style_json)

		if err != nil {
			return fmt.Errorf("Failed to unmarshal point style, %w", err)
		}

		run_opts.PointStyle = point_style
	}

	// Derive features to show

	fc := geojson.NewFeatureCollection()

	cb := func(ctx context.Context, uri string) error {

		body, err := reader.BytesFromURI(ctx, uri)

		if err != nil {
			return fmt.Errorf("Failed to open '%s' for reading, %w", uri, err)
		}

		if is_featurecollection {

			f, err := geojson.UnmarshalFeatureCollection(body)

			if err != nil {

				os.WriteFile("wtf.json", body, 0644)

				return fmt.Errorf("Failed to unmarshal '%s' as GeoJSON FeatureCollection, %w", uri, err)
			}

			fc = f
		} else {

			f, err := geojson.UnmarshalFeature(body)

			if err != nil {
				return fmt.Errorf("Failed to unmarshal '%s' as GeoJSON, %w", uri, err)
			}

			fc.Append(f)
		}

		return nil
	}

	err = uris.ExpandURIsWithCallback(ctx, cb, fs_uris...)

	if err != nil {
		return fmt.Errorf("Failed to run, %w", err)
	}

	run_opts.Features = fc.Features

	// Show the map

	return sfom_show.RunWithOptions(ctx, run_opts)
}
