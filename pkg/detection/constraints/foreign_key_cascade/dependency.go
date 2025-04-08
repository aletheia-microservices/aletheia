package foreign_key_cascade

import (
	"fmt"
	"strings"

	"analyzer/pkg/datastores"
	"analyzer/pkg/service"
)

type deleteDependency struct {
	services  []*service.Service
	datastore *datastores.Datastore
	cascading bool
}

func newDeleteDependency(datastore *datastores.Datastore, services []*service.Service) *deleteDependency {
	return &deleteDependency{
		datastore: datastore,
		services:  services,
	}
}

func (dep *deleteDependency) setCascading(v bool) {
	dep.cascading = v
}

func (dep *deleteDependency) hasCascadingAction() bool {
	return dep.cascading
}

func (dep *deleteDependency) hasDatastore(datastore *datastores.Datastore) bool {
	return dep.datastore == datastore
}

func (dep *deleteDependency) String() string {
	var svcs []string
	for _, service := range dep.services {
		svcs = append(svcs, service.GetName())
	}
	return fmt.Sprintf("(%s, %s)", dep.datastore.GetName(), strings.Join(svcs, ", "))
}

func (dep *deleteDependency) LongString() string {
	return dep.String()
}

func (dep *deleteDependency) GetDatastoreName() string {
	return dep.datastore.GetName()
}
