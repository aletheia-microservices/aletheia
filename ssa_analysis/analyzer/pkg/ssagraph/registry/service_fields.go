package registry

import (
	"fmt"
	"go/types"

	"golang.org/x/tools/go/ssa"

	"analyzer/pkg/app"
	"analyzer/pkg/app/services"
	"analyzer/pkg/ssagraph"
)

func RegisterFields(app *app.App, graphs []*ssagraph.SSAGraph) {
	for _, graph := range graphs {
		service := app.GetServiceWithConstructorShortPathIfExists(graph.GetFunctionShortPath())
		if service == nil {
			continue
		}
		lastNode := graph.GetNodes()[0]
		fmt.Printf("lastNode instr: [%T] %s\n", lastNode.GetInstruction(), lastNode.GetInstruction())
		if _, ok := lastNode.GetInstruction().(*ssa.Return); ok {
			for _, edge := range graph.GetEdgesToNode(lastNode) {
				fromNode := edge.GetFromNode()
				fmt.Printf("fromNode instr: [%T] %s\n", fromNode.GetInstruction(), fromNode.GetInstruction())
				if iface, ok := fromNode.GetInstruction().(*ssa.MakeInterface); ok {
					fmt.Printf("iface: %s\n", iface)
					fmt.Printf("iface X: [%T] %s\n", iface.X, iface.X)
					if alloc, ok := iface.X.(*ssa.Alloc); ok {
						fmt.Printf("alloc: %s\n", alloc)
						fmt.Printf("alloc type: [%T] %s\n", alloc.Type(), alloc.Type())
						if typesPointer, ok := alloc.Type().(*types.Pointer); ok {
							if typesNamed, ok := typesPointer.Elem().(*types.Named); ok {
								returnedService := app.GetServiceWithImplPathIfExists(typesNamed.String())
								fmt.Printf("returned service: %s\n", returnedService)
								if returnedService == service {
									allocNode := graph.GetNodeByName(alloc.Name())
									registerFieldsFromAlloc(app, graph, service, allocNode)
								}

							}
						}
					}
				}
			}
		}
	}
}

// check postnotification_simple.NewUploadServiceImpl.dot
// for now this is a bit hardcoded
func registerFieldsFromAlloc(app *app.App, graph *ssagraph.SSAGraph, service *services.Service, allocNode *ssagraph.SSANode) {
	for _, edge := range graph.GetEdgesFromNode(allocNode) {
		if edge.GetType() == ssagraph.EDGE_FIELD {
			fieldNode := edge.GetToNode()
			if ssaFieldAddr, ok := fieldNode.GetValue().(*ssa.FieldAddr); ok {
				fieldIdx := ssaFieldAddr.Field
				for _, fieldEdge := range graph.GetEdgesFromNode(fieldNode) {
					if fieldEdge.GetType() == ssagraph.EDGE_STORE_ADDRESS {
						storeNode := fieldEdge.GetToNode()
						for _, storeEdge := range graph.GetEdgesToNode(storeNode) {
							if storeEdge.GetType() == ssagraph.EDGE_STORE_VALUE {
								valNode := storeEdge.GetFromNode()
								if _, ok := valNode.GetValue().(*ssa.Parameter); ok {
									paramIdx := graph.GetIndexOfParameter(valNode)
									fmt.Printf("[REGISTRY] [%s] val node (index=%d): %v\n", service.GetName(), paramIdx, valNode.String())
									wiringName := service.GetWiringNameAt(paramIdx)
									field := service.GetFieldAt(fieldIdx)
									field.SetWiringName(wiringName)
								}
							}
						}
					}
				}
			}
		}
	}
}
