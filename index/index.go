package index

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	sql_index "github.com/whosonfirst/go-whosonfirst-database/app/sql/tables/index"
	"github.com/whosonfirst/wof"
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

func usage() {

	fmt.Println("Usage: wof index [TARGET] [OPTIONS]")
	fmt.Println("Valid commands are:")

	targets := []string{
		"sql",
	}

	for _, cmd := range targets {
		fmt.Printf("* %s\n", cmd)
	}

	os.Exit(0)
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

	args := os.Args

	if len(args) < 3 {
		usage()
	}

	switch args[2] {
	case "sql":
		fs := sql_index.DefaultFlagSet()
		fs.Parse(args[3:])
		return sql_index.RunWithFlagSet(ctx, fs)
	default:
		return fmt.Errorf("Invalid database target")
	}

	return nil
}
