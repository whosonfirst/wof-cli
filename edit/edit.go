package edit

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/whosonfirst/wof"
	"github.com/whosonfirst/wof/reader"
	"github.com/whosonfirst/wof/uris"
	// "github.com/whosonfirst/wof/writer"
	"github.com/sfomuseum/go-www-show"
)

type RunOptions struct {
	URIs   []string
	Stdout bool
}

type EditCommand struct {
	wof.Command
}

func init() {
	ctx := context.Background()
	wof.RegisterCommand(ctx, "edit", NewEditCommand)
}

func NewEditCommand(ctx context.Context, cmd string) (wof.Command, error) {

	c := &EditCommand{}

	return c, nil
}

func (c *EditCommand) Run(ctx context.Context, args []string) error {

	fs := DefaultFlagSet()
	fs.Parse(args)

	uris := fs.Args()

	opts := &RunOptions{
		URIs: uris,
	}

	return RunWithOptions(ctx, opts)
}

func RunWithOptions(ctx context.Context, opts *RunOptions) error {

	cb := func(ctx context.Context, uri string) error {

		r, is_stdin, err := reader.ReadCloserFromURI(ctx, uri)

		if err != nil {
			return fmt.Errorf("Failed to open '%s' for reading, %w", uri, err)
		}

		if !is_stdin {
			defer r.Close()
		}

		_, err = io.ReadAll(r)

		if err != nil {
			return fmt.Errorf("Failed to read '%s', %w", uri, err)
		}

		return nil
	}

	err := uris.ExpandURIsWithCallback(ctx, cb, opts.URIs...)

	if err != nil {
		return fmt.Errorf("Failed to run, %w", err)
	}

	mux := http.NewServeMux()

	browser, err := show.NewBrowser(ctx, "web://")

	if err != nil {
		return err
	}

	show_opts := &show.RunOptions{
		Browser: browser,
		Mux:     mux,
	}

	return show.RunWithOptions(ctx, show_opts)

	return nil
}
