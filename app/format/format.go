package format

import (
	"context"
	"io"
	"flag"
	"fmt"
	
	"github.com/whosonfirst/go-whosonfirst-format"
	"github.com/whosonfirst/wof/reader"
	"github.com/whosonfirst/wof/writer"		
)

func Run(ctx context.Context) error {

	fs := DefaultFlagSet()
	return RunWithFlagSet(ctx, fs)
}

func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet) error {

	opts, err := RunOptionsFromFlagSet(ctx, fs)

	if err != nil {
		return fmt.Errorf("Failed to create run options, %w", err)
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
