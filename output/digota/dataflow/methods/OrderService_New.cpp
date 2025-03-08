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
[___3] (SliceObject SliceType) orderItems [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (Reference ArrayType) ref <orderItems [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ getUpdatedOrderItems
         --> w-tainted: write(orders_db.Order.Items) {1}
[_____5] (PointerObject PointerType) v (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
          --> w-tainted: write(orders_db.Order.Items) {1}
[______6] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
           --> w-tainted: write(orders_db.Order.Items) {1}
[_______7] (Reference ArrayType) ref <items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ OrderService
          --> w-tainted: write(orders_db.Order.Items) {1}
[______6] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
           --> w-tainted: write(orders_db.Order.Items) {1}
[_______7] (FieldObject FieldType) Parent string
            --> w-tainted: write(orders_db.Order.Items) {1}
[________8] (BasicObject BasicType) Parent string
           --> w-tainted: write(orders_db.Order.Items) {1}
[_______7] (FieldObject FieldType) Quantity int64
            --> w-tainted: write(orders_db.Order.Items) {1}
[________8] (BasicObject BasicType) 1 int
         --> w-tainted: write(orders_db.Order.Items) {1}
[_____5] (PointerObject PointerType) v (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
          --> w-tainted: write(orders_db.Order.Items) {1}
[______6] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
           --> w-tainted: write(orders_db.Order.Items) {1}
[_______7] (Reference ArrayType) ref <items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ OrderService
          --> w-tainted: write(orders_db.Order.Items) {1}
[______6] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
           --> w-tainted: write(orders_db.Order.Items) {1}
[_______7] (FieldObject FieldType) Parent string
            --> w-tainted: write(orders_db.Order.Items) {1}
[________8] (BasicObject BasicType) Parent string
           --> w-tainted: write(orders_db.Order.Items) {1}
[_______7] (FieldObject FieldType) Quantity int64
            --> w-tainted: write(orders_db.Order.Items) {1}
[________8] (BasicObject BasicType) 1 int
         --> w-tainted: write(orders_db.Order.Items) {1}
[_____5] (PointerObject PointerType) v (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
          --> w-tainted: write(orders_db.Order.Items) {1}
[______6] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
           --> w-tainted: write(orders_db.Order.Items) {1}
[_______7] (Reference ArrayType) ref <items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ OrderService
          --> w-tainted: write(orders_db.Order.Items) {1}
[______6] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
           --> w-tainted: write(orders_db.Order.Items) {1}
[_______7] (FieldObject FieldType) Parent string
            --> w-tainted: write(orders_db.Order.Items) {1}
[________8] (BasicObject BasicType) Parent string
           --> w-tainted: write(orders_db.Order.Items) {1}
[_______7] (FieldObject FieldType) Quantity int64
            --> w-tainted: write(orders_db.Order.Items) {1}
[________8] (BasicObject BasicType) 1 int
         --> w-tainted: write(orders_db.Order.Items) {1}
[_____5] (PointerObject PointerType) v (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
          --> w-tainted: write(orders_db.Order.Items) {1}
[______6] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
           --> w-tainted: write(orders_db.Order.Items) {1}
[_______7] (Reference ArrayType) ref <items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ OrderService
          --> w-tainted: write(orders_db.Order.Items) {1}
[______6] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
           --> w-tainted: write(orders_db.Order.Items) {1}
[_______7] (FieldObject FieldType) Parent string
            --> w-tainted: write(orders_db.Order.Items) {1}
[________8] (BasicObject BasicType) Parent string
           --> w-tainted: write(orders_db.Order.Items) {1}
[_______7] (FieldObject FieldType) Quantity int64
            --> w-tainted: write(orders_db.Order.Items) {1}
[________8] (BasicObject BasicType) 1 int
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

    --> w-tainted: write(orders_db.Order.Items) {1}
[0] (SliceObject SliceType) orderItems [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
     --> w-tainted: write(orders_db.Order.Items) {1}
[_1] (Reference ArrayType) ref <orderItems [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ getUpdatedOrderItems
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

[0] (InterfaceObject UserType) err .error
[_1] (Reference BasicType) ref <nil> @ getUpdatedOrderItems

    --> w-tainted: write(orders_db.Order.Amount) {1}
[0] (BasicObject BasicType) amount int64

[0] (InterfaceObject UserType) err .error

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = orders, collection = orders}

[0] (InterfaceObject UserType) err .error

[0] (InterfaceObject UserType) err .error

