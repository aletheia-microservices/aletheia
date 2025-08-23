package ssagraph

type EdgeType int

const (
	EDGE_USAGE EdgeType = iota
	EDGE_STORE_ADDRESS
	EDGE_STORE_VALUE // usually pointed by element that is used as copy in store
	EDGE_ARG_ON_CALL
	EDGE_RECEIVER_ON_CALL
	EDGE_RETURN_ON
	EDGE_EXTRACT
	EDGE_PHI_ON
	EDGE_BINOP_X
	EDGE_BINOP_Y
	EDGE_MAP_TARGET
	EDGE_MAP_KEY
	EDGE_MAP_VALUE
	EDGE_LOOKUP_TARGET
	EDGE_LOOKUP_INDEX
	EDGE_LOAD
	EDGE_FIELD
	EDGE_INDEX
	EDGE_PARAMETER
	EDGE_POINTS_TO
)

type SSAEdge struct {
	edgeType EdgeType
	from     *SSANode
	to       *SSANode

	index int
	param string

	path string //pointer only
}

func (edge *SSAEdge) GetParam() string {
	return edge.param
}

func (edge *SSAEdge) GetIndex() int {
	return edge.index
}

func (edge *SSAEdge) GetType() EdgeType {
	return edge.edgeType
}

func (edge *SSAEdge) IsType(t EdgeType) bool {
	return edge.edgeType == t
}

func (edge *SSAEdge) HasFromNode(node *SSANode) bool {
	return edge.from == node
}

func (edge *SSAEdge) GetFromNode() *SSANode {
	return edge.from
}

func (edge *SSAEdge) GetToNode() *SSANode {
	return edge.to
}

func (edge *SSAEdge) GetPath() string {
	return edge.path
}

func (edge *SSAEdge) SetPath(path string) {
	edge.path = path
}

func (edge *SSAEdge) GetTypeString() string {
	switch edge.GetType() {
	case EDGE_USAGE:
		return "USAGE"
	case EDGE_STORE_ADDRESS:
		return "STORE_ADDRESS"
	case EDGE_STORE_VALUE:
		return "STORE_VALUE"
	case EDGE_ARG_ON_CALL:
		return "CALL_ON"
	case EDGE_RECEIVER_ON_CALL:
		return "CALL_ON"
	case EDGE_RETURN_ON:
		return "RETURN_ON"
	case EDGE_EXTRACT:
		return "EXTRACT"
	case EDGE_PHI_ON:
		return "PHI_ON"
	case EDGE_BINOP_X:
		return "BINOP_X"
	case EDGE_BINOP_Y:
		return "BINOP_Y"
	case EDGE_LOAD:
		return "LOAD"
	case EDGE_FIELD:
		return "FIELD"
	case EDGE_INDEX:
		return "INDEX"
	case EDGE_PARAMETER:
		return "PARAMETER"
	case EDGE_POINTS_TO:
		return "POINTS_TO"
	case EDGE_MAP_TARGET:
		return "EDGE_MAP_TARGET"
	case EDGE_MAP_KEY:
		return "EDGE_MAP_KEY"
	case EDGE_MAP_VALUE:
		return "EDGE_MAP_VALUE"
	case EDGE_LOOKUP_TARGET:
		return "EDGE_LOOKUP_TARGET"
	case EDGE_LOOKUP_INDEX:
		return "EDGE_LOOKUP_INDEX"
	default:
		return "UNKNOWN"
	}
}
