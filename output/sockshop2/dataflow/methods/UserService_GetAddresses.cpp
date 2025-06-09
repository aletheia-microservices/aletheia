[0] (PointerObject PointerType) s (*user.userServiceImpl struct{UserService user.UserService, db NoSQLDatabase})
[_1] (StructObject UserType) user.userServiceImpl struct{UserService user.UserService, db NoSQLDatabase}
[__2] (FieldObject FieldType) db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ OrderService
[__2] (Reference UserType) ref <ctx context.Context> @ FrontendService
[_1] (Reference UserType) ref <ctx context.Context> @ FrontendService

    --> r-tainted: read(user_db._.addresses) {1}
[0] (BasicObject BasicType) addressid string
     --> r-tainted: read(user_db._.addresses) {1}
[_1] (Reference BasicType) ref <addressID string> @ OrderService
      --> r-tainted: read(user_db._.addresses) {1}
[__2] (Reference BasicType) ref <addressID string> @ FrontendService
     --> r-tainted: read(user_db._.addresses) {1}
[_1] (Reference BasicType) ref <addressID string> @ FrontendService

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = users, collection = users}

[0] (InterfaceObject UserType) _ .error

    --> r-tainted: read(user_db.[]user.Address struct{Street string, Number string, Country string, City string, PostCode string, ID string}) {1}
[0] (ArrayObject ArrayType) addresses []user.Address struct{Street string, Number string, Country string, City string, PostCode string, ID string}

[0] (SliceObject UserType) filter primitive.D
[_1] (StructObject StructType) struct{Key "addresses" string, Key "addresses" string, Value string, Value string}
[__2] (FieldObject FieldType) Key "addresses" string
[___3] (BasicObject BasicType) "addresses" string
[__2] (FieldObject FieldType) Value string
       --> r-tainted: read(user_db._.addresses) {1}
[___3] (BasicObject BasicType) addressid string
        --> r-tainted: read(user_db._.addresses) {1}
[____4] (Reference BasicType) ref <addressID string> @ OrderService
         --> r-tainted: read(user_db._.addresses) {1}
[_____5] (Reference BasicType) ref <addressID string> @ FrontendService
        --> r-tainted: read(user_db._.addresses) {1}
[____4] (Reference BasicType) ref <addressID string> @ FrontendService

[0] (SliceObject UserType) projection primitive.D
[_1] (StructObject StructType) struct{Key "addresses" string, Key "addresses" string, Value true bool, Value true bool}
[__2] (FieldObject FieldType) Key "addresses" string
[___3] (BasicObject BasicType) "addresses" string
[__2] (FieldObject FieldType) Value true bool
[___3] (BasicObject BasicType) true bool

    --> r-tainted: read(user_db._) {1}
[0] (BlueprintBackendObject BlueprintBackendType) result NoSQLCursor {database = users, collection = users}
     --> r-tainted: read(user_db.[]user.Address struct{Street string, Number string, Country string, City string, PostCode string, ID string}) {1}
[_1] (ArrayObject ArrayType) addresses []user.Address struct{Street string, Number string, Country string, City string, PostCode string, ID string}

[0] (InterfaceObject UserType) _ .error

