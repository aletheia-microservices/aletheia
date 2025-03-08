[0] (PointerObject PointerType) s (*digota.ProductServiceImpl struct{db NoSQLDatabase})
[_1] (StructObject UserType) digota.ProductServiceImpl struct{db NoSQLDatabase}
[__2] (FieldObject FieldType) db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context

    --> r-tainted: read(products_db.Product.Id) {1}
[0] (BasicObject BasicType) id string

    --> w-tainted: write(products_db.Product.Name) {1}       --> w-tainted: write(products_db.Product.Name) {1} --> r-tainted: read(products_db.Product.Name) {1}
[0] (BasicObject BasicType) name string

    --> w-tainted: write(products_db.Product.Active) {1}       --> w-tainted: write(products_db.Product.Active) {1} --> r-tainted: read(products_db.Product.Active) {1}
[0] (BasicObject BasicType) active bool

    --> w-tainted: write(products_db.Product.Attributes) {1}       --> w-tainted: write(products_db.Product.Attributes) {1} --> r-tainted: read(products_db.Product.Attributes) {1}
[0] (ArrayObject SliceType) attributes []string

    --> w-tainted: write(products_db.Product.Description) {1}       --> w-tainted: write(products_db.Product.Description) {1} --> r-tainted: read(products_db.Product.Description) {1}
[0] (BasicObject BasicType) description string

    --> w-tainted: write(products_db.Product.Images) {1}       --> w-tainted: write(products_db.Product.Images) {1} --> r-tainted: read(products_db.Product.Images) {1}
[0] (ArrayObject SliceType) images []string

    --> w-tainted: write(products_db.Product.Metadata) {1}       --> w-tainted: write(products_db.Product.Metadata) {1} --> r-tainted: read(products_db.Product.Metadata) {1}
[0] (MapObject MapType) metadata map[string]string

    --> w-tainted: write(products_db.Product.Shippable) {1}       --> w-tainted: write(products_db.Product.Shippable) {1} --> r-tainted: read(products_db.Product.Shippable) {1}
[0] (BasicObject BasicType) shippable bool

    --> w-tainted: write(products_db.Product.Url) {1}       --> w-tainted: write(products_db.Product.Url) {1} --> r-tainted: read(products_db.Product.Url) {1}
[0] (BasicObject BasicType) url string

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
     --> w-tainted: write(products_db.Product) {1}         --> w-tainted: write(products_db.Product) {1} --> r-tainted: read(products_db.Product, products_db.Product.Active, products_db.Product.Attributes, products_db.Product.Description, products_db.Product.Images, products_db.Product.Metadata, products_db.Product.Name, products_db.Product.Shippable, products_db.Product.Url) {9}
[_1] (StructObject UserType) digota.Product struct{Id string, Name string, Active bool, Attributes []string, Description string, Images []string, Metadata map[string]string, Shippable bool, Url string, Skus [](*digota.Sku struct{Id string, Name string, Price uint64, Currency int32, Active bool, Parent string, Metadata map[string]string, Attributes map[string]string, Image string, PackageDimensions (*digota.PackageDimensions struct{Height float64, Length float64, Weight float64, Width float64}), Inventory (*digota.Inventory struct{Quantity int64, Type int32}), Created int64, Updated int64}), Created int64, Updated int64}
      --> w-tainted: write(products_db.Product.Active) {1}           --> w-tainted: write(products_db.Product.Active) {1} --> r-tainted: read(products_db.Product.Active) {1}
[__2] (FieldObject FieldType) Active bool
       --> w-tainted: write(products_db.Product.Active) {1}             --> w-tainted: write(products_db.Product.Active) {1} --> r-tainted: read(products_db.Product.Active) {1}
[___3] (BasicObject BasicType) active bool
      --> w-tainted: write(products_db.Product.Attributes) {1}           --> w-tainted: write(products_db.Product.Attributes) {1} --> r-tainted: read(products_db.Product.Attributes) {1}
[__2] (FieldObject FieldType) Attributes []string
       --> w-tainted: write(products_db.Product.Attributes) {1}             --> w-tainted: write(products_db.Product.Attributes) {1} --> r-tainted: read(products_db.Product.Attributes) {1}
[___3] (ArrayObject SliceType) attributes []string
      --> w-tainted: write(products_db.Product.Description) {1}           --> w-tainted: write(products_db.Product.Description) {1} --> r-tainted: read(products_db.Product.Description) {1}
[__2] (FieldObject FieldType) Description string
       --> w-tainted: write(products_db.Product.Description) {1}             --> w-tainted: write(products_db.Product.Description) {1} --> r-tainted: read(products_db.Product.Description) {1}
[___3] (BasicObject BasicType) description string
      --> w-tainted: write(products_db.Product.Images) {1}           --> w-tainted: write(products_db.Product.Images) {1} --> r-tainted: read(products_db.Product.Images) {1}
[__2] (FieldObject FieldType) Images []string
       --> w-tainted: write(products_db.Product.Images) {1}             --> w-tainted: write(products_db.Product.Images) {1} --> r-tainted: read(products_db.Product.Images) {1}
[___3] (ArrayObject SliceType) images []string
      --> w-tainted: write(products_db.Product.Metadata) {1}           --> w-tainted: write(products_db.Product.Metadata) {1} --> r-tainted: read(products_db.Product.Metadata) {1}
[__2] (FieldObject FieldType) Metadata map[string]string
       --> w-tainted: write(products_db.Product.Metadata) {1}             --> w-tainted: write(products_db.Product.Metadata) {1} --> r-tainted: read(products_db.Product.Metadata) {1}
[___3] (MapObject MapType) metadata map[string]string
      --> w-tainted: write(products_db.Product.Name) {1}           --> w-tainted: write(products_db.Product.Name) {1} --> r-tainted: read(products_db.Product.Name) {1}
[__2] (FieldObject FieldType) Name string
       --> w-tainted: write(products_db.Product.Name) {1}             --> w-tainted: write(products_db.Product.Name) {1} --> r-tainted: read(products_db.Product.Name) {1}
[___3] (BasicObject BasicType) name string
      --> w-tainted: write(products_db.Product.Shippable) {1}           --> w-tainted: write(products_db.Product.Shippable) {1} --> r-tainted: read(products_db.Product.Shippable) {1}
[__2] (FieldObject FieldType) Shippable bool
       --> w-tainted: write(products_db.Product.Shippable) {1}             --> w-tainted: write(products_db.Product.Shippable) {1} --> r-tainted: read(products_db.Product.Shippable) {1}
[___3] (BasicObject BasicType) shippable bool
      --> w-tainted: write(products_db.Product.Url) {1}           --> w-tainted: write(products_db.Product.Url) {1} --> r-tainted: read(products_db.Product.Url) {1}
[__2] (FieldObject FieldType) Url string
       --> w-tainted: write(products_db.Product.Url) {1}             --> w-tainted: write(products_db.Product.Url) {1} --> r-tainted: read(products_db.Product.Url) {1}
[___3] (BasicObject BasicType) url string

[0] (InterfaceObject UserType) err .error

[0] (PointerObject PointerType) product (*digota.Product struct{Id string, Name string, Active bool, Attributes []string, Description string, Images []string, Metadata map[string]string, Shippable bool, Url string, Skus [](*digota.Sku struct{Id string, Name string, Price uint64, Currency int32, Active bool, Parent string, Metadata map[string]string, Attributes map[string]string, Image string, PackageDimensions (*digota.PackageDimensions struct{Height float64, Length float64, Weight float64, Width float64}), Inventory (*digota.Inventory struct{Quantity int64, Type int32}), Created int64, Updated int64}), Created int64, Updated int64})
     --> w-tainted: write(products_db.Product) {1}         --> w-tainted: write(products_db.Product) {1} --> r-tainted: read(products_db.Product, products_db.Product.Active, products_db.Product.Attributes, products_db.Product.Description, products_db.Product.Images, products_db.Product.Metadata, products_db.Product.Name, products_db.Product.Shippable, products_db.Product.Url) {9}
[_1] (StructObject UserType) digota.Product struct{Id string, Name string, Active bool, Attributes []string, Description string, Images []string, Metadata map[string]string, Shippable bool, Url string, Skus [](*digota.Sku struct{Id string, Name string, Price uint64, Currency int32, Active bool, Parent string, Metadata map[string]string, Attributes map[string]string, Image string, PackageDimensions (*digota.PackageDimensions struct{Height float64, Length float64, Weight float64, Width float64}), Inventory (*digota.Inventory struct{Quantity int64, Type int32}), Created int64, Updated int64}), Created int64, Updated int64}
      --> w-tainted: write(products_db.Product.Active) {1}           --> w-tainted: write(products_db.Product.Active) {1} --> r-tainted: read(products_db.Product.Active) {1}
[__2] (FieldObject FieldType) Active bool
       --> w-tainted: write(products_db.Product.Active) {1}             --> w-tainted: write(products_db.Product.Active) {1} --> r-tainted: read(products_db.Product.Active) {1}
[___3] (BasicObject BasicType) active bool
      --> w-tainted: write(products_db.Product.Attributes) {1}           --> w-tainted: write(products_db.Product.Attributes) {1} --> r-tainted: read(products_db.Product.Attributes) {1}
[__2] (FieldObject FieldType) Attributes []string
       --> w-tainted: write(products_db.Product.Attributes) {1}             --> w-tainted: write(products_db.Product.Attributes) {1} --> r-tainted: read(products_db.Product.Attributes) {1}
[___3] (ArrayObject SliceType) attributes []string
      --> w-tainted: write(products_db.Product.Description) {1}           --> w-tainted: write(products_db.Product.Description) {1} --> r-tainted: read(products_db.Product.Description) {1}
[__2] (FieldObject FieldType) Description string
       --> w-tainted: write(products_db.Product.Description) {1}             --> w-tainted: write(products_db.Product.Description) {1} --> r-tainted: read(products_db.Product.Description) {1}
[___3] (BasicObject BasicType) description string
      --> w-tainted: write(products_db.Product.Images) {1}           --> w-tainted: write(products_db.Product.Images) {1} --> r-tainted: read(products_db.Product.Images) {1}
[__2] (FieldObject FieldType) Images []string
       --> w-tainted: write(products_db.Product.Images) {1}             --> w-tainted: write(products_db.Product.Images) {1} --> r-tainted: read(products_db.Product.Images) {1}
[___3] (ArrayObject SliceType) images []string
      --> w-tainted: write(products_db.Product.Metadata) {1}           --> w-tainted: write(products_db.Product.Metadata) {1} --> r-tainted: read(products_db.Product.Metadata) {1}
[__2] (FieldObject FieldType) Metadata map[string]string
       --> w-tainted: write(products_db.Product.Metadata) {1}             --> w-tainted: write(products_db.Product.Metadata) {1} --> r-tainted: read(products_db.Product.Metadata) {1}
[___3] (MapObject MapType) metadata map[string]string
      --> w-tainted: write(products_db.Product.Name) {1}           --> w-tainted: write(products_db.Product.Name) {1} --> r-tainted: read(products_db.Product.Name) {1}
[__2] (FieldObject FieldType) Name string
       --> w-tainted: write(products_db.Product.Name) {1}             --> w-tainted: write(products_db.Product.Name) {1} --> r-tainted: read(products_db.Product.Name) {1}
[___3] (BasicObject BasicType) name string
      --> w-tainted: write(products_db.Product.Shippable) {1}           --> w-tainted: write(products_db.Product.Shippable) {1} --> r-tainted: read(products_db.Product.Shippable) {1}
[__2] (FieldObject FieldType) Shippable bool
       --> w-tainted: write(products_db.Product.Shippable) {1}             --> w-tainted: write(products_db.Product.Shippable) {1} --> r-tainted: read(products_db.Product.Shippable) {1}
[___3] (BasicObject BasicType) shippable bool
      --> w-tainted: write(products_db.Product.Url) {1}           --> w-tainted: write(products_db.Product.Url) {1} --> r-tainted: read(products_db.Product.Url) {1}
[__2] (FieldObject FieldType) Url string
       --> w-tainted: write(products_db.Product.Url) {1}             --> w-tainted: write(products_db.Product.Url) {1} --> r-tainted: read(products_db.Product.Url) {1}
[___3] (BasicObject BasicType) url string

[0] (BasicObject BasicType) found bool

[0] (InterfaceObject UserType) err .error

[0] (PointerObject PointerType) s (*digota.ProductServiceImpl struct{db NoSQLDatabase})
[_1] (StructObject UserType) digota.ProductServiceImpl struct{db NoSQLDatabase}
[__2] (FieldObject FieldType) db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context

    --> r-tainted: read(products_db.Product.Id) {1}
[0] (BasicObject BasicType) id string

    --> w-tainted: write(products_db.Product.Name) {1}       --> w-tainted: write(products_db.Product.Name) {1} --> r-tainted: read(products_db.Product.Name) {1}
[0] (BasicObject BasicType) name string

    --> w-tainted: write(products_db.Product.Active) {1}       --> w-tainted: write(products_db.Product.Active) {1} --> r-tainted: read(products_db.Product.Active) {1}
[0] (BasicObject BasicType) active bool

    --> w-tainted: write(products_db.Product.Attributes) {1}       --> w-tainted: write(products_db.Product.Attributes) {1} --> r-tainted: read(products_db.Product.Attributes) {1}
[0] (ArrayObject SliceType) attributes []string

    --> w-tainted: write(products_db.Product.Description) {1}       --> w-tainted: write(products_db.Product.Description) {1} --> r-tainted: read(products_db.Product.Description) {1}
[0] (BasicObject BasicType) description string

    --> w-tainted: write(products_db.Product.Images) {1}       --> w-tainted: write(products_db.Product.Images) {1} --> r-tainted: read(products_db.Product.Images) {1}
[0] (ArrayObject SliceType) images []string

    --> w-tainted: write(products_db.Product.Metadata) {1}       --> w-tainted: write(products_db.Product.Metadata) {1} --> r-tainted: read(products_db.Product.Metadata) {1}
[0] (MapObject MapType) metadata map[string]string

    --> w-tainted: write(products_db.Product.Shippable) {1}       --> w-tainted: write(products_db.Product.Shippable) {1} --> r-tainted: read(products_db.Product.Shippable) {1}
[0] (BasicObject BasicType) shippable bool

    --> w-tainted: write(products_db.Product.Url) {1}       --> w-tainted: write(products_db.Product.Url) {1} --> r-tainted: read(products_db.Product.Url) {1}
[0] (BasicObject BasicType) url string

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
     --> w-tainted: write(products_db.Product) {1}         --> w-tainted: write(products_db.Product) {1} --> r-tainted: read(products_db.Product, products_db.Product.Active, products_db.Product.Attributes, products_db.Product.Description, products_db.Product.Images, products_db.Product.Metadata, products_db.Product.Name, products_db.Product.Shippable, products_db.Product.Url) {9}
[_1] (StructObject UserType) digota.Product struct{Id string, Name string, Active bool, Attributes []string, Description string, Images []string, Metadata map[string]string, Shippable bool, Url string, Skus [](*digota.Sku struct{Id string, Name string, Price uint64, Currency int32, Active bool, Parent string, Metadata map[string]string, Attributes map[string]string, Image string, PackageDimensions (*digota.PackageDimensions struct{Height float64, Length float64, Weight float64, Width float64}), Inventory (*digota.Inventory struct{Quantity int64, Type int32}), Created int64, Updated int64}), Created int64, Updated int64}
      --> w-tainted: write(products_db.Product.Active) {1}           --> w-tainted: write(products_db.Product.Active) {1} --> r-tainted: read(products_db.Product.Active) {1}
[__2] (FieldObject FieldType) Active bool
       --> w-tainted: write(products_db.Product.Active) {1}             --> w-tainted: write(products_db.Product.Active) {1} --> r-tainted: read(products_db.Product.Active) {1}
[___3] (BasicObject BasicType) active bool
      --> w-tainted: write(products_db.Product.Attributes) {1}           --> w-tainted: write(products_db.Product.Attributes) {1} --> r-tainted: read(products_db.Product.Attributes) {1}
[__2] (FieldObject FieldType) Attributes []string
       --> w-tainted: write(products_db.Product.Attributes) {1}             --> w-tainted: write(products_db.Product.Attributes) {1} --> r-tainted: read(products_db.Product.Attributes) {1}
[___3] (ArrayObject SliceType) attributes []string
      --> w-tainted: write(products_db.Product.Description) {1}           --> w-tainted: write(products_db.Product.Description) {1} --> r-tainted: read(products_db.Product.Description) {1}
[__2] (FieldObject FieldType) Description string
       --> w-tainted: write(products_db.Product.Description) {1}             --> w-tainted: write(products_db.Product.Description) {1} --> r-tainted: read(products_db.Product.Description) {1}
[___3] (BasicObject BasicType) description string
      --> w-tainted: write(products_db.Product.Images) {1}           --> w-tainted: write(products_db.Product.Images) {1} --> r-tainted: read(products_db.Product.Images) {1}
[__2] (FieldObject FieldType) Images []string
       --> w-tainted: write(products_db.Product.Images) {1}             --> w-tainted: write(products_db.Product.Images) {1} --> r-tainted: read(products_db.Product.Images) {1}
[___3] (ArrayObject SliceType) images []string
      --> w-tainted: write(products_db.Product.Metadata) {1}           --> w-tainted: write(products_db.Product.Metadata) {1} --> r-tainted: read(products_db.Product.Metadata) {1}
[__2] (FieldObject FieldType) Metadata map[string]string
       --> w-tainted: write(products_db.Product.Metadata) {1}             --> w-tainted: write(products_db.Product.Metadata) {1} --> r-tainted: read(products_db.Product.Metadata) {1}
[___3] (MapObject MapType) metadata map[string]string
      --> w-tainted: write(products_db.Product.Name) {1}           --> w-tainted: write(products_db.Product.Name) {1} --> r-tainted: read(products_db.Product.Name) {1}
[__2] (FieldObject FieldType) Name string
       --> w-tainted: write(products_db.Product.Name) {1}             --> w-tainted: write(products_db.Product.Name) {1} --> r-tainted: read(products_db.Product.Name) {1}
[___3] (BasicObject BasicType) name string
      --> w-tainted: write(products_db.Product.Shippable) {1}           --> w-tainted: write(products_db.Product.Shippable) {1} --> r-tainted: read(products_db.Product.Shippable) {1}
[__2] (FieldObject FieldType) Shippable bool
       --> w-tainted: write(products_db.Product.Shippable) {1}             --> w-tainted: write(products_db.Product.Shippable) {1} --> r-tainted: read(products_db.Product.Shippable) {1}
[___3] (BasicObject BasicType) shippable bool
      --> w-tainted: write(products_db.Product.Url) {1}           --> w-tainted: write(products_db.Product.Url) {1} --> r-tainted: read(products_db.Product.Url) {1}
[__2] (FieldObject FieldType) Url string
       --> w-tainted: write(products_db.Product.Url) {1}             --> w-tainted: write(products_db.Product.Url) {1} --> r-tainted: read(products_db.Product.Url) {1}
[___3] (BasicObject BasicType) url string

[0] (InterfaceObject UserType) err .error

[0] (PointerObject PointerType) product (*digota.Product struct{Id string, Name string, Active bool, Attributes []string, Description string, Images []string, Metadata map[string]string, Shippable bool, Url string, Skus [](*digota.Sku struct{Id string, Name string, Price uint64, Currency int32, Active bool, Parent string, Metadata map[string]string, Attributes map[string]string, Image string, PackageDimensions (*digota.PackageDimensions struct{Height float64, Length float64, Weight float64, Width float64}), Inventory (*digota.Inventory struct{Quantity int64, Type int32}), Created int64, Updated int64}), Created int64, Updated int64})
     --> w-tainted: write(products_db.Product) {1}         --> w-tainted: write(products_db.Product) {1} --> r-tainted: read(products_db.Product, products_db.Product.Active, products_db.Product.Attributes, products_db.Product.Description, products_db.Product.Images, products_db.Product.Metadata, products_db.Product.Name, products_db.Product.Shippable, products_db.Product.Url) {9}
[_1] (StructObject UserType) digota.Product struct{Id string, Name string, Active bool, Attributes []string, Description string, Images []string, Metadata map[string]string, Shippable bool, Url string, Skus [](*digota.Sku struct{Id string, Name string, Price uint64, Currency int32, Active bool, Parent string, Metadata map[string]string, Attributes map[string]string, Image string, PackageDimensions (*digota.PackageDimensions struct{Height float64, Length float64, Weight float64, Width float64}), Inventory (*digota.Inventory struct{Quantity int64, Type int32}), Created int64, Updated int64}), Created int64, Updated int64}
      --> w-tainted: write(products_db.Product.Active) {1}           --> w-tainted: write(products_db.Product.Active) {1} --> r-tainted: read(products_db.Product.Active) {1}
[__2] (FieldObject FieldType) Active bool
       --> w-tainted: write(products_db.Product.Active) {1}             --> w-tainted: write(products_db.Product.Active) {1} --> r-tainted: read(products_db.Product.Active) {1}
[___3] (BasicObject BasicType) active bool
      --> w-tainted: write(products_db.Product.Attributes) {1}           --> w-tainted: write(products_db.Product.Attributes) {1} --> r-tainted: read(products_db.Product.Attributes) {1}
[__2] (FieldObject FieldType) Attributes []string
       --> w-tainted: write(products_db.Product.Attributes) {1}             --> w-tainted: write(products_db.Product.Attributes) {1} --> r-tainted: read(products_db.Product.Attributes) {1}
[___3] (ArrayObject SliceType) attributes []string
      --> w-tainted: write(products_db.Product.Description) {1}           --> w-tainted: write(products_db.Product.Description) {1} --> r-tainted: read(products_db.Product.Description) {1}
[__2] (FieldObject FieldType) Description string
       --> w-tainted: write(products_db.Product.Description) {1}             --> w-tainted: write(products_db.Product.Description) {1} --> r-tainted: read(products_db.Product.Description) {1}
[___3] (BasicObject BasicType) description string
      --> w-tainted: write(products_db.Product.Images) {1}           --> w-tainted: write(products_db.Product.Images) {1} --> r-tainted: read(products_db.Product.Images) {1}
[__2] (FieldObject FieldType) Images []string
       --> w-tainted: write(products_db.Product.Images) {1}             --> w-tainted: write(products_db.Product.Images) {1} --> r-tainted: read(products_db.Product.Images) {1}
[___3] (ArrayObject SliceType) images []string
      --> w-tainted: write(products_db.Product.Metadata) {1}           --> w-tainted: write(products_db.Product.Metadata) {1} --> r-tainted: read(products_db.Product.Metadata) {1}
[__2] (FieldObject FieldType) Metadata map[string]string
       --> w-tainted: write(products_db.Product.Metadata) {1}             --> w-tainted: write(products_db.Product.Metadata) {1} --> r-tainted: read(products_db.Product.Metadata) {1}
[___3] (MapObject MapType) metadata map[string]string
      --> w-tainted: write(products_db.Product.Name) {1}           --> w-tainted: write(products_db.Product.Name) {1} --> r-tainted: read(products_db.Product.Name) {1}
[__2] (FieldObject FieldType) Name string
       --> w-tainted: write(products_db.Product.Name) {1}             --> w-tainted: write(products_db.Product.Name) {1} --> r-tainted: read(products_db.Product.Name) {1}
[___3] (BasicObject BasicType) name string
      --> w-tainted: write(products_db.Product.Shippable) {1}           --> w-tainted: write(products_db.Product.Shippable) {1} --> r-tainted: read(products_db.Product.Shippable) {1}
[__2] (FieldObject FieldType) Shippable bool
       --> w-tainted: write(products_db.Product.Shippable) {1}             --> w-tainted: write(products_db.Product.Shippable) {1} --> r-tainted: read(products_db.Product.Shippable) {1}
[___3] (BasicObject BasicType) shippable bool
      --> w-tainted: write(products_db.Product.Url) {1}           --> w-tainted: write(products_db.Product.Url) {1} --> r-tainted: read(products_db.Product.Url) {1}
[__2] (FieldObject FieldType) Url string
       --> w-tainted: write(products_db.Product.Url) {1}             --> w-tainted: write(products_db.Product.Url) {1} --> r-tainted: read(products_db.Product.Url) {1}
[___3] (BasicObject BasicType) url string

[0] (BasicObject BasicType) found bool

[0] (InterfaceObject UserType) err .error

[0] (InterfaceObject UserType) err .error

