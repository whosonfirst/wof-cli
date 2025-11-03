package indexer

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log/slog"
	"sync"

	database_sql "github.com/sfomuseum/go-database/sql"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-reader/v2"
	wof_tables "github.com/whosonfirst/go-whosonfirst-database/sql/tables"
	"github.com/whosonfirst/go-whosonfirst-feature/geometry"
	"github.com/whosonfirst/go-whosonfirst-feature/properties"
	"github.com/whosonfirst/go-whosonfirst-uri"
)

// LoadRecordFuncOptions is a struct to define options when loading Who's On First feature records.
type LoadRecordFuncOptions struct {
	// StrictAltFiles is a boolean flag indicating whether the failure to load or parse an alternate geometry file should trigger a critical error.
	StrictAltFiles bool
}

// IndexRelationsFuncOptions
type IndexRelationsFuncOptions struct {
	// Reader is a valid `whosonfirst/go-reader` instance used to load Who's On First feature data
	Reader reader.Reader
	// Strict is a boolean flag indicating whether the failure to load or parse feature record should trigger a critical error.
	Strict bool
}

// LoadRecordFunc returns a `go-whosonfirst-sql/indexer/IndexerLoadRecordFunc` callback
// function that will ensure the the record being processed is a valid Who's On First GeoJSON Feature record.
func LoadRecordFunc(opts *LoadRecordFuncOptions) IndexerLoadRecordFunc {

	cb := func(ctx context.Context, path string, r io.ReadSeeker, args ...interface{}) (interface{}, error) {

		select {

		case <-ctx.Done():
			return nil, nil
		default:
			// pass
		}

		body, err := io.ReadAll(r)

		if err != nil {
			return nil, fmt.Errorf("Failed read %s, %w", path, err)
		}

		_, err = properties.Id(body)

		if err != nil {
			return nil, fmt.Errorf("Failed to derive wof:id for %s, %w", path, err)
		}

		_, err = geometry.Geometry(body)

		if err != nil {
			return nil, fmt.Errorf("Failed to derive geometry for %s, %w", path, err)
		}

		return body, nil
	}

	return cb
}

// IndexRelationsFunc returns a `go-whosonfirst-sql/indexer/IndexerPostIndexFunc` callback
// function used to index relations for a WOF record after that record has been successfully indexed.
func IndexRelationsFunc(r reader.Reader) IndexerPostIndexFunc {

	opts := &IndexRelationsFuncOptions{}
	opts.Reader = r

	return IndexRelationsFuncWithOptions(opts)
}

// IndexRelationsFuncWithOptions returns a `go-whosonfirst-sql/indexer/IndexerPostIndexFunc` callback
// function used to index relations for a WOF record after that record has been successfully indexed, but with custom
// `IndexRelationsFuncOptions` options defined in 'opts'.
func IndexRelationsFuncWithOptions(opts *IndexRelationsFuncOptions) IndexerPostIndexFunc {

	seen := new(sync.Map)

	cb := func(ctx context.Context, db *sql.DB, tables []database_sql.Table, record interface{}) error {

		geojson_t, err := wof_tables.NewGeoJSONTableWithDatabase(ctx, db)

		if err != nil {
			return fmt.Errorf("Failed to create new GeoJSON table, %w", err)
		}

		body := record.([]byte)

		relations := make(map[int64]bool)

		candidates := []string{
			"properties.wof:belongsto",
			"properties.wof:involves",
			"properties.wof:depicts",
		}

		for _, path := range candidates {

			// log.Println("RELATIONS", path)

			rsp := gjson.GetBytes(body, path)

			if !rsp.Exists() {
				// log.Println("MISSING", path)
				continue
			}

			for _, r := range rsp.Array() {

				id := r.Int()

				// skip -1, -4, etc.
				// (20201224/thisisaaronland)

				if id <= 0 {
					continue
				}

				relations[id] = true
			}
		}

		for id, _ := range relations {

			_, ok := seen.Load(id)

			if ok {
				continue
			}

			seen.Store(id, true)

			sql := fmt.Sprintf("SELECT COUNT(id) FROM %s WHERE id=?", geojson_t.Name())
			row := db.QueryRow(sql, id)

			var count int
			err = row.Scan(&count)

			if err != nil {
				return fmt.Errorf("Failed to count records for ID %d, %v", id, err)
			}

			if count != 0 {
				continue
			}

			rel_path, err := uri.Id2RelPath(id)

			if err != nil {
				return fmt.Errorf("Failed to determine relative path for %d, %v", id, err)
			}

			fh, err := opts.Reader.Read(ctx, rel_path)

			if err != nil {

				if opts.Strict {
					return fmt.Errorf("Failed to open %s, %v", rel_path, err)
				}

				slog.Debug("Failed to read '%s' because '%v'. Strict mode is disabled so skipping\n", rel_path, err)
				continue
			}

			defer fh.Close()

			ancestor, err := io.ReadAll(fh)

			if err != nil {
				return fmt.Errorf("Failed to read data for %s, %v", rel_path, err)
			}

			err = database_sql.IndexRecord(ctx, db, ancestor, tables...)

			if err != nil {
				return fmt.Errorf("Failed to index ancestor (%s), %v", rel_path, err)
			}
		}

		return nil
	}

	return cb
}
