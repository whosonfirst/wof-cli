package show

import (
	"context"
	"fmt"
	"io"

	"github.com/paulmach/orb/geojson"
	"github.com/whosonfirst/wof"
	"github.com/whosonfirst/wof/reader"
	"github.com/whosonfirst/wof/uris"
	sfom_show "github.com/sfomuseum/go-geojson-show/app/show"
)

type RunOptions struct {
	URIs []string
	MapProvider string
	MapTileURI  string
	ProtomapsTheme string
	Port int
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

	fs := DefaultFlagSet()
	fs.Parse(args)

	uris := fs.Args()

	opts := &RunOptions{
		URIs: uris,
		MapProvider: map_provider,
		MapTileURI: map_tile_uri,
		ProtomapsTheme: protomaps_theme,
		Port: port,
	}

	return RunWithOptions(ctx, opts)
}

func RunWithOptions(ctx context.Context, opts *RunOptions) error {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

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

	err := uris.ExpandURIsWithCallback(ctx, cb, opts.URIs...)

	if err != nil {
		return fmt.Errorf("Failed to run, %w", err)
	}

	run_opts := &sfom_show.RunOptions{
		MapProvider: opts.MapProvider,
		MapTileURI: opts.MapTileURI,
		ProtomapsTheme: opts.ProtomapsTheme,
		Port: opts.Port,		
		Features: fc.Features,		
	}

	return sfom_show.RunWithOptions(ctx, run_opts)
}
