package ssagraph

import (
	"github.com/sirupsen/logrus"

	"analyzer/pkg/common"
)

func ComputeCallID(graph *SSAGraph, node *SSANode) string {
	return graph.GetServiceWithMethod() + "." + node.GetName()
}

// Call is implemented by every call tracked in the SSA graph
type Call interface {
	GetID() string
	GetT() string
	GetMethod() string
	GetNode() *SSANode
	GetArguments() []*SSANode
	String() string
}

// baseCall holds the fields common to all calls
type baseCall struct {
	// --- input for abstract call graph ---
	callID string     // unique call identifier; format: [<func_short_path>_]<ssa_instr_unique_name> (func_short_path does not exist for database calls)
	callTS string     // timestamp to order calls within function; format: <ssa_variable_name>
	args   []*SSANode // arguments passed to the call
	method string

	// --- extra info for SSA graph ---
	node *SSANode
}

func newBaseCall(id string, ts string, node *SSANode, args []*SSANode, method string) baseCall {
	return baseCall{
		callID: id,
		callTS: ts,
		node:   node,
		args:   args,
		method: method,
	}
}

func (call *baseCall) GetT() string {
	return call.callTS
}

func (call *baseCall) GetID() string {
	return call.callID
}

func (call *baseCall) GetMethod() string {
	return call.method
}

func (call *baseCall) GetNode() *SSANode {
	return call.node
}

func (call *baseCall) GetArguments() []*SSANode {
	return call.args
}

type ServiceCall struct {
	// --- input for abstract call graph ---
	baseCall      // args do not include any receiver
	rets          []*SSANode
	funcShortPath string //TODO: rename to funcID
	service       string
}

func NewServiceCall(id string, node *SSANode, args []*SSANode, rets []*SSANode, service string, method string, funcShortPath string) *ServiceCall {
	return &ServiceCall{
		baseCall:      newBaseCall(id, node.GetValue().Name(), node, args, method),
		rets:          rets,
		service:       service,
		funcShortPath: funcShortPath,
	}
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

func (call *ServiceCall) String() string {
	return call.GetService() + "." + call.GetMethod()
}

type MethodCall struct {
	// --- input for abstract call graph ---
	baseCall                 // args include receiver if exists
	binds         []*SSANode // go routine
	rets          []*SSANode
	funcShortPath string
	goroutine     bool
}

func NewMethodCall(id string, node *SSANode, args []*SSANode, rets []*SSANode, method string, funcShortPath string) *MethodCall {
	return &MethodCall{
		baseCall:      newBaseCall(id, node.GetValue().Name(), node, args, method),
		rets:          rets,
		funcShortPath: funcShortPath,
	}
}

func NewMethodCallGoRoutine(id string, instrID string, node *SSANode, binds []*SSANode, args []*SSANode, rets []*SSANode, method string, funcShortPath string) *MethodCall {
	return &MethodCall{
		baseCall:      newBaseCall(id, instrID, node, args, method),
		binds:         binds,
		rets:          rets,
		funcShortPath: funcShortPath,
		goroutine:     true,
	}
}

func (call *MethodCall) GetReturns() []*SSANode {
	return call.rets
}

func (call *MethodCall) GetBindAt(idx int) *SSANode {
	return call.binds[idx]
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

func (call *MethodCall) GetArgumentAt(idx int) *SSANode {
	return call.args[idx]
}

func (call *MethodCall) String() string {
	return call.GetFuncShortPath()
}

type DatabaseCall struct {
	// --- input for abstract call graph ---
	baseCall // args do not include any receiver
	opType   common.DatabaseOperationType
	database string
	schema   string // can be e.g., collection, topic, table, etc.
}

func NewDatabaseCall(id string, node *SSANode, args []*SSANode, database string, schema string, method string, opType common.DatabaseOperationType) *DatabaseCall {
	return &DatabaseCall{
		baseCall: newBaseCall(id, node.GetValue().Name(), node, args, method),
		database: database,
		schema:   schema,
		opType:   opType,
	}
}

func (call *DatabaseCall) GetOpType() common.DatabaseOperationType {
	return call.opType
}

func (call *DatabaseCall) GetDatabasePath() string {
	return call.database + "." + call.schema
}

func (call *DatabaseCall) GetDatabaseName() string {
	return call.database
}

func (call *DatabaseCall) GetSchemaName() string {
	return call.schema
}

func (call *DatabaseCall) String() string {
	return call.GetDatabasePath() + "." + call.GetMethod() + "(...)"
}
