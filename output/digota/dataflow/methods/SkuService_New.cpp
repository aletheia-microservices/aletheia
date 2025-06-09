[0] (PointerObject PointerType) s (*digota.SkuServiceImpl struct{db NoSQLDatabase})
[_1] (StructObject UserType) digota.SkuServiceImpl struct{db NoSQLDatabase}
[__2] (FieldObject FieldType) db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ ProductService

    --> w-tainted: write(skus_db.Sku.Name, products_db.Product.Name) {2}
[0] (BasicObject BasicType) name string
     --> w-tainted: write(products_db.Product.Name, skus_db.Sku.Name) {2}
[_1] (Reference BasicType) ref <name string> @ ProductService

    --> w-tainted: write(skus_db.Sku.Currency) {1}
[0] (BasicObject BasicType) currency int32
     --> w-tainted: write(skus_db.Sku.Currency) {1}
[_1] (Reference BasicType) ref <0 int> @ ProductService

    --> w-tainted: write(skus_db.Sku.Active) {1}
[0] (BasicObject BasicType) active bool
     --> w-tainted: write(skus_db.Sku.Active) {1}
[_1] (Reference BasicType) ref <true bool> @ ProductService

    --> w-tainted: write(skus_db.Sku.Price) {1}
[0] (BasicObject BasicType) price uint64
     --> w-tainted: write(skus_db.Sku.Price) {1}
[_1] (Reference BasicType) ref <0 int> @ ProductService

    --> w-tainted: write(skus_db.Sku.Parent) {1}
[0] (BasicObject BasicType) parent string
     --> w-tainted: write(skus_db.Sku.Parent) {1}
[_1] (Reference BasicType) ref <"parent" string> @ ProductService

    --> w-tainted: write(skus_db.Sku.Metadata) {1}
[0] (MapObject MapType) metadata map[string]string
     --> w-tainted: write(skus_db.Sku.Metadata) {1}
[_1] (Reference BasicType) ref <nil> @ ProductService

[0] (BasicObject BasicType) image string
[_1] (Reference BasicType) ref <"image" string> @ ProductService

    --> w-tainted: write(skus_db.Sku.PackageDimensions) {1}
[0] (PointerObject PointerType) packageDimensions (*digota.PackageDimensions struct{Height float64, Length float64, Weight float64, Width float64})
     --> w-tainted: write(skus_db.Sku.PackageDimensions) {1}
[_1] (StructObject UserType) digota.PackageDimensions struct{Height float64, Length float64, Weight float64, Width float64}

    --> w-tainted: write(skus_db.Sku.Inventory) {1}
[0] (PointerObject PointerType) inventory (*digota.Inventory struct{Quantity int64, Type int32})
     --> w-tainted: write(skus_db.Sku.Inventory) {1}
[_1] (StructObject UserType) digota.Inventory struct{Quantity int64, Type int32}

    --> w-tainted: write(skus_db.Sku.Attributes) {1}
[0] (MapObject MapType) attributes map[string]string
     --> w-tainted: write(skus_db.Sku.Attributes) {1}
[_1] (Reference BasicType) ref <nil> @ ProductService

[0] (PointerObject PointerType) sku (*digota.Sku struct{Id string, Name string, Price uint64, Currency int32, Active bool, Parent string, Metadata map[string]string, Attributes map[string]string, Image string, PackageDimensions (*digota.PackageDimensions struct{Height float64, Length float64, Weight float64, Width float64}), Inventory (*digota.Inventory struct{Quantity int64, Type int32}), Created int64, Updated int64})
     --> w-tainted: write(skus_db.Sku) {1}
[_1] (StructObject UserType) digota.Sku struct{Id string, Name string, Price uint64, Currency int32, Active bool, Parent string, Metadata map[string]string, Attributes map[string]string, Image string, PackageDimensions (*digota.PackageDimensions struct{Height float64, Length float64, Weight float64, Width float64}), Inventory (*digota.Inventory struct{Quantity int64, Type int32}), Created int64, Updated int64}
      --> w-tainted: write(skus_db.Sku.Active) {1}
[__2] (FieldObject FieldType) Active bool
       --> w-tainted: write(skus_db.Sku.Active) {1}
[___3] (BasicObject BasicType) active bool
        --> w-tainted: write(skus_db.Sku.Active) {1}
[____4] (Reference BasicType) ref <true bool> @ ProductService
      --> w-tainted: write(skus_db.Sku.Attributes) {1}
[__2] (FieldObject FieldType) Attributes map[string]string
       --> w-tainted: write(skus_db.Sku.Attributes) {1}
[___3] (MapObject MapType) attributes map[string]string
        --> w-tainted: write(skus_db.Sku.Attributes) {1}
[____4] (Reference BasicType) ref <nil> @ ProductService
      --> w-tainted: write(skus_db.Sku.Currency) {1}
[__2] (FieldObject FieldType) Currency int32
       --> w-tainted: write(skus_db.Sku.Currency) {1}
[___3] (BasicObject BasicType) currency int32
        --> w-tainted: write(skus_db.Sku.Currency) {1}
[____4] (Reference BasicType) ref <0 int> @ ProductService
      --> w-tainted: write(skus_db.Sku.Inventory) {1}
[__2] (FieldObject FieldType) Inventory (*digota.Inventory struct{Quantity int64, Type int32})
       --> w-tainted: write(skus_db.Sku.Inventory) {1}
[___3] (PointerObject PointerType) inventory (*digota.Inventory struct{Quantity int64, Type int32})
        --> w-tainted: write(skus_db.Sku.Inventory) {1}
[____4] (StructObject UserType) digota.Inventory struct{Quantity int64, Type int32}
      --> w-tainted: write(skus_db.Sku.Metadata) {1}
[__2] (FieldObject FieldType) Metadata map[string]string
       --> w-tainted: write(skus_db.Sku.Metadata) {1}
[___3] (MapObject MapType) metadata map[string]string
        --> w-tainted: write(skus_db.Sku.Metadata) {1}
[____4] (Reference BasicType) ref <nil> @ ProductService
      --> w-tainted: write(skus_db.Sku.Name) {1}
[__2] (FieldObject FieldType) Name string
       --> w-tainted: write(skus_db.Sku.Name, products_db.Product.Name) {2}
[___3] (BasicObject BasicType) name string
        --> w-tainted: write(products_db.Product.Name, skus_db.Sku.Name) {2}
[____4] (Reference BasicType) ref <name string> @ ProductService
      --> w-tainted: write(skus_db.Sku.PackageDimensions) {1}
[__2] (FieldObject FieldType) PackageDimensions (*digota.PackageDimensions struct{Height float64, Length float64, Weight float64, Width float64})
       --> w-tainted: write(skus_db.Sku.PackageDimensions) {1}
[___3] (PointerObject PointerType) packageDimensions (*digota.PackageDimensions struct{Height float64, Length float64, Weight float64, Width float64})
        --> w-tainted: write(skus_db.Sku.PackageDimensions) {1}
[____4] (StructObject UserType) digota.PackageDimensions struct{Height float64, Length float64, Weight float64, Width float64}
      --> w-tainted: write(skus_db.Sku.Parent) {1}
[__2] (FieldObject FieldType) Parent string
       --> w-tainted: write(skus_db.Sku.Parent) {1}
[___3] (BasicObject BasicType) parent string
        --> w-tainted: write(skus_db.Sku.Parent) {1}
[____4] (Reference BasicType) ref <"parent" string> @ ProductService
      --> w-tainted: write(skus_db.Sku.Price) {1}
[__2] (FieldObject FieldType) Price uint64
       --> w-tainted: write(skus_db.Sku.Price) {1}
[___3] (BasicObject BasicType) price uint64
        --> w-tainted: write(skus_db.Sku.Price) {1}
[____4] (Reference BasicType) ref <0 int> @ ProductService

[0] (InterfaceObject UserType) err .error
[_1] (PointerObject PointerType) sku (*digota.Sku struct{Id string, Name string, Price uint64, Currency int32, Active bool, Parent string, Metadata map[string]string, Attributes map[string]string, Image string, PackageDimensions (*digota.PackageDimensions struct{Height float64, Length float64, Weight float64, Width float64}), Inventory (*digota.Inventory struct{Quantity int64, Type int32}), Created int64, Updated int64})
      --> w-tainted: write(skus_db.Sku) {1}
[__2] (StructObject UserType) digota.Sku struct{Id string, Name string, Price uint64, Currency int32, Active bool, Parent string, Metadata map[string]string, Attributes map[string]string, Image string, PackageDimensions (*digota.PackageDimensions struct{Height float64, Length float64, Weight float64, Width float64}), Inventory (*digota.Inventory struct{Quantity int64, Type int32}), Created int64, Updated int64}
       --> w-tainted: write(skus_db.Sku.Active) {1}
[___3] (FieldObject FieldType) Active bool
        --> w-tainted: write(skus_db.Sku.Active) {1}
[____4] (BasicObject BasicType) active bool
         --> w-tainted: write(skus_db.Sku.Active) {1}
[_____5] (Reference BasicType) ref <true bool> @ ProductService
       --> w-tainted: write(skus_db.Sku.Attributes) {1}
[___3] (FieldObject FieldType) Attributes map[string]string
        --> w-tainted: write(skus_db.Sku.Attributes) {1}
[____4] (MapObject MapType) attributes map[string]string
         --> w-tainted: write(skus_db.Sku.Attributes) {1}
[_____5] (Reference BasicType) ref <nil> @ ProductService
       --> w-tainted: write(skus_db.Sku.Currency) {1}
[___3] (FieldObject FieldType) Currency int32
        --> w-tainted: write(skus_db.Sku.Currency) {1}
[____4] (BasicObject BasicType) currency int32
         --> w-tainted: write(skus_db.Sku.Currency) {1}
[_____5] (Reference BasicType) ref <0 int> @ ProductService
       --> w-tainted: write(skus_db.Sku.Inventory) {1}
[___3] (FieldObject FieldType) Inventory (*digota.Inventory struct{Quantity int64, Type int32})
        --> w-tainted: write(skus_db.Sku.Inventory) {1}
[____4] (PointerObject PointerType) inventory (*digota.Inventory struct{Quantity int64, Type int32})
         --> w-tainted: write(skus_db.Sku.Inventory) {1}
[_____5] (StructObject UserType) digota.Inventory struct{Quantity int64, Type int32}
       --> w-tainted: write(skus_db.Sku.Metadata) {1}
[___3] (FieldObject FieldType) Metadata map[string]string
        --> w-tainted: write(skus_db.Sku.Metadata) {1}
[____4] (MapObject MapType) metadata map[string]string
         --> w-tainted: write(skus_db.Sku.Metadata) {1}
[_____5] (Reference BasicType) ref <nil> @ ProductService
       --> w-tainted: write(skus_db.Sku.Name) {1}
[___3] (FieldObject FieldType) Name string
        --> w-tainted: write(skus_db.Sku.Name, products_db.Product.Name) {2}
[____4] (BasicObject BasicType) name string
         --> w-tainted: write(products_db.Product.Name, skus_db.Sku.Name) {2}
[_____5] (Reference BasicType) ref <name string> @ ProductService
       --> w-tainted: write(skus_db.Sku.PackageDimensions) {1}
[___3] (FieldObject FieldType) PackageDimensions (*digota.PackageDimensions struct{Height float64, Length float64, Weight float64, Width float64})
        --> w-tainted: write(skus_db.Sku.PackageDimensions) {1}
[____4] (PointerObject PointerType) packageDimensions (*digota.PackageDimensions struct{Height float64, Length float64, Weight float64, Width float64})
         --> w-tainted: write(skus_db.Sku.PackageDimensions) {1}
[_____5] (StructObject UserType) digota.PackageDimensions struct{Height float64, Length float64, Weight float64, Width float64}
       --> w-tainted: write(skus_db.Sku.Parent) {1}
[___3] (FieldObject FieldType) Parent string
        --> w-tainted: write(skus_db.Sku.Parent) {1}
[____4] (BasicObject BasicType) parent string
         --> w-tainted: write(skus_db.Sku.Parent) {1}
[_____5] (Reference BasicType) ref <"parent" string> @ ProductService
       --> w-tainted: write(skus_db.Sku.Price) {1}
[___3] (FieldObject FieldType) Price uint64
        --> w-tainted: write(skus_db.Sku.Price) {1}
[____4] (BasicObject BasicType) price uint64
         --> w-tainted: write(skus_db.Sku.Price) {1}
[_____5] (Reference BasicType) ref <0 int> @ ProductService

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = skus, collection = skus}

[0] (InterfaceObject UserType) err .error

[0] (InterfaceObject UserType) err .error

