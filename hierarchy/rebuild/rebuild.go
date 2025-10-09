package rebuild

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"sync"

	_ "github.com/whosonfirst/go-whosonfirst-iterate-reader/v3"

	"github.com/whosonfirst/go-reader/v2"
	export "github.com/whosonfirst/go-whosonfirst-export/v3"
	"github.com/whosonfirst/go-whosonfirst-feature/properties"
	"github.com/whosonfirst/go-whosonfirst-iterate/v3"
	wof_reader "github.com/whosonfirst/go-whosonfirst-reader/v2"
	uri "github.com/whosonfirst/go-whosonfirst-uri"
	wof_writer "github.com/whosonfirst/go-whosonfirst-writer/v3"
	"github.com/whosonfirst/go-writer/v3"
	"github.com/whosonfirst/wof"
)

type RebuildHierarchyCommand struct {
	wof.Command
}

func init() {
	ctx := context.Background()
	wof.RegisterCommand(ctx, "rebuild-hierarchy", NewRebuildHierarchyCommand)
}

func NewRebuildHierarchyCommand(ctx context.Context, cmd string) (wof.Command, error) {

	c := &RebuildHierarchyCommand{}
	return c, nil
}

func (c *RebuildHierarchyCommand) Run(ctx context.Context, args []string) error {

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

	parent_r, err := reader.NewReader(ctx, parent_reader_uri)

	if err != nil {
		return fmt.Errorf("Failed to create new (parent) reader, %w", err)
	}

	parent_map := new(sync.Map)

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

		// Do stuff here

		var with_parents []int64

		if len(parent_ids) > 0 {
			with_parents = parent_ids
		} else {
			parent_id, err := properties.ParentId(body)

			if err != nil {
				logger.Error("Failed to derive parent ID for record", "error", err)
				return err
			}

			if parent_id < 0 {
				logger.Warn("Record has special-case parent ID, skipping. Please run again with one or more -parent-id flags", "parent id", parent_id)
				return nil
			}

			with_parents = []int64{parent_id}
		}

		new_hiers := make([]map[string]int64, 0)

		for _, parent_id := range with_parents {

			var parent_hiers []map[string]int64
			v, exists := parent_map.Load(parent_id)

			if !exists {

				parent_body, err := wof_reader.LoadBytes(ctx, parent_r, parent_id)

				if err != nil {
					logger.Error("Failed to load body for parent", "parent id", parent_id, "error", err)
					return err
				}

				parent_hiers = properties.Hierarchies(parent_body)
			} else {
				parent_hiers = v.([]map[string]int64)
				parent_map.Store(parent_id, parent_hiers)
			}

			for _, h := range parent_hiers {
				new_hiers = append(new_hiers, h)
			}
		}

		updates := map[string]any{
			"properties.wof:hierarchy": new_hiers,
		}

		has_changes, new_body, err := export.AssignPropertiesIfChanged(ctx, body, updates)

		if !has_changes {
			logger.Debug("No changes after update")
			return nil
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
