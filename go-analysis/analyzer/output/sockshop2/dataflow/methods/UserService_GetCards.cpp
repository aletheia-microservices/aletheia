[0] (PointerObject PointerType) s (*user.userServiceImpl struct{UserService user.UserService, db NoSQLDatabase})
[_1] (StructObject UserType) user.userServiceImpl struct{UserService user.UserService, db NoSQLDatabase}
[__2] (FieldObject FieldType) db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ FrontendService
[_1] (Reference UserType) ref <ctx context.Context> @ OrderService
[__2] (Reference UserType) ref <ctx context.Context> @ FrontendService

    --> r-tainted: read(user_db._.cards) {1}
[0] (BasicObject BasicType) cardid string
     --> r-tainted: read(user_db._.cards) {1}
[_1] (Reference BasicType) ref <cardID string> @ FrontendService
     --> r-tainted: read(user_db._.cards) {1}
[_1] (Reference BasicType) ref <cardID string> @ OrderService
      --> r-tainted: read(user_db._.cards) {1}
[__2] (Reference BasicType) ref <cardID string> @ FrontendService

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = users, collection = users}

[0] (InterfaceObject UserType) _ .error

    --> r-tainted: read(user_db.[]user.Card struct{LongNum string, Expires string, CCV string, ID string}) {1}
[0] (ArrayObject ArrayType) cards []user.Card struct{LongNum string, Expires string, CCV string, ID string}

[0] (SliceObject UserType) filter primitive.D
[_1] (StructObject StructType) struct{Key "cards" string, Key "cards" string, Value string, Value string}
[__2] (FieldObject FieldType) Key "cards" string
[___3] (BasicObject BasicType) "cards" string
[__2] (FieldObject FieldType) Value string
       --> r-tainted: read(user_db._.cards) {1}
[___3] (BasicObject BasicType) cardid string
        --> r-tainted: read(user_db._.cards) {1}
[____4] (Reference BasicType) ref <cardID string> @ FrontendService
        --> r-tainted: read(user_db._.cards) {1}
[____4] (Reference BasicType) ref <cardID string> @ OrderService
         --> r-tainted: read(user_db._.cards) {1}
[_____5] (Reference BasicType) ref <cardID string> @ FrontendService

[0] (SliceObject UserType) projection primitive.D
[_1] (StructObject StructType) struct{Key "cards" string, Key "cards" string, Value true bool, Value true bool}
[__2] (FieldObject FieldType) Key "cards" string
[___3] (BasicObject BasicType) "cards" string
[__2] (FieldObject FieldType) Value true bool
[___3] (BasicObject BasicType) true bool

    --> r-tainted: read(user_db._) {1}
[0] (BlueprintBackendObject BlueprintBackendType) result NoSQLCursor {database = users, collection = users}
     --> r-tainted: read(user_db.[]user.Card struct{LongNum string, Expires string, CCV string, ID string}) {1}
[_1] (ArrayObject ArrayType) cards []user.Card struct{LongNum string, Expires string, CCV string, ID string}

[0] (InterfaceObject UserType) _ .error

