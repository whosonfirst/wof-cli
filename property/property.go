package property

import (
	"context"
	"fmt"
	_ "log/slog"
	"strings"

	"github.com/sfomuseum/go-csvdict"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-whosonfirst-export/v2"
	"github.com/whosonfirst/wof"
	"github.com/whosonfirst/wof/reader"
	"github.com/whosonfirst/wof/uris"
	"github.com/whosonfirst/wof/writer"
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

	if prefix != "" {

		prefix = strings.TrimRight(prefix, ".")

		for idx, p := range paths {
			paths[idx] = fmt.Sprintf("%s.%s", prefix, p)
		}
	}

	switch action {
	case "remove":
		return c.removeProperties(ctx, fs_uris, paths)
	default:
		return c.listProperties(ctx, fs_uris, paths)
	}

}

func (c *PropertyCommand) listProperties(ctx context.Context, fs_uris []string, paths []string) error {

	var csv_wr *csvdict.Writer

	cb := func(ctx context.Context, cb_uri string) error {

		body, err := reader.BytesFromURI(ctx, cb_uri)

		if err != nil {
			return fmt.Errorf("Failed to open %s for reading, %w", cb_uri, err)
		}

		switch format {
		case "csv":

			out := map[string]string{
				"uri": cb_uri,
			}

			for _, path := range paths {
				rsp := gjson.GetBytes(body, path)
				out[path] = rsp.String()
			}

			if csv_wr == nil {

				fieldnames := make([]string, 0)

				for k, _ := range out {
					fieldnames = append(fieldnames, k)
				}

				wr, err := csvdict.NewWriter(os.Stdout, fieldnames)

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

func (c *PropertyCommand) removeProperties(ctx context.Context, fs_uris []string, paths []string) error {

	ex, err := export.NewExporter(ctx, "whosonfirst://")

	if err != nil {
		return fmt.Errorf("Failed to create new exporter, %w", err)
	}

	cb := func(ctx context.Context, cb_uri string) error {

		body, err := reader.BytesFromURI(ctx, cb_uri)

		if err != nil {
			return fmt.Errorf("Failed to read %s, %w", cb_uri, err)
		}

		new_body, err := export.RemoveProperties(ctx, body, paths)

		if err != nil {
			return fmt.Errorf("Failed to remove properties from %s, %w", cb_uri, err)
		}

		new_body, err = ex.Export(ctx, new_body)

		if err != nil {
			return fmt.Errorf("Failed to export %s, %w", cb_uri, err)
		}

		err = writer.Write(ctx, cb_uri, new_body)

		if err != nil {
			return fmt.Errorf("Failed to write changes to %s, %w", cb_uri, err)
		}

		return nil
	}

	return uris.ExpandURIsWithCallback(ctx, cb, fs_uris...)
}
