package supersede

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"time"

	"github.com/whosonfirst/go-whosonfirst-export/v2"
	"github.com/whosonfirst/go-whosonfirst-feature/properties"
	wof_id "github.com/whosonfirst/go-whosonfirst-id"
	"github.com/whosonfirst/wof"
	"github.com/whosonfirst/wof/reader"
	"github.com/whosonfirst/wof/uris"
	"github.com/whosonfirst/wof/writer"
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

	if superseding_id == -1 {

		pr, err := wof_id.NewProvider(ctx)

		if err != nil {
			return fmt.Errorf("Failed to create new wof:id provider, %w", err)
		}

		id_provider = pr
	}

	cb := func(ctx context.Context, cb_uri string) error {

		body, err := reader.BytesFromURI(ctx, cb_uri)

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
				"id": new_id,
				"properties.wof:id": new_id,
				"properties.edtf:inception":    now.Format("2006-01-02"),
				"properties.edtf:cessation":    "..",
				"properties.mz:is_current":     1,
				"properties.wof:supersedes":    []int64{old_id},
				"properties.wof:superseded_by": []int64{},
			}

			_, new_body, err := export.AssignPropertiesIfChanged(ctx, body, updates)

			if err != nil {
				return fmt.Errorf("Failed to assign new properties to new record derived from %s, %w", cb_uri, err)
			}

			// write stuff TBD...
			slog.Debug(string(new_body))

			cb_superseding_id = new_id
		} else {
			
			// Load other ID and assign old_id to wof:supersedes
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

			err := writer.Write(ctx, cb_uri, new_body)

			if err != nil {
				return fmt.Errorf("Failed to write new properties to %s, %w", cb_uri, err)
			}
		}

		return nil
	}

	return uris.ExpandURIsWithCallback(ctx, cb, fs_uris...)
}
