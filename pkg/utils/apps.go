package utils

import (
	specs_app_constraints_referential_integrity "github.com/blueprint-uservices/blueprint/examples/app_constraints_referential_integrity/wiring/specs"
	specs_coupons_app "github.com/blueprint-uservices/blueprint/examples/coupons_app/wiring/specs"
	specs_coupons_app_sql "github.com/blueprint-uservices/blueprint/examples/coupons_app_sql/wiring/specs"
	dsb_hotel "github.com/blueprint-uservices/blueprint/examples/dsb_hotel/wiring/specs"
	dsb_sn "github.com/blueprint-uservices/blueprint/examples/dsb_sn/wiring/specs"
	specs_employee_app "github.com/blueprint-uservices/blueprint/examples/employee_app/wiring/specs"
	specs_foobar "github.com/blueprint-uservices/blueprint/examples/foobar/wiring/specs"
	specs_postnotification "github.com/blueprint-uservices/blueprint/examples/postnotification/wiring/specs"
	specs_postnotification_simple "github.com/blueprint-uservices/blueprint/examples/postnotification_simple/wiring/specs"
	specs_shopping_app "github.com/blueprint-uservices/blueprint/examples/shopping_app/wiring/specs"
	specs_shopping_simple "github.com/blueprint-uservices/blueprint/examples/shopping_simple/wiring/specs"
	specs_sockshop2 "github.com/blueprint-uservices/blueprint/examples/sockshop2/wiring/specs"
	specs_trainticket "github.com/blueprint-uservices/blueprint/examples/train_ticket/wiring/specs"
	"github.com/blueprint-uservices/blueprint/plugins/cmdbuilder"

	"analyzer/pkg/logger"
)

const (
	PATH_BLUEPRINT_EXAMPLES     string = "github.com/blueprint-uservices/blueprint/examples/"
	PATH_BLUEPRINT_CORE_BACKEND string = "github.com/blueprint-uservices/blueprint/runtime/core/backend"
)

var Apps = []string{"foobar", "shopping_simple", "shopping_app", "postnotification_simple", "postnotification", "sockshop2", "trainticket", "app_constraints_referential_integrity", "employee_app", "dsb_sn", "dsb_hotel", "coupons_app", "coupons_app_sql"}

type AppInfo struct {
	PackagePath   string
	BlueprintSpec cmdbuilder.SpecOption
}

var APPS_INFO = map[string]AppInfo{
	"postnotification":                      {PATH_BLUEPRINT_EXAMPLES + "postnotification/workflow/postnotification", specs_postnotification.Docker},
	"postnotification_simple":               {PATH_BLUEPRINT_EXAMPLES + "postnotification_simple/workflow/postnotification_simple", specs_postnotification_simple.Docker},
	"app_constraints_referential_integrity": {PATH_BLUEPRINT_EXAMPLES + "app_constraints_referential_integrity/workflow/app_constraints_referential_integrity", specs_app_constraints_referential_integrity.Docker},
	"employee_app":                          {PATH_BLUEPRINT_EXAMPLES + "employee_app/workflow/employee_app", specs_employee_app.Docker},
	"coupons_app":                           {PATH_BLUEPRINT_EXAMPLES + "coupons_app/workflow/coupons_app", specs_coupons_app.Docker},
	"coupons_app_sql":                       {PATH_BLUEPRINT_EXAMPLES + "coupons_app_sql/workflow/coupons_app_sql", specs_coupons_app_sql.Docker},
	"foobar":                                {PATH_BLUEPRINT_EXAMPLES + "foobar/workflow/foobar", specs_foobar.Docker},
	"sockshop2":                             {PATH_BLUEPRINT_EXAMPLES + "sockshop2/workflow", specs_sockshop2.Docker},
	"trainticket":                           {PATH_BLUEPRINT_EXAMPLES + "train_ticket/workflow", specs_trainticket.Docker},
	"shopping_app":                          {PATH_BLUEPRINT_EXAMPLES + "shopping_app/workflow", specs_shopping_app.Docker},
	"shopping_simple":                       {PATH_BLUEPRINT_EXAMPLES + "shopping_simple/workflow", specs_shopping_simple.Docker},
	"dsb_hotel":                             {PATH_BLUEPRINT_EXAMPLES + "dsb_hotel/workflow/hotelreservation", dsb_hotel.Original},
	"dsb_sn":                                {PATH_BLUEPRINT_EXAMPLES + "dsb_sn/workflow/socialnetwork", dsb_sn.Docker},
}

func LoadAppPath(app string) string {
	if info, ok := APPS_INFO[app]; ok {
		return info.PackagePath
	}
	logger.Logger.Fatalf("unknown application name %s", app)
	return ""
}

func LoadAppSpec(app string) cmdbuilder.SpecOption {
	if info, ok := APPS_INFO[app]; ok {
		return info.BlueprintSpec
	}
	logger.Logger.Fatalf("unknown application %s", app)
	return cmdbuilder.SpecOption{}
}
