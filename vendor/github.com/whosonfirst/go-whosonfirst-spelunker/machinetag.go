package spelunker

type MachineTag interface {
	Namespace() string
	Predicate() string
	Value() any
	String() string
}
