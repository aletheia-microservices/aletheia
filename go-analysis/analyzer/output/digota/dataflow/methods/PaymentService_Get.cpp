[0] (PointerObject PointerType) s (*digota.PaymentServiceImpl struct{db NoSQLDatabase})
[_1] (StructObject UserType) digota.PaymentServiceImpl struct{db NoSQLDatabase}
[__2] (FieldObject FieldType) db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context

[0] (BasicObject BasicType) id string

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = payments, collection = payments}

[0] (InterfaceObject UserType) err .error

[0] (SliceObject UserType) query primitive.D
[_1] (StructObject StructType) struct{Key "id" string, Key "id" string, Value string, Value string}
[__2] (FieldObject FieldType) Key "id" string
[___3] (BasicObject BasicType) "id" string
[__2] (FieldObject FieldType) Value string
[___3] (BasicObject BasicType) id string

[0] (BlueprintBackendObject BlueprintBackendType) result NoSQLCursor {database = payments, collection = payments}
[_1] (StructObject UserType) digota.Charge struct{Id string, Statement string, ChargeAmount uint64, RefundAmount uint64, Refunds [](*digota.Refund struct{RefundAmount uint64, ProviderRefundId string, Reason int32, Created int64}), Currency int32, Email string, Paid bool, Refunded bool, ProviderId int32, ProviderChargeId string, Created int64, Updated int64}

[0] (InterfaceObject UserType) err .error

[0] (PointerObject PointerType) charge (*digota.Charge struct{Id string, Statement string, ChargeAmount uint64, RefundAmount uint64, Refunds [](*digota.Refund struct{RefundAmount uint64, ProviderRefundId string, Reason int32, Created int64}), Currency int32, Email string, Paid bool, Refunded bool, ProviderId int32, ProviderChargeId string, Created int64, Updated int64})
[_1] (StructObject UserType) digota.Charge struct{Id string, Statement string, ChargeAmount uint64, RefundAmount uint64, Refunds [](*digota.Refund struct{RefundAmount uint64, ProviderRefundId string, Reason int32, Created int64}), Currency int32, Email string, Paid bool, Refunded bool, ProviderId int32, ProviderChargeId string, Created int64, Updated int64}

[0] (BasicObject BasicType) found bool

[0] (InterfaceObject UserType) err .error

