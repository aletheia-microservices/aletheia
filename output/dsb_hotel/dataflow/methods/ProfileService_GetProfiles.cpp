[0] (PointerObject PointerType) p (*hotelreservation.ProfileServiceImpl struct{profileCache Cache, profileDB NoSQLDatabase})
[_1] (StructObject UserType) hotelreservation.ProfileServiceImpl struct{profileCache Cache, profileDB NoSQLDatabase}
[__2] (FieldObject FieldType) profileCache Cache
[___3] (BlueprintBackendObject BlueprintBackendType) profileCache Cache
[__2] (FieldObject FieldType) profileDB NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) profileDB NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context

[0] (ArrayObject ArrayType) hotelIds []string

[0] (BasicObject BasicType) locale string

[0] (ArrayObject ArrayType) profiles []hotelreservation.HotelProfile struct{ID string, Name string, PhoneNumber string, Description string, Address hotelreservation.Address struct{StreetNumber string, StreetName string, City string, State string, Country string, PostalCode string, Lat float64, Lon float64}}
[_1] (StructObject UserType) profile hotelreservation.HotelProfile struct{ID string, Name string, PhoneNumber string, Description string, Address hotelreservation.Address struct{StreetNumber string, StreetName string, City string, State string, Country string, PostalCode string, Lat float64, Lon float64}}

[0] (BasicObject BasicType) hid string
[_1] (ArrayObject ArrayType) hotelIds []string

[0] (StructObject UserType) profile hotelreservation.HotelProfile struct{ID string, Name string, PhoneNumber string, Description string, Address hotelreservation.Address struct{StreetNumber string, StreetName string, City string, State string, Country string, PostalCode string, Lat float64, Lon float64}}

[0] (BasicObject BasicType) exists bool

[0] (InterfaceObject UserType) err .error

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = profile-db, collection = hotels}

[0] (InterfaceObject UserType) err .error

[0] (SliceObject UserType) query primitive.D
[_1] (StructObject StructType) struct{ string,  "id" string,  string}

[0] (BlueprintBackendObject BlueprintBackendType) res NoSQLCursor {database = profile-db, collection = hotels}
[_1] (StructObject UserType) profile hotelreservation.HotelProfile struct{ID string, Name string, PhoneNumber string, Description string, Address hotelreservation.Address struct{StreetNumber string, StreetName string, City string, State string, Country string, PostalCode string, Lat float64, Lon float64}}

[0] (InterfaceObject UserType) err .error

[0] (InterfaceObject UserType) err .error

