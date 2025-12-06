package xcy

import (
	"fmt"
)

type Lineage struct {
	ID         int
	Operations []*Operation
	Previous   *Lineage
}

func (l *Lineage) String() string {
	opsStr := "{"
	for i, op := range l.Operations {
		opsStr += op.String()
		if i < len(l.Operations)-1 {
			opsStr += ", "
		}
	}
	opsStr += "}"
	return fmt.Sprintf("lineage #%d operations = %s", l.ID, opsStr)
}

func (l *Lineage) AddOperation(op *Operation) {
	l.Operations = append(l.Operations, op)
}

func (l *Lineage) IsRoot() bool {
	return l.Previous == nil
}

func (l *Lineage) GetXCYDependenciesByMostRecent() []*Operation {
	var dependencies []*Operation
	// reverse - start from most recent
	for i := len(l.Operations) - 1; i >= 0; i-- {
		if l.Operations[i].Write {
			dependencies = append(dependencies, l.Operations[i])
		}
	}
	if !l.IsRoot() {
		dependencies = append(dependencies, l.Previous.GetXCYDependenciesByMostRecent()...)
	}
	return dependencies
}

func (l *Lineage) GetDependencies() []*Operation {
	var dependencies []*Operation
	if !l.IsRoot() {
		dependencies = append(dependencies, l.Previous.GetDependencies()...)
	}
	for _, op := range l.Operations {
		if op.Write {
			dependencies = append(dependencies, op)
		}
	}
	return dependencies
}

func (l *Lineage) GetOperations() []*Operation {
	return l.Operations
}
