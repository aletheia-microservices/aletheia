package abstractcallgraph

type AbstractArgument struct {
	ssaStr        string
	directTaints  map[string][]string
	indirectTaints map[string][]string
}

func NewAbstractArgument(directTaints map[string][]string, indirectTaints map[string][]string, ssaStr string) *AbstractArgument {
	arg := &AbstractArgument{
		ssaStr:        ssaStr,
		directTaints:  make(map[string][]string, len(directTaints)),
		indirectTaints: make(map[string][]string),
	}
	for k, v := range directTaints {
		copied := make([]string, len(v))
		copy(copied, v)
		arg.directTaints[k] = copied
	}
	
	for k, v := range indirectTaints {
		copied := make([]string, len(v))
		copy(copied, v)
		arg.indirectTaints[k] = copied
	}
	return arg
}

func (arg *AbstractArgument) SSAString() string {
	return arg.ssaStr
}

func (arg *AbstractArgument) GetDirectTaints() map[string][]string {
	return arg.directTaints
}

func (arg *AbstractArgument) GetIndirectTaints() map[string][]string {
	return arg.indirectTaints
}
