[0] (PointerObject PointerType) s (*digota.OrderServiceImpl struct{skuService digota.SkuService, db NoSQLDatabase})
[_1] (StructObject UserType) digota.OrderServiceImpl struct{skuService digota.SkuService, db NoSQLDatabase}
[__2] (FieldObject FieldType) db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase
[__2] (FieldObject FieldType) skuService digota.SkuService
[___3] (ServiceObject ServiceType) skuService digota.SkuService

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ OrderService

    --> w-tainted: write(orders_db.Order.Items) {1}
[0] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
     --> w-tainted: write(orders_db.Order.Items) {1}
[_1] (Reference ArrayType) ref <items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ OrderService

[0] (MapObject MapType) skuMap map[string](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})

    --> w-tainted: write(orders_db.Order.Items) {1}
[0] (ArrayObject ArrayType) orderItems [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
     --> w-tainted: write(orders_db.Order.Items) {1}
[_1] (PointerObject PointerType) v (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (Reference ArrayType) ref <items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ OrderService
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (FieldObject FieldType) Parent string
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (BasicObject BasicType) Parent string
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (FieldObject FieldType) Quantity int64
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (BasicObject BasicType) 1 int
     --> w-tainted: write(orders_db.Order.Items) {1}
[_1] (PointerObject PointerType) v (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (Reference ArrayType) ref <items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ OrderService
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (FieldObject FieldType) Parent string
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (BasicObject BasicType) Parent string
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (FieldObject FieldType) Quantity int64
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (BasicObject BasicType) 1 int
     --> w-tainted: write(orders_db.Order.Items) {1}
[_1] (PointerObject PointerType) v (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (Reference ArrayType) ref <items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ OrderService
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (FieldObject FieldType) Parent string
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (BasicObject BasicType) Parent string
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (FieldObject FieldType) Quantity int64
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (BasicObject BasicType) 1 int
     --> w-tainted: write(orders_db.Order.Items) {1}
[_1] (PointerObject PointerType) v (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (Reference ArrayType) ref <items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ OrderService
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (FieldObject FieldType) Parent string
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (BasicObject BasicType) Parent string
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (FieldObject FieldType) Quantity int64
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (BasicObject BasicType) 1 int

[0] (ArrayObject ArrayType) errs []error

[0] (StructObject UserType) wg sync.WaitGroup

    --> w-tainted: write(orders_db.Order.Items) {1}
[0] (PointerObject PointerType) v (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
     --> w-tainted: write(orders_db.Order.Items) {1}
[_1] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (Reference ArrayType) ref <items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ OrderService
     --> w-tainted: write(orders_db.Order.Items) {1}
[_1] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (FieldObject FieldType) Parent string
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (BasicObject BasicType) Parent string
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (FieldObject FieldType) Quantity int64
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (BasicObject BasicType) 1 int

[0] (PointerObject PointerType) skuItem (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
[_1] (PointerObject PointerType) * (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
[__2] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
[___3] (FieldObject FieldType) Quantity int64
[____4] (BasicObject BasicType) Quantity int64
[_1] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
[__2] (FieldObject FieldType) Quantity int64
[___3] (BasicObject BasicType) Quantity int64

[0] (ArrayObject ArrayType) orderItems [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
     --> w-tainted: write(orders_db.Order.Items) {1}
[_1] (PointerObject PointerType) v (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (Reference ArrayType) ref <items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ OrderService
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (FieldObject FieldType) Parent string
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (BasicObject BasicType) Parent string
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (FieldObject FieldType) Quantity int64
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (BasicObject BasicType) 1 int

[0] (PointerObject PointerType) s (*digota.OrderServiceImpl struct{skuService digota.SkuService, db NoSQLDatabase})
[_1] (StructObject UserType) digota.OrderServiceImpl struct{skuService digota.SkuService, db NoSQLDatabase}
[__2] (FieldObject FieldType) db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase
[__2] (FieldObject FieldType) skuService digota.SkuService
[___3] (ServiceObject ServiceType) skuService digota.SkuService

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ OrderService

    --> w-tainted: write(orders_db.Order.Items) {1}
[0] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
     --> w-tainted: write(orders_db.Order.Items) {1}
[_1] (Reference ArrayType) ref <items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ OrderService

[0] (MapObject MapType) skuMap map[string](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})

    --> w-tainted: write(orders_db.Order.Items) {1}
[0] (ArrayObject ArrayType) orderItems [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
     --> w-tainted: write(orders_db.Order.Items) {1}
[_1] (PointerObject PointerType) v (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (Reference ArrayType) ref <items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ OrderService
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (FieldObject FieldType) Parent string
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (BasicObject BasicType) Parent string
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (FieldObject FieldType) Quantity int64
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (BasicObject BasicType) 1 int
     --> w-tainted: write(orders_db.Order.Items) {1}
[_1] (PointerObject PointerType) v (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (Reference ArrayType) ref <items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ OrderService
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (FieldObject FieldType) Parent string
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (BasicObject BasicType) Parent string
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (FieldObject FieldType) Quantity int64
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (BasicObject BasicType) 1 int
     --> w-tainted: write(orders_db.Order.Items) {1}
[_1] (PointerObject PointerType) v (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (Reference ArrayType) ref <items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ OrderService
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (FieldObject FieldType) Parent string
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (BasicObject BasicType) Parent string
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (FieldObject FieldType) Quantity int64
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (BasicObject BasicType) 1 int
     --> w-tainted: write(orders_db.Order.Items) {1}
[_1] (PointerObject PointerType) v (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (Reference ArrayType) ref <items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ OrderService
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (FieldObject FieldType) Parent string
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (BasicObject BasicType) Parent string
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (FieldObject FieldType) Quantity int64
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (BasicObject BasicType) 1 int

[0] (ArrayObject ArrayType) errs []error

[0] (StructObject UserType) wg sync.WaitGroup

    --> w-tainted: write(orders_db.Order.Items) {1}
[0] (PointerObject PointerType) v (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
     --> w-tainted: write(orders_db.Order.Items) {1}
[_1] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (Reference ArrayType) ref <items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ OrderService
     --> w-tainted: write(orders_db.Order.Items) {1}
[_1] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (FieldObject FieldType) Parent string
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (BasicObject BasicType) Parent string
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (FieldObject FieldType) Quantity int64
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (BasicObject BasicType) 1 int

[0] (PointerObject PointerType) skuItem (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
[_1] (PointerObject PointerType) * (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
[__2] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
[___3] (FieldObject FieldType) Quantity int64
[____4] (BasicObject BasicType) Quantity int64
[_1] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
[__2] (FieldObject FieldType) Quantity int64
[___3] (BasicObject BasicType) Quantity int64

[0] (PointerObject PointerType) s (*digota.OrderServiceImpl struct{skuService digota.SkuService, db NoSQLDatabase})
[_1] (StructObject UserType) digota.OrderServiceImpl struct{skuService digota.SkuService, db NoSQLDatabase}
[__2] (FieldObject FieldType) db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase
[__2] (FieldObject FieldType) skuService digota.SkuService
[___3] (ServiceObject ServiceType) skuService digota.SkuService

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ OrderService

    --> w-tainted: write(orders_db.Order.Items) {1}
[0] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
     --> w-tainted: write(orders_db.Order.Items) {1}
[_1] (Reference ArrayType) ref <items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ OrderService

[0] (MapObject MapType) skuMap map[string](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})

    --> w-tainted: write(orders_db.Order.Items) {1}
[0] (ArrayObject ArrayType) orderItems [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
     --> w-tainted: write(orders_db.Order.Items) {1}
[_1] (PointerObject PointerType) v (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (Reference ArrayType) ref <items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ OrderService
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (FieldObject FieldType) Parent string
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (BasicObject BasicType) Parent string
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (FieldObject FieldType) Quantity int64
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (BasicObject BasicType) 1 int
     --> w-tainted: write(orders_db.Order.Items) {1}
[_1] (PointerObject PointerType) v (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (Reference ArrayType) ref <items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ OrderService
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (FieldObject FieldType) Parent string
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (BasicObject BasicType) Parent string
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (FieldObject FieldType) Quantity int64
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (BasicObject BasicType) 1 int
     --> w-tainted: write(orders_db.Order.Items) {1}
[_1] (PointerObject PointerType) v (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (Reference ArrayType) ref <items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ OrderService
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (FieldObject FieldType) Parent string
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (BasicObject BasicType) Parent string
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (FieldObject FieldType) Quantity int64
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (BasicObject BasicType) 1 int
     --> w-tainted: write(orders_db.Order.Items) {1}
[_1] (PointerObject PointerType) v (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (Reference ArrayType) ref <items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ OrderService
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (FieldObject FieldType) Parent string
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (BasicObject BasicType) Parent string
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (FieldObject FieldType) Quantity int64
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (BasicObject BasicType) 1 int

[0] (ArrayObject ArrayType) errs []error

[0] (StructObject UserType) wg sync.WaitGroup

    --> w-tainted: write(orders_db.Order.Items) {1}
[0] (PointerObject PointerType) v (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
     --> w-tainted: write(orders_db.Order.Items) {1}
[_1] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (Reference ArrayType) ref <items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ OrderService
     --> w-tainted: write(orders_db.Order.Items) {1}
[_1] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (FieldObject FieldType) Parent string
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (BasicObject BasicType) Parent string
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (FieldObject FieldType) Quantity int64
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (BasicObject BasicType) 1 int

[0] (PointerObject PointerType) v (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
     --> w-tainted: write(orders_db.Order.Items) {1}
[_1] (ArrayObject ArrayType) orderItems [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (PointerObject PointerType) v (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (Reference ArrayType) ref <items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ OrderService
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (FieldObject FieldType) Parent string
         --> w-tainted: write(orders_db.Order.Items) {1}
[_____5] (BasicObject BasicType) Parent string
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (FieldObject FieldType) Quantity int64
         --> w-tainted: write(orders_db.Order.Items) {1}
[_____5] (BasicObject BasicType) 1 int
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (PointerObject PointerType) v (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (Reference ArrayType) ref <items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ OrderService
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (FieldObject FieldType) Parent string
         --> w-tainted: write(orders_db.Order.Items) {1}
[_____5] (BasicObject BasicType) Parent string
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (FieldObject FieldType) Quantity int64
         --> w-tainted: write(orders_db.Order.Items) {1}
[_____5] (BasicObject BasicType) 1 int
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (PointerObject PointerType) v (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (Reference ArrayType) ref <items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ OrderService
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (FieldObject FieldType) Parent string
         --> w-tainted: write(orders_db.Order.Items) {1}
[_____5] (BasicObject BasicType) Parent string
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (FieldObject FieldType) Quantity int64
         --> w-tainted: write(orders_db.Order.Items) {1}
[_____5] (BasicObject BasicType) 1 int
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (PointerObject PointerType) v (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (Reference ArrayType) ref <items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ OrderService
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (FieldObject FieldType) Parent string
         --> w-tainted: write(orders_db.Order.Items) {1}
[_____5] (BasicObject BasicType) Parent string
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (FieldObject FieldType) Quantity int64
         --> w-tainted: write(orders_db.Order.Items) {1}
[_____5] (BasicObject BasicType) 1 int
[_1] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
[__2] (FieldObject FieldType) Amount int64
[___3] (BasicObject BasicType) int64
[____4] (BasicObject BasicType) Price uint64
[__2] (FieldObject FieldType) Currency int32
[___3] (BasicObject BasicType) Currency int32
[__2] (FieldObject FieldType) Description string
[___3] (BasicObject BasicType) github.com/blueprint-uservices/blueprint/examples/digota/workflow/digota.defaultTaxDescription "Tax" untyped string
[__2] (FieldObject FieldType) Parent string
       --> r-tainted: read(skus_db._.id) {1}
[___3] (BasicObject BasicType) Parent string

[0] (PointerObject PointerType) item (*digota.Sku struct{Id string, Name string, Price uint64, Currency int32, Active bool, Parent string, Metadata map[string]string, Attributes map[string]string, Image string, PackageDimensions (*digota.PackageDimensions struct{Height float64, Length float64, Weight float64, Width float64}), Inventory (*digota.Inventory struct{Quantity int64, Type int32}), Created int64, Updated int64})
[_1] (StructObject UserType) digota.Sku struct{Id string, Name string, Price uint64, Currency int32, Active bool, Parent string, Metadata map[string]string, Attributes map[string]string, Image string, PackageDimensions (*digota.PackageDimensions struct{Height float64, Length float64, Weight float64, Width float64}), Inventory (*digota.Inventory struct{Quantity int64, Type int32}), Created int64, Updated int64}
      --> r-tainted: read(skus_db.Sku) {1}
[__2] (Reference UserType) ref <digota.Sku struct{Id string, Name string, Price uint64, Currency int32, Active bool, Parent string, Metadata map[string]string, Attributes map[string]string, Image string, PackageDimensions (*digota.PackageDimensions struct{Height float64, Length float64, Weight float64, Width float64}), Inventory (*digota.Inventory struct{Quantity int64, Type int32}), Created int64, Updated int64}> @ SkuService
[__2] (FieldObject FieldType) Currency int32
[___3] (BasicObject BasicType) Currency int32
[__2] (FieldObject FieldType) Name string
[___3] (BasicObject BasicType) Name string
[__2] (FieldObject FieldType) Price uint64
[___3] (BasicObject BasicType) Price uint64

[0] (InterfaceObject UserType) err .error
[_1] (Reference BasicType) ref <nil> @ SkuService

[0] (PointerObject PointerType) s (*digota.OrderServiceImpl struct{skuService digota.SkuService, db NoSQLDatabase})
[_1] (StructObject UserType) digota.OrderServiceImpl struct{skuService digota.SkuService, db NoSQLDatabase}
[__2] (FieldObject FieldType) db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase
[__2] (FieldObject FieldType) skuService digota.SkuService
[___3] (ServiceObject ServiceType) skuService digota.SkuService

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ OrderService

    --> w-tainted: write(orders_db.Order.Items) {1}
[0] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
     --> w-tainted: write(orders_db.Order.Items) {1}
[_1] (Reference ArrayType) ref <items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ OrderService

[0] (MapObject MapType) skuMap map[string](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})

    --> w-tainted: write(orders_db.Order.Items) {1}
[0] (ArrayObject ArrayType) orderItems [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
     --> w-tainted: write(orders_db.Order.Items) {1}
[_1] (PointerObject PointerType) v (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (Reference ArrayType) ref <items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ OrderService
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (FieldObject FieldType) Parent string
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (BasicObject BasicType) Parent string
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (FieldObject FieldType) Quantity int64
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (BasicObject BasicType) 1 int
     --> w-tainted: write(orders_db.Order.Items) {1}
[_1] (PointerObject PointerType) v (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (Reference ArrayType) ref <items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ OrderService
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (FieldObject FieldType) Parent string
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (BasicObject BasicType) Parent string
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (FieldObject FieldType) Quantity int64
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (BasicObject BasicType) 1 int
     --> w-tainted: write(orders_db.Order.Items) {1}
[_1] (PointerObject PointerType) v (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (Reference ArrayType) ref <items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ OrderService
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (FieldObject FieldType) Parent string
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (BasicObject BasicType) Parent string
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (FieldObject FieldType) Quantity int64
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (BasicObject BasicType) 1 int
     --> w-tainted: write(orders_db.Order.Items) {1}
[_1] (PointerObject PointerType) v (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (Reference ArrayType) ref <items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ OrderService
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (FieldObject FieldType) Parent string
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (BasicObject BasicType) Parent string
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (FieldObject FieldType) Quantity int64
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (BasicObject BasicType) 1 int

[0] (ArrayObject ArrayType) errs []error

[0] (StructObject UserType) wg sync.WaitGroup

    --> w-tainted: write(orders_db.Order.Items) {1}
[0] (PointerObject PointerType) v (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
     --> w-tainted: write(orders_db.Order.Items) {1}
[_1] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (Reference ArrayType) ref <items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ OrderService
     --> w-tainted: write(orders_db.Order.Items) {1}
[_1] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (FieldObject FieldType) Parent string
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (BasicObject BasicType) Parent string
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (FieldObject FieldType) Quantity int64
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (BasicObject BasicType) 1 int

[0] (PointerObject PointerType) skuItem (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
[_1] (PointerObject PointerType) * (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
[__2] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
[___3] (FieldObject FieldType) Quantity int64
[____4] (BasicObject BasicType) Quantity int64
[_1] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
[__2] (FieldObject FieldType) Quantity int64
[___3] (BasicObject BasicType) Quantity int64

[0] (ArrayObject ArrayType) orderItems [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
     --> w-tainted: write(orders_db.Order.Items) {1}
[_1] (PointerObject PointerType) v (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (Reference ArrayType) ref <items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ OrderService
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (FieldObject FieldType) Parent string
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (BasicObject BasicType) Parent string
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (FieldObject FieldType) Quantity int64
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (BasicObject BasicType) 1 int

[0] (PointerObject PointerType) s (*digota.OrderServiceImpl struct{skuService digota.SkuService, db NoSQLDatabase})
[_1] (StructObject UserType) digota.OrderServiceImpl struct{skuService digota.SkuService, db NoSQLDatabase}
[__2] (FieldObject FieldType) db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase
[__2] (FieldObject FieldType) skuService digota.SkuService
[___3] (ServiceObject ServiceType) skuService digota.SkuService

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ OrderService

    --> w-tainted: write(orders_db.Order.Items) {1}
[0] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
     --> w-tainted: write(orders_db.Order.Items) {1}
[_1] (Reference ArrayType) ref <items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ OrderService

[0] (MapObject MapType) skuMap map[string](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})

    --> w-tainted: write(orders_db.Order.Items) {1}
[0] (ArrayObject ArrayType) orderItems [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
     --> w-tainted: write(orders_db.Order.Items) {1}
[_1] (PointerObject PointerType) v (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (Reference ArrayType) ref <items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ OrderService
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (FieldObject FieldType) Parent string
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (BasicObject BasicType) Parent string
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (FieldObject FieldType) Quantity int64
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (BasicObject BasicType) 1 int
     --> w-tainted: write(orders_db.Order.Items) {1}
[_1] (PointerObject PointerType) v (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (Reference ArrayType) ref <items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ OrderService
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (FieldObject FieldType) Parent string
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (BasicObject BasicType) Parent string
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (FieldObject FieldType) Quantity int64
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (BasicObject BasicType) 1 int
     --> w-tainted: write(orders_db.Order.Items) {1}
[_1] (PointerObject PointerType) v (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (Reference ArrayType) ref <items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ OrderService
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (FieldObject FieldType) Parent string
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (BasicObject BasicType) Parent string
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (FieldObject FieldType) Quantity int64
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (BasicObject BasicType) 1 int
     --> w-tainted: write(orders_db.Order.Items) {1}
[_1] (PointerObject PointerType) v (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (Reference ArrayType) ref <items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ OrderService
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (FieldObject FieldType) Parent string
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (BasicObject BasicType) Parent string
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (FieldObject FieldType) Quantity int64
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (BasicObject BasicType) 1 int

[0] (ArrayObject ArrayType) errs []error

[0] (StructObject UserType) wg sync.WaitGroup

    --> w-tainted: write(orders_db.Order.Items) {1}
[0] (PointerObject PointerType) v (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
     --> w-tainted: write(orders_db.Order.Items) {1}
[_1] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (Reference ArrayType) ref <items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ OrderService
     --> w-tainted: write(orders_db.Order.Items) {1}
[_1] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (FieldObject FieldType) Parent string
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (BasicObject BasicType) Parent string
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (FieldObject FieldType) Quantity int64
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (BasicObject BasicType) 1 int

[0] (PointerObject PointerType) skuItem (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
[_1] (PointerObject PointerType) * (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
[__2] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
[___3] (FieldObject FieldType) Quantity int64
[____4] (BasicObject BasicType) Quantity int64
[_1] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
[__2] (FieldObject FieldType) Quantity int64
[___3] (BasicObject BasicType) Quantity int64

[0] (PointerObject PointerType) s (*digota.OrderServiceImpl struct{skuService digota.SkuService, db NoSQLDatabase})
[_1] (StructObject UserType) digota.OrderServiceImpl struct{skuService digota.SkuService, db NoSQLDatabase}
[__2] (FieldObject FieldType) db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase
[__2] (FieldObject FieldType) skuService digota.SkuService
[___3] (ServiceObject ServiceType) skuService digota.SkuService

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ OrderService

    --> w-tainted: write(orders_db.Order.Items) {1}
[0] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
     --> w-tainted: write(orders_db.Order.Items) {1}
[_1] (Reference ArrayType) ref <items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ OrderService

[0] (MapObject MapType) skuMap map[string](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})

    --> w-tainted: write(orders_db.Order.Items) {1}
[0] (ArrayObject ArrayType) orderItems [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
     --> w-tainted: write(orders_db.Order.Items) {1}
[_1] (PointerObject PointerType) v (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (Reference ArrayType) ref <items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ OrderService
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (FieldObject FieldType) Parent string
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (BasicObject BasicType) Parent string
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (FieldObject FieldType) Quantity int64
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (BasicObject BasicType) 1 int
     --> w-tainted: write(orders_db.Order.Items) {1}
[_1] (PointerObject PointerType) v (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (Reference ArrayType) ref <items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ OrderService
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (FieldObject FieldType) Parent string
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (BasicObject BasicType) Parent string
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (FieldObject FieldType) Quantity int64
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (BasicObject BasicType) 1 int
     --> w-tainted: write(orders_db.Order.Items) {1}
[_1] (PointerObject PointerType) v (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (Reference ArrayType) ref <items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ OrderService
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (FieldObject FieldType) Parent string
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (BasicObject BasicType) Parent string
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (FieldObject FieldType) Quantity int64
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (BasicObject BasicType) 1 int
     --> w-tainted: write(orders_db.Order.Items) {1}
[_1] (PointerObject PointerType) v (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (Reference ArrayType) ref <items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ OrderService
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (FieldObject FieldType) Parent string
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (BasicObject BasicType) Parent string
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (FieldObject FieldType) Quantity int64
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (BasicObject BasicType) 1 int

[0] (ArrayObject ArrayType) errs []error

[0] (StructObject UserType) wg sync.WaitGroup

    --> w-tainted: write(orders_db.Order.Items) {1}
[0] (PointerObject PointerType) v (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
     --> w-tainted: write(orders_db.Order.Items) {1}
[_1] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (Reference ArrayType) ref <items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ OrderService
     --> w-tainted: write(orders_db.Order.Items) {1}
[_1] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (FieldObject FieldType) Parent string
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (BasicObject BasicType) Parent string
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (FieldObject FieldType) Quantity int64
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (BasicObject BasicType) 1 int

[0] (PointerObject PointerType) v (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
     --> w-tainted: write(orders_db.Order.Items) {1}
[_1] (ArrayObject ArrayType) orderItems [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (PointerObject PointerType) v (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (Reference ArrayType) ref <items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ OrderService
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (FieldObject FieldType) Parent string
         --> w-tainted: write(orders_db.Order.Items) {1}
[_____5] (BasicObject BasicType) Parent string
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (FieldObject FieldType) Quantity int64
         --> w-tainted: write(orders_db.Order.Items) {1}
[_____5] (BasicObject BasicType) 1 int
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (PointerObject PointerType) v (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (Reference ArrayType) ref <items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ OrderService
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (FieldObject FieldType) Parent string
         --> w-tainted: write(orders_db.Order.Items) {1}
[_____5] (BasicObject BasicType) Parent string
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (FieldObject FieldType) Quantity int64
         --> w-tainted: write(orders_db.Order.Items) {1}
[_____5] (BasicObject BasicType) 1 int
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (PointerObject PointerType) v (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (Reference ArrayType) ref <items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ OrderService
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (FieldObject FieldType) Parent string
         --> w-tainted: write(orders_db.Order.Items) {1}
[_____5] (BasicObject BasicType) Parent string
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (FieldObject FieldType) Quantity int64
         --> w-tainted: write(orders_db.Order.Items) {1}
[_____5] (BasicObject BasicType) 1 int
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (PointerObject PointerType) v (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (Reference ArrayType) ref <items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ OrderService
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (FieldObject FieldType) Parent string
         --> w-tainted: write(orders_db.Order.Items) {1}
[_____5] (BasicObject BasicType) Parent string
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (FieldObject FieldType) Quantity int64
         --> w-tainted: write(orders_db.Order.Items) {1}
[_____5] (BasicObject BasicType) 1 int
[_1] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
[__2] (FieldObject FieldType) Amount int64
[___3] (BasicObject BasicType) int64
[____4] (BasicObject BasicType) Price uint64
[__2] (FieldObject FieldType) Currency int32
[___3] (BasicObject BasicType) Currency int32
[__2] (FieldObject FieldType) Description string
[___3] (BasicObject BasicType) github.com/blueprint-uservices/blueprint/examples/digota/workflow/digota.defaultTaxDescription "Tax" untyped string
[__2] (FieldObject FieldType) Parent string
       --> r-tainted: read(skus_db._.id) {1}
[___3] (BasicObject BasicType) Parent string

