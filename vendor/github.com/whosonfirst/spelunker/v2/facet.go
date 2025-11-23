package spelunker

// type Facet is a struct for representing faceted properties.
type Facet struct {
	// Property is the name (label) of the facet.
	Property string `json:"property"`
}

// String returns the string representation for 'f'
func (f *Facet) String() string {
	return f.Property
}

// NewFacet returns a `Facet` instance representing 'p'.
func NewFacet(p string) *Facet {
	f := &Facet{
		Property: p,
	}

	return f
}

// FacetCount is a struct for representing the results of a faceting operation.
type FacetCount struct {
	// Key is the value of the property associated with a faceting.
	Key string `json:"key"`
	// The number of records associated with 'Key'.
	Count int64 `json:"count"`
}

// Faceting is a struct representing a faceting operation.
type Faceting struct {
	// The `Facet` instance being faceted.
	Facet *Facet `json:"facet"`
	// Results is an array of `FacetCount` instances representing the values of a faceting operation.
	Results []*FacetCount `json:"results"`
}
