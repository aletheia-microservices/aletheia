[0] (PointerObject PointerType) ts (*train.TrainServiceImpl struct{db NoSQLDatabase})
[_1] (StructObject UserType) train.TrainServiceImpl struct{db NoSQLDatabase}
[__2] (FieldObject FieldType) db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context

    --> r-tainted: read(train_db.[]train.TrainType struct{ID string, Name string, EconomyClass int64, ComfortClass int64, AvgSpeed int64}) {1}
[0] (ArrayObject ArrayType) trains []train.TrainType struct{ID string, Name string, EconomyClass int64, ComfortClass int64, AvgSpeed int64}

[0] (BlueprintBackendObject BlueprintBackendType) coll NoSQLCollection {database = train, collection = train}

[0] (InterfaceObject UserType) err .error

    --> r-tainted: read(train_db.TrainType) {1}
[0] (BlueprintBackendObject BlueprintBackendType) res NoSQLCursor {database = train, collection = train}
     --> r-tainted: read(train_db.[]train.TrainType struct{ID string, Name string, EconomyClass int64, ComfortClass int64, AvgSpeed int64}) {1}
[_1] (ArrayObject ArrayType) trains []train.TrainType struct{ID string, Name string, EconomyClass int64, ComfortClass int64, AvgSpeed int64}

[0] (InterfaceObject UserType) err .error

[0] (InterfaceObject UserType) err .error

