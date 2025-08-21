package blueprint

import (
	"log"

	specs_digota "github.com/blueprint-uservices/blueprint/examples/digota/wiring/specs"
	specs_dsb_hotel2 "github.com/blueprint-uservices/blueprint/examples/dsb_hotel2/wiring/specs"
	specs_dsb_media_sql "github.com/blueprint-uservices/blueprint/examples/dsb_media_sql/wiring/specs"
	specs_dsb_sn "github.com/blueprint-uservices/blueprint/examples/dsb_sn/wiring/specs"
	specs_foobar "github.com/blueprint-uservices/blueprint/examples/foobar/wiring/specs"
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
	"sockshop3":               {utils.APP_PATH_SOCKSHOP3, specs_sockshop3.Docker},
	"dsb_sn":                  {utils.APP_PATH_DSB_SN, specs_dsb_sn.Docker},
	"dsb_hotel2":              {utils.APP_PATH_DSB_HOTEL2, specs_dsb_hotel2.Original},
	"train_ticket2":           {utils.APP_PATH_TRAIN_TICKET2, specs_trainticket.Docker},
	"foobar":                  {utils.APP_PATH_FOO_BAR, specs_foobar.Docker},
}

func loadAppSpec(app string) cmdbuilder.SpecOption {
	if info, ok := APPS_INFO[app]; ok {
		return info.BlueprintSpec
	}
	log.Fatalf("unknown application %s", app)
	return cmdbuilder.SpecOption{}
}
