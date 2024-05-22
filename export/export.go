package export

import (
	"context"
	"fmt"
	"io"

	"github.com/whosonfirst/go-whosonfirst-export/v2"
	"github.com/whosonfirst/wof"
	"github.com/whosonfirst/wof/reader"
	"github.com/whosonfirst/wof/uris"
	"github.com/whosonfirst/wof/writer"
)

type RunOptions struct {
	URIs     []string
	Stdout   bool
	Exporter export.Exporter
}

type ExportCommand struct {
	wof.Command
	exporter export.Exporter
}

func init() {
	ctx := context.Background()
	wof.RegisterCommand(ctx, "export", NewExportCommand)
}

func NewExportCommand(ctx context.Context, cmd string) (wof.Command, error) {

	ex, err := export.NewExporter(ctx, "whosonfirst://")

	if err != nil {
		return nil, fmt.Errorf("Failed to create new exporter, %w", err)
	}

	c := &ExportCommand{
		exporter: ex,
	}

	return c, nil
}

func (c *ExportCommand) Run(ctx context.Context, args []string) error {

	fs := DefaultFlagSet()
	fs.Parse(args)

	uris := fs.Args()

	opts := &RunOptions{
		URIs:     uris,
		Stdout:   stdout,
		Exporter: c.exporter,
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

		body, err := io.ReadAll(r)

		if err != nil {
			return fmt.Errorf("Failed to read '%s', %w", uri, err)
		}

		new_body, err := opts.Exporter.Export(ctx, body)

		if err != nil {
			return fmt.Errorf("Failed to export body for '%s', %w", uri, err)
		}

		wr_uri := uri

		if opts.Stdout {
			wr_uri = writer.STDOUT
		}

		err = writer.Write(ctx, wr_uri, new_body)

		if err != nil {
			return fmt.Errorf("Failed to write body for '%s', %w", uri, err)
		}

		return nil
	}

	err := uris.ExpandURIsWithCallback(ctx, cb, opts.URIs...)

	if err != nil {
		return fmt.Errorf("Failed to run, %w", err)
	}

	return nil
}
