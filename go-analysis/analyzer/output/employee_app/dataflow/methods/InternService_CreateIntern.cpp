[0] (PointerObject PointerType) s (*employee_app.InternServiceImpl struct{internsDB NoSQLDatabase})
[_1] (StructObject UserType) employee_app.InternServiceImpl struct{internsDB NoSQLDatabase}
[__2] (FieldObject FieldType) internsDB NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) internsDB NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ EmployeeService
[__2] (Reference UserType) ref <ctx context.Context> @ Frontend

    --> w-tainted: write(interns_db.Intern.EmployeeID, employees_db.Employee.EmployeeID) {2}
[0] (BasicObject BasicType) employeeID string
     --> w-tainted: write(employees_db.Employee.EmployeeID, interns_db.Intern.EmployeeID) {2}
[_1] (Reference BasicType) ref <employeeID string> @ EmployeeService
      --> w-tainted: write(employees_db.Employee.EmployeeID, interns_db.Intern.EmployeeID) {2}
[__2] (Reference BasicType) ref <employeeID string> @ Frontend

    --> w-tainted: write(interns_db.Intern.InternID) {1}
[0] (BasicObject BasicType) internID string
     --> w-tainted: write(interns_db.Intern.InternID) {1}
[_1] (Reference BasicType) ref <internID string> @ EmployeeService
      --> w-tainted: write(interns_db.Intern.InternID) {1}
[__2] (Reference BasicType) ref <internID string> @ Frontend

    --> w-tainted: write(interns_db.Intern.Mentor) {1}
[0] (BasicObject BasicType) mentorID string
     --> w-tainted: write(interns_db.Intern.Mentor) {1}
[_1] (Reference BasicType) ref <mentorID string> @ EmployeeService
      --> w-tainted: write(interns_db.Intern.Mentor) {1}
[__2] (Reference BasicType) ref <mentorID string> @ Frontend

    --> w-tainted: write(interns_db.Intern.Stipend) {1}
[0] (BasicObject BasicType) stipend string
     --> w-tainted: write(interns_db.Intern.Stipend) {1}
[_1] (Reference BasicType) ref <stipend string> @ EmployeeService
      --> w-tainted: write(interns_db.Intern.Stipend) {1}
[__2] (Reference BasicType) ref <stipend string> @ Frontend

    --> w-tainted: write(interns_db.Intern.Duration) {1}
[0] (BasicObject BasicType) duration string
     --> w-tainted: write(interns_db.Intern.Duration) {1}
[_1] (Reference BasicType) ref <duration string> @ EmployeeService
      --> w-tainted: write(interns_db.Intern.Duration) {1}
[__2] (Reference BasicType) ref <duration string> @ Frontend

    --> w-tainted: write(interns_db.Intern) {1}
[0] (StructObject UserType) employee employee_app.Intern struct{EmployeeID string, InternID string, Stipend string, Duration string, Mentor string}
     --> w-tainted: write(interns_db.Intern.Duration) {1}
[_1] (FieldObject FieldType) Duration string
      --> w-tainted: write(interns_db.Intern.Duration) {1}
[__2] (BasicObject BasicType) duration string
       --> w-tainted: write(interns_db.Intern.Duration) {1}
[___3] (Reference BasicType) ref <duration string> @ EmployeeService
        --> w-tainted: write(interns_db.Intern.Duration) {1}
[____4] (Reference BasicType) ref <duration string> @ Frontend
     --> w-tainted: write(interns_db.Intern.EmployeeID) {1}
[_1] (FieldObject FieldType) EmployeeID string
      --> w-tainted: write(interns_db.Intern.EmployeeID, employees_db.Employee.EmployeeID) {2}
[__2] (BasicObject BasicType) employeeID string
       --> w-tainted: write(employees_db.Employee.EmployeeID, interns_db.Intern.EmployeeID) {2}
[___3] (Reference BasicType) ref <employeeID string> @ EmployeeService
        --> w-tainted: write(employees_db.Employee.EmployeeID, interns_db.Intern.EmployeeID) {2}
[____4] (Reference BasicType) ref <employeeID string> @ Frontend
     --> w-tainted: write(interns_db.Intern.InternID) {1}
[_1] (FieldObject FieldType) InternID string
      --> w-tainted: write(interns_db.Intern.InternID) {1}
[__2] (BasicObject BasicType) internID string
       --> w-tainted: write(interns_db.Intern.InternID) {1}
[___3] (Reference BasicType) ref <internID string> @ EmployeeService
        --> w-tainted: write(interns_db.Intern.InternID) {1}
[____4] (Reference BasicType) ref <internID string> @ Frontend
     --> w-tainted: write(interns_db.Intern.Mentor) {1}
[_1] (FieldObject FieldType) Mentor string
      --> w-tainted: write(interns_db.Intern.Mentor) {1}
[__2] (BasicObject BasicType) mentorID string
       --> w-tainted: write(interns_db.Intern.Mentor) {1}
[___3] (Reference BasicType) ref <mentorID string> @ EmployeeService
        --> w-tainted: write(interns_db.Intern.Mentor) {1}
[____4] (Reference BasicType) ref <mentorID string> @ Frontend
     --> w-tainted: write(interns_db.Intern.Stipend) {1}
[_1] (FieldObject FieldType) Stipend string
      --> w-tainted: write(interns_db.Intern.Stipend) {1}
[__2] (BasicObject BasicType) stipend string
       --> w-tainted: write(interns_db.Intern.Stipend) {1}
[___3] (Reference BasicType) ref <stipend string> @ EmployeeService
        --> w-tainted: write(interns_db.Intern.Stipend) {1}
[____4] (Reference BasicType) ref <stipend string> @ Frontend

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = interns, collection = interns}

[0] (InterfaceObject UserType) err .error

[0] (InterfaceObject UserType) err .error

