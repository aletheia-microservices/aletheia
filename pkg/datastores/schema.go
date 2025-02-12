package datastores

import (
	"fmt"
	"slices"
	"strings"

	"analyzer/pkg/logger"
)

const UNKNOWN_FIELD_TYPE = "<unknown type>"

const ROOT_FIELD_NAME_NOSQL = "_"
const ROOT_FIELD_NAME_QUEUE = "_"

const ROOT_FIELD_NAME_CACHE_KEY = "key"
const ROOT_FIELD_NAME_CACHE_VALUE = "value"

type Schema struct {
	Fields            []*Field            `json:"fields"`
	UnfoldedFields    []*Field            `json:"unfolded_fields"`
	Constraints       []*Constraint       `json:"constraints"`
}

func (s *Schema) AddField(field *Field) {
	s.Fields = append(s.Fields, field)
}

func NewEntry(name string, t string, id int64, datastore *Datastore) *Field {
	return &Field{
		Name:      name,
		Type:      t,
		Id:        id,
		Datastore: datastore,
	}
}

func (s *Schema) GetConstraintsUniqueForFieldName(fieldName string) []*Constraint {
	var constraints []*Constraint
	for _, constraintUnique := range s.GetConstraintsUnique() {
		for _, field := range constraintUnique.GetFields() {
			if field.GetName() == fieldName {
				constraints = append(constraints, constraintUnique)
			}
		}
	}
	return constraints
}

func (s *Schema) GetConstraintsUnique() []*Constraint {
	var constraints []*Constraint
	for _, c := range s.Constraints {
		if c.unique {
			constraints = append(constraints, c)
		}
	}
	return constraints
}

func (s *Schema) HasConstraintsUnique() bool {
	for _, c := range s.Constraints {
		if c.unique {
			return true
		}
	}
	return false
}

func (s *Schema) GetAllFields() []*Field {
	var fields []*Field
	for _, field := range s.Fields {
		if !slices.Contains(fields, field) {
			fields = append(fields, field)
		}
	}
	for _, unfoldedField := range s.UnfoldedFields {
		if !slices.Contains(fields, unfoldedField) {
			fields = append(fields, unfoldedField)
		}
	}
	return fields
}

func (s *Schema) GetRootFieldName() string {
	//logger.Logger.Infof("SCHEMA: %v", s.Fields)
	// FIXME: better to have an additional bool for the fields that state if they are root or no, but for now we have:
	// index 0 is for "_" root field (that can be created in reads if no fields exists) and index 1 is when there was a previous write
	if len(s.Fields) >= 2 {
		return s.Fields[1].GetName()
	}
	if len(s.Fields) == 1 {
		return s.Fields[0].GetName()
	}
	// should never happen because prior to this is only called for query objects and prior to this
	// the cursor must have already created a new root field if it did not exist yet
	logger.Logger.Fatalf("[SCHEMA] no fields for schema: %s", s.String())
	return ""
}

func (s *Schema) GetOrCreateField(name string, t string, id int64, datastore *Datastore) *Field {
	for _, field := range s.Fields {
		if field.GetName() == name && field.GetDatastoreName() == datastore.GetName() { // last condition of datastore is just for sanity check
			// upgrade type if type is unknown
			if field.HasUnknownType() && t != UNKNOWN_FIELD_TYPE {
				field.SetType(t)
			}
			return field
		}
	}

	e := NewEntry(name, t, id, datastore)
	s.Fields = append(s.Fields, e)
	return e
}

func (s *Schema) GetOrCreateUnfoldedField(name string, t string, id int64, datastore *Datastore) *Field {
	for _, field := range s.UnfoldedFields {
		if field.IsNamed(name) && field.GetDatastoreName() == datastore.GetName() { // last condition of datastore is just for sanity check
			// upgrade type if type is unknown
			if field.HasUnknownType() && t != UNKNOWN_FIELD_TYPE {
				field.SetType(t)
			}
			return field
		}
	}

	e := NewEntry(name, t, id, datastore)
	s.UnfoldedFields = append(s.UnfoldedFields, e)
	return e
}

func (s *Schema) GetOrCreateRootField(name string, t string, id int64, datastore *Datastore) *Field {
	// FIXME: better to have an additional bool for the fields that state if they are root or no, but for now we have:
	// index 0 is for "_" root field (that can be created in reads if no fields exists) and index 1 is when there was a previous write
	if len(s.Fields) >= 2 {
		return s.Fields[1]
	}
	if len(s.Fields) == 1 {
		return s.Fields[0]
	}
	e := NewEntry(name, t, id, datastore)
	s.Fields = append(s.Fields, e)
	return e
}

func (s *Schema) AddForeignReferenceToField(current *Field, reference *Field) {
	if !slices.Contains(current.References, reference) {
		current.References = append(current.References, reference)
	}
}

func (s *Schema) String() string {
	fieldsStr := "fields = {"
	for i, field := range s.Fields {
		fieldsStr += fmt.Sprintf("[field #%d] %s", i, field.String())
		if i < len(s.Fields)-1 {
			fieldsStr += ", "
		}
	}
	fieldsStr += "}"

	unfoldedFieldsStr := "unfolded fields = {"
	for i, field := range s.UnfoldedFields {
		unfoldedFieldsStr += fmt.Sprintf("[unfolded field #%d] %s", i, field.String())
		if i < len(s.Fields)-1 {
			unfoldedFieldsStr += ", "
		}
	}
	unfoldedFieldsStr += "}"

	return fieldsStr + ", " + unfoldedFieldsStr
}

func (s *Schema) GetRootUnfoldedField() *Field {
	if len(s.UnfoldedFields) > 1 {
		return s.UnfoldedFields[0]
	}
	logger.Logger.Fatalf("[SCHEMA] no root unfolded field for schema: %s", s.String())
	return nil
}

func (s *Schema) GetFieldByFullName(str string) *Field {
	for _, unfoldedField := range s.UnfoldedFields {
		if unfoldedField.GetFullName() == str {
			return unfoldedField
		}
	}
	logger.Logger.Fatalf("[SCHEMA] no unfolded field (%s) for schema: %s", str, s.String())
	return nil
}

func (s *Schema) GetField(name string) *Field {
	for _, field := range s.Fields {
		if field.IsNamed(name) {
			return field
		}
	}
	for _, field := range s.UnfoldedFields {
		if field.IsNamed(name) {
			return field
		}
	}
	logger.Logger.Fatalf("[FIXME] no field for name %s in datastore schema %s", name, s.String())
	return nil
}

func (s *Schema) GetFieldById(id int64) *Field {
	for _, field := range s.Fields {
		if field.HasId(id) {
			return field
		}
	}
	for _, field := range s.UnfoldedFields {
		if field.HasId(id) {
			return field
		}
	}
	logger.Logger.Warnf("[FIXME] no field for id %d in datastore schema %v", id, s.String())
	return nil
}

type Constraint struct {
	// default constraint 	-> fields size = 1
	// composed constraint 	-> fields size > 1
	fields  []*Field
	unique  bool
	primary bool
	foreign bool
}

func (constraint *Constraint) AddField(field *Field) {
	constraint.fields = append(constraint.fields, field)
}

func (constraint *Constraint) GetFields() []*Field {
	return constraint.fields
}

func (constraint *Constraint) String() string {
	var fieldsStr string
	/* if len(constraint.fields) > 1 {
		fieldsStr += "("
	} */
	fieldsStr += "("
	for i, field := range constraint.fields {
		fieldsStr += field.GetName()
		if i < len(constraint.fields)-1 {
			fieldsStr += ", "
		}
	}
	/* if len(constraint.fields) > 1 {
		fieldsStr += ")"
	} */

	fieldsStr += ")"
	if constraint.unique {
		return "UNIQUE" + fieldsStr
	} else if constraint.primary {
		return "PRIMARY KEY" + fieldsStr
	} else if constraint.foreign {
		return "FOREIGN KEY" + fieldsStr + "REFERENCES" + " [TODO]"
	}
	return ""
}

func NewConstraintUnique(field ...*Field) *Constraint {
	return &Constraint{
		fields: field,
		unique: true,
	}
}

func NewConstraintPrimary(field ...*Field) *Constraint {
	return &Constraint{
		fields:  field,
		primary: true,
	}
}

func NewConstraintForeign(field ...*Field) *Constraint {
	return &Constraint{
		fields:  field,
		foreign: true,
	}
}

func (s *Schema) AddConstraint(constraint *Constraint) {
	s.Constraints = append(s.Constraints, constraint)
}

func (s *Schema) GetConstraints() []*Constraint {
	return s.Constraints
}

type Field struct {
	Name              string
	Type              string
	Datastore         *Datastore
	References        []*Field
	MandatoryRefs     []*Field // aka Total Participation
	Id                int64
	constraints       []*Constraint
}

func (field *Field) GetName() string {
	return field.Name
}
func (field *Field) IsNamed(other string) bool {
	return strings.EqualFold(field.GetName(), other) // FIXME NOSQL MONGODB
}

func (field *Field) GetDatastoreName() string {
	return field.Datastore.GetName()
}
func (field *Field) GetDatastore() *Datastore {
	return field.Datastore
}
func (field *Field) GetFullName() string {
	return strings.ToUpper(field.Datastore.GetName()) + "." + field.Name
}
func (field *Field) GetType() string {
	return field.Type
}
func (field *Field) String() string {
	return field.Name + " " + field.Type
}
func (field *Field) HasId(id int64) bool {
	return field.Id == id
}
func (field *Field) HasUnknownType() bool {
	return field.Type == UNKNOWN_FIELD_TYPE
}
func (field *Field) SetType(t string) {
	field.Type = t
}
func (field *Field) AddMandatoryReference(ref *Field) {
	field.MandatoryRefs = append(field.MandatoryRefs, ref)
}
func (field *Field) GetMandatoryReferences() []*Field {
	return field.MandatoryRefs
}
func (field *Field) AddConstraint(constraint *Constraint) {
	field.constraints = append(field.constraints, constraint)
}
