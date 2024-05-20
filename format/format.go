package format

import (
	"context"
	"fmt"
	"io"

	"github.com/whosonfirst/go-whosonfirst-format"
	"github.com/whosonfirst/wof"
	"github.com/whosonfirst/wof/reader"
	"github.com/whosonfirst/wof/writer"
)

type RunOptions struct {
	URIs      []string
	Overwrite bool
}

type FormatCommand struct {
	wof.Command
}

func init() {
	ctx := context.Background()
	wof.RegisterCommand(ctx, "fmt", NewFormatCommand)
}

func NewFormatCommand(ctx context.Context, cmd string) (wof.Command, error) {
	c := &FormatCommand{}
	return c, nil
}

func (c *FormatCommand) Run(ctx context.Context, args []string) error {

	fs := DefaultFlagSet()
	fs.Parse(args)

	uris := fs.Args()

	opts := &RunOptions{
		URIs:      uris,
		Overwrite: overwrite,
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

		new_body, err := format.FormatBytes(body)

		if err != nil {
			return fmt.Errorf("Failed to format body for '%s', %w", uri, err)
		}

		err = writer.Write(ctx, uri, new_body, opts.Overwrite)

		if err != nil {
			return fmt.Errorf("Failed to write body for '%s', %w", uri, err)
		}
	}

	return nil
}
