package edit

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"sync"
	"strings"
	
	_ "github.com/whosonfirst/go-reader-findingaid/v2"
	_ "github.com/whosonfirst/go-reader-github/v2"

	"github.com/aaronland/gocloud/runtimevar"
	"github.com/sfomuseum/go-www-show"
	go_reader "github.com/whosonfirst/go-reader/v2"
	wof_uri "github.com/whosonfirst/go-whosonfirst-uri"
	gh_writer "github.com/whosonfirst/go-writer-github/v3"
	"github.com/whosonfirst/wof"
	"github.com/whosonfirst/wof/edit/http/api"
	"github.com/whosonfirst/wof/edit/http/www"
	"github.com/whosonfirst/wof/edit/static"
	"github.com/whosonfirst/wof/reader"
	"github.com/whosonfirst/wof/uris"
	edit_writer "github.com/whosonfirst/wof/writer"
)

const WOF_FINDINGAID_READER_URI string = "wof-findingaid://"
const WOF_PR_WRITER_URI string = "wof-pr://"

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

	if reader_uri == WOF_FINDINGAID_READER_URI {

		logger.Debug("Automatically configuring data.whosonfirst.org reader URI")

		reader_uri = "findingaid://https/data.whosonfirst.org/findingaid?template=https://raw.githubusercontent.com/whosonfirst-data/{repo}/master/data/"
		ensure_rel_path = true
	}

	if strings.HasPrefix(writer_uri, WOF_PR_WRITER_URI) {

		logger.Debug("Automatically configuring whosonfirst-data Github API PR writer URI")
		
		if access_token_uri == "" {
			return fmt.Errorf("-gh-access-token-uri may not be empty if -writer-uri flag is '%s'", WOF_PR_WRITER_URI)
		}

		pr_access_token, err := runtimevar.StringVar(ctx, access_token_uri)

		if err != nil {
			return fmt.Errorf("Failed to derive access token for use with %s writer, %w", WOF_PR_WRITER_URI, err)
		}

		pr_access_token = strings.TrimSpace(pr_access_token)
			
		wr_q := url.Values{}
		wr_q.Set("access_token", pr_access_token)
		wr_q.Set("prefix", "data")
		wr_q.Set("pr-branch", "wof-cli-edit-{UUID}")
		wr_q.Set("pr-title", "Update {PLACETYPE} {NAME}")
		wr_q.Set("pr-description", "Updating {URI} using wof-cli edit")
		wr_q.Set("pr-ensure-repo", "true")

		wr_u := url.URL{}
		wr_u.Scheme = gh_writer.GITHUBAPI_PR_SCHEME
		wr_u.Host = "whosonfirst-data"
		wr_u.Path = "{REPO}"
		wr_u.RawQuery = wr_q.Encode()

		writer_uri = wr_u.String()
		ensure_rel_path = true
	}

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
		err := os.RemoveAll(tmpdir)

		if err != nil {
			logger.Error("Failed to remove tmpdir on program exit", "error", err)
		}
	}()

	root, err := os.OpenRoot(tmpdir)

	if err != nil {
		return fmt.Errorf("Failed to open root, %w", err)
	}

	// START OF can we update all the tools in wof-cli to use this
	// (wrapped in a function or something)

	// START OF hoop-jumping around readers and writers

	// The "tl;dr" is that everything else in the wof-cli package
	// uses the internal "reader" and "writer" package to process
	// files on the local filesystem or STDIN/STDIN without any
	// additional syntax. The reality of the "edit" tool is that
	// people may want to simply read documents from alternate sources
	// (like data.whosonfirst.org) and/or write them directly back
	// to a GitHub PR so this is the kind of thing we need to do.
	// The internal "writer" package already implements the go-writer
	// interface but the "reader" package does not (yet) implement
	// the go-reader interface. That may happen shortly but for now
	// this is how things are handled.

	var rdr go_reader.Reader

	if reader_uri != "" {

		if writer_uri == "" {
			return fmt.Errorf("-writer-uri must be specified if -reader-uri is not-empty")
		}

		rdr, err = go_reader.NewReader(ctx, reader_uri)

		if err != nil {
			return fmt.Errorf("Failed to create reader, %w", err)
		}
	}

	if writer_uri == "" {
		writer_uri = fmt.Sprintf("%s://", edit_writer.WRITER_SCHEME)
	}

	// (NOT QUITE) END OF hoop-jumping around readers and writers

	uri_map := new(sync.Map)

	cb := func(ctx context.Context, uri string) error {

		logger := slog.Default()
		logger = logger.With("uri", uri)
		logger.Debug("Process record")

		id, uri_args, err := wof_uri.ParseURI(uri)

		if err != nil {
			return fmt.Errorf("Failed to parse '%s', %w", uri, err)
		}

		fname, err := wof_uri.Id2Fname(id, uri_args)

		if err != nil {
			return fmt.Errorf("Failed to create fname from URI, %w", err)
		}

		// START OF MORE hoop-jumping around readers and writers
		// See note above wrt/ hoop-jumping. If the internal "reader"
		// package implemented the go-reader interface then we could
		// get rid of some of this code.

		var uri_r io.ReadCloser

		if rdr != nil {

			if ensure_rel_path {

				rel_path, err := wof_uri.Id2RelPath(id, uri_args)

				if err != nil {
					return fmt.Errorf("Failed to derive relative path for %s, %w", uri, err)
				}

				uri = rel_path
			}

			logger.Debug("Read URI", "uri", uri)
			uri_r, err = rdr.Read(ctx, uri)

			if err != nil {
				return err
			}

			defer uri_r.Close()

		} else {

			r, is_stdin, err := reader.ReadCloserFromURI(ctx, uri)

			if err != nil {
				return fmt.Errorf("Failed to open '%s' for reading, %w", uri, err)
			}

			if !is_stdin {
				defer r.Close()
			}

			uri_r = r
		}

		// END OF MORE hoop-jumping around readers and writers
		// END OF hoop-jumping around readers and writers

		uri_wr, err := root.OpenFile(fname, os.O_RDWR|os.O_CREATE, 0644)

		if err != nil {
			return fmt.Errorf("Failed to open %s in root, %w", fname, err)
		}

		_, err = io.Copy(uri_wr, uri_r)

		if err != nil {
			return fmt.Errorf("Failed to copy %s to root, %w", uri, err)
		}

		err = uri_wr.Close()

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

	// END OF can we update all the tools in wof-cli to use this
	// (modulo the copying files around part...

	// END OF copy all the records to a temporary directory

	// See this? This is not the way the rest of the wof-cli package works (at least
	// not yet) but we're doing it this way so that the http/mux code (below) can be
	// used by other projects and can just expect to write changes to a generic Writer
	// implementation.

	// START OF make this a function... maybe?
	// Remember (because I've forgotten at least once already) that "root" and
	// "uri_map" are local clones of the data being edited, independent of its
	// source, and "wtr" is the tool we use to write that data to a target which may
	// or may not be the same as the source.

	mux := http.NewServeMux()

	list_handler := api.ListHandler(root)
	mux.Handle("/api/list", list_handler)

	slog.Info("Derbug", "writer uri", writer_uri)
	save_handler := api.SaveHandler(root, uri_map, writer_uri)
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
