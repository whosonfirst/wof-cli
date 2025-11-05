package index

import (
	"context"
	"fmt"
	"os"

	sql_index "github.com/whosonfirst/go-whosonfirst-database/app/sql/tables/index"
	"github.com/whosonfirst/wof"
)

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

	if len(args) < 2 {
		usage()
	}

	switch args[0] {
	case "sql":

		fs := sql_index.DefaultFlagSet()
		fs.Parse(args[1:])
		args := fs.Args()

		opts, err := sql_index.RunOptionsFromParsedFlags(args...)

		if err != nil {
			return fmt.Errorf("Failed to derive run options, %w", err)
		}

		return sql_index.RunWithOptions(ctx, opts)
	default:
		return fmt.Errorf("Invalid database target")
	}

	return nil
}
