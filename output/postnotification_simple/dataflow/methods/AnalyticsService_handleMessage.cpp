[0] (PointerObject PointerType) a (*postnotification_simple.AnalyticsServiceImpl struct{analyticsQueue Queue, analyticsDb NoSQLDatabase, numWorkers int})
[_1] (StructObject UserType) postnotification_simple.AnalyticsServiceImpl struct{analyticsQueue Queue, analyticsDb NoSQLDatabase, numWorkers 4 int}
[__2] (FieldObject FieldType) analyticsDb NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) analyticsDb NoSQLDatabase
[__2] (FieldObject FieldType) analyticsQueue Queue
[___3] (BlueprintBackendObject BlueprintBackendType) analyticsQueue Queue
[__2] (FieldObject FieldType) numWorkers 4 int
[___3] (BasicObject BasicType) 4 int

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ AnalyticsService

    --> w-tainted: write(analytics_queue.TriggerAnalyticsMessage) {1}       --> w-tainted: write(analytics_queue.TriggerAnalyticsMessage) {1} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage) {1}
[0] (StructObject UserType) message postnotification_simple.TriggerAnalyticsMessage struct{PostID string}
     --> w-tainted: write(analytics_queue.TriggerAnalyticsMessage) {1}         --> w-tainted: write(analytics_queue.TriggerAnalyticsMessage) {1} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage) {1}
[_1] (Reference UserType) ref <analyticsEvent postnotification_simple.TriggerAnalyticsMessage struct{PostID string}> @ AnalyticsService
      --> w-tainted: write(analytics_queue.TriggerAnalyticsMessage) {1}           --> w-tainted: write(analytics_queue.TriggerAnalyticsMessage) {1} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage) {1}
[__2] (Reference UserType) ref <message postnotification_simple.TriggerAnalyticsMessage struct{PostID string}> @ AnalyticsService
       --> w-tainted: write(analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID) {2}             --> w-tainted: write(analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID) {2} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage.PostID) {1}
[___3] (FieldObject FieldType) PostID string
        --> w-tainted: write(analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID) {2}               --> w-tainted: write(analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID) {2} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage.PostID) {1}
[____4] (BasicObject BasicType) string
         --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID_MESSAGE) {4}                 --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID_MESSAGE) {4} --> r-tainted: read(posts_db._.postid, analytics_queue.TriggerAnalyticsMessage.PostID, notifications_queue.Message.PostID_MESSAGE, posts_db.Post.PostID) {4}
[_____5] (BasicObject BasicType) postID_STORAGE_SVC int64
         --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID_MESSAGE) {4}                 --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID_MESSAGE) {4} --> r-tainted: read(posts_db._.postid, analytics_queue.TriggerAnalyticsMessage.PostID, notifications_queue.Message.PostID_MESSAGE, posts_db.Post.PostID) {4}
[_____5] (BasicObject BasicType) postID_STORAGE_SVC int64
     --> w-tainted: write(analytics_queue.TriggerAnalyticsMessage) {1}         --> w-tainted: write(analytics_queue.TriggerAnalyticsMessage) {1} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage) {1}
[_1] (Reference UserType) ref <message postnotification_simple.TriggerAnalyticsMessage struct{PostID string}> @ AnalyticsService
      --> w-tainted: write(analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID) {2}           --> w-tainted: write(analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID) {2} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage.PostID) {1}
[__2] (FieldObject FieldType) PostID string
       --> w-tainted: write(analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID) {2}             --> w-tainted: write(analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID) {2} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage.PostID) {1}
[___3] (BasicObject BasicType) string
        --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID_MESSAGE) {4}               --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID_MESSAGE) {4} --> r-tainted: read(posts_db._.postid, analytics_queue.TriggerAnalyticsMessage.PostID, notifications_queue.Message.PostID_MESSAGE, posts_db.Post.PostID) {4}
[____4] (BasicObject BasicType) postID_STORAGE_SVC int64
        --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID_MESSAGE) {4}               --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID_MESSAGE) {4} --> r-tainted: read(posts_db._.postid, analytics_queue.TriggerAnalyticsMessage.PostID, notifications_queue.Message.PostID_MESSAGE, posts_db.Post.PostID) {4}
[____4] (BasicObject BasicType) postID_STORAGE_SVC int64
     --> w-tainted: write(analytics_queue.TriggerAnalyticsMessage.PostID) {1}
[_1] (FieldObject FieldType) PostID string
      --> w-tainted: write(analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID) {2}           --> w-tainted: write(analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID) {2} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage.PostID) {1}
[__2] (Reference FieldType) ref <PostID string> @ AnalyticsService
       --> w-tainted: write(analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID) {2}             --> w-tainted: write(analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID) {2} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage.PostID) {1}
[___3] (BasicObject BasicType) string
        --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID_MESSAGE) {4}               --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID_MESSAGE) {4} --> r-tainted: read(posts_db._.postid, analytics_queue.TriggerAnalyticsMessage.PostID, notifications_queue.Message.PostID_MESSAGE, posts_db.Post.PostID) {4}
[____4] (BasicObject BasicType) postID_STORAGE_SVC int64
        --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID_MESSAGE) {4}               --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID_MESSAGE) {4} --> r-tainted: read(posts_db._.postid, analytics_queue.TriggerAnalyticsMessage.PostID, notifications_queue.Message.PostID_MESSAGE, posts_db.Post.PostID) {4}
[____4] (BasicObject BasicType) postID_STORAGE_SVC int64
      --> w-tainted: write(analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID) {2}
[__2] (BasicObject BasicType) PostID string
       --> w-tainted: write(analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID) {2}             --> w-tainted: write(analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID) {2} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage.PostID) {1}
[___3] (Reference FieldType) ref <PostID string> @ AnalyticsService
        --> w-tainted: write(analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID) {2}               --> w-tainted: write(analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID) {2} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage.PostID) {1}
[____4] (BasicObject BasicType) string
         --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID_MESSAGE) {4}                 --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID_MESSAGE) {4} --> r-tainted: read(posts_db._.postid, analytics_queue.TriggerAnalyticsMessage.PostID, notifications_queue.Message.PostID_MESSAGE, posts_db.Post.PostID) {4}
[_____5] (BasicObject BasicType) postID_STORAGE_SVC int64
         --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID_MESSAGE) {4}                 --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID_MESSAGE) {4} --> r-tainted: read(posts_db._.postid, analytics_queue.TriggerAnalyticsMessage.PostID, notifications_queue.Message.PostID_MESSAGE, posts_db.Post.PostID) {4}
[_____5] (BasicObject BasicType) postID_STORAGE_SVC int64

    --> w-tainted: write(analytics_db.Analytics.PostID) {1}
[0] (BasicObject BasicType) postID int64
     --> w-tainted: write(analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID) {2}
[_1] (BasicObject BasicType) PostID string
      --> w-tainted: write(analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID) {2}           --> w-tainted: write(analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID) {2} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage.PostID) {1}
[__2] (Reference FieldType) ref <PostID string> @ AnalyticsService
       --> w-tainted: write(analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID) {2}             --> w-tainted: write(analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID) {2} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage.PostID) {1}
[___3] (BasicObject BasicType) string
        --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID_MESSAGE) {4}               --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID_MESSAGE) {4} --> r-tainted: read(posts_db._.postid, analytics_queue.TriggerAnalyticsMessage.PostID, notifications_queue.Message.PostID_MESSAGE, posts_db.Post.PostID) {4}
[____4] (BasicObject BasicType) postID_STORAGE_SVC int64
        --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID_MESSAGE) {4}               --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID_MESSAGE) {4} --> r-tainted: read(posts_db._.postid, analytics_queue.TriggerAnalyticsMessage.PostID, notifications_queue.Message.PostID_MESSAGE, posts_db.Post.PostID) {4}
[____4] (BasicObject BasicType) postID_STORAGE_SVC int64
     --> w-tainted: write(analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID) {2}
[_1] (BasicObject BasicType) PostID string
      --> w-tainted: write(analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID) {2}           --> w-tainted: write(analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID) {2} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage.PostID) {1}
[__2] (Reference FieldType) ref <PostID string> @ AnalyticsService
       --> w-tainted: write(analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID) {2}             --> w-tainted: write(analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID) {2} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage.PostID) {1}
[___3] (BasicObject BasicType) string
        --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID_MESSAGE) {4}               --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID_MESSAGE) {4} --> r-tainted: read(posts_db._.postid, analytics_queue.TriggerAnalyticsMessage.PostID, notifications_queue.Message.PostID_MESSAGE, posts_db.Post.PostID) {4}
[____4] (BasicObject BasicType) postID_STORAGE_SVC int64
        --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID_MESSAGE) {4}               --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID_MESSAGE) {4} --> r-tainted: read(posts_db._.postid, analytics_queue.TriggerAnalyticsMessage.PostID, notifications_queue.Message.PostID_MESSAGE, posts_db.Post.PostID) {4}
[____4] (BasicObject BasicType) postID_STORAGE_SVC int64

[0] (InterfaceObject UserType) err .error
     --> w-tainted: write(analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID) {2}
[_1] (BasicObject BasicType) PostID string
      --> w-tainted: write(analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID) {2}           --> w-tainted: write(analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID) {2} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage.PostID) {1}
[__2] (Reference FieldType) ref <PostID string> @ AnalyticsService
       --> w-tainted: write(analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID) {2}             --> w-tainted: write(analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID) {2} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage.PostID) {1}
[___3] (BasicObject BasicType) string
        --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID_MESSAGE) {4}               --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID_MESSAGE) {4} --> r-tainted: read(posts_db._.postid, analytics_queue.TriggerAnalyticsMessage.PostID, notifications_queue.Message.PostID_MESSAGE, posts_db.Post.PostID) {4}
[____4] (BasicObject BasicType) postID_STORAGE_SVC int64
        --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID_MESSAGE) {4}               --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID_MESSAGE) {4} --> r-tainted: read(posts_db._.postid, analytics_queue.TriggerAnalyticsMessage.PostID, notifications_queue.Message.PostID_MESSAGE, posts_db.Post.PostID) {4}
[____4] (BasicObject BasicType) postID_STORAGE_SVC int64

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = analyticsDb, collection = analytics_collection}

[0] (InterfaceObject UserType) err .error

    --> w-tainted: write(analytics_db.Analytics) {1}
[0] (StructObject UserType) analytics postnotification_simple.Analytics struct{PostID int64}
     --> w-tainted: write(analytics_db.Analytics.PostID) {1}
[_1] (FieldObject FieldType) PostID int64
      --> w-tainted: write(analytics_db.Analytics.PostID) {1}
[__2] (BasicObject BasicType) postID int64
       --> w-tainted: write(analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID) {2}
[___3] (BasicObject BasicType) PostID string
        --> w-tainted: write(analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID) {2}               --> w-tainted: write(analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID) {2} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage.PostID) {1}
[____4] (Reference FieldType) ref <PostID string> @ AnalyticsService
         --> w-tainted: write(analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID) {2}                 --> w-tainted: write(analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID) {2} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage.PostID) {1}
[_____5] (BasicObject BasicType) string
          --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID_MESSAGE) {4}                   --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID_MESSAGE) {4} --> r-tainted: read(posts_db._.postid, analytics_queue.TriggerAnalyticsMessage.PostID, notifications_queue.Message.PostID_MESSAGE, posts_db.Post.PostID) {4}
[______6] (BasicObject BasicType) postID_STORAGE_SVC int64
          --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID_MESSAGE) {4}                   --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID_MESSAGE) {4} --> r-tainted: read(posts_db._.postid, analytics_queue.TriggerAnalyticsMessage.PostID, notifications_queue.Message.PostID_MESSAGE, posts_db.Post.PostID) {4}
[______6] (BasicObject BasicType) postID_STORAGE_SVC int64
       --> w-tainted: write(analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID) {2}
[___3] (BasicObject BasicType) PostID string
        --> w-tainted: write(analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID) {2}               --> w-tainted: write(analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID) {2} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage.PostID) {1}
[____4] (Reference FieldType) ref <PostID string> @ AnalyticsService
         --> w-tainted: write(analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID) {2}                 --> w-tainted: write(analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID) {2} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage.PostID) {1}
[_____5] (BasicObject BasicType) string
          --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID_MESSAGE) {4}                   --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID_MESSAGE) {4} --> r-tainted: read(posts_db._.postid, analytics_queue.TriggerAnalyticsMessage.PostID, notifications_queue.Message.PostID_MESSAGE, posts_db.Post.PostID) {4}
[______6] (BasicObject BasicType) postID_STORAGE_SVC int64
          --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID_MESSAGE) {4}                   --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID_MESSAGE) {4} --> r-tainted: read(posts_db._.postid, analytics_queue.TriggerAnalyticsMessage.PostID, notifications_queue.Message.PostID_MESSAGE, posts_db.Post.PostID) {4}
[______6] (BasicObject BasicType) postID_STORAGE_SVC int64

