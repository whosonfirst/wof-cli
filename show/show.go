package show

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/paulmach/orb/geojson"
	"github.com/pkg/browser"
	"github.com/whosonfirst/go-whosonfirst-format-wasm/static/javascript"
	"github.com/whosonfirst/go-whosonfirst-format-wasm/static/wasm"
	"github.com/whosonfirst/wof"
	"github.com/whosonfirst/wof/reader"
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

	if port == 0 {

		listener, err := net.Listen("tcp", "localhost:0")

		if err != nil {
			return fmt.Errorf("Failed to determine next available port, %w", err)
		}

		port = listener.Addr().(*net.TCPAddr).Port
	}

	opts := &RunOptions{
		URIs: uris,
		Port: port,
	}

	return RunWithOptions(ctx, opts)
}

func RunWithOptions(ctx context.Context, opts *RunOptions) error {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	fc := geojson.NewFeatureCollection()

	cb := func(ctx context.Context, uri string) error {

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

		f, err := geojson.UnmarshalFeature(body)

		if err != nil {
			return fmt.Errorf("Failed to unmarshal '%s' as GeoJSON, %w", uri, err)
		}

		fc.Append(f)
		return nil
	}

	err := uris.ExpandURIsWithCallback(ctx, cb, opts.URIs...)

	if err != nil {
		return fmt.Errorf("Failed to run, %w", err)
	}

	data_handler := dataHandler(fc)

	http_fs := http.FS(www.FS)
	fs_handler := http.FileServer(http_fs)

	wasm_fs := http.FS(wasm.FS)
	wasm_handler := http.FileServer(wasm_fs)

	javascript_fs := http.FS(javascript.FS)
	javascript_handler := http.FileServer(javascript_fs)

	mux := http.NewServeMux()
	mux.Handle("/features.geojson", data_handler)

	mux.Handle("/javascript/wasm/", http.StripPrefix("/javascript/wasm/", javascript_handler))
	mux.Handle("/wasm/", http.StripPrefix("/wasm/", wasm_handler))

	mux.Handle("/", fs_handler)

	addr := fmt.Sprintf("localhost:%d", opts.Port)
	url := fmt.Sprintf("http://%s", addr)

	http_server := http.Server{
		Addr: addr,
	}

	http_server.Handler = mux

	done_ch := make(chan bool)
	err_ch := make(chan error)

	go func() {

		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		slog.Info("Shutting server down")
		err := http_server.Shutdown(ctx)

		if err != nil {
			slog.Error("HTTP server shutdown error", "error", err)
		}

		close(done_ch)
	}()

	go func() {

		err := http_server.ListenAndServe()

		if err != nil {
			err_ch <- fmt.Errorf("Failed to start server, %w", err)
		}
	}()

	server_ready := false

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case err := <-err_ch:
			return err
		case <-ticker.C:

			rsp, err := http.Head(url)

			if err != nil {
				slog.Warn("HEAD request failed", "url", url, "error", err)
			} else {

				defer rsp.Body.Close()

				if rsp.StatusCode != 200 {
					slog.Warn("HEAD request did not return expected status code", "url", url, "code", rsp.StatusCode)
				} else {
					slog.Debug("HEAD request succeeded", "url", url)
					server_ready = true
				}
			}
		}

		if server_ready {
			break
		}
	}

	err = browser.OpenURL(url)

	if err != nil {
		return fmt.Errorf("Failed to open '%s', %w", url, err)
	}

	fmt.Printf("Records are viewable at %s\n", url)

	<-done_ch
	return nil
}
