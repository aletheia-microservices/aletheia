package ssagraph

import (
	"github.com/sirupsen/logrus"

	"analyzer/pkg/common"
)

func ComputeCallID(graph *SSAGraph, node *SSANode) string {
	return graph.GetServiceWithMethod() + "." + node.GetName()
}

type ServiceCall struct {
	id   string // format: <func_short_path>_<ssa_instr_name>
	t    string // format: <ssa_variable_name>
	node *SSANode
	args []*SSANode // does not include any receiver
	rets []*SSANode

	service       string
	method        string
	funcShortPath string
}

func NewServiceCall(id string, node *SSANode, args []*SSANode, rets []*SSANode, service string, method string, funcShortPath string) *ServiceCall {
	return &ServiceCall{
		id:            id,
		t:             node.GetValue().Name(),
		node:          node,
		args:          args,
		rets:          rets,
		service:       service,
		method:        method,
		funcShortPath: funcShortPath,
	}
}

func (call *ServiceCall) Copy() *ServiceCall {
	if call == nil {
		return nil
	}
	var copyArgs []*SSANode
	for _, arg := range call.args {
		copyArgs = append(copyArgs, arg.SimpleCopy())
	}
	var copyRets []*SSANode
	for _, ret := range call.rets {
		copyRets = append(copyRets, ret.SimpleCopy())
	}
	return &ServiceCall{
		id:            call.id,
		t:             call.t,
		node:          call.node.SimpleCopy(),
		args:          copyArgs,
		rets:          copyRets,
		service:       call.service,
		method:        call.method,
		funcShortPath: call.funcShortPath,
	}
}

func (call *ServiceCall) GetT() string {
	return call.t
}

func (call *ServiceCall) GetID() string {
	return call.id
}

func (call *ServiceCall) GetReturns() []*SSANode {
	return call.rets
}

func (call *ServiceCall) GetServiceWithMethod() string {
	return call.service + "." + call.method
}

func (call *ServiceCall) GetService() string {
	return call.service
}

func (call *ServiceCall) GetFuncShortPath() string {
	return call.funcShortPath
}

func (call *ServiceCall) GetMethod() string {
	return call.method
}

func (call *ServiceCall) GetNode() *SSANode {
	return call.node
}

func (call *ServiceCall) GetArguments() []*SSANode {
	return call.args
}

func (call *ServiceCall) String() string {
	return call.GetService() + "." + call.GetMethod()
}

type MethodCall struct {
	id   string // format: <func_short_path>_<ssa_instr_name>
	t    string // format: <ssa_variable_name>
	node *SSANode
	args []*SSANode // includes receiver if exists
	rets []*SSANode

	method        string
	funcShortPath string
}

func NewMethodCall(id string, node *SSANode, args []*SSANode, rets []*SSANode, method string, funcShortPath string) *MethodCall {
	return &MethodCall{
		id:            id,
		t:             node.GetValue().Name(),
		node:          node,
		args:          args,
		rets:          rets,
		method:        method,
		funcShortPath: funcShortPath,
	}
}

func (call *MethodCall) Copy() *MethodCall {
	if call == nil {
		return nil
	}
	var copyArgs []*SSANode
	for _, arg := range call.args {
		copyArgs = append(copyArgs, arg.SimpleCopy())
	}
	var copyRets []*SSANode
	for _, ret := range call.rets {
		copyRets = append(copyRets, ret.SimpleCopy())
	}
	return &MethodCall{
		id:            call.id,
		t:             call.t,
		node:          call.node.SimpleCopy(),
		args:          copyArgs,
		rets:          copyRets,
		method:        call.method,
		funcShortPath: call.funcShortPath,
	}
}

func (call *MethodCall) GetT() string {
	return call.t
}

func (call *MethodCall) GetID() string {
	return call.id
}

func (call *MethodCall) GetReturns() []*SSANode {
	return call.rets
}

func (call *MethodCall) GetReturnAt(idx int) *SSANode {
	if idx >= len(call.rets) {
		logrus.Fatalf("index (%d) out of range for call (%s) with returns lst: %v\n", idx, call.String(), call.rets)
	}
	return call.rets[idx]
}

func (call *MethodCall) TryGetReturnAt(idx int) *SSANode {
	if idx >= len(call.rets) {
		return nil
	}
	return call.rets[idx]
}

func (call *MethodCall) GetFuncShortPath() string {
	return call.funcShortPath
}

func (call *MethodCall) GetMethod() string {
	return call.method
}

func (call *MethodCall) GetNode() *SSANode {
	return call.node
}

func (call *MethodCall) GetArguments() []*SSANode {
	return call.args
}

func (call *MethodCall) GetArgumentAt(idx int) *SSANode {
	return call.args[idx]
}

func (call *MethodCall) String() string {
	return call.GetFuncShortPath()
}

type DatabaseCall struct {
	id     string // the ssa instr name for the db call on the callee side
	t      string // format: <ssa_variable_name>
	node   *SSANode
	args   []*SSANode // does not include any receiver
	opType common.DatabaseOperationType

	database string
	schema   string // can be e.g., collection, topic, table, etc.
	method   string
}

func NewDatabaseCall(id string, node *SSANode, args []*SSANode, database string, schema string, method string, opType common.DatabaseOperationType) *DatabaseCall {
	return &DatabaseCall{
		id:       id,
		t:        node.GetValue().Name(),
		node:     node,
		args:     args,
		database: database,
		schema:   schema,
		method:   method,
		opType:   opType,
	}
}

func (call *DatabaseCall) Copy() *DatabaseCall {
	if call == nil {
		return nil
	}
	var copyArgs []*SSANode
	for _, arg := range call.args {
		copyArgs = append(copyArgs, arg.SimpleCopy())
	}
	return &DatabaseCall{
		id:       call.id,
		t:        call.t,
		node:     call.node.SimpleCopy(),
		args:     copyArgs,
		opType:   call.opType,
		database: call.database,
		schema:   call.schema,
		method:   call.method,
	}
}

func (call *DatabaseCall) GetT() string {
	return call.t
}

func (call *DatabaseCall) GetID() string {
	return call.id
}

func (call *DatabaseCall) GetOpType() common.DatabaseOperationType {
	return call.opType
}

func (call *DatabaseCall) GetDatabasePath() string {
	return call.database + "." + call.schema
}

func (call *DatabaseCall) GetMethod() string {
	return call.method
}

func (call *DatabaseCall) GetDatabaseName() string {
	return call.database
}

func (call *DatabaseCall) GetSchemaName() string {
	return call.schema
}

func (call *DatabaseCall) GetNode() *SSANode {
	return call.node
}

func (call *DatabaseCall) GetArguments() []*SSANode {
	return call.args
}

func (call *DatabaseCall) String() string {
	return call.GetDatabasePath() + "." + call.GetMethod() + "(...)"
}
