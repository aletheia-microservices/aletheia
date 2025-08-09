package ssagraph

type EdgeType int

const (
	EDGE_USAGE EdgeType = iota
	EDGE_STORE_ADDRESS
	EDGE_STORE_VALUE // usually pointed by element that is used as copy in store
	EDGE_CALL_ON
	EDGE_RETURN_ON
	EDGE_EXTRACT
	EDGE_PHI_ON
	EDGE_BINOP_X
	EDGE_BINOP_Y
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

func (edge *SSAEdge) GetType() EdgeType {
	return edge.edgeType
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
	case EDGE_CALL_ON:
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
	default:
		return "UNKNOWN"
	}
}
