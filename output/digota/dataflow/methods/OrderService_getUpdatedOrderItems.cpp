[0] (PointerObject PointerType) s (*digota.OrderServiceImpl struct{skuService digota.SkuService, db NoSQLDatabase})
[_1] (StructObject UserType) digota.OrderServiceImpl struct{skuService digota.SkuService, db NoSQLDatabase}
[__2] (FieldObject FieldType) db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase
[__2] (FieldObject FieldType) skuService digota.SkuService
[___3] (ServiceObject ServiceType) skuService digota.SkuService

[0] (InterfaceObject UserType) ctx context.Context

[0] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})

[0] (ArrayObject ArrayType) orderItems [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
[_1] (PointerObject PointerType) myitem1 (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
[__2] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
[__2] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
[___3] (FieldObject FieldType) Quantity int64
[____4] (BasicObject BasicType) 1 int

[0] (ArrayObject ArrayType) errs []error

[0] (StructObject UserType) wg sync.WaitGroup

[0] (PointerObject PointerType) myitem1 (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
[_1] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
[_1] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
[__2] (FieldObject FieldType) Quantity int64
[___3] (BasicObject BasicType) 1 int

[0] (ArrayObject ArrayType) orderItems [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
[_1] (PointerObject PointerType) myitem1 (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
[__2] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
[__2] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
[___3] (FieldObject FieldType) Quantity int64
[____4] (BasicObject BasicType) 1 int

[0] (PointerObject PointerType) s (*digota.OrderServiceImpl struct{skuService digota.SkuService, db NoSQLDatabase})
[_1] (StructObject UserType) digota.OrderServiceImpl struct{skuService digota.SkuService, db NoSQLDatabase}
[__2] (FieldObject FieldType) db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase
[__2] (FieldObject FieldType) skuService digota.SkuService
[___3] (ServiceObject ServiceType) skuService digota.SkuService

[0] (InterfaceObject UserType) ctx context.Context

[0] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})

[0] (ArrayObject ArrayType) orderItems [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
[_1] (PointerObject PointerType) myitem1 (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
[__2] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
[__2] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
[___3] (FieldObject FieldType) Quantity int64
[____4] (BasicObject BasicType) 1 int

[0] (ArrayObject ArrayType) errs []error

[0] (StructObject UserType) wg sync.WaitGroup

[0] (PointerObject PointerType) myitem1 (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
[_1] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
[_1] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
[__2] (FieldObject FieldType) Quantity int64
[___3] (BasicObject BasicType) 1 int

[0] (PointerObject PointerType) myitem2 (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
[_1] (ArrayObject ArrayType) orderItems [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
[__2] (PointerObject PointerType) myitem1 (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
[___3] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
[___3] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
[____4] (FieldObject FieldType) Quantity int64
[_____5] (BasicObject BasicType) 1 int
[_1] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
[__2] (FieldObject FieldType) Amount int64
[___3] (BasicObject BasicType) int64
[____4] (BasicObject BasicType) Price uint64
[__2] (FieldObject FieldType) Currency int32
[___3] (BasicObject BasicType) Currency int32
[__2] (FieldObject FieldType) Description string
[___3] (BasicObject BasicType) Name string
[__2] (FieldObject FieldType) Parent string
[___3] (BasicObject BasicType) Parent string

[0] (PointerObject PointerType) item (*digota.Sku struct{Id string, Name string, Price uint64, Currency int32, Active bool, Parent string, Metadata map[string]string, Attributes map[string]string, Image string, PackageDimensions (*digota.PackageDimensions struct{Height float64, Length float64, Weight float64, Width float64}), Inventory (*digota.Inventory struct{Quantity int64, Type int32}), Created int64, Updated int64})
[_1] (StructObject UserType) digota.Sku struct{Id string, Name string, Price uint64, Currency int32, Active bool, Parent string, Metadata map[string]string, Attributes map[string]string, Image string, PackageDimensions (*digota.PackageDimensions struct{Height float64, Length float64, Weight float64, Width float64}), Inventory (*digota.Inventory struct{Quantity int64, Type int32}), Created int64, Updated int64}
[__2] (FieldObject FieldType) Currency int32
[___3] (BasicObject BasicType) Currency int32
[__2] (FieldObject FieldType) Name string
[___3] (BasicObject BasicType) Name string
[__2] (FieldObject FieldType) Price uint64
[___3] (BasicObject BasicType) Price uint64

[0] (InterfaceObject UserType) err .error

[0] (PointerObject PointerType) s (*digota.OrderServiceImpl struct{skuService digota.SkuService, db NoSQLDatabase})
[_1] (StructObject UserType) digota.OrderServiceImpl struct{skuService digota.SkuService, db NoSQLDatabase}
[__2] (FieldObject FieldType) db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase
[__2] (FieldObject FieldType) skuService digota.SkuService
[___3] (ServiceObject ServiceType) skuService digota.SkuService

[0] (InterfaceObject UserType) ctx context.Context

[0] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})

[0] (ArrayObject ArrayType) orderItems [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
[_1] (PointerObject PointerType) myitem1 (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
[__2] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
[__2] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
[___3] (FieldObject FieldType) Quantity int64
[____4] (BasicObject BasicType) 1 int

[0] (ArrayObject ArrayType) errs []error

[0] (StructObject UserType) wg sync.WaitGroup

[0] (PointerObject PointerType) myitem1 (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
[_1] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
[_1] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
[__2] (FieldObject FieldType) Quantity int64
[___3] (BasicObject BasicType) 1 int

[0] (ArrayObject ArrayType) orderItems [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
[_1] (PointerObject PointerType) myitem1 (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
[__2] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
[__2] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
[___3] (FieldObject FieldType) Quantity int64
[____4] (BasicObject BasicType) 1 int

[0] (PointerObject PointerType) s (*digota.OrderServiceImpl struct{skuService digota.SkuService, db NoSQLDatabase})
[_1] (StructObject UserType) digota.OrderServiceImpl struct{skuService digota.SkuService, db NoSQLDatabase}
[__2] (FieldObject FieldType) db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase
[__2] (FieldObject FieldType) skuService digota.SkuService
[___3] (ServiceObject ServiceType) skuService digota.SkuService

[0] (InterfaceObject UserType) ctx context.Context

[0] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})

[0] (ArrayObject ArrayType) orderItems [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
[_1] (PointerObject PointerType) myitem1 (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
[__2] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
[__2] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
[___3] (FieldObject FieldType) Quantity int64
[____4] (BasicObject BasicType) 1 int

[0] (ArrayObject ArrayType) errs []error

[0] (StructObject UserType) wg sync.WaitGroup

[0] (PointerObject PointerType) myitem1 (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
[_1] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
[_1] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
[__2] (FieldObject FieldType) Quantity int64
[___3] (BasicObject BasicType) 1 int

[0] (PointerObject PointerType) myitem2 (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
[_1] (ArrayObject ArrayType) orderItems [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
[__2] (PointerObject PointerType) myitem1 (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
[___3] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
[___3] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
[____4] (FieldObject FieldType) Quantity int64
[_____5] (BasicObject BasicType) 1 int
[_1] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
[__2] (FieldObject FieldType) Amount int64
[___3] (BasicObject BasicType) int64
[____4] (BasicObject BasicType) Price uint64
[__2] (FieldObject FieldType) Currency int32
[___3] (BasicObject BasicType) Currency int32
[__2] (FieldObject FieldType) Description string
[___3] (BasicObject BasicType) Name string
[__2] (FieldObject FieldType) Parent string
[___3] (BasicObject BasicType) Parent string

