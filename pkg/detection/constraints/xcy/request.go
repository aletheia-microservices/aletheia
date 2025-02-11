package xcy

import (
	"encoding/json"
	"fmt"
	"slices"

	"github.com/golang-collections/collections/stack"

	"analyzer/pkg/abstractgraph"
	"analyzer/pkg/datastores"
	"analyzer/pkg/frameworks/blueprint"
	"analyzer/pkg/logger"
	"analyzer/pkg/types/objects"
)

type Request struct {
	EntryNode       abstractgraph.AbstractNode
	Inconsistencies []*XCYInconsistency
	Lineages        []*Lineage
	Operations      []*Operation
	LineagesStack   *stack.Stack
}

func (request *Request) GetReadOperations() []*Operation {
	var operations []*Operation
	for _, op := range request.Operations {
		if !op.Write {
			operations = append(operations, op)
		}
	}
	return operations
}

func (request *Request) GetWriteOperationsForDatastore(datastore *datastores.Datastore) []*Operation {
	var operations []*Operation
	for _, op := range request.Operations {
		if op.Datastore == datastore {
			operations = append(operations, op)
		}
	}
	return operations
}

func (request *Request) AddInconsistency(inconsistency *XCYInconsistency) {
	request.Inconsistencies = append(request.Inconsistencies, inconsistency)
}

func (request *Request) AddWrite(op *Operation) {
	op.ID = len(request.CurrentLineage().Operations)
	op.LineageID = request.CurrentLineage().ID
	request.Operations = append(request.Operations, op)
}

func (request *Request) AddRead(op *Operation) {
	op.ID = len(request.CurrentLineage().Operations)
	op.LineageID = request.CurrentLineage().ID
	request.Operations = append(request.Operations, op)
}

func (request *Request) PushLineage() {
	lineage := &Lineage{
		ID: len(request.Lineages),
	}
	if lineage.ID != 0 {
		lineage.Previous = request.CurrentLineage()
	}
	request.LineagesStack.Push(lineage)
	logger.Logger.Infof("[XCY LINEAGES - PUSH] %s", lineage.String())
	request.Lineages = append(request.Lineages, lineage)
}

func (request *Request) CurrentLineage() *Lineage {
	lineage := request.LineagesStack.Peek().(*Lineage)
	logger.Logger.Tracef("[XCY LINEAGES - PEEK] %s", lineage.String())
	return lineage
}

func (request *Request) PopLineage() *Lineage {
	lineage := request.LineagesStack.Pop().(*Lineage)
	logger.Logger.Infof("[XCY LINEAGES - POP] %s", lineage.String())
	return lineage
}

func (request *Request) GetXCYDependencies() []*Operation {
	var dependencies []*Operation
	for _, op := range request.Operations {
		if op.Write {
			dependencies = append(dependencies, op)
		}
	}
	return dependencies
}

func (request *Request) createOperation(key objects.Object, object objects.Object, call *abstractgraph.AbstractDatabaseCall, write bool) *Operation {
	op := &Operation{
		Key:       key,
		Object:    object,
		Service:   call.Service,
		Datastore: call.DbInstance.GetDatastore(),
		Method:    call.GetMethodStr(),
		Write:     write,
	}

	op.AppendVisibleDependencies(request.CurrentLineage().GetDependencies())
	return op
}

func (request *Request) SaveWriteOperation(call *abstractgraph.AbstractDatabaseCall, backend *blueprint.BackendMethod) *Operation {
	keyIndex := backend.GetWrittenKeyIndex()
	objIndex := backend.GetWrittenObjectIndex()
	key := call.GetParam(keyIndex)
	object := call.GetParam(objIndex)

	write := request.createOperation(key, object, call, true)

	request.AddWrite(write)
	request.CurrentLineage().AddOperation(write)

	return write
}

func (request *Request) SaveReadOperation(call *abstractgraph.AbstractDatabaseCall, backend *blueprint.BackendMethod) *Operation {
	keyIndex := backend.GetReadKeyIndex()
	objIndex := backend.GetReadObjectIndex()

	var object objects.Object
	key := call.GetParam(keyIndex)

	if objIndex >= 0 {
		object = call.GetParam(objIndex)
	} else {
		// gives index in inverse
		// e.g. -1 is 0 | -2 is 1 | -3 is 2 |
		objIndex = -objIndex - 1
		object = call.GetReturn(objIndex)
	}

	read := request.createOperation(key, object, call, false)
	request.AddRead(read)
	request.CurrentLineage().AddOperation(read)

	return read
}

type XCYInconsistency struct {
	Write             *Operation
	Read              *Operation
	MissingDependency bool
	Dataflows         []*objects.ObjectDataflow
}

func NewXCYInconsistency(write *Operation, read *Operation, dataflows []*objects.ObjectDataflow) *XCYInconsistency {
	return &XCYInconsistency{
		Write:     write,
		Read:      read,
		Dataflows: dataflows,
	}
}

func (inconsistency *XCYInconsistency) AppendDataflows(dataflows []*objects.ObjectDataflow) {
	inconsistency.Dataflows = append(inconsistency.Dataflows, dataflows...)
}

func (inconsistency *XCYInconsistency) MarshalJSON() ([]byte, error) {
	var dataflows []string
	for _, df := range inconsistency.Dataflows {
		prefix := ""
		if !df.Direct {
			prefix = "[indirect] "
		}
		repr := fmt.Sprintf("%s%s(%s, %s.%s)", prefix, df.GetOpString(), df.Service, df.Datastore, df.Field.GetName())
		if !slices.Contains(dataflows, repr) {
			dataflows = append(dataflows, repr)
		}
	}

	return json.Marshal(&struct {
		Write     *Operation `json:"write"`
		Read      *Operation `json:"read"`
		Dataflows []string   `json:"-"`
	}{
		Write:     inconsistency.Write,
		Read:      inconsistency.Read,
		Dataflows: dataflows,
	})
}

func (i *XCYInconsistency) String() string {
	return fmt.Sprintf("READ = %v / WRITE = %v", i.Read.String(), i.Write.String())
}
