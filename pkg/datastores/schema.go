package datastores

import (
	"fmt"
	"slices"
	"strings"

	"analyzer/pkg/logger"
	"analyzer/pkg/utils"
)

const UNKNOWN_FIELD_TYPE = "<unknown type>"

const ROOT_FIELD_NAME_NOSQL = "_"
const ROOT_FIELD_NAME_QUEUE = "_"
const ROOT_FIELD_NAME_RELATIONALDB = "*"

const ROOT_FIELD_NAME_CACHE_KEY = "key"
const ROOT_FIELD_NAME_CACHE_VALUE = "value"

type Schema struct {
	Fields         []*Field      `json:"fields"`
	UnfoldedFields []*Field      `json:"unfolded_fields"`
	Constraints    []*Constraint `json:"constraints"`
}

func (s *Schema) AddField(field *Field) {
	s.Fields = append(s.Fields, field)
}

func NewField(name string, t string, id int64, datastore *Datastore) *Field {
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

func (s *Schema) GetConstraintsNumericalForFieldName(fieldName string) []*Constraint {
	var constraints []*Constraint
	for _, constraintUnique := range s.GetConstraintsNumerical() {
		for _, field := range constraintUnique.GetFields() {
			if field.GetName() == fieldName {
				constraints = append(constraints, constraintUnique)
			}
		}
	}
	return constraints
}

// includes any PK
func (s *Schema) GetConstraintsUnique() []*Constraint {
	var constraints []*Constraint
	for _, c := range s.Constraints {
		if c.unique {
			constraints = append(constraints, c)
		} else if c.primary {
			constraints = append(constraints, c)
		}
	}
	return constraints
}

func (s *Schema) GetConstraintsNumerical() []*Constraint {
	var constraints []*Constraint
	for _, c := range s.Constraints {
		if c.numerical != nil {
			constraints = append(constraints, c)
		}
	}
	return constraints
}

// includes any PK
func (s *Schema) HasConstraintsUnique() bool {
	for _, c := range s.Constraints {
		if c.unique {
			return true
		} else if c.primary {
			return true
		}
	}
	return false
}

func (s *Schema) HasConstraintsNumerical() bool {
	for _, c := range s.Constraints {
		if c.numerical != nil {
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
		logger.Logger.Debugf("[SCHEMA] getting root field (%s) at index 1 for fields list: %v", s.Fields[1].GetName(), s.Fields)
		return s.Fields[1].GetName()
	}
	if len(s.Fields) == 1 {
		logger.Logger.Debugf("[SCHEMA] getting root field (%s) at index 0 for fields list: %v", s.Fields[0].GetName(), s.Fields)
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

	e := NewField(name, t, id, datastore)
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

	e := NewField(name, t, id, datastore)
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
	e := NewField(name, t, id, datastore)
	s.Fields = append(s.Fields, e)
	return e
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
	if field := s.GetFieldIfExists(name); field != nil {
		return field
	}
	logger.Logger.Fatalf("[FIXME] no field for name %s in datastore schema %s", name, s.String())
	return nil
}

func (s *Schema) GetFieldIfExists(name string) *Field {
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
	fields    []*Field
	unique    bool
	primary   bool
	reference bool
	mandatory bool
	numerical *NumericalConstraint
}

type ComparisonOperator int

const (
	EQ ComparisonOperator = iota // ==
	LT // <
	GT // >
	LE // <=
	GE // >=
	NE // !=
)

type NumericalConstraint struct {
	value    string
	operator ComparisonOperator
}

func NewNumericalConstraint(value string, operator ComparisonOperator) *NumericalConstraint {
	return &NumericalConstraint{value: value, operator: operator}
}

func ConstraintOperatorToString(operator ComparisonOperator) string {
	switch operator {
	case EQ:
		return "=="
	case LT:
		return "<"
	case LE:
		return "<="
	case GT:
		return ">"
	case GE:
		return ">="
	case NE:
		return "!="
	}
	return ""
}

func ConstraintValueToString(value int) string {
	return fmt.Sprintf("%d", value)
}

func (constraint *Constraint) AddField(field *Field) {
	constraint.fields = append(constraint.fields, field)
}

func (constraint *Constraint) GetFields() []*Field {
	return constraint.fields
}

func (constraint *Constraint) String() string {
	var fieldsStr, delimiter string

	if constraint.reference {
		delimiter = " --> "
	} else {
		delimiter = ", "
	}
	/* if len(constraint.fields) > 1 {
		fieldsStr += "("
	} */
	fieldsStr += "("
	for i, field := range constraint.fields {
		fieldsStr += field.GetName()
		if i < len(constraint.fields)-1 {
			fieldsStr += delimiter
		}
	}
	/* if len(constraint.fields) > 1 {
		fieldsStr += ")"
	} */

	fieldsStr += ")"
	if constraint.unique {
		return "UNIQUE" + fieldsStr
	} else if constraint.primary {
		return "PRIMARY" + fieldsStr
	} else if constraint.reference && constraint.mandatory {
		return "REFERENCE" + fieldsStr + " [+]"
	} else if constraint.reference {
		return "REFERENCE" + fieldsStr
	} else if constraint.mandatory {
		return "MANDATORY" + fieldsStr
	} else if constraint.numerical != nil {
		return "CHECK(" + constraint.fields[0].GetName() + " " + ConstraintOperatorToString(constraint.numerical.operator) + " " + constraint.numerical.value + ")"
	}
	return ""
}

func NewConstraintNumerical(numerical *NumericalConstraint, fields ...*Field) *Constraint {
	return &Constraint{
		fields:    fields,
		numerical: numerical,
	}
}

func NewConstraintUnique(fields ...*Field) *Constraint {
	return &Constraint{
		fields: fields,
		unique: true,
	}
}

func NewConstraintPrimary(fields ...*Field) *Constraint {
	return &Constraint{
		fields:  fields,
		primary: true,
	}
}

func NewConstraintReference(mandatory bool, fields ...*Field) *Constraint {
	return &Constraint{
		fields:    fields,
		reference: true,
		mandatory: mandatory,
	}
}
func (c *Constraint) FieldIsReferencing(field *Field) bool {
	if len(c.fields) != 2 {
		logger.Logger.Fatalf("[CONSTRAINT] unexpected length (%d) for constraint REFERENCE", len(c.fields))
	}
	return c.fields[0] == field
}
func (c *Constraint) GetReferencingField() *Field {
	if len(c.fields) != 2 {
		logger.Logger.Fatalf("[CONSTRAINT] unexpected length (%d) for constraint REFERENCE", len(c.fields))
	}
	return c.fields[0]
}
func (c *Constraint) FieldIsReferencedBy(field *Field) bool {
	if len(c.fields) != 2 {
		logger.Logger.Fatalf("[CONSTRAINT] unexpected length (%d) for constraint REFERENCE", len(c.fields))
	}
	return c.fields[1] == field
}
func (c *Constraint) GetReferencedByField() *Field {
	if len(c.fields) != 2 {
		logger.Logger.Fatalf("[CONSTRAINT] unexpected length (%d) for constraint REFERENCE", len(c.fields))
	}
	return c.fields[1]
}

func (s *Schema) AddConstraint(constraint *Constraint) {
	s.Constraints = append(s.Constraints, constraint)
}

func (s *Schema) GetConstraints(filter ConstraintFilter) []*Constraint {
	var constraints []*Constraint
	for _, c := range s.Constraints {
		if (filter.Unique == nil || *filter.Unique == c.unique) &&
			(filter.Primary == nil || *filter.Primary == c.primary) &&
			(filter.Reference == nil || *filter.Reference == c.reference) &&
			(filter.Mandatory == nil || *filter.Mandatory == c.mandatory) {
			constraints = append(constraints, c)
		}
	}
	return constraints
}
func (s *Schema) GetAllConstraints() []*Constraint {
	return s.Constraints
}

type Field struct {
	Name        string
	Type        string
	Datastore   *Datastore
	References  []*Field
	Id          int64
	constraints []*Constraint
}

func (field *Field) GetName() string {
	return field.Name
}

type ConstraintFilter struct {
	Unique    *bool
	Primary   *bool
	Reference *bool
	Mandatory *bool
}

func (field *Field) GetConstraints(filter ConstraintFilter) []*Constraint {
	var constraints []*Constraint
	for _, c := range field.constraints {
		if (filter.Unique == nil || *filter.Unique == c.unique) &&
			(filter.Primary == nil || *filter.Primary == c.primary) &&
			(filter.Reference == nil || *filter.Reference == c.reference) &&
			(filter.Mandatory == nil || *filter.Mandatory == c.mandatory) {
			constraints = append(constraints, c)
		}
	}
	return constraints
}

func (field *Field) GetAllConstraints() []*Constraint {
	return field.constraints
}
func (field *Field) GetReferences() []*Field {
	return field.References
}
func (field *Field) HasReference(ref *Field) bool {
	return slices.Contains(field.References, ref)
}
func (field *Field) CreateAndAddReference(ref *Field, mandatory bool) *Constraint {
	field.References = append(field.References, ref)
	constraint := NewConstraintReference(mandatory, field, ref)
	field.AddConstraint(constraint)
	return constraint
}

func (field *Field) CompactReferences(refsToKeep []*Field) {
	field.References = refsToKeep
	var constraintsToKeep []*Constraint
	for _, ref := range refsToKeep {
		for _, constraint := range field.GetConstraints(ConstraintFilter{Reference: utils.BoolPtr(true)}) {
			if constraint.FieldIsReferencedBy(ref) {
				constraintsToKeep = append(constraintsToKeep, constraint)
			}
		}
	}
	field.constraints = constraintsToKeep
}
func (field *Field) HasReferences() bool {
	return len(field.References) > 0
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
func (field *Field) AddConstraint(constraint *Constraint) {
	field.constraints = append(field.constraints, constraint)
}
