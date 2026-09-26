package components

import (
	"encoding/json"
	"fmt"
)

// ---------
// DATASTORE
// ---------

type DatastoreType int

const (
	DATASTORE_TYPE_RELATIONALDB DatastoreType = iota
	DATASTORE_TYPE_CACHE
	DATASTORE_TYPE_NOSQL
	DATASTORE_TYPE_QUEUE
)

type DatastoreKind int

const (
	DATASTORE_KIND_MYSQL DatastoreKind = iota
	DATASTORE_KIND_REDIS
	DATASTORE_KIND_MEMCACHED
	DATASTORE_KIND_MONGODB
	DATASTORE_KIND_RABBITMQ
)

type DatastoreInfo struct {
	Name string
	Type DatastoreType
	Kind DatastoreKind
}

func (info *DatastoreInfo) GetTypeLongName() string {
	return fmt.Sprintf("%s (%s)", info.GetKindString(), info.GetTypeString())
}

func (info *DatastoreInfo) GetTypeString() string {
	var typeToString = map[DatastoreType]string{
		DATASTORE_TYPE_RELATIONALDB: "RelationalDB",
		DATASTORE_TYPE_CACHE:        "Cache",
		DATASTORE_TYPE_NOSQL:        "NoSQLDatabase",
		DATASTORE_TYPE_QUEUE:        "Queue",
	}
	return typeToString[info.Type]
}

func (info *DatastoreInfo) GetKindString() string {
	var kindToString = map[DatastoreKind]string{
		DATASTORE_KIND_MYSQL:    "MySQL",
		DATASTORE_KIND_REDIS:    "Redis",
		DATASTORE_KIND_MONGODB:  "MongoDB",
		DATASTORE_KIND_RABBITMQ: "RabbitMQ",
	}
	return kindToString[info.Kind]
}

func (info *DatastoreInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Name string `json:"name"`
		Type string `json:"type"`
		Kind string `json:"kind"`
	}{
		Name: info.Name,
		Type: info.GetTypeString(),
		Kind: info.GetKindString(),
	})
}

func (info *DatastoreInfo) GetName() string {
	return info.Name
}

func (info *DatastoreInfo) IsRelationalDB() bool {
	return info.GetTypeString() == "RelationalDB"
}

func (info *DatastoreInfo) IsCache() bool {
	return info.GetTypeString() == "Cache"
}

func (info *DatastoreInfo) IsNoSQLDatabase() bool {
	return info.GetTypeString() == "NoSQLDatabase"
}

func (info *DatastoreInfo) IsQueue() bool {
	return info.GetTypeString() == "Queue"
}
