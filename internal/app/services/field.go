package services

import (
	"fmt"
)

type Field struct {
	idx        int    // index in service struct
	fieldName  string // name in service struct
	wiringName string // identifier passed in blueprint wiring (if field is DB then it is the db name)
}

func NewField(idx int, name string) *Field {
	return &Field{
		idx:       idx,
		fieldName: name,
	}
}

func (field *Field) String() string {
	str := fmt.Sprintf("%s #%d", field.fieldName, field.idx)
	if field.wiringName != "" {
		str += " (" + field.wiringName + ")"
	}
	return str
}

func (field *Field) SetWiringName(id string) {
	field.wiringName = id
}

func (field *Field) GetWiringName() string {
	return field.wiringName
}

func (field *Field) GetFieldName() string {
	return field.fieldName
}
