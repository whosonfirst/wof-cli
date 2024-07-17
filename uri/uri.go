package uri

import (
	"context"
	"fmt"
	_ "log/slog"
	"path/filepath"

	wof_uri "github.com/whosonfirst/go-whosonfirst-uri"
	"github.com/whosonfirst/wof"
	"github.com/whosonfirst/wof/uris"
)

type URICommand struct {
	wof.Command
}

func init() {
	ctx := context.Background()
	wof.RegisterCommand(ctx, "uri", NewURICommand)
}

func NewURICommand(ctx context.Context, cmd string) (wof.Command, error) {

	c := &URICommand{}
	return c, nil
}

func (c *URICommand) Run(ctx context.Context, args []string) error {

	fs := DefaultFlagSet()
	fs.Parse(args)

	fs_uris := fs.Args()

	cb := func(ctx context.Context, cb_uri string) error {

		id, uri_args, err := wof_uri.ParseURI(cb_uri)

		if err != nil {
			return err
		}

		rel_path, err := wof_uri.Id2RelPath(id, uri_args)

		if err != nil {
			return err
		}

		if prefix != "" {
			rel_path = filepath.Join(prefix, rel_path)
		}

		fmt.Println(rel_path)
		return nil
	}

	return uris.ExpandURIsWithCallback(ctx, cb, fs_uris...)
}
