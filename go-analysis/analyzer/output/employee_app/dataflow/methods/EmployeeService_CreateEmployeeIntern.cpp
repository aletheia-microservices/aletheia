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

    --> w-tainted: write(employees_db.Employee.EmployeeID, interns_db.Intern.EmployeeID) {2}
[0] (BasicObject BasicType) employeeID string
     --> w-tainted: write(employees_db.Employee.EmployeeID, interns_db.Intern.EmployeeID) {2}
[_1] (Reference BasicType) ref <employeeID string> @ Frontend

    --> w-tainted: write(employees_db.Employee.Name) {1}
[0] (BasicObject BasicType) name string
     --> w-tainted: write(employees_db.Employee.Name) {1}
[_1] (Reference BasicType) ref <name string> @ Frontend

    --> w-tainted: write(employees_db.Employee.IBAN) {1}
[0] (BasicObject BasicType) IBAN string
     --> w-tainted: write(employees_db.Employee.IBAN) {1}
[_1] (Reference BasicType) ref <IBAN string> @ Frontend

    --> w-tainted: write(interns_db.Intern.InternID) {1}
[0] (BasicObject BasicType) internID string
     --> w-tainted: write(interns_db.Intern.InternID) {1}
[_1] (Reference BasicType) ref <internID string> @ Frontend

    --> w-tainted: write(interns_db.Intern.Mentor) {1}
[0] (BasicObject BasicType) mentorID string
     --> w-tainted: write(interns_db.Intern.Mentor) {1}
[_1] (Reference BasicType) ref <mentorID string> @ Frontend

    --> w-tainted: write(interns_db.Intern.Stipend) {1}
[0] (BasicObject BasicType) stipend string
     --> w-tainted: write(interns_db.Intern.Stipend) {1}
[_1] (Reference BasicType) ref <stipend string> @ Frontend

    --> w-tainted: write(interns_db.Intern.Duration) {1}
[0] (BasicObject BasicType) duration string
     --> w-tainted: write(interns_db.Intern.Duration) {1}
[_1] (Reference BasicType) ref <duration string> @ Frontend

    --> w-tainted: write(employees_db.Employee) {1}
[0] (StructObject UserType) employee employee_app.Employee struct{EmployeeID string, Name string, IBAN string, SpecFlag "intern" string}
     --> w-tainted: write(employees_db.Employee.EmployeeID) {1}
[_1] (FieldObject FieldType) EmployeeID string
      --> w-tainted: write(employees_db.Employee.EmployeeID, interns_db.Intern.EmployeeID) {2}
[__2] (BasicObject BasicType) employeeID string
       --> w-tainted: write(employees_db.Employee.EmployeeID, interns_db.Intern.EmployeeID) {2}
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
[_1] (FieldObject FieldType) SpecFlag "intern" string
      --> w-tainted: write(employees_db.Employee.SpecFlag) {1}
[__2] (BasicObject BasicType) "intern" string

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = employees, collection = employees}

[0] (InterfaceObject UserType) err .error

[0] (InterfaceObject UserType) err .error

    --> w-tainted: write(interns_db.Intern) {1}
[0] (StructObject UserType) intern employee_app.Intern struct{EmployeeID string, InternID string, Stipend string, Duration string, Mentor string}
     --> w-tainted: write(interns_db.Intern) {1}
[_1] (Reference UserType) ref <employee employee_app.Intern struct{EmployeeID string, InternID string, Stipend string, Duration string, Mentor string}> @ InternService
      --> w-tainted: write(interns_db.Intern.Duration) {1}
[__2] (FieldObject FieldType) Duration string
       --> w-tainted: write(interns_db.Intern.Duration) {1}
[___3] (BasicObject BasicType) duration string
        --> w-tainted: write(interns_db.Intern.Duration) {1}
[____4] (Reference BasicType) ref <duration string> @ EmployeeService
         --> w-tainted: write(interns_db.Intern.Duration) {1}
[_____5] (Reference BasicType) ref <duration string> @ Frontend
      --> w-tainted: write(interns_db.Intern.EmployeeID) {1}
[__2] (FieldObject FieldType) EmployeeID string
       --> w-tainted: write(interns_db.Intern.EmployeeID, employees_db.Employee.EmployeeID) {2}
[___3] (BasicObject BasicType) employeeID string
        --> w-tainted: write(employees_db.Employee.EmployeeID, interns_db.Intern.EmployeeID) {2}
[____4] (Reference BasicType) ref <employeeID string> @ EmployeeService
         --> w-tainted: write(employees_db.Employee.EmployeeID, interns_db.Intern.EmployeeID) {2}
[_____5] (Reference BasicType) ref <employeeID string> @ Frontend
      --> w-tainted: write(interns_db.Intern.InternID) {1}
[__2] (FieldObject FieldType) InternID string
       --> w-tainted: write(interns_db.Intern.InternID) {1}
[___3] (BasicObject BasicType) internID string
        --> w-tainted: write(interns_db.Intern.InternID) {1}
[____4] (Reference BasicType) ref <internID string> @ EmployeeService
         --> w-tainted: write(interns_db.Intern.InternID) {1}
[_____5] (Reference BasicType) ref <internID string> @ Frontend
      --> w-tainted: write(interns_db.Intern.Mentor) {1}
[__2] (FieldObject FieldType) Mentor string
       --> w-tainted: write(interns_db.Intern.Mentor) {1}
[___3] (BasicObject BasicType) mentorID string
        --> w-tainted: write(interns_db.Intern.Mentor) {1}
[____4] (Reference BasicType) ref <mentorID string> @ EmployeeService
         --> w-tainted: write(interns_db.Intern.Mentor) {1}
[_____5] (Reference BasicType) ref <mentorID string> @ Frontend
      --> w-tainted: write(interns_db.Intern.Stipend) {1}
[__2] (FieldObject FieldType) Stipend string
       --> w-tainted: write(interns_db.Intern.Stipend) {1}
[___3] (BasicObject BasicType) stipend string
        --> w-tainted: write(interns_db.Intern.Stipend) {1}
[____4] (Reference BasicType) ref <stipend string> @ EmployeeService
         --> w-tainted: write(interns_db.Intern.Stipend) {1}
[_____5] (Reference BasicType) ref <stipend string> @ Frontend

[0] (InterfaceObject UserType) err .error
[_1] (Reference UserType) ref <err .error> @ InternService

