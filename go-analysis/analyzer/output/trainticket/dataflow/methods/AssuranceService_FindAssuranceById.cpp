[0] (PointerObject PointerType) a (*assurance.AssuranceServiceImpl struct{db NoSQLDatabase})
[_1] (StructObject UserType) assurance.AssuranceServiceImpl struct{db NoSQLDatabase}
[__2] (FieldObject FieldType) db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ AssuranceService

    --> r-tainted: read(assurance_db._.id) {1}
[0] (BasicObject BasicType) id string
     --> r-tainted: read(assurance_db._.id) {1}
[_1] (Reference BasicType) ref <ID string> @ AssuranceService

[0] (BlueprintBackendObject BlueprintBackendType) coll NoSQLCollection {database = assurance, collection = assurance}

[0] (InterfaceObject UserType) err .error

[0] (SliceObject UserType) query primitive.D
[_1] (StructObject StructType) struct{ string,  "id" string,  string}

    --> r-tainted: read(assurance_db._) {1}
[0] (BlueprintBackendObject BlueprintBackendType) res NoSQLCursor {database = assurance, collection = assurance}
     --> r-tainted: read(assurance_db.Assurance) {1}
[_1] (StructObject UserType) ass assurance.Assurance struct{ID string, OrderID string, AT assurance.AssuranceType struct{Index int64, Name string, Price float64}}

[0] (InterfaceObject UserType) err .error

    --> r-tainted: read(assurance_db.Assurance) {1}
[0] (StructObject UserType) ass assurance.Assurance struct{ID string, OrderID string, AT assurance.AssuranceType struct{Index int64, Name string, Price float64}}

[0] (BasicObject BasicType) exists bool

[0] (InterfaceObject UserType) err .error

