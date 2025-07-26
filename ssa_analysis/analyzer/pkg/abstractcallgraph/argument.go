package abstractcallgraph

type AbstractArgument struct {
	ssaStr         string
	directTaints   map[string][]*AbstractTaint
	indirectTaints map[string][]*AbstractTaint
}

func NewAbstractArgument(directTaints map[string][]*AbstractTaint, indirectTaints map[string][]*AbstractTaint, ssaStr string) *AbstractArgument {
	arg := &AbstractArgument{
		ssaStr:          ssaStr,
		directTaints:   directTaints,
		indirectTaints: indirectTaints,
	}
	return arg
}

func (arg *AbstractArgument) SSAString() string {
	return arg.ssaStr
}

func (arg *AbstractArgument) GetDirectTaints() map[string][]*AbstractTaint {
	return arg.directTaints
}

func (arg *AbstractArgument) GetIndirectTaints() map[string][]*AbstractTaint {
	return arg.indirectTaints
}
