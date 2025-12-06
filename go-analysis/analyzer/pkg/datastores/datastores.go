package datastores

import (
	"encoding/json"
	"fmt"
	"slices"
)

type DatastoreType int

const (
	RelationalDB DatastoreType = iota
	Cache
	NoSQL
	Queue
)

type DatastoreKind int

const (
	MySQL DatastoreKind = iota
	Redis
	Memcached
	MongoDB
	RabbitMQ
)

type Datastore struct {
	Name                  string
	Type                  DatastoreType
	Kind                  DatastoreKind
	Schema                *Schema
	ReferencingDatastores []*Datastore
}

func (ds *Datastore) GetTypeLongName() string {
	return fmt.Sprintf("%s (%s)", ds.GetKindString(), ds.GetTypeString())
}

func (ds *Datastore) GetSchema() *Schema {
	return ds.Schema
}

func (ds *Datastore) IsReferencingDatastore(dependency *Datastore) bool {
	return slices.Contains(ds.ReferencingDatastores, dependency)
}

func (ds *Datastore) AddReferencingDatastoreIfNotExists(newDep *Datastore) {
	for _, dep := range ds.ReferencingDatastores {
		if dep == newDep {
			return
		}
	}
	ds.ReferencingDatastores = append(ds.ReferencingDatastores, newDep)
}

func (ds *Datastore) GetTypeString() string {
	var typeToString = map[DatastoreType]string{
		RelationalDB: "RelationalDB",
		Cache:        "Cache",
		NoSQL:        "NoSQLDatabase",
		Queue:        "Queue",
	}
	return typeToString[ds.Type]
}

func (ds *Datastore) GetKindString() string {
	var kindToString = map[DatastoreKind]string{
		MySQL:    "MySQL",
		Redis:    "Redis",
		MongoDB:  "MongoDB",
		RabbitMQ: "RabbitMQ",
	}
	return kindToString[ds.Kind]
}

func (ds *Datastore) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Name   string  `json:"name"`
		Type   string  `json:"type"`
		Kind   string  `json:"kind"`
		Schema *Schema `json:"-"`
	}{
		Name:   ds.Name,
		Type:   ds.GetTypeString(),
		Kind:   ds.GetKindString(),
		Schema: ds.Schema,
	})
}

func (ds *Datastore) GetName() string {
	return ds.Name
}
func (ds *Datastore) IsRelationalDB() bool {
	return ds.GetTypeString() == "RelationalDB"
}
func (ds *Datastore) IsCache() bool {
	return ds.GetTypeString() == "Cache"
}
func (ds *Datastore) IsNoSQLDatabase() bool {
	return ds.GetTypeString() == "NoSQLDatabase"
}
func (ds *Datastore) IsQueue() bool {
	return ds.GetTypeString() == "Queue"
}
