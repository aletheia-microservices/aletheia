package common

type DatabaseOperationType int

const (
	OP_UNDEFINED DatabaseOperationType = iota
	OP_READ
	OP_READ_MANY
	OP_WRITE
	OP_DELETE
	OP_UPDATE
)

func OperationTypeToString(opType DatabaseOperationType) string {
	switch opType {
	case OP_UNDEFINED:
		return "undefined"
	case OP_READ, OP_READ_MANY:
		return "read"
	case OP_WRITE:
		return "write"
	case OP_UPDATE:
		return "update"
	case OP_DELETE:
		return "delete"
	}
	return ""
}
