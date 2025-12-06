[0] (PointerObject PointerType) s (*payment.paymentImpl struct{declineOverAmount float32})
[_1] (StructObject UserType) payment.paymentImpl struct{declineOverAmount 50 float32}
[__2] (FieldObject FieldType) declineOverAmount 50 float32
[___3] (BasicObject BasicType) 50 float32

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ OrderService
[__2] (Reference UserType) ref <ctx context.Context> @ FrontendService

    --> w-tainted: write(order_db.Order.Total) {1}
[0] (BasicObject BasicType) amount float32
     --> w-tainted: write(order_db.Order.Total) {1}
[_1] (Reference BasicType) ref <amount 10 float32> @ OrderService

[0] (StructObject UserType) authorization payment.Authorisation struct{Authorised false bool, Message "Payment declined: amount exceeds" string}
[_1] (FieldObject FieldType) Authorised false bool
[__2] (BasicObject BasicType) false bool
[_1] (FieldObject FieldType) Message "Payment declined: amount exceeds" string
[__2] (BasicObject BasicType) "Payment declined: amount exceeds" string

