package blueprint

import (
	"log"

	specs_digota "github.com/blueprint-uservices/blueprint/examples/digota/wiring/specs"
	specs_dsb_hotel2 "github.com/blueprint-uservices/blueprint/examples/dsb_hotel2/wiring/specs"
	specs_dsb_media_nosql "github.com/blueprint-uservices/blueprint/examples/dsb_media_nosql/wiring/specs"
	specs_dsb_media_sql "github.com/blueprint-uservices/blueprint/examples/dsb_media_sql/wiring/specs"
	specs_dsb_sn2 "github.com/blueprint-uservices/blueprint/examples/dsb_sn2/wiring/specs"
	specs_foobar "github.com/blueprint-uservices/blueprint/examples/foobar/wiring/specs"
	specs_largescaleapp "github.com/blueprint-uservices/blueprint/examples/large_scale_app/wiring/specs"
	specs_largescaleapp_A "github.com/blueprint-uservices/blueprint/examples/large_scale_app_A/wiring/specs"
	specs_largescaleapp_B "github.com/blueprint-uservices/blueprint/examples/large_scale_app_B/wiring/specs"
	specs_largescaleapp_C "github.com/blueprint-uservices/blueprint/examples/large_scale_app_C/wiring/specs"
	specs_largescaleapp_D "github.com/blueprint-uservices/blueprint/examples/large_scale_app_D/wiring/specs"
	specs_largescaleapp_E "github.com/blueprint-uservices/blueprint/examples/large_scale_app_E/wiring/specs"
	specs_postnotification_simple "github.com/blueprint-uservices/blueprint/examples/postnotification_simple/wiring/specs"
	specs_sockshop3 "github.com/blueprint-uservices/blueprint/examples/sockshop3/wiring/specs"
	specs_trainticket "github.com/blueprint-uservices/blueprint/examples/train_ticket2/wiring/specs"
	"github.com/blueprint-uservices/blueprint/plugins/cmdbuilder"

	"analyzer/pkg/utils"
)

const BLUEPRINT_PATH_CORE_BACKEND string = "github.com/blueprint-uservices/blueprint/runtime/core/backend"

type AppInfo struct {
	PackagePath   string
	BlueprintSpec cmdbuilder.SpecOption
}

var APPS_INFO = map[string]AppInfo{
	"postnotification_simple": {utils.APP_PATH_POSTNOTIFICATION_SIMPLE, specs_postnotification_simple.Docker},
	"digota":                  {utils.APP_PATH_DIGOTA, specs_digota.Docker},
	"dsb_media_sql":           {utils.APP_PATH_DSB_MEDIA_SQL, specs_dsb_media_sql.Docker},
	"dsb_media_nosql":         {utils.APP_PATH_DSB_MEDIA_NOSQL, specs_dsb_media_nosql.Docker},
	"sockshop3":               {utils.APP_PATH_SOCKSHOP3, specs_sockshop3.Docker},
	"dsb_sn2":                 {utils.APP_PATH_DSB_SN2, specs_dsb_sn2.Docker},
	"dsb_hotel2":              {utils.APP_PATH_DSB_HOTEL2, specs_dsb_hotel2.Original},
	"train_ticket2":           {utils.APP_PATH_TRAIN_TICKET2, specs_trainticket.Docker},
	"foobar":                  {utils.APP_PATH_FOO_BAR, specs_foobar.Docker},
	"large_scale_app":         {utils.APP_PATH_LARGE_SCALE_APP, specs_largescaleapp.Docker},
	"large_scale_app_A":       {utils.APP_PATH_LARGE_SCALE_APP_A, specs_largescaleapp_A.Docker},
	"large_scale_app_B":       {utils.APP_PATH_LARGE_SCALE_APP_B, specs_largescaleapp_B.Docker},
	"large_scale_app_C":       {utils.APP_PATH_LARGE_SCALE_APP_C, specs_largescaleapp_C.Docker},
	"large_scale_app_D":       {utils.APP_PATH_LARGE_SCALE_APP_D, specs_largescaleapp_D.Docker},
	"large_scale_app_E":       {utils.APP_PATH_LARGE_SCALE_APP_E, specs_largescaleapp_E.Docker},
}

func loadAppSpec(app string) cmdbuilder.SpecOption {
	if info, ok := APPS_INFO[app]; ok {
		return info.BlueprintSpec
	}
	log.Fatalf("unknown application %s", app)
	return cmdbuilder.SpecOption{}
}
