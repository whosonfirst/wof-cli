package remove

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"

	_ "github.com/whosonfirst/go-whosonfirst-iterate-reader/v3"

	"github.com/tidwall/gjson"
	export "github.com/whosonfirst/go-whosonfirst-export/v3"
	"github.com/whosonfirst/go-whosonfirst-iterate/v3"
	uri "github.com/whosonfirst/go-whosonfirst-uri"
	wof_writer "github.com/whosonfirst/go-whosonfirst-writer/v3"
	"github.com/whosonfirst/go-writer/v3"
	"github.com/whosonfirst/wof"
)

type RemovePropertyCommand struct {
	wof.Command
}

func init() {
	ctx := context.Background()
	wof.RegisterCommand(ctx, "remove-property", NewRemovePropertyCommand)
}

func NewRemovePropertyCommand(ctx context.Context, cmd string) (wof.Command, error) {

	c := &RemovePropertyCommand{}
	return c, nil
}

func (c *RemovePropertyCommand) Run(ctx context.Context, args []string) error {

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

		to_remove := make([]string, 0)

		for _, path := range properties {

			if !strings.Contains(path, ".#.") {
				to_remove = append(to_remove, path)
				continue
			}

			parts := strings.Split(path, ".#.")

			if len(parts) != 2 {
				logger.Error("Support for complex parts only supports a single '.#.' condition")
				return fmt.Errorf("Support for complex parts only supports a single '.#.' condition")
			}

			rsp := gjson.GetBytes(body, parts[0])
			count := len(rsp.Array())

			for i := 0; i < count; i++ {
				path = fmt.Sprintf("%s.%d.%s", parts[0], i, parts[1])
				to_remove = append(to_remove, path)
			}
		}

		if len(to_remove) == 0 {
			continue
		}

		has_changes, new_body, err := export.RemovePropertiesIfChanged(ctx, body, to_remove)

		if err != nil {
			logger.Error("Failed to remove feature properties", "error", err)
			return err
		}

		if !has_changes {
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
