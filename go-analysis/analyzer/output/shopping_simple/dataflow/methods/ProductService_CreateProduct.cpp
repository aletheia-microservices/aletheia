[0] (PointerObject PointerType) s (*shopping_simple.ProductServiceImpl struct{product_db NoSQLDatabase, product_queue Queue, num_workers int})
[_1] (StructObject UserType) shopping_simple.ProductServiceImpl struct{product_db NoSQLDatabase, product_queue Queue, num_workers int}
[__2] (FieldObject FieldType) product_db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) product_db NoSQLDatabase
[__2] (FieldObject FieldType) product_queue Queue
[___3] (BlueprintBackendObject BlueprintBackendType) product_queue Queue

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ Frontend

    --> w-tainted: write(product_db.Product.ProductID) {1}
[0] (BasicObject BasicType) productID string
     --> w-tainted: write(product_db.Product.ProductID) {1}
[_1] (Reference BasicType) ref <productID string> @ Frontend

    --> w-tainted: write(product_db.Product.Description) {1}
[0] (BasicObject BasicType) description string
     --> w-tainted: write(product_db.Product.Description) {1}
[_1] (Reference BasicType) ref <description string> @ Frontend

    --> w-tainted: write(product_db.Product.PricePerUnit) {1}
[0] (BasicObject BasicType) pricePerUnit int
     --> w-tainted: write(product_db.Product.PricePerUnit) {1}
[_1] (Reference BasicType) ref <pricePerUnit int> @ Frontend

    --> w-tainted: write(product_db.Product.Category) {1}
[0] (BasicObject BasicType) category string
     --> w-tainted: write(product_db.Product.Category) {1}
[_1] (Reference BasicType) ref <category string> @ Frontend

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = product_database, collection = product_collection}

[0] (InterfaceObject UserType) _ .error

    --> w-tainted: write(product_db.Product) {1}
[0] (StructObject UserType) product shopping_simple.Product struct{ProductID string, Description string, PricePerUnit int, Category string}
     --> w-tainted: write(product_db.Product.Category) {1}
[_1] (FieldObject FieldType) Category string
      --> w-tainted: write(product_db.Product.Category) {1}
[__2] (BasicObject BasicType) category string
       --> w-tainted: write(product_db.Product.Category) {1}
[___3] (Reference BasicType) ref <category string> @ Frontend
     --> w-tainted: write(product_db.Product.Description) {1}
[_1] (FieldObject FieldType) Description string
      --> w-tainted: write(product_db.Product.Description) {1}
[__2] (BasicObject BasicType) description string
       --> w-tainted: write(product_db.Product.Description) {1}
[___3] (Reference BasicType) ref <description string> @ Frontend
     --> w-tainted: write(product_db.Product.PricePerUnit) {1}
[_1] (FieldObject FieldType) PricePerUnit int
      --> w-tainted: write(product_db.Product.PricePerUnit) {1}
[__2] (BasicObject BasicType) pricePerUnit int
       --> w-tainted: write(product_db.Product.PricePerUnit) {1}
[___3] (Reference BasicType) ref <pricePerUnit int> @ Frontend
     --> w-tainted: write(product_db.Product.ProductID) {1}
[_1] (FieldObject FieldType) ProductID string
      --> w-tainted: write(product_db.Product.ProductID) {1}
[__2] (BasicObject BasicType) productID string
       --> w-tainted: write(product_db.Product.ProductID) {1}
[___3] (Reference BasicType) ref <productID string> @ Frontend

[0] (InterfaceObject UserType) err .error

