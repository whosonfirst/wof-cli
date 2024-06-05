package emit

import (
	"context"
	_ "fmt"
	_ "io"
	"log/slog"

	"github.com/whosonfirst/wof"	
	_ "github.com/whosonfirst/go-writer-featurecollection/v3"
	_ "github.com/whosonfirst/go-writer-jsonl/v3"	
	app "github.com/whosonfirst/go-whosonfirst-iterwriter/app/iterwriter"
)

type EmitCommand struct {
	wof.Command
}

func init() {
	ctx := context.Background()
	wof.RegisterCommand(ctx, "emit", NewEmitCommand)
}

func NewEmitCommand(ctx context.Context, cmd string) (wof.Command, error) {
	c := &EmitCommand{}
	return c, nil
}

func (c *EmitCommand) Run(ctx context.Context, args []string) error {

	logger := slog.Default()

	fs := app.DefaultFlagSet()

	opts, err := app.DefaultOptionsFromFlagSet(fs, false)

	if err != nil {
		return err
	}
	
	return app.RunWithOptions(ctx, opts, logger)
}
