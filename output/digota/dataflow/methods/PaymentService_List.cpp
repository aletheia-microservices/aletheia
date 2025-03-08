[0] (PointerObject PointerType) s (*digota.PaymentServiceImpl struct{db NoSQLDatabase})
[_1] (StructObject UserType) digota.PaymentServiceImpl struct{db NoSQLDatabase}
[__2] (FieldObject FieldType) db NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) db NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context

[0] (BasicObject BasicType) page int64

[0] (BasicObject BasicType) limit int64

[0] (BasicObject BasicType) sort int32

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = payments, collection = payments}

[0] (InterfaceObject UserType) err .error

[0] (BlueprintBackendObject BlueprintBackendType) result NoSQLCursor {database = payments, collection = payments}

[0] (InterfaceObject UserType) err .error

[0] (ArrayObject ArrayType) charges [](*digota.Charge struct{Id string, Statement string, ChargeAmount uint64, RefundAmount uint64, Refunds [](*digota.Refund struct{RefundAmount uint64, ProviderRefundId string, Reason int32, Created int64}), Currency int32, Email string, Paid bool, Refunded bool, ProviderId int32, ProviderChargeId string, Created int64, Updated int64})

[0] (InterfaceObject UserType) err .error

[0] (PointerObject PointerType) chargeList (*digota.ChargeList struct{Charges [](*digota.Charge struct{Id string, Statement string, ChargeAmount uint64, RefundAmount uint64, Refunds [](*digota.Refund struct{RefundAmount uint64, ProviderRefundId string, Reason int32, Created int64}), Currency int32, Email string, Paid bool, Refunded bool, ProviderId int32, ProviderChargeId string, Created int64, Updated int64}), Total len(charges [](*digota.Charge struct{Id string, Statement string, ChargeAmount uint64, RefundAmount uint64, Refunds [](*digota.Refund struct{RefundAmount uint64, ProviderRefundId string, Reason int32, Created int64}), Currency int32, Email string, Paid bool, Refunded bool, ProviderId int32, ProviderChargeId string, Created int64, Updated int64})) int32})
[_1] (StructObject UserType) digota.ChargeList struct{Charges [](*digota.Charge struct{Id string, Statement string, ChargeAmount uint64, RefundAmount uint64, Refunds [](*digota.Refund struct{RefundAmount uint64, ProviderRefundId string, Reason int32, Created int64}), Currency int32, Email string, Paid bool, Refunded bool, ProviderId int32, ProviderChargeId string, Created int64, Updated int64}), Total len(charges [](*digota.Charge struct{Id string, Statement string, ChargeAmount uint64, RefundAmount uint64, Refunds [](*digota.Refund struct{RefundAmount uint64, ProviderRefundId string, Reason int32, Created int64}), Currency int32, Email string, Paid bool, Refunded bool, ProviderId int32, ProviderChargeId string, Created int64, Updated int64})) int32}
[__2] (FieldObject FieldType) Charges [](*digota.Charge struct{Id string, Statement string, ChargeAmount uint64, RefundAmount uint64, Refunds [](*digota.Refund struct{RefundAmount uint64, ProviderRefundId string, Reason int32, Created int64}), Currency int32, Email string, Paid bool, Refunded bool, ProviderId int32, ProviderChargeId string, Created int64, Updated int64})
[___3] (ArrayObject ArrayType) charges [](*digota.Charge struct{Id string, Statement string, ChargeAmount uint64, RefundAmount uint64, Refunds [](*digota.Refund struct{RefundAmount uint64, ProviderRefundId string, Reason int32, Created int64}), Currency int32, Email string, Paid bool, Refunded bool, ProviderId int32, ProviderChargeId string, Created int64, Updated int64})
[__2] (FieldObject FieldType) Total len(charges [](*digota.Charge struct{Id string, Statement string, ChargeAmount uint64, RefundAmount uint64, Refunds [](*digota.Refund struct{RefundAmount uint64, ProviderRefundId string, Reason int32, Created int64}), Currency int32, Email string, Paid bool, Refunded bool, ProviderId int32, ProviderChargeId string, Created int64, Updated int64})) int32
[___3] (BasicObject BasicType) len(charges [](*digota.Charge struct{Id string, Statement string, ChargeAmount uint64, RefundAmount uint64, Refunds [](*digota.Refund struct{RefundAmount uint64, ProviderRefundId string, Reason int32, Created int64}), Currency int32, Email string, Paid bool, Refunded bool, ProviderId int32, ProviderChargeId string, Created int64, Updated int64})) int32

