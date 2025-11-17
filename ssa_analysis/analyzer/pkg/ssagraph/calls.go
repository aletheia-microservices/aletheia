package ssagraph

import "analyzer/pkg/common"

func ComputeCallID(graph *SSAGraph, node *SSANode) string {
	return graph.GetServiceWithMethod() + "." + node.GetName()
}

type ServiceCall struct {
	id   string // format: <func_short_path>_<ssa_instr_name>
	t    string // format: <ssa_variable_name>
	node *SSANode
	args []*SSANode
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

type DatabaseCall struct {
	id     string // the ssa instr name for the db call on the callee side
	t      string // format: <ssa_variable_name>
	node   *SSANode
	args   []*SSANode
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
