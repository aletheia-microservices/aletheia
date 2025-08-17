package blueprint

import (
	"log"

	specs_digota "github.com/blueprint-uservices/blueprint/examples/digota/wiring/specs"
	specs_dsb_media_sql "github.com/blueprint-uservices/blueprint/examples/dsb_media_sql/wiring/specs"
	specs_postnotification_simple "github.com/blueprint-uservices/blueprint/examples/postnotification_simple/wiring/specs"
	specs_sockshop3 "github.com/blueprint-uservices/blueprint/examples/sockshop3/wiring/specs"
	/* specs_app_constraints_referential_integrity "github.com/blueprint-uservices/blueprint/examples/app_constraints_referential_integrity/wiring/specs"
	specs_coupons_app "github.com/blueprint-uservices/blueprint/examples/coupons_app/wiring/specs"
	specs_coupons_app_cache "github.com/blueprint-uservices/blueprint/examples/coupons_app_cache/wiring/specs"
	specs_coupons_app_sql "github.com/blueprint-uservices/blueprint/examples/coupons_app_sql/wiring/specs"
	specs_dsb_hotel "github.com/blueprint-uservices/blueprint/examples/dsb_hotel/wiring/specs"
	specs_dsb_media "github.com/blueprint-uservices/blueprint/examples/dsb_media/wiring/specs"
	specs_dsb_sn "github.com/blueprint-uservices/blueprint/examples/dsb_sn/wiring/specs"
	specs_employee_app "github.com/blueprint-uservices/blueprint/examples/employee_app/wiring/specs"
	specs_foobar "github.com/blueprint-uservices/blueprint/examples/foobar/wiring/specs"
	specs_postnotification "github.com/blueprint-uservices/blueprint/examples/postnotification/wiring/specs"
	specs_shopping_app "github.com/blueprint-uservices/blueprint/examples/shopping_app/wiring/specs"
	specs_shopping_simple "github.com/blueprint-uservices/blueprint/examples/shopping_simple/wiring/specs"
	specs_sockshop2 "github.com/blueprint-uservices/blueprint/examples/sockshop2/wiring/specs"
	specs_trainticket "github.com/blueprint-uservices/blueprint/examples/train_ticket/wiring/specs" */
	"github.com/blueprint-uservices/blueprint/plugins/cmdbuilder"
)

const (
	BLUEPRINT_PATH_EXAMPLES     string = "github.com/blueprint-uservices/blueprint/examples/"
	BLUEPRINT_PATH_CORE_BACKEND string = "github.com/blueprint-uservices/blueprint/runtime/core/backend"
)

type AppInfo struct {
	PackagePath   string
	BlueprintSpec cmdbuilder.SpecOption
}

var APPS_INFO = map[string]AppInfo{
	"postnotification_simple": {BLUEPRINT_PATH_EXAMPLES + "postnotification_simple/workflow/postnotification_simple", specs_postnotification_simple.Docker},
	"digota":                  {BLUEPRINT_PATH_EXAMPLES + "digota/workflow/digota", specs_digota.Docker},
	"dsb_media_sql":           {BLUEPRINT_PATH_EXAMPLES + "dsb_media_sql/workflow/mediamicroservices_sql", specs_dsb_media_sql.Docker},
	"sockshop3":               {BLUEPRINT_PATH_EXAMPLES + "sockshop3/workflow/sockshop3", specs_sockshop3.Docker},
	/* "postnotification":                      {BLUEPRINT_PATH_EXAMPLES + "postnotification/workflow/postnotification", specs_postnotification.Docker},
	"app_constraints_referential_integrity": {BLUEPRINT_PATH_EXAMPLES + "app_constraints_referential_integrity/workflow/app_constraints_referential_integrity", specs_app_constraints_referential_integrity.Docker},
	"employee_app":                          {BLUEPRINT_PATH_EXAMPLES + "employee_app/workflow/employee_app", specs_employee_app.Docker},
	"coupons_app":                           {BLUEPRINT_PATH_EXAMPLES + "coupons_app/workflow/coupons_app", specs_coupons_app.Docker},
	"coupons_app_sql":                       {BLUEPRINT_PATH_EXAMPLES + "coupons_app_sql/workflow/coupons_app_sql", specs_coupons_app_sql.Docker},
	"coupons_app_cache":                     {BLUEPRINT_PATH_EXAMPLES + "coupons_app_cache/workflow/coupons_app_cache", specs_coupons_app_cache.Docker},
	"foobar":                                {BLUEPRINT_PATH_EXAMPLES + "foobar/workflow/foobar", specs_foobar.Docker},
	"sockshop2":                             {BLUEPRINT_PATH_EXAMPLES + "sockshop2/workflow", specs_sockshop2.Docker},
	"trainticket":                           {BLUEPRINT_PATH_EXAMPLES + "train_ticket/workflow", specs_trainticket.Docker},
	"shopping_app":                          {BLUEPRINT_PATH_EXAMPLES + "shopping_app/workflow", specs_shopping_app.Docker},
	"shopping_simple":                       {BLUEPRINT_PATH_EXAMPLES + "shopping_simple/workflow", specs_shopping_simple.Docker},
	"dsb_hotel":                             {BLUEPRINT_PATH_EXAMPLES + "dsb_hotel/workflow/hotelreservation", specs_dsb_hotel.Original},
	"dsb_sn":                                {BLUEPRINT_PATH_EXAMPLES + "dsb_sn/workflow/socialnetwork", specs_dsb_sn.Docker},
	"dsb_media":                             {BLUEPRINT_PATH_EXAMPLES + "dsb_media/workflow/mediamicroservices", specs_dsb_media.Docker}, */
}

func loadAppSpec(app string) cmdbuilder.SpecOption {
	if info, ok := APPS_INFO[app]; ok {
		return info.BlueprintSpec
	}
	log.Fatalf("unknown application %s", app)
	return cmdbuilder.SpecOption{}
}
