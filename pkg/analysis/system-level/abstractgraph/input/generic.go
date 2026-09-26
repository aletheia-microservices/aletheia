package abstractgraphinput

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v2"

	"analyzer/pkg/analysis/common"
)

// InputModel is the language-independent list of calls made by a graph node.
type InputModel struct {
	Calls []*Call `yaml:"calls"`
}

type CallType string

const (
	CallTypeRPC CallType = "RPC"
	CallTypeDB  CallType = "DB"
)

type Call struct {
	CallID       string        `yaml:"call_id"`
	CallTS       string        `yaml:"call_ts"`
	Arguments    []*Node       `yaml:"arguments,omitempty"`
	CallType     CallType      `yaml:"call_type"`
	ServiceCall  *ServiceCall  `yaml:"service_call,omitempty"`
	DatabaseCall *DatabaseCall `yaml:"database_call,omitempty"`
}

type ServiceCall struct {
	Returns       []*Node `yaml:"returns,omitempty"`
	Service       string  `yaml:"service"`
	Method        string  `yaml:"method"`
	FuncShortPath string  `yaml:"func_short_path"`
}

type DatabaseCall struct {
	OperationType string `yaml:"operation_type"`
	Database      string `yaml:"database"`
	Schema        string `yaml:"schema"`
	Method        string `yaml:"method"`
}

type Node struct {
	Name   string              `yaml:"name"`
	Taints map[string][]*Taint `yaml:"taints,omitempty"`
}

type TaintType string

const (
	TaintService  TaintType = "TAINT_SERVICE"
	TaintDatabase TaintType = "TAINT_DATABASE"
)

type Taint struct {
	TaintType     TaintType      `yaml:"taint_type"`
	CallID        string         `yaml:"call_id"`
	CallerT       string         `yaml:"caller_t"`
	Path          string         `yaml:"path"`
	DatabaseTaint *DatabaseTaint `yaml:"database_taint,omitempty"`
}

type DatabaseTaint struct {
	ReadKey   bool `yaml:"read_key"`
	ReadValue bool `yaml:"read_value"`
}

// LoadInputModel reads and validates a YAML call model.
func LoadInputModel(path string) (*InputModel, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var model InputModel
	if err := yaml.UnmarshalStrict(data, &model); err != nil {
		return nil, fmt.Errorf("decode input model: %w", err)
	}
	if _, err := model.Index(); err != nil {
		return nil, err
	}
	return &model, nil
}

func OperationType(value string) (common.DatabaseOperationType, error) {
	switch strings.ToLower(value) {
	case "undefined":
		return common.OP_UNDEFINED, nil
	case "write":
		return common.OP_WRITE, nil
	case "update":
		return common.OP_UPDATE, nil
	case "read":
		return common.OP_READ, nil
	case "read_many":
		return common.OP_READ_MANY, nil
	case "delete":
		return common.OP_DELETE, nil
	default:
		return common.OP_UNDEFINED, fmt.Errorf("unknown database operation %q", value)
	}
}

func (model *InputModel) Index() (map[string]*Call, error) {
	if model == nil {
		return nil, fmt.Errorf("input model is nil")
	}
	calls := make(map[string]*Call, len(model.Calls))
	for i, call := range model.Calls {
		if call == nil || call.CallID == "" {
			return nil, fmt.Errorf("call %d: call_id is required", i)
		}
		if calls[call.CallID] != nil {
			return nil, fmt.Errorf("duplicate call_id %q", call.CallID)
		}
		calls[call.CallID] = call
		switch call.CallType {
		case CallTypeRPC:
			if call.ServiceCall == nil || call.DatabaseCall != nil {
				return nil, fmt.Errorf("call %q: RPC requires only service_call", call.CallID)
			}
			if call.ServiceCall.Service == "" || call.ServiceCall.Method == "" {
				return nil, fmt.Errorf("call %q: service and method are required", call.CallID)
			}
		case CallTypeDB:
			if call.DatabaseCall == nil || call.ServiceCall != nil {
				return nil, fmt.Errorf("call %q: DB requires only database_call", call.CallID)
			}
			db := call.DatabaseCall
			if db.Database == "" || db.Schema == "" {
				return nil, fmt.Errorf("call %q: database and schema are required", call.CallID)
			}
			if _, err := OperationType(db.OperationType); err != nil {
				return nil, fmt.Errorf("call %q: %w", call.CallID, err)
			}
		default:
			return nil, fmt.Errorf("call %q: unknown call_type %q", call.CallID, call.CallType)
		}
	}
	for _, call := range model.Calls {
		nodes := append([]*Node(nil), call.Arguments...)
		if call.ServiceCall != nil {
			nodes = append(nodes, call.ServiceCall.Returns...)
		}
		for _, node := range nodes {
			if node == nil {
				return nil, fmt.Errorf("call %q: nil node", call.CallID)
			}
			for _, taints := range node.Taints {
				for _, taint := range taints {
					if taint == nil {
						return nil, fmt.Errorf("call %q: nil taint", call.CallID)
					}
					ref := calls[taint.CallID]
					if ref == nil {
						return nil, fmt.Errorf("call %q: unknown taint call_id %q", call.CallID, taint.CallID)
					}
					switch taint.TaintType {
					case TaintDatabase:
						if ref.CallType != CallTypeDB || taint.DatabaseTaint == nil {
							return nil, fmt.Errorf("call %q: invalid database taint", call.CallID)
						}
						if !strings.HasPrefix(taint.Path, ref.DatabaseCall.Database+".") {
							return nil, fmt.Errorf("call %q: taint path %q does not belong to database %q", call.CallID, taint.Path, ref.DatabaseCall.Database)
						}
					case TaintService:
						if ref.CallType != CallTypeRPC || taint.DatabaseTaint != nil {
							return nil, fmt.Errorf("call %q: invalid service taint", call.CallID)
						}
					default:
						return nil, fmt.Errorf("call %q: unknown taint_type %q", call.CallID, taint.TaintType)
					}
				}
			}
		}
	}
	return calls, nil
}
