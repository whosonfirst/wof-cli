package supersede

import (
	"context"
	"fmt"
	_ "log/slog"
	"slices"
	"time"

	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-export/v2"
	"github.com/whosonfirst/go-whosonfirst-feature/properties"
	wof_id "github.com/whosonfirst/go-whosonfirst-id"
	wof_reader "github.com/whosonfirst/go-whosonfirst-reader"
	wof_writer "github.com/whosonfirst/go-whosonfirst-writer/v3"
	"github.com/whosonfirst/go-writer/v3"
	"github.com/whosonfirst/wof"
	cli_reader "github.com/whosonfirst/wof/reader"
	"github.com/whosonfirst/wof/uris"
	cli_writer "github.com/whosonfirst/wof/writer"
)

type SupersedeCommand struct {
	wof.Command
}

func init() {
	ctx := context.Background()
	wof.RegisterCommand(ctx, "supersede", NewSupersedeCommand)
}

func NewSupersedeCommand(ctx context.Context, cmd string) (wof.Command, error) {

	c := &SupersedeCommand{}
	return c, nil
}

func (c *SupersedeCommand) Run(ctx context.Context, args []string) error {

	fs := DefaultFlagSet()
	fs.Parse(args)

	fs_uris := fs.Args()

	var id_provider wof_id.Provider
	var superseding_reader reader.Reader
	var superseding_writer writer.Writer
	var parent_reader reader.Reader

	if superseding_id == -1 {

		pr, err := wof_id.NewProvider(ctx)

		if err != nil {
			return fmt.Errorf("Failed to create new wof:id provider, %w", err)
		}

		id_provider = pr

		r, err := reader.NewReader(ctx, superseding_reader_uri)

		if err != nil {
			return fmt.Errorf("Failed to create superseding reader, %w", err)
		}

		superseding_reader = r

		wr, err := writer.NewWriter(ctx, superseding_writer_uri)

		if err != nil {
			return fmt.Errorf("Failed to create superseding writer, %w", err)
		}

		superseding_writer = wr
	}

	if parent_id != -1 {

		r, err := reader.NewReader(ctx, parent_reader_uri)

		if err != nil {
			return fmt.Errorf("Failed to create parent reader, %w", err)
		}

		parent_reader = r
	}

	cb := func(ctx context.Context, cb_uri string) error {

		body, err := cli_reader.BytesFromURI(ctx, cb_uri)

		if err != nil {
			return fmt.Errorf("Failed to open %s for reading, %w", cb_uri, err)
		}

		old_id, err := properties.Id(body)

		if err != nil {
			return fmt.Errorf("Failed to derive ID for %s, %w", cb_uri, err)
		}

		cb_superseding_id := superseding_id

		now := time.Now()

		// If -1 then clone record and use the new ID as the superseding_id

		if cb_superseding_id == -1 {

			new_id, err := id_provider.NewID(ctx)

			if err != nil {
				return fmt.Errorf("Failed to derive new ID, %w", err)
			}

			updates := map[string]interface{}{
				"id":                           new_id,
				"properties.wof:id":            new_id,
				"properties.edtf:inception":    now.Format("2006-01-02"),
				"properties.edtf:cessation":    "..",
				"properties.mz:is_current":     1,
				"properties.wof:supersedes":    []int64{old_id},
				"properties.wof:superseded_by": []int64{},
			}

			// Read properties from parent?

			if parent_id != -1 {

				parent_body, err := wof_reader.LoadBytes(ctx, parent_reader, parent_id)

				if err != nil {
					return fmt.Errorf("Failed to load record for parent ID, %w", err)
				}

				parent_hierarchies := properties.Hierarchies(parent_body)
				parent_country := properties.Country(parent_body)

				updates["properties.wof:parent_id"] = parent_id
				updates["properties.wof:hierarchy"] = parent_hierarchies
				updates["properties.wof:country"] = parent_country
			}

			// Else, re-PIP?

			_, new_body, err := export.AssignPropertiesIfChanged(ctx, body, updates)

			if err != nil {
				return fmt.Errorf("Failed to assign new properties to new record derived from %s, %w", cb_uri, err)
			}

			_, err = wof_writer.WriteBytes(ctx, superseding_writer, new_body)

			if err != nil {
				return fmt.Errorf("Failed to write new record (%d) derived from %s, %w", new_id, cb_uri, err)
			}

			cb_superseding_id = new_id

		} else {

			superseding_body, err := wof_reader.LoadBytes(ctx, superseding_reader, cb_superseding_id)

			if err != nil {
				return fmt.Errorf("Failed to load record for superseding ID, %w", err)
			}

			supersedes := properties.Supersedes(superseding_body)

			if !slices.Contains(supersedes, old_id) {
				supersedes = append(supersedes, old_id)
			}

			updates := map[string]interface{}{
				"properties.wof:supersede": supersedes,
			}

			has_changes, new_body, err := export.AssignPropertiesIfChanged(ctx, superseding_body, updates)

			if err != nil {
				return fmt.Errorf("Failed to assign new properties to superseding ID for %s, %w", cb_uri, err)
			}

			if has_changes {

				_, err := wof_writer.WriteBytes(ctx, superseding_writer, new_body)

				if err != nil {
					return fmt.Errorf("Failed to write changes to superseding ID for %s, %w", cb_uri, err)
				}
			}
		}

		superseded_by := properties.SupersededBy(body)

		if !slices.Contains(superseded_by, cb_superseding_id) {
			superseded_by = append(superseded_by, cb_superseding_id)
		}

		updates := map[string]interface{}{
			"properties.mz:is_current":     0,
			"properties.edtf:cessation":    now.Format("2006-01-02"),
			"properties.wof:superseded_by": superseded_by,
		}

		if is_deprecated {
			updates["properties.edtf:deprecated"] = now.Format("2006-01-02")
		}

		has_changes, new_body, err := export.AssignPropertiesIfChanged(ctx, body, updates)

		if err != nil {
			return fmt.Errorf("Failed to assign new properties to %s, %w", cb_uri, err)
		}

		if has_changes {

			err := cli_writer.Write(ctx, cb_uri, new_body)

			if err != nil {
				return fmt.Errorf("Failed to write new properties to %s, %w", cb_uri, err)
			}
		}

		return nil
	}

	return uris.ExpandURIsWithCallback(ctx, cb, fs_uris...)
}
