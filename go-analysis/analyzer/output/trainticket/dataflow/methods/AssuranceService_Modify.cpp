[0] (PointerObject PointerType) a (*assurance.AssuranceServiceImpl struct{db NoSQLDatabase})
[_1] (StructObject UserType) assurance.AssuranceServiceImpl struct{db NoSQLDatabase}
[__2] (FieldObject FieldType) db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context

[0] (StructObject UserType) assurance assurance.Assurance struct{ID string, OrderID string, AT assurance.AssuranceType struct{Index int64, Name string, Price float64}}
[_1] (FieldObject FieldType) ID string
      --> r-tainted: read(assurance_db._.id) {1}
[__2] (BasicObject BasicType) ID string

[0] (StructObject UserType) _ assurance.Assurance struct{ID string, OrderID string, AT assurance.AssuranceType struct{Index int64, Name string, Price float64}}
     --> r-tainted: read(assurance_db.Assurance) {1}
[_1] (Reference UserType) ref <ass assurance.Assurance struct{ID string, OrderID string, AT assurance.AssuranceType struct{Index int64, Name string, Price float64}}> @ FindAssuranceById

[0] (InterfaceObject UserType) err .error
[_1] (Reference BasicType) ref <nil> @ FindAssuranceById

[0] (BlueprintBackendObject BlueprintBackendType) coll NoSQLCollection {database = assurance, collection = assurance}

[0] (InterfaceObject UserType) err .error

[0] (SliceObject UserType) query primitive.D
[_1] (StructObject StructType) struct{ string,  "id" string,  string}

[0] (BasicObject BasicType) ok bool

[0] (InterfaceObject UserType) err .error

