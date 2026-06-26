package supersede

import (
	"context"
	"fmt"
	_ "log/slog"
	"slices"
	"time"

	"github.com/whosonfirst/go-reader/v2"
	"github.com/whosonfirst/go-whosonfirst/v4/export"
	"github.com/whosonfirst/go-whosonfirst/v4/feature/properties"
	wof_reader "github.com/whosonfirst/go-whosonfirst/v4/reader"
	wof_writer "github.com/whosonfirst/go-whosonfirst/v4/writer"
	"github.com/whosonfirst/go-writer/v3"
	"github.com/whosonfirst/wof"
	cli_reader "github.com/whosonfirst/wof/reader"
	"github.com/whosonfirst/wof/uris"
	cli_writer "github.com/whosonfirst/wof/writer"
)

type DeprecateCommand struct {
	wof.Command
}

func init() {
	ctx := context.Background()
	wof.RegisterCommand(ctx, "deprecate", NewDeprecateCommand)
}

func NewDeprecateCommand(ctx context.Context, cmd string) (wof.Command, error) {

	c := &DeprecateCommand{}
	return c, nil
}

func (c *DeprecateCommand) Run(ctx context.Context, args []string) error {

	fs := DefaultFlagSet()
	fs.Parse(args)

	fs_uris := fs.Args()

	var superseded_r reader.Reader
	var superseded_wr writer.Writer

	if superseded_by_id != -1 {

		r, err := reader.NewReader(ctx, superseded_by_reader_uri)

		if err != nil {
			return fmt.Errorf("Failed to create new superseded by reader, %w", err)
		}

		wr, err := writer.NewWriter(ctx, superseded_by_writer_uri)

		if err != nil {
			return fmt.Errorf("Failed to create new superseded by writer, %w", err)
		}

		superseded_r = r
		superseded_wr = wr
	}

	cb := func(ctx context.Context, cb_uri string) error {

		body, err := cli_reader.BytesFromURI(ctx, cb_uri)

		if err != nil {
			return fmt.Errorf("Failed to open %s for reading, %w", cb_uri, err)
		}

		id, err := properties.Id(body)

		if err != nil {
			return fmt.Errorf("Failed to derive ID for %s, %w", cb_uri, err)
		}

		now := time.Now()

		updates := map[string]any{
			"properties.edtf:deprecated": now.Format("2006-01-02"),
			"properties.mz:is_current":   0,
		}

		if superseded_by_id != -1 {

			superseded_by := properties.SupersededBy(body)

			if !slices.Contains(superseded_by, superseded_by_id) {
				superseded_by = append(superseded_by, superseded_by_id)
				updates["properties.wof:superseded_by"] = superseded_by
			}
		}

		has_changes, new_body, err := export.AssignPropertiesIfChanged(ctx, body, updates)

		if err != nil {
			return fmt.Errorf("Failed to assign new properties to new record derived from %s, %w", cb_uri, err)
		}

		if has_changes {

			err := cli_writer.Write(ctx, cb_uri, new_body)

			if err != nil {
				return fmt.Errorf("Failed to deprecate record derived from %s, %w", cb_uri, err)
			}
		}

		if superseded_by_id != -1 {

			superseding_body, err := wof_reader.LoadBytes(ctx, superseded_r, superseded_by_id)

			if err != nil {
				return fmt.Errorf("Failed to load data for superseded by ID %d, %w", superseded_by_id, err)
			}

			supersedes := properties.Supersedes(superseding_body)

			if slices.Contains(supersedes, id) {
				supersedes = append(supersedes, id)
			}

			updates := map[string]any{
				"properties.wof:supersedes": supersedes,
			}

			has_changes, new_body, err := export.AssignPropertiesIfChanged(ctx, superseding_body, updates)

			if err != nil {
				return fmt.Errorf("Failed to assign new properties to new record derived from %s, %w", cb_uri, err)
			}

			if has_changes {

				_, err := wof_writer.WriteBytes(ctx, superseded_wr, new_body)

				if err != nil {
					return fmt.Errorf("Failed to update superseding for %s, %w", cb_uri, err)
				}
			}

		}

		return nil
	}

	return uris.ExpandURIsWithCallback(ctx, cb, fs_uris...)
}
