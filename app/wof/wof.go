package wof

import (
	"context"
	"fmt"
	"os"

	"github.com/whosonfirst/wof"
	_ "github.com/whosonfirst/wof/emit"
	_ "github.com/whosonfirst/wof/export"
	_ "github.com/whosonfirst/wof/format"
	_ "github.com/whosonfirst/wof/geometry"
	_ "github.com/whosonfirst/wof/open"
	_ "github.com/whosonfirst/wof/pip"
	_ "github.com/whosonfirst/wof/property"
	_ "github.com/whosonfirst/wof/show"
	_ "github.com/whosonfirst/wof/uri"
	_ "github.com/whosonfirst/wof/validate"
)

func usage() {

	fmt.Println("Usage: wof [CMD] [OPTIONS]")
	fmt.Println("Valid commands are:")

	for _, cmd := range wof.Commands() {
		fmt.Printf("* %s\n", cmd)
	}

	os.Exit(0)
}

func Run(ctx context.Context) error {

	if len(os.Args) < 2 {
		usage()
	}

	cmd := os.Args[1]

	c, err := wof.NewCommand(ctx, cmd)

	if err != nil {
		usage()
	}

	args := make([]string, 0)

	if len(os.Args) > 2 {
		args = os.Args[2:]
	}

	err = c.Run(ctx, args)

	if err != nil {
		return fmt.Errorf("Failed to run '%s' command, %w", cmd, err)
	}

	return nil
}
