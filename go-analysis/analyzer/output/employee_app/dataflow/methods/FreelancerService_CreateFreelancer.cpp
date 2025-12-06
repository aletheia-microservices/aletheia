[0] (PointerObject PointerType) s (*employee_app.FreelancerServiceImpl struct{freelancersDB NoSQLDatabase})
[_1] (StructObject UserType) employee_app.FreelancerServiceImpl struct{freelancersDB NoSQLDatabase}
[__2] (FieldObject FieldType) freelancersDB NoSQLDatabase
[___3] (BlueprintBackendObject BlueprintBackendType) freelancersDB NoSQLDatabase

[0] (InterfaceObject UserType) ctx context.Context
[_1] (Reference UserType) ref <ctx context.Context> @ EmployeeService
[__2] (Reference UserType) ref <ctx context.Context> @ Frontend

    --> w-tainted: write(freelancers_db.Freelancer.EmployeeID, employees_db.Employee.EmployeeID) {2}
[0] (BasicObject BasicType) employeeID string
     --> w-tainted: write(employees_db.Employee.EmployeeID, freelancers_db.Freelancer.EmployeeID) {2}
[_1] (Reference BasicType) ref <employeeID string> @ EmployeeService
      --> w-tainted: write(employees_db.Employee.EmployeeID, freelancers_db.Freelancer.EmployeeID) {2}
[__2] (Reference BasicType) ref <employeeID string> @ Frontend

    --> w-tainted: write(freelancers_db.Freelancer.FreelancerID) {1}
[0] (BasicObject BasicType) freelancerID string
     --> w-tainted: write(freelancers_db.Freelancer.FreelancerID) {1}
[_1] (Reference BasicType) ref <freelancerID string> @ EmployeeService
      --> w-tainted: write(freelancers_db.Freelancer.FreelancerID) {1}
[__2] (Reference BasicType) ref <freelancerID string> @ Frontend

    --> w-tainted: write(freelancers_db.Freelancer.Rate) {1}
[0] (BasicObject BasicType) rate string
     --> w-tainted: write(freelancers_db.Freelancer.Rate) {1}
[_1] (Reference BasicType) ref <rate string> @ EmployeeService
      --> w-tainted: write(freelancers_db.Freelancer.Rate) {1}
[__2] (Reference BasicType) ref <rate string> @ Frontend

    --> w-tainted: write(freelancers_db.Freelancer.Terms) {1}
[0] (BasicObject BasicType) terms string
     --> w-tainted: write(freelancers_db.Freelancer.Terms) {1}
[_1] (Reference BasicType) ref <terms string> @ EmployeeService
      --> w-tainted: write(freelancers_db.Freelancer.Terms) {1}
[__2] (Reference BasicType) ref <terms string> @ Frontend

    --> w-tainted: write(freelancers_db.Freelancer) {1}
[0] (StructObject UserType) employee employee_app.Freelancer struct{EmployeeID string, FreelancerID string, Rate string, Terms string}
     --> w-tainted: write(freelancers_db.Freelancer.EmployeeID) {1}
[_1] (FieldObject FieldType) EmployeeID string
      --> w-tainted: write(freelancers_db.Freelancer.EmployeeID, employees_db.Employee.EmployeeID) {2}
[__2] (BasicObject BasicType) employeeID string
       --> w-tainted: write(employees_db.Employee.EmployeeID, freelancers_db.Freelancer.EmployeeID) {2}
[___3] (Reference BasicType) ref <employeeID string> @ EmployeeService
        --> w-tainted: write(employees_db.Employee.EmployeeID, freelancers_db.Freelancer.EmployeeID) {2}
[____4] (Reference BasicType) ref <employeeID string> @ Frontend
     --> w-tainted: write(freelancers_db.Freelancer.FreelancerID) {1}
[_1] (FieldObject FieldType) FreelancerID string
      --> w-tainted: write(freelancers_db.Freelancer.FreelancerID) {1}
[__2] (BasicObject BasicType) freelancerID string
       --> w-tainted: write(freelancers_db.Freelancer.FreelancerID) {1}
[___3] (Reference BasicType) ref <freelancerID string> @ EmployeeService
        --> w-tainted: write(freelancers_db.Freelancer.FreelancerID) {1}
[____4] (Reference BasicType) ref <freelancerID string> @ Frontend
     --> w-tainted: write(freelancers_db.Freelancer.Rate) {1}
[_1] (FieldObject FieldType) Rate string
      --> w-tainted: write(freelancers_db.Freelancer.Rate) {1}
[__2] (BasicObject BasicType) rate string
       --> w-tainted: write(freelancers_db.Freelancer.Rate) {1}
[___3] (Reference BasicType) ref <rate string> @ EmployeeService
        --> w-tainted: write(freelancers_db.Freelancer.Rate) {1}
[____4] (Reference BasicType) ref <rate string> @ Frontend
     --> w-tainted: write(freelancers_db.Freelancer.Terms) {1}
[_1] (FieldObject FieldType) Terms string
      --> w-tainted: write(freelancers_db.Freelancer.Terms) {1}
[__2] (BasicObject BasicType) terms string
       --> w-tainted: write(freelancers_db.Freelancer.Terms) {1}
[___3] (Reference BasicType) ref <terms string> @ EmployeeService
        --> w-tainted: write(freelancers_db.Freelancer.Terms) {1}
[____4] (Reference BasicType) ref <terms string> @ Frontend

[0] (BlueprintBackendObject BlueprintBackendType) collection NoSQLCollection {database = freelancers, collection = freelancers}

[0] (InterfaceObject UserType) err .error

[0] (InterfaceObject UserType) err .error

