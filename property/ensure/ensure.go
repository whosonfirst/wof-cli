package ensure

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"

	_ "github.com/whosonfirst/go-whosonfirst-iterate-reader/v3"

	"github.com/paulmach/orb/encoding/wkt"
	"github.com/paulmach/orb/geojson"
	"github.com/sfomuseum/go-flags/multi"
	"github.com/tidwall/gjson"
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

	var geom *geojson.Geometry

	if geom_property.String() != "" {

		geom_k := geom_property.Key()
		geom_v := geom_property.Value().(string)

		switch geom_k {
		case "wkt":

			orb_geom, err := wkt.Unmarshal(geom_v)

			if err != nil {
				return fmt.Errorf("Failed to unmarshal WKT geometry, %w", err)
			}

			geom = geojson.NewGeometry(orb_geom)

		case "geojson":

			geojson_geom, err := geojson.UnmarshalGeometry([]byte(geom_v))

			if err != nil {
				return fmt.Errorf("Failed to unmarshal GeoJSON geometry, %w", err)
			}

			geom = geojson_geom

		case "file":

			body, err := os.ReadFile(geom_v)

			if err != nil {
				return fmt.Errorf("Failed to read geometry file, %w", err)
			}

			f, err := geojson.UnmarshalFeature(body)

			if err != nil {
				return fmt.Errorf("Failed to unmarshal GeoJSON from geometry file, %w", err)
			}

			orb_geom := f.Geometry
			geom = geojson.NewGeometry(orb_geom)

		default:
			return fmt.Errorf("Invalid or unsupported geom property key")
		}
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

		if str_properties_from != "" && len(str_properties) > 0 {

			rsp := gjson.GetBytes(body, str_properties_from)

			if !rsp.Exists() {
				return fmt.Errorf("Record is missing '%s' property", str_properties_from)
			}

			if rsp.Type != gjson.String {
				return fmt.Errorf("%s property for record is not a string", str_properties_from)
			}

			derived_props := multi.KeyValueString{}

			for _, pr := range str_properties {
				str_fl := fmt.Sprintf("%s=%s", pr.Key(), rsp.String())
				derived_props.Set(str_fl)
			}

			str_properties = derived_props
		}

		if int_properties_from != "" {

			rsp := gjson.GetBytes(body, int_properties_from)

			if !rsp.Exists() {
				return fmt.Errorf("Record is missing '%s' property", int_properties_from)
			}

			if rsp.Type != gjson.Number {
				return fmt.Errorf("%s property for record is not an int", int_properties_from)
			}

			derived_props := multi.KeyValueInt64{}

			for _, pr := range int_properties {
				int_fl := fmt.Sprintf("%s=%d", pr.Key(), rsp.Int())
				derived_props.Set(int_fl)
			}

			int_properties = derived_props
		}

		if float_properties_from != "" {

			rsp := gjson.GetBytes(body, float_properties_from)

			if !rsp.Exists() {
				return fmt.Errorf("Record is missing '%s' property", float_properties_from)
			}

			if rsp.Type != gjson.Number {
				return fmt.Errorf("%s property for record is not a number", float_properties_from)
			}

			derived_props := multi.KeyValueFloat64{}

			for _, pr := range float_properties {
				float_fl := fmt.Sprintf("%s=%v", pr.Key(), rsp.Float())
				derived_props.Set(float_fl)
			}

			float_properties = derived_props
		}

		if bool_properties_from != "" {

			rsp := gjson.GetBytes(body, bool_properties_from)

			if !rsp.Exists() {
				return fmt.Errorf("Record is missing '%s' property", bool_properties_from)
			}

			if rsp.Type != gjson.True && rsp.Type != gjson.False {
				return fmt.Errorf("%s property for record is not a boolean", bool_properties_from)
			}

			derived_props := multi.KeyValueBool{}

			for _, pr := range bool_properties {
				bool_fl := fmt.Sprintf("%s=%v", pr.Key(), rsp.Bool())
				derived_props.Set(bool_fl)
			}

			bool_properties = derived_props
		}

		opts := &update.UpdateFeatureOptions{
			StringProperties:  str_properties,
			Int64Properties:   int_properties,
			Float64Properties: float_properties,
			BooleanProperties: bool_properties,
			IfMissing:         if_missing,
		}

		if geom != nil {
			opts.Geometry = geom
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
