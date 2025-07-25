package abstractcallgraph

type AbstractArgument struct {
	ssaStr string
	taints map[string][]string
}

func NewAbstractArgument(taints map[string][]string, ssaStr string) *AbstractArgument {
	arg := &AbstractArgument{
		ssaStr: ssaStr,
		taints: make(map[string][]string, len(taints)),
	}
	for k, v := range taints {
		copied := make([]string, len(v))
		copy(copied, v)
		arg.taints[k] = copied
	}
	return arg
}

func (arg *AbstractArgument) SSAString() string {
	return arg.ssaStr
}

func (arg *AbstractArgument) GetTaints() map[string][]string {
	return arg.taints
}
