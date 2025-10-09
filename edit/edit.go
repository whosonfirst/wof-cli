package edit

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	_ "path/filepath"
	"sync"

	_ "github.com/whosonfirst/go-reader-findingaid/v2"
	_ "github.com/whosonfirst/go-reader-github/v2"
	_ "github.com/whosonfirst/go-writer-github/v3"

	"github.com/sfomuseum/go-www-show"
	go_reader "github.com/whosonfirst/go-reader/v2"
	wof_uri "github.com/whosonfirst/go-whosonfirst-uri"
	go_writer "github.com/whosonfirst/go-writer/v3"
	"github.com/whosonfirst/wof"
	"github.com/whosonfirst/wof/edit/http/api"
	"github.com/whosonfirst/wof/edit/http/www"
	"github.com/whosonfirst/wof/edit/static"
	_ "github.com/whosonfirst/wof/reader"
	"github.com/whosonfirst/wof/uris"
	_ "github.com/whosonfirst/wof/writer"
)

type RunOptions struct {
	URIs    []string
	Verbose bool
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
		URIs:    uris,
		Verbose: verbose,
	}

	return RunWithOptions(ctx, opts)
}

func RunWithOptions(ctx context.Context, opts *RunOptions) error {

	if opts.Verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("Verbose logging enabled")
	}

	logger := slog.Default()

	// START OF copy all the records to a temporary directory which
	// will then be used to serve an os.Root instance from the web server
	tmpdir, err := os.MkdirTemp("", "wof-edit")

	if err != nil {
		return fmt.Errorf("Failed to create temp dir, %w", err)
	}

	logger = logger.With("tmpdir", tmpdir)
	logger.Debug("Tmpdir created")

	defer func() {
		logger.Debug("Remove tmpdir")
		os.RemoveAll(tmpdir)
	}()

	root, err := os.OpenRoot(tmpdir)

	if err != nil {
		return fmt.Errorf("Failed to open root, %w", err)
	}

	rdr, err := go_reader.NewReader(ctx, reader_uri)

	if err != nil {
		return fmt.Errorf("Failed to create reader, %w", err)
	}

	wtr, err := go_writer.NewWriter(ctx, writer_uri)

	if err != nil {
		return fmt.Errorf("Failed to create writer, %w", err)
	}

	uri_map := new(sync.Map)

	cb := func(ctx context.Context, uri string) error {

		logger.Debug("Process record", "uri", uri)

		id, uri_args, err := wof_uri.ParseURI(uri)

		if err != nil {
			return fmt.Errorf("Failed to parse '%s', %w", uri, err)
		}

		fname, err := wof_uri.Id2Fname(id, uri_args)

		if err != nil {
			return fmt.Errorf("Failed to create fname from URI, %w", err)
		}

		r, err := rdr.Read(ctx, uri)

		if err != nil {
			return err
		}

		defer r.Close()

		/*
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
		*/

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

		uri_map.Store(fname, uri)
		logger.Debug("Copy record to tmpdir root", "uri", uri, "fname", fname)
		return nil
	}

	err = uris.ExpandURIsWithCallback(ctx, cb, opts.URIs...)

	if err != nil {
		return fmt.Errorf("Failed to run, %w", err)
	}

	// END OF copy all the records to a temporary directory

	// See this? This is not the way the rest of the wof-cli package works (at least
	// not yet) but we're doing it this way so that the http/mux code (below) can be
	// used by other projects and can just expect to write changes to a generic Writer
	// implementation.

	// START OF make this a function... maybe?

	mux := http.NewServeMux()

	list_handler := api.ListHandler(root)
	mux.Handle("/api/list", list_handler)

	save_handler := api.SaveHandler(root, uri_map, wtr)
	mux.Handle("/api/save/", save_handler)

	data_handler := www.DataHandler(root)
	data_handler = http.StripPrefix("/data/", data_handler)
	mux.Handle("/data/", data_handler)

	// The actual web application (HTML + CSS + JavaScript)

	static_handler := www.StaticHandler(static.FS)
	mux.Handle("/", static_handler)

	// END OF make this a function... maybe?

	browser, err := show.NewBrowser(ctx, "web://")

	if err != nil {
		return err
	}

	show_opts := &show.RunOptions{
		Browser: browser,
		Mux:     mux,
	}

	return show.RunWithOptions(ctx, show_opts)
}
