[0] (PointerObject PointerType) t (*trainfood.TrainFoodServiceImpl struct{db NoSQLDatabase})
[_1] (StructObject UserType) trainfood.TrainFoodServiceImpl struct{db NoSQLDatabase}
[__2] (FieldObject FieldType) db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context

    --> w-tainted: write(trainfood_db.TrainFood) {1}
[0] (StructObject UserType) tf trainfood.TrainFood struct{ID string, TripID string, Foods []food.Food struct{Name string, Price float64}}
     --> w-tainted: write(trainfood_db.TrainFood.TripID) {1}
[_1] (FieldObject FieldType) TripID string
      --> w-tainted: write(trainfood_db.TrainFood.TripID) {1}           --> w-tainted: write(trainfood_db.TrainFood.TripID) {1} --> r-tainted: read(trainfood_db._.tripid) {1}
[__2] (BasicObject BasicType) TripID string

[0] (BlueprintBackendObject BlueprintBackendType) coll NoSQLCollection {database = trainfood, collection = trainfood}

[0] (InterfaceObject UserType) err .error

[0] (SliceObject UserType) query primitive.D
[_1] (StructObject StructType) struct{ string,  "tripid" string,  string}

    --> r-tainted: read(trainfood_db._) {1}
[0] (BlueprintBackendObject BlueprintBackendType) res NoSQLCursor {database = trainfood, collection = trainfood}
     --> r-tainted: read(trainfood_db.TrainFood) {1}
[_1] (StructObject UserType) stored_tf trainfood.TrainFood struct{ID string, TripID string, Foods []food.Food struct{Name string, Price float64}}

[0] (InterfaceObject UserType) err .error

    --> r-tainted: read(trainfood_db.TrainFood) {1}
[0] (StructObject UserType) stored_tf trainfood.TrainFood struct{ID string, TripID string, Foods []food.Food struct{Name string, Price float64}}

[0] (BasicObject BasicType) exists bool

[0] (InterfaceObject UserType) err .error

[0] (BasicObject BasicType) ok bool

[0] (InterfaceObject UserType) err .error

