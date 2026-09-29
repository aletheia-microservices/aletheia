package parser_test

import (
	"strings"
	"testing"

	abstractgraphinput "github.com/aletheia-microservices/aletheia/internal/analysis/system-level/abstractgraph/input"
	"github.com/aletheia-microservices/aletheia/internal/app"
	appparser "github.com/aletheia-microservices/aletheia/internal/app/parser"
)

func TestInitFromInput(t *testing.T) {
	model := &abstractgraphinput.InputModel{
		Databases: []*abstractgraphinput.Database{
			{Name: "users_db", Type: "RelationalDB", Schemas: []*abstractgraphinput.Schema{
				{Name: "user", PrimaryKey: []string{"TenantID", "UserID"}, Unique: []string{"Email"}},
			}},
			{Name: "events_queue", Type: "Queue"},
		},
		Functions: []*abstractgraphinput.Function{
			{Service: "Frontend", Method: "Register", Calls: []*abstractgraphinput.Call{
				{CallID: "rpc", ServiceCall: &abstractgraphinput.ServiceCall{Service: "UserService", Method: "Create"}},
				{CallID: "rpc2", ServiceCall: &abstractgraphinput.ServiceCall{Service: "UserService", Method: "Create"}},
			}},
			{Service: "UserService", Method: "Create"},
			{Service: "UserService", Method: "Run"},
		},
	}
	a := app.NewApp("input")
	appparser.InitFromInput(a, model)

	if !a.HasDatabase("users_db") || !a.GetDatabaseByName("events_queue").IsQueue() {
		t.Fatal("missing databases")
	}
	schema := a.GetDatabaseByName("users_db").GetSchemaByNameIfExists("user")
	if schema == nil || len(schema.GetAllConstraints()) != 2 {
		t.Fatal("missing schema constraints")
	}
	// fields of a composite primary key are not primary keys by themselves
	if pk := schema.GetFieldByPath("users_db.user.TenantID").GetConstraintPrimaryKey(); pk != nil || len(schema.GetAllConstraints()[0].GetFields()) != 2 {
		t.Fatal("incorrect composite primary key")
	}
	if !schema.GetFieldByPath("users_db.user.Email").IsUnique() {
		t.Fatal("missing unique key")
	}

	if a.NumberOfMicroservices() != 2 {
		t.Fatalf("unexpected services: %v", a.GetAllServices())
	}
	users := a.GetServiceByName("UserService")
	if !users.HasMethod("Create") || !users.HasInitializerMethod() {
		t.Fatal("incorrect service methods")
	}
	if data, _ := a.GetServiceByName("Frontend").MarshalJSON(); !strings.Contains(string(data), `"services":["UserService"]`) {
		t.Fatalf("incorrect service dependencies: %s", data)
	}
}
