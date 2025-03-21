[0] (PointerObject PointerType) s (*bar.BarServiceImpl struct{barDb NoSQLDatabase})
[_1] (StructObject UserType) bar.BarServiceImpl struct{barDb NoSQLDatabase}
[__2] (FieldObject FieldType) barDb NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) barDb NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ FrontendService

    --> w-tainted: write(foo_db.Foo.Text, bar_db.Bar.Text) {2}
[0] (BasicObject BasicType) text string
     --> w-tainted: write(foo_db.Foo.Text, bar_db.Bar.Text) {2}
[_1] (Reference BasicType) ref <Text string> @ FrontendService
      --> w-tainted: write(foo_db.Foo.Text, bar_db.Bar.Text) {2}
[__2] (Reference FieldType) ref <Text string> @ FooService
       --> w-tainted: write(foo_db.Foo.Text, bar_db.Bar.Text) {2}
[___3] (BasicObject BasicType) text string
        --> w-tainted: write(foo_db.Foo.Text, bar_db.Bar.Text) {2}
[____4] (Reference BasicType) ref <"Frontend" string> @ FrontendService

    --> w-tainted: write(bar_db.Bar.Text) {1}
[0] (BasicObject BasicType) newText string
     --> w-tainted: write(foo_db.Foo.Text, bar_db.Bar.Text) {2}
[_1] (BasicObject BasicType) text string
      --> w-tainted: write(foo_db.Foo.Text, bar_db.Bar.Text) {2}
[__2] (Reference BasicType) ref <Text string> @ FrontendService
       --> w-tainted: write(foo_db.Foo.Text, bar_db.Bar.Text) {2}
[___3] (Reference FieldType) ref <Text string> @ FooService
        --> w-tainted: write(foo_db.Foo.Text, bar_db.Bar.Text) {2}
[____4] (BasicObject BasicType) text string
         --> w-tainted: write(foo_db.Foo.Text, bar_db.Bar.Text) {2}
[_____5] (Reference BasicType) ref <"Frontend" string> @ FrontendService

    --> w-tainted: write(bar_db.Bar) {1}
[0] (StructObject UserType) bar bar.Bar struct{ID "id" string, Text string}
     --> w-tainted: write(bar_db.Bar.ID) {1}
[_1] (FieldObject FieldType) ID "id" string
      --> w-tainted: write(bar_db.Bar.ID) {1}
[__2] (BasicObject BasicType) "id" string
     --> w-tainted: write(bar_db.Bar.Text) {1}
[_1] (FieldObject FieldType) Text string
      --> w-tainted: write(bar_db.Bar.Text) {1}
[__2] (BasicObject BasicType) newText string
       --> w-tainted: write(foo_db.Foo.Text, bar_db.Bar.Text) {2}
[___3] (BasicObject BasicType) text string
        --> w-tainted: write(foo_db.Foo.Text, bar_db.Bar.Text) {2}
[____4] (Reference BasicType) ref <Text string> @ FrontendService
         --> w-tainted: write(foo_db.Foo.Text, bar_db.Bar.Text) {2}
[_____5] (Reference FieldType) ref <Text string> @ FooService
          --> w-tainted: write(foo_db.Foo.Text, bar_db.Bar.Text) {2}
[______6] (BasicObject BasicType) text string
           --> w-tainted: write(foo_db.Foo.Text, bar_db.Bar.Text) {2}
[_______7] (Reference BasicType) ref <"Frontend" string> @ FrontendService

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = bar, collection = bar}

[0] (InterfaceObject UserType) err .error

[0] (InterfaceObject UserType) err .error

