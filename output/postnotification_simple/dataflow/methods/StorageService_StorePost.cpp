[0] (PointerObject PointerType) s (*postnotification_simple.StorageServiceImpl struct{postsDb NoSQLDatabase, analyticsQueue Queue})
[_1] (StructObject UserType) postnotification_simple.StorageServiceImpl struct{postsDb NoSQLDatabase, analyticsQueue Queue}
[__2] (FieldObject FieldType) analyticsQueue Queue
[___3] (BlueprintBackendObject BlueprintBackendType) analyticsQueue Queue
[__2] (FieldObject FieldType) postsDb NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) postsDb NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ UploadService

    --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID) {2}
[0] (BasicObject BasicType) reqID int64
     --> w-tainted: write(notifications_queue.Message.ReqID, posts_db.Post.ReqID) {2}         --> w-tainted: write(notifications_queue.Message.ReqID, posts_db.Post.ReqID) {2} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[_1] (Reference BasicType) ref <reqID int64> @ UploadService

    --> w-tainted: write(posts_db.Post.Text) {1}
[0] (BasicObject BasicType) text string
     --> w-tainted: write(posts_db.Post.Text) {1}
[_1] (Reference BasicType) ref <text string> @ UploadService

    --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID_MESSAGE) {4}       --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID_MESSAGE) {4} --> r-tainted: read(posts_db._.postid, analytics_queue.TriggerAnalyticsMessage.PostID, notifications_queue.Message.PostID_MESSAGE, posts_db.Post.PostID) {4}
[0] (BasicObject BasicType) postID_STORAGE_SVC int64

    --> w-tainted: write(posts_db.Post.Timestamp) {1}
[0] (BasicObject BasicType) timestamp int64

    --> w-tainted: write(posts_db.Post.Mentions) {1}
[0] (ArrayObject ArrayType) mentions []string
     --> w-tainted: write(posts_db.Post.Mentions) {1}
[_1] (BasicObject BasicType) "alice" string
     --> w-tainted: write(posts_db.Post.Mentions) {1}
[_1] (BasicObject BasicType) "bob" string

    --> w-tainted: write(posts_db.Post) {1}
[0] (StructObject UserType) post postnotification_simple.Post struct{ReqID int64, PostID int64, MediaID int64, Text string, Mentions []string, Timestamp int64, Creator postnotification_simple.Creator struct{Username "some username" string}}
     --> w-tainted: write(posts_db.Post.Creator) {1}
[_1] (FieldObject FieldType) Creator postnotification_simple.Creator struct{Username "some username" string}
      --> w-tainted: write(posts_db.Post.Creator) {1}
[__2] (StructObject UserType) postnotification_simple.Creator struct{Username "some username" string}
       --> w-tainted: write(posts_db.Post.Creator.Username) {1}
[___3] (FieldObject FieldType) Username "some username" string
        --> w-tainted: write(posts_db.Post.Creator.Username) {1}
[____4] (BasicObject BasicType) "some username" string
     --> w-tainted: write(posts_db.Post.Mentions) {1}
[_1] (FieldObject FieldType) Mentions []string
      --> w-tainted: write(posts_db.Post.Mentions) {1}
[__2] (ArrayObject ArrayType) mentions []string
       --> w-tainted: write(posts_db.Post.Mentions) {1}
[___3] (BasicObject BasicType) "alice" string
       --> w-tainted: write(posts_db.Post.Mentions) {1}
[___3] (BasicObject BasicType) "bob" string
     --> w-tainted: write(posts_db.Post.PostID) {1}
[_1] (FieldObject FieldType) PostID int64
      --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID_MESSAGE) {4}           --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID_MESSAGE) {4} --> r-tainted: read(posts_db._.postid, analytics_queue.TriggerAnalyticsMessage.PostID, notifications_queue.Message.PostID_MESSAGE, posts_db.Post.PostID) {4}
[__2] (BasicObject BasicType) postID_STORAGE_SVC int64
     --> w-tainted: write(posts_db.Post.ReqID) {1}
[_1] (FieldObject FieldType) ReqID int64
      --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID) {2}
[__2] (BasicObject BasicType) reqID int64
       --> w-tainted: write(notifications_queue.Message.ReqID, posts_db.Post.ReqID) {2}             --> w-tainted: write(notifications_queue.Message.ReqID, posts_db.Post.ReqID) {2} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[___3] (Reference BasicType) ref <reqID int64> @ UploadService
     --> w-tainted: write(posts_db.Post.Text) {1}
[_1] (FieldObject FieldType) Text string
      --> w-tainted: write(posts_db.Post.Text) {1}
[__2] (BasicObject BasicType) text string
       --> w-tainted: write(posts_db.Post.Text) {1}
[___3] (Reference BasicType) ref <text string> @ UploadService
     --> w-tainted: write(posts_db.Post.Timestamp) {1}
[_1] (FieldObject FieldType) Timestamp int64
      --> w-tainted: write(posts_db.Post.Timestamp) {1}
[__2] (BasicObject BasicType) timestamp int64

[0] (BasicObject BasicType) myval 0 int

[0] (ArrayObject ArrayType) mymentions []string
[_1] (BasicObject BasicType) mention string
      --> w-tainted: write(posts_db.Post.Mentions) {1}
[__2] (ArrayObject ArrayType) mentions []string
       --> w-tainted: write(posts_db.Post.Mentions) {1}
[___3] (BasicObject BasicType) "alice" string
       --> w-tainted: write(posts_db.Post.Mentions) {1}
[___3] (BasicObject BasicType) "bob" string

[0] (BasicObject BasicType) idx 0 int

[0] (BasicObject BasicType) mention string
     --> w-tainted: write(posts_db.Post.Mentions) {1}
[_1] (ArrayObject ArrayType) mentions []string
      --> w-tainted: write(posts_db.Post.Mentions) {1}
[__2] (BasicObject BasicType) "alice" string
      --> w-tainted: write(posts_db.Post.Mentions) {1}
[__2] (BasicObject BasicType) "bob" string

[0] (ArrayObject ArrayType) mymentions []string
[_1] (BasicObject BasicType) mention string
      --> w-tainted: write(posts_db.Post.Mentions) {1}
[__2] (ArrayObject ArrayType) mentions []string
       --> w-tainted: write(posts_db.Post.Mentions) {1}
[___3] (BasicObject BasicType) "alice" string
       --> w-tainted: write(posts_db.Post.Mentions) {1}
[___3] (BasicObject BasicType) "bob" string

[0] (PointerObject PointerType) s (*postnotification_simple.StorageServiceImpl struct{postsDb NoSQLDatabase, analyticsQueue Queue})
[_1] (StructObject UserType) postnotification_simple.StorageServiceImpl struct{postsDb NoSQLDatabase, analyticsQueue Queue}
[__2] (FieldObject FieldType) analyticsQueue Queue
[___3] (BlueprintBackendObject BlueprintBackendType) analyticsQueue Queue
[__2] (FieldObject FieldType) postsDb NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) postsDb NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ UploadService

    --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID) {2}
[0] (BasicObject BasicType) reqID int64
     --> w-tainted: write(notifications_queue.Message.ReqID, posts_db.Post.ReqID) {2}         --> w-tainted: write(notifications_queue.Message.ReqID, posts_db.Post.ReqID) {2} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[_1] (Reference BasicType) ref <reqID int64> @ UploadService

    --> w-tainted: write(posts_db.Post.Text) {1}
[0] (BasicObject BasicType) text string
     --> w-tainted: write(posts_db.Post.Text) {1}
[_1] (Reference BasicType) ref <text string> @ UploadService

    --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID_MESSAGE) {4}       --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID_MESSAGE) {4} --> r-tainted: read(posts_db._.postid, analytics_queue.TriggerAnalyticsMessage.PostID, notifications_queue.Message.PostID_MESSAGE, posts_db.Post.PostID) {4}
[0] (BasicObject BasicType) postID_STORAGE_SVC int64

    --> w-tainted: write(posts_db.Post.Timestamp) {1}
[0] (BasicObject BasicType) timestamp int64

    --> w-tainted: write(posts_db.Post.Mentions) {1}
[0] (ArrayObject ArrayType) mentions []string
     --> w-tainted: write(posts_db.Post.Mentions) {1}
[_1] (BasicObject BasicType) "alice" string
     --> w-tainted: write(posts_db.Post.Mentions) {1}
[_1] (BasicObject BasicType) "bob" string

    --> w-tainted: write(posts_db.Post) {1}
[0] (StructObject UserType) post postnotification_simple.Post struct{ReqID int64, PostID int64, MediaID int64, Text string, Mentions []string, Timestamp int64, Creator postnotification_simple.Creator struct{Username "some username" string}}
     --> w-tainted: write(posts_db.Post.Creator) {1}
[_1] (FieldObject FieldType) Creator postnotification_simple.Creator struct{Username "some username" string}
      --> w-tainted: write(posts_db.Post.Creator) {1}
[__2] (StructObject UserType) postnotification_simple.Creator struct{Username "some username" string}
       --> w-tainted: write(posts_db.Post.Creator.Username) {1}
[___3] (FieldObject FieldType) Username "some username" string
        --> w-tainted: write(posts_db.Post.Creator.Username) {1}
[____4] (BasicObject BasicType) "some username" string
     --> w-tainted: write(posts_db.Post.Mentions) {1}
[_1] (FieldObject FieldType) Mentions []string
      --> w-tainted: write(posts_db.Post.Mentions) {1}
[__2] (ArrayObject ArrayType) mentions []string
       --> w-tainted: write(posts_db.Post.Mentions) {1}
[___3] (BasicObject BasicType) "alice" string
       --> w-tainted: write(posts_db.Post.Mentions) {1}
[___3] (BasicObject BasicType) "bob" string
     --> w-tainted: write(posts_db.Post.PostID) {1}
[_1] (FieldObject FieldType) PostID int64
      --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID_MESSAGE) {4}           --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID_MESSAGE) {4} --> r-tainted: read(posts_db._.postid, analytics_queue.TriggerAnalyticsMessage.PostID, notifications_queue.Message.PostID_MESSAGE, posts_db.Post.PostID) {4}
[__2] (BasicObject BasicType) postID_STORAGE_SVC int64
     --> w-tainted: write(posts_db.Post.ReqID) {1}
[_1] (FieldObject FieldType) ReqID int64
      --> w-tainted: write(posts_db.Post.ReqID, notifications_queue.Message.ReqID) {2}
[__2] (BasicObject BasicType) reqID int64
       --> w-tainted: write(notifications_queue.Message.ReqID, posts_db.Post.ReqID) {2}             --> w-tainted: write(notifications_queue.Message.ReqID, posts_db.Post.ReqID) {2} --> r-tainted: read(notifications_queue.Message.ReqID) {1}
[___3] (Reference BasicType) ref <reqID int64> @ UploadService
     --> w-tainted: write(posts_db.Post.Text) {1}
[_1] (FieldObject FieldType) Text string
      --> w-tainted: write(posts_db.Post.Text) {1}
[__2] (BasicObject BasicType) text string
       --> w-tainted: write(posts_db.Post.Text) {1}
[___3] (Reference BasicType) ref <text string> @ UploadService
     --> w-tainted: write(posts_db.Post.Timestamp) {1}
[_1] (FieldObject FieldType) Timestamp int64
      --> w-tainted: write(posts_db.Post.Timestamp) {1}
[__2] (BasicObject BasicType) timestamp int64

[0] (BasicObject BasicType) myval 0 int

[0] (ArrayObject ArrayType) mymentions []string
[_1] (BasicObject BasicType) mention string
      --> w-tainted: write(posts_db.Post.Mentions) {1}
[__2] (ArrayObject ArrayType) mentions []string
       --> w-tainted: write(posts_db.Post.Mentions) {1}
[___3] (BasicObject BasicType) "alice" string
       --> w-tainted: write(posts_db.Post.Mentions) {1}
[___3] (BasicObject BasicType) "bob" string

[0] (BasicObject BasicType) idx 0 int

[0] (BasicObject BasicType) mention string
     --> w-tainted: write(posts_db.Post.Mentions) {1}
[_1] (ArrayObject ArrayType) mentions []string
      --> w-tainted: write(posts_db.Post.Mentions) {1}
[__2] (BasicObject BasicType) "alice" string
      --> w-tainted: write(posts_db.Post.Mentions) {1}
[__2] (BasicObject BasicType) "bob" string

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = post, collection = post}

[0] (InterfaceObject UserType) err .error

[0] (InterfaceObject UserType) err .error

    --> w-tainted: write(analytics_queue.TriggerAnalyticsMessage) {1}       --> w-tainted: write(analytics_queue.TriggerAnalyticsMessage) {1} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage) {1}
[0] (StructObject UserType) message postnotification_simple.TriggerAnalyticsMessage struct{PostID string}
     --> w-tainted: write(analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID) {2}         --> w-tainted: write(analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID) {2} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage.PostID) {1}
[_1] (FieldObject FieldType) PostID string
      --> w-tainted: write(analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID) {2}           --> w-tainted: write(analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID) {2} --> r-tainted: read(analytics_queue.TriggerAnalyticsMessage.PostID) {1}
[__2] (BasicObject BasicType) string
       --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID_MESSAGE) {4}             --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID_MESSAGE) {4} --> r-tainted: read(posts_db._.postid, analytics_queue.TriggerAnalyticsMessage.PostID, notifications_queue.Message.PostID_MESSAGE, posts_db.Post.PostID) {4}
[___3] (BasicObject BasicType) postID_STORAGE_SVC int64
       --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID_MESSAGE) {4}             --> w-tainted: write(posts_db.Post.PostID, analytics_queue.TriggerAnalyticsMessage.PostID, analytics_db.Analytics.PostID, notifications_queue.Message.PostID_MESSAGE) {4} --> r-tainted: read(posts_db._.postid, analytics_queue.TriggerAnalyticsMessage.PostID, notifications_queue.Message.PostID_MESSAGE, posts_db.Post.PostID) {4}
[___3] (BasicObject BasicType) postID_STORAGE_SVC int64

[0] (BasicObject BasicType) _ bool

[0] (InterfaceObject UserType) err .error

