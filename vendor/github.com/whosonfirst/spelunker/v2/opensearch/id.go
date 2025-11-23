package opensearch

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"

	opensearchapi "github.com/opensearch-project/opensearch-go/v4/opensearchapi"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-reader-cachereader/v2"
	"github.com/whosonfirst/go-reader/v2"
	wof_spr "github.com/whosonfirst/go-whosonfirst-spr/v2"
	"github.com/whosonfirst/go-whosonfirst-uri"
)

// GetRecordForId retrieves properties (or more specifically the "document") for a given ID in an OpenSearchSpelunker index.
func (s *OpenSearchSpelunker) GetRecordForId(ctx context.Context, id int64, uri_args *uri.URIArgs) ([]byte, error) {

	q := fmt.Sprintf(`{"query": { "ids": { "values": [ %d ] } } }`, id)

	req := &opensearchapi.SearchReq{
		Indices: []string{
			s.index,
		},
		Body: strings.NewReader(q),
	}

	body, err := s.searchWithIndex(ctx, req)

	if err != nil {
		slog.Error("Get by ID query failed", "q", q)
		return nil, fmt.Errorf("Failed to retrieve %d, %w", id, err)
	}

	r := gjson.GetBytes(body, "hits.hits.0._source")

	if !r.Exists() {
		return nil, fmt.Errorf("First hit missing for ID '%d'", id)
	}

	return []byte(r.String()), nil
}

// GetSPRForId retrieves the `spr.StandardPlaceResult` instance for a given ID in an OpenSearchSpelunker index.
func (s *OpenSearchSpelunker) GetSPRForId(ctx context.Context, id int64, uri_args *uri.URIArgs) (wof_spr.StandardPlacesResult, error) {

	r, err := s.GetRecordForId(ctx, id, uri_args)

	if err != nil {
		return nil, err
	}

	return NewSpelunkerRecordSPR(r)
}

// GetFeatureForId retrieves the GeoJSON Feature record for a given ID in an OpenSearchSpelunker index.
func (s *OpenSearchSpelunker) GetFeatureForId(ctx context.Context, id int64, uri_args *uri.URIArgs) ([]byte, error) {

	rel_path, err := uri.Id2RelPath(id, uri_args)

	if err != nil {
		return nil, err
	}

	f_reader := s.reader

	if f_reader == nil {

		record, err := s.GetRecordForId(ctx, id, uri_args)

		if err != nil {
			return nil, err
		}

		repo_name := gjson.GetBytes(record, "wof:repo")
		reader_uri := fmt.Sprintf("https://raw.githubusercontent.com/whosonfirst-data/%s/master/data", repo_name)

		r, err := reader.NewReader(ctx, reader_uri)

		if err != nil {
			return nil, err
		}

		f_reader = r
	}

	r := f_reader

	if s.cache != nil {

		cr_opts := &cachereader.CacheReaderOptions{
			Reader: f_reader,
			Cache:  s.cache,
		}

		cr, err := cachereader.NewCacheReaderWithOptions(ctx, cr_opts)

		if err != nil {
			return nil, fmt.Errorf("Failed to create cache reader, %w", err)
		}

		r = cr
	}

	rsp, err := r.Read(ctx, rel_path)

	if err != nil {
		return nil, err
	}

	defer rsp.Close()
	return io.ReadAll(rsp)
}
