package specialization

import (
	"analyzer/pkg/datastores"
)

type mandatoryField struct {
	field        datastores.Field
	mandatoryRef datastores.Field
}

func newMandatoryField(field datastores.Field, mandatoryRef datastores.Field) *mandatoryField {
	return &mandatoryField{
		field:        field,
		mandatoryRef: mandatoryRef,
	}
}
