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

    --> w-tainted: write(employees_db.Employee.EmployeeID, freelancers_db.Freelancer.EmployeeID) {2}
[0] (BasicObject BasicType) employeeID string

    --> w-tainted: write(employees_db.Employee.Name) {1}
[0] (BasicObject BasicType) name string

    --> w-tainted: write(employees_db.Employee.IBAN) {1}
[0] (BasicObject BasicType) IBAN string

    --> w-tainted: write(freelancers_db.Freelancer.FreelancerID) {1}
[0] (BasicObject BasicType) freelancerID string

    --> w-tainted: write(freelancers_db.Freelancer.Rate) {1}
[0] (BasicObject BasicType) rate string

    --> w-tainted: write(freelancers_db.Freelancer.Terms) {1}
[0] (BasicObject BasicType) terms string

    --> w-tainted: write(employees_db.Employee) {1}
[0] (StructObject UserType) employee employee_app.Employee struct{EmployeeID string, Name string, IBAN string, SpecFlag string}
     --> w-tainted: write(employees_db.Employee) {1}
[_1] (Reference UserType) ref <employee employee_app.Employee struct{EmployeeID string, Name string, IBAN string, SpecFlag "freelancer" string}> @ EmployeeService
      --> w-tainted: write(employees_db.Employee.EmployeeID) {1}
[__2] (FieldObject FieldType) EmployeeID string
       --> w-tainted: write(employees_db.Employee.EmployeeID, freelancers_db.Freelancer.EmployeeID) {2}
[___3] (BasicObject BasicType) employeeID string
        --> w-tainted: write(employees_db.Employee.EmployeeID, freelancers_db.Freelancer.EmployeeID) {2}
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
[__2] (FieldObject FieldType) SpecFlag "freelancer" string
       --> w-tainted: write(employees_db.Employee.SpecFlag) {1}
[___3] (BasicObject BasicType) "freelancer" string

    --> w-tainted: write(freelancers_db.Freelancer) {1}
[0] (StructObject UserType) freelancer employee_app.Freelancer struct{EmployeeID string, FreelancerID string, Rate string, Terms string}
     --> w-tainted: write(freelancers_db.Freelancer) {1}
[_1] (Reference UserType) ref <freelancer employee_app.Freelancer struct{EmployeeID string, FreelancerID string, Rate string, Terms string}> @ EmployeeService
      --> w-tainted: write(freelancers_db.Freelancer) {1}
[__2] (Reference UserType) ref <employee employee_app.Freelancer struct{EmployeeID string, FreelancerID string, Rate string, Terms string}> @ FreelancerService
       --> w-tainted: write(freelancers_db.Freelancer.EmployeeID) {1}
[___3] (FieldObject FieldType) EmployeeID string
        --> w-tainted: write(freelancers_db.Freelancer.EmployeeID, employees_db.Employee.EmployeeID) {2}
[____4] (BasicObject BasicType) employeeID string
         --> w-tainted: write(employees_db.Employee.EmployeeID, freelancers_db.Freelancer.EmployeeID) {2}
[_____5] (Reference BasicType) ref <employeeID string> @ EmployeeService
          --> w-tainted: write(employees_db.Employee.EmployeeID, freelancers_db.Freelancer.EmployeeID) {2}
[______6] (Reference BasicType) ref <employeeID string> @ Frontend
       --> w-tainted: write(freelancers_db.Freelancer.FreelancerID) {1}
[___3] (FieldObject FieldType) FreelancerID string
        --> w-tainted: write(freelancers_db.Freelancer.FreelancerID) {1}
[____4] (BasicObject BasicType) freelancerID string
         --> w-tainted: write(freelancers_db.Freelancer.FreelancerID) {1}
[_____5] (Reference BasicType) ref <freelancerID string> @ EmployeeService
          --> w-tainted: write(freelancers_db.Freelancer.FreelancerID) {1}
[______6] (Reference BasicType) ref <freelancerID string> @ Frontend
       --> w-tainted: write(freelancers_db.Freelancer.Rate) {1}
[___3] (FieldObject FieldType) Rate string
        --> w-tainted: write(freelancers_db.Freelancer.Rate) {1}
[____4] (BasicObject BasicType) rate string
         --> w-tainted: write(freelancers_db.Freelancer.Rate) {1}
[_____5] (Reference BasicType) ref <rate string> @ EmployeeService
          --> w-tainted: write(freelancers_db.Freelancer.Rate) {1}
[______6] (Reference BasicType) ref <rate string> @ Frontend
       --> w-tainted: write(freelancers_db.Freelancer.Terms) {1}
[___3] (FieldObject FieldType) Terms string
        --> w-tainted: write(freelancers_db.Freelancer.Terms) {1}
[____4] (BasicObject BasicType) terms string
         --> w-tainted: write(freelancers_db.Freelancer.Terms) {1}
[_____5] (Reference BasicType) ref <terms string> @ EmployeeService
          --> w-tainted: write(freelancers_db.Freelancer.Terms) {1}
[______6] (Reference BasicType) ref <terms string> @ Frontend

[0] (InterfaceObject UserType) err .error
[_1] (Reference BasicType) ref <nil> @ EmployeeService

