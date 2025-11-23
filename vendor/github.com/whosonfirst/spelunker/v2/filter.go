package spelunker

// Filter is an interface for property-specific filtering criteria.
type Filter interface {
	// Scheme is the URI scheme used to create the `Filter` instance.
	Scheme() string
	// Value is the value of the filter being applied.
	Value() any
}
