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

    --> w-tainted: write(fulltimes_db.Fulltime.EmployeeID) {1}
[0] (BasicObject BasicType) employeeID string

[0] (BasicObject BasicType) name string

[0] (BasicObject BasicType) IBAN string

    --> w-tainted: write(fulltimes_db.Fulltime.FulltimeID) {1}
[0] (BasicObject BasicType) fulltimeID string

    --> w-tainted: write(fulltimes_db.Fulltime.Salary) {1}
[0] (BasicObject BasicType) salary string

    --> w-tainted: write(fulltimes_db.Fulltime.Position) {1}
[0] (BasicObject BasicType) position string

[0] (StructObject UserType) fulltime employee_app.Fulltime struct{EmployeeID string, FulltimeID string, Salary string, Position string}
[_1] (Reference UserType) ref <fulltime employee_app.Fulltime struct{EmployeeID string, FulltimeID string, Salary string, Position string}> @ EmployeeService

[0] (InterfaceObject UserType) err .error
[_1] (Reference BasicType) ref <nil> @ EmployeeService

