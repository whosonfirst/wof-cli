package ensure

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	_ "github.com/whosonfirst/go-whosonfirst-iterate-reader/v3"

	export "github.com/whosonfirst/go-whosonfirst-export/v3"
	"github.com/whosonfirst/go-whosonfirst-iterate/v3"
	uri "github.com/whosonfirst/go-whosonfirst-uri"
	wof_writer "github.com/whosonfirst/go-whosonfirst-writer/v3"
	"github.com/whosonfirst/go-writer/v3"
	"github.com/whosonfirst/wof"
	"github.com/whosonfirst/wof/internal/update"
)

type EnsurePropertyCommand struct {
	wof.Command
}

func init() {
	ctx := context.Background()
	wof.RegisterCommand(ctx, "ensure-property", NewEnsurePropertyCommand)
}

func NewEnsurePropertyCommand(ctx context.Context, cmd string) (wof.Command, error) {

	c := &EnsurePropertyCommand{}
	return c, nil
}

func (c *EnsurePropertyCommand) Run(ctx context.Context, args []string) error {

	fs := DefaultFlagSet()
	fs.Parse(args)

	fs_uris := fs.Args()

	if verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("Verbose logging enabled")
	}
	
	ex, err := export.NewExporter(ctx, exporter_uri)

	if err != nil {
		return fmt.Errorf("Failed to create exporter for '%s', %v", exporter_uri, err)
	}

	wr, err := writer.NewWriter(ctx, writer_uri)

	if err != nil {
		return fmt.Errorf("Failed to create writer for '%s', %v", writer_uri, err)
	}

	iter, err := iterate.NewIterator(ctx, iterator_uri)

	if err != nil {
		return fmt.Errorf("Failed to create new iterator, %w", err)
	}

	for rec, err := range iter.Iterate(ctx, fs_uris...) {

		if err != nil {
			return fmt.Errorf("Iterator returned an error, %w", err)
		}

		logger := slog.Default()
		logger = logger.With("path", rec.Path)

		defer rec.Body.Close()

		_, uri_args, err := uri.ParseURI(rec.Path)

		if err != nil {
			return fmt.Errorf("Failed to parse URI, %w", err)
		}

		if uri_args.IsAlternate {
			logger.Debug("Alternate files are not supported yet, skipping")
			continue
		}

		body, err := io.ReadAll(rec.Body)

		if err != nil {
			logger.Error("Failed to read body", "error", err)
			return err
		}

		opts := &update.UpdateFeatureOptions{
			StringProperties:  str_properties,
			Int64Properties:   int_properties,
			Float64Properties: float_properties,
			BooleanProperties: bool_properties,
		}

		has_changes, new_body, err := update.UpdateFeature(ctx, body, opts)

		if err != nil {
			logger.Error("Failed to update feature properties", "error", err)
			return err
		}

		if !has_changes {
			logger.Debug("No changes, skipping")
			continue
		}

		_, new_body, err = ex.Export(ctx, new_body)

		if err != nil {
			logger.Error("Failed to export updated record", "error", err)
			return err
		}
		
		_, err = wof_writer.WriteBytes(ctx, wr, new_body)

		if err != nil {
			logger.Error("Failed to write updates", "error", err)
			return err
		}

		logger.Info("Updated record")
	}

	return nil
}
