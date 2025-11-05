package httpd

import (
	"context"
	"fmt"
	_ "log/slog"
	go_http "net/http"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/aaronland/go-http-sanitize"
	"github.com/aaronland/go-pagination"
	"github.com/aaronland/go-pagination/countable"
	"github.com/aaronland/go-pagination/cursor"
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-spelunker"
	"github.com/whosonfirst/go-whosonfirst-spelunker-httpd/webfinger"
	"github.com/whosonfirst/go-whosonfirst-spr/v2"
	wof_uri "github.com/whosonfirst/go-whosonfirst-uri"
)

var re_path_id = regexp.MustCompile(`/id/(\d+)/.*$`)

type URI struct {
	Id          int64
	URI         string
	Feature     []byte
	URIArgs     *wof_uri.URIArgs
	IsAlternate bool
}

func ParseURIFromRequest(req *go_http.Request, r reader.Reader) (*URI, error, int) {

	ctx := req.Context()

	path, err := sanitize.GetString(req, "id")

	if err != nil {
		return nil, fmt.Errorf("Failed to derive ?id= parameter, %w", err), go_http.StatusBadRequest
	}

	resource, err := sanitize.GetString(req, "resource")

	if err != nil {
		return nil, fmt.Errorf("Failed to derive ?resource= parameter, %w", err), go_http.StatusBadRequest
	}

	if path == "" && resource != "" {

		wof_uri, err := webfinger.DeriveWhosOnFirstURIFromResource(resource)

		if err != nil {

		}

		path = wof_uri
	}

	// Y U NO WORK...
	// https://pkg.go.dev/net/http@master#hdr-Patterns

	if path == "" {
		path = req.PathValue("id")
	}

	// Oh well, at least the ServeMux recognizes wildcards now...
	if path == "" {

		path = req.URL.Path

		if re_path_id.MatchString(path) {
			m := re_path_id.FindStringSubmatch(path)
			path = m[1]
		}
	}

	return ParseURIFromPath(ctx, path, r)
}

func ParseURIFromPath(ctx context.Context, path string, r reader.Reader) (*URI, error, int) {

	wofid, uri_args, err := wof_uri.ParseURI(path)

	if err != nil {
		return nil, fmt.Errorf("Error locating record for %s, %w", path, err), go_http.StatusNotFound
	}

	if wofid == -1 {
		return nil, fmt.Errorf("Failed to locate record for %s", path), go_http.StatusNotFound
	}

	/*
		rel_path, err := wof_uri.Id2RelPath(wofid, uri_args)

		if err != nil {
			return nil, fmt.Errorf("Failed to derive relative path from %d (%s), %w", wofid, path, err), go_http.StatusBadRequest // StatusInternalServerError
		}

		fh, err := r.Read(ctx, rel_path)

		if err != nil {
			return nil, fmt.Errorf("Failed to read %s, %w", rel_path, err), go_http.StatusBadRequest // StatusInternalServerError
		}

		f, err := io.ReadAll(fh)

		if err != nil {
			return nil, fmt.Errorf("Failed to read feature for %s, %w", rel_path, err), go_http.StatusInternalServerError
		}

	*/

	fname, err := wof_uri.Id2Fname(wofid, uri_args)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive filename from %d (%s), %w", wofid, path, err), go_http.StatusInternalServerError
	}

	ext := filepath.Ext(fname)
	fname = strings.Replace(fname, ext, "", 1)

	parsed_uri := &URI{
		Id:      wofid,
		URI:     fname,
		URIArgs: uri_args,
		// Feature:     f,
		IsAlternate: uri_args.IsAlternate,
	}

	return parsed_uri, nil, 0
}

func PaginationOptionsFromRequest(req *go_http.Request) (pagination.Options, error) {

	q_cursor, err := sanitize.GetString(req, "cursor")

	if err != nil {
		return nil, fmt.Errorf("Failed to derive ?cursor= parameter, %w", err)
	}

	if q_cursor != "" {

		pg_opts, err := cursor.NewCursorOptions()

		if err != nil {
			return nil, fmt.Errorf("Failed to create cursor options, %w", err)
		}

		pg_opts.Pointer(q_cursor)
		return pg_opts, nil
	}

	page, err := sanitize.GetInt64(req, "page")

	if err != nil {
		return nil, fmt.Errorf("Failed to derive ?page= parameter, %w", err)
	}

	if page == 0 {
		page = 1
	}

	pg_opts, err := countable.NewCountableOptions()

	if err != nil {
		return nil, fmt.Errorf("Failed to create countable options, %w", err)
	}

	pg_opts.Pointer(page)
	return pg_opts, nil
}

func ParsePageNumberFromRequest(req *go_http.Request) (int64, error) {

	page, err := sanitize.GetInt64(req, "page")

	if err != nil {
		return 0, fmt.Errorf("Failed to derive ?page= parameter, %w", err)
	}

	if page == 0 {
		page = 1
	}

	return page, nil
}

func FeatureFromRequestURI(ctx context.Context, sp spelunker.Spelunker, req_uri *URI) ([]byte, error) {

	wof_id := req_uri.Id

	f, err := sp.GetFeatureForId(ctx, wof_id, req_uri.URIArgs)

	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve feature for %d, %w", wof_id, err)
	}

	return f, nil
}

func RecordFromRequestURI(ctx context.Context, sp spelunker.Spelunker, req_uri *URI) ([]byte, error) {

	wof_id := req_uri.Id

	f, err := sp.GetRecordForId(ctx, wof_id, req_uri.URIArgs)

	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve record for %d, %w", wof_id, err)
	}

	return f, nil
}

func SPRFromRequestURI(ctx context.Context, sp spelunker.Spelunker, req_uri *URI) (spr.StandardPlacesResult, error) {

	wof_id := req_uri.Id

	f, err := sp.GetSPRForId(ctx, wof_id, req_uri.URIArgs)

	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve record for %d, %w", wof_id, err)
	}

	return f, nil
}
