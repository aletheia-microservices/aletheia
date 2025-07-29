package backends

import "log"

type Schema struct {
	fields map[string]*Field
}

func NewSchema() *Schema {
	return &Schema{
		fields: make(map[string]*Field),
	}
}

func (schema *Schema) HasField(fieldname string) bool {
	_, ok := schema.fields[fieldname]
	return ok
}

func (schema *Schema) AddField(field *Field) {
	schema.fields[field.path] = field
}

func (schema *Schema) GetFields() map[string]*Field {
	return schema.fields
}

func (schema *Schema) GetFieldByPath(path string) *Field {
	field, ok := schema.fields[path]
	if !ok {
		log.Fatalf("field with path (%s) not found for schema: %s\n", path, schema.String())
	}
	return field
}

func (schema *Schema) String() string {
	var str string
	for _, field := range schema.fields {
		str += "\t " + field.String() + "\n"
	}
	return str
}
