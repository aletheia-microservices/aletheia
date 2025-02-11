package cascade

import (
	"fmt"

	"analyzer/pkg/datastores"
	"analyzer/pkg/service"
)

type deleteDependency struct {
	service   *service.Service
	datastore *datastores.Datastore
	cascading bool
}

func (dep *deleteDependency) hasServiceAndDatastore(serviceName string, datastore *datastores.Datastore) bool {
	return dep.service.GetName() == serviceName && dep.datastore == datastore
}

func (dep *deleteDependency) isEqual(other *deleteDependency) bool {
	return dep.service == other.service && dep.datastore == other.datastore
}

func (dep *deleteDependency) String() string {
	return fmt.Sprintf("(%s, %s)", dep.service.Name, dep.datastore.GetName())
}

func (dep *deleteDependency) LongString() string {
	return fmt.Sprintf("(%s, %s)", dep.service.Name, dep.datastore.GetName())
}

func (dep *deleteDependency) GetServiceName() string {
	return dep.service.Name
}

func (dep *deleteDependency) GetDatastoreName() string {
	return dep.datastore.GetName()
}
