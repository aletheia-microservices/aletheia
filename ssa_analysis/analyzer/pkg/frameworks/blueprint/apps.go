package blueprint

import (
	"log"

	specs_digota "github.com/blueprint-uservices/blueprint/examples/digota/wiring/specs"
	specs_dsb_media_sql "github.com/blueprint-uservices/blueprint/examples/dsb_media_sql/wiring/specs"
	specs_dsb_sn "github.com/blueprint-uservices/blueprint/examples/dsb_sn/wiring/specs"
	specs_postnotification_simple "github.com/blueprint-uservices/blueprint/examples/postnotification_simple/wiring/specs"
	specs_sockshop3 "github.com/blueprint-uservices/blueprint/examples/sockshop3/wiring/specs"
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
	"dsb_sn":                  {BLUEPRINT_PATH_EXAMPLES + "dsb_sn/workflow/socialnetwork", specs_dsb_sn.Docker},
}

func loadAppSpec(app string) cmdbuilder.SpecOption {
	if info, ok := APPS_INFO[app]; ok {
		return info.BlueprintSpec
	}
	log.Fatalf("unknown application %s", app)
	return cmdbuilder.SpecOption{}
}
