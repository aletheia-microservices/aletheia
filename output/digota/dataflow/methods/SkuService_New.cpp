[0] (PointerObject PointerType) s (*digota.SkuServiceImpl struct{db NoSQLDatabase})
[_1] (StructObject UserType) digota.SkuServiceImpl struct{db NoSQLDatabase}
[__2] (FieldObject FieldType) db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context

[0] (BasicObject BasicType) name string

[0] (BasicObject BasicType) currency int32

[0] (BasicObject BasicType) active bool

[0] (BasicObject BasicType) price uint64

[0] (BasicObject BasicType) parent string

[0] (MapObject MapType) metadata map[string]string

[0] (BasicObject BasicType) image string

[0] (PointerObject PointerType) packageDimensions (*digota.PackageDimensions struct{Height float64, Length float64, Weight float64, Width float64})
[_1] (StructObject UserType) digota.PackageDimensions struct{Height float64, Length float64, Weight float64, Width float64}

[0] (PointerObject PointerType) inventory (*digota.Inventory struct{Quantity int64, Type int32})
[_1] (StructObject UserType) digota.Inventory struct{Quantity int64, Type int32}

[0] (MapObject MapType) attributes map[string]string

[0] (PointerObject PointerType) sku (*digota.Sku struct{Id string, Name string, Price uint64, Currency int32, Active bool, Parent string, Metadata map[string]string, Attributes map[string]string, Image string, PackageDimensions (*digota.PackageDimensions struct{Height float64, Length float64, Weight float64, Width float64}), Inventory (*digota.Inventory struct{Quantity int64, Type int32}), Created int64, Updated int64})
[_1] (StructObject UserType) digota.Sku struct{Id string, Name string, Price uint64, Currency int32, Active bool, Parent string, Metadata map[string]string, Attributes map[string]string, Image string, PackageDimensions (*digota.PackageDimensions struct{Height float64, Length float64, Weight float64, Width float64}), Inventory (*digota.Inventory struct{Quantity int64, Type int32}), Created int64, Updated int64}
[__2] (FieldObject FieldType) Active bool
[___3] (BasicObject BasicType) active bool
[__2] (FieldObject FieldType) Attributes map[string]string
[___3] (MapObject MapType) attributes map[string]string
[__2] (FieldObject FieldType) Currency int32
[___3] (BasicObject BasicType) currency int32
[__2] (FieldObject FieldType) Inventory (*digota.Inventory struct{Quantity int64, Type int32})
[___3] (PointerObject PointerType) inventory (*digota.Inventory struct{Quantity int64, Type int32})
[____4] (StructObject UserType) digota.Inventory struct{Quantity int64, Type int32}
[__2] (FieldObject FieldType) Metadata map[string]string
[___3] (MapObject MapType) metadata map[string]string
[__2] (FieldObject FieldType) Name string
[___3] (BasicObject BasicType) name string
[__2] (FieldObject FieldType) PackageDimensions (*digota.PackageDimensions struct{Height float64, Length float64, Weight float64, Width float64})
[___3] (PointerObject PointerType) packageDimensions (*digota.PackageDimensions struct{Height float64, Length float64, Weight float64, Width float64})
[____4] (StructObject UserType) digota.PackageDimensions struct{Height float64, Length float64, Weight float64, Width float64}
[__2] (FieldObject FieldType) Parent string
[___3] (BasicObject BasicType) parent string
[__2] (FieldObject FieldType) Price uint64
[___3] (BasicObject BasicType) price uint64

[0] (InterfaceObject UserType) err .error
[_1] (PointerObject PointerType) sku (*digota.Sku struct{Id string, Name string, Price uint64, Currency int32, Active bool, Parent string, Metadata map[string]string, Attributes map[string]string, Image string, PackageDimensions (*digota.PackageDimensions struct{Height float64, Length float64, Weight float64, Width float64}), Inventory (*digota.Inventory struct{Quantity int64, Type int32}), Created int64, Updated int64})
[__2] (StructObject UserType) digota.Sku struct{Id string, Name string, Price uint64, Currency int32, Active bool, Parent string, Metadata map[string]string, Attributes map[string]string, Image string, PackageDimensions (*digota.PackageDimensions struct{Height float64, Length float64, Weight float64, Width float64}), Inventory (*digota.Inventory struct{Quantity int64, Type int32}), Created int64, Updated int64}
[___3] (FieldObject FieldType) Active bool
[____4] (BasicObject BasicType) active bool
[___3] (FieldObject FieldType) Attributes map[string]string
[____4] (MapObject MapType) attributes map[string]string
[___3] (FieldObject FieldType) Currency int32
[____4] (BasicObject BasicType) currency int32
[___3] (FieldObject FieldType) Inventory (*digota.Inventory struct{Quantity int64, Type int32})
[____4] (PointerObject PointerType) inventory (*digota.Inventory struct{Quantity int64, Type int32})
[_____5] (StructObject UserType) digota.Inventory struct{Quantity int64, Type int32}
[___3] (FieldObject FieldType) Metadata map[string]string
[____4] (MapObject MapType) metadata map[string]string
[___3] (FieldObject FieldType) Name string
[____4] (BasicObject BasicType) name string
[___3] (FieldObject FieldType) PackageDimensions (*digota.PackageDimensions struct{Height float64, Length float64, Weight float64, Width float64})
[____4] (PointerObject PointerType) packageDimensions (*digota.PackageDimensions struct{Height float64, Length float64, Weight float64, Width float64})
[_____5] (StructObject UserType) digota.PackageDimensions struct{Height float64, Length float64, Weight float64, Width float64}
[___3] (FieldObject FieldType) Parent string
[____4] (BasicObject BasicType) parent string
[___3] (FieldObject FieldType) Price uint64
[____4] (BasicObject BasicType) price uint64

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = skus, collection = skus}

[0] (InterfaceObject UserType) err .error

[0] (InterfaceObject UserType) err .error

