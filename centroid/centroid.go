package centroid

import (
	"context"
	"fmt"
	_ "log/slog"

	"github.com/whosonfirst/go-whosonfirst-feature/properties"
	"github.com/whosonfirst/wof"
	"github.com/whosonfirst/wof/reader"
	"github.com/whosonfirst/wof/uris"
)

type CentroidCommand struct {
	wof.Command
}

func init() {
	ctx := context.Background()
	wof.RegisterCommand(ctx, "centroid", NewCentroidCommand)
}

func NewCentroidCommand(ctx context.Context, cmd string) (wof.Command, error) {
	c := &CentroidCommand{}
	return c, nil
}

func (c *CentroidCommand) Run(ctx context.Context, args []string) error {

	fs := DefaultFlagSet()
	fs.Parse(args)

	fs_uris := fs.Args()

	cb := func(ctx context.Context, cb_uri string) error {

		body, err := reader.BytesFromURI(ctx, cb_uri)

		if err != nil {
			return fmt.Errorf("Failed to open %s for reading, %w", cb_uri, err)
		}

		c, src, err := properties.Centroid(body)

		if err != nil {
			return fmt.Errorf("Failed to derive centroid for %s, %w", cb_uri, err)
		}

		fmt.Printf("%s,%f,%f\n", src, c.X(), c.Y())
		return nil
	}

	return uris.ExpandURIsWithCallback(ctx, cb, fs_uris...)
}
