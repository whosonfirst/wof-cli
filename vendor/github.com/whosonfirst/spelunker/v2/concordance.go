package spelunker

import (
	"fmt"
)

// Concordance implements the `MachineTag` interface for records from external sources that Who's On First records have concordances with.
type Concordance struct {
	MachineTag
	namespace string
	predicate string
	value     any
}

// NewConcordanceFromString derives a new `Concordance` instance from 'str_concordance' which is expected to take the form of:
//
//	{NAMESPACE}:{PREDICATE}={VALUE}
func NewConcordanceFromString(str_concordance string) (*Concordance, error) {
	return nil, ErrNotImplemented
}

// NewConcordanceFromString derives a new `Concordance` instance from 'namespace' and 'predicate' and 'value'.
func NewConcordanceFromTriple(namespace string, predicate string, value any) *Concordance {

	c := &Concordance{
		namespace: namespace,
		predicate: predicate,
		value:     value,
	}

	return c
}

// Namespace returns the namespace associated with 'c'.
func (c *Concordance) Namespace() string {
	return c.namespace
}

// Predicate returns the predicate associated with 'c'.
func (c *Concordance) Predicate() string {
	return c.predicate
}

// Concordance returns the value associated with 'c'.
func (c *Concordance) Value() any {
	return c.value
}

// String returns the string representation of 'c'.
func (c *Concordance) String() string {

	var str_concordance string

	switch {
	case c.namespace != "" && c.predicate != "" && c.value != "":
		str_concordance = fmt.Sprintf("%s:%s=%v", c.namespace, c.predicate, c.value)
	case c.namespace != "" && c.predicate != "":
		str_concordance = fmt.Sprintf("%s:%s=", c.namespace, c.predicate)
	case c.namespace != "" && c.value != "":
		str_concordance = fmt.Sprintf("%s:=%v", c.predicate, c.value)
	case c.predicate != "" && c.value != "":
		str_concordance = fmt.Sprintf(":%s=%v", c.predicate, c.value)
	case c.namespace != "":
		str_concordance = fmt.Sprintf("%s:", c.namespace)
	case c.predicate != "":
		str_concordance = fmt.Sprintf(":%s=", c.predicate)
	case c.value != "":
		str_concordance = fmt.Sprintf(":=%v", c.value)
	}

	return str_concordance
}
