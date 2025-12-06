[0] (PointerObject PointerType) c (*socialnetwork.ComposePostServiceImpl struct{postStorageService socialnetwork.PostStorageService, userTimelineService socialnetwork.UserTimelineService, userService socialnetwork.UserService, uniqueIDService socialnetwork.UniqueIdService, mediaService socialnetwork.MediaService, textService socialnetwork.TextService, homeTimelineService socialnetwork.HomeTimelineService})
[_1] (StructObject UserType) socialnetwork.ComposePostServiceImpl struct{postStorageService socialnetwork.PostStorageService, userTimelineService socialnetwork.UserTimelineService, userService socialnetwork.UserService, uniqueIDService socialnetwork.UniqueIdService, mediaService socialnetwork.MediaService, textService socialnetwork.TextService, homeTimelineService socialnetwork.HomeTimelineService}
[__2] (FieldObject FieldType) homeTimelineService socialnetwork.HomeTimelineService
[___3] (ServiceObject ServiceType) homeTimelineService socialnetwork.HomeTimelineService
[__2] (FieldObject FieldType) mediaService socialnetwork.MediaService
[___3] (ServiceObject ServiceType) mediaService socialnetwork.MediaService
[__2] (FieldObject FieldType) postStorageService socialnetwork.PostStorageService
[___3] (ServiceObject ServiceType) postStorageService socialnetwork.PostStorageService
[__2] (FieldObject FieldType) textService socialnetwork.TextService
[___3] (ServiceObject ServiceType) textService socialnetwork.TextService
[__2] (FieldObject FieldType) uniqueIDService socialnetwork.UniqueIdService
[___3] (ServiceObject ServiceType) uniqueIDService socialnetwork.UniqueIdService
[__2] (FieldObject FieldType) userService socialnetwork.UserService
[___3] (ServiceObject ServiceType) userService socialnetwork.UserService
[__2] (FieldObject FieldType) userTimelineService socialnetwork.UserTimelineService
[___3] (ServiceObject ServiceType) userTimelineService socialnetwork.UserTimelineService

[0] (InterfaceObject UserType) ctx context.Context

[0] (BasicObject BasicType) reqID int64

[0] (BasicObject BasicType) username string

[0] (BasicObject BasicType) userID int64

[0] (BasicObject BasicType) text string

[0] (ArrayObject ArrayType) mediaIDs []int64

[0] (ArrayObject ArrayType) mediaTypes []string

[0] (BasicObject BasicType) postType int64

[0] (BasicObject BasicType) timestamp int64

[0] (BasicObject BasicType) err1 error

[0] (BasicObject BasicType) err2 error

[0] (BasicObject BasicType) err3 error

[0] (BasicObject BasicType) err4 error

[0] (BasicObject BasicType) uniqueID int64

[0] (StructObject UserType) creator socialnetwork.Creator struct{UserID int64, Username string}

[0] (BasicObject BasicType) up_text string

[0] (ArrayObject SliceType) medias []socialnetwork.Media struct{MediaID int64, MediaType string}

[0] (ArrayObject SliceType) urls []socialnetwork.URL struct{ShortenedUrl string, ExpandedUrl string}

[0] (ArrayObject SliceType) usermentions []socialnetwork.UserMention struct{UserID int64, Username string}

[0] (StructObject UserType) wg sync.WaitGroup

[0] (StructObject UserType) post socialnetwork.Post struct{PostID int64, Creator socialnetwork.Creator struct{UserID int64, Username string}, ReqID int64, Text string, UserMentions []socialnetwork.UserMention struct{UserID int64, Username string}, Medias []socialnetwork.Media struct{MediaID int64, MediaType string}, Urls []socialnetwork.URL struct{ShortenedUrl string, ExpandedUrl string}, Timestamp int64, PostType int64}
[_1] (FieldObject FieldType) Creator socialnetwork.Creator struct{UserID int64, Username string}
[__2] (StructObject UserType) creator socialnetwork.Creator struct{UserID int64, Username string}
[_1] (FieldObject FieldType) Medias []socialnetwork.Media struct{MediaID int64, MediaType string}
[__2] (ArrayObject SliceType) medias []socialnetwork.Media struct{MediaID int64, MediaType string}
[_1] (FieldObject FieldType) PostID int64
[__2] (BasicObject BasicType) uniqueID int64
[_1] (FieldObject FieldType) PostType int64
[__2] (BasicObject BasicType) postType int64
[_1] (FieldObject FieldType) ReqID int64
[__2] (BasicObject BasicType) reqID int64
[_1] (FieldObject FieldType) Text string
[__2] (BasicObject BasicType) up_text string
[_1] (FieldObject FieldType) Urls []socialnetwork.URL struct{ShortenedUrl string, ExpandedUrl string}
[__2] (ArrayObject SliceType) urls []socialnetwork.URL struct{ShortenedUrl string, ExpandedUrl string}
[_1] (FieldObject FieldType) UserMentions []socialnetwork.UserMention struct{UserID int64, Username string}
[__2] (ArrayObject SliceType) usermentions []socialnetwork.UserMention struct{UserID int64, Username string}

[0] (ArrayObject ArrayType) usermentionIds []int64
[_1] (BasicObject BasicType) UserID int64

[0] (StructObject UserType) um socialnetwork.UserMention struct{UserID int64, Username string}
[_1] (ArrayObject SliceType) usermentions []socialnetwork.UserMention struct{UserID int64, Username string}
[_1] (FieldObject FieldType) UserID int64
[__2] (BasicObject BasicType) UserID int64

[0] (ArrayObject ArrayType) usermentionIds []int64
[_1] (BasicObject BasicType) UserID int64

[0] (PointerObject PointerType) c (*socialnetwork.ComposePostServiceImpl struct{postStorageService socialnetwork.PostStorageService, userTimelineService socialnetwork.UserTimelineService, userService socialnetwork.UserService, uniqueIDService socialnetwork.UniqueIdService, mediaService socialnetwork.MediaService, textService socialnetwork.TextService, homeTimelineService socialnetwork.HomeTimelineService})
[_1] (StructObject UserType) socialnetwork.ComposePostServiceImpl struct{postStorageService socialnetwork.PostStorageService, userTimelineService socialnetwork.UserTimelineService, userService socialnetwork.UserService, uniqueIDService socialnetwork.UniqueIdService, mediaService socialnetwork.MediaService, textService socialnetwork.TextService, homeTimelineService socialnetwork.HomeTimelineService}
[__2] (FieldObject FieldType) homeTimelineService socialnetwork.HomeTimelineService
[___3] (ServiceObject ServiceType) homeTimelineService socialnetwork.HomeTimelineService
[__2] (FieldObject FieldType) mediaService socialnetwork.MediaService
[___3] (ServiceObject ServiceType) mediaService socialnetwork.MediaService
[__2] (FieldObject FieldType) postStorageService socialnetwork.PostStorageService
[___3] (ServiceObject ServiceType) postStorageService socialnetwork.PostStorageService
[__2] (FieldObject FieldType) textService socialnetwork.TextService
[___3] (ServiceObject ServiceType) textService socialnetwork.TextService
[__2] (FieldObject FieldType) uniqueIDService socialnetwork.UniqueIdService
[___3] (ServiceObject ServiceType) uniqueIDService socialnetwork.UniqueIdService
[__2] (FieldObject FieldType) userService socialnetwork.UserService
[___3] (ServiceObject ServiceType) userService socialnetwork.UserService
[__2] (FieldObject FieldType) userTimelineService socialnetwork.UserTimelineService
[___3] (ServiceObject ServiceType) userTimelineService socialnetwork.UserTimelineService

[0] (InterfaceObject UserType) ctx context.Context

[0] (BasicObject BasicType) reqID int64

[0] (BasicObject BasicType) username string

[0] (BasicObject BasicType) userID int64

[0] (BasicObject BasicType) text string

[0] (ArrayObject ArrayType) mediaIDs []int64

[0] (ArrayObject ArrayType) mediaTypes []string

[0] (BasicObject BasicType) postType int64

[0] (BasicObject BasicType) timestamp int64

[0] (BasicObject BasicType) err1 error

[0] (BasicObject BasicType) err2 error

[0] (BasicObject BasicType) err3 error

[0] (BasicObject BasicType) err4 error

[0] (BasicObject BasicType) uniqueID int64

[0] (StructObject UserType) creator socialnetwork.Creator struct{UserID int64, Username string}

[0] (BasicObject BasicType) up_text string

[0] (ArrayObject SliceType) medias []socialnetwork.Media struct{MediaID int64, MediaType string}

[0] (ArrayObject SliceType) urls []socialnetwork.URL struct{ShortenedUrl string, ExpandedUrl string}

[0] (ArrayObject SliceType) usermentions []socialnetwork.UserMention struct{UserID int64, Username string}

[0] (StructObject UserType) wg sync.WaitGroup

[0] (StructObject UserType) post socialnetwork.Post struct{PostID int64, Creator socialnetwork.Creator struct{UserID int64, Username string}, ReqID int64, Text string, UserMentions []socialnetwork.UserMention struct{UserID int64, Username string}, Medias []socialnetwork.Media struct{MediaID int64, MediaType string}, Urls []socialnetwork.URL struct{ShortenedUrl string, ExpandedUrl string}, Timestamp int64, PostType int64}
[_1] (FieldObject FieldType) Creator socialnetwork.Creator struct{UserID int64, Username string}
[__2] (StructObject UserType) creator socialnetwork.Creator struct{UserID int64, Username string}
[_1] (FieldObject FieldType) Medias []socialnetwork.Media struct{MediaID int64, MediaType string}
[__2] (ArrayObject SliceType) medias []socialnetwork.Media struct{MediaID int64, MediaType string}
[_1] (FieldObject FieldType) PostID int64
[__2] (BasicObject BasicType) uniqueID int64
[_1] (FieldObject FieldType) PostType int64
[__2] (BasicObject BasicType) postType int64
[_1] (FieldObject FieldType) ReqID int64
[__2] (BasicObject BasicType) reqID int64
[_1] (FieldObject FieldType) Text string
[__2] (BasicObject BasicType) up_text string
[_1] (FieldObject FieldType) Urls []socialnetwork.URL struct{ShortenedUrl string, ExpandedUrl string}
[__2] (ArrayObject SliceType) urls []socialnetwork.URL struct{ShortenedUrl string, ExpandedUrl string}
[_1] (FieldObject FieldType) UserMentions []socialnetwork.UserMention struct{UserID int64, Username string}
[__2] (ArrayObject SliceType) usermentions []socialnetwork.UserMention struct{UserID int64, Username string}

[0] (ArrayObject ArrayType) usermentionIds []int64
[_1] (BasicObject BasicType) UserID int64

[0] (StructObject UserType) um socialnetwork.UserMention struct{UserID int64, Username string}
[_1] (ArrayObject SliceType) usermentions []socialnetwork.UserMention struct{UserID int64, Username string}
[_1] (FieldObject FieldType) UserID int64
[__2] (BasicObject BasicType) UserID int64

[0] (StructObject UserType) wg2 sync.WaitGroup

