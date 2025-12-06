[0] (PointerObject PointerType) s (*shopping_app.AnalyticsServiceImpl struct{analytics_db NoSQLDatabase, analytics_queue Queue, num_workers int})
[_1] (StructObject UserType) shopping_app.AnalyticsServiceImpl struct{analytics_db NoSQLDatabase, analytics_queue Queue, num_workers 4 int}
[__2] (FieldObject FieldType) analytics_db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) analytics_db NoSQLDatabase
[__2] (FieldObject FieldType) analytics_queue Queue
[___3] (BlueprintBackendObject BlueprintBackendType) analytics_queue Queue
[__2] (FieldObject FieldType) num_workers 4 int
[___3] (BasicObject BasicType) 4 int

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ AnalyticsService

[0] (StructObject UserType) message shopping_app.AnalyticsMessage struct{UserID string, ProductCategory string}
[_1] (Reference UserType) ref <workerMessage shopping_app.AnalyticsMessage struct{UserID string, ProductCategory string}> @ AnalyticsService
      --> w-tainted: write(analytics_db.Analytics.Categories) {1}           --> w-tainted: write(analytics_db.Analytics.Categories) {1} --> r-tainted: read(analytics_db.Analytics.ProductCategory) {1}
[__2] (FieldObject FieldType) ProductCategory string
       --> w-tainted: write(analytics_db.Analytics.Categories) {1}             --> w-tainted: write(analytics_db.Analytics.Categories) {1} --> r-tainted: read(analytics_db.Analytics.ProductCategory) {1}
[___3] (BasicObject BasicType) * string
        --> w-tainted: write(analytics_db.Analytics.Categories) {1}               --> w-tainted: write(analytics_db.Analytics.Categories) {1} --> r-tainted: read(analytics_db.Analytics.interface{}) {1}
[____4] (InterfaceObject InterfaceType) * interface{}
      --> w-tainted: write(analytics_db.Analytics.UserID) {1}           --> w-tainted: write(analytics_db.Analytics.UserID) {1} --> r-tainted: read(analytics_db._.userID) {1}
[__2] (FieldObject FieldType) UserID string
       --> w-tainted: write(analytics_db.Analytics.UserID) {1}             --> w-tainted: write(analytics_db.Analytics.UserID) {1} --> r-tainted: read(analytics_db._.userID) {1}
[___3] (BasicObject BasicType) * string
        --> w-tainted: write(analytics_db.Analytics.UserID) {1}               --> w-tainted: write(analytics_db.Analytics.UserID) {1} --> r-tainted: read(analytics_db._.userID) {1}
[____4] (InterfaceObject InterfaceType) * interface{}
[_1] (FieldObject FieldType) ProductCategory string
      --> w-tainted: write(analytics_db.Analytics.Categories) {1}           --> w-tainted: write(analytics_db.Analytics.Categories) {1} --> r-tainted: read(analytics_db.Analytics.ProductCategory) {1}
[__2] (Reference FieldType) ref <ProductCategory string> @ AnalyticsService
       --> w-tainted: write(analytics_db.Analytics.Categories) {1}             --> w-tainted: write(analytics_db.Analytics.Categories) {1} --> r-tainted: read(analytics_db.Analytics.ProductCategory) {1}
[___3] (BasicObject BasicType) * string
        --> w-tainted: write(analytics_db.Analytics.Categories) {1}               --> w-tainted: write(analytics_db.Analytics.Categories) {1} --> r-tainted: read(analytics_db.Analytics.interface{}) {1}
[____4] (InterfaceObject InterfaceType) * interface{}
      --> w-tainted: write(analytics_db.Analytics.Categories) {1}           --> w-tainted: write(analytics_db.Analytics.Categories) {1} --> r-tainted: read(analytics_db.Analytics.string) {1}
[__2] (BasicObject BasicType) ProductCategory string
       --> w-tainted: write(analytics_db.Analytics.Categories) {1}             --> w-tainted: write(analytics_db.Analytics.Categories) {1} --> r-tainted: read(analytics_db.Analytics.ProductCategory) {1}
[___3] (Reference FieldType) ref <ProductCategory string> @ AnalyticsService
        --> w-tainted: write(analytics_db.Analytics.Categories) {1}               --> w-tainted: write(analytics_db.Analytics.Categories) {1} --> r-tainted: read(analytics_db.Analytics.ProductCategory) {1}
[____4] (BasicObject BasicType) * string
         --> w-tainted: write(analytics_db.Analytics.Categories) {1}                 --> w-tainted: write(analytics_db.Analytics.Categories) {1} --> r-tainted: read(analytics_db.Analytics.interface{}) {1}
[_____5] (InterfaceObject InterfaceType) * interface{}
[_1] (FieldObject FieldType) UserID string
      --> w-tainted: write(analytics_db.Analytics.UserID) {1}           --> w-tainted: write(analytics_db.Analytics.UserID) {1} --> r-tainted: read(analytics_db._.userID) {1}
[__2] (Reference FieldType) ref <UserID string> @ AnalyticsService
       --> w-tainted: write(analytics_db.Analytics.UserID) {1}             --> w-tainted: write(analytics_db.Analytics.UserID) {1} --> r-tainted: read(analytics_db._.userID) {1}
[___3] (BasicObject BasicType) * string
        --> w-tainted: write(analytics_db.Analytics.UserID) {1}               --> w-tainted: write(analytics_db.Analytics.UserID) {1} --> r-tainted: read(analytics_db._.userID) {1}
[____4] (InterfaceObject InterfaceType) * interface{}
      --> w-tainted: write(analytics_db.Analytics.UserID) {1}           --> w-tainted: write(analytics_db.Analytics.UserID) {1} --> r-tainted: read(analytics_db._.userID) {1}
[__2] (BasicObject BasicType) UserID string
       --> w-tainted: write(analytics_db.Analytics.UserID) {1}             --> w-tainted: write(analytics_db.Analytics.UserID) {1} --> r-tainted: read(analytics_db._.userID) {1}
[___3] (Reference FieldType) ref <UserID string> @ AnalyticsService
        --> w-tainted: write(analytics_db.Analytics.UserID) {1}               --> w-tainted: write(analytics_db.Analytics.UserID) {1} --> r-tainted: read(analytics_db._.userID) {1}
[____4] (BasicObject BasicType) * string
         --> w-tainted: write(analytics_db.Analytics.UserID) {1}                 --> w-tainted: write(analytics_db.Analytics.UserID) {1} --> r-tainted: read(analytics_db._.userID) {1}
[_____5] (InterfaceObject InterfaceType) * interface{}

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = analytics_database, collection = analytics_collection}

[0] (InterfaceObject UserType) _ .error

    --> r-tainted: read(analytics_db.Analytics, analytics_db.Analytics.Categories, analytics_db.Analytics.string, analytics_db.Analytics.ProductCategory, analytics_db.Analytics.interface{}) {5}
[0] (StructObject UserType) analytics shopping_app.Analytics struct{UserID string, Categories []string}
     --> r-tainted: read(analytics_db.Analytics.Categories) {1}
[_1] (FieldObject FieldType) Categories []string
      --> w-tainted: write(analytics_db.Analytics.Categories) {1}           --> w-tainted: write(analytics_db.Analytics.Categories) {1} --> r-tainted: read(analytics_db.Analytics.Categories) {1}
[__2] (SliceObject SliceType) Categories []string
       --> w-tainted: write(analytics_db.Analytics.Categories) {1}             --> w-tainted: write(analytics_db.Analytics.Categories) {1} --> r-tainted: read(analytics_db.Analytics.string) {1}
[___3] (BasicObject BasicType) ProductCategory string
        --> w-tainted: write(analytics_db.Analytics.Categories) {1}               --> w-tainted: write(analytics_db.Analytics.Categories) {1} --> r-tainted: read(analytics_db.Analytics.ProductCategory) {1}
[____4] (Reference FieldType) ref <ProductCategory string> @ AnalyticsService
         --> w-tainted: write(analytics_db.Analytics.Categories) {1}                 --> w-tainted: write(analytics_db.Analytics.Categories) {1} --> r-tainted: read(analytics_db.Analytics.ProductCategory) {1}
[_____5] (BasicObject BasicType) * string
          --> w-tainted: write(analytics_db.Analytics.Categories) {1}                   --> w-tainted: write(analytics_db.Analytics.Categories) {1} --> r-tainted: read(analytics_db.Analytics.interface{}) {1}
[______6] (InterfaceObject InterfaceType) * interface{}

[0] (SliceObject UserType) query primitive.D
[_1] (StructObject StructType) struct{Key "userID" string, Key "userID" string, Value string, Value string}
[__2] (FieldObject FieldType) Key "userID" string
[___3] (BasicObject BasicType) "userID" string
[__2] (FieldObject FieldType) Value string
       --> w-tainted: write(analytics_db.Analytics.UserID) {1}             --> w-tainted: write(analytics_db.Analytics.UserID) {1} --> r-tainted: read(analytics_db._.userID) {1}
[___3] (BasicObject BasicType) UserID string
        --> w-tainted: write(analytics_db.Analytics.UserID) {1}               --> w-tainted: write(analytics_db.Analytics.UserID) {1} --> r-tainted: read(analytics_db._.userID) {1}
[____4] (Reference FieldType) ref <UserID string> @ AnalyticsService
         --> w-tainted: write(analytics_db.Analytics.UserID) {1}                 --> w-tainted: write(analytics_db.Analytics.UserID) {1} --> r-tainted: read(analytics_db._.userID) {1}
[_____5] (BasicObject BasicType) * string
          --> w-tainted: write(analytics_db.Analytics.UserID) {1}                   --> w-tainted: write(analytics_db.Analytics.UserID) {1} --> r-tainted: read(analytics_db._.userID) {1}
[______6] (InterfaceObject InterfaceType) * interface{}

    --> r-tainted: read(analytics_db._) {1}
[0] (BlueprintBackendObject BlueprintBackendType) result NoSQLCursor {database = analytics_database, collection = analytics_collection}
     --> r-tainted: read(analytics_db.Analytics, analytics_db.Analytics.Categories, analytics_db.Analytics.string, analytics_db.Analytics.ProductCategory, analytics_db.Analytics.interface{}) {5}
[_1] (StructObject UserType) analytics shopping_app.Analytics struct{UserID string, Categories []string}
      --> r-tainted: read(analytics_db.Analytics.Categories) {1}
[__2] (FieldObject FieldType) Categories []string
       --> w-tainted: write(analytics_db.Analytics.Categories) {1}             --> w-tainted: write(analytics_db.Analytics.Categories) {1} --> r-tainted: read(analytics_db.Analytics.Categories) {1}
[___3] (SliceObject SliceType) Categories []string
        --> w-tainted: write(analytics_db.Analytics.Categories) {1}               --> w-tainted: write(analytics_db.Analytics.Categories) {1} --> r-tainted: read(analytics_db.Analytics.string) {1}
[____4] (BasicObject BasicType) ProductCategory string
         --> w-tainted: write(analytics_db.Analytics.Categories) {1}                 --> w-tainted: write(analytics_db.Analytics.Categories) {1} --> r-tainted: read(analytics_db.Analytics.ProductCategory) {1}
[_____5] (Reference FieldType) ref <ProductCategory string> @ AnalyticsService
          --> w-tainted: write(analytics_db.Analytics.Categories) {1}                   --> w-tainted: write(analytics_db.Analytics.Categories) {1} --> r-tainted: read(analytics_db.Analytics.ProductCategory) {1}
[______6] (BasicObject BasicType) * string
           --> w-tainted: write(analytics_db.Analytics.Categories) {1}                     --> w-tainted: write(analytics_db.Analytics.Categories) {1} --> r-tainted: read(analytics_db.Analytics.interface{}) {1}
[_______7] (InterfaceObject InterfaceType) * interface{}

[0] (InterfaceObject UserType) _ .error

    --> w-tainted: write(analytics_db.Analytics) {1}
[0] (StructObject UserType) updatedAnalytics shopping_app.Analytics struct{UserID string, Categories []string}
     --> w-tainted: write(analytics_db.Analytics.Categories) {1}
[_1] (FieldObject FieldType) Categories []string
      --> w-tainted: write(analytics_db.Analytics.Categories) {1}           --> w-tainted: write(analytics_db.Analytics.Categories) {1} --> r-tainted: read(analytics_db.Analytics.Categories) {1}
[__2] (SliceObject SliceType) Categories []string
       --> w-tainted: write(analytics_db.Analytics.Categories) {1}             --> w-tainted: write(analytics_db.Analytics.Categories) {1} --> r-tainted: read(analytics_db.Analytics.string) {1}
[___3] (BasicObject BasicType) ProductCategory string
        --> w-tainted: write(analytics_db.Analytics.Categories) {1}               --> w-tainted: write(analytics_db.Analytics.Categories) {1} --> r-tainted: read(analytics_db.Analytics.ProductCategory) {1}
[____4] (Reference FieldType) ref <ProductCategory string> @ AnalyticsService
         --> w-tainted: write(analytics_db.Analytics.Categories) {1}                 --> w-tainted: write(analytics_db.Analytics.Categories) {1} --> r-tainted: read(analytics_db.Analytics.ProductCategory) {1}
[_____5] (BasicObject BasicType) * string
          --> w-tainted: write(analytics_db.Analytics.Categories) {1}                   --> w-tainted: write(analytics_db.Analytics.Categories) {1} --> r-tainted: read(analytics_db.Analytics.interface{}) {1}
[______6] (InterfaceObject InterfaceType) * interface{}
     --> w-tainted: write(analytics_db.Analytics.UserID) {1}
[_1] (FieldObject FieldType) UserID string
      --> w-tainted: write(analytics_db.Analytics.UserID) {1}           --> w-tainted: write(analytics_db.Analytics.UserID) {1} --> r-tainted: read(analytics_db._.userID) {1}
[__2] (BasicObject BasicType) UserID string
       --> w-tainted: write(analytics_db.Analytics.UserID) {1}             --> w-tainted: write(analytics_db.Analytics.UserID) {1} --> r-tainted: read(analytics_db._.userID) {1}
[___3] (Reference FieldType) ref <UserID string> @ AnalyticsService
        --> w-tainted: write(analytics_db.Analytics.UserID) {1}               --> w-tainted: write(analytics_db.Analytics.UserID) {1} --> r-tainted: read(analytics_db._.userID) {1}
[____4] (BasicObject BasicType) * string
         --> w-tainted: write(analytics_db.Analytics.UserID) {1}                 --> w-tainted: write(analytics_db.Analytics.UserID) {1} --> r-tainted: read(analytics_db._.userID) {1}
[_____5] (InterfaceObject InterfaceType) * interface{}

