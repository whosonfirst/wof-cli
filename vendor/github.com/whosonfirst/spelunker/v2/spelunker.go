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

// Spelunker is an interface for reading and querying Who's On First style data from an "index" (a database or queryable datafile).
type Spelunker interface {

	// Retrieve properties (or more specifically the "document") for a given ID.
	GetRecordForId(context.Context, int64, *uri.URIArgs) ([]byte, error)
	// Retrieve the `spr.StandardPlaceResult` instance for a given ID.
	GetSPRForId(context.Context, int64, *uri.URIArgs) (spr.StandardPlacesResult, error)
	// Retrieve the GeoJSON Feature record for a given ID.
	GetFeatureForId(context.Context, int64, *uri.URIArgs) ([]byte, error)

	// Retrieve all the Who's On First record that are a descendant of a specific Who's On First ID.
	GetDescendants(context.Context, pagination.Options, int64, []Filter) (spr.StandardPlacesResults, pagination.Results, error)
	// Retrieve faceted properties for records that are a descendant of a specific Who's On First ID.
	GetDescendantsFaceted(context.Context, int64, []Filter, []*Facet) ([]*Faceting, error)
	// Return the total number of Who's On First records that are a descendant of a specific Who's On First ID.
	CountDescendants(context.Context, int64) (int64, error)

	// Retrieve all the Who's On First records that match a search criteria.
	Search(context.Context, pagination.Options, *SearchOptions, []Filter) (spr.StandardPlacesResults, pagination.Results, error)
	// Retrieve faceted properties for records match a search criteria.
	SearchFaceted(context.Context, *SearchOptions, []Filter, []*Facet) ([]*Faceting, error)

	// Retrieve all the Who's On First records that have been modified with a window of time.
	GetRecent(context.Context, pagination.Options, time.Duration, []Filter) (spr.StandardPlacesResults, pagination.Results, error)
	// Retrieve faceted properties for records that have been modified with a window of time.
	GetRecentFaceted(context.Context, time.Duration, []Filter, []*Facet) ([]*Faceting, error)

	// Retrieve the list of unique placetypes in a Spleunker index.
	GetPlacetypes(context.Context) (*Faceting, error)
	// Retrieve the list of records with a given placetype.
	HasPlacetype(context.Context, pagination.Options, *placetypes.WOFPlacetype, []Filter) (spr.StandardPlacesResults, pagination.Results, error)
	// Retrieve faceted properties for records with a given placetype.
	HasPlacetypeFaceted(context.Context, *placetypes.WOFPlacetype, []Filter, []*Facet) ([]*Faceting, error)

	// Retrieve the list of alternate placetype ("wof:placetype_alt") in a SQLSpelunker database.
	GetAlternatePlacetypes(context.Context) (*Faceting, error)
	// Retrieve the list of Who's On First records with a given alternate placetype ("wof:placetype_alt") in a SQLSpelunker database.
	HasAlternatePlacetype(context.Context, pagination.Options, string, []Filter) (spr.StandardPlacesResults, pagination.Results, error)
	// Retrieve faceted properties for records with a given alternate placetype ("wof:placetype_alt") in a SQLSpelunker database.
	HasAlternatePlacetypeFaceted(context.Context, string, []Filter, []*Facet) ([]*Faceting, error)

	// Retrieve the list of unique concordances in a Spelunker index.
	GetConcordances(context.Context) (*Faceting, error)
	// Retrieve the list of records with a given concordance.
	HasConcordance(context.Context, pagination.Options, string, string, any, []Filter) (spr.StandardPlacesResults, pagination.Results, error)
	// Retrieve faceted properties for records with a given concordance.
	HasConcordanceFaceted(context.Context, string, string, any, []Filter, []*Facet) ([]*Faceting, error)

	// Retrieve the list of unique tags in a Spelunker index.
	GetTags(context.Context) (*Faceting, error)
	// Retrieve the list of records that have a given tag.
	HasTag(context.Context, pagination.Options, string, []Filter) (spr.StandardPlacesResults, pagination.Results, error)
	// Retrieve faceted properties for records that have a given tag.
	HasTagFaceted(context.Context, string, []Filter, []*Facet) ([]*Faceting, error)

	// Retrieve the list of records that are "visiting Null Island" (have a latitude, longitude value of "0.0, 0.0".
	VisitingNullIsland(context.Context, pagination.Options, []Filter) (spr.StandardPlacesResults, pagination.Results, error)
	// Retrieve faceted properties for records that are "visiting Null Island" (have a latitude, longitude value of "0.0, 0.0".
	VisitingNullIslandFaceted(context.Context, []Filter, []*Facet) ([]*Faceting, error)
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

	if i == nil {
		return nil, fmt.Errorf("Scheme not implemented")
	}

	init_func := i.(SpelunkerInitializationFunc)
	return init_func(ctx, uri)
}

// SpelunkerSchemes returns the list of schemes that have been registered.
func SpelunkerSchemes() []string {

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
