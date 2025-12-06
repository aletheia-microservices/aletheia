[0] (PointerObject PointerType) s (*stationfood.StationFoodServiceImpl struct{db NoSQLDatabase})
[_1] (StructObject UserType) stationfood.StationFoodServiceImpl struct{db NoSQLDatabase}
[__2] (FieldObject FieldType) db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context

    --> w-tainted: write(stationfood_db.StationFoodStore) {1}
[0] (StructObject UserType) store stationfood.StationFoodStore struct{ID string, StationName string, StoreName string, Telephone string, BusinessTime string, DeliveryFee float64, Foods []food.Food struct{Name string, Price float64}}
     --> w-tainted: write(stationfood_db.StationFoodStore.ID) {1}
[_1] (FieldObject FieldType) ID string
      --> w-tainted: write(stationfood_db.StationFoodStore.ID) {1}           --> w-tainted: write(stationfood_db.StationFoodStore.ID) {1} --> r-tainted: read(stationfood_db._.id) {1}
[__2] (BasicObject BasicType) ID string

[0] (BlueprintBackendObject BlueprintBackendType) coll NoSQLCollection {database = stationfood, collection = stationfood}

[0] (InterfaceObject UserType) err .error

[0] (SliceObject UserType) query primitive.D
[_1] (StructObject StructType) struct{ string,  "id" string,  string}

    --> r-tainted: read(stationfood_db._) {1}
[0] (BlueprintBackendObject BlueprintBackendType) res NoSQLCursor {database = stationfood, collection = stationfood}
     --> r-tainted: read(stationfood_db.StationFoodStore) {1}
[_1] (StructObject UserType) st stationfood.StationFoodStore struct{ID string, StationName string, StoreName string, Telephone string, BusinessTime string, DeliveryFee float64, Foods []food.Food struct{Name string, Price float64}}

[0] (InterfaceObject UserType) err .error

    --> r-tainted: read(stationfood_db.StationFoodStore) {1}
[0] (StructObject UserType) st stationfood.StationFoodStore struct{ID string, StationName string, StoreName string, Telephone string, BusinessTime string, DeliveryFee float64, Foods []food.Food struct{Name string, Price float64}}

[0] (BasicObject BasicType) exists bool

[0] (InterfaceObject UserType) err .error

