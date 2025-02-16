package abstractgraph

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"analyzer/pkg/app"
	"analyzer/pkg/logger"
	"analyzer/pkg/types"
)

func (graph *AbstractGraph) SetRequestsIDs() {
	for index, entryNode := range graph.Nodes {
		id := fmt.Sprintf("%d", index)
		entryNode.SetID(id)
		graph.setRequestsIDsHelper(entryNode, id)
	}
	graph.outputServicesExposedMethods()
}

func (graph *AbstractGraph) setRequestsIDsHelper(node AbstractNode, id string) {
	logger.Logger.Infof("[%s: %s -> %s] %s", id, node.GetCallerStr(), node.GetCallee(), node.GetMethodStr())
	for childIdx, childNode := range node.GetChildren() {
		id = fmt.Sprintf("%s.%d", id, childIdx)
		childNode.SetID(id)
		graph.setRequestsIDsHelper(childNode, id)
	}
}

func (graph *AbstractGraph) outputServicesExposedMethods() {
	/* for _, service := range graph.Services {
		for _, method := range service.ExposedMethods {
			logger.Logger.Infof("%s.%s", service.GetName(), method.SimpleString())
		}
	} */
	for _, entryNode := range graph.Nodes {
		method := entryNode.(*AbstractServiceCall).GetParsedCall().GetMethod().(*types.ParsedMethod)
		logger.Logger.Infof("%s.%s", entryNode.(*AbstractServiceCall).GetCallee(), method.SimpleString())

	}
}

type attachedDatabaseField struct {
	service string
	method  string
	argName string
	dbField string
}

func (graph *AbstractGraph) AttachDatabaseFieldsToEntryArgs(app *app.App, autofill bool) {
	graph.outputServicesExposedMethods()

	var input string
	var err error
	var attachedFields []attachedDatabaseField
	/* var targetRequestIDs []string */

	fmt.Printf("Please specify any database fields to be associated with services exposed methods.\nFormat (delimiter is ';'): <service>:<method_name>:<argument_name:database_field> (e.g., FrontendService:AddItem:itemID:CATALOGUE_DB.Sock.ID\n> ")

	// FORMAT: FrontendService:AddItem:itemID:CATALOGUE_DB.Sock.ID
	if autofill {
		/* if graph.AppName == "sockshop2" {
			input += "FrontendService:AddItem:itemID:CATALOGUE_DB.Sock.ID"
		} */
	} else {
		reader := bufio.NewReader(os.Stdin)
		input, err = reader.ReadString('\n')
		if err != nil {
			logger.Logger.Fatalf("error attached database fields to entry args %s", err.Error())
			return
		}
	}

	if input == "" || input == " " || input == "\n" {
		return
	}

	for _, target := range strings.Split(input, ";") {
		logger.Logger.Debugf("[INPUT] parsing target: %s", target)
		splits := strings.Split(target, ":")
		service, method, arg, dbField := splits[0], splits[1], splits[2], splits[3]
		attachedFields = append(attachedFields, attachedDatabaseField{service, method, arg, dbField})
	}

	for _, attachedField := range attachedFields {
		service := app.GetService(attachedField.service)
		method := service.GetExportedMethod(attachedField.method)

		dbStrings := strings.SplitN(attachedField.dbField, ".", 2)
		dbName, dbFieldName := strings.ToLower(dbStrings[0]), dbStrings[1]
		dbInstance := app.GetDatastoreInstance(dbName)
		dbField := dbInstance.GetDatastore().GetSchema().GetField(dbFieldName)

		var methodParamIdx int
		for idx, param := range method.GetParams() {
			if param.GetName() == attachedField.argName {
				methodParamIdx = idx
			}
		}

		for _, abstractServiceCall := range graph.getAbstractServiceCallsToMethod(method) {
			object := abstractServiceCall.GetParam(methodParamIdx)
			dataflow := object.GetVariableInfo().SetDirectDataflow(dbName, service.GetName(), object, dbField, true, -1)
			dataflow.EnablePermanent()
			app.AddDataflow(dataflow, nil) //FIXME: can we actually do this????

			logger.Logger.Warnf("ADDED DATAFLOW FOR OBJECT WITH DATAFLOWS: %v", object.GetVariableInfo().GetDataflows())
		}
	}
}

func (graph *AbstractGraph) getAbstractServiceCallsToMethod(method types.Method) []*AbstractServiceCall {
	var abstractServiceCalls []*AbstractServiceCall
	for _, entryNode := range graph.Nodes {
		if serviceCall, ok := entryNode.(*AbstractServiceCall); ok {
			if serviceCall.GetParsedCall().GetMethod() == method {
				abstractServiceCalls = append(abstractServiceCalls, serviceCall)
			}
		}
	}
	return abstractServiceCalls
}
