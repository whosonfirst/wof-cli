package index

import (
	"context"
	"log/slog"

	"github.com/whosonfirst/wof"	
	sql_index "github.com/whosonfirst/go-whosonfirst-database/app/sql/tables/index"
)

type RunOptions struct {
	Verbose bool
}

type IndexCommand struct {
	wof.Command
}

func init() {
	ctx := context.Background()
	wof.RegisterCommand(ctx, "index", NewIndexCommand)
}

func NewIndexCommand(ctx context.Context, cmd string) (wof.Command, error) {
	c := &IndexCommand{}
	return c, nil
}

func (c *IndexCommand) Run(ctx context.Context, args []string) error {

	fs := DefaultFlagSet()
	fs.Parse(args)

	opts := &RunOptions{
		Verbose: verbose,
	}

	return RunWithOptions(ctx, opts)
}

func RunWithOptions(ctx context.Context, opts *RunOptions) error {

	if opts.Verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("Verbose logging enabled")
	}

	return sql_index.Run(ctx)
}
