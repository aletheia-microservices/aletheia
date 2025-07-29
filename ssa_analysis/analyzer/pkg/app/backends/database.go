package backends

type Database struct {
	name   string
	schema *Schema
}

func NewDatabase(name string, schema *Schema) *Database {
	return &Database{
		name:   name,
		schema: schema,
	}
}

func (database *Database) GetName() string {
	return database.name
}

func (database *Database) GetSchema() *Schema {
	return database.schema
}

func (database *Database) String() string {
	return database.name + " // schema: \n" + database.schema.String()
}
