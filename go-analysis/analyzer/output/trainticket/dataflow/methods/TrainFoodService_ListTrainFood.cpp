[0] (PointerObject PointerType) t (*trainfood.TrainFoodServiceImpl struct{db NoSQLDatabase})
[_1] (StructObject UserType) trainfood.TrainFoodServiceImpl struct{db NoSQLDatabase}
[__2] (FieldObject FieldType) db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context

[0] (BlueprintBackendObject BlueprintBackendType) coll NoSQLCollection {database = trainfood, collection = trainfood}

[0] (InterfaceObject UserType) err .error

    --> r-tainted: read(trainfood_db._) {1}
[0] (BlueprintBackendObject BlueprintBackendType) res NoSQLCursor {database = trainfood, collection = trainfood}
     --> r-tainted: read(trainfood_db.[]trainfood.TrainFood struct{ID string, TripID string, Foods []food.Food struct{Name string, Price float64}}) {1}
[_1] (ArrayObject ArrayType) all_foods []trainfood.TrainFood struct{ID string, TripID string, Foods []food.Food struct{Name string, Price float64}}

[0] (InterfaceObject UserType) err .error

    --> r-tainted: read(trainfood_db.[]trainfood.TrainFood struct{ID string, TripID string, Foods []food.Food struct{Name string, Price float64}}) {1}
[0] (ArrayObject ArrayType) all_foods []trainfood.TrainFood struct{ID string, TripID string, Foods []food.Food struct{Name string, Price float64}}

[0] (InterfaceObject UserType) err .error

