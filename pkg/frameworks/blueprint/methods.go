package blueprint

import (
	"fmt"

	"analyzer/pkg/datastores"
	"analyzer/pkg/logger"
	"analyzer/pkg/types"
	"analyzer/pkg/types/gotypes"
)

type OperationType int

const (
	OP_WRITE OperationType = iota
	OP_READ
	OP_DELETE
	OP_UPDATE
	OP_NOSQL_COLLECTION
)

type BackendMethod struct {
	types.Method
	Name      string
	Backend   string
	Component string
	Params    []*types.MethodField
	Returns   []*types.MethodField
	Operation OperationType
}

func (b *BackendMethod) Copy() *BackendMethod {
	return &BackendMethod{
		Name:      b.Name,
		Backend:   b.Backend,
		Component: b.Component,
		Params:    b.Params,
		Returns:   b.Returns,
		Operation: b.Operation,
	}
}

func (b *BackendMethod) DeepCopy() types.Method {
	return &BackendMethod{
		Name:      b.Name,
		Backend:   b.Backend,
		Component: b.Component,
		Params:    b.Params,
		Returns:   b.Returns,
		Operation: b.Operation,
	}
}

func (b *BackendMethod) GetName() string {
	return b.Name
}

func (b *BackendMethod) String() string {
	var repr string
	if b.Component != "" {
		repr = fmt.Sprintf("%s.%s.%s(", b.Backend, b.Component, b.Name)
	} else {
		repr = fmt.Sprintf("%s.%s(", b.Backend, b.Name)
	}
	for i, param := range b.Params {
		repr += param.String()
		if i < len(b.Params)-1 {
			repr += ", "
		}
	}
	repr += ")"
	return repr
}

func (b *BackendMethod) LongString() string {
	var repr string
	if b.Component != "" {
		repr = fmt.Sprintf("%s.%s.%s(", b.Backend, b.Component, b.Name)
	} else {
		repr = fmt.Sprintf("%s.%s(", b.Backend, b.Name)
	}
	for i, param := range b.Params {
		repr += param.String()
		if i < len(b.Params)-1 {
			repr += ", "
		}
	}
	repr += ")"
	if len(b.Returns) == 1 {
		repr += " " + b.Returns[0].String()
	}
	if len(b.Returns) > 1 {
		repr += " ("
		for i, ret := range b.Returns {
			repr += ret.String()
			if i < len(b.Returns)-1 {
				repr += ", "
			}
		}
		repr += ")"
	}
	return repr
}

func (b *BackendMethod) GetParams() []*types.MethodField {
	return b.Params
}

func (b *BackendMethod) GetReturns() []*types.MethodField {
	return b.Returns
}

func (b *BackendMethod) FullName() string {
	if b.Component != "" {
		return b.Backend + "." + b.Component + "." + b.Name
	}
	return b.Backend + "." + b.Name
}

// FIXME: this is messing up with previous assignments!
func (b *BackendMethod) SetNoSQLDatabaseCollection(databaseName string, collectionName string, dbInstance datastores.DatabaseInstance) {
	if len(b.Returns) > 0 {
		collection := b.Returns[0]
		collection.GetType().(*BlueprintBackendType).DatastoreInstance = dbInstance
		collection.GetType().(*BlueprintBackendType).NoSQLComponent = &NoSQLComponent{
			Database:   databaseName,
			Collection: collectionName,
			Type:       NoSQLCollectionType,
		}
		logger.Logger.Warnf("[BLUEPRINT] setting NoSQL collection for (%s, %s)", databaseName, collectionName)
	}
}

func (b *BackendMethod) IsNoSQLBackendCall() bool {
	return b.Backend == "NoSQLDatabase" && b.Component == ""
}

func (b *BackendMethod) IsRelationalDBSelectCall() bool {
	return b.Backend == "RelationalDB" && b.Component == "" && b.Name == "Select"
}

func (b *BackendMethod) IsRelationalDBQueryCall() bool {
	return b.Backend == "RelationalDB" && b.Component == "" && b.Name == "Query"
}

func (b *BackendMethod) IsRelationalDBExecCall() bool {
	return b.Backend == "RelationalDB" && b.Component == "" && b.Name == "Exec"
}

func (b *BackendMethod) IsNoSQLComponentCall() bool {
	return b.Backend == "NoSQLDatabase" && (b.Component == "NoSQLCollection" || b.Component == "NoSQLCursor")
}

func (b *BackendMethod) IsNoSQLCollectionCall() bool {
	return b.Backend == "NoSQLDatabase" && b.Component == "NoSQLCollection"
}

func (b *BackendMethod) IsNoSQLCursorCall() bool {
	return b.Backend == "NoSQLDatabase" && b.Component == "NoSQLCursor"
}

func (b *BackendMethod) ReturnsNoSQLCollection() (bool, int) {
	return b.Backend == "NoSQLDatabase" && b.Name == "GetCollection", 0 // index of collection
}

func (b *BackendMethod) ReturnsNoSQLCursor() (bool, int) {
	return b.Backend == "NoSQLDatabase" && b.Component == "NoSQLCollection" && (b.Name == "FindOne" || b.Name == "FindMany"), 0 // index of cursor
}

func (b *BackendMethod) IsWrite() bool {
	return b.Operation == OP_WRITE
}
func (b *BackendMethod) IsRead() bool {
	return b.Operation == OP_READ
}
func (b *BackendMethod) IsDelete() bool {
	return b.Operation == OP_DELETE
}
func (b *BackendMethod) IsUpdate() bool {
	return b.Operation == OP_UPDATE
}

func (b *BackendMethod) IsQueueRead() bool {
	return b.IsRead() && b.FullName() == "Queue.Pop"
}

func (b *BackendMethod) IsQueueWrite() bool {
	return b.IsWrite() && b.FullName() == "Queue.Push"
}

func (b *BackendMethod) MatchQueueIdentifiers() map[int]int {
	var matches map[int]int
	if b.FullName() == "Queue.Push" {
		matches = make(map[int]int, 0)
		// we have Queue.Push(ctx, src) and Queue.Pop(ctx, dst)
		// so the (src @ index 1) matches (dst @ index 1)
		matches[1] = 1
	}
	return matches
}

func (b *BackendMethod) GetWrittenObjectIndex() int {
	switch b.FullName() {
	case "Cache.Put":
		return 2
	case "NoSQLDatabase.NoSQLCollection.InsertOne", "NoSQLDatabase.NoSQLCollection.InsertMany", "NoSQLDatabase.NoSQLCollection.UpdateOne", "NoSQLDatabase.NoSQLCollection.UpdateMany":
		return 1
	case "NoSQLDatabase.NoSQLCollection.Upsert", "NoSQLDatabase.NoSQLCollection.ReplaceOne":
		return 2
	case "Queue.Push":
		return 1
	default:
		logger.Logger.Fatalf("unknown backend %s", b.FullName())
	}
	return -1
}

func (b *BackendMethod) GetReadObjectIndex() int {
	switch b.FullName() {
	case "RelationalDB.Select":
		return 1
	case "Cache.Get", "Cache.Mget":
		return 2
	case "NoSQLDatabase.NoSQLCollection.FindOne", "NoSQLDatabase.NoSQLCollection.FindMany":
		return -1
	case "Queue.Pop":
		return 1
	default:
		logger.Logger.Fatalf("unknown backend %s", b.FullName())
	}

	return -1
}

func (b *BackendMethod) GetWrittenKeyIndex() int {
	switch b.FullName() {
	case "Cache.Put":
		return 1
	case "NoSQLDatabase.NoSQLCollection.InsertOne", "NoSQLDatabase.NoSQLCollection.InsertMany", "NoSQLDatabase.NoSQLCollection.UpdateOne", "NoSQLDatabase.NoSQLCollection.UpdateMany",
		"NoSQLDatabase.NoSQLCollection.Upsert", "NoSQLDatabase.NoSQLCollection.ReplaceOne":
		return 1
	case "Queue.Push":
		return 1
	default:
		logger.Logger.Fatalf("unknown backend %s", b.FullName())
	}
	return -1
}

func (b *BackendMethod) GetReadKeyIndex() int {
	switch b.FullName() {
	case "Cache.Get", "Cache.Mget":
		return 1
	case "NoSQLDatabase.NoSQLCollection.FindOne", "NoSQLDatabase.NoSQLCollection.FindMany":
		return 1
	case "Queue.Pop":
		return 1
	default:
		logger.Logger.Fatalf("unknown backend %s", b.FullName())
	}
	return -1
}

func BuildBackendComponentMethods(name string) []*BackendMethod {
	var methods []*BackendMethod
	switch name {
	// --------
	// Backends
	// --------
	case "Cache":
		// Put(ctx context.Context, key string, value interface{}) error
		methods = append(methods, &BackendMethod{Name: "Put", Backend: "Cache", Operation: OP_WRITE,
			Params:  []*types.MethodField{&ctxParam, &keyParam, &valueParam},
			Returns: []*types.MethodField{&errorReturn},
		})
		// Get(ctx context.Context, key string, val interface{}) (bool, error)
		methods = append(methods, &BackendMethod{Name: "Get", Backend: "Cache", Operation: OP_READ,
			Params:  []*types.MethodField{&ctxParam, &keyParam, &valueParam},
			Returns: []*types.MethodField{&boolReturn, &errorReturn},
		})
		// Mget(ctx context.Context, keys []string, values []interface{}) error
		methods = append(methods, &BackendMethod{Name: "Mget", Backend: "Cache", Operation: OP_READ,
			Params:  []*types.MethodField{&ctxParam, &keysParam, &valuesParam},
			Returns: []*types.MethodField{&errorReturn},
		})
	case "Queue":
		// Push(ctx context.Context, item interface{}) (bool, error)
		methods = append(methods, &BackendMethod{Name: "Push", Backend: "Queue", Operation: OP_WRITE,
			Params:  []*types.MethodField{&ctxParam, &itemParam},
			Returns: []*types.MethodField{&boolReturn, &errorReturn},
		})
		// // Pop(ctx context.Context, dst interface{}) (bool, error)
		methods = append(methods, &BackendMethod{Name: "Pop", Backend: "Queue", Operation: OP_READ,
			Params:  []*types.MethodField{&ctxParam, &itemParam},
			Returns: []*types.MethodField{&boolReturn, &errorReturn},
		})
	case "NoSQLDatabase":
		// GetCollection(ctx context.Context, db_name string, collection_name string) (NoSQLCollection, error)
		methods = append(methods, &BackendMethod{Name: "GetCollection", Backend: "NoSQLDatabase", Operation: OP_NOSQL_COLLECTION,
			Params:  []*types.MethodField{&ctxParam, &dbNameParam, &collectionNameParam},
			Returns: []*types.MethodField{&NoSQLCollectionReturn, &errorReturn},
		})
	case "RelationalDB":
		// Select(ctx context.Context, dst interface{}, query string, args ...any) error
		methods = append(methods, &BackendMethod{Name: "Select", Backend: "RelationalDB", Operation: OP_READ,
			Params:  []*types.MethodField{&ctxParam, &dstParam, &queryParam, &argsParam},
			Returns: []*types.MethodField{&errorReturn},
		})
		// Exec(ctx context.Context, query string, args ...any) (sql.Result, error)
		methods = append(methods, &BackendMethod{Name: "Exec", Backend: "RelationalDB", Operation: OP_WRITE,
			Params:  []*types.MethodField{&ctxParam, &queryParam, &argsParam},
			Returns: []*types.MethodField{&errorReturn},
		})

	// ----------------
	// NoSQL Components
	// ----------------
	case "NoSQLCollection":
		methods = buildBackendNoSQLCollectionMethods()
	case "NoSQLCursor":
		methods = buildBackendNoSQLCursorMethods()
	default:
		logger.Logger.Fatalf("could not build methods for backend %s", name)
	}
	return methods
}

func buildBackendNoSQLCollectionMethods() []*BackendMethod {
	var methods []*BackendMethod
	// FindOne(ctx context.Context, filter bson.D, projection ...bson.D) (NoSQLCursor, error)
	methods = append(methods, &BackendMethod{Name: "FindOne", Backend: "NoSQLDatabase", Component: "NoSQLCollection", Operation: OP_READ,
		Params:  []*types.MethodField{&ctxParam, &filterParam, &projectionParam},
		Returns: []*types.MethodField{&NoSQLCursorReturn, &errorReturn},
	})
	// FindMany(ctx context.Context, filter bson.D, projection ...bson.D) (NoSQLCursor, error)
	methods = append(methods, &BackendMethod{Name: "FindMany", Backend: "NoSQLDatabase", Component: "NoSQLCollection", Operation: OP_READ,
		Params:  []*types.MethodField{&ctxParam, &filterParam, &projectionParam},
		Returns: []*types.MethodField{&NoSQLCursorReturn, &errorReturn},
	})
	// Upsert(ctx context.Context, filter bson.D, document interface{}) (bool, error)
	methods = append(methods, &BackendMethod{Name: "Upsert", Backend: "NoSQLDatabase", Component: "NoSQLCollection", Operation: OP_UPDATE,
		Params:  []*types.MethodField{&ctxParam, &filterParam, &docParam},
		Returns: []*types.MethodField{&boolReturn, &errorReturn},
	})
	// UpdateOne(ctx context.Context, filter bson.D, update bson.D) (int, error)
	methods = append(methods, &BackendMethod{Name: "UpdateOne", Backend: "NoSQLDatabase", Component: "NoSQLCollection", Operation: OP_UPDATE,
		Params:  []*types.MethodField{&ctxParam, &filterParam, &updateParam},
		Returns: []*types.MethodField{&intReturn, &errorReturn},
	})
	// UpdateMany(ctx context.Context, filter bson.D, update bson.D) (int, error)
	methods = append(methods, &BackendMethod{Name: "UpdateMany", Backend: "NoSQLDatabase", Component: "NoSQLCollection", Operation: OP_UPDATE,
		Params:  []*types.MethodField{&ctxParam, &filterParam, &updateParam},
		Returns: []*types.MethodField{&intReturn, &errorReturn},
	})
	// ReplaceOne(ctx context.Context, filter bson.D, replacement{}) (int, error)
	methods = append(methods, &BackendMethod{Name: "ReplaceOne", Backend: "NoSQLDatabase", Component: "NoSQLCollection", Operation: OP_UPDATE,
		Params:  []*types.MethodField{&ctxParam, &filterParam, &replacementParam},
		Returns: []*types.MethodField{&intReturn, &errorReturn},
	})
	// InsertOne(ctx context.Context, document interface{}) error
	methods = append(methods, &BackendMethod{Name: "InsertOne", Backend: "NoSQLDatabase", Component: "NoSQLCollection", Operation: OP_WRITE,
		Params:  []*types.MethodField{&ctxParam, &docParam},
		Returns: []*types.MethodField{&errorReturn},
	})
	// InsertMany(ctx context.Context, documents []interface{}) error
	methods = append(methods, &BackendMethod{Name: "InsertMany", Backend: "NoSQLDatabase", Component: "NoSQLCollection", Operation: OP_WRITE,
		Params:  []*types.MethodField{&ctxParam, &docsParam},
		Returns: []*types.MethodField{&errorReturn},
	})
	// DeleteOne(ctx context.Context, filter bson.D) error
	methods = append(methods, &BackendMethod{Name: "DeleteOne", Backend: "NoSQLDatabase", Component: "NoSQLCollection", Operation: OP_DELETE,
		Params:  []*types.MethodField{&ctxParam, &filterParam},
		Returns: []*types.MethodField{&errorReturn},
	})
	// DeleteMany(ctx context.Context, filter bson.D) error
	methods = append(methods, &BackendMethod{Name: "DeleteMany", Backend: "NoSQLDatabase", Component: "NoSQLCollection", Operation: OP_DELETE,
		Params:  []*types.MethodField{&ctxParam, &filterParam},
		Returns: []*types.MethodField{&errorReturn},
	})
	return methods
}

func buildBackendNoSQLCursorMethods() []*BackendMethod {
	var methods []*BackendMethod
	// One(ctx context.Context, obj interface{}) (bool, error)
	methods = append(methods, &BackendMethod{Name: "One", Backend: "NoSQLDatabase", Component: "NoSQLCursor", Operation: OP_READ,
		Params:  []*types.MethodField{&ctxParam, &objParam},
		Returns: []*types.MethodField{&boolReturn, &errorReturn},
	})
	// All(ctx context.Context, obj interface{}) error //similar logic to Decode, but for multiple documents
	methods = append(methods, &BackendMethod{Name: "All", Backend: "NoSQLDatabase", Component: "NoSQLCursor", Operation: OP_WRITE,
		Params:  []*types.MethodField{&ctxParam, &objParam},
		Returns: []*types.MethodField{&errorReturn},
	})
	return methods
}

var ctxParam = types.MethodField{
	FieldInfo: types.FieldInfo{
		Name: "ctx",
		Type: &gotypes.UserType{
			Name:        "Context",
			PackagePath: "context",
			UserType:    nil,
		},
	},
}
var keyParam = types.MethodField{
	FieldInfo: types.FieldInfo{
		Name: "key",
		Type: &gotypes.BasicType{
			Name: "string",
		},
	},
}
var keysParam = types.MethodField{
	FieldInfo: types.FieldInfo{
		Name: "key",
		Type: &gotypes.SliceType{
			UnderlyingType: &gotypes.BasicType{
				Name: "string",
			},
		},
	},
}
var valueParam = types.MethodField{
	FieldInfo: types.FieldInfo{
		Name: "value",
		Type: &gotypes.InterfaceType{Methods: make(map[string]string)},
	},
}
var valuesParam = types.MethodField{
	FieldInfo: types.FieldInfo{
		Name: "key",
		Type: &gotypes.SliceType{
			UnderlyingType: &gotypes.InterfaceType{Methods: make(map[string]string)},
		},
	},
}
var itemParam = types.MethodField{
	FieldInfo: types.FieldInfo{
		Name: "item",
		Type: &gotypes.InterfaceType{Methods: make(map[string]string)},
	},
}
var docParam = types.MethodField{
	FieldInfo: types.FieldInfo{
		Name: "document",
		Type: &gotypes.InterfaceType{Methods: make(map[string]string)},
	},
}
var dstParam = types.MethodField{
	FieldInfo: types.FieldInfo{
		Name: "dst",
		Type: &gotypes.InterfaceType{Methods: make(map[string]string)},
	},
}
var argsParam = types.MethodField{
	FieldInfo: types.FieldInfo{
		Name: "dst",
		Type: &gotypes.SliceType{
			UnderlyingType:  &gotypes.InterfaceType{Methods: make(map[string]string)},
		},
	},
}
var docsParam = types.MethodField{
	FieldInfo: types.FieldInfo{
		Name: "key",
		Type: &gotypes.SliceType{
			UnderlyingType: &gotypes.InterfaceType{Methods: make(map[string]string)},
		},
	},
}
var objParam = types.MethodField{
	FieldInfo: types.FieldInfo{
		Name: "obj",
		Type: &gotypes.InterfaceType{Methods: make(map[string]string)},
	},
}
var filterParam = types.MethodField{
	FieldInfo: types.FieldInfo{
		Name: "filter",
		Type: &gotypes.UserType{
			Name:        "D",
			PackagePath: "bson",
		},
	},
}
var updateParam = types.MethodField{
	FieldInfo: types.FieldInfo{
		Name: "update",
		Type: &gotypes.UserType{
			Name:        "D",
			PackagePath: "bson",
		},
	},
}
var projectionParam = types.MethodField{
	FieldInfo: types.FieldInfo{
		Name: "projection",
		Type: &gotypes.UserType{
			Name:        "D",
			PackagePath: "bson",
		},
	},
}
var dbNameParam = types.MethodField{
	FieldInfo: types.FieldInfo{
		Name: "db_name",
		Type: &gotypes.BasicType{
			Name: "string",
		},
	},
}
var collectionNameParam = types.MethodField{
	FieldInfo: types.FieldInfo{
		Name: "collection_name",
		Type: &gotypes.BasicType{
			Name: "string",
		},
	},
}
var queryParam = types.MethodField{
	FieldInfo: types.FieldInfo{
		Name: "query",
		Type: &gotypes.BasicType{
			Name: "string",
		},
	},
}
var boolReturn = types.MethodField{
	FieldInfo: types.FieldInfo{
		Type: &gotypes.BasicType{
			Name: "err",
		},
	},
}
var NoSQLCursorReturn = types.MethodField{
	FieldInfo: types.FieldInfo{
		Type: &BlueprintBackendType{
			Name:    "NoSQLCursor",
			Package: "github.com/blueprint-uservices/blueprint/runtime/core/backend",
			Methods: buildBackendNoSQLCursorMethods(),
		},
	},
}
var NoSQLCollectionReturn = types.MethodField{
	FieldInfo: types.FieldInfo{
		Type: &BlueprintBackendType{
			Name:    "NoSQLCollection",
			Package: "github.com/blueprint-uservices/blueprint/runtime/core/backend",
			Methods: buildBackendNoSQLCollectionMethods(),
		},
	},
}
var errorReturn = types.MethodField{
	FieldInfo: types.FieldInfo{
		// error is actually an interface
		Type: &gotypes.InterfaceType{Methods: make(map[string]string)},
	},
}
var replacementParam = types.MethodField{
	FieldInfo: types.FieldInfo{
		Type: &gotypes.InterfaceType{Methods: make(map[string]string)},
	},
}
var intReturn = types.MethodField{
	FieldInfo: types.FieldInfo{
		// error is actually an interface
		Type: &gotypes.BasicType{
			Name: "int",
		},
	},
}
