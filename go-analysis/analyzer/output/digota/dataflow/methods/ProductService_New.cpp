[0] (PointerObject PointerType) s (*digota.ProductServiceImpl struct{skuService digota.SkuService, db NoSQLDatabase})
[_1] (StructObject UserType) digota.ProductServiceImpl struct{skuService digota.SkuService, db NoSQLDatabase}
[__2] (FieldObject FieldType) db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase
[__2] (FieldObject FieldType) skuService digota.SkuService
[___3] (ServiceObject ServiceType) skuService digota.SkuService

[0] (InterfaceObject UserType) ctx context.Context

    --> w-tainted: write(products_db.Product.Name, skus_db.Sku.Name) {2}
[0] (BasicObject BasicType) name string

    --> w-tainted: write(products_db.Product.Active) {1}
[0] (BasicObject BasicType) active bool

    --> w-tainted: write(products_db.Product.Attributes) {1}
[0] (ArrayObject ArrayType) attributes []string

    --> w-tainted: write(products_db.Product.Description) {1}
[0] (BasicObject BasicType) description string

    --> w-tainted: write(products_db.Product.Images) {1}
[0] (ArrayObject ArrayType) images []string

    --> w-tainted: write(products_db.Product.Metadata) {1}
[0] (MapObject MapType) metadata map[string]string

    --> w-tainted: write(products_db.Product.Shippable) {1}
[0] (BasicObject BasicType) shippable bool

    --> w-tainted: write(products_db.Product.Url) {1}
[0] (BasicObject BasicType) url string

[0] (PointerObject PointerType) product (*digota.Product struct{Id string, Name string, Active bool, Attributes []string, Description string, Images []string, Metadata map[string]string, Shippable bool, Url string, Skus [](*digota.Sku struct{Id string, Name string, Price uint64, Currency int32, Active bool, Parent string, Metadata map[string]string, Attributes map[string]string, Image string, PackageDimensions (*digota.PackageDimensions struct{Height float64, Length float64, Weight float64, Width float64}), Inventory (*digota.Inventory struct{Quantity int64, Type int32}), Created int64, Updated int64}), Created int64, Updated int64})
     --> w-tainted: write(products_db.Product) {1}
[_1] (StructObject UserType) digota.Product struct{Id string, Name string, Active bool, Attributes []string, Description string, Images []string, Metadata map[string]string, Shippable bool, Url string, Skus [](*digota.Sku struct{Id string, Name string, Price uint64, Currency int32, Active bool, Parent string, Metadata map[string]string, Attributes map[string]string, Image string, PackageDimensions (*digota.PackageDimensions struct{Height float64, Length float64, Weight float64, Width float64}), Inventory (*digota.Inventory struct{Quantity int64, Type int32}), Created int64, Updated int64}), Created int64, Updated int64}
      --> w-tainted: write(products_db.Product.Active) {1}
[__2] (FieldObject FieldType) Active bool
       --> w-tainted: write(products_db.Product.Active) {1}
[___3] (BasicObject BasicType) active bool
      --> w-tainted: write(products_db.Product.Attributes) {1}
[__2] (FieldObject FieldType) Attributes []string
       --> w-tainted: write(products_db.Product.Attributes) {1}
[___3] (ArrayObject ArrayType) attributes []string
      --> w-tainted: write(products_db.Product.Description) {1}
[__2] (FieldObject FieldType) Description string
       --> w-tainted: write(products_db.Product.Description) {1}
[___3] (BasicObject BasicType) description string
      --> w-tainted: write(products_db.Product.Images) {1}
[__2] (FieldObject FieldType) Images []string
       --> w-tainted: write(products_db.Product.Images) {1}
[___3] (ArrayObject ArrayType) images []string
      --> w-tainted: write(products_db.Product.Metadata) {1}
[__2] (FieldObject FieldType) Metadata map[string]string
       --> w-tainted: write(products_db.Product.Metadata) {1}
[___3] (MapObject MapType) metadata map[string]string
      --> w-tainted: write(products_db.Product.Name) {1}
[__2] (FieldObject FieldType) Name string
       --> w-tainted: write(products_db.Product.Name, skus_db.Sku.Name) {2}
[___3] (BasicObject BasicType) name string
      --> w-tainted: write(products_db.Product.Shippable) {1}
[__2] (FieldObject FieldType) Shippable bool
       --> w-tainted: write(products_db.Product.Shippable) {1}
[___3] (BasicObject BasicType) shippable bool
      --> w-tainted: write(products_db.Product.Url) {1}
[__2] (FieldObject FieldType) Url string
       --> w-tainted: write(products_db.Product.Url) {1}
[___3] (BasicObject BasicType) url string

[0] (InterfaceObject UserType) err .error
[_1] (PointerObject PointerType) product (*digota.Product struct{Id string, Name string, Active bool, Attributes []string, Description string, Images []string, Metadata map[string]string, Shippable bool, Url string, Skus [](*digota.Sku struct{Id string, Name string, Price uint64, Currency int32, Active bool, Parent string, Metadata map[string]string, Attributes map[string]string, Image string, PackageDimensions (*digota.PackageDimensions struct{Height float64, Length float64, Weight float64, Width float64}), Inventory (*digota.Inventory struct{Quantity int64, Type int32}), Created int64, Updated int64}), Created int64, Updated int64})
      --> w-tainted: write(products_db.Product) {1}
[__2] (StructObject UserType) digota.Product struct{Id string, Name string, Active bool, Attributes []string, Description string, Images []string, Metadata map[string]string, Shippable bool, Url string, Skus [](*digota.Sku struct{Id string, Name string, Price uint64, Currency int32, Active bool, Parent string, Metadata map[string]string, Attributes map[string]string, Image string, PackageDimensions (*digota.PackageDimensions struct{Height float64, Length float64, Weight float64, Width float64}), Inventory (*digota.Inventory struct{Quantity int64, Type int32}), Created int64, Updated int64}), Created int64, Updated int64}
       --> w-tainted: write(products_db.Product.Active) {1}
[___3] (FieldObject FieldType) Active bool
        --> w-tainted: write(products_db.Product.Active) {1}
[____4] (BasicObject BasicType) active bool
       --> w-tainted: write(products_db.Product.Attributes) {1}
[___3] (FieldObject FieldType) Attributes []string
        --> w-tainted: write(products_db.Product.Attributes) {1}
[____4] (ArrayObject ArrayType) attributes []string
       --> w-tainted: write(products_db.Product.Description) {1}
[___3] (FieldObject FieldType) Description string
        --> w-tainted: write(products_db.Product.Description) {1}
[____4] (BasicObject BasicType) description string
       --> w-tainted: write(products_db.Product.Images) {1}
[___3] (FieldObject FieldType) Images []string
        --> w-tainted: write(products_db.Product.Images) {1}
[____4] (ArrayObject ArrayType) images []string
       --> w-tainted: write(products_db.Product.Metadata) {1}
[___3] (FieldObject FieldType) Metadata map[string]string
        --> w-tainted: write(products_db.Product.Metadata) {1}
[____4] (MapObject MapType) metadata map[string]string
       --> w-tainted: write(products_db.Product.Name) {1}
[___3] (FieldObject FieldType) Name string
        --> w-tainted: write(products_db.Product.Name, skus_db.Sku.Name) {2}
[____4] (BasicObject BasicType) name string
       --> w-tainted: write(products_db.Product.Shippable) {1}
[___3] (FieldObject FieldType) Shippable bool
        --> w-tainted: write(products_db.Product.Shippable) {1}
[____4] (BasicObject BasicType) shippable bool
       --> w-tainted: write(products_db.Product.Url) {1}
[___3] (FieldObject FieldType) Url string
        --> w-tainted: write(products_db.Product.Url) {1}
[____4] (BasicObject BasicType) url string

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = products, collection = products}

[0] (InterfaceObject UserType) err .error

[0] (InterfaceObject UserType) err .error

