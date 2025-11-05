package spelunker

import (
	"fmt"
)

type Concordance struct {
	MachineTag
	namespace string
	predicate string
	value     any
}

func NewConcordanceFromString(str_concordance string) (*Concordance, error) {
	return nil, ErrNotImplemented
}

func NewConcordanceFromTriple(namespace string, predicate string, value any) *Concordance {

	c := &Concordance{
		namespace: namespace,
		predicate: predicate,
		value:     value,
	}

	return c
}

func (c *Concordance) Namespace() string {
	return c.namespace
}

func (c *Concordance) Predicate() string {
	return c.predicate
}

func (c *Concordance) Value() any {
	return c.value
}

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
