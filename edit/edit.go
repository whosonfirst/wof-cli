package edit

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/sfomuseum/go-www-show"
	wof_uri "github.com/whosonfirst/go-whosonfirst-uri"
	"github.com/whosonfirst/wof"
	"github.com/whosonfirst/wof/reader"
	"github.com/whosonfirst/wof/uris"
)

type RunOptions struct {
	URIs   []string
	Stdout bool
}

type EditCommand struct {
	wof.Command
}

func init() {
	ctx := context.Background()
	wof.RegisterCommand(ctx, "edit", NewEditCommand)
}

func NewEditCommand(ctx context.Context, cmd string) (wof.Command, error) {

	c := &EditCommand{}

	return c, nil
}

func (c *EditCommand) Run(ctx context.Context, args []string) error {

	fs := DefaultFlagSet()
	fs.Parse(args)

	uris := fs.Args()

	opts := &RunOptions{
		URIs: uris,
	}

	return RunWithOptions(ctx, opts)
}

func RunWithOptions(ctx context.Context, opts *RunOptions) error {

	logger := slog.Default()

	tmpdir, err := os.MkdirTemp("", "wof-edit")

	if err != nil {
		return fmt.Errorf("Failed to create temp dir, %w", err)
	}

	logger = logger.With("tmpdir", tmpdir)
	logger.Info("Tmpdir created")

	defer func() {
		logger.Info("Remove tmpdir")
		os.RemoveAll(tmpdir)
	}()

	root, err := os.OpenRoot(tmpdir)

	if err != nil {
		return fmt.Errorf("Failed to open root, %w", err)
	}

	cb := func(ctx context.Context, uri string) error {

		logger.Info("Process record", "uri", uri)

		id, uri_args, err := wof_uri.ParseURI(uri)

		if err != nil {
			return fmt.Errorf("Failed to parse '%s', %w", uri, err)
		}

		r, is_stdin, err := reader.ReadCloserFromURI(ctx, uri)

		if err != nil {
			return fmt.Errorf("Failed to open '%s' for reading, %w", uri, err)
		}

		if !is_stdin {
			defer r.Close()
		}

		fname, err := wof_uri.Id2Fname(id, uri_args)

		if err != nil {
			return fmt.Errorf("Failed to create fname from URI, %w", err)
		}

		wr, err := root.OpenFile(fname, os.O_RDWR|os.O_CREATE, 0644)

		if err != nil {
			return fmt.Errorf("Failed to open %s in root, %w", fname, err)
		}

		_, err = io.Copy(wr, r)

		if err != nil {
			return fmt.Errorf("Failed to copy %s to root, %w", uri, err)
		}

		err = wr.Close()

		if err != nil {
			return fmt.Errorf("Failed to close %s after writing, %w", uri, err)
		}

		logger.Info("Copy record to tmpdir root", "uri", uri, "fname", fname)
		return nil
	}

	err = uris.ExpandURIsWithCallback(ctx, cb, opts.URIs...)

	if err != nil {
		return fmt.Errorf("Failed to run, %w", err)
	}

	root_fs := root.FS()

	mux := http.NewServeMux()

	list_handler := apiListHandler(root_fs)
	mux.Handle("/api/list", list_handler)

	data_handler := dataHandler(root_fs)
	data_handler = http.StripPrefix("/data/", data_handler)
	mux.Handle("/data/", data_handler)

	static_handler := staticHandler()
	mux.Handle("/", static_handler)

	browser, err := show.NewBrowser(ctx, "web://")

	if err != nil {
		return err
	}

	show_opts := &show.RunOptions{
		Browser: browser,
		Mux:     mux,
	}

	return show.RunWithOptions(ctx, show_opts)

	return nil
}
