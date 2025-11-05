package spelunker

type Facet struct {
	Property string `json:"property"`
}

func (f *Facet) String() string {
	return f.Property
}

func NewFacet(p string) *Facet {
	f := &Facet{
		Property: p,
	}

	return f
}

type FacetCount struct {
	Key   string `json:"key"`
	Count int64  `json:"count"`
}

type Faceting struct {
	Facet   *Facet        `json:"facet"`
	Results []*FacetCount `json:"results"`
}
