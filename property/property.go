package property

import (
	"context"
	"fmt"
	_ "log/slog"

	"github.com/tidwall/gjson"
	"github.com/whosonfirst/wof"
	"github.com/whosonfirst/wof/reader"
	"github.com/whosonfirst/wof/uris"
)

type PropertyCommand struct {
	wof.Command
}

func init() {
	ctx := context.Background()
	wof.RegisterCommand(ctx, "property", NewPropertyCommand)
}

func NewPropertyCommand(ctx context.Context, cmd string) (wof.Command, error) {

	c := &PropertyCommand{}
	return c, nil
}

func (c *PropertyCommand) Run(ctx context.Context, args []string) error {

	fs := DefaultFlagSet()
	fs.Parse(args)

	fs_uris := fs.Args()

	cb := func(ctx context.Context, cb_uri string) error {

		body, err := reader.BytesFromURI(ctx, cb_uri)

		if err != nil {
			return fmt.Errorf("Failed to open %s for reading, %w", cb_uri, err)
		}
		
		for _, path := range paths {
			rsp := gjson.GetBytes(body, path)
			fmt.Println(rsp.String())
		}

		return nil
	}

	return uris.ExpandURIsWithCallback(ctx, cb, fs_uris...)
}
