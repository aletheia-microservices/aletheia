[0] (PointerObject PointerType) s (*digota.OrderServiceImpl struct{skuService digota.SkuService, db NoSQLDatabase})
[_1] (StructObject UserType) digota.OrderServiceImpl struct{skuService digota.SkuService, db NoSQLDatabase}
[__2] (FieldObject FieldType) db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase
[__2] (FieldObject FieldType) skuService digota.SkuService
[___3] (ServiceObject ServiceType) skuService digota.SkuService

[0] (InterfaceObject UserType) ctx context.Context

    --> w-tainted: write(orders_db.Order.Currency) {1}
[0] (BasicObject BasicType) currency int32

    --> w-tainted: write(orders_db.Order.Items) {1}
[0] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})

    --> w-tainted: write(orders_db.Order.Metadata) {1}
[0] (MapObject MapType) metadata map[string]string

    --> w-tainted: write(orders_db.Order.Email) {1}
[0] (BasicObject BasicType) email string

    --> w-tainted: write(orders_db.Order.Shipping) {1}
[0] (PointerObject PointerType) shipping (*digota.Shipping struct{Name string, Phone string, Address (*digota.Shipping_Address struct{Line1 string, City string, Country string, Line2 string, PostalCode string, State string}), Carrier string, TrackingNumber string})
     --> w-tainted: write(orders_db.Order.Shipping) {1}
[_1] (StructObject UserType) digota.Shipping struct{Name string, Phone string, Address (*digota.Shipping_Address struct{Line1 string, City string, Country string, Line2 string, PostalCode string, State string}), Carrier string, TrackingNumber string}

[0] (PointerObject PointerType) order (*digota.Order struct{Id string, Amount int64, Currency int32, Items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}), Metadata map[string]string, Email string, ChargeId string, Status int32, Shipping (*digota.Shipping struct{Name string, Phone string, Address (*digota.Shipping_Address struct{Line1 string, City string, Country string, Line2 string, PostalCode string, State string}), Carrier string, TrackingNumber string}), Created int64, Updated int64})
     --> w-tainted: write(orders_db.Order) {1}
[_1] (StructObject UserType) digota.Order struct{Id string, Amount int64, Currency int32, Items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}), Metadata map[string]string, Email string, ChargeId string, Status int32, Shipping (*digota.Shipping struct{Name string, Phone string, Address (*digota.Shipping_Address struct{Line1 string, City string, Country string, Line2 string, PostalCode string, State string}), Carrier string, TrackingNumber string}), Created int64, Updated int64}
      --> w-tainted: write(orders_db.Order.Amount) {1}
[__2] (FieldObject FieldType) Amount int64
       --> w-tainted: write(orders_db.Order.Amount) {1}
[___3] (BasicObject BasicType) amount int64
      --> w-tainted: write(orders_db.Order.Currency) {1}
[__2] (FieldObject FieldType) Currency int32
       --> w-tainted: write(orders_db.Order.Currency) {1}
[___3] (BasicObject BasicType) currency int32
      --> w-tainted: write(orders_db.Order.Email) {1}
[__2] (FieldObject FieldType) Email string
       --> w-tainted: write(orders_db.Order.Email) {1}
[___3] (BasicObject BasicType) email string
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (FieldObject FieldType) Items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
      --> w-tainted: write(orders_db.Order.Metadata) {1}
[__2] (FieldObject FieldType) Metadata map[string]string
       --> w-tainted: write(orders_db.Order.Metadata) {1}
[___3] (MapObject MapType) metadata map[string]string
      --> w-tainted: write(orders_db.Order.Shipping) {1}
[__2] (FieldObject FieldType) Shipping (*digota.Shipping struct{Name string, Phone string, Address (*digota.Shipping_Address struct{Line1 string, City string, Country string, Line2 string, PostalCode string, State string}), Carrier string, TrackingNumber string})
       --> w-tainted: write(orders_db.Order.Shipping) {1}
[___3] (PointerObject PointerType) shipping (*digota.Shipping struct{Name string, Phone string, Address (*digota.Shipping_Address struct{Line1 string, City string, Country string, Line2 string, PostalCode string, State string}), Carrier string, TrackingNumber string})
        --> w-tainted: write(orders_db.Order.Shipping) {1}
[____4] (StructObject UserType) digota.Shipping struct{Name string, Phone string, Address (*digota.Shipping_Address struct{Line1 string, City string, Country string, Line2 string, PostalCode string, State string}), Carrier string, TrackingNumber string}

[0] (ArrayObject ArrayType) orderItems [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})

[0] (PointerObject PointerType) myitem1 (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
     --> w-tainted: write(orders_db.Order.Items) {1}
[_1] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
[_1] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
[__2] (FieldObject FieldType) Amount int64
[___3] (BasicObject BasicType) int64
[____4] (BasicObject BasicType) Price uint64
[__2] (FieldObject FieldType) Currency int32
[___3] (BasicObject BasicType) Currency int32
[__2] (FieldObject FieldType) Description string
[___3] (BasicObject BasicType) Name string
[__2] (FieldObject FieldType) Parent string
       --> r-tainted: read(skus_db._.id) {1}
[___3] (BasicObject BasicType) Parent string

[0] (PointerObject PointerType) itemFromSku (*digota.Sku struct{Id string, Name string, Price uint64, Currency int32, Active bool, Parent string, Metadata map[string]string, Attributes map[string]string, Image string, PackageDimensions (*digota.PackageDimensions struct{Height float64, Length float64, Weight float64, Width float64}), Inventory (*digota.Inventory struct{Quantity int64, Type int32}), Created int64, Updated int64})
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

    --> w-tainted: write(orders_db.Order.Currency) {1}
[0] (BasicObject BasicType) currency int32

    --> w-tainted: write(orders_db.Order.Items) {1}
[0] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})

    --> w-tainted: write(orders_db.Order.Metadata) {1}
[0] (MapObject MapType) metadata map[string]string

    --> w-tainted: write(orders_db.Order.Email) {1}
[0] (BasicObject BasicType) email string

    --> w-tainted: write(orders_db.Order.Shipping) {1}
[0] (PointerObject PointerType) shipping (*digota.Shipping struct{Name string, Phone string, Address (*digota.Shipping_Address struct{Line1 string, City string, Country string, Line2 string, PostalCode string, State string}), Carrier string, TrackingNumber string})
     --> w-tainted: write(orders_db.Order.Shipping) {1}
[_1] (StructObject UserType) digota.Shipping struct{Name string, Phone string, Address (*digota.Shipping_Address struct{Line1 string, City string, Country string, Line2 string, PostalCode string, State string}), Carrier string, TrackingNumber string}

[0] (PointerObject PointerType) order (*digota.Order struct{Id string, Amount int64, Currency int32, Items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}), Metadata map[string]string, Email string, ChargeId string, Status int32, Shipping (*digota.Shipping struct{Name string, Phone string, Address (*digota.Shipping_Address struct{Line1 string, City string, Country string, Line2 string, PostalCode string, State string}), Carrier string, TrackingNumber string}), Created int64, Updated int64})
     --> w-tainted: write(orders_db.Order) {1}
[_1] (StructObject UserType) digota.Order struct{Id string, Amount int64, Currency int32, Items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}), Metadata map[string]string, Email string, ChargeId string, Status int32, Shipping (*digota.Shipping struct{Name string, Phone string, Address (*digota.Shipping_Address struct{Line1 string, City string, Country string, Line2 string, PostalCode string, State string}), Carrier string, TrackingNumber string}), Created int64, Updated int64}
      --> w-tainted: write(orders_db.Order.Amount) {1}
[__2] (FieldObject FieldType) Amount int64
       --> w-tainted: write(orders_db.Order.Amount) {1}
[___3] (BasicObject BasicType) amount int64
      --> w-tainted: write(orders_db.Order.Currency) {1}
[__2] (FieldObject FieldType) Currency int32
       --> w-tainted: write(orders_db.Order.Currency) {1}
[___3] (BasicObject BasicType) currency int32
      --> w-tainted: write(orders_db.Order.Email) {1}
[__2] (FieldObject FieldType) Email string
       --> w-tainted: write(orders_db.Order.Email) {1}
[___3] (BasicObject BasicType) email string
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (FieldObject FieldType) Items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
      --> w-tainted: write(orders_db.Order.Metadata) {1}
[__2] (FieldObject FieldType) Metadata map[string]string
       --> w-tainted: write(orders_db.Order.Metadata) {1}
[___3] (MapObject MapType) metadata map[string]string
      --> w-tainted: write(orders_db.Order.Shipping) {1}
[__2] (FieldObject FieldType) Shipping (*digota.Shipping struct{Name string, Phone string, Address (*digota.Shipping_Address struct{Line1 string, City string, Country string, Line2 string, PostalCode string, State string}), Carrier string, TrackingNumber string})
       --> w-tainted: write(orders_db.Order.Shipping) {1}
[___3] (PointerObject PointerType) shipping (*digota.Shipping struct{Name string, Phone string, Address (*digota.Shipping_Address struct{Line1 string, City string, Country string, Line2 string, PostalCode string, State string}), Carrier string, TrackingNumber string})
        --> w-tainted: write(orders_db.Order.Shipping) {1}
[____4] (StructObject UserType) digota.Shipping struct{Name string, Phone string, Address (*digota.Shipping_Address struct{Line1 string, City string, Country string, Line2 string, PostalCode string, State string}), Carrier string, TrackingNumber string}

[0] (ArrayObject ArrayType) orderItems [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})

[0] (PointerObject PointerType) myitem1 (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
     --> w-tainted: write(orders_db.Order.Items) {1}
[_1] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
[_1] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
[__2] (FieldObject FieldType) Amount int64
[___3] (BasicObject BasicType) int64
[____4] (BasicObject BasicType) Price uint64
[__2] (FieldObject FieldType) Currency int32
[___3] (BasicObject BasicType) Currency int32
[__2] (FieldObject FieldType) Description string
[___3] (BasicObject BasicType) Name string
[__2] (FieldObject FieldType) Parent string
       --> r-tainted: read(skus_db._.id) {1}
[___3] (BasicObject BasicType) Parent string

    --> w-tainted: write(orders_db.Order.Amount) {1}
[0] (BasicObject BasicType) amount int64

[0] (InterfaceObject UserType) err .error

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = orders, collection = orders}

[0] (InterfaceObject UserType) err .error

[0] (InterfaceObject UserType) err .error

