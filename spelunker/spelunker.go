package spelunker

import (
	"context"
	"log/slog"
	
	_ "github.com/whosonfirst/go-whosonfirst-spelunker-sql"
	
	"github.com/whosonfirst/go-whosonfirst-spelunker-httpd/app/server"
	"github.com/whosonfirst/wof"
)

type SpelunkerCommand struct {
	wof.Command
}

func init() {
	ctx := context.Background()
	wof.RegisterCommand(ctx, "spelunker", NewSpelunkerCommand)
}


func NewSpelunkerCommand(ctx context.Context, cmd string) (wof.Command, error) {
	c := &SpelunkerCommand{}
	return c, nil
}

func (c *SpelunkerCommand) Run(ctx context.Context, args []string) error {

	if len(args) < 2 {
		return nil
	}

	fs := server.DefaultFlagSet()
	fs.Parse(args[2:])
	
	logger := slog.Default()
	
	opts, err := server.RunOptionsFromFlagSet(ctx, fs)

	if err != nil {
		return err
	}

	return server.RunWithOptions(ctx, opts, logger)
}
