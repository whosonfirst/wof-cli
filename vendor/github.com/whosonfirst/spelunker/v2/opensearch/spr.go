package opensearch

import (
	"log/slog"

	"github.com/sfomuseum/go-edtf"
	"github.com/sfomuseum/go-edtf/common"
	"github.com/sfomuseum/go-edtf/parser"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-whosonfirst-flags"
	"github.com/whosonfirst/go-whosonfirst-flags/existential"
	wof_spr "github.com/whosonfirst/go-whosonfirst-spr/v2"
	"github.com/whosonfirst/go-whosonfirst-uri"
)

// SpelunkerRecordSPR implements the `whosonfirst/go-whosonfirst-spr/v2.StandardPlacesResult` interface for Who's On First records stored in an OpenSearch index.
type SpelunkerRecordSPR struct {
	wof_spr.StandardPlacesResult

	// The problem with this is that we can't use it with api/spr.go
	// endpoint if it is using the Spelunker interface GetSPRWithId
	// method.

	props []byte
}

// SpelunkerStandardPlacesResults implements the `whosonfirst/go-whosonfirst-spr/v2.StandardPlacesResults` interface for Who's On First records stored in an OpenSearch index.
type SpelunkerStandardPlacesResults struct {
	wof_spr.StandardPlacesResults
	results []wof_spr.StandardPlacesResult
}

// Results returns the list of `whosonfirst/go-whosonfirst-spr/v2.StandardPlacesResult` instances stored in 'r'.
func (r *SpelunkerStandardPlacesResults) Results() []wof_spr.StandardPlacesResult {
	return r.results
}

// NewSpelunkerStandardPlacesResults returns a new `whosonfirst/go-whosonfirst-spr/v2.StandardPlacesResults` instance for 'results'.
func NewSpelunkerStandardPlacesResults(results []wof_spr.StandardPlacesResult) wof_spr.StandardPlacesResults {

	r := &SpelunkerStandardPlacesResults{
		results: results,
	}

	return r
}

// NewSpelunkerRecordSPR returns a new `whosonfirst/go-whosonfirst-spr/v2.StandardPlacesResult` instance derived from 'props'.
func NewSpelunkerRecordSPR(props []byte) (wof_spr.StandardPlacesResult, error) {

	s := &SpelunkerRecordSPR{
		props: props,
	}

	return s, nil
}

// Return the unique ID of the place result.
func (s *SpelunkerRecordSPR) Id() string {
	return gjson.GetBytes(s.props, "wof:id").String()
}

// Return the unique parent ID of the place result.
func (s *SpelunkerRecordSPR) ParentId() string {
	return gjson.GetBytes(s.props, "wof:parent_id").String()
}

// Return the name of the place result.
func (s *SpelunkerRecordSPR) Name() string {
	return gjson.GetBytes(s.props, "wof:name").String()
}

// Return the Who's On First placetype of the place result.
func (s *SpelunkerRecordSPR) Placetype() string {
	return gjson.GetBytes(s.props, "wof:placetype").String()
}

// Return the two-letter country code of the place result.
func (s *SpelunkerRecordSPR) Country() string {
	return gjson.GetBytes(s.props, "wof:country").String()
}

// Return the (Git) repository name where the source record for the place result is stored.
func (s *SpelunkerRecordSPR) Repo() string {
	return gjson.GetBytes(s.props, "wof:repo").String()
}

// Return the relative path for the Who's On First record associated with the place result.
func (s *SpelunkerRecordSPR) Path() string {

	id := gjson.GetBytes(s.props, "wof:id").Int()
	path, _ := uri.Id2RelPath(id)
	return path
}

// Return the fully-qualified URI (URL) for the Who's On First record associated with the place result.
func (s *SpelunkerRecordSPR) URI() string {
	return s.Path()
}

// Return the EDTF inception date of the place result.
func (s *SpelunkerRecordSPR) Inception() *edtf.EDTFDate {
	return s.edtfDate("edtf:inception")
}

// Return the EDTF cessation date of the place result.
func (s *SpelunkerRecordSPR) Cessation() *edtf.EDTFDate {
	return s.edtfDate("edtf:cessation")
}

// Return the latitude for the principal centroid (typically "label") of the place result.
func (s *SpelunkerRecordSPR) Latitude() float64 {
	return gjson.GetBytes(s.props, "geom:latitude").Float()
}

// Return the longitude for the principal centroid (typically "label") of the place result.
func (s *SpelunkerRecordSPR) Longitude() float64 {
	return gjson.GetBytes(s.props, "geom:longitude").Float()
}

// Return the minimum latitude of the bounding box of the place result.
func (s *SpelunkerRecordSPR) MinLatitude() float64 {
	return gjson.GetBytes(s.props, "geom:bbox.1").Float()
}

// Return the minimum longitude of the bounding box of the place result.
func (s *SpelunkerRecordSPR) MinLongitude() float64 {
	return gjson.GetBytes(s.props, "geom:bbox.0").Float()
}

// Return the maximum latitude of the bounding box of the place result.
func (s *SpelunkerRecordSPR) MaxLatitude() float64 {
	return gjson.GetBytes(s.props, "geom:bbox.3").Float()
}

// Return the maximum longitude of the bounding box of the place result.
func (s *SpelunkerRecordSPR) MaxLongitude() float64 {
	return gjson.GetBytes(s.props, "geom:bbox.2").Float()
}

// Return the Who's On First "existential" flag denoting whether the place result is "current" or not.
func (s *SpelunkerRecordSPR) IsCurrent() flags.ExistentialFlag {
	fl_i := gjson.GetBytes(s.props, "mz:is_current").Int()
	return s.existentialFlag(fl_i)
}

// Return the Who's On First "existential" flag denoting whether the place result is "ceased" or not`.
func (s *SpelunkerRecordSPR) IsCeased() flags.ExistentialFlag {

	fl_i := int64(0)

	r := gjson.GetBytes(s.props, "edtf:cessation")

	if r.Exists() {

		switch r.String() {
		case edtf.UNKNOWN, edtf.UNKNOWN_2012:
			fl_i = -1
		default:
			fl_i = 1
		}
	}

	return s.existentialFlag(fl_i)
}

// Return the Who's On First "existential" flag denoting whether the place result is superseded or not.
func (s *SpelunkerRecordSPR) IsDeprecated() flags.ExistentialFlag {

	fl_i := int64(0)

	r := gjson.GetBytes(s.props, "edtf:deprecated")

	if r.Exists() && r.String() != "" {
		fl_i = 1
	}

	return s.existentialFlag(fl_i)
}

// Return the Who's On First "existential" flag denoting whether the place result has been superseded.
func (s *SpelunkerRecordSPR) IsSuperseded() flags.ExistentialFlag {

	fl_i := int64(0)

	if len(s.SupersededBy()) > 0 {
		fl_i = 1
	}

	return s.existentialFlag(fl_i)
}

// Return the Who's On First "existential" flag denoting whether the place result supersedes other records.
func (s *SpelunkerRecordSPR) IsSuperseding() flags.ExistentialFlag {

	fl_i := int64(0)

	if len(s.Supersedes()) > 0 {
		fl_i = 1
	}

	return s.existentialFlag(fl_i)
}

// Return the list of Who's On First IDs that supersede the place result.
func (s *SpelunkerRecordSPR) SupersededBy() []int64 {

	return s.gatherIds("wof:superseded_by")
}

// Return the list of Who's On First IDs that are superseded by the place result.
func (s *SpelunkerRecordSPR) Supersedes() []int64 {

	return s.gatherIds("wof:supersedes")
}

// Return the list of Who's On First IDs that are ancestors of the place result.
func (s *SpelunkerRecordSPR) BelongsTo() []int64 {

	return s.gatherIds("wof:belongsto")
}

// Return the Unix timestamp indicating when the place result was last modified.
func (s *SpelunkerRecordSPR) LastModified() int64 {

	return gjson.GetBytes(s.props, "wof:lastmodified").Int()
}

func (s *SpelunkerRecordSPR) edtfDate(path string) *edtf.EDTFDate {

	str_dt := gjson.GetBytes(s.props, path).String()
	dt, err := parser.ParseString(str_dt)

	if err != nil {
		slog.Error("Failed to parse date", "id", s.Id(), "path", path, "date", str_dt, "error", err)
		return s.unknownEDTF()
	}

	return dt
}

func (s *SpelunkerRecordSPR) unknownEDTF() *edtf.EDTFDate {

	sp := common.UnknownDateSpan()

	d := &edtf.EDTFDate{
		Start:   sp.Start,
		End:     sp.End,
		EDTF:    edtf.UNKNOWN,
		Level:   -1,
		Feature: "Unknown",
	}

	return d
}

func (s *SpelunkerRecordSPR) existentialFlag(fl_i int64) flags.ExistentialFlag {

	fl, err := existential.NewKnownUnknownFlag(fl_i)

	if err != nil {
		fl, _ = existential.NewNullFlag()
	}

	return fl
}

func (s *SpelunkerRecordSPR) gatherIds(path string) []int64 {

	ids := make([]int64, 0)

	r := gjson.GetBytes(s.props, "wof:superseded_by")

	if !r.Exists() {
		return ids
	}

	ra := r.Array()
	count := len(ra)

	if count == 0 {
		return ids
	}

	ids = make([]int64, count)

	for idx, a := range ra {
		ids[idx] = a.Int()
	}

	return ids
}
