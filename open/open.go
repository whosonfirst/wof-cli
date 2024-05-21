package open

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/whosonfirst/wof"
	"github.com/whosonfirst/wof/uris"
)

type RunOptions struct {
	URIs   []string
	Editor string
}

type OpenCommand struct {
	wof.Command
}

func init() {
	ctx := context.Background()
	wof.RegisterCommand(ctx, "open", NewOpenCommand)
}

func NewOpenCommand(ctx context.Context, cmd string) (wof.Command, error) {

	c := &OpenCommand{}
	return c, nil
}

func (c *OpenCommand) Run(ctx context.Context, args []string) error {

	fs := DefaultFlagSet()
	fs.Parse(args)

	uris := fs.Args()

	if editor == "" {
		editor = os.Getenv("EDITOR")
	}

	if editor == "" {
		return fmt.Errorf("Undefined editor")
	}

	opts := &RunOptions{
		URIs:   uris,
		Editor: editor,
	}

	return RunWithOptions(ctx, opts)
}

func RunWithOptions(ctx context.Context, opts *RunOptions) error {

	cb := func(ctx context.Context, uri string) error {
		return open(opts.Editor, uri)
	}

	err := uris.ExpandURIsWithCallback(ctx, cb, opts.URIs...)

	if err != nil {
		return fmt.Errorf("Failed to run, %w", err)
	}

	return nil
}

func open(editor string, path string) error {

	cmd := exec.Command("sh", "-c", editor+" "+path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
