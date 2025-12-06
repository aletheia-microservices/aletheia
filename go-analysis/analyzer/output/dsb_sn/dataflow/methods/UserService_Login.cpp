[0] (PointerObject PointerType) u (*socialnetwork.UserServiceImpl struct{machine_id string, counter int64, current_timestamp int64, secret string, userCache Cache, userDB NoSQLDatabase, socialGraphService socialnetwork.SocialGraphService})
[_1] (StructObject UserType) socialnetwork.UserServiceImpl struct{machine_id string, counter 0 int, current_timestamp int64, secret string, userCache Cache, userDB NoSQLDatabase, socialGraphService socialnetwork.SocialGraphService}
[__2] (FieldObject FieldType) counter 0 int
[___3] (BasicObject BasicType) 1 int
[__2] (FieldObject FieldType) current_timestamp int64
[___3] (BasicObject BasicType) timestamp int64
[__2] (FieldObject FieldType) machine_id string
[___3] (BasicObject BasicType) string
[__2] (FieldObject FieldType) secret string
[___3] (BasicObject BasicType) secret string
[__2] (FieldObject FieldType) socialGraphService socialnetwork.SocialGraphService
[___3] (ServiceObject ServiceType) socialGraphService socialnetwork.SocialGraphService
[__2] (FieldObject FieldType) userCache Cache
[___3] (BlueprintBackendObject BlueprintBackendType) userCache Cache
[__2] (FieldObject FieldType) userDB NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) userDB NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context

[0] (BasicObject BasicType) reqID int64

[0] (BasicObject BasicType) username string

[0] (BasicObject BasicType) password string

[0] (BasicObject BasicType) timestamp int64

[0] (StructObject UserType) login socialnetwork.LoginObj struct{UserID int64, Password string, Salt string}
[_1] (FieldObject FieldType) Password string
[__2] (BasicObject BasicType) PwdHashed string
[_1] (FieldObject FieldType) Salt string
[__2] (BasicObject BasicType) Salt string
[_1] (FieldObject FieldType) UserID int64
[__2] (BasicObject BasicType) UserID int64

[0] (StructObject UserType) user socialnetwork.User struct{UserID int64, FirstName string, LastName string, Username string, PwdHashed string, Salt string}
[_1] (FieldObject FieldType) PwdHashed string
[__2] (BasicObject BasicType) PwdHashed string
[_1] (FieldObject FieldType) Salt string
[__2] (BasicObject BasicType) Salt string
[_1] (FieldObject FieldType) UserID int64
[__2] (BasicObject BasicType) UserID int64

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = user, collection = user}

[0] (InterfaceObject UserType) err .error

[0] (SliceObject UserType) query primitive.D
[_1] (StructObject StructType) struct{ string,  "username" string,  string}

[0] (BlueprintBackendObject BlueprintBackendType) res NoSQLCursor {database = user, collection = user}
[_1] (StructObject UserType) user socialnetwork.User struct{UserID int64, FirstName string, LastName string, Username string, PwdHashed string, Salt string}
[__2] (FieldObject FieldType) PwdHashed string
[___3] (BasicObject BasicType) PwdHashed string
[__2] (FieldObject FieldType) Salt string
[___3] (BasicObject BasicType) Salt string
[__2] (FieldObject FieldType) UserID int64
[___3] (BasicObject BasicType) UserID int64

[0] (InterfaceObject UserType) err .error

[0] (BasicObject BasicType) tokenStr string

[0] (BasicObject BasicType) hashed_pwd string

[0] (StructObject UserType) expiration_time time.Time struct{wall uint64, ext int64, loc (*time.Location struct{name string, zone []time.zone struct{name string, offset int, isDST bool}, tx []time.zoneTrans struct{when int64, index uint8, isstd bool, isutc bool}, extend string, cacheStart int64, cacheEnd int64, cacheZone (*time.zone struct{name string, offset int, isDST bool})})}
[_1] (BasicObject UserType) time.Duration int64
[_1] (BasicObject UserType) time.Duration int64

[0] (PointerObject PointerType) claims (*socialnetwork.claimsT struct{Username string, UserID string, Timestamp int64, StandardClaims jwt-go.StandardClaims struct{Audience string, ExpiresAt int64, Id string, IssuedAt int64, Issuer string, NotBefore int64, Subject string}})
[_1] (StructObject UserType) socialnetwork.claimsT struct{Username string, UserID string, Timestamp int64, StandardClaims jwt-go.StandardClaims struct{Audience string, ExpiresAt int64, Id string, IssuedAt int64, Issuer string, NotBefore int64, Subject string}}
[__2] (FieldObject FieldType) StandardClaims jwt-go.StandardClaims struct{Audience string, ExpiresAt int64, Id string, IssuedAt int64, Issuer string, NotBefore int64, Subject string}
[___3] (StructObject UserType) jwt-go.StandardClaims struct{Audience string, ExpiresAt int64, Id string, IssuedAt int64, Issuer string, NotBefore int64, Subject string}
[____4] (FieldObject FieldType) ExpiresAt int64
[_____5] (BasicObject BasicType) int64
[__2] (FieldObject FieldType) Timestamp int64
[___3] (BasicObject BasicType) timestamp int64
[__2] (FieldObject FieldType) UserID string
[___3] (BasicObject BasicType) string
[____4] (BasicObject BasicType) UserID int64
[____4] (BasicObject BasicType) 10 int
[____4] (BasicObject BasicType) UserID int64
[____4] (BasicObject BasicType) 10 int
[__2] (FieldObject FieldType) Username string
[___3] (BasicObject BasicType) username string

[0] (BasicObject BasicType) err error

[0] (PointerObject PointerType) token (*jwt-go.Token struct{Raw string, Method jwt-go.SigningMethod interface{ interface{Alg() string; Sign(signingString string, key interface{}) (string, error); Verify(signingString string, signature string, key interface{}) error} }, Header map[string]interface{ interface{} }, Claims jwt-go.Claims interface{ interface{Valid() error} }, Signature string, Valid bool})
[_1] (PointerObject PointerType) (*jwt-go.SigningMethodHMAC struct{Name string, Hash crypto.Hash uint})
[__2] (StructObject UserType) jwt-go.SigningMethodHMAC struct{Name string, Hash crypto.Hash uint}
[_1] (PointerObject PointerType) claims (*socialnetwork.claimsT struct{Username string, UserID string, Timestamp int64, StandardClaims jwt-go.StandardClaims struct{Audience string, ExpiresAt int64, Id string, IssuedAt int64, Issuer string, NotBefore int64, Subject string}})
[__2] (StructObject UserType) socialnetwork.claimsT struct{Username string, UserID string, Timestamp int64, StandardClaims jwt-go.StandardClaims struct{Audience string, ExpiresAt int64, Id string, IssuedAt int64, Issuer string, NotBefore int64, Subject string}}
[___3] (FieldObject FieldType) StandardClaims jwt-go.StandardClaims struct{Audience string, ExpiresAt int64, Id string, IssuedAt int64, Issuer string, NotBefore int64, Subject string}
[____4] (StructObject UserType) jwt-go.StandardClaims struct{Audience string, ExpiresAt int64, Id string, IssuedAt int64, Issuer string, NotBefore int64, Subject string}
[_____5] (FieldObject FieldType) ExpiresAt int64
[______6] (BasicObject BasicType) int64
[___3] (FieldObject FieldType) Timestamp int64
[____4] (BasicObject BasicType) timestamp int64
[___3] (FieldObject FieldType) UserID string
[____4] (BasicObject BasicType) string
[_____5] (BasicObject BasicType) UserID int64
[_____5] (BasicObject BasicType) 10 int
[_____5] (BasicObject BasicType) UserID int64
[_____5] (BasicObject BasicType) 10 int
[___3] (FieldObject FieldType) Username string
[____4] (BasicObject BasicType) username string
[_1] (PointerObject PointerType) (*jwt-go.SigningMethodHMAC struct{Name string, Hash crypto.Hash uint})
[__2] (StructObject UserType) jwt-go.SigningMethodHMAC struct{Name string, Hash crypto.Hash uint}
[_1] (PointerObject PointerType) claims (*socialnetwork.claimsT struct{Username string, UserID string, Timestamp int64, StandardClaims jwt-go.StandardClaims struct{Audience string, ExpiresAt int64, Id string, IssuedAt int64, Issuer string, NotBefore int64, Subject string}})
[__2] (StructObject UserType) socialnetwork.claimsT struct{Username string, UserID string, Timestamp int64, StandardClaims jwt-go.StandardClaims struct{Audience string, ExpiresAt int64, Id string, IssuedAt int64, Issuer string, NotBefore int64, Subject string}}
[___3] (FieldObject FieldType) StandardClaims jwt-go.StandardClaims struct{Audience string, ExpiresAt int64, Id string, IssuedAt int64, Issuer string, NotBefore int64, Subject string}
[____4] (StructObject UserType) jwt-go.StandardClaims struct{Audience string, ExpiresAt int64, Id string, IssuedAt int64, Issuer string, NotBefore int64, Subject string}
[_____5] (FieldObject FieldType) ExpiresAt int64
[______6] (BasicObject BasicType) int64
[___3] (FieldObject FieldType) Timestamp int64
[____4] (BasicObject BasicType) timestamp int64
[___3] (FieldObject FieldType) UserID string
[____4] (BasicObject BasicType) string
[_____5] (BasicObject BasicType) UserID int64
[_____5] (BasicObject BasicType) 10 int
[_____5] (BasicObject BasicType) UserID int64
[_____5] (BasicObject BasicType) 10 int
[___3] (FieldObject FieldType) Username string
[____4] (BasicObject BasicType) username string
[_1] (StructObject UserType) jwt-go.Token struct{Raw string, Method jwt-go.SigningMethod interface{ interface{Alg() string; Sign(signingString string, key interface{}) (string, error); Verify(signingString string, signature string, key interface{}) error} }, Header map[string]interface{ interface{} }, Claims jwt-go.Claims interface{ interface{Valid() error} }, Signature string, Valid bool}

[0] (BasicObject BasicType) tokenStr string
[_1] (BasicObject BasicType) byte
[__2] (BasicObject BasicType) secret string
[_1] (BasicObject BasicType) byte
[__2] (BasicObject BasicType) secret string

[0] (InterfaceObject UserType) err .error
[_1] (BasicObject BasicType) byte
[__2] (BasicObject BasicType) secret string

[0] (InterfaceObject UserType) err .error

