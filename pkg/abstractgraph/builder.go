package abstractgraph

import (
	"fmt"

	"analyzer/pkg/app"
	"analyzer/pkg/frameworks/blueprint"
	"analyzer/pkg/logger"
	"analyzer/pkg/lookup"
	"analyzer/pkg/service"
	"analyzer/pkg/types"
	"analyzer/pkg/types/objects"
	"analyzer/pkg/utils"
)

func Build(app *app.App, frontends []string) *AbstractGraph {
	graph := &AbstractGraph{
		AppName:  app.Name,
		Nodes:    make([]AbstractNode, 0),
		Services: app.Services,
		GIndex:   1,
	}
	for _, frontend := range frontends {
		service := app.Services[frontend]
		for _, method := range service.ExposedMethods {
			graph.initBuild(app, app.Services[frontend], method)
		}
	}

	for _, node := range graph.Nodes {
		graph.cleanVisited(node)
	}
	graph.GIndex = 1
	return graph
}

func (graph *AbstractGraph) cleanVisited(node AbstractNode) {
	node.SetVisited(false)
	for _, edge := range node.GetChildren() {
		graph.cleanVisited(edge)
	}
}

func (graph *AbstractGraph) initBuild(app *app.App, serviceNode *service.Service, targetMethod *types.ParsedMethod) {
	graph.addEntry(serviceNode, targetMethod)
	for _, abstractEntry := range graph.Nodes {
		if abstractServiceEntry, ok := abstractEntry.(*AbstractServiceCall); ok { // sanity check
			graph.recurseBuild(app, abstractServiceEntry)
		}
	}
}

func (graph *AbstractGraph) createDummyAbstractServiceCall(node *service.Service, method *types.ParsedMethod, parent *AbstractDatabaseCall, depth int) AbstractServiceCall {
	logger.Logger.Debugf("[GRAPH - DUMMY] create dummy abstract service call for node %s and method %s with params %v", node.Name, method.Name, method.GetParams())

	var callerStr = "Client"
	if parent != nil {
		callerStr = parent.DbInstance.String()
	}

	parsedCall := types.ParsedCall{
		Name:    method.Name,
		CallStr: node.Name,
		Method:  method,
		Params:  nil, // just a reminder that it is not set
		Returns: nil, // just a reminder that it is not set
	}

	call := AbstractServiceCall{
		ParsedCall: &types.ParsedServiceCall{
			ParsedCall: parsedCall.DeepCopy(),
		},
		Caller: callerStr,
		Callee: node.Name,
		Depth:  depth,
	}
	for _, p := range method.GetParams() {
		v := lookup.CreateObjectFromType(p.GetName(), p.GetType())
		call.ParsedCall.Params = append(call.ParsedCall.Params, v)
		logger.Logger.Debugf("[GRAPH - DUMMY] created variable %s (%s)", v.String(), utils.GetType(v))
	}
	return call
}

func (graph *AbstractGraph) addEntry(node *service.Service, method *types.ParsedMethod) {
	// add entry node to graph
	entryCall := graph.createDummyAbstractServiceCall(node, method, nil, 0)
	graph.Nodes = append(graph.Nodes, &entryCall)
}

func (graph *AbstractGraph) recurseBuild(app *app.App, abstractNode AbstractNode) {
	if abstractNode.IsVisited() {
		return
	}
	abstractNode.SetVisited(true)

	switch node := abstractNode.(type) {
	case *AbstractServiceCall:
		targetMethod := graph.Services[node.Callee].GetExportedMethod(node.ParsedCall.Name)
		if targetMethod == nil {
			logger.Logger.Fatalf("unexpected nil target method for call (%s) in node (%s)", node.ParsedCall.Name, node.Callee)
		}
		fmt.Println("------------------------------------------------------------------------------------------------------------------------------------------------")
		logger.Logger.Infof("[ABSTRACT SERVICE CALL (%d) calls] %s", len(targetMethod.Calls), node.String())
		for _, c := range targetMethod.Calls {
			logger.Logger.Debugf("\t\t- call: %v", c.String())
		}
		fmt.Println("------------------------------------------------------------------------------------------------------------------------------------------------")
		graph.appendAbstractEdges(node, targetMethod)
		fmt.Println()
	case *AbstractDatabaseCall:
		if node.ParsedCall.Method.IsQueueWrite() {
			fmt.Println("------------------------------------------------------------------------------------------------------------------------------------------------")
			logger.Logger.Infof("[ABSTRACT DATABASE CALL] %s", node.String())
			fmt.Println("------------------------------------------------------------------------------------------------------------------------------------------------")
			graph.createQueueHandlerNodes(app, node)
		}
		fmt.Println()
	case *AbstractQueueHandler:
		targetMethod := graph.Services[node.Callee].GetQueueHandlerMethod(node.ParsedCall.Name)
		if targetMethod == nil {
			logger.Logger.Fatalf("unexpected nil target method for call (%s) in node (%s)", node.ParsedCall.Name, node.Callee)
		}
		fmt.Println("------------------------------------------------------------------------------------------------------------------------------------------------")
		logger.Logger.Infof("[ABSTRACT QUEUE HANDLER (%d calls)] %s", len(targetMethod.Calls), node.String())
		for _, c := range targetMethod.Calls {
			logger.Logger.Debugf("\t\t- call: %v", c.String())
		}
		fmt.Println("------------------------------------------------------------------------------------------------------------------------------------------------")
		graph.appendAbstractEdges(node, targetMethod)
		fmt.Println()
	default:
		logger.Logger.Fatalf("Error recursing build for %s\nUnknown node type: %s", node, utils.GetType(node))
	}
	for _, edge := range abstractNode.GetChildren() {
		graph.recurseBuild(app, edge)
	}
}

func getPushMethodIfPublisherChild(parent AbstractNode, child *AbstractDatabaseCall) *blueprint.BackendMethod {
	if queueHandler, ok := parent.(*AbstractQueueHandler); ok && !queueHandler.IsEnabled() {
		if child.ParsedCall.Method.IsQueueRead() && queueHandler.Publisher.DbInstance == child.ParsedCall.DbInstance {
			pushMethod, ok := queueHandler.Publisher.ParsedCall.Method.(*blueprint.BackendMethod)
			if !ok {
				return nil
			}
			return pushMethod
		}
	}
	return nil
}

func (graph *AbstractGraph) assignRemainingObjectIDs(objects []objects.Object) {
	for _, dep := range objects {
		if dep.GetVariableInfo().IsUnassigned() {
			dep.GetVariableInfo().AssignID(graph.getAndIncGIndex())
			logger.Logger.Debugf("\t\t\t[GID DEP] assigned new gid (%d) to (%s)", dep.GetId(), dep.String())
		}
	}
}

// FIXME: we should actually create a new deep copy for the CFG
// to avoid changing it if we have another abstract node that will use the same parsed method
func (graph *AbstractGraph) referenceCallObjects(parsedCall types.Call, child AbstractNode) {
	entryBlock := parsedCall.GetMethod().(*types.ParsedMethod).GetParsedCfg().GetEntryParsedBlock()
	for i, parentObj := range parsedCall.GetParams() {
		// variable at block index 0 is the receiver so we need to skip that one
		if parsedCall.GetMethod().(*types.ParsedMethod).HasReceiver() {
			i = i + 1
		}
		childObj := entryBlock.GetVariableAt(i)
		childObj.AddReferenceWithID(parentObj, child.GetCallerStr())
		logger.Logger.Infof("\t\t[REF CALL OBJS - PARAMS] added reference (%d) to parent %s: (%s) -> (%s)", childObj.GetId(), child.GetCallerStr(), childObj.GetVariableInfo().String(), parentObj.GetVariableInfo().String())
		graph.assignRemainingObjectIDs(parentObj.GetNestedDependencies(false))
	}
	liveBlock := parsedCall.GetMethod().(*types.ParsedMethod).GetParsedCfg().GetLastLiveBlock() //FIXME: check comment in GetLastLiveBlock function
	for i, childObj := range liveBlock.GetResults() {
		if tupleChildObj, ok := childObj.(*objects.TupleObject); ok { //FIXME: sometimes this is giving a tuple (WITH NO OBJECT INFO!!!)!!
			for j, tupleElem := range tupleChildObj.Objects {
				parentObj := parsedCall.GetReturn(j)
				parentObj.AddReferenceWithID(tupleElem, child.GetCallee())
				logger.Logger.Infof("\t\t[REF CALL OBJS - RESULTS] added reference (%d) to child %s: (%v) -> (%v)", parentObj.GetId(), child.GetCallee(), parentObj.GetVariableInfo(), tupleElem.GetVariableInfo())
			}
		} else {
			parentObj := parsedCall.GetReturn(i)
			parentObj.AddReferenceWithID(childObj, child.GetCallee())
			//FIXME: for some reason the chid has no info!!
			if childObj.GetVariableInfo() == nil {
				logger.Logger.Fatal("????")
			}
			logger.Logger.Infof("\t\t[REF CALL OBJS - RESULTS] added reference (%d) to child %s: (%s) -> (%s)", parentObj.GetId(), child.GetCallee(), parentObj.GetVariableInfo().String(), childObj.GetVariableInfo().String())
		}
	}

	/* if parsedCall.GetName() == "StorePostNoSQL" {
		logger.Logger.Fatalf("HERE!")
	} */
}

func (graph *AbstractGraph) referenceQueuePopObjects(queueHandler *AbstractQueueHandler, queuePopCall *AbstractDatabaseCall) {
	logger.Logger.Infof("[GRAPH REF POP] referencing published vars to pop method call: %s", queuePopCall.GetName())
	queuePushMethod := getPushMethodIfPublisherChild(queueHandler, queuePopCall)
	if queuePushMethod == nil {
		logger.Logger.Warnf("[GRAPH] nil queue push method for queue handler (%s) and queue pop call (%s)", queueHandler.GetName(), queuePopCall.GetName())
		return
	}

	for pushParamIdx, popParamIdx := range queuePushMethod.MatchQueueIdentifiers() {
		pushParam := queueHandler.Publisher.GetParam(pushParamIdx)
		popParam := queuePopCall.GetParam(popParamIdx)

		// FIXME: do we really need this?
		popParam = objects.UnwrapPointerVariable(popParam)

		logger.Logger.Warnf("FIXMEEEEEEEEEE! IDK IF IT IS WORKING IN NOTIFY BECAUSE OF ASSIGNMENTS AND ASSERTS AFTERWARDS")
		for _, dep := range pushParam.GetNestedDependencies(true) {
			if dep.GetVariableInfo().IsUnassigned() {
				dep.GetVariableInfo().AssignID(graph.getAndIncGIndex())
				logger.Logger.Debugf("\t\t\t[QUEUE POP - GID DEP] assigned new gid (%d) to (%s)", dep.GetId(), dep.String())
			}
		}
		popParam.AddReferenceWithID(pushParam, queuePopCall.GetCallerStr())
		logger.Logger.Infof("\t\t[QUEUE POP - REF BLOCK VAR] added reference (%d) from creator (%s): (%s) -> (%s)", popParam.GetId(), queuePopCall.GetCallerStr(), popParam.GetType().GetName(), popParam.GetVariableInfo().GetName())
	}
	//logger.Logger.Fatalf("EXIT!")
}

func (graph *AbstractGraph) appendAbstractEdges(parent AbstractNode, targetMethod *types.ParsedMethod) {
	for _, call := range targetMethod.Calls {
		_, rootIsQueueHandler := parent.(*AbstractQueueHandler)
		switch parsedCall := call.(type) {
		case *types.ParsedDatabaseCall:
			child := &AbstractDatabaseCall{
				ParsedCall: parsedCall,
				Service:    parsedCall.CallerTypeName.GetName(),
				DbInstance: parsedCall.DbInstance,
				Subscriber: rootIsQueueHandler,
				Depth:      parent.GetNextDepth(),
				Caller:     parent.GetMethodStr(),
			}
			parent.AddChild(child)
			logger.Logger.Infof("[GRAPH] addding node for abstract database call with parent (%s): %s", child.GetCallerStr(), child.String())

			// if root parent is a queue handler (e.g. workerThread) and child is queue pop call
			if queueHandler, ok := parent.(*AbstractQueueHandler); ok && !queueHandler.IsEnabled() {
				graph.referenceQueuePopObjects(queueHandler, child)
				queueHandler.Enable()
			}
			logger.Logger.Infof("[GRAPH] added node for abstract database call with parent (%s): %s", child.GetCallerStr(), child.String())

		case *types.ParsedServiceCall:
			child := &AbstractServiceCall{
				ParsedCall: parsedCall.DeepCopy(),
				Caller:     parsedCall.CallerTypeName.GetName(),
				Callee:     parsedCall.CalleeTypeName.GetName(),
				Depth:      parent.GetNextDepth(),
			}
			parent.AddChild(child)
			logger.Logger.Infof("[GRAPH] adding node for abstract service call: %s", child.String())
			graph.referenceCallObjects(parsedCall, child)
			logger.Logger.Infof("[GRAPH] added node for abstract service call: %s", child.String())

			/* if parsedCall.GetName() == "StorePostNoSQL" {
				logger.Logger.Fatalf("HERE!")
			} */

		case *types.ParsedInternalCall:
			tempChild := &AbstractTempInternalCall{
				ParsedCall: parsedCall,
				Service:    parsedCall.ServiceTypeName.GetName(),
				Depth:      parent.GetNextDepth(),
			}

			tempMethod := graph.Services[tempChild.Service].GetInternalOrPackageMethod(tempChild.ParsedCall.Name)
			graph.appendAbstractEdges(parent, tempMethod)

			logger.Logger.Infof("[GRAPH] adding temporary node for abstract service call: %s", tempChild.String())
			graph.referenceCallObjects(parsedCall, tempChild)
			logger.Logger.Infof("[GRAPH] added temporary node for abstract service call: %s", tempChild.String())
		}
	}
}

func (graph *AbstractGraph) createQueueHandlerNodes(app *app.App, publisher *AbstractDatabaseCall) {
	instance := publisher.ParsedCall.GetTargetedDatabaseInstance()
	if instance == nil {
		logger.Logger.Fatalf("received unexpected nil database for queue publisher parsed call: %s", publisher.ParsedCall.String())
	}
	for _, node := range app.Services {
		for _, handlerMethod := range node.GetQueueHandlersForDatabaseInstance(instance) {
			logger.Logger.Debugf("[GRAPH - QUEUE] found worker %s on instance '%s'", handlerMethod.String(), instance.GetName())
			abstractHandler := &AbstractQueueHandler{
				AbstractServiceCall: graph.createDummyAbstractServiceCall(node, handlerMethod, publisher, publisher.GetNextDepth()),
				DbInstance:          instance,
				Publisher:           publisher,
			}
			publisher.Children = append(publisher.Children, abstractHandler)
		}
	}
}
