package abstractgraphinput

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/aletheia-microservices/aletheia/internal/analysis/common"
)

// InputModel is the language-independent equivalent of the Blueprint wiring and SSA graphs: the databases
// of the app, and its functions with the entrypoints from which they are parsed into the abstract graph
type InputModel struct {
	App         string      `yaml:"app,omitempty"`
	Databases   []*Database `yaml:"databases,omitempty"`
	Entrypoints []string    `yaml:"entrypoints"`
	Functions   []*Function `yaml:"functions"`
}

// DatabaseTypes are the valid types of databases (same as the Blueprint backends)
var DatabaseTypes = []string{"NoSQLDatabase", "RelationalDB", "Cache", "Queue"}

type Database struct {
	Name    string    `yaml:"name"`
	Type    string    `yaml:"type"`
	Schemas []*Schema `yaml:"schemas,omitempty"`
}

// Schema declares the keys of a database schema (e.g., collection, table or topic) that cannot be
// inferred from the calls (the equivalent of bson _id tags and SQL files), with fields named relative to it
type Schema struct {
	Name       string   `yaml:"name"`
	PrimaryKey []string `yaml:"primary_key,omitempty"` // one (possibly composite) key
	Unique     []string `yaml:"unique,omitempty"`      // each field is unique by itself
}

// Function is a service method identified by its func_short_path, where the taints of its params,
// returns and calls can only reference its own calls
type Function struct {
	FuncShortPath string    `yaml:"func_short_path"`
	Service       string    `yaml:"service"`
	Method        string    `yaml:"method"`
	Params        []*Node   `yaml:"params,omitempty"`  // excluding receiver and context
	Returns       [][]*Node `yaml:"returns,omitempty"` // one list of objects for each return
	Calls         []*Call   `yaml:"calls,omitempty"`
}

// Index is a validated input model, with functions by func_short_path and their calls by call_id
type Index struct {
	Functions map[string]*Function
	Calls     map[string]map[string]*Call
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
	CallerT       string         `yaml:"caller_t,omitempty"`
	Path          string         `yaml:"path"`
	DatabaseTaint *DatabaseTaint `yaml:"database_taint,omitempty"` // only for database taints (default: not read)
}

type DatabaseTaint struct {
	ReadKey   bool `yaml:"read_key"`
	ReadValue bool `yaml:"read_value"`
}

// LoadInputModel reads and validates the YAML input model at path, which is either a file or a folder
// (usually one per app) whose files (*.yaml or *.yml) are combined into one model (see MergeInputModels)
func LoadInputModel(path string) (*InputModel, error) {
	paths := []string{path}
	if info, err := os.Stat(path); err != nil {
		return nil, err
	} else if info.IsDir() {
		yamlPaths, _ := filepath.Glob(filepath.Join(path, "*.yaml"))
		ymlPaths, _ := filepath.Glob(filepath.Join(path, "*.yml"))
		paths = slices.Sorted(slices.Values(append(yamlPaths, ymlPaths...)))
		if len(paths) == 0 {
			return nil, fmt.Errorf("no input models (*.yaml or *.yml) in folder %s", path)
		}
	}
	var models []*InputModel
	for _, path := range paths {
		model, err := decodeInputModel(path)
		if err != nil {
			return nil, err
		}
		models = append(models, model)
	}
	// only the combined model is validated, since a model can reference functions of another one
	model, err := MergeInputModels(models...)
	if err != nil {
		return nil, err
	}
	// indexing validates the model (the index itself is built again by the parser)
	if _, err := model.Index(); err != nil {
		return nil, err
	}
	return model, nil
}

func decodeInputModel(path string) (*InputModel, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	// strict decoding fails on unknown or duplicate fields
	var model InputModel
	if err := yaml.UnmarshalStrict(data, &model); err != nil {
		return nil, fmt.Errorf("decode input model %s: %w", path, err)
	}
	return &model, nil
}

// MergeInputModels combines models (e.g., one for the frontend services and another for the backend services)
// into one model with all their databases, entrypoints and functions, where:
//   - the app must be the same in all models that name it
//   - a database can be declared in multiple models with the same type, and its schemas are combined
func MergeInputModels(models ...*InputModel) (*InputModel, error) {
	merged := &InputModel{}
	databases := make(map[string]*Database)
	for _, model := range models {
		if model.App != "" {
			if merged.App != "" && merged.App != model.App {
				return nil, fmt.Errorf("input models belong to different apps (%q and %q)", merged.App, model.App)
			}
			merged.App = model.App
		}
		for _, db := range model.Databases {
			// invalid databases are kept as they are, so that the validation reports them
			if db == nil || db.Name == "" {
				merged.Databases = append(merged.Databases, db)
				continue
			}
			existing := databases[db.Name]
			if existing == nil {
				existing = &Database{Name: db.Name, Type: db.Type, Schemas: slices.Clone(db.Schemas)}
				databases[db.Name] = existing
				merged.Databases = append(merged.Databases, existing)
				continue
			}
			if existing.Type != db.Type {
				return nil, fmt.Errorf("database %q is declared with different types (%q and %q)", db.Name, existing.Type, db.Type)
			}
			existing.Schemas = append(existing.Schemas, db.Schemas...)
		}
		merged.Entrypoints = append(merged.Entrypoints, model.Entrypoints...)
		merged.Functions = append(merged.Functions, model.Functions...)
	}
	return merged, nil
}

func OperationType(value string) (common.DatabaseOperationType, error) {
	// operation types are case-insensitive
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

// tsPattern is the format of call timestamps, e.g., t4 or t4.t7 (see utils.LessT)
var tsPattern = regexp.MustCompile(`^t[0-9]+(\.t[0-9]+)*$`)

func (model *InputModel) Index() (*Index, error) {
	if model == nil {
		return nil, fmt.Errorf("input model is nil")
	}
	if err := model.validateDatabases(); err != nil {
		return nil, err
	}
	// validate each function and index it (and its calls) by func_short_path
	index := &Index{
		Functions: make(map[string]*Function, len(model.Functions)),
		Calls:     make(map[string]map[string]*Call, len(model.Functions)),
	}
	for i, fn := range model.Functions {
		if fn == nil || fn.FuncShortPath == "" {
			return nil, fmt.Errorf("function %d: func_short_path is required", i)
		}
		if index.Functions[fn.FuncShortPath] != nil {
			return nil, fmt.Errorf("duplicate func_short_path %q", fn.FuncShortPath)
		}
		if fn.Service == "" || fn.Method == "" {
			return nil, fmt.Errorf("function %q: service and method are required", fn.FuncShortPath)
		}
		calls, err := fn.index()
		if err != nil {
			return nil, fmt.Errorf("function %q: %w", fn.FuncShortPath, err)
		}
		index.Functions[fn.FuncShortPath] = fn
		index.Calls[fn.FuncShortPath] = calls
	}
	// calls are matched by call_id across functions of the same request
	callers := make(map[string]string)
	for _, fn := range model.Functions {
		for _, call := range fn.Calls {
			if caller, ok := callers[call.CallID]; ok {
				return nil, fmt.Errorf("function %q: call_id %q is already used in function %q", fn.FuncShortPath, call.CallID, caller)
			}
			callers[call.CallID] = fn.FuncShortPath
		}
	}
	// each service call must reference an existing function with the same service and method
	for _, fn := range model.Functions {
		for _, call := range fn.Calls {
			svc := call.ServiceCall
			if svc == nil {
				continue
			}
			callee := index.Functions[svc.FuncShortPath]
			if callee == nil {
				return nil, fmt.Errorf("function %q: call %q: unknown func_short_path %q", fn.FuncShortPath, call.CallID, svc.FuncShortPath)
			}
			if callee.Service != svc.Service || callee.Method != svc.Method {
				return nil, fmt.Errorf("function %q: call %q: %s.%s does not match function %q", fn.FuncShortPath, call.CallID, svc.Service, svc.Method, svc.FuncShortPath)
			}
		}
	}
	// there must be at least one entrypoint, and all of them must be existing functions
	if len(model.Entrypoints) == 0 {
		return nil, fmt.Errorf("at least one entrypoint is required")
	}
	for i, entrypoint := range model.Entrypoints {
		if index.Functions[entrypoint] == nil {
			return nil, fmt.Errorf("unknown entrypoint %q", entrypoint)
		}
		if slices.Contains(model.Entrypoints[:i], entrypoint) {
			return nil, fmt.Errorf("duplicate entrypoint %q", entrypoint)
		}
	}
	return index, nil
}

func (model *InputModel) validateDatabases() error {
	// databases must have a unique name and a valid type
	databases := make(map[string]bool, len(model.Databases))
	for i, db := range model.Databases {
		if db == nil || db.Name == "" {
			return fmt.Errorf("database %d: name is required", i)
		}
		if databases[db.Name] {
			return fmt.Errorf("duplicate database %q", db.Name)
		}
		databases[db.Name] = true
		if !slices.Contains(DatabaseTypes, db.Type) {
			return fmt.Errorf("database %q: type must be one of %v", db.Name, DatabaseTypes)
		}
		// schemas must have a unique name in the database, and no empty key fields
		schemas := make(map[string]bool, len(db.Schemas))
		for j, schema := range db.Schemas {
			if schema == nil || schema.Name == "" {
				return fmt.Errorf("database %q: schema %d: name is required", db.Name, j)
			}
			if schemas[schema.Name] {
				return fmt.Errorf("database %q: duplicate schema %q", db.Name, schema.Name)
			}
			schemas[schema.Name] = true
			if slices.Contains(schema.PrimaryKey, "") || slices.Contains(schema.Unique, "") {
				return fmt.Errorf("database %q: schema %q: empty key field", db.Name, schema.Name)
			}
		}
	}
	return nil
}

// index validates the calls of fn and the taints of its objects, and returns its calls by call_id
func (fn *Function) index() (map[string]*Call, error) {
	// validate each call and index it by call_id
	calls := make(map[string]*Call, len(fn.Calls))
	for i, call := range fn.Calls {
		if call == nil || call.CallID == "" {
			return nil, fmt.Errorf("call %d: call_id is required", i)
		}
		if calls[call.CallID] != nil {
			return nil, fmt.Errorf("duplicate call_id %q", call.CallID)
		}
		calls[call.CallID] = call
		if !tsPattern.MatchString(call.CallTS) {
			return nil, fmt.Errorf("call %q: call_ts %q must have the format t<number>[.t<number>...]", call.CallID, call.CallTS)
		}
		// RPCs only have a service_call, and DB calls only have a database_call
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
	// validate the taints of the arguments and returns of each call (only once all calls are indexed,
	// since taints can reference any call of the function)
	for _, call := range fn.Calls {
		nodes := append([]*Node(nil), call.Arguments...)
		if call.ServiceCall != nil {
			nodes = append(nodes, call.ServiceCall.Returns...)
		}
		if err := validateNodes(nodes, calls); err != nil {
			return nil, fmt.Errorf("call %q: %w", call.CallID, err)
		}
	}
	// validate the taints of the params and returns of the function
	if err := validateNodes(fn.Params, calls); err != nil {
		return nil, fmt.Errorf("params: %w", err)
	}
	for i, rets := range fn.Returns {
		// all return statements must have the same number of objects, since they are merged by position
		if len(rets) != len(fn.Returns[0]) {
			return nil, fmt.Errorf("returns %d: expected %d objects but got %d", i, len(fn.Returns[0]), len(rets))
		}
		if err := validateNodes(rets, calls); err != nil {
			return nil, fmt.Errorf("returns %d: %w", i, err)
		}
	}
	return calls, nil
}

func validateNodes(nodes []*Node, calls map[string]*Call) error {
	for _, node := range nodes {
		if node == nil {
			return fmt.Errorf("nil node")
		}
		for _, taints := range node.Taints {
			for _, taint := range taints {
				if taint == nil {
					return fmt.Errorf("nil taint")
				}
				// taints must reference a call of the same function
				ref := calls[taint.CallID]
				if ref == nil {
					return fmt.Errorf("unknown taint call_id %q", taint.CallID)
				}
				if taint.CallerT != "" && !tsPattern.MatchString(taint.CallerT) {
					return fmt.Errorf("caller_t %q must have the format t<number>[.t<number>...]", taint.CallerT)
				}
				switch taint.TaintType {
				// database taints come from a DB call and their path is a field of its database
				case TaintDatabase:
					if ref.CallType != CallTypeDB {
						return fmt.Errorf("invalid database taint")
					}
					if !strings.HasPrefix(taint.Path, ref.DatabaseCall.Database+".") {
						return fmt.Errorf("taint path %q does not belong to database %q", taint.Path, ref.DatabaseCall.Database)
					}
				// service taints come from an RPC and their path is an argument or return of it
				case TaintService:
					if ref.CallType != CallTypeRPC || taint.DatabaseTaint != nil {
						return fmt.Errorf("invalid service taint")
					}
					if err := validateServiceTaintPath(taint.Path, ref); err != nil {
						return err
					}
				default:
					return fmt.Errorf("unknown taint_type %q", taint.TaintType)
				}
			}
		}
	}
	return nil
}

// validateServiceTaintPath checks that path has the format <service>.<method>.<name>[.<subpath>] (see
// abstractgraph.AbstractTrace), where name is an argument or return of call, which is used to match them
func validateServiceTaintPath(path string, call *Call) error {
	svc := call.ServiceCall
	prefix := svc.Service + "." + svc.Method + "."
	if !strings.HasPrefix(path, prefix) {
		return fmt.Errorf("taint path %q does not belong to %s.%s", path, svc.Service, svc.Method)
	}
	// extract the name of the object from the rest of the path (e.g., req from req.Field or req[*].Field)
	name, _, _ := strings.Cut(strings.TrimPrefix(path, prefix), ".")
	name, _, _ = strings.Cut(name, "[*]")
	// look for it in the arguments and returns of the call
	for _, obj := range append(append([]*Node(nil), call.Arguments...), svc.Returns...) {
		if obj != nil && obj.Name == name {
			return nil
		}
	}
	return fmt.Errorf("taint path %q does not name an argument or return of call %q", path, call.CallID)
}
