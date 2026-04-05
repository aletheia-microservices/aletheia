package blueprint

import (
	"fmt"

	"github.com/sirupsen/logrus"
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
	Name              string
	Backend           string
	Component         string
	Params            []*MethodField
	Returns           []*MethodField
	Operation         OperationType
	calledBackendType *BlueprintBackendType
}

func (b *BackendMethod) SetCalledBackendType(t *BlueprintBackendType) {
	b.calledBackendType = t
}

func (b *BackendMethod) GetCalledBackendType() *BlueprintBackendType {
	return b.calledBackendType
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

func (b *BackendMethod) GetParams() []*MethodField {
	return b.Params
}

func (b *BackendMethod) GetReturns() []*MethodField {
	return b.Returns
}

func (b *BackendMethod) FullName() string {
	if b.Component != "" {
		return b.Backend + "." + b.Component + "." + b.Name
	}
	return b.Backend + "." + b.Name
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
		logrus.Fatalf("unknown backend %s", b.FullName())
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
		logrus.Fatalf("unknown backend %s", b.FullName())
	}

	return -1
}

func (b *BackendMethod) GetWrittenKeyIndex() int {
	switch b.FullName() {
	case "Cache.Put":
		return 1
	case "Queue.Push":
		return 1
	default:
		logrus.Fatalf("unknown backend %s", b.FullName())
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
		logrus.Fatalf("unknown backend %s", b.FullName())
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
			Params:  []*MethodField{&ctxParam, &keyParam, &valueParam},
			Returns: []*MethodField{&errorReturn},
		})
		// Get(ctx context.Context, key string, val interface{}) (bool, error)
		methods = append(methods, &BackendMethod{Name: "Get", Backend: "Cache", Operation: OP_READ,
			Params:  []*MethodField{&ctxParam, &keyParam, &valueParam},
			Returns: []*MethodField{&boolReturn, &errorReturn},
		})
		// Mget(ctx context.Context, keys []string, values []interface{}) error
		methods = append(methods, &BackendMethod{Name: "Mget", Backend: "Cache", Operation: OP_READ,
			Params:  []*MethodField{&ctxParam, &keysParam, &valuesParam},
			Returns: []*MethodField{&errorReturn},
		})
	case "Queue":
		// Push(ctx context.Context, item interface{}) (bool, error)
		methods = append(methods, &BackendMethod{Name: "Push", Backend: "Queue", Operation: OP_WRITE,
			Params:  []*MethodField{&ctxParam, &itemParam},
			Returns: []*MethodField{&boolReturn, &errorReturn},
		})
		// // Pop(ctx context.Context, dst interface{}) (bool, error)
		methods = append(methods, &BackendMethod{Name: "Pop", Backend: "Queue", Operation: OP_READ,
			Params:  []*MethodField{&ctxParam, &itemParam},
			Returns: []*MethodField{&boolReturn, &errorReturn},
		})
	case "NoSQLDatabase":
		// GetCollection(ctx context.Context, db_name string, collection_name string) (NoSQLCollection, error)
		methods = append(methods, &BackendMethod{Name: "GetCollection", Backend: "NoSQLDatabase", Operation: OP_NOSQL_COLLECTION,
			Params:  []*MethodField{&ctxParam, &dbNameParam, &collectionNameParam},
			Returns: []*MethodField{&NoSQLCollectionReturn, &errorReturn},
		})
	case "RelationalDB":
		// Select(ctx context.Context, dst interface{}, query string, args ...any) error
		methods = append(methods, &BackendMethod{Name: "Select", Backend: "RelationalDB", Operation: OP_READ,
			Params:  []*MethodField{&ctxParam, &dstParam, &queryParam, &argsParam},
			Returns: []*MethodField{&errorReturn},
		})
		// Exec(ctx context.Context, query string, args ...any) (sql.Result, error)
		methods = append(methods, &BackendMethod{Name: "Exec", Backend: "RelationalDB", Operation: OP_WRITE,
			Params:  []*MethodField{&ctxParam, &queryParam, &argsParam},
			Returns: []*MethodField{&errorReturn},
		})

	// ----------------
	// NoSQL Components
	// ----------------
	case "NoSQLCollection":
		methods = buildBackendNoSQLCollectionMethods()
	case "NoSQLCursor":
		methods = buildBackendNoSQLCursorMethods()
	default:
		logrus.Fatalf("could not build methods for backend %s", name)
	}
	return methods
}

func buildBackendNoSQLCollectionMethods() []*BackendMethod {
	var methods []*BackendMethod
	// FindOne(ctx context.Context, filter bson.D, projection ...bson.D) (NoSQLCursor, error)
	methods = append(methods, &BackendMethod{Name: "FindOne", Backend: "NoSQLDatabase", Component: "NoSQLCollection", Operation: OP_READ,
		Params:  []*MethodField{&ctxParam, &filterParam, &projectionParam},
		Returns: []*MethodField{&NoSQLCursorReturn, &errorReturn},
	})
	// FindMany(ctx context.Context, filter bson.D, projection ...bson.D) (NoSQLCursor, error)
	methods = append(methods, &BackendMethod{Name: "FindMany", Backend: "NoSQLDatabase", Component: "NoSQLCollection", Operation: OP_READ,
		Params:  []*MethodField{&ctxParam, &filterParam, &projectionParam},
		Returns: []*MethodField{&NoSQLCursorReturn, &errorReturn},
	})
	// Upsert(ctx context.Context, filter bson.D, document interface{}) (bool, error)
	methods = append(methods, &BackendMethod{Name: "Upsert", Backend: "NoSQLDatabase", Component: "NoSQLCollection", Operation: OP_UPDATE,
		Params:  []*MethodField{&ctxParam, &filterParam, &docParam},
		Returns: []*MethodField{&boolReturn, &errorReturn},
	})
	// UpdateOne(ctx context.Context, filter bson.D, update bson.D) (int, error)
	methods = append(methods, &BackendMethod{Name: "UpdateOne", Backend: "NoSQLDatabase", Component: "NoSQLCollection", Operation: OP_UPDATE,
		Params:  []*MethodField{&ctxParam, &filterParam, &updateParam},
		Returns: []*MethodField{&intReturn, &errorReturn},
	})
	// UpdateMany(ctx context.Context, filter bson.D, update bson.D) (int, error)
	methods = append(methods, &BackendMethod{Name: "UpdateMany", Backend: "NoSQLDatabase", Component: "NoSQLCollection", Operation: OP_UPDATE,
		Params:  []*MethodField{&ctxParam, &filterParam, &updateParam},
		Returns: []*MethodField{&intReturn, &errorReturn},
	})
	// ReplaceOne(ctx context.Context, filter bson.D, replacement{}) (int, error)
	methods = append(methods, &BackendMethod{Name: "ReplaceOne", Backend: "NoSQLDatabase", Component: "NoSQLCollection", Operation: OP_UPDATE,
		Params:  []*MethodField{&ctxParam, &filterParam, &replacementParam},
		Returns: []*MethodField{&intReturn, &errorReturn},
	})
	// InsertOne(ctx context.Context, document interface{}) error
	methods = append(methods, &BackendMethod{Name: "InsertOne", Backend: "NoSQLDatabase", Component: "NoSQLCollection", Operation: OP_WRITE,
		Params:  []*MethodField{&ctxParam, &docParam},
		Returns: []*MethodField{&errorReturn},
	})
	// InsertMany(ctx context.Context, documents []interface{}) error
	methods = append(methods, &BackendMethod{Name: "InsertMany", Backend: "NoSQLDatabase", Component: "NoSQLCollection", Operation: OP_WRITE,
		Params:  []*MethodField{&ctxParam, &docsParam},
		Returns: []*MethodField{&errorReturn},
	})
	// DeleteOne(ctx context.Context, filter bson.D) error
	methods = append(methods, &BackendMethod{Name: "DeleteOne", Backend: "NoSQLDatabase", Component: "NoSQLCollection", Operation: OP_DELETE,
		Params:  []*MethodField{&ctxParam, &filterParam},
		Returns: []*MethodField{&errorReturn},
	})
	// DeleteMany(ctx context.Context, filter bson.D) error
	methods = append(methods, &BackendMethod{Name: "DeleteMany", Backend: "NoSQLDatabase", Component: "NoSQLCollection", Operation: OP_DELETE,
		Params:  []*MethodField{&ctxParam, &filterParam},
		Returns: []*MethodField{&errorReturn},
	})
	return methods
}

func buildBackendNoSQLCursorMethods() []*BackendMethod {
	var methods []*BackendMethod
	// One(ctx context.Context, obj interface{}) (bool, error)
	methods = append(methods, &BackendMethod{Name: "One", Backend: "NoSQLDatabase", Component: "NoSQLCursor", Operation: OP_READ,
		Params:  []*MethodField{&ctxParam, &objParam},
		Returns: []*MethodField{&boolReturn, &errorReturn},
	})
	// All(ctx context.Context, obj interface{}) error //similar logic to Decode, but for multiple documents
	methods = append(methods, &BackendMethod{Name: "All", Backend: "NoSQLDatabase", Component: "NoSQLCursor", Operation: OP_WRITE,
		Params:  []*MethodField{&ctxParam, &objParam},
		Returns: []*MethodField{&errorReturn},
	})
	return methods
}

var ctxParam = MethodField{
	FieldInfo: FieldInfo{
		Name: "ctx",
	},
}
var keyParam = MethodField{
	FieldInfo: FieldInfo{
		Name: "key",
	},
}
var keysParam = MethodField{
	FieldInfo: FieldInfo{
		Name: "key",
	},
}
var valueParam = MethodField{
	FieldInfo: FieldInfo{
		Name: "value",
	},
}
var valuesParam = MethodField{
	FieldInfo: FieldInfo{
		Name: "key",
	},
}
var itemParam = MethodField{
	FieldInfo: FieldInfo{
		Name: "item",
	},
}
var docParam = MethodField{
	FieldInfo: FieldInfo{
		Name: "document",
	},
}
var dstParam = MethodField{
	FieldInfo: FieldInfo{
		Name: "dst",
	},
}
var argsParam = MethodField{
	FieldInfo: FieldInfo{
		Name: "dst",
	},
}
var docsParam = MethodField{
	FieldInfo: FieldInfo{
		Name: "key",
	},
}
var objParam = MethodField{
	FieldInfo: FieldInfo{
		Name: "obj",
	},
}
var filterParam = MethodField{
	FieldInfo: FieldInfo{
		Name: "filter",
	},
}
var updateParam = MethodField{
	FieldInfo: FieldInfo{
		Name: "update",
	},
}
var projectionParam = MethodField{
	FieldInfo: FieldInfo{
		Name: "projection",
	},
}
var dbNameParam = MethodField{
	FieldInfo: FieldInfo{
		Name: "db_name",
	},
}
var collectionNameParam = MethodField{
	FieldInfo: FieldInfo{
		Name: "collection_name",
	},
}
var queryParam = MethodField{
	FieldInfo: FieldInfo{
		Name: "query",
	},
}
var boolReturn = MethodField{
	FieldInfo: FieldInfo{
		Name: "err",
	},
}
var NoSQLCursorReturn = MethodField{
	FieldInfo: FieldInfo{},
}
var NoSQLCollectionReturn = MethodField{
	FieldInfo: FieldInfo{},
}
var errorReturn = MethodField{}
var replacementParam = MethodField{}
var intReturn = MethodField{
	FieldInfo: FieldInfo{
		// error is actually an interface
	},
}
