[0] (PointerObject PointerType) s (*foo.FooServiceImpl struct{fooDb NoSQLDatabase})
[_1] (StructObject UserType) foo.FooServiceImpl struct{fooDb NoSQLDatabase}
[__2] (FieldObject FieldType) fooDb NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) fooDb NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ FrontendService

    --> w-tainted: write(foo_db.Foo.Text, bar_db.Bar.Text) {2}
[0] (BasicObject BasicType) text string
     --> w-tainted: write(foo_db.Foo.Text, bar_db.Bar.Text) {2}
[_1] (Reference BasicType) ref <"Frontend" string> @ FrontendService

    --> w-tainted: write(foo_db.Foo) {1}
[0] (StructObject UserType) foo foo.Foo struct{ID "id" string, Text string}
     --> w-tainted: write(foo_db.Foo.ID) {1}
[_1] (FieldObject FieldType) ID "id" string
      --> w-tainted: write(foo_db.Foo.ID) {1}
[__2] (BasicObject BasicType) "id" string
     --> w-tainted: write(foo_db.Foo.Text, bar_db.Bar.Text) {2}
[_1] (FieldObject FieldType) Text string
      --> w-tainted: write(foo_db.Foo.Text, bar_db.Bar.Text) {2}
[__2] (BasicObject BasicType) text string
       --> w-tainted: write(foo_db.Foo.Text, bar_db.Bar.Text) {2}
[___3] (Reference BasicType) ref <"Frontend" string> @ FrontendService

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = foo, collection = foo}

[0] (InterfaceObject UserType) err .error

[0] (InterfaceObject UserType) err .error

