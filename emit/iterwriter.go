package emit

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"sort"
	"strconv"
	"strings"

	"github.com/aaronland/go-json-query"
	"github.com/sfomuseum/go-csvdict/v2"
	"github.com/sfomuseum/go-timings"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"github.com/whosonfirst/go-whosonfirst-iterate/v3"
	"github.com/whosonfirst/go-whosonfirst-iterwriter/v4"
	wof_spr "github.com/whosonfirst/go-whosonfirst-spr/v2"
	wof_spr_util "github.com/whosonfirst/go-whosonfirst-spr/v2/util"
	"github.com/whosonfirst/go-whosonfirst-uri"
	"github.com/whosonfirst/go-writer/v3"
)

type iterwriterCallbackOptions struct {
	Format              string
	Forgiving           bool
	QuerySet            *query.QuerySet
	IncludeAltGeoms     bool
	CSVAppendProperties map[string]string
}

func iterwriterCallbackFunc(opts *iterwriterCallbackOptions) iterwriter.IterwriterCallback {

	csv_header := false

	iter_cb := func(ctx context.Context, rec *iterate.Record) error {

		logger := slog.Default()
		logger = logger.With("path", rec.Path)

		// See this? It's important. We are rewriting path to a normalized WOF relative path
		// That means this will only work with WOF documents

		id, uri_args, err := uri.ParseURI(rec.Path)

		if err != nil {
			slog.Error("Failed to parse URI", "error", err)
			return fmt.Errorf("Unable to parse %s, %w", rec.Path, err)
		}

		logger = logger.With("id", id)

		if !opts.IncludeAltGeoms && uri_args.IsAlternate {
			return nil
		}

		logger = logger.With("alt geom", uri_args.IsAlternate)

		rel_path, err := uri.Id2RelPath(id, uri_args)

		if err != nil {
			logger.Error("Failed to derive URI", "error", err)
			return fmt.Errorf("Unable to derive relative (WOF) path for %s, %w", rec.Path, err)
		}

		logger = logger.With("rel_path", rel_path)

		if opts.QuerySet != nil {

			body, err := io.ReadAll(rec.Body)

			if err != nil {
				logger.Error("Failed to read body", "error", err)
				return fmt.Errorf("Failed to read body for %s, %w", rec.Path, err)
			}

			matches, err := query.Matches(ctx, opts.QuerySet, body)

			if err != nil {
				logger.Error("Failed to derive query matches", "error", err)
				return fmt.Errorf("Failed to derive query matches for %s, %w", rec.Path, err)
			}

			if !matches {
				return nil
			}

			_, err = rec.Body.Seek(0, 0)

			if err != nil {
				logger.Error("Failed to rewind body", "error", err)
				return fmt.Errorf("Failed to rewind body for %s, %w", path, err)
			}
		}

		body_r := rec.Body

		switch opts.Format {
		case "csv", "spr", "spr-geojson":

			body, err := io.ReadAll(rec.Body)

			if err != nil {
				logger.Error("Failed to read body to derive SPR", "error", err)
				return fmt.Errorf("Failed to read body for %s, %w", rec.Path, err)
			}

			var spr_rsp wof_spr.StandardPlacesResult

			if uri_args.IsAlternate {

				rsp, err := wof_spr.WhosOnFirstAltSPR(body)

				if err != nil {
					logger.Error("Failed to derive SPR", "error", err)

					if !opts.Forgiving {
						return fmt.Errorf("Failed to derive (alt) SPR for %s, %w", rec.Path, err)
					}
				}

				spr_rsp = rsp
			} else {

				rsp, err := wof_spr.WhosOnFirstSPR(body)

				if err != nil {
					logger.Error("Failed to derive SPR", "error", err)

					if !opts.Forgiving {
						return fmt.Errorf("Failed to derive SPR for %s, %w", rec.Path, err)
					}
				}

				spr_rsp = rsp
			}

			switch format {
			case "csv":

				spr_row, err := wof_spr_util.SPRToMap(spr_rsp)

				if err != nil {

					logger.Error("Failed to derive SPR map", "error", err)

					if !opts.Forgiving {
						return fmt.Errorf("Failed to derive SPR map for %s, %w", rec.Path, err)
					}
				}

				spr_row["mz:uri"] = strings.Replace(spr_row["mz:uri"], "https://data.whosonfirst.org/", "", 1)

				for col_name, path := range opts.CSVAppendProperties {
					rsp := gjson.GetBytes(body, path)
					spr_row[col_name] = rsp.String()
				}

				var buf bytes.Buffer
				wr := bufio.NewWriter(&buf)

				csv_wr, err := csvdict.NewWriter(wr)

				if err != nil {
					return fmt.Errorf("Failed to create CSV writer for %s, %w", rec.Path, err)
				}

				err = csv_wr.WriteRow(spr_row)

				if err != nil {
					return fmt.Errorf("Failed to write CSV row for %s, %w", rec.Path, err)
				}

				wr.Flush()
				csv_wr.Flush()

				body_r = bytes.NewReader(buf.Bytes())

			case "spr-geojson":

				body, err = sjson.SetBytes(body, "properties", spr_rsp)

				if err != nil {
					logger.Error("Failed to assign SPR as properties hash", "error", err)
					return fmt.Errorf("Failed to assign SPR properties for %s, %w", rec.Path, err)
				}

				wof_id, err := strconv.ParseInt(spr_rsp.Id(), 10, 64)

				if err != nil {
					logger.Error("Failed to parse SPR ID", "error", err)
					return fmt.Errorf("Failed to parse SPR ID for %s, %w", rec.Path, err)
				}

				body, err = sjson.SetBytes(body, "properties.wof:id", wof_id)

				if err != nil {
					logger.Error("Failed to assign wof:id to properties", "error", err)
					return fmt.Errorf("Failed to assign wof:id for %s, %w", rec.Path, err)
				}

				body_r = bytes.NewReader(body)

			default:

				enc_spr, err := json.Marshal(spr_rsp)

				if err != nil {
					logger.Error("Failed to marshal SPR", "error", err)
					return fmt.Errorf("Failed to marshal SPR for %s, %w", rec.Path, err)
				}

				body_r = bytes.NewReader(enc_spr)
			}
		default:
			//
		}

		_, err = wr.Write(ctx, rel_path, body_r)

		if err != nil {

			logger.Error("Failed to write record %s (%s), %w", rel_path, rec.Path, "error", err)

			if !opts.Forgiving {
				return fmt.Errorf("Failed to write record for %s, %w", rel_path, err)
			}
		}

		return nil
	}

	return iter_cb
}
