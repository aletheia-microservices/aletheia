package abstracttypes

import (
	"encoding/json"

	"analyzer/pkg/datastores"
	"analyzer/pkg/types"
	"analyzer/pkg/types/objects"
)

type AbstractEdgeType int

const (
	RPC AbstractEdgeType = iota
	MESSAGE_QUEUE
	INTER_PROCESS
)

type AbstractEdge struct {
	Type       AbstractEdgeType
	Source     AbstractNode
	Target     AbstractNode
	Datastore  datastores.DatabaseInstance
	ParsedCall *types.ParsedCall
}

func (edge *AbstractEdge) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Type      string           `json:"type"`
		Params    []objects.Object `json:"params"`
		Returns   []objects.Object `json:"returns"`
		Source    AbstractNode     `json:"source"`
		Target    AbstractNode     `json:"target"`
		Datastore string           `json:"datastore"`
	}{
		Type:      edge.GetTypeName(),
		Params:    edge.ParsedCall.Params,
		Returns:   edge.ParsedCall.Returns,
		Source:    edge.Source,
		Target:    edge.Target,
		Datastore: edge.Datastore.GetName(),
	})
}

func (edge *AbstractEdge) IsRPCType() bool {
	return edge.Type == RPC
}

func (edge *AbstractEdge) IsMessageQueueType() bool {
	return edge.Type == MESSAGE_QUEUE
}

func (edge *AbstractEdge) IsInterProcessType() bool {
	return edge.Type == INTER_PROCESS
}

func (edge *AbstractEdge) GetType() AbstractEdgeType {
	return edge.Type
}

func (edge *AbstractEdge) GetTypeName() string {
	var detectionMap = map[AbstractEdgeType]string{
		RPC:           "rpc",
		MESSAGE_QUEUE: "message_queue",
		INTER_PROCESS: "inter_process",
	}
	return detectionMap[edge.Type]
}

func (edge *AbstractEdge) GetMethod() types.Method {
	return edge.ParsedCall.GetMethod()
}

func (edge *AbstractEdge) GetMethodString() string {
	return edge.ParsedCall.GetMethod().String()
}

func (edge *AbstractEdge) GetParams() []objects.Object {
	return edge.ParsedCall.GetParams()
}

func (edge *AbstractEdge) GetReturns() []objects.Object {
	return edge.ParsedCall.GetReturns()
}

func (edge *AbstractEdge) GetSourceNode() AbstractNode {
	return edge.Source
}

func (edge *AbstractEdge) GetTargetNode() AbstractNode {
	return edge.Target
}

func (edge *AbstractEdge) GetDatastoreInstance() datastores.DatabaseInstance {
	return edge.Datastore
}

func (edge *AbstractEdge) GetParsedCall() *types.ParsedCall {
	return edge.ParsedCall
}
