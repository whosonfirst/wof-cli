package show

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"github.com/pkg/browser"
	"github.com/whosonfirst/wof"
	"github.com/whosonfirst/wof/show/www"
	"github.com/whosonfirst/wof/uris"
)

type RunOptions struct {
	URIs []string
	Port int
}

type ShowCommand struct {
	wof.Command
}

func init() {
	ctx := context.Background()
	wof.RegisterCommand(ctx, "show", NewShowCommand)
}

func NewShowCommand(ctx context.Context, cmd string) (wof.Command, error) {

	c := &ShowCommand{}
	return c, nil
}

func (c *ShowCommand) Run(ctx context.Context, args []string) error {

	fs := DefaultFlagSet()
	fs.Parse(args)

	uris := fs.Args()

	opts := &RunOptions{
		URIs: uris,
		Port: port,
	}

	return RunWithOptions(ctx, opts)
}

func RunWithOptions(ctx context.Context, opts *RunOptions) error {

	cb := func(ctx context.Context, uri string) error {
		return nil
	}

	err := uris.ExpandURIsWithCallback(ctx, cb, opts.URIs...)

	if err != nil {
		return fmt.Errorf("Failed to run, %w", err)
	}

	http_fs := http.FS(www.FS)
	fs_handler := http.FileServer(http_fs)

	mux := http.NewServeMux()
	mux.Handle("/", fs_handler)

	addr := fmt.Sprintf("localhost:%d", opts.Port)

	http_server := http.Server{
		Addr: addr,
	}

	http_server.Handler = mux

	done_ch := make(chan bool)

	go func() {

		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		err := http_server.Shutdown(ctx)

		if err != nil {
			slog.Error("HTTP server shutdown error", "error", err)
		}

		close(done_ch)
	}()

	go func() {
		http_server.ListenAndServe()
	}()

	url := fmt.Sprintf("http://%s", addr)

	err = browser.OpenURL(url)

	if err != nil {
		return fmt.Errorf("Failed to open '%s', %w", url, err)
	}

	fmt.Printf("Records are viewable at %s\n", url)

	<-done_ch
	return nil
}
