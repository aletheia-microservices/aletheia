[0] (PointerObject PointerType) r (*hotelreservation.ReservationServiceImpl struct{reserveCache Cache, reserveDB NoSQLDatabase, CacheHits int64, NumRequests int64})
[_1] (StructObject UserType) hotelreservation.ReservationServiceImpl struct{reserveCache Cache, reserveDB NoSQLDatabase, CacheHits int64, NumRequests int64}
[__2] (FieldObject FieldType) CacheHits int64
[___3] (BasicObject BasicType) 1 int
[__2] (FieldObject FieldType) NumRequests int64
[___3] (BasicObject BasicType) 1 int
[__2] (FieldObject FieldType) reserveCache Cache
[___3] (BlueprintBackendObject BlueprintBackendType) reserveCache Cache
[__2] (FieldObject FieldType) reserveDB NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) reserveDB NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context

[0] (BasicObject BasicType) customerName string

[0] (ArrayObject ArrayType) hotelIds []string

[0] (BasicObject BasicType) inDate string

[0] (BasicObject BasicType) outDate string

[0] (BasicObject BasicType) roomNumber int64

[0] (BlueprintBackendObject BlueprintBackendType) reservation_collection NoSQLCollection {database = reservation-db, collection = reservation}

[0] (InterfaceObject UserType) err .error

[0] (BlueprintBackendObject BlueprintBackendType) hnumber_collection NoSQLCollection {database = reservation-db, collection = number}

[0] (InterfaceObject UserType) err .error

[0] (StructObject UserType) newInDate time.Time struct{wall uint64, ext int64, loc (*time.Location struct{name string, zone []time.zone struct{name string, offset int, isDST bool}, tx []time.zoneTrans struct{when int64, index uint8, isstd bool, isutc bool}, extend string, cacheStart int64, cacheEnd int64, cacheZone (*time.zone struct{name string, offset int, isDST bool})})}
[_1] (BasicObject BasicType) untyped string
[_1] (BasicObject BasicType) T12:00:00+00:00 
[__2] (BasicObject BasicType) inDate string
[__2] (BasicObject BasicType) "T12:00:00+00:00" string
[_1] (BasicObject BasicType) untyped string
[_1] (BasicObject BasicType) T12:00:00+00:00 
[__2] (BasicObject BasicType) inDate string
[__2] (BasicObject BasicType) "T12:00:00+00:00" string

[0] (InterfaceObject UserType) _ .error
[_1] (BasicObject BasicType) untyped string
[_1] (BasicObject BasicType) T12:00:00+00:00 
[__2] (BasicObject BasicType) inDate string
[__2] (BasicObject BasicType) "T12:00:00+00:00" string

[0] (StructObject UserType) newOutDate time.Time struct{wall uint64, ext int64, loc (*time.Location struct{name string, zone []time.zone struct{name string, offset int, isDST bool}, tx []time.zoneTrans struct{when int64, index uint8, isstd bool, isutc bool}, extend string, cacheStart int64, cacheEnd int64, cacheZone (*time.zone struct{name string, offset int, isDST bool})})}
[_1] (BasicObject BasicType) untyped string
[_1] (BasicObject BasicType) T12:00:00+00:00 
[__2] (BasicObject BasicType) outDate string
[__2] (BasicObject BasicType) "T12:00:00+00:00" string
[_1] (BasicObject BasicType) untyped string
[_1] (BasicObject BasicType) T12:00:00+00:00 
[__2] (BasicObject BasicType) outDate string
[__2] (BasicObject BasicType) "T12:00:00+00:00" string

[0] (InterfaceObject UserType) _ .error
[_1] (BasicObject BasicType) untyped string
[_1] (BasicObject BasicType) T12:00:00+00:00 
[__2] (BasicObject BasicType) outDate string
[__2] (BasicObject BasicType) "T12:00:00+00:00" string

[0] (BasicObject BasicType) hotelId string
[_1] (BasicObject BasicType) * string

[0] (BasicObject BasicType) indate string

[0] (MapObject MapType) reservation_update_map map[string]int64

[0] (StructObject UserType) newInDate time.Time struct{wall uint64, ext int64, loc (*time.Location struct{name string, zone []time.zone struct{name string, offset int, isDST bool}, tx []time.zoneTrans struct{when int64, index uint8, isstd bool, isutc bool}, extend string, cacheStart int64, cacheEnd int64, cacheZone (*time.zone struct{name string, offset int, isDST bool})})}
[_1] (BasicObject BasicType) 0 int
[_1] (BasicObject BasicType) 0 int
[_1] (BasicObject BasicType) 1 int
[_1] (BasicObject BasicType) 0 int
[_1] (BasicObject BasicType) 0 int
[_1] (BasicObject BasicType) 1 int

[0] (BasicObject BasicType) outdate string

[0] (BasicObject BasicType) key string

[0] (BasicObject BasicType) room_number int64

[0] (BasicObject BasicType) exists bool

[0] (InterfaceObject UserType) err .error

[0] (ArrayObject ArrayType) reservations []hotelreservation.Reservation struct{HotelId string, CustomerName string, InDate string, OutDate string, Number int64}

[0] (SliceObject UserType) query primitive.D
[_1] (StructObject StructType) struct{ string,  "hotelid" string,  string}
[_1] (StructObject StructType) struct{ string,  "indate" string,  string}
[_1] (StructObject StructType) struct{ string,  "outdate" string,  string}

[0] (BlueprintBackendObject BlueprintBackendType) res NoSQLCursor {database = reservation-db, collection = reservation}
[_1] (ArrayObject ArrayType) reservations []hotelreservation.Reservation struct{HotelId string, CustomerName string, InDate string, OutDate string, Number int64}

[0] (InterfaceObject UserType) err .error

[0] (StructObject UserType) reservation hotelreservation.Reservation struct{HotelId string, CustomerName string, InDate string, OutDate string, Number int64}
[_1] (ArrayObject ArrayType) reservations []hotelreservation.Reservation struct{HotelId string, CustomerName string, InDate string, OutDate string, Number int64}
[_1] (FieldObject FieldType) Number int64
[__2] (BasicObject BasicType) Number int64

[0] (PointerObject PointerType) r (*hotelreservation.ReservationServiceImpl struct{reserveCache Cache, reserveDB NoSQLDatabase, CacheHits int64, NumRequests int64})
[_1] (StructObject UserType) hotelreservation.ReservationServiceImpl struct{reserveCache Cache, reserveDB NoSQLDatabase, CacheHits int64, NumRequests int64}
[__2] (FieldObject FieldType) CacheHits int64
[___3] (BasicObject BasicType) 1 int
[__2] (FieldObject FieldType) NumRequests int64
[___3] (BasicObject BasicType) 1 int
[__2] (FieldObject FieldType) reserveCache Cache
[___3] (BlueprintBackendObject BlueprintBackendType) reserveCache Cache
[__2] (FieldObject FieldType) reserveDB NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) reserveDB NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context

[0] (BasicObject BasicType) customerName string

[0] (ArrayObject ArrayType) hotelIds []string

[0] (BasicObject BasicType) inDate string

[0] (BasicObject BasicType) outDate string

[0] (BasicObject BasicType) roomNumber int64

[0] (BlueprintBackendObject BlueprintBackendType) reservation_collection NoSQLCollection {database = reservation-db, collection = reservation}

[0] (InterfaceObject UserType) err .error

[0] (BlueprintBackendObject BlueprintBackendType) hnumber_collection NoSQLCollection {database = reservation-db, collection = number}

[0] (InterfaceObject UserType) err .error

[0] (StructObject UserType) newInDate time.Time struct{wall uint64, ext int64, loc (*time.Location struct{name string, zone []time.zone struct{name string, offset int, isDST bool}, tx []time.zoneTrans struct{when int64, index uint8, isstd bool, isutc bool}, extend string, cacheStart int64, cacheEnd int64, cacheZone (*time.zone struct{name string, offset int, isDST bool})})}
[_1] (BasicObject BasicType) untyped string
[_1] (BasicObject BasicType) T12:00:00+00:00 
[__2] (BasicObject BasicType) inDate string
[__2] (BasicObject BasicType) "T12:00:00+00:00" string
[_1] (BasicObject BasicType) untyped string
[_1] (BasicObject BasicType) T12:00:00+00:00 
[__2] (BasicObject BasicType) inDate string
[__2] (BasicObject BasicType) "T12:00:00+00:00" string

[0] (InterfaceObject UserType) _ .error
[_1] (BasicObject BasicType) untyped string
[_1] (BasicObject BasicType) T12:00:00+00:00 
[__2] (BasicObject BasicType) inDate string
[__2] (BasicObject BasicType) "T12:00:00+00:00" string

[0] (StructObject UserType) newOutDate time.Time struct{wall uint64, ext int64, loc (*time.Location struct{name string, zone []time.zone struct{name string, offset int, isDST bool}, tx []time.zoneTrans struct{when int64, index uint8, isstd bool, isutc bool}, extend string, cacheStart int64, cacheEnd int64, cacheZone (*time.zone struct{name string, offset int, isDST bool})})}
[_1] (BasicObject BasicType) untyped string
[_1] (BasicObject BasicType) T12:00:00+00:00 
[__2] (BasicObject BasicType) outDate string
[__2] (BasicObject BasicType) "T12:00:00+00:00" string
[_1] (BasicObject BasicType) untyped string
[_1] (BasicObject BasicType) T12:00:00+00:00 
[__2] (BasicObject BasicType) outDate string
[__2] (BasicObject BasicType) "T12:00:00+00:00" string

[0] (InterfaceObject UserType) _ .error
[_1] (BasicObject BasicType) untyped string
[_1] (BasicObject BasicType) T12:00:00+00:00 
[__2] (BasicObject BasicType) outDate string
[__2] (BasicObject BasicType) "T12:00:00+00:00" string

[0] (BasicObject BasicType) hotelId string
[_1] (BasicObject BasicType) * string

[0] (BasicObject BasicType) indate string

[0] (MapObject MapType) reservation_update_map map[string]int64

[0] (StructObject UserType) newInDate time.Time struct{wall uint64, ext int64, loc (*time.Location struct{name string, zone []time.zone struct{name string, offset int, isDST bool}, tx []time.zoneTrans struct{when int64, index uint8, isstd bool, isutc bool}, extend string, cacheStart int64, cacheEnd int64, cacheZone (*time.zone struct{name string, offset int, isDST bool})})}
[_1] (BasicObject BasicType) 0 int
[_1] (BasicObject BasicType) 0 int
[_1] (BasicObject BasicType) 1 int
[_1] (BasicObject BasicType) 0 int
[_1] (BasicObject BasicType) 0 int
[_1] (BasicObject BasicType) 1 int

[0] (BasicObject BasicType) outdate string

[0] (BasicObject BasicType) key string

[0] (BasicObject BasicType) room_number int64

[0] (BasicObject BasicType) exists bool

[0] (InterfaceObject UserType) err .error

[0] (ArrayObject ArrayType) reservations []hotelreservation.Reservation struct{HotelId string, CustomerName string, InDate string, OutDate string, Number int64}

[0] (SliceObject UserType) query primitive.D
[_1] (StructObject StructType) struct{ string,  "hotelid" string,  string}
[_1] (StructObject StructType) struct{ string,  "indate" string,  string}
[_1] (StructObject StructType) struct{ string,  "outdate" string,  string}

[0] (BlueprintBackendObject BlueprintBackendType) res NoSQLCursor {database = reservation-db, collection = reservation}
[_1] (ArrayObject ArrayType) reservations []hotelreservation.Reservation struct{HotelId string, CustomerName string, InDate string, OutDate string, Number int64}

[0] (InterfaceObject UserType) err .error

[0] (StructObject UserType) reservation hotelreservation.Reservation struct{HotelId string, CustomerName string, InDate string, OutDate string, Number int64}
[_1] (ArrayObject ArrayType) reservations []hotelreservation.Reservation struct{HotelId string, CustomerName string, InDate string, OutDate string, Number int64}
[_1] (FieldObject FieldType) Number int64
[__2] (BasicObject BasicType) Number int64

[0] (BasicObject BasicType) cap_key _cap 
[_1] (BasicObject BasicType) hotelId string
[__2] (BasicObject BasicType) * string
[_1] (BasicObject BasicType) "_cap" string

[0] (StructObject UserType) hotelNumber hotelreservation.HotelNumber struct{HotelId string, Number int64}
[_1] (FieldObject FieldType) Number int64
[__2] (BasicObject BasicType) Number int64

[0] (BasicObject BasicType) capacity int64

[0] (BasicObject BasicType) exists bool

[0] (InterfaceObject UserType) err .error

[0] (SliceObject UserType) query primitive.D
[_1] (StructObject StructType) struct{ string,  "hotelid" string,  string}

[0] (BlueprintBackendObject BlueprintBackendType) res NoSQLCursor {database = reservation-db, collection = number}
[_1] (StructObject UserType) hotelNumber hotelreservation.HotelNumber struct{HotelId string, Number int64}
[__2] (FieldObject FieldType) Number int64
[___3] (BasicObject BasicType) Number int64

[0] (InterfaceObject UserType) err .error

[0] (FieldObject FieldType) capacity int64
[_1] (FieldObject FieldType) Number int64
[__2] (BasicObject BasicType) Number int64
[_1] (BasicObject BasicType) int64
[__2] (BasicObject BasicType) Number int64

[0] (InterfaceObject UserType) err .error

[0] (BasicObject BasicType) indate string
[_1] (BasicObject BasicType) outdate string

[0] (PointerObject PointerType) r (*hotelreservation.ReservationServiceImpl struct{reserveCache Cache, reserveDB NoSQLDatabase, CacheHits int64, NumRequests int64})
[_1] (StructObject UserType) hotelreservation.ReservationServiceImpl struct{reserveCache Cache, reserveDB NoSQLDatabase, CacheHits int64, NumRequests int64}
[__2] (FieldObject FieldType) CacheHits int64
[___3] (BasicObject BasicType) 1 int
[__2] (FieldObject FieldType) NumRequests int64
[___3] (BasicObject BasicType) 1 int
[__2] (FieldObject FieldType) reserveCache Cache
[___3] (BlueprintBackendObject BlueprintBackendType) reserveCache Cache
[__2] (FieldObject FieldType) reserveDB NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) reserveDB NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context

[0] (BasicObject BasicType) customerName string

[0] (ArrayObject ArrayType) hotelIds []string

[0] (BasicObject BasicType) inDate string

[0] (BasicObject BasicType) outDate string

[0] (BasicObject BasicType) roomNumber int64

[0] (BlueprintBackendObject BlueprintBackendType) reservation_collection NoSQLCollection {database = reservation-db, collection = reservation}

[0] (InterfaceObject UserType) err .error

[0] (BlueprintBackendObject BlueprintBackendType) hnumber_collection NoSQLCollection {database = reservation-db, collection = number}

[0] (InterfaceObject UserType) err .error

[0] (StructObject UserType) newInDate time.Time struct{wall uint64, ext int64, loc (*time.Location struct{name string, zone []time.zone struct{name string, offset int, isDST bool}, tx []time.zoneTrans struct{when int64, index uint8, isstd bool, isutc bool}, extend string, cacheStart int64, cacheEnd int64, cacheZone (*time.zone struct{name string, offset int, isDST bool})})}
[_1] (BasicObject BasicType) untyped string
[_1] (BasicObject BasicType) T12:00:00+00:00 
[__2] (BasicObject BasicType) inDate string
[__2] (BasicObject BasicType) "T12:00:00+00:00" string
[_1] (BasicObject BasicType) untyped string
[_1] (BasicObject BasicType) T12:00:00+00:00 
[__2] (BasicObject BasicType) inDate string
[__2] (BasicObject BasicType) "T12:00:00+00:00" string

[0] (InterfaceObject UserType) _ .error
[_1] (BasicObject BasicType) untyped string
[_1] (BasicObject BasicType) T12:00:00+00:00 
[__2] (BasicObject BasicType) inDate string
[__2] (BasicObject BasicType) "T12:00:00+00:00" string

[0] (StructObject UserType) newOutDate time.Time struct{wall uint64, ext int64, loc (*time.Location struct{name string, zone []time.zone struct{name string, offset int, isDST bool}, tx []time.zoneTrans struct{when int64, index uint8, isstd bool, isutc bool}, extend string, cacheStart int64, cacheEnd int64, cacheZone (*time.zone struct{name string, offset int, isDST bool})})}
[_1] (BasicObject BasicType) untyped string
[_1] (BasicObject BasicType) T12:00:00+00:00 
[__2] (BasicObject BasicType) outDate string
[__2] (BasicObject BasicType) "T12:00:00+00:00" string
[_1] (BasicObject BasicType) untyped string
[_1] (BasicObject BasicType) T12:00:00+00:00 
[__2] (BasicObject BasicType) outDate string
[__2] (BasicObject BasicType) "T12:00:00+00:00" string

[0] (InterfaceObject UserType) _ .error
[_1] (BasicObject BasicType) untyped string
[_1] (BasicObject BasicType) T12:00:00+00:00 
[__2] (BasicObject BasicType) outDate string
[__2] (BasicObject BasicType) "T12:00:00+00:00" string

[0] (BasicObject BasicType) hotelId string
[_1] (BasicObject BasicType) * string

[0] (BasicObject BasicType) indate string

[0] (MapObject MapType) reservation_update_map map[string]int64

[0] (StructObject UserType) newInDate time.Time struct{wall uint64, ext int64, loc (*time.Location struct{name string, zone []time.zone struct{name string, offset int, isDST bool}, tx []time.zoneTrans struct{when int64, index uint8, isstd bool, isutc bool}, extend string, cacheStart int64, cacheEnd int64, cacheZone (*time.zone struct{name string, offset int, isDST bool})})}
[_1] (BasicObject BasicType) untyped string
[_1] (BasicObject BasicType) T12:00:00+00:00 
[__2] (BasicObject BasicType) inDate string
[__2] (BasicObject BasicType) "T12:00:00+00:00" string
[_1] (BasicObject BasicType) untyped string
[_1] (BasicObject BasicType) T12:00:00+00:00 
[__2] (BasicObject BasicType) inDate string
[__2] (BasicObject BasicType) "T12:00:00+00:00" string

[0] (InterfaceObject UserType) _ .error
[_1] (BasicObject BasicType) untyped string
[_1] (BasicObject BasicType) T12:00:00+00:00 
[__2] (BasicObject BasicType) inDate string
[__2] (BasicObject BasicType) "T12:00:00+00:00" string

[0] (BasicObject BasicType) indate string

[0] (StructObject UserType) newInDate time.Time struct{wall uint64, ext int64, loc (*time.Location struct{name string, zone []time.zone struct{name string, offset int, isDST bool}, tx []time.zoneTrans struct{when int64, index uint8, isstd bool, isutc bool}, extend string, cacheStart int64, cacheEnd int64, cacheZone (*time.zone struct{name string, offset int, isDST bool})})}
[_1] (BasicObject BasicType) 0 int
[_1] (BasicObject BasicType) 0 int
[_1] (BasicObject BasicType) 1 int
[_1] (BasicObject BasicType) 0 int
[_1] (BasicObject BasicType) 0 int
[_1] (BasicObject BasicType) 1 int

[0] (BasicObject BasicType) outdate string

[0] (StructObject UserType) reservation hotelreservation.Reservation struct{HotelId string, CustomerName string, InDate string, OutDate string, Number int64}
[_1] (FieldObject FieldType) CustomerName string
[__2] (BasicObject BasicType) customerName string
[_1] (FieldObject FieldType) HotelId string
[__2] (BasicObject BasicType) hotelId string
[___3] (BasicObject BasicType) * string
[_1] (FieldObject FieldType) InDate string
[__2] (BasicObject BasicType) indate string
[_1] (FieldObject FieldType) Number int64
[__2] (BasicObject BasicType) roomNumber int64
[_1] (FieldObject FieldType) OutDate string
[__2] (BasicObject BasicType) outdate string

[0] (InterfaceObject UserType) err .error

[0] (PointerObject PointerType) r (*hotelreservation.ReservationServiceImpl struct{reserveCache Cache, reserveDB NoSQLDatabase, CacheHits int64, NumRequests int64})
[_1] (StructObject UserType) hotelreservation.ReservationServiceImpl struct{reserveCache Cache, reserveDB NoSQLDatabase, CacheHits int64, NumRequests int64}
[__2] (FieldObject FieldType) CacheHits int64
[___3] (BasicObject BasicType) 1 int
[__2] (FieldObject FieldType) NumRequests int64
[___3] (BasicObject BasicType) 1 int
[__2] (FieldObject FieldType) reserveCache Cache
[___3] (BlueprintBackendObject BlueprintBackendType) reserveCache Cache
[__2] (FieldObject FieldType) reserveDB NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) reserveDB NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context

[0] (BasicObject BasicType) customerName string

[0] (ArrayObject ArrayType) hotelIds []string

[0] (BasicObject BasicType) inDate string

[0] (BasicObject BasicType) outDate string

[0] (BasicObject BasicType) roomNumber int64

[0] (BlueprintBackendObject BlueprintBackendType) reservation_collection NoSQLCollection {database = reservation-db, collection = reservation}

[0] (InterfaceObject UserType) err .error

[0] (BlueprintBackendObject BlueprintBackendType) hnumber_collection NoSQLCollection {database = reservation-db, collection = number}

[0] (InterfaceObject UserType) err .error

[0] (StructObject UserType) newInDate time.Time struct{wall uint64, ext int64, loc (*time.Location struct{name string, zone []time.zone struct{name string, offset int, isDST bool}, tx []time.zoneTrans struct{when int64, index uint8, isstd bool, isutc bool}, extend string, cacheStart int64, cacheEnd int64, cacheZone (*time.zone struct{name string, offset int, isDST bool})})}
[_1] (BasicObject BasicType) untyped string
[_1] (BasicObject BasicType) T12:00:00+00:00 
[__2] (BasicObject BasicType) inDate string
[__2] (BasicObject BasicType) "T12:00:00+00:00" string
[_1] (BasicObject BasicType) untyped string
[_1] (BasicObject BasicType) T12:00:00+00:00 
[__2] (BasicObject BasicType) inDate string
[__2] (BasicObject BasicType) "T12:00:00+00:00" string

[0] (InterfaceObject UserType) _ .error
[_1] (BasicObject BasicType) untyped string
[_1] (BasicObject BasicType) T12:00:00+00:00 
[__2] (BasicObject BasicType) inDate string
[__2] (BasicObject BasicType) "T12:00:00+00:00" string

[0] (StructObject UserType) newOutDate time.Time struct{wall uint64, ext int64, loc (*time.Location struct{name string, zone []time.zone struct{name string, offset int, isDST bool}, tx []time.zoneTrans struct{when int64, index uint8, isstd bool, isutc bool}, extend string, cacheStart int64, cacheEnd int64, cacheZone (*time.zone struct{name string, offset int, isDST bool})})}
[_1] (BasicObject BasicType) untyped string
[_1] (BasicObject BasicType) T12:00:00+00:00 
[__2] (BasicObject BasicType) outDate string
[__2] (BasicObject BasicType) "T12:00:00+00:00" string
[_1] (BasicObject BasicType) untyped string
[_1] (BasicObject BasicType) T12:00:00+00:00 
[__2] (BasicObject BasicType) outDate string
[__2] (BasicObject BasicType) "T12:00:00+00:00" string

[0] (InterfaceObject UserType) _ .error
[_1] (BasicObject BasicType) untyped string
[_1] (BasicObject BasicType) T12:00:00+00:00 
[__2] (BasicObject BasicType) outDate string
[__2] (BasicObject BasicType) "T12:00:00+00:00" string

[0] (BasicObject BasicType) hotelId string
[_1] (BasicObject BasicType) * string

[0] (BasicObject BasicType) indate string

[0] (MapObject MapType) reservation_update_map map[string]int64

[0] (StructObject UserType) newInDate time.Time struct{wall uint64, ext int64, loc (*time.Location struct{name string, zone []time.zone struct{name string, offset int, isDST bool}, tx []time.zoneTrans struct{when int64, index uint8, isstd bool, isutc bool}, extend string, cacheStart int64, cacheEnd int64, cacheZone (*time.zone struct{name string, offset int, isDST bool})})}
[_1] (BasicObject BasicType) 0 int
[_1] (BasicObject BasicType) 0 int
[_1] (BasicObject BasicType) 1 int
[_1] (BasicObject BasicType) 0 int
[_1] (BasicObject BasicType) 0 int
[_1] (BasicObject BasicType) 1 int

[0] (BasicObject BasicType) outdate string

[0] (BasicObject BasicType) key string

[0] (BasicObject BasicType) room_number int64

[0] (BasicObject BasicType) exists bool

[0] (InterfaceObject UserType) err .error

[0] (ArrayObject ArrayType) reservations []hotelreservation.Reservation struct{HotelId string, CustomerName string, InDate string, OutDate string, Number int64}

[0] (SliceObject UserType) query primitive.D
[_1] (StructObject StructType) struct{ string,  "hotelid" string,  string}
[_1] (StructObject StructType) struct{ string,  "indate" string,  string}
[_1] (StructObject StructType) struct{ string,  "outdate" string,  string}

[0] (BlueprintBackendObject BlueprintBackendType) res NoSQLCursor {database = reservation-db, collection = reservation}
[_1] (ArrayObject ArrayType) reservations []hotelreservation.Reservation struct{HotelId string, CustomerName string, InDate string, OutDate string, Number int64}

[0] (InterfaceObject UserType) err .error

[0] (StructObject UserType) reservation hotelreservation.Reservation struct{HotelId string, CustomerName string, InDate string, OutDate string, Number int64}
[_1] (ArrayObject ArrayType) reservations []hotelreservation.Reservation struct{HotelId string, CustomerName string, InDate string, OutDate string, Number int64}
[_1] (FieldObject FieldType) Number int64
[__2] (BasicObject BasicType) Number int64

[0] (PointerObject PointerType) r (*hotelreservation.ReservationServiceImpl struct{reserveCache Cache, reserveDB NoSQLDatabase, CacheHits int64, NumRequests int64})
[_1] (StructObject UserType) hotelreservation.ReservationServiceImpl struct{reserveCache Cache, reserveDB NoSQLDatabase, CacheHits int64, NumRequests int64}
[__2] (FieldObject FieldType) CacheHits int64
[___3] (BasicObject BasicType) 1 int
[__2] (FieldObject FieldType) NumRequests int64
[___3] (BasicObject BasicType) 1 int
[__2] (FieldObject FieldType) reserveCache Cache
[___3] (BlueprintBackendObject BlueprintBackendType) reserveCache Cache
[__2] (FieldObject FieldType) reserveDB NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) reserveDB NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context

[0] (BasicObject BasicType) customerName string

[0] (ArrayObject ArrayType) hotelIds []string

[0] (BasicObject BasicType) inDate string

[0] (BasicObject BasicType) outDate string

[0] (BasicObject BasicType) roomNumber int64

[0] (BlueprintBackendObject BlueprintBackendType) reservation_collection NoSQLCollection {database = reservation-db, collection = reservation}

[0] (InterfaceObject UserType) err .error

[0] (BlueprintBackendObject BlueprintBackendType) hnumber_collection NoSQLCollection {database = reservation-db, collection = number}

[0] (InterfaceObject UserType) err .error

[0] (StructObject UserType) newInDate time.Time struct{wall uint64, ext int64, loc (*time.Location struct{name string, zone []time.zone struct{name string, offset int, isDST bool}, tx []time.zoneTrans struct{when int64, index uint8, isstd bool, isutc bool}, extend string, cacheStart int64, cacheEnd int64, cacheZone (*time.zone struct{name string, offset int, isDST bool})})}
[_1] (BasicObject BasicType) untyped string
[_1] (BasicObject BasicType) T12:00:00+00:00 
[__2] (BasicObject BasicType) inDate string
[__2] (BasicObject BasicType) "T12:00:00+00:00" string
[_1] (BasicObject BasicType) untyped string
[_1] (BasicObject BasicType) T12:00:00+00:00 
[__2] (BasicObject BasicType) inDate string
[__2] (BasicObject BasicType) "T12:00:00+00:00" string

[0] (InterfaceObject UserType) _ .error
[_1] (BasicObject BasicType) untyped string
[_1] (BasicObject BasicType) T12:00:00+00:00 
[__2] (BasicObject BasicType) inDate string
[__2] (BasicObject BasicType) "T12:00:00+00:00" string

[0] (StructObject UserType) newOutDate time.Time struct{wall uint64, ext int64, loc (*time.Location struct{name string, zone []time.zone struct{name string, offset int, isDST bool}, tx []time.zoneTrans struct{when int64, index uint8, isstd bool, isutc bool}, extend string, cacheStart int64, cacheEnd int64, cacheZone (*time.zone struct{name string, offset int, isDST bool})})}
[_1] (BasicObject BasicType) untyped string
[_1] (BasicObject BasicType) T12:00:00+00:00 
[__2] (BasicObject BasicType) outDate string
[__2] (BasicObject BasicType) "T12:00:00+00:00" string
[_1] (BasicObject BasicType) untyped string
[_1] (BasicObject BasicType) T12:00:00+00:00 
[__2] (BasicObject BasicType) outDate string
[__2] (BasicObject BasicType) "T12:00:00+00:00" string

[0] (InterfaceObject UserType) _ .error
[_1] (BasicObject BasicType) untyped string
[_1] (BasicObject BasicType) T12:00:00+00:00 
[__2] (BasicObject BasicType) outDate string
[__2] (BasicObject BasicType) "T12:00:00+00:00" string

[0] (BasicObject BasicType) hotelId string
[_1] (BasicObject BasicType) * string

[0] (BasicObject BasicType) indate string

[0] (MapObject MapType) reservation_update_map map[string]int64

[0] (StructObject UserType) newInDate time.Time struct{wall uint64, ext int64, loc (*time.Location struct{name string, zone []time.zone struct{name string, offset int, isDST bool}, tx []time.zoneTrans struct{when int64, index uint8, isstd bool, isutc bool}, extend string, cacheStart int64, cacheEnd int64, cacheZone (*time.zone struct{name string, offset int, isDST bool})})}
[_1] (BasicObject BasicType) 0 int
[_1] (BasicObject BasicType) 0 int
[_1] (BasicObject BasicType) 1 int
[_1] (BasicObject BasicType) 0 int
[_1] (BasicObject BasicType) 0 int
[_1] (BasicObject BasicType) 1 int

[0] (BasicObject BasicType) outdate string

[0] (BasicObject BasicType) key string

[0] (BasicObject BasicType) room_number int64

[0] (BasicObject BasicType) exists bool

[0] (InterfaceObject UserType) err .error

[0] (ArrayObject ArrayType) reservations []hotelreservation.Reservation struct{HotelId string, CustomerName string, InDate string, OutDate string, Number int64}

[0] (SliceObject UserType) query primitive.D
[_1] (StructObject StructType) struct{ string,  "hotelid" string,  string}
[_1] (StructObject StructType) struct{ string,  "indate" string,  string}
[_1] (StructObject StructType) struct{ string,  "outdate" string,  string}

[0] (BlueprintBackendObject BlueprintBackendType) res NoSQLCursor {database = reservation-db, collection = reservation}
[_1] (ArrayObject ArrayType) reservations []hotelreservation.Reservation struct{HotelId string, CustomerName string, InDate string, OutDate string, Number int64}

[0] (InterfaceObject UserType) err .error

[0] (StructObject UserType) reservation hotelreservation.Reservation struct{HotelId string, CustomerName string, InDate string, OutDate string, Number int64}
[_1] (ArrayObject ArrayType) reservations []hotelreservation.Reservation struct{HotelId string, CustomerName string, InDate string, OutDate string, Number int64}
[_1] (FieldObject FieldType) Number int64
[__2] (BasicObject BasicType) Number int64

[0] (BasicObject BasicType) cap_key _cap 
[_1] (BasicObject BasicType) hotelId string
[__2] (BasicObject BasicType) * string
[_1] (BasicObject BasicType) "_cap" string

[0] (StructObject UserType) hotelNumber hotelreservation.HotelNumber struct{HotelId string, Number int64}
[_1] (FieldObject FieldType) Number int64
[__2] (BasicObject BasicType) Number int64

[0] (BasicObject BasicType) capacity int64

[0] (BasicObject BasicType) exists bool

[0] (InterfaceObject UserType) err .error

[0] (SliceObject UserType) query primitive.D
[_1] (StructObject StructType) struct{ string,  "hotelid" string,  string}

[0] (BlueprintBackendObject BlueprintBackendType) res NoSQLCursor {database = reservation-db, collection = number}
[_1] (StructObject UserType) hotelNumber hotelreservation.HotelNumber struct{HotelId string, Number int64}
[__2] (FieldObject FieldType) Number int64
[___3] (BasicObject BasicType) Number int64

[0] (InterfaceObject UserType) err .error

[0] (FieldObject FieldType) capacity int64
[_1] (FieldObject FieldType) Number int64
[__2] (BasicObject BasicType) Number int64
[_1] (BasicObject BasicType) int64
[__2] (BasicObject BasicType) Number int64

[0] (InterfaceObject UserType) err .error

[0] (BasicObject BasicType) indate string
[_1] (BasicObject BasicType) outdate string

[0] (PointerObject PointerType) r (*hotelreservation.ReservationServiceImpl struct{reserveCache Cache, reserveDB NoSQLDatabase, CacheHits int64, NumRequests int64})
[_1] (StructObject UserType) hotelreservation.ReservationServiceImpl struct{reserveCache Cache, reserveDB NoSQLDatabase, CacheHits int64, NumRequests int64}
[__2] (FieldObject FieldType) CacheHits int64
[___3] (BasicObject BasicType) 1 int
[__2] (FieldObject FieldType) NumRequests int64
[___3] (BasicObject BasicType) 1 int
[__2] (FieldObject FieldType) reserveCache Cache
[___3] (BlueprintBackendObject BlueprintBackendType) reserveCache Cache
[__2] (FieldObject FieldType) reserveDB NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) reserveDB NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context

[0] (BasicObject BasicType) customerName string

[0] (ArrayObject ArrayType) hotelIds []string

[0] (BasicObject BasicType) inDate string

[0] (BasicObject BasicType) outDate string

[0] (BasicObject BasicType) roomNumber int64

[0] (BlueprintBackendObject BlueprintBackendType) reservation_collection NoSQLCollection {database = reservation-db, collection = reservation}

[0] (InterfaceObject UserType) err .error

[0] (BlueprintBackendObject BlueprintBackendType) hnumber_collection NoSQLCollection {database = reservation-db, collection = number}

[0] (InterfaceObject UserType) err .error

[0] (StructObject UserType) newInDate time.Time struct{wall uint64, ext int64, loc (*time.Location struct{name string, zone []time.zone struct{name string, offset int, isDST bool}, tx []time.zoneTrans struct{when int64, index uint8, isstd bool, isutc bool}, extend string, cacheStart int64, cacheEnd int64, cacheZone (*time.zone struct{name string, offset int, isDST bool})})}
[_1] (BasicObject BasicType) untyped string
[_1] (BasicObject BasicType) T12:00:00+00:00 
[__2] (BasicObject BasicType) inDate string
[__2] (BasicObject BasicType) "T12:00:00+00:00" string
[_1] (BasicObject BasicType) untyped string
[_1] (BasicObject BasicType) T12:00:00+00:00 
[__2] (BasicObject BasicType) inDate string
[__2] (BasicObject BasicType) "T12:00:00+00:00" string

[0] (InterfaceObject UserType) _ .error
[_1] (BasicObject BasicType) untyped string
[_1] (BasicObject BasicType) T12:00:00+00:00 
[__2] (BasicObject BasicType) inDate string
[__2] (BasicObject BasicType) "T12:00:00+00:00" string

[0] (StructObject UserType) newOutDate time.Time struct{wall uint64, ext int64, loc (*time.Location struct{name string, zone []time.zone struct{name string, offset int, isDST bool}, tx []time.zoneTrans struct{when int64, index uint8, isstd bool, isutc bool}, extend string, cacheStart int64, cacheEnd int64, cacheZone (*time.zone struct{name string, offset int, isDST bool})})}
[_1] (BasicObject BasicType) untyped string
[_1] (BasicObject BasicType) T12:00:00+00:00 
[__2] (BasicObject BasicType) outDate string
[__2] (BasicObject BasicType) "T12:00:00+00:00" string
[_1] (BasicObject BasicType) untyped string
[_1] (BasicObject BasicType) T12:00:00+00:00 
[__2] (BasicObject BasicType) outDate string
[__2] (BasicObject BasicType) "T12:00:00+00:00" string

[0] (InterfaceObject UserType) _ .error
[_1] (BasicObject BasicType) untyped string
[_1] (BasicObject BasicType) T12:00:00+00:00 
[__2] (BasicObject BasicType) outDate string
[__2] (BasicObject BasicType) "T12:00:00+00:00" string

[0] (BasicObject BasicType) hotelId string
[_1] (BasicObject BasicType) * string

[0] (BasicObject BasicType) indate string

[0] (MapObject MapType) reservation_update_map map[string]int64

[0] (StructObject UserType) newInDate time.Time struct{wall uint64, ext int64, loc (*time.Location struct{name string, zone []time.zone struct{name string, offset int, isDST bool}, tx []time.zoneTrans struct{when int64, index uint8, isstd bool, isutc bool}, extend string, cacheStart int64, cacheEnd int64, cacheZone (*time.zone struct{name string, offset int, isDST bool})})}
[_1] (BasicObject BasicType) untyped string
[_1] (BasicObject BasicType) T12:00:00+00:00 
[__2] (BasicObject BasicType) inDate string
[__2] (BasicObject BasicType) "T12:00:00+00:00" string
[_1] (BasicObject BasicType) untyped string
[_1] (BasicObject BasicType) T12:00:00+00:00 
[__2] (BasicObject BasicType) inDate string
[__2] (BasicObject BasicType) "T12:00:00+00:00" string

[0] (InterfaceObject UserType) _ .error
[_1] (BasicObject BasicType) untyped string
[_1] (BasicObject BasicType) T12:00:00+00:00 
[__2] (BasicObject BasicType) inDate string
[__2] (BasicObject BasicType) "T12:00:00+00:00" string

[0] (BasicObject BasicType) indate string

[0] (BasicObject BasicType) key int64
[_1] (MapObject MapType) reservation_update_map map[string]int64

[0] (BasicObject BasicType) val string

[0] (InterfaceObject UserType) err .error

