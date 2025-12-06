    --> w-tainted: write(orders_db.Order.Currency) {1}
[0] (BasicObject BasicType) currency int32
     --> w-tainted: write(orders_db.Order.Currency) {1}
[_1] (Reference BasicType) ref <currency int32> @ OrderService

    --> w-tainted: write(orders_db.Order.Items) {1}
[0] (ArrayObject ArrayType) orderItems [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
     --> w-tainted: write(orders_db.Order.Items) {1}
[_1] (Reference ArrayType) ref <items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ OrderService

[0] (BasicObject BasicType) err error

[0] (BasicObject BasicType) currencyString string
[_1] (BasicObject BasicType) * string

[0] (PointerObject PointerType) m (*go-money.Money struct{amount int64, currency (*go-money.Currency struct{Code string, NumericCode string, Fraction int, Grapheme string, Template string, Decimal string, Thousand string})})
[_1] (BasicObject BasicType) 0 int
[_1] (BasicObject BasicType) currencyString string
[__2] (BasicObject BasicType) * string
[_1] (BasicObject BasicType) 0 int
[_1] (BasicObject BasicType) currencyString string
[__2] (BasicObject BasicType) * string
[_1] (StructObject UserType) go-money.Money struct{amount int64, currency (*go-money.Currency struct{Code string, NumericCode string, Fraction int, Grapheme string, Template string, Decimal string, Thousand string})}

[0] (PointerObject PointerType) v (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
     --> w-tainted: write(orders_db.Order.Items) {1}
[_1] (ArrayObject ArrayType) orderItems [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (Reference ArrayType) ref <items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ OrderService
[_1] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
[__2] (FieldObject FieldType) Currency int32
[___3] (BasicObject BasicType) Currency int32
[__2] (FieldObject FieldType) Quantity int64
[___3] (BasicObject BasicType) 1 int

[0] (BasicObject BasicType) vCurrencyString string
[_1] (BasicObject BasicType) * string

[0] (PointerObject PointerType) m (*go-money.Money struct{amount int64, currency (*go-money.Currency struct{Code string, NumericCode string, Fraction int, Grapheme string, Template string, Decimal string, Thousand string})})
[_1] (PointerObject PointerType) (*go-money.Money struct{amount int64, currency (*go-money.Currency struct{Code string, NumericCode string, Fraction int, Grapheme string, Template string, Decimal string, Thousand string})})
[__2] (BasicObject BasicType) int64
[__2] (BasicObject BasicType) vCurrencyString string
[___3] (BasicObject BasicType) * string
[__2] (BasicObject BasicType) int64
[__2] (BasicObject BasicType) vCurrencyString string
[___3] (BasicObject BasicType) * string
[__2] (StructObject UserType) go-money.Money struct{amount int64, currency (*go-money.Currency struct{Code string, NumericCode string, Fraction int, Grapheme string, Template string, Decimal string, Thousand string})}
[_1] (PointerObject PointerType) (*go-money.Money struct{amount int64, currency (*go-money.Currency struct{Code string, NumericCode string, Fraction int, Grapheme string, Template string, Decimal string, Thousand string})})
[__2] (BasicObject BasicType) int64
[__2] (BasicObject BasicType) vCurrencyString string
[___3] (BasicObject BasicType) * string
[__2] (BasicObject BasicType) int64
[__2] (BasicObject BasicType) vCurrencyString string
[___3] (BasicObject BasicType) * string
[__2] (StructObject UserType) go-money.Money struct{amount int64, currency (*go-money.Currency struct{Code string, NumericCode string, Fraction int, Grapheme string, Template string, Decimal string, Thousand string})}
[_1] (StructObject UserType) go-money.Money struct{amount int64, currency (*go-money.Currency struct{Code string, NumericCode string, Fraction int, Grapheme string, Template string, Decimal string, Thousand string})}

[0] (InterfaceObject UserType) err .error
[_1] (PointerObject PointerType) (*go-money.Money struct{amount int64, currency (*go-money.Currency struct{Code string, NumericCode string, Fraction int, Grapheme string, Template string, Decimal string, Thousand string})})
[__2] (BasicObject BasicType) int64
[__2] (BasicObject BasicType) vCurrencyString string
[___3] (BasicObject BasicType) * string
[__2] (BasicObject BasicType) int64
[__2] (BasicObject BasicType) vCurrencyString string
[___3] (BasicObject BasicType) * string
[__2] (StructObject UserType) go-money.Money struct{amount int64, currency (*go-money.Currency struct{Code string, NumericCode string, Fraction int, Grapheme string, Template string, Decimal string, Thousand string})}

