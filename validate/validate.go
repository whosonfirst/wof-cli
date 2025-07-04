package validate

import (
	"context"
	"fmt"
	_ "log/slog"

	"github.com/whosonfirst/go-whosonfirst-validate"
	"github.com/whosonfirst/wof"
	"github.com/whosonfirst/wof/reader"
	"github.com/whosonfirst/wof/uris"
)

type RunOptions struct {
	URIs            []string
	Overwrite       bool
	ValidateOptions *validate.Options
}

type ValidateCommand struct {
	wof.Command
}

func init() {
	ctx := context.Background()
	wof.RegisterCommand(ctx, "validate", NewValidateCommand)
}

func NewValidateCommand(ctx context.Context, cmd string) (wof.Command, error) {

	c := &ValidateCommand{}
	return c, nil
}

func (c *ValidateCommand) Run(ctx context.Context, args []string) error {

	fs := DefaultFlagSet()
	fs.Parse(args)

	uris := fs.Args()

	// To do: read from flags

	validate_opts := &validate.Options{
		ValidateName:      true,
		ValidatePlacetype: true,
		ValidateRepo:      true,
		ValidateEDTF:      true,
		ValidateIsCurrent: true,
		ValidateNames:     true,
	}

	opts := &RunOptions{
		URIs:            uris,
		ValidateOptions: validate_opts,
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

		body, err := validate.EnsureValidGeoJSON(r)

		if err != nil {
			return fmt.Errorf("Failed to read '%s', %w", uri, err)
		}

		err = validate.ValidateWithOptions(body, opts.ValidateOptions)

		if err != nil {
			return fmt.Errorf("Failed to validate body for '%s', %w", uri, err)
		}

		return nil
	}

	err := uris.ExpandURIsWithCallback(ctx, cb, opts.URIs...)

	if err != nil {
		return fmt.Errorf("Failed to run, %w", err)
	}

	return nil
}
