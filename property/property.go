package property

import (
	"context"
	"fmt"
	_ "log/slog"

	"github.com/sfomuseum/go-csvdict"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/wof"
	"github.com/whosonfirst/wof/reader"
	"github.com/whosonfirst/wof/uris"
	"os"
)

type PropertyCommand struct {
	wof.Command
}

func init() {
	ctx := context.Background()
	wof.RegisterCommand(ctx, "property", NewPropertyCommand)
}

func NewPropertyCommand(ctx context.Context, cmd string) (wof.Command, error) {

	c := &PropertyCommand{}
	return c, nil
}

func (c *PropertyCommand) Run(ctx context.Context, args []string) error {

	fs := DefaultFlagSet()
	fs.Parse(args)

	fs_uris := fs.Args()

	var csv_wr *csvdict.Writer

	cb := func(ctx context.Context, cb_uri string) error {

		body, err := reader.BytesFromURI(ctx, cb_uri)

		if err != nil {
			return fmt.Errorf("Failed to open %s for reading, %w", cb_uri, err)
		}

		switch format {
		case "csv":

			out := make(map[string]string)

			for _, path := range paths {
				rsp := gjson.GetBytes(body, path)
				out[path] = rsp.String()
			}

			if csv_wr == nil {

				wr, err := csvdict.NewWriter(os.Stdout, paths)

				if err != nil {
					return fmt.Errorf("Failed to create new CSV writer, %w", err)
				}

				err = wr.WriteHeader()

				if err != nil {
					return fmt.Errorf("Failed to write CSV header, %w", err)
				}

				csv_wr = wr
			}

			err = csv_wr.WriteRow(out)

			if err != nil {
				return fmt.Errorf("Failed to write CSV row for %s, %w", cb_uri, err)
			}

			csv_wr.Flush()

		default:

			for _, path := range paths {
				rsp := gjson.GetBytes(body, path)
				fmt.Println(rsp.String())
			}
		}

		return nil
	}

	return uris.ExpandURIsWithCallback(ctx, cb, fs_uris...)
}
