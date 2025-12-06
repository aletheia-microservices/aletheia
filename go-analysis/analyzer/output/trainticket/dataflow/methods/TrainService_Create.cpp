[0] (PointerObject PointerType) ts (*train.TrainServiceImpl struct{db NoSQLDatabase})
[_1] (StructObject UserType) train.TrainServiceImpl struct{db NoSQLDatabase}
[__2] (FieldObject FieldType) db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context

    --> w-tainted: write(train_db.TrainType) {1}
[0] (StructObject UserType) tt train.TrainType struct{ID string, Name string, EconomyClass int64, ComfortClass int64, AvgSpeed int64}
     --> w-tainted: write(train_db.TrainType.Name) {1}
[_1] (FieldObject FieldType) Name string
      --> w-tainted: write(train_db.TrainType.Name) {1}           --> w-tainted: write(train_db.TrainType.Name) {1} --> r-tainted: read(train_db._.name) {1}
[__2] (BasicObject BasicType) Name string

[0] (BlueprintBackendObject BlueprintBackendType) coll NoSQLCollection {database = train, collection = train}

[0] (InterfaceObject UserType) err .error

[0] (SliceObject UserType) query primitive.D
[_1] (StructObject StructType) struct{ string,  "name" string,  string}

    --> r-tainted: read(train_db._) {1}
[0] (BlueprintBackendObject BlueprintBackendType) res NoSQLCursor {database = train, collection = train}
     --> r-tainted: read(train_db.TrainType) {1}
[_1] (StructObject UserType) saved_tt train.TrainType struct{ID string, Name string, EconomyClass int64, ComfortClass int64, AvgSpeed int64}

[0] (InterfaceObject UserType) err .error

    --> r-tainted: read(train_db.TrainType) {1}
[0] (StructObject UserType) saved_tt train.TrainType struct{ID string, Name string, EconomyClass int64, ComfortClass int64, AvgSpeed int64}

[0] (BasicObject BasicType) exists bool

[0] (InterfaceObject UserType) err .error

[0] (InterfaceObject UserType) err .error

