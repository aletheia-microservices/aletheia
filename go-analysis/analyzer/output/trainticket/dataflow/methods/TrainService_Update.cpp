[0] (PointerObject PointerType) ts (*train.TrainServiceImpl struct{db NoSQLDatabase})
[_1] (StructObject UserType) train.TrainServiceImpl struct{db NoSQLDatabase}
[__2] (FieldObject FieldType) db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context

[0] (StructObject UserType) ttype train.TrainType struct{ID string, Name string, EconomyClass int64, ComfortClass int64, AvgSpeed int64}
[_1] (FieldObject FieldType) ID string
[__2] (BasicObject BasicType) ID string

[0] (BlueprintBackendObject BlueprintBackendType) coll NoSQLCollection {database = train, collection = train}

[0] (InterfaceObject UserType) err .error

[0] (SliceObject UserType) query primitive.D
[_1] (StructObject StructType) struct{ string,  "id" string,  string}

