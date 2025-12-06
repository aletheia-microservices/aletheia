[0] (PointerObject PointerType) t (*trainfood.TrainFoodServiceImpl struct{db NoSQLDatabase})
[_1] (StructObject UserType) trainfood.TrainFoodServiceImpl struct{db NoSQLDatabase}
[__2] (FieldObject FieldType) db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context

    --> r-tainted: read(trainfood_db._.tripid) {1}
[0] (BasicObject BasicType) tripid string

[0] (BlueprintBackendObject BlueprintBackendType) coll NoSQLCollection {database = trainfood, collection = trainfood}

[0] (InterfaceObject UserType) err .error

    --> r-tainted: read(trainfood_db._) {1}
[0] (BlueprintBackendObject BlueprintBackendType) res NoSQLCursor {database = trainfood, collection = trainfood}
     --> r-tainted: read(trainfood_db.TrainFood, trainfood_db.TrainFood.Foods) {2}
[_1] (StructObject UserType) tf trainfood.TrainFood struct{ID string, TripID string, Foods []food.Food struct{Name string, Price float64}}
      --> r-tainted: read(trainfood_db.TrainFood.Foods) {1}
[__2] (FieldObject FieldType) Foods []food.Food struct{Name string, Price float64}
       --> r-tainted: read(trainfood_db.TrainFood.Foods) {1}
[___3] (SliceObject SliceType) Foods []food.Food struct{Name string, Price float64}

[0] (InterfaceObject UserType) err .error

    --> r-tainted: read(trainfood_db.TrainFood, trainfood_db.TrainFood.Foods) {2}
[0] (StructObject UserType) tf trainfood.TrainFood struct{ID string, TripID string, Foods []food.Food struct{Name string, Price float64}}
     --> r-tainted: read(trainfood_db.TrainFood.Foods) {1}
[_1] (FieldObject FieldType) Foods []food.Food struct{Name string, Price float64}
      --> r-tainted: read(trainfood_db.TrainFood.Foods) {1}
[__2] (SliceObject SliceType) Foods []food.Food struct{Name string, Price float64}

[0] (BasicObject BasicType) exists bool

[0] (InterfaceObject UserType) err .error

