[0] (PointerObject PointerType) u (*employee_app.FrontendImpl struct{employeeService employee_app.EmployeeService, freelancerService employee_app.FreelancerService, fulltimeService employee_app.FulltimeService, internService employee_app.InternService})
[_1] (StructObject UserType) employee_app.FrontendImpl struct{employeeService employee_app.EmployeeService, freelancerService employee_app.FreelancerService, fulltimeService employee_app.FulltimeService, internService employee_app.InternService}
[__2] (FieldObject FieldType) employeeService employee_app.EmployeeService
[___3] (ServiceObject ServiceType) employeeService employee_app.EmployeeService
[__2] (FieldObject FieldType) freelancerService employee_app.FreelancerService
[___3] (ServiceObject ServiceType) freelancerService employee_app.FreelancerService
[__2] (FieldObject FieldType) fulltimeService employee_app.FulltimeService
[___3] (ServiceObject ServiceType) fulltimeService employee_app.FulltimeService
[__2] (FieldObject FieldType) internService employee_app.InternService
[___3] (ServiceObject ServiceType) internService employee_app.InternService

[0] (InterfaceObject UserType) ctx context.Context

    --> w-tainted: write(employees_db.Employee.EmployeeID, fulltimes_db.Fulltime.EmployeeID) {2}
[0] (BasicObject BasicType) employeeID string

    --> w-tainted: write(employees_db.Employee.Name) {1}
[0] (BasicObject BasicType) name string

    --> w-tainted: write(employees_db.Employee.IBAN) {1}
[0] (BasicObject BasicType) IBAN string

    --> w-tainted: write(fulltimes_db.Fulltime.FulltimeID) {1}
[0] (BasicObject BasicType) fulltimeID string

    --> w-tainted: write(fulltimes_db.Fulltime.Salary) {1}
[0] (BasicObject BasicType) salary string

    --> w-tainted: write(fulltimes_db.Fulltime.Position) {1}
[0] (BasicObject BasicType) position string

    --> w-tainted: write(employees_db.Employee) {1}
[0] (StructObject UserType) employee employee_app.Employee struct{EmployeeID string, Name string, IBAN string, SpecFlag string}
     --> w-tainted: write(employees_db.Employee) {1}
[_1] (Reference UserType) ref <employee employee_app.Employee struct{EmployeeID string, Name string, IBAN string, SpecFlag "fulltime" string}> @ EmployeeService
      --> w-tainted: write(employees_db.Employee.EmployeeID) {1}
[__2] (FieldObject FieldType) EmployeeID string
       --> w-tainted: write(employees_db.Employee.EmployeeID, fulltimes_db.Fulltime.EmployeeID) {2}
[___3] (BasicObject BasicType) employeeID string
        --> w-tainted: write(employees_db.Employee.EmployeeID, fulltimes_db.Fulltime.EmployeeID) {2}
[____4] (Reference BasicType) ref <employeeID string> @ Frontend
      --> w-tainted: write(employees_db.Employee.IBAN) {1}
[__2] (FieldObject FieldType) IBAN string
       --> w-tainted: write(employees_db.Employee.IBAN) {1}
[___3] (BasicObject BasicType) IBAN string
        --> w-tainted: write(employees_db.Employee.IBAN) {1}
[____4] (Reference BasicType) ref <IBAN string> @ Frontend
      --> w-tainted: write(employees_db.Employee.Name) {1}
[__2] (FieldObject FieldType) Name string
       --> w-tainted: write(employees_db.Employee.Name) {1}
[___3] (BasicObject BasicType) name string
        --> w-tainted: write(employees_db.Employee.Name) {1}
[____4] (Reference BasicType) ref <name string> @ Frontend
      --> w-tainted: write(employees_db.Employee.SpecFlag) {1}
[__2] (FieldObject FieldType) SpecFlag "fulltime" string
       --> w-tainted: write(employees_db.Employee.SpecFlag) {1}
[___3] (BasicObject BasicType) "fulltime" string

    --> w-tainted: write(fulltimes_db.Fulltime) {1}
[0] (StructObject UserType) fulltime employee_app.Fulltime struct{EmployeeID string, FulltimeID string, Salary string, Position string}
     --> w-tainted: write(fulltimes_db.Fulltime) {1}
[_1] (Reference UserType) ref <fulltime employee_app.Fulltime struct{EmployeeID string, FulltimeID string, Salary string, Position string}> @ EmployeeService
      --> w-tainted: write(fulltimes_db.Fulltime) {1}
[__2] (Reference UserType) ref <employee employee_app.Fulltime struct{EmployeeID string, FulltimeID string, Salary string, Position string}> @ FulltimeService
       --> w-tainted: write(fulltimes_db.Fulltime.EmployeeID) {1}
[___3] (FieldObject FieldType) EmployeeID string
        --> w-tainted: write(fulltimes_db.Fulltime.EmployeeID, employees_db.Employee.EmployeeID) {2}
[____4] (BasicObject BasicType) employeeID string
         --> w-tainted: write(employees_db.Employee.EmployeeID, fulltimes_db.Fulltime.EmployeeID) {2}
[_____5] (Reference BasicType) ref <employeeID string> @ EmployeeService
          --> w-tainted: write(employees_db.Employee.EmployeeID, fulltimes_db.Fulltime.EmployeeID) {2}
[______6] (Reference BasicType) ref <employeeID string> @ Frontend
         --> w-tainted: write(fulltimes_db.Fulltime.EmployeeID) {1}
[_____5] (Reference BasicType) ref <employeeID string> @ EmployeeService
          --> w-tainted: write(fulltimes_db.Fulltime.EmployeeID) {1}
[______6] (Reference BasicType) ref <employeeID string> @ Frontend
       --> w-tainted: write(fulltimes_db.Fulltime.FulltimeID) {1}
[___3] (FieldObject FieldType) FulltimeID string
        --> w-tainted: write(fulltimes_db.Fulltime.FulltimeID) {1}
[____4] (BasicObject BasicType) fulltimeID string
         --> w-tainted: write(fulltimes_db.Fulltime.FulltimeID) {1}
[_____5] (Reference BasicType) ref <fulltimeID string> @ EmployeeService
          --> w-tainted: write(fulltimes_db.Fulltime.FulltimeID) {1}
[______6] (Reference BasicType) ref <fulltimeID string> @ Frontend
         --> w-tainted: write(fulltimes_db.Fulltime.FulltimeID) {1}
[_____5] (Reference BasicType) ref <fulltimeID string> @ EmployeeService
          --> w-tainted: write(fulltimes_db.Fulltime.FulltimeID) {1}
[______6] (Reference BasicType) ref <fulltimeID string> @ Frontend
       --> w-tainted: write(fulltimes_db.Fulltime.Position) {1}
[___3] (FieldObject FieldType) Position string
        --> w-tainted: write(fulltimes_db.Fulltime.Position) {1}
[____4] (BasicObject BasicType) position string
         --> w-tainted: write(fulltimes_db.Fulltime.Position) {1}
[_____5] (Reference BasicType) ref <position string> @ EmployeeService
          --> w-tainted: write(fulltimes_db.Fulltime.Position) {1}
[______6] (Reference BasicType) ref <position string> @ Frontend
         --> w-tainted: write(fulltimes_db.Fulltime.Position) {1}
[_____5] (Reference BasicType) ref <position string> @ EmployeeService
          --> w-tainted: write(fulltimes_db.Fulltime.Position) {1}
[______6] (Reference BasicType) ref <position string> @ Frontend
       --> w-tainted: write(fulltimes_db.Fulltime.Salary) {1}
[___3] (FieldObject FieldType) Salary string
        --> w-tainted: write(fulltimes_db.Fulltime.Salary) {1}
[____4] (BasicObject BasicType) salary string
         --> w-tainted: write(fulltimes_db.Fulltime.Salary) {1}
[_____5] (Reference BasicType) ref <salary string> @ EmployeeService
          --> w-tainted: write(fulltimes_db.Fulltime.Salary) {1}
[______6] (Reference BasicType) ref <salary string> @ Frontend
         --> w-tainted: write(fulltimes_db.Fulltime.Salary) {1}
[_____5] (Reference BasicType) ref <salary string> @ EmployeeService
          --> w-tainted: write(fulltimes_db.Fulltime.Salary) {1}
[______6] (Reference BasicType) ref <salary string> @ Frontend

[0] (InterfaceObject UserType) err .error
[_1] (Reference BasicType) ref <nil> @ EmployeeService

