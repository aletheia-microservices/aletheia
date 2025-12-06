[0] (PointerObject PointerType) a (*assurance.AssuranceServiceImpl struct{db NoSQLDatabase})
[_1] (StructObject UserType) assurance.AssuranceServiceImpl struct{db NoSQLDatabase}
[__2] (FieldObject FieldType) db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context

[0] (BlueprintBackendObject BlueprintBackendType) coll NoSQLCollection {database = assurance, collection = assurance}

[0] (InterfaceObject UserType) err .error

    --> r-tainted: read(assurance_db._) {1}
[0] (BlueprintBackendObject BlueprintBackendType) res NoSQLCursor {database = assurance, collection = assurance}
     --> r-tainted: read(assurance_db.[]assurance.Assurance struct{ID string, OrderID string, AT assurance.AssuranceType struct{Index int64, Name string, Price float64}}) {1}
[_1] (ArrayObject ArrayType) assurances []assurance.Assurance struct{ID string, OrderID string, AT assurance.AssuranceType struct{Index int64, Name string, Price float64}}

[0] (InterfaceObject UserType) err .error

    --> r-tainted: read(assurance_db.[]assurance.Assurance struct{ID string, OrderID string, AT assurance.AssuranceType struct{Index int64, Name string, Price float64}}) {1}
[0] (ArrayObject ArrayType) assurances []assurance.Assurance struct{ID string, OrderID string, AT assurance.AssuranceType struct{Index int64, Name string, Price float64}}

[0] (InterfaceObject UserType) err .error

