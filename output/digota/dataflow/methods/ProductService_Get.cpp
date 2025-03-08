[0] (PointerObject PointerType) s (*digota.ProductServiceImpl struct{db NoSQLDatabase})
[_1] (StructObject UserType) digota.ProductServiceImpl struct{db NoSQLDatabase}
[__2] (FieldObject FieldType) db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context

    --> r-tainted: read(products_db.Product.Id) {1}
[0] (BasicObject BasicType) id string

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = products, collection = products}

[0] (InterfaceObject UserType) err .error

[0] (SliceObject UserType) query primitive.D
[_1] (StructObject StructType) struct{Key "id" string, Key "id" string, Value string, Value string}
[__2] (FieldObject FieldType) Key "id" string
[___3] (BasicObject BasicType) "id" string
[__2] (FieldObject FieldType) Value string
       --> r-tainted: read(products_db.Product.Id) {1}
[___3] (BasicObject BasicType) id string

    --> r-tainted: read(products_db.Product) {1}
[0] (BlueprintBackendObject BlueprintBackendType) result NoSQLCursor {database = products, collection = products}
     --> r-tainted: read(products_db.Product) {1}
[_1] (StructObject UserType) digota.Product struct{Id string, Name string, Active bool, Attributes []string, Description string, Images []string, Metadata map[string]string, Shippable bool, Url string, Skus [](*digota.Sku struct{Id string, Name string, Price uint64, Currency int32, Active bool, Parent string, Metadata map[string]string, Attributes map[string]string, Image string, PackageDimensions (*digota.PackageDimensions struct{Height float64, Length float64, Weight float64, Width float64}), Inventory (*digota.Inventory struct{Quantity int64, Type int32}), Created int64, Updated int64}), Created int64, Updated int64}

[0] (InterfaceObject UserType) err .error

[0] (PointerObject PointerType) product (*digota.Product struct{Id string, Name string, Active bool, Attributes []string, Description string, Images []string, Metadata map[string]string, Shippable bool, Url string, Skus [](*digota.Sku struct{Id string, Name string, Price uint64, Currency int32, Active bool, Parent string, Metadata map[string]string, Attributes map[string]string, Image string, PackageDimensions (*digota.PackageDimensions struct{Height float64, Length float64, Weight float64, Width float64}), Inventory (*digota.Inventory struct{Quantity int64, Type int32}), Created int64, Updated int64}), Created int64, Updated int64})
     --> r-tainted: read(products_db.Product) {1}
[_1] (StructObject UserType) digota.Product struct{Id string, Name string, Active bool, Attributes []string, Description string, Images []string, Metadata map[string]string, Shippable bool, Url string, Skus [](*digota.Sku struct{Id string, Name string, Price uint64, Currency int32, Active bool, Parent string, Metadata map[string]string, Attributes map[string]string, Image string, PackageDimensions (*digota.PackageDimensions struct{Height float64, Length float64, Weight float64, Width float64}), Inventory (*digota.Inventory struct{Quantity int64, Type int32}), Created int64, Updated int64}), Created int64, Updated int64}

[0] (BasicObject BasicType) found bool

[0] (InterfaceObject UserType) err .error

