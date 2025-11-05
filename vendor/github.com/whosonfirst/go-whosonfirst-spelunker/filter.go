package spelunker

type Filter interface {
	Scheme() string
	Value() any
}
