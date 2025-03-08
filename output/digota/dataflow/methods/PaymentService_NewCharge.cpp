[0] (PointerObject PointerType) s (*digota.PaymentServiceImpl struct{db NoSQLDatabase})
[_1] (StructObject UserType) digota.PaymentServiceImpl struct{db NoSQLDatabase}
[__2] (FieldObject FieldType) db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context

[0] (BasicObject BasicType) currency int32

[0] (BasicObject BasicType) total uint64

[0] (PointerObject PointerType) card (*digota.Card struct{Number string, ExpireMonth string, ExpireYear string, FirstName string, LastName string, CVC string, Type int32})
[_1] (StructObject UserType) digota.Card struct{Number string, ExpireMonth string, ExpireYear string, FirstName string, LastName string, CVC string, Type int32}

[0] (BasicObject BasicType) email string

[0] (BasicObject BasicType) statement string

[0] (BasicObject BasicType) paymentProviderId int32

[0] (MapObject MapType) metadata map[string]string

[0] (PointerObject PointerType) charge (*digota.Charge struct{Id string, Statement string, ChargeAmount uint64, RefundAmount uint64, Refunds [](*digota.Refund struct{RefundAmount uint64, ProviderRefundId string, Reason int32, Created int64}), Currency int32, Email string, Paid bool, Refunded bool, ProviderId int32, ProviderChargeId string, Created int64, Updated int64})
[_1] (StructObject UserType) digota.Charge struct{Id string, Statement string, ChargeAmount uint64, RefundAmount uint64, Refunds [](*digota.Refund struct{RefundAmount uint64, ProviderRefundId string, Reason int32, Created int64}), Currency int32, Email string, Paid bool, Refunded bool, ProviderId int32, ProviderChargeId string, Created int64, Updated int64}
[__2] (FieldObject FieldType) ChargeAmount uint64
[___3] (BasicObject BasicType) total uint64
[__2] (FieldObject FieldType) Currency int32
[___3] (BasicObject BasicType) currency int32
[__2] (FieldObject FieldType) Email string
[___3] (BasicObject BasicType) email string
[__2] (FieldObject FieldType) Statement string
[___3] (BasicObject BasicType) statement string

[0] (InterfaceObject UserType) err .error
[_1] (PointerObject PointerType) charge (*digota.Charge struct{Id string, Statement string, ChargeAmount uint64, RefundAmount uint64, Refunds [](*digota.Refund struct{RefundAmount uint64, ProviderRefundId string, Reason int32, Created int64}), Currency int32, Email string, Paid bool, Refunded bool, ProviderId int32, ProviderChargeId string, Created int64, Updated int64})
[__2] (StructObject UserType) digota.Charge struct{Id string, Statement string, ChargeAmount uint64, RefundAmount uint64, Refunds [](*digota.Refund struct{RefundAmount uint64, ProviderRefundId string, Reason int32, Created int64}), Currency int32, Email string, Paid bool, Refunded bool, ProviderId int32, ProviderChargeId string, Created int64, Updated int64}
[___3] (FieldObject FieldType) ChargeAmount uint64
[____4] (BasicObject BasicType) total uint64
[___3] (FieldObject FieldType) Currency int32
[____4] (BasicObject BasicType) currency int32
[___3] (FieldObject FieldType) Email string
[____4] (BasicObject BasicType) email string
[___3] (FieldObject FieldType) Statement string
[____4] (BasicObject BasicType) statement string

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = payments, collection = payments}

[0] (InterfaceObject UserType) err .error

[0] (InterfaceObject UserType) err .error

