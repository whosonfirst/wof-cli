package spelunker

// MachineTag defines an interface for working with machine tags.
type MachineTag interface {
	// Return the namespace for a given machine tag.
	Namespace() string
	// Return the predicate for a given machine tag.
	Predicate() string
	// Return the value for a given machine tag.
	Value() any
	// Return the string representation for a machine tag (namespace:predicate=value).
	String() string
}
