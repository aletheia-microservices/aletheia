# Microservices Consistency Analyzer (MSCA)

**Microservices Consistency Analyzer** is a toolkit that statically analyzes microservice applications written in frameworks (e.g. [Service Weaver](https://dl.acm.org/doi/10.1145/3593856.3595909) and [Blueprint](https://dl.acm.org/doi/10.1145/3600006.3613138)) that introduce a new modular approach for building and deploying microservices. The goal is to detect potential consistency violations at configuration time by leveraging the global view over the entire application, including the relationship between services and operations to underlying databases. At the moment, the toolkit supports Blueprint ([code](https://github.com/Blueprint-uServices/blueprint)), and is capable of capturing _cross-service inconsistencies_, introduced in [Antipode](https://dl.acm.org/doi/10.1145/3600006.3613176). In the future, we aim to support a wide variety of similar frameworks and capture additional consistency anomalies.

## Getting Started

After cloning the repository, fetch submodules for `blueprint` repository

```zsh
cd ms-consistency-analyzer
git submodule update --init --recursive
```

## Requirements

- [Golang](https://go.dev/doc/install) >= 1.22.2

## Getting Started

Install python requirements (use venv if needed):

```zsh
python3 -m venv ~/.venv
source ~/.venv/bin/activate
pip3 install -r requirements.txt
```

## Running the Tool

Available applications:

- `foobar`
- `shopping_app`
- `shopping_simple`
- `postnotification`
- `postnotification_simple`
- `trainticket`
- `sockshop2`
- `dsb_sn`
- `dsb_hotel`
- `app_constraints_referential_integrity`
- `employee_app`

Run the code analyzer:

```zsh
go run main.go --help

# usage
go run main.go -app=APP_NAME [--auto] [--xcy] [--fk] [--cascade] [--specialization]

# Cross-Service Inconsistency (XCY)
go run main.go --app postnotification --auto --xcy
go run main.go --app sockshop2 --auto --xcy

# Referential Integrity Constraint: Invalid Foreign Key
go run main.go --app postnotification --auto --fk_coordination
go run main.go --app postnotification_simple --auto --fk_coordination
go run main.go --app sockshop2 --auto --fk_coordination

# Referential Integrity Constraint: Absence of Cascading Deletes
go run main.go --app shopping_simple --auto --fk_cascade
go run main.go --app sockshop2 --auto --fk_cascade
go run main.go --app trainticket --auto --fk_cascade
go run main.go --app employee_app --auto --fk_cascade

# Specialization Constraint: Violation of Mandatory and Disjoint Constraints
go run main.go --app employee_app --auto --specialization

# Unicity Constraints
# > (STUDENTS_DB.Student.StudentID);(COUPONS_DB.Coupon.CouponID);(COUPONS_DB.ClaimedCoupon.CouponID,COUPONS_DB.ClaimedCoupon.UserID)
go run main.go --app coupons_app --auto --unicity
go run main.go --app coupons_app_sql --auto --unicity
go run main.go --app coupons_app_cache --auto --unicity

# Domain (Numerical) Constraints
go run main.go --app coupons_app --auto --numerical
go run main.go --app coupons_app_sql --auto --numerical
go run main.go --app coupons_app_cache --auto --numerical

# all applications and detection patterns
go run main.go -all=true --auto --xcy --fk_coordination --fk_cascade

# others
go run main.go --app digota --auto
go run main.go --app postnotification_simple --auto --compact_schema --fk_coordination --fk_concurrency --fk_cascade
go run main.go --app postnotification --auto --compact_schema --fk_coordination --fk_concurrency --fk_cascade
go run main.go --app sockshop2 --auto --fk_coordination
```

Run the code analyzer:

```zsh
## Primary Key Coordination
go run main.go --auto --pk_coordination --app postnotification_simple
go run main.go --auto --pk_coordination --app postnotification
go run main.go --auto --pk_coordination --app dsb_media
go run main.go --auto --pk_coordination --app dsb_media_sql
go run main.go --auto --pk_coordination --app digota

## Foreign Key Coordination
go run main.go --auto --fk_coordination --app digota
go run main.go --auto --fk_coordination --app shopping_simple
go run main.go --auto --fk_coordination --app shopping_app
go run main.go --auto --fk_coordination --app postnotification_simple
go run main.go --auto --fk_coordination --app postnotification
go run main.go --auto --fk_coordination --app dsb_media
go run main.go --auto --fk_coordination --app dsb_media_sql
go run main.go --auto --fk_coordination --app dsb_sn
go run main.go --auto --fk_coordination --app dsb_hotel
go run main.go --auto --fk_coordination --app trainticket

## Foreign Key Concurrency
go run main.go --auto --fk_concurrency --app digota
go run main.go --auto --fk_concurrency --app shopping_simple
go run main.go --auto --fk_concurrency --app shopping_app
go run main.go --auto --fk_concurrency --app postnotification_simple
go run main.go --auto --fk_concurrency --app postnotification
go run main.go --auto --fk_concurrency --app dsb_media
go run main.go --auto --fk_concurrency --app dsb_media_sql
go run main.go --auto --fk_concurrency --app dsb_sn
go run main.go --auto --fk_concurrency --app dsb_hotel
go run main.go --auto --fk_concurrency --app trainticket

## Foreign Key Cascade
go run main.go --auto --fk_cascade --app digota
go run main.go --auto --fk_cascade --app sockshop2
go run main.go --auto --fk_cascade --app shopping_simple
go run main.go --auto --fk_cascade --app shopping_app
go run main.go --auto --fk_cascade --app postnotification_simple
go run main.go --auto --fk_cascade --app postnotification
go run main.go --auto --fk_cascade --app dsb_media
go run main.go --auto --fk_cascade --app dsb_media_sql
go run main.go --auto --fk_cascade --app dsb_sn
go run main.go --auto --fk_cascade --app dsb_hotel
go run main.go --auto --fk_cascade --app trainticket

## Unicity Constraints
go run main.go --auto --unicity --app digota
go run main.go --auto --unicity --app postnotification_simple
go run main.go --auto --unicity --app postnotification
go run main.go --auto --unicity --app dsb_media
go run main.go --auto --unicity --app dsb_media_sql
go run main.go --auto --unicity --app dsb_sn
go run main.go --auto --unicity --app dsb_hotel
go run main.go --auto --unicity --app trainticket
```

Run the graph builder:

```zsh
source ~/.venv/bin/activate
./graphs.py --help

# usage
graphs.py [-h] [--app {postnotification,postnotification_simple,trainticket,shopping_app,shopping_simple,sockshop2,foobar}] [--graph {app,call}] [--labeled] [--all]

# example
./graphs.py --app foobar
./graphs.py --app postnotification
./graphs.py --app employee_app
./graphs.py --app coupons_app
./graphs.py --app digota

# all applications
./graphs.py --all
```

## Results

| Application Name                        | Runs    | XCY         | Cascade Delete | Foreign Key | Specialization | Unicity (Ind.) |
| --------------------------------------- | ------- | ----------- | -------------- | ----------- | -------------- | -------------- |
| `foobar`                                | **YES** | 0           | 0;0            | 0           | 0;0            | 0              |
| `shopping_simple`                       | **YES** | 0           | 1;1            | 0           | 0;0            | 0              |
| `shopping_app`                          | **YES** | 0 (PREV=1!) | 0;0            | 0 (PREV=1!) | 0;0            | 0              |
| `postnotification_simple`               | **YES** | 1           | 0;0            | 1           | 0;0            | 0              |
| `postnotification`                      | **YES** | 2           | 0;0            | 2           | 0;0            | 0              |
| `sockshop2`                             | **YES** | 0           | 4;0            | 0           | 0;0            | 0              |
| `trainticket`                           | **YES** | 0           | 13;0           | 0           | 0;0            | 0              |
| `dsb_sn`                                | **YES** | 0           | 0;0            | 0           | 0;0            | 0              |
| `dsb_hotel`                             | **NO**  | -           | -              | -           | -              | 0              |
| `app_constraints_referential_integrity` | **YES** | 0           | 2;1            | 1           | 0;0            | 0              |
| `employee_app`                          | **YES** | 0           | 2;3            | 0           | 1;0            | 0              |
| `coupons_app`                           | **YES** | 0           | 0;0            | 0           | 0;0            | 1;2            |
| `coupons_app_sql`                       | **YES** | 0           | 0;0            | -           | 0;0            | 1;2            |

## Structure

The toolkit takes as input applications in the `examples` folder from the [blueprint](https://github.com/mafaldacf/blueprint) git module forked from the original repository.

The representation of the application and the abstract graph are saved in the `assets` folder in JSON format.

The source code is located in the `pkg` folder, divided into the following modules: `abstractgraph`, `app`, `controlflow`, `frameworks`, `logger`, `service`, `types`, `utils`.

| Module                 | Description                                                                                                                                                                                                                                            |
| ---------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
|                        |                                                                                                                                                                                                                                                        |
| **Application**        | Initializes the application and registers **service nodes** and **database instances** with the application spec generated by Blueprint                                                                                                                |
| **Frameworks**         | Contains utilities specific to each framework. In the context of Blueprint, it contains (i) helper functions to compute the service spec from Intermediate Representation (IR), and (ii) methods that define backends like databases and queues        |
| **Service**            | Parses the **abstract syntax tree (ast)** of each service, capturing all service, database, or even internal calls within each exposed method                                                                                                          |
| **Control Flow Graph** | Parses the **control flow graph (cfg)** of each service, capturing all dependencies in declared and assigned variables within each exposed method block from service, database, or internal calls                                                      |
| **Abstract Graph**     | Generates the final abstract graph containing nodes (**abstract service calls** or **abstract database calls**) and edges representing the relationship and order between service and database calls, including queue handlers triggered by publishers |
