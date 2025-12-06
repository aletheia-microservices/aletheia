[0] (PointerObject PointerType) c (*contacts.ContactsServiceImpl struct{contactsDB NoSQLDatabase})
[_1] (StructObject UserType) contacts.ContactsServiceImpl struct{contactsDB NoSQLDatabase}
[__2] (FieldObject FieldType) contactsDB NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context

    --> r-tainted: read(contacts_db._.id) {1}
[0] (BasicObject BasicType) id string

[0] (BlueprintBackendObject BlueprintBackendType) coll NoSQLCollection {database = contacts, collection = contacts}

[0] (InterfaceObject UserType) err .error

[0] (SliceObject UserType) query primitive.D
[_1] (StructObject StructType) struct{ string,  "id" string,  string}

    --> r-tainted: read(contacts_db._) {1}
[0] (BlueprintBackendObject BlueprintBackendType) res NoSQLCursor {database = contacts, collection = contacts}
     --> r-tainted: read(contacts_db.Contact) {1}
[_1] (StructObject UserType) contact contacts.Contact struct{ID string, AccountID string, Name string, DocumentType int, DocumentNumber string, PhoneNumber string}

[0] (InterfaceObject UserType) err .error

    --> r-tainted: read(contacts_db.Contact) {1}
[0] (StructObject UserType) contact contacts.Contact struct{ID string, AccountID string, Name string, DocumentType int, DocumentNumber string, PhoneNumber string}

[0] (BasicObject BasicType) exists bool

[0] (InterfaceObject UserType) err .error

