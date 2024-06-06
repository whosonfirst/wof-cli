package emit

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"strconv"

	"github.com/aaronland/go-json-query"
	"github.com/sfomuseum/go-timings"
	"github.com/tidwall/sjson"
	"github.com/whosonfirst/go-whosonfirst-iterate/v2/emitter"
	"github.com/whosonfirst/go-whosonfirst-iterwriter"
	wof_spr "github.com/whosonfirst/go-whosonfirst-spr/v2"
	"github.com/whosonfirst/go-whosonfirst-uri"
	"github.com/whosonfirst/go-writer/v3"
)

type iterwriterCallbackOptions struct {
	AsSPR           bool
	AsSPRGeoJSON    bool
	Forgiving       bool
	QuerySet        *query.QuerySet
	IncludeAltGeoms bool
}

func iterwriterCallbackFunc(opts *iterwriterCallbackOptions) iterwriter.IterwriterCallbackFunc {

	return func(wr writer.Writer, monitor timings.Monitor) emitter.EmitterCallbackFunc {

		iter_cb := func(ctx context.Context, path string, r io.ReadSeeker, args ...interface{}) error {

			logger := slog.Default()
			logger = logger.With("path", path)

			// See this? It's important. We are rewriting path to a normalized WOF relative path
			// That means this will only work with WOF documents

			id, uri_args, err := uri.ParseURI(path)

			if err != nil {
				slog.Error("Failed to parse URI", "error", err)
				return fmt.Errorf("Unable to parse %s, %w", path, err)
			}

			logger = logger.With("id", id)

			if !opts.IncludeAltGeoms && uri_args.IsAlternate {
				return nil
			}

			logger = logger.With("alt geom", uri_args.IsAlternate)

			rel_path, err := uri.Id2RelPath(id, uri_args)

			if err != nil {
				logger.Error("Failed to derive URI", "error", err)
				return fmt.Errorf("Unable to derive relative (WOF) path for %s, %w", path, err)
			}

			logger = logger.With("rel_path", rel_path)

			if opts.QuerySet != nil {

				body, err := io.ReadAll(r)

				if err != nil {
					logger.Error("Failed to read body", "error", err)
					return fmt.Errorf("Failed to read body for %s, %w", path, err)
				}

				matches, err := query.Matches(ctx, opts.QuerySet, body)

				if err != nil {
					logger.Error("Failed to derive query matches", "error", err)
					return fmt.Errorf("Failed to derive query matches for %s, %w", path, err)
				}

				if !matches {
					return nil
				}

				_, err = r.Seek(0, 0)

				if err != nil {
					logger.Error("Failed to rewind body", "error", err)
					return fmt.Errorf("Failed to rewind body for %s, %w", path, err)
				}
			}

			body_r := r

			if opts.AsSPR {

				body, err := io.ReadAll(r)

				if err != nil {
					logger.Error("Failed to read body to derive SPR", "error", err)
					return fmt.Errorf("Failed to read body for %s, %w", path, err)
				}

				var spr_rsp wof_spr.StandardPlacesResult

				if uri_args.IsAlternate {

					rsp, err := wof_spr.WhosOnFirstAltSPR(body)

					if err != nil {
						logger.Error("Failed to derive SPR", "error", err)

						if !opts.Forgiving {
							return fmt.Errorf("Failed to derive (alt) SPR for %s, %w", path, err)
						}
					}

					spr_rsp = rsp
				} else {

					rsp, err := wof_spr.WhosOnFirstSPR(body)

					if err != nil {
						logger.Error("Failed to derive SPR", "error", err)

						if !opts.Forgiving {
							return fmt.Errorf("Failed to derive SPR for %s, %w", path, err)
						}
					}

					spr_rsp = rsp
				}

				if opts.AsSPRGeoJSON {

					body, err = sjson.SetBytes(body, "properties", spr_rsp)

					if err != nil {
						logger.Error("Failed to assign SPR as properties hash", "error", err)
						return fmt.Errorf("Failed to assign SPR properties for %s, %w", path, err)
					}

					wof_id, err := strconv.ParseInt(spr_rsp.Id(), 10, 64)

					if err != nil {
						logger.Error("Failed to parse SPR ID", "error", err)
						return fmt.Errorf("Failed to parse SPR ID for %s, %w", path, err)
					}

					body, err = sjson.SetBytes(body, "properties.wof:id", wof_id)

					if err != nil {
						logger.Error("Failed to assign wof:id to properties", "error", err)
						return fmt.Errorf("Failed to assign wof:id for %s, %w", path, err)
					}

					body_r = bytes.NewReader(body)

				} else {

					enc_spr, err := json.Marshal(spr_rsp)

					if err != nil {
						logger.Error("Failed to marshal SPR", "error", err)
						return fmt.Errorf("Failed to marshal SPR for %s, %w", path, err)
					}

					body_r = bytes.NewReader(enc_spr)
				}
			}

			_, err = wr.Write(ctx, rel_path, body_r)

			if err != nil {

				logger.Error("Failed to write record %s (%s), %w", rel_path, path, err)

				if !opts.Forgiving {
					return fmt.Errorf("Failed to write record for %s, %w", rel_path, err)
				}
			}

			go monitor.Signal(ctx)
			return nil
		}

		return iter_cb
	}
}
