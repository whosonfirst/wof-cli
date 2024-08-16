package show

import (
	"context"
	"fmt"
	"io"
	"slices"

	"github.com/paulmach/orb/geojson"
	sfom_show "github.com/sfomuseum/go-geojson-show"
	"github.com/whosonfirst/wof"
	"github.com/whosonfirst/wof/reader"
	"github.com/whosonfirst/wof/uris"
)

type RunOptions struct {
	URIs           []string
	MapProvider    string
	MapTileURI     string
	ProtomapsTheme string
	Port           int
}

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

	fs := sfom_show.DefaultFlagSet()
	fs.Parse(args)

	fs_uris := fs.Args()

	run_opts, err := sfom_show.RunOptionsFromFlagSet(fs)

	if err != nil {
		return fmt.Errorf("Failed to derive run options, %w", err)
	}

	label_props := []string{
		"wof:name",
		"wof:id",
		"wof:placetype",
		"src:geom",
	}

	for _, prop := range run_opts.LabelProperties {

		if !slices.Contains(label_props, prop) {
			label_props = append(label_props, prop)
		}
	}

	run_opts.LabelProperties = label_props

	// Add styles here...

	fc := geojson.NewFeatureCollection()

	cb := func(ctx context.Context, uri string) error {

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

		f, err := geojson.UnmarshalFeature(body)

		if err != nil {
			return fmt.Errorf("Failed to unmarshal '%s' as GeoJSON, %w", uri, err)
		}

		fc.Append(f)
		return nil
	}

	err = uris.ExpandURIsWithCallback(ctx, cb, fs_uris...)

	if err != nil {
		return fmt.Errorf("Failed to run, %w", err)
	}

	run_opts.Features = fc.Features

	return sfom_show.RunWithOptions(ctx, run_opts)
}
