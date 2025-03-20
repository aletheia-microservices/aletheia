[0] (PointerObject PointerType) s (*employee_app.EmployeeServiceImpl struct{freelancerService employee_app.FreelancerService, fulltimeService employee_app.FulltimeService, internService employee_app.InternService, employeesDB NoSQLDatabase})
[_1] (StructObject UserType) employee_app.EmployeeServiceImpl struct{freelancerService employee_app.FreelancerService, fulltimeService employee_app.FulltimeService, internService employee_app.InternService, employeesDB NoSQLDatabase}
[__2] (FieldObject FieldType) employeesDB NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) employeesDB NoSQLDatabase
[__2] (FieldObject FieldType) freelancerService employee_app.FreelancerService
[___3] (ServiceObject ServiceType) freelancerService employee_app.FreelancerService
[__2] (FieldObject FieldType) fulltimeService employee_app.FulltimeService
[___3] (ServiceObject ServiceType) fulltimeService employee_app.FulltimeService
[__2] (FieldObject FieldType) internService employee_app.InternService
[___3] (ServiceObject ServiceType) internService employee_app.InternService

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ Frontend

    --> w-tainted: write(employees_db.Employee.EmployeeID, fulltimes_db.Fulltime.EmployeeID) {2}
[0] (BasicObject BasicType) employeeID string
     --> w-tainted: write(employees_db.Employee.EmployeeID, fulltimes_db.Fulltime.EmployeeID) {2}
[_1] (Reference BasicType) ref <employeeID string> @ Frontend

    --> w-tainted: write(employees_db.Employee.Name) {1}
[0] (BasicObject BasicType) name string
     --> w-tainted: write(employees_db.Employee.Name) {1}
[_1] (Reference BasicType) ref <name string> @ Frontend

    --> w-tainted: write(employees_db.Employee.IBAN) {1}
[0] (BasicObject BasicType) IBAN string
     --> w-tainted: write(employees_db.Employee.IBAN) {1}
[_1] (Reference BasicType) ref <IBAN string> @ Frontend

    --> w-tainted: write(fulltimes_db.Fulltime.FulltimeID) {1}
[0] (BasicObject BasicType) fulltimeID string
     --> w-tainted: write(fulltimes_db.Fulltime.FulltimeID) {1}
[_1] (Reference BasicType) ref <fulltimeID string> @ Frontend

    --> w-tainted: write(fulltimes_db.Fulltime.Salary) {1}
[0] (BasicObject BasicType) salary string
     --> w-tainted: write(fulltimes_db.Fulltime.Salary) {1}
[_1] (Reference BasicType) ref <salary string> @ Frontend

    --> w-tainted: write(fulltimes_db.Fulltime.Position) {1}
[0] (BasicObject BasicType) position string
     --> w-tainted: write(fulltimes_db.Fulltime.Position) {1}
[_1] (Reference BasicType) ref <position string> @ Frontend

    --> w-tainted: write(employees_db.Employee) {1}
[0] (StructObject UserType) employee employee_app.Employee struct{EmployeeID string, Name string, IBAN string, SpecFlag "fulltime" string}
     --> w-tainted: write(employees_db.Employee.EmployeeID) {1}
[_1] (FieldObject FieldType) EmployeeID string
      --> w-tainted: write(employees_db.Employee.EmployeeID, fulltimes_db.Fulltime.EmployeeID) {2}
[__2] (BasicObject BasicType) employeeID string
       --> w-tainted: write(employees_db.Employee.EmployeeID, fulltimes_db.Fulltime.EmployeeID) {2}
[___3] (Reference BasicType) ref <employeeID string> @ Frontend
     --> w-tainted: write(employees_db.Employee.IBAN) {1}
[_1] (FieldObject FieldType) IBAN string
      --> w-tainted: write(employees_db.Employee.IBAN) {1}
[__2] (BasicObject BasicType) IBAN string
       --> w-tainted: write(employees_db.Employee.IBAN) {1}
[___3] (Reference BasicType) ref <IBAN string> @ Frontend
     --> w-tainted: write(employees_db.Employee.Name) {1}
[_1] (FieldObject FieldType) Name string
      --> w-tainted: write(employees_db.Employee.Name) {1}
[__2] (BasicObject BasicType) name string
       --> w-tainted: write(employees_db.Employee.Name) {1}
[___3] (Reference BasicType) ref <name string> @ Frontend
     --> w-tainted: write(employees_db.Employee.SpecFlag) {1}
[_1] (FieldObject FieldType) SpecFlag "fulltime" string
      --> w-tainted: write(employees_db.Employee.SpecFlag) {1}
[__2] (BasicObject BasicType) "fulltime" string

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = employees, collection = employees}

[0] (InterfaceObject UserType) err .error

[0] (InterfaceObject UserType) err .error

    --> w-tainted: write(fulltimes_db.Fulltime) {1}
[0] (StructObject UserType) fulltime employee_app.Fulltime struct{EmployeeID string, FulltimeID string, Salary string, Position string}
     --> w-tainted: write(fulltimes_db.Fulltime) {1}
[_1] (Reference UserType) ref <employee employee_app.Fulltime struct{EmployeeID string, FulltimeID string, Salary string, Position string}> @ FulltimeService
      --> w-tainted: write(fulltimes_db.Fulltime.EmployeeID) {1}
[__2] (FieldObject FieldType) EmployeeID string
       --> w-tainted: write(fulltimes_db.Fulltime.EmployeeID, employees_db.Employee.EmployeeID) {2}
[___3] (BasicObject BasicType) employeeID string
        --> w-tainted: write(employees_db.Employee.EmployeeID, fulltimes_db.Fulltime.EmployeeID) {2}
[____4] (Reference BasicType) ref <employeeID string> @ EmployeeService
         --> w-tainted: write(employees_db.Employee.EmployeeID, fulltimes_db.Fulltime.EmployeeID) {2}
[_____5] (Reference BasicType) ref <employeeID string> @ Frontend
        --> w-tainted: write(fulltimes_db.Fulltime.EmployeeID) {1}
[____4] (Reference BasicType) ref <employeeID string> @ EmployeeService
         --> w-tainted: write(fulltimes_db.Fulltime.EmployeeID) {1}
[_____5] (Reference BasicType) ref <employeeID string> @ Frontend
      --> w-tainted: write(fulltimes_db.Fulltime.FulltimeID) {1}
[__2] (FieldObject FieldType) FulltimeID string
       --> w-tainted: write(fulltimes_db.Fulltime.FulltimeID) {1}
[___3] (BasicObject BasicType) fulltimeID string
        --> w-tainted: write(fulltimes_db.Fulltime.FulltimeID) {1}
[____4] (Reference BasicType) ref <fulltimeID string> @ EmployeeService
         --> w-tainted: write(fulltimes_db.Fulltime.FulltimeID) {1}
[_____5] (Reference BasicType) ref <fulltimeID string> @ Frontend
        --> w-tainted: write(fulltimes_db.Fulltime.FulltimeID) {1}
[____4] (Reference BasicType) ref <fulltimeID string> @ EmployeeService
         --> w-tainted: write(fulltimes_db.Fulltime.FulltimeID) {1}
[_____5] (Reference BasicType) ref <fulltimeID string> @ Frontend
      --> w-tainted: write(fulltimes_db.Fulltime.Position) {1}
[__2] (FieldObject FieldType) Position string
       --> w-tainted: write(fulltimes_db.Fulltime.Position) {1}
[___3] (BasicObject BasicType) position string
        --> w-tainted: write(fulltimes_db.Fulltime.Position) {1}
[____4] (Reference BasicType) ref <position string> @ EmployeeService
         --> w-tainted: write(fulltimes_db.Fulltime.Position) {1}
[_____5] (Reference BasicType) ref <position string> @ Frontend
        --> w-tainted: write(fulltimes_db.Fulltime.Position) {1}
[____4] (Reference BasicType) ref <position string> @ EmployeeService
         --> w-tainted: write(fulltimes_db.Fulltime.Position) {1}
[_____5] (Reference BasicType) ref <position string> @ Frontend
      --> w-tainted: write(fulltimes_db.Fulltime.Salary) {1}
[__2] (FieldObject FieldType) Salary string
       --> w-tainted: write(fulltimes_db.Fulltime.Salary) {1}
[___3] (BasicObject BasicType) salary string
        --> w-tainted: write(fulltimes_db.Fulltime.Salary) {1}
[____4] (Reference BasicType) ref <salary string> @ EmployeeService
         --> w-tainted: write(fulltimes_db.Fulltime.Salary) {1}
[_____5] (Reference BasicType) ref <salary string> @ Frontend
        --> w-tainted: write(fulltimes_db.Fulltime.Salary) {1}
[____4] (Reference BasicType) ref <salary string> @ EmployeeService
         --> w-tainted: write(fulltimes_db.Fulltime.Salary) {1}
[_____5] (Reference BasicType) ref <salary string> @ Frontend

[0] (InterfaceObject UserType) err .error
[_1] (Reference UserType) ref <err .error> @ FulltimeService

