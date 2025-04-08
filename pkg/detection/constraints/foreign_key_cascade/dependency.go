package foreign_key_cascade

import (
	"fmt"
	"strings"

	"analyzer/pkg/datastores"
	"analyzer/pkg/service"
)

type pendingDelete struct {
	datastore  *datastores.Datastore
	constraint *datastores.Constraint
	services   []*service.Service
	cascading  bool
}

func newPendingDelete(datastore *datastores.Datastore, constraint *datastores.Constraint, services []*service.Service) *pendingDelete {
	return &pendingDelete{
		datastore:  datastore,
		constraint: constraint,
		services:   services,
	}
}

func (pd *pendingDelete) setCascading(v bool) {
	pd.cascading = v
}

func (pd *pendingDelete) hasCascadingAction() bool {
	return pd.cascading
}

func (pd *pendingDelete) isOnDatastore(datastore *datastores.Datastore) bool {
	return pd.datastore == datastore
}

func (pd *pendingDelete) isOnConstraint(constraint *datastores.Constraint) bool {
	return pd.constraint == constraint
}

func (pd *pendingDelete) String() string {
	var svcs []string
	for _, service := range pd.services {
		svcs = append(svcs, service.GetName())
	}
	return fmt.Sprintf("(%s, %s) on constraint %s", pd.datastore.GetName(), strings.Join(svcs, ", "), pd.constraint.String())
}

func (pd *pendingDelete) LongString() string {
	return pd.String()
}

func (pd *pendingDelete) GetDatastoreName() string {
	return pd.datastore.GetName()
}
