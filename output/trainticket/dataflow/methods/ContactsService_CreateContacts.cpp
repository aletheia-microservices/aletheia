[0] (PointerObject PointerType) c (*contacts.ContactsServiceImpl struct{contactsDB NoSQLDatabase})
[_1] (StructObject UserType) contacts.ContactsServiceImpl struct{contactsDB NoSQLDatabase}
[__2] (FieldObject FieldType) contactsDB NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context

    --> w-tainted: write(contacts_db.Contact) {1}
[0] (StructObject UserType) contact contacts.Contact struct{ID string, AccountID string, Name string, DocumentType int, DocumentNumber string, PhoneNumber string}
     --> w-tainted: write(contacts_db.Contact.AccountID) {1}
[_1] (FieldObject FieldType) AccountID string
      --> w-tainted: write(contacts_db.Contact.AccountID) {1}
[__2] (BasicObject BasicType) userID string
     --> w-tainted: write(contacts_db.Contact.DocumentNumber) {1}
[_1] (FieldObject FieldType) DocumentNumber string
      --> w-tainted: write(contacts_db.Contact.DocumentNumber) {1}           --> w-tainted: write(contacts_db.Contact.DocumentNumber) {1} --> r-tainted: read(contacts_db._.documentnumber) {1}
[__2] (BasicObject BasicType) DocumentNumber string
     --> w-tainted: write(contacts_db.Contact.DocumentType) {1}
[_1] (FieldObject FieldType) DocumentType int
      --> w-tainted: write(contacts_db.Contact.DocumentType) {1}           --> w-tainted: write(contacts_db.Contact.DocumentType) {1} --> r-tainted: read(contacts_db._.documenttype) {1}
[__2] (BasicObject BasicType) DocumentType int

    --> w-tainted: write(contacts_db.Contact.AccountID) {1}
[0] (BasicObject BasicType) userID string

[0] (BlueprintBackendObject BlueprintBackendType) coll NoSQLCollection {database = contacts, collection = contacts}

[0] (InterfaceObject UserType) err .error

[0] (SliceObject UserType) query primitive.D
[_1] (StructObject StructType) struct{ string,  "accountid" string,  string}
[_1] (StructObject StructType) struct{ string,  "documentnumber" string,  string}
[_1] (StructObject StructType) struct{ int,  "documenttype" string,  int}

    --> r-tainted: read(contacts_db._) {1}
[0] (BlueprintBackendObject BlueprintBackendType) res NoSQLCursor {database = contacts, collection = contacts}

[0] (InterfaceObject UserType) err .error

    --> r-tainted: read(contacts_db.Contact) {1}
[0] (StructObject UserType) existing contacts.Contact struct{ID string, AccountID string, Name string, DocumentType int, DocumentNumber string, PhoneNumber string}

[0] (BasicObject BasicType) exists bool

[0] (InterfaceObject UserType) err .error

