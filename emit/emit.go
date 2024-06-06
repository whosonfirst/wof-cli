package emit

/*

> ./bin/wof emit -writer-uri 'jsonl://?writer=stdout://' /usr/local/data/sfomuseum-data-whosonfirst/ | /usr/local/sfomuseum/gpq/bin/gpq convert --from geojson --to geoparquet
2024/06/05 16:07:47 INFO time to index paths (1) 9.355815916s
gpq: error: failed to create schema after reading 1 features

*/

import (
	"context"
	"log/slog"

	"github.com/aaronland/go-json-query"
	app "github.com/whosonfirst/go-whosonfirst-iterwriter/app/iterwriter"
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

	defer wr.Close(ctx)

	iterwr_opts := &iterwriterCallbackOptions{
		AsSPR:           as_spr,
		AsSPRGeoJSON:    as_spr_geojson,
		Forgiving:       forgiving,
		IncludeAltGeoms: include_alt_geoms,
	}

	if len(queries) > 0 {

		qs := &query.QuerySet{
			Queries: queries,
			Mode:    query_mode,
		}

		iterwr_opts.QuerySet = qs
	}

	iterwr_cb := iterwriterCallbackFunc(iterwr_opts)

	uris := fs.Args()

	opts := &app.RunOptions{
		Writer:        wr,
		IteratorURI:   iterator_uri,
		IteratorPaths: uris,
		CallbackFunc:  iterwr_cb,
		MonitorURI:    "counter://PT60S",
	}

	logger := slog.Default()
	return app.RunWithOptions(ctx, opts, logger)
}
