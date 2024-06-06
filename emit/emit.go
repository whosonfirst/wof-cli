package emit

/*

> ./bin/wof emit -writer-uri 'jsonl://?writer=stdout://' /usr/local/data/sfomuseum-data-whosonfirst/ | /usr/local/sfomuseum/gpq/bin/gpq convert --from geojson --to geoparquet
2024/06/05 16:07:47 INFO time to index paths (1) 9.355815916s
gpq: error: failed to create schema after reading 1 features

*/

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	_ "os"

	"github.com/sfomuseum/go-timings"
	"github.com/tidwall/sjson"
	"github.com/whosonfirst/go-whosonfirst-iterate/v2/emitter"
	"github.com/whosonfirst/go-whosonfirst-iterwriter"
	app "github.com/whosonfirst/go-whosonfirst-iterwriter/app/iterwriter"
	wof_spr "github.com/whosonfirst/go-whosonfirst-spr/v2"
	"github.com/whosonfirst/go-whosonfirst-uri"
	_ "github.com/whosonfirst/go-writer-featurecollection/v3"
	_ "github.com/whosonfirst/go-writer-jsonl/v3"
	"github.com/whosonfirst/go-writer/v3"
	"github.com/whosonfirst/wof"
)

type EmitCommand struct {
	wof.Command
}

func init() {
	ctx := context.Background()
	wof.RegisterCommand(ctx, "emit", NewEmitCommand)
}

func NewEmitCommand(ctx context.Context, cmd string) (wof.Command, error) {
	c := &EmitCommand{}
	return c, nil
}

func (c *EmitCommand) Run(ctx context.Context, args []string) error {

	fs := DefaultFlagSet()
	fs.Parse(args)

	wr, err := writer.NewWriter(ctx, writer_uri)

	if err != nil {
		return err
	}

	cb_func := iterwriter.DefaultIterwriterCallback

	if forgiving {
		cb_func = iterwriter.ForgivingIterwriterCallback
	}

	if as_spr {

		if as_spr_geojson {
			cb_func = sprGeoJSONIterwriterCallback
		} else {
			cb_func = sprIterwriterCallback
		}
	}

	uris := fs.Args()

	opts := &app.RunOptions{
		Writer:        wr,
		IteratorURI:   iterator_uri,
		IteratorPaths: uris,
		CallbackFunc:  cb_func,
		MonitorURI:    "counter://PT60S",
	}

	logger := slog.Default()
	return app.RunWithOptions(ctx, opts, logger)
}

func sprIterwriterCallback(wr writer.Writer, monitor timings.Monitor) emitter.EmitterCallbackFunc {

	iter_cb := func(ctx context.Context, path string, r io.ReadSeeker, args ...interface{}) error {

		logger := slog.Default()
		logger = logger.With("path", path)

		// See this? It's important. We are rewriting path to a normalized WOF relative path
		// That means this will only work with WOF documents

		id, uri_args, err := uri.ParseURI(path)

		if err != nil {
			slog.Error("Failed to parse URI", "error", err)
			return fmt.Errorf("Unable to parse %s, %w", path, err)
		}

		logger = logger.With("id", id)

		rel_path, err := uri.Id2RelPath(id, uri_args)

		if err != nil {
			slog.Error("Failed to derive URI", "error", err)
			return fmt.Errorf("Unable to derive relative (WOF) path for %s, %w", path, err)
		}

		logger = logger.With("rel_path", rel_path)

		body, err := io.ReadAll(r)

		if err != nil {
			return fmt.Errorf("Failed to read body", "error", err)
		}

		spr_rsp, err := wof_spr.WhosOnFirstSPR(body)

		if err != nil {
			return fmt.Errorf("Failed to derive SPR", "error", err)
		}

		enc_spr, err := json.Marshal(spr_rsp)

		if err != nil {
			return fmt.Errorf("Failed to marshal SPR", "error", err)
		}

		spr_r := bytes.NewReader(enc_spr)

		_, err = wr.Write(ctx, rel_path, spr_r)

		if err != nil {

			slog.Error("Failed to write record", "error", err)
		}

		go monitor.Signal(ctx)
		return nil
	}

	return iter_cb
}

func sprGeoJSONIterwriterCallback(wr writer.Writer, monitor timings.Monitor) emitter.EmitterCallbackFunc {

	iter_cb := func(ctx context.Context, path string, r io.ReadSeeker, args ...interface{}) error {

		logger := slog.Default()
		logger = logger.With("path", path)

		// See this? It's important. We are rewriting path to a normalized WOF relative path
		// That means this will only work with WOF documents

		id, uri_args, err := uri.ParseURI(path)

		if err != nil {
			slog.Error("Failed to parse URI", "error", err)
			return fmt.Errorf("Unable to parse %s, %w", path, err)
		}

		logger = logger.With("id", id)

		rel_path, err := uri.Id2RelPath(id, uri_args)

		if err != nil {
			slog.Error("Failed to derive URI", "error", err)
			return fmt.Errorf("Unable to derive relative (WOF) path for %s, %w", path, err)
		}

		logger = logger.With("rel_path", rel_path)

		body, err := io.ReadAll(r)

		if err != nil {
			return fmt.Errorf("Failed to read body", "error", err)
		}

		spr_rsp, err := wof_spr.WhosOnFirstSPR(body)

		if err != nil {
			return fmt.Errorf("Failed to derive SPR", "error", err)
		}

		body, err = sjson.SetBytes(body, "properties", spr_rsp)

		if err != nil {
			return err
		}

		spr_r := bytes.NewReader(body)

		_, err = wr.Write(ctx, rel_path, spr_r)

		if err != nil {

			slog.Error("Failed to write record", "error", err)
		}

		go monitor.Signal(ctx)
		return nil
	}

	return iter_cb
}
