package spelunker

import (
	"context"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/aaronland/go-pagination"
	"github.com/aaronland/go-roster"
	"github.com/whosonfirst/go-whosonfirst-placetypes"
	"github.com/whosonfirst/go-whosonfirst-spr/v2"
	"github.com/whosonfirst/go-whosonfirst-uri"
)

var spelunker_roster roster.Roster

// SpelunkerInitializationFunc is a function defined by individual spelunker package and used to create
// an instance of that spelunker
type SpelunkerInitializationFunc func(ctx context.Context, uri string) (Spelunker, error)

// Spelunker is an interface for writing data to multiple sources or targets.
type Spelunker interface {

	// Retrieve properties (or more specifically the "document") for...
	GetRecordForId(context.Context, int64, *uri.URIArgs) ([]byte, error)
	GetSPRForId(context.Context, int64, *uri.URIArgs) (spr.StandardPlacesResult, error)
	// Retrive GeoJSON Feature for...
	GetFeatureForId(context.Context, int64, *uri.URIArgs) ([]byte, error)

	// Retrieve all the Who's On First record that are a descendant of a specific Who's On First ID.
	GetDescendants(context.Context, pagination.Options, int64, []Filter) (spr.StandardPlacesResults, pagination.Results, error)
	GetDescendantsFaceted(context.Context, int64, []Filter, []*Facet) ([]*Faceting, error)
	// Return the total number of Who's On First records that are a descendant of a specific Who's On First ID.
	CountDescendants(context.Context, int64) (int64, error)

	// Retrieve all the Who's On First records that match a search criteria.
	Search(context.Context, pagination.Options, *SearchOptions, []Filter) (spr.StandardPlacesResults, pagination.Results, error)
	SearchFaceted(context.Context, *SearchOptions, []Filter, []*Facet) ([]*Faceting, error)

	// Retrieve all the Who's On First records that have been modified with a window of time.
	GetRecent(context.Context, pagination.Options, time.Duration, []Filter) (spr.StandardPlacesResults, pagination.Results, error)
	GetRecentFaceted(context.Context, time.Duration, []Filter, []*Facet) ([]*Faceting, error)

	GetPlacetypes(context.Context) (*Faceting, error)
	HasPlacetype(context.Context, pagination.Options, *placetypes.WOFPlacetype, []Filter) (spr.StandardPlacesResults, pagination.Results, error)
	HasPlacetypeFaceted(context.Context, *placetypes.WOFPlacetype, []Filter, []*Facet) ([]*Faceting, error)

	GetConcordances(context.Context) (*Faceting, error)
	HasConcordance(context.Context, pagination.Options, string, string, any, []Filter) (spr.StandardPlacesResults, pagination.Results, error)
	HasConcordanceFaceted(context.Context, string, string, any, []Filter, []*Facet) ([]*Faceting, error)

	GetTags(context.Context) (*Faceting, error)
	HasTag(context.Context, pagination.Options, string, []Filter) (spr.StandardPlacesResults, pagination.Results, error)
	HasTagFaceted(context.Context, string, []Filter, []*Facet) ([]*Faceting, error)

	VisitingNullIsland(context.Context, pagination.Options, []Filter) (spr.StandardPlacesResults, pagination.Results, error)
	VisitingNullIslandFaceted(context.Context, []Filter, []*Facet) ([]*Faceting, error)

	// TBD...
	// Unclear whether this should implement all of https://github.com/whosonfirst/go-whosonfirst-spatial/blob/main/spatial.go#L11
	// or https://github.com/whosonfirst/go-whosonfirst-spatial/blob/main/database/database.go#L16
	//
	// See also:
	// https://github.com/whosonfirst/go-whosonfirst-spatial-pip/blob/main/http/api/pointinpolygon.go
	// which in turns requires implementing https://github.com/whosonfirst/go-whosonfirst-spatial/blob/main/app/app.go#L21
	//
	// So it all starts to be a bit much...
	//
	// Maybe all we want are the structs and helper methods from this
	// https://github.com/whosonfirst/go-whosonfirst-spatial-pip/blob/main/pip.go
	//
	// But as twisty as all the spatial database stuff is when you start trying to make a simpler
	// version you just always end up with the same problems and questions...
	//
	// See also:
	// https://github.com/whosonfirst/go-whosonfirst-spatial-pmtiles/blob/main/database.go
	//
	// PointInPolygon(context.Context, orb.Point) (spr.StandardPlacesResults, error)

}

// RegisterSpelunker registers 'scheme' as a key pointing to 'init_func' in an internal lookup table
// used to create new `Spelunker` instances by the `NewSpelunker` method.
func RegisterSpelunker(ctx context.Context, scheme string, init_func SpelunkerInitializationFunc) error {

	err := ensureSpelunkerRoster()

	if err != nil {
		return err
	}

	return spelunker_roster.Register(ctx, scheme, init_func)
}

func ensureSpelunkerRoster() error {

	if spelunker_roster == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		spelunker_roster = r
	}

	return nil
}

// NewSpelunker returns a new `Spelunker` instance configured by 'uri'. The value of 'uri' is parsed
// as a `url.URL` and its scheme is used as the key for a corresponding `SpelunkerInitializationFunc`
// function used to instantiate the new `Spelunker`. It is assumed that the scheme (and initialization
// function) have been registered by the `RegisterSpelunker` method.
func NewSpelunker(ctx context.Context, uri string) (Spelunker, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := spelunker_roster.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	init_func := i.(SpelunkerInitializationFunc)
	return init_func(ctx, uri)
}

// Schemes returns the list of schemes that have been registered.
func Schemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensureSpelunkerRoster()

	if err != nil {
		return schemes
	}

	for _, dr := range spelunker_roster.Drivers(ctx) {
		scheme := fmt.Sprintf("%s://", strings.ToLower(dr))
		schemes = append(schemes, scheme)
	}

	sort.Strings(schemes)
	return schemes
}
