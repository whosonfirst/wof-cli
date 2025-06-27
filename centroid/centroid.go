package centroid

import (
	"context"
	"fmt"
	_ "log/slog"
	"os"
	"strconv"

	"github.com/sfomuseum/go-csvdict/v2"
	"github.com/whosonfirst/go-whosonfirst-feature/properties"
	"github.com/whosonfirst/wof"
	"github.com/whosonfirst/wof/reader"
	"github.com/whosonfirst/wof/uris"
)

type CentroidCommand struct {
	wof.Command
}

func init() {
	ctx := context.Background()
	wof.RegisterCommand(ctx, "centroid", NewCentroidCommand)
}

func NewCentroidCommand(ctx context.Context, cmd string) (wof.Command, error) {
	c := &CentroidCommand{}
	return c, nil
}

func (c *CentroidCommand) Run(ctx context.Context, args []string) error {

	fs := DefaultFlagSet()
	fs.Parse(args)

	fs_uris := fs.Args()

	var csv_wr *csvdict.Writer

	cb := func(ctx context.Context, cb_uri string) error {

		body, err := reader.BytesFromURI(ctx, cb_uri)

		if err != nil {
			return fmt.Errorf("Failed to open %s for reading, %w", cb_uri, err)
		}

		c, src, err := properties.Centroid(body)

		if err != nil {
			return fmt.Errorf("Failed to derive centroid for %s, %w", cb_uri, err)
		}

		out := map[string]string{
			"uri":       cb_uri,
			"source":    src,
			"latitude":  strconv.FormatFloat(c.Y(), 'g', -1, 64),
			"longitude": strconv.FormatFloat(c.X(), 'g', -1, 64),
		}

		if csv_wr == nil {

			wr, err := csvdict.NewWriter(os.Stdout)

			if err != nil {
				return fmt.Errorf("Failed to create new CSV writer, %w", err)
			}

			csv_wr = wr
		}

		err = csv_wr.WriteRow(out)

		if err != nil {
			return fmt.Errorf("Failed to write CSV row for %s, %w", cb_uri, err)
		}

		csv_wr.Flush()

		return nil
	}

	return uris.ExpandURIsWithCallback(ctx, cb, fs_uris...)
}
