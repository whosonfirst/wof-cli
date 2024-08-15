package emit

/*

> ./bin/wof emit -writer-uri 'jsonl://?writer=stdout://' /usr/local/data/sfomuseum-data-whosonfirst/ | /usr/local/sfomuseum/gpq/bin/gpq convert --from geojson --to geoparquet
2024/06/05 16:07:47 INFO time to index paths (1) 9.355815916s
gpq: error: failed to create schema after reading 1 features

*/

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aaronland/go-json-query"
	app "github.com/whosonfirst/go-whosonfirst-iterwriter/app/iterwriter"
	_ "github.com/whosonfirst/go-writer-featurecollection/v3"
	_ "github.com/whosonfirst/go-writer-jsonl/v3"
	"github.com/whosonfirst/go-writer/v3"
	"github.com/whosonfirst/wof"
	"github.com/whosonfirst/wof/uris"
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

	defer wr.Close(ctx)

	iter_uris := fs.Args()

	if iterator_uri == "-" {

		uris_expanded := make([]string, 0)

		uris_cb := func(ctx context.Context, uri string) error {
			uris_expanded = append(uris_expanded, uri)
			return nil
		}

		err := uris.ExpandURIsWithCallback(ctx, uris_cb, iter_uris...)

		if err != nil {
			return fmt.Errorf("Failed to expand URIs, %w", err)
		}

		iter_uris = uris_expanded
		iterator_uri = "file://"
	}

	if format == "" {

		if as_spr {
			slog.Warn("-as-spr flag is deprecated. Please use -format spr instead.")
			format = "spr"
		}

		if as_spr_geojson {
			slog.Warn("-as-spr-geojson flag is deprecated. Please use -format geojson instead.")
			format = "geojson"
		}
	}

	iterwr_opts := &iterwriterCallbackOptions{
		Forgiving:       forgiving,
		IncludeAltGeoms: include_alt_geoms,
		Format:          format,
	}

	if format == "csv" {

		append_properties_map := make(map[string]string)

		for _, kv := range csv_append_properties {
			append_properties_map[kv.Key()] = kv.Value().(string)
		}

		iterwr_opts.CSVAppendProperties = append_properties_map
	}

	if len(queries) > 0 {

		qs := &query.QuerySet{
			Queries: queries,
			Mode:    query_mode,
		}

		iterwr_opts.QuerySet = qs
	}

	iterwr_cb := iterwriterCallbackFunc(iterwr_opts)

	opts := &app.RunOptions{
		Writer:        wr,
		IteratorURI:   iterator_uri,
		IteratorPaths: iter_uris,
		CallbackFunc:  iterwr_cb,
		MonitorURI:    "counter://PT60S",
	}

	logger := slog.Default()
	return app.RunWithOptions(ctx, opts, logger)
}
