[0] (PointerObject PointerType) a (*assurance.AssuranceServiceImpl struct{db NoSQLDatabase})
[_1] (StructObject UserType) assurance.AssuranceServiceImpl struct{db NoSQLDatabase}
[__2] (FieldObject FieldType) db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context

[0] (BasicObject BasicType) typeindex int64

    --> w-tainted: write(assurance_db.Assurance.OrderID) {1}
[0] (BasicObject BasicType) orderid string

    --> w-tainted: write(assurance_db.Assurance.AT) {1}
[0] (StructObject UserType) at assurance.AssuranceType struct{Index int64, Name string, Price float64}
     --> w-tainted: write(assurance_db.Assurance.AT) {1}
[_1] (Reference UserType) ref <assurance.AssuranceType struct{Index int64, Name string, Price float64}> @ getAssuranceType

[0] (InterfaceObject UserType) err .error
[_1] (Reference BasicType) ref <nil> @ getAssuranceType

    --> w-tainted: write(assurance_db.Assurance.ID) {1}
[0] (BasicObject BasicType) id string

    --> w-tainted: write(assurance_db.Assurance) {1}
[0] (StructObject UserType) assurance assurance.Assurance struct{ID string, OrderID string, AT assurance.AssuranceType struct{Index int64, Name string, Price float64}}
     --> w-tainted: write(assurance_db.Assurance.AT) {1}
[_1] (FieldObject FieldType) AT assurance.AssuranceType struct{Index int64, Name string, Price float64}
      --> w-tainted: write(assurance_db.Assurance.AT) {1}
[__2] (StructObject UserType) at assurance.AssuranceType struct{Index int64, Name string, Price float64}
       --> w-tainted: write(assurance_db.Assurance.AT) {1}
[___3] (Reference UserType) ref <assurance.AssuranceType struct{Index int64, Name string, Price float64}> @ getAssuranceType
     --> w-tainted: write(assurance_db.Assurance.ID) {1}
[_1] (FieldObject FieldType) ID string
      --> w-tainted: write(assurance_db.Assurance.ID) {1}
[__2] (BasicObject BasicType) id string
     --> w-tainted: write(assurance_db.Assurance.OrderID) {1}
[_1] (FieldObject FieldType) OrderID string
      --> w-tainted: write(assurance_db.Assurance.OrderID) {1}
[__2] (BasicObject BasicType) orderid string

[0] (BlueprintBackendObject BlueprintBackendType) coll NoSQLCollection {database = assurance, collection = assurance}

[0] (InterfaceObject UserType) err .error

