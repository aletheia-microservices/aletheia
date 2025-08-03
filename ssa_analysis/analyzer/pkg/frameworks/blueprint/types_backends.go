package blueprint

import (
	"fmt"
	"log"

	"analyzer/pkg/frameworks/components"
)

func IsBackendComponent(name string) bool {
	return IsBackend(name) || IsNoSQLComponent(name)
}

func IsBackend(name string) bool {
	switch name {
	case "Queue", "Cache", "NoSQLDatabase", "RelationalDB":
		return true
	}
	return false
}

func IsNoSQLComponent(name string) bool {
	switch name {
	case "NoSQLCollection", "NoSQLCursor":
		return true
	}
	return false
}

type NoSQLComponentType int

const (
	NoSQLCollectionType NoSQLComponentType = iota
	NoSQLCursorType
)

type NoSQLComponent struct {
	Type       NoSQLComponentType
	Projection []string // fields to be projected in queries - applies only if Type is NoSQLCursorType
	Database   string
	Collection string
}

func (t *NoSQLComponent) HasProjection() bool {
	return t.Projection != nil
}

func (t *NoSQLComponent) GetProjection() []string {
	return t.Projection
}

func (t *NoSQLComponent) Copy() *NoSQLComponent {
	return &NoSQLComponent{
		Type:       t.Type,
		Database:   t.Database,
		Collection: t.Collection,
	}
}

func (t *NoSQLComponent) String() string {
	prefix := ""
	if t.Type == NoSQLCollectionType {
		prefix = "NoSQLCollection"
	} else if t.Type == NoSQLCursorType {
		prefix = "NoSQLCursor"
	}
	return fmt.Sprintf("%s {database = %s, collection = %s}", prefix, t.Database, t.Collection)
}

func (t *NoSQLComponent) LongString() string {
	return t.String()
}

type BlueprintBackendType struct {
	Name           string
	Package        string
	Methods        []*BackendMethod
	Datastore      *components.DatastoreInfo
	NoSQLComponent *NoSQLComponent
}

func (t *BlueprintBackendType) SetInstance(instance *components.DatastoreInfo) {
	t.Datastore = instance
}

func (t *BlueprintBackendType) String() string {
	if t.NoSQLComponent != nil {
		return t.NoSQLComponent.String()
	}
	return t.Name
}

func (t *BlueprintBackendType) StringWithInstance() string {
	if t.NoSQLComponent != nil {
		return t.NoSQLComponent.String()
	}
	instance := "<nil>"
	if t.Datastore != nil {
		instance = t.Datastore.GetName()
	}
	return fmt.Sprintf("%s {instance = %s}", t.Name, instance)
}

func (t *BlueprintBackendType) LongString() string {
	if t.NoSQLComponent != nil {
		return t.NoSQLComponent.LongString()
	}
	s := t.Name
	if t.Datastore != nil {
		s += fmt.Sprintf(" {instance = %s}", t.Datastore.GetName())
	}
	if len(t.Methods) == 0 {
		return s + " interface{}"
	}
	s += " interface{"
	for i, m := range t.Methods {
		s += m.String()
		if i < len(t.Methods)-1 {
			s += ", "
		}
	}
	return s + "}"
}

func (t *BlueprintBackendType) StringWithMethodsList() string {
	if t.NoSQLComponent != nil {
		return t.NoSQLComponent.String()
	}

	s := t.Name
	if len(t.Methods) == 0 {
		return s + " interface{}"
	}
	s += "\n" + t.Name + " interface{\n"
	for _, m := range t.Methods {
		s += "\t" + m.String() + "\n"
	}
	return s + "}"
}

func (t *BlueprintBackendType) GetName() string {
	return t.Name
}

func (t *BlueprintBackendType) GetLongName() string {
	return t.Package + "." + t.Name
}

func (t *BlueprintBackendType) GetPackage() string {
	return t.Package
}

func (t *BlueprintBackendType) GetBasicValue() string {
	log.Fatalf("[TYPES BLUEPRINT] unable to get value for blueprint backend type type %s", t.String())
	return ""
}

func (t *BlueprintBackendType) AddValue(value string) {
	log.Fatalf("[TYPES BLUEPRINT] unable to add value for blueprint backend type type %s", t.String())
}

func (t *BlueprintBackendType) IsNoSQLComponent() bool {
	return t.NoSQLComponent != nil
}

func (t *BlueprintBackendType) IsNoSQLCollection() bool {
	return t.NoSQLComponent != nil && t.NoSQLComponent.Type == NoSQLCollectionType
}

func (t *BlueprintBackendType) IsNoSQLCursor() bool {
	return t.NoSQLComponent != nil && t.NoSQLComponent.Type == NoSQLCursorType
}

func (t *BlueprintBackendType) IsQueue() bool {
	return t.Name == "Queue"
}

func (t *BlueprintBackendType) IsNoSQLDatabase() bool {
	return t.Name == "NoSQLDatabase"
}

func (t *BlueprintBackendType) GetMethods() []*BackendMethod {
	return t.Methods
}

func (t *BlueprintBackendType) GetMethod(name string) *BackendMethod {
	for _, m := range t.Methods {
		if m.Name == name {
			return m
		}
	}
	log.Fatalf("[TYPES BLUEPRINT] could not find method (%s) for backend type (%s) with methods (%v)", name, t.String(), t.Methods)
	return nil
}

func (t *BlueprintBackendType) SetNoSQLDatabaseCollection(databaseName string, collectionName string, dbInstance *components.DatastoreInfo) {
	t.NoSQLComponent = &NoSQLComponent{
		Type:       NoSQLCollectionType,
		Database:   databaseName,
		Collection: collectionName,
	}
	t.Datastore = dbInstance
}

func (t *BlueprintBackendType) SetNoSQLDatabaseCursor(databaseName string, collectionName string, dbInstance *components.DatastoreInfo) {
	t.NoSQLComponent = &NoSQLComponent{
		Type:       NoSQLCursorType,
		Database:   databaseName,
		Collection: collectionName,
	}
	t.Datastore = dbInstance
}

func (t *BlueprintBackendType) Copy(force bool) *BlueprintBackendType {
	var methods []*BackendMethod
	for _, m := range t.Methods {
		methods = append(methods, m.Copy())
	}
	var noSQLComponent *NoSQLComponent
	if t.NoSQLComponent != nil {
		noSQLComponent = t.NoSQLComponent.Copy()
	}
	return &BlueprintBackendType{
		Name:           t.Name,
		Package:        t.Package,
		Methods:        methods,
		Datastore:      t.Datastore,
		NoSQLComponent: noSQLComponent,
	}
}
