package ssagraph

func ComputeCallID(graph *SSAGraph, node *SSANode) string {
	return graph.GetService() + "." + graph.GetMethodName() + "." + node.GetName()
}

type ServiceCall struct {
	id   string // format: <func_short_path>_<ssa_instr_name>
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
		node:          node,
		args:          args,
		rets:          rets,
		service:       service,
		method:        method,
		funcShortPath: funcShortPath,
	}
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
	return call.GetService() + "." + call.GetMethod() + "(...)"
}

type DatabaseCall struct {
	id    string // the ssa instr name for the db call on the callee side
	node  *SSANode
	args  []*SSANode
	write bool

	database          string
	collectionOrTopic string
	method            string
}

func NewDatabaseCall(id string, node *SSANode, args []*SSANode, database string, collectionOrTopic string, method string, write bool) *DatabaseCall {
	return &DatabaseCall{
		id:                id,
		node:              node,
		args:              args,
		database:          database,
		collectionOrTopic: collectionOrTopic,
		method:            method,
		write:             write,
	}
}

func (call *DatabaseCall) GetID() string {
	return call.id
}

func (call *DatabaseCall) IsWrite() bool {
	return call.write
}

func (call *DatabaseCall) GetDatabasePath() string {
	return call.database + "." + call.collectionOrTopic
}

func (call *DatabaseCall) GetMethod() string {
	return call.method
}

func (call *DatabaseCall) GetDatabaseName() string {
	return call.database
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
