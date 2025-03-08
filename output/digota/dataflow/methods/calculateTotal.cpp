    --> w-tainted: write(orders_db.Order.Currency) {1}
[0] (BasicObject BasicType) currency int32
     --> w-tainted: write(orders_db.Order.Currency) {1}
[_1] (Reference BasicType) ref <currency int32> @ OrderService

    --> w-tainted: write(orders_db.Order.Items) {1}
[0] (ArrayObject ArrayType) orderItems [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
     --> w-tainted: write(orders_db.Order.Items) {1}
[_1] (Reference SliceType) ref <orderItems [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ OrderService
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (Reference ArrayType) ref <orderItems [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ getUpdatedOrderItems
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (PointerObject PointerType) v (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
         --> w-tainted: write(orders_db.Order.Items) {1}
[_____5] (Reference ArrayType) ref <items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ OrderService
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
         --> w-tainted: write(orders_db.Order.Items) {1}
[_____5] (FieldObject FieldType) Parent string
          --> w-tainted: write(orders_db.Order.Items) {1}
[______6] (BasicObject BasicType) Parent string
         --> w-tainted: write(orders_db.Order.Items) {1}
[_____5] (FieldObject FieldType) Quantity int64
          --> w-tainted: write(orders_db.Order.Items) {1}
[______6] (BasicObject BasicType) 1 int
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (PointerObject PointerType) v (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
         --> w-tainted: write(orders_db.Order.Items) {1}
[_____5] (Reference ArrayType) ref <items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ OrderService
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
         --> w-tainted: write(orders_db.Order.Items) {1}
[_____5] (FieldObject FieldType) Parent string
          --> w-tainted: write(orders_db.Order.Items) {1}
[______6] (BasicObject BasicType) Parent string
         --> w-tainted: write(orders_db.Order.Items) {1}
[_____5] (FieldObject FieldType) Quantity int64
          --> w-tainted: write(orders_db.Order.Items) {1}
[______6] (BasicObject BasicType) 1 int
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (PointerObject PointerType) v (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
         --> w-tainted: write(orders_db.Order.Items) {1}
[_____5] (Reference ArrayType) ref <items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ OrderService
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
         --> w-tainted: write(orders_db.Order.Items) {1}
[_____5] (FieldObject FieldType) Parent string
          --> w-tainted: write(orders_db.Order.Items) {1}
[______6] (BasicObject BasicType) Parent string
         --> w-tainted: write(orders_db.Order.Items) {1}
[_____5] (FieldObject FieldType) Quantity int64
          --> w-tainted: write(orders_db.Order.Items) {1}
[______6] (BasicObject BasicType) 1 int
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (PointerObject PointerType) v (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
         --> w-tainted: write(orders_db.Order.Items) {1}
[_____5] (Reference ArrayType) ref <items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ OrderService
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
         --> w-tainted: write(orders_db.Order.Items) {1}
[_____5] (FieldObject FieldType) Parent string
          --> w-tainted: write(orders_db.Order.Items) {1}
[______6] (BasicObject BasicType) Parent string
         --> w-tainted: write(orders_db.Order.Items) {1}
[_____5] (FieldObject FieldType) Quantity int64
          --> w-tainted: write(orders_db.Order.Items) {1}
[______6] (BasicObject BasicType) 1 int

[0] (BasicObject BasicType) err error

[0] (BasicObject BasicType) currencyString string

[0] (PointerObject PointerType) m (*go-money.Money struct{amount int64, currency (*go-money.Currency struct{Code string, NumericCode string, Fraction int, Grapheme string, Template string, Decimal string, Thousand string})})
[_1] (BasicObject BasicType) 0 int
[_1] (BasicObject BasicType) currencyString string
[_1] (BasicObject BasicType) 0 int
[_1] (BasicObject BasicType) currencyString string
[_1] (StructObject UserType) go-money.Money struct{amount int64, currency (*go-money.Currency struct{Code string, NumericCode string, Fraction int, Grapheme string, Template string, Decimal string, Thousand string})}

[0] (PointerObject PointerType) v (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
     --> w-tainted: write(orders_db.Order.Items) {1}
[_1] (ArrayObject ArrayType) orderItems [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
      --> w-tainted: write(orders_db.Order.Items) {1}
[__2] (Reference SliceType) ref <orderItems [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ OrderService
       --> w-tainted: write(orders_db.Order.Items) {1}
[___3] (Reference ArrayType) ref <orderItems [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ getUpdatedOrderItems
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (PointerObject PointerType) v (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
         --> w-tainted: write(orders_db.Order.Items) {1}
[_____5] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
          --> w-tainted: write(orders_db.Order.Items) {1}
[______6] (Reference ArrayType) ref <items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ OrderService
         --> w-tainted: write(orders_db.Order.Items) {1}
[_____5] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
          --> w-tainted: write(orders_db.Order.Items) {1}
[______6] (FieldObject FieldType) Parent string
           --> w-tainted: write(orders_db.Order.Items) {1}
[_______7] (BasicObject BasicType) Parent string
          --> w-tainted: write(orders_db.Order.Items) {1}
[______6] (FieldObject FieldType) Quantity int64
           --> w-tainted: write(orders_db.Order.Items) {1}
[_______7] (BasicObject BasicType) 1 int
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (PointerObject PointerType) v (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
         --> w-tainted: write(orders_db.Order.Items) {1}
[_____5] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
          --> w-tainted: write(orders_db.Order.Items) {1}
[______6] (Reference ArrayType) ref <items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ OrderService
         --> w-tainted: write(orders_db.Order.Items) {1}
[_____5] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
          --> w-tainted: write(orders_db.Order.Items) {1}
[______6] (FieldObject FieldType) Parent string
           --> w-tainted: write(orders_db.Order.Items) {1}
[_______7] (BasicObject BasicType) Parent string
          --> w-tainted: write(orders_db.Order.Items) {1}
[______6] (FieldObject FieldType) Quantity int64
           --> w-tainted: write(orders_db.Order.Items) {1}
[_______7] (BasicObject BasicType) 1 int
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (PointerObject PointerType) v (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
         --> w-tainted: write(orders_db.Order.Items) {1}
[_____5] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
          --> w-tainted: write(orders_db.Order.Items) {1}
[______6] (Reference ArrayType) ref <items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ OrderService
         --> w-tainted: write(orders_db.Order.Items) {1}
[_____5] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
          --> w-tainted: write(orders_db.Order.Items) {1}
[______6] (FieldObject FieldType) Parent string
           --> w-tainted: write(orders_db.Order.Items) {1}
[_______7] (BasicObject BasicType) Parent string
          --> w-tainted: write(orders_db.Order.Items) {1}
[______6] (FieldObject FieldType) Quantity int64
           --> w-tainted: write(orders_db.Order.Items) {1}
[_______7] (BasicObject BasicType) 1 int
        --> w-tainted: write(orders_db.Order.Items) {1}
[____4] (PointerObject PointerType) v (*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
         --> w-tainted: write(orders_db.Order.Items) {1}
[_____5] (ArrayObject ArrayType) items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})
          --> w-tainted: write(orders_db.Order.Items) {1}
[______6] (Reference ArrayType) ref <items [](*digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string})> @ OrderService
         --> w-tainted: write(orders_db.Order.Items) {1}
[_____5] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
          --> w-tainted: write(orders_db.Order.Items) {1}
[______6] (FieldObject FieldType) Parent string
           --> w-tainted: write(orders_db.Order.Items) {1}
[_______7] (BasicObject BasicType) Parent string
          --> w-tainted: write(orders_db.Order.Items) {1}
[______6] (FieldObject FieldType) Quantity int64
           --> w-tainted: write(orders_db.Order.Items) {1}
[_______7] (BasicObject BasicType) 1 int
[_1] (StructObject UserType) digota.OrderItem struct{Type int32, Quantity int64, Amount int64, Currency int32, Parent string, Description string}
[__2] (FieldObject FieldType) Currency int32
[___3] (BasicObject BasicType) Currency int32
[__2] (FieldObject FieldType) Quantity int64
[___3] (BasicObject BasicType) 1 int

[0] (BasicObject BasicType) vCurrencyString string

[0] (PointerObject PointerType) m (*go-money.Money struct{amount int64, currency (*go-money.Currency struct{Code string, NumericCode string, Fraction int, Grapheme string, Template string, Decimal string, Thousand string})})
[_1] (PointerObject PointerType) (*go-money.Money struct{amount int64, currency (*go-money.Currency struct{Code string, NumericCode string, Fraction int, Grapheme string, Template string, Decimal string, Thousand string})})
[__2] (BasicObject BasicType) int64
[__2] (BasicObject BasicType) vCurrencyString string
[__2] (BasicObject BasicType) int64
[__2] (BasicObject BasicType) vCurrencyString string
[__2] (StructObject UserType) go-money.Money struct{amount int64, currency (*go-money.Currency struct{Code string, NumericCode string, Fraction int, Grapheme string, Template string, Decimal string, Thousand string})}
[_1] (PointerObject PointerType) (*go-money.Money struct{amount int64, currency (*go-money.Currency struct{Code string, NumericCode string, Fraction int, Grapheme string, Template string, Decimal string, Thousand string})})
[__2] (BasicObject BasicType) int64
[__2] (BasicObject BasicType) vCurrencyString string
[__2] (BasicObject BasicType) int64
[__2] (BasicObject BasicType) vCurrencyString string
[__2] (StructObject UserType) go-money.Money struct{amount int64, currency (*go-money.Currency struct{Code string, NumericCode string, Fraction int, Grapheme string, Template string, Decimal string, Thousand string})}
[_1] (StructObject UserType) go-money.Money struct{amount int64, currency (*go-money.Currency struct{Code string, NumericCode string, Fraction int, Grapheme string, Template string, Decimal string, Thousand string})}

[0] (InterfaceObject UserType) err .error
[_1] (PointerObject PointerType) (*go-money.Money struct{amount int64, currency (*go-money.Currency struct{Code string, NumericCode string, Fraction int, Grapheme string, Template string, Decimal string, Thousand string})})
[__2] (BasicObject BasicType) int64
[__2] (BasicObject BasicType) vCurrencyString string
[__2] (BasicObject BasicType) int64
[__2] (BasicObject BasicType) vCurrencyString string
[__2] (StructObject UserType) go-money.Money struct{amount int64, currency (*go-money.Currency struct{Code string, NumericCode string, Fraction int, Grapheme string, Template string, Decimal string, Thousand string})}

