package sql

import (
	"context"
	db_sql "database/sql"
	"fmt"
	"net/url"

	"github.com/whosonfirst/go-whosonfirst-spelunker"
)

type SQLSpelunker struct {
	spelunker.Spelunker
	engine string
	db     *db_sql.DB
}

func init() {
	ctx := context.Background()
	spelunker.RegisterSpelunker(ctx, "sql", NewSQLSpelunker)
}

func NewSQLSpelunker(ctx context.Context, uri string) (spelunker.Spelunker, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	engine := u.Host

	q := u.Query()

	dsn := q.Get("dsn")

	if dsn == "" {
		return nil, fmt.Errorf("Missing ?dsn= parameter")
	}

	db, err := db_sql.Open(engine, dsn)

	if err != nil {
		return nil, fmt.Errorf("Failed to open database connection, %w", err)
	}

	// db.SetMaxOpenConns(1)

	s := &SQLSpelunker{
		engine: engine,
		db:     db,
	}

	return s, nil
}

// concordances.go
// GetConcordances(context.Context) (*Faceting, error)
// HasConcordance(context.Context, pagination.Options, string, string, any, []Filter) (spr.StandardPlacesResults, pagination.Results, error)
// HasConcordanceFaceted(context.Context, string, string, any, []Filter, []*Facet) ([]*Faceting, error)

// descendants.go
// GetDescendants(context.Context, pagination.Options, int64, []Filter) (spr.StandardPlacesResults, pagination.Results, error)
// GetDescendantsFaceted(context.Context, int64, []Filter, []*Facet) ([]*Faceting, error)
// CountDescendants(context.Context, int64) (int64, error)

// id.go
// GetRecordForId(context.Context, int64, *uri.URIArgs) ([]byte, error)
// GetFeatureForId(context.Context, int64, *uri.URIArgs) ([]byte, error)
// GetSPRForId(context.Context, int64, *uri.URIArgs) (spr.StandardPlacesResult, error)

// nullisland.go
// VisitingNullIsland(context.Context, pagination.Options, []Filter) (spr.StandardPlacesResults, pagination.Results, error)
// VisitingNullIslandFaceted(context.Context, []Filter, []*Facet) ([]*Faceting, error)

// placetypes.go
// GetPlacetypes(context.Context) (*Faceting, error)
// HasPlacetype(context.Context, pagination.Options, *placetypes.WOFPlacetype, []Filter) (spr.StandardPlacesResults, pagination.Results, error)
// HasPlacetypeFaceted(context.Context, *placetypes.WOFPlacetype, []Filter, []*Facet) ([]*Faceting, error)

// recent.go
// GetRecent(context.Context, pagination.Options, time.Duration, []Filter) (spr.StandardPlacesResults, pagination.Results, error)
// GetRecentFaceted(context.Context, time.Duration, []Filter, []*Facet) ([]*Faceting, error)

// search.go
// Search(context.Context, pagination.Options, *SearchOptions, []Filter) (spr.StandardPlacesResults, pagination.Results, error)
// SearchFaceted(context.Context, *SearchOptions, []Filter, []*Facet) ([]*Faceting, error)
