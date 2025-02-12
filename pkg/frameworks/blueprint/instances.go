package blueprint

import (
	"encoding/json"

	"analyzer/pkg/datastores"
)

type BlueprintDatabaseInstance struct {
	datastores.DatabaseInstance `json:"-"`
	Name                        string                `json:"name"`
	Datastore                   *datastores.Datastore `json:"datastore"`
}

func (i *BlueprintDatabaseInstance) GetDatastore() *datastores.Datastore {
	return i.Datastore
}

type QueueInstance struct {
	BlueprintDatabaseInstance
}

func (q *QueueInstance) GetName() string {
	return q.Name
}

func (q *QueueInstance) String() string {
	return q.Name + " " + q.Datastore.GetKindString()
}

func (q *QueueInstance) LongString() string {
	return q.Name + " " + q.GetTypeLongName()
}

func (q *QueueInstance) IsQueue() bool {
	return true
}

func (q *QueueInstance) GetTypeName() string {
	return "Queue"
}

func (q *QueueInstance) GetTypeLongName() string {
	return q.GetDatastore().GetTypeLongName()
}

// MarshalJSON is used by app.Save()
func (q *QueueInstance) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Name      string                `json:"name"`
		Datastore *datastores.Datastore `json:"datastore"`
	}{
		Name:      q.GetName(),
		Datastore: q.GetDatastore(),
	})
}

type CacheInstance struct {
	BlueprintDatabaseInstance
}

func (cache *CacheInstance) GetName() string {
	return cache.Name
}

func (cache *CacheInstance) String() string {
	return cache.Name + " " + cache.Datastore.GetKindString()
}

func (cache *CacheInstance) LongString() string {
	return cache.Name + " " + cache.GetTypeLongName()
}

func (cache *CacheInstance) IsQueue() bool {
	return false
}

func (cache *CacheInstance) GetTypeName() string {
	return "Cache"
}

func (cache *CacheInstance) GetTypeLongName() string {
	return cache.GetDatastore().GetTypeLongName()
}

// MarshalJSON is used by app.Save()
func (cache *CacheInstance) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Name      string                `json:"name"`
		Datastore *datastores.Datastore `json:"datastore"`
	}{
		Name:      cache.GetName(),
		Datastore: cache.GetDatastore(),
	})
}

type NoSQLInstance struct {
	BlueprintDatabaseInstance
}

func (nosql *NoSQLInstance) GetName() string {
	return nosql.Name
}

func (nosql *NoSQLInstance) String() string {
	return nosql.Name + " " + nosql.Datastore.GetKindString()
}

func (nosql *NoSQLInstance) LongString() string {
	return nosql.Name + " " + nosql.GetTypeLongName()
}

func (nosql *NoSQLInstance) IsQueue() bool {
	return false
}

func (nosql *NoSQLInstance) GetTypeName() string {
	return "NoSQLDatabase"
}

func (nosql *NoSQLInstance) GetTypeLongName() string {
	return nosql.GetDatastore().GetTypeLongName()
}

// MarshalJSON is used by app.Save()
func (nosql *NoSQLInstance) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Name      string                `json:"name"`
		Datastore *datastores.Datastore `json:"datastore"`
	}{
		Name:      nosql.GetName(),
		Datastore: nosql.GetDatastore(),
	})
}

type RelationalDBInstance struct {
	BlueprintDatabaseInstance
}

func (reldb *RelationalDBInstance) GetName() string {
	return reldb.Name
}

func (reldb *RelationalDBInstance) String() string {
	return reldb.Name + " " + reldb.Datastore.GetKindString()
}

func (reldb *RelationalDBInstance) LongString() string {
	return reldb.Name + " " + reldb.GetTypeLongName()
}

func (reldb *RelationalDBInstance) IsQueue() bool {
	return false
}

func (reldb *RelationalDBInstance) GetTypeName() string {
	return "RelationalDB"
}

func (reldb *RelationalDBInstance) GetTypeLongName() string {
	return reldb.GetDatastore().GetTypeLongName()
}

// MarshalJSON is used by app.Save()
func (reldb *RelationalDBInstance) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Name      string                `json:"name"`
		Datastore *datastores.Datastore `json:"datastore"`
	}{
		Name:      reldb.GetName(),
		Datastore: reldb.GetDatastore(),
	})
}
