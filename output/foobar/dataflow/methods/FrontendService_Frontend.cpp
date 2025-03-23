[0] (PointerObject PointerType) d (*foobar.FrontendServiceImpl struct{barService bar.BarService, fooService foo.FooService})
[_1] (StructObject UserType) foobar.FrontendServiceImpl struct{barService bar.BarService, fooService foo.FooService}
[__2] (FieldObject FieldType) barService bar.BarService
[___3] (ServiceObject ServiceType) barService bar.BarService
[__2] (FieldObject FieldType) fooService foo.FooService
[___3] (ServiceObject ServiceType) fooService foo.FooService

[0] (InterfaceObject UserType) ctx context.Context

    --> w-tainted: write(foo_db.Foo) {1}
[0] (StructObject UserType) foo foo.Foo struct{ID string, Text string}
     --> w-tainted: write(foo_db.Foo) {1}
[_1] (Reference UserType) ref <foo foo.Foo struct{ID "id" string, Text string}> @ FooService
      --> w-tainted: write(foo_db.Foo.ID) {1}
[__2] (FieldObject FieldType) ID "id" string
       --> w-tainted: write(foo_db.Foo.ID) {1}
[___3] (BasicObject BasicType) "id" string
      --> w-tainted: write(foo_db.Foo.Text, bar_db.Bar.Text) {2}
[__2] (FieldObject FieldType) Text string
       --> w-tainted: write(foo_db.Foo.Text, bar_db.Bar.Text) {2}
[___3] (BasicObject BasicType) text string
        --> w-tainted: write(foo_db.Foo.Text, bar_db.Bar.Text) {2}
[____4] (Reference BasicType) ref <"Frontend" string> @ FrontendService
     --> w-tainted: write(foo_db.Foo.Text) {1}
[_1] (FieldObject FieldType) Text string
      --> w-tainted: write(foo_db.Foo.Text, bar_db.Bar.Text) {2}
[__2] (Reference FieldType) ref <Text string> @ FooService
       --> w-tainted: write(foo_db.Foo.Text, bar_db.Bar.Text) {2}
[___3] (BasicObject BasicType) text string
        --> w-tainted: write(foo_db.Foo.Text, bar_db.Bar.Text) {2}
[____4] (Reference BasicType) ref <"Frontend" string> @ FrontendService
      --> w-tainted: write(foo_db.Foo.Text, bar_db.Bar.Text) {2}
[__2] (BasicObject BasicType) Text string
       --> w-tainted: write(foo_db.Foo.Text, bar_db.Bar.Text) {2}
[___3] (Reference FieldType) ref <Text string> @ FooService
        --> w-tainted: write(foo_db.Foo.Text, bar_db.Bar.Text) {2}
[____4] (BasicObject BasicType) text string
         --> w-tainted: write(foo_db.Foo.Text, bar_db.Bar.Text) {2}
[_____5] (Reference BasicType) ref <"Frontend" string> @ FrontendService

[0] (InterfaceObject UserType) err1 .error
[_1] (Reference BasicType) ref <nil> @ FooService

    --> w-tainted: write(bar_db.Bar) {1}
[0] (StructObject UserType) bar bar.Bar struct{ID string, Text string, Flag bool}
     --> w-tainted: write(bar_db.Bar) {1}
[_1] (Reference UserType) ref <bar bar.Bar struct{ID "id" string, Text string, Flag bool}> @ BarService
      --> w-tainted: write(bar_db.Bar.ID) {1}
[__2] (FieldObject FieldType) ID "id" string
       --> w-tainted: write(bar_db.Bar.ID) {1}
[___3] (BasicObject BasicType) "id" string
      --> w-tainted: write(bar_db.Bar.Text) {1}
[__2] (FieldObject FieldType) Text string
       --> w-tainted: write(bar_db.Bar.Text) {1}
[___3] (BasicObject BasicType) newText string
        --> w-tainted: write(foo_db.Foo.Text, bar_db.Bar.Text) {2}
[____4] (BasicObject BasicType) text string
         --> w-tainted: write(foo_db.Foo.Text, bar_db.Bar.Text) {2}
[_____5] (Reference BasicType) ref <Text string> @ FrontendService
          --> w-tainted: write(foo_db.Foo.Text, bar_db.Bar.Text) {2}
[______6] (Reference FieldType) ref <Text string> @ FooService
           --> w-tainted: write(foo_db.Foo.Text, bar_db.Bar.Text) {2}
[_______7] (BasicObject BasicType) text string
            --> w-tainted: write(foo_db.Foo.Text, bar_db.Bar.Text) {2}
[________8] (Reference BasicType) ref <"Frontend" string> @ FrontendService
     --> w-tainted: write(bar_db.Bar.Text) {1}
[_1] (FieldObject FieldType) Text string
      --> w-tainted: write(bar_db.Bar.Text) {1}
[__2] (Reference FieldType) ref <Text string> @ BarService
       --> w-tainted: write(bar_db.Bar.Text) {1}
[___3] (BasicObject BasicType) newText string
        --> w-tainted: write(foo_db.Foo.Text, bar_db.Bar.Text) {2}
[____4] (BasicObject BasicType) text string
         --> w-tainted: write(foo_db.Foo.Text, bar_db.Bar.Text) {2}
[_____5] (Reference BasicType) ref <Text string> @ FrontendService
          --> w-tainted: write(foo_db.Foo.Text, bar_db.Bar.Text) {2}
[______6] (Reference FieldType) ref <Text string> @ FooService
           --> w-tainted: write(foo_db.Foo.Text, bar_db.Bar.Text) {2}
[_______7] (BasicObject BasicType) text string
            --> w-tainted: write(foo_db.Foo.Text, bar_db.Bar.Text) {2}
[________8] (Reference BasicType) ref <"Frontend" string> @ FrontendService
      --> w-tainted: write(bar_db.Bar.Text) {1}
[__2] (BasicObject BasicType) Text string
       --> w-tainted: write(bar_db.Bar.Text) {1}
[___3] (Reference FieldType) ref <Text string> @ BarService
        --> w-tainted: write(bar_db.Bar.Text) {1}
[____4] (BasicObject BasicType) newText string
         --> w-tainted: write(foo_db.Foo.Text, bar_db.Bar.Text) {2}
[_____5] (BasicObject BasicType) text string
          --> w-tainted: write(foo_db.Foo.Text, bar_db.Bar.Text) {2}
[______6] (Reference BasicType) ref <Text string> @ FrontendService
           --> w-tainted: write(foo_db.Foo.Text, bar_db.Bar.Text) {2}
[_______7] (Reference FieldType) ref <Text string> @ FooService
            --> w-tainted: write(foo_db.Foo.Text, bar_db.Bar.Text) {2}
[________8] (BasicObject BasicType) text string
             --> w-tainted: write(foo_db.Foo.Text, bar_db.Bar.Text) {2}
[_________9] (Reference BasicType) ref <"Frontend" string> @ FrontendService

[0] (InterfaceObject UserType) err2 .error
[_1] (Reference BasicType) ref <nil> @ BarService

[0] (BasicObject BasicType) out string
[_1] (BasicObject BasicType) "%s, %s" string
     --> w-tainted: write(foo_db.Foo.Text, bar_db.Bar.Text) {2}
[_1] (BasicObject BasicType) Text string
      --> w-tainted: write(foo_db.Foo.Text, bar_db.Bar.Text) {2}
[__2] (Reference FieldType) ref <Text string> @ FooService
       --> w-tainted: write(foo_db.Foo.Text, bar_db.Bar.Text) {2}
[___3] (BasicObject BasicType) text string
        --> w-tainted: write(foo_db.Foo.Text, bar_db.Bar.Text) {2}
[____4] (Reference BasicType) ref <"Frontend" string> @ FrontendService
     --> w-tainted: write(bar_db.Bar.Text) {1}
[_1] (BasicObject BasicType) Text string
      --> w-tainted: write(bar_db.Bar.Text) {1}
[__2] (Reference FieldType) ref <Text string> @ BarService
       --> w-tainted: write(bar_db.Bar.Text) {1}
[___3] (BasicObject BasicType) newText string
        --> w-tainted: write(foo_db.Foo.Text, bar_db.Bar.Text) {2}
[____4] (BasicObject BasicType) text string
         --> w-tainted: write(foo_db.Foo.Text, bar_db.Bar.Text) {2}
[_____5] (Reference BasicType) ref <Text string> @ FrontendService
          --> w-tainted: write(foo_db.Foo.Text, bar_db.Bar.Text) {2}
[______6] (Reference FieldType) ref <Text string> @ FooService
           --> w-tainted: write(foo_db.Foo.Text, bar_db.Bar.Text) {2}
[_______7] (BasicObject BasicType) text string
            --> w-tainted: write(foo_db.Foo.Text, bar_db.Bar.Text) {2}
[________8] (Reference BasicType) ref <"Frontend" string> @ FrontendService
[_1] (BasicObject BasicType) "%s, %s" string
     --> w-tainted: write(foo_db.Foo.Text, bar_db.Bar.Text) {2}
[_1] (BasicObject BasicType) Text string
      --> w-tainted: write(foo_db.Foo.Text, bar_db.Bar.Text) {2}
[__2] (Reference FieldType) ref <Text string> @ FooService
       --> w-tainted: write(foo_db.Foo.Text, bar_db.Bar.Text) {2}
[___3] (BasicObject BasicType) text string
        --> w-tainted: write(foo_db.Foo.Text, bar_db.Bar.Text) {2}
[____4] (Reference BasicType) ref <"Frontend" string> @ FrontendService
     --> w-tainted: write(bar_db.Bar.Text) {1}
[_1] (BasicObject BasicType) Text string
      --> w-tainted: write(bar_db.Bar.Text) {1}
[__2] (Reference FieldType) ref <Text string> @ BarService
       --> w-tainted: write(bar_db.Bar.Text) {1}
[___3] (BasicObject BasicType) newText string
        --> w-tainted: write(foo_db.Foo.Text, bar_db.Bar.Text) {2}
[____4] (BasicObject BasicType) text string
         --> w-tainted: write(foo_db.Foo.Text, bar_db.Bar.Text) {2}
[_____5] (Reference BasicType) ref <Text string> @ FrontendService
          --> w-tainted: write(foo_db.Foo.Text, bar_db.Bar.Text) {2}
[______6] (Reference FieldType) ref <Text string> @ FooService
           --> w-tainted: write(foo_db.Foo.Text, bar_db.Bar.Text) {2}
[_______7] (BasicObject BasicType) text string
            --> w-tainted: write(foo_db.Foo.Text, bar_db.Bar.Text) {2}
[________8] (Reference BasicType) ref <"Frontend" string> @ FrontendService

