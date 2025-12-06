package blueprint

import (
	specs_digota "github.com/blueprint-uservices/blueprint/examples/digota/wiring/specs"
	specs_dsb_hotel2 "github.com/blueprint-uservices/blueprint/examples/dsb_hotel2/wiring/specs"
	specs_dsb_mediamicroservices "github.com/blueprint-uservices/blueprint/examples/dsb_mediamicroservices/wiring/specs"
	specs_dsb_socialnetwork "github.com/blueprint-uservices/blueprint/examples/dsb_socialnetwork/wiring/specs"
	specs_eshopmicroservices "github.com/blueprint-uservices/blueprint/examples/eshopmicroservices/wiring/specs"
	specs_foobar "github.com/blueprint-uservices/blueprint/examples/foobar/wiring/specs"
	specs_foobar2 "github.com/blueprint-uservices/blueprint/examples/foobar2/wiring/specs"
	specs_postnotification "github.com/blueprint-uservices/blueprint/examples/postnotification/wiring/specs"
	specs_sockshop "github.com/blueprint-uservices/blueprint/examples/sockshop/wiring/specs"
	specs_synthetic_app "github.com/blueprint-uservices/blueprint/examples/synthetic_app/wiring/specs"
	specs_synthetic_app1 "github.com/blueprint-uservices/blueprint/examples/synthetic_app1/wiring/specs"
	specs_synthetic_app2 "github.com/blueprint-uservices/blueprint/examples/synthetic_app2/wiring/specs"
	specs_synthetic_app3 "github.com/blueprint-uservices/blueprint/examples/synthetic_app3/wiring/specs"
	specs_synthetic_app4 "github.com/blueprint-uservices/blueprint/examples/synthetic_app4/wiring/specs"
	specs_synthetic_app5 "github.com/blueprint-uservices/blueprint/examples/synthetic_app5/wiring/specs"
	specs_synthetic_app6 "github.com/blueprint-uservices/blueprint/examples/synthetic_app6/wiring/specs"
	specs_synthetic_app7 "github.com/blueprint-uservices/blueprint/examples/synthetic_app7/wiring/specs"
	specs_synthetic_app8 "github.com/blueprint-uservices/blueprint/examples/synthetic_app8/wiring/specs"
	specs_synthetic_appA "github.com/blueprint-uservices/blueprint/examples/synthetic_appA/wiring/specs"
	specs_synthetic_appB "github.com/blueprint-uservices/blueprint/examples/synthetic_appB/wiring/specs"
	specs_trainticket "github.com/blueprint-uservices/blueprint/examples/trainticket/wiring/specs"
	"github.com/blueprint-uservices/blueprint/plugins/cmdbuilder"
	"github.com/sirupsen/logrus"

	"analyzer/pkg/utils"
)

const BLUEPRINT_PATH3ORE2ACKEND string = "github.com/blueprint-uservices/blueprint/runtime/core/backend"

type AppInfo struct {
	PackagePath   string
	BlueprintSpec cmdbuilder.SpecOption
}

var APPS_INFO = map[string]AppInfo{
	"postnotification":       {utils.APP_PATH_POSTNOTIFICATION, specs_postnotification.Docker},
	"digota":                 {utils.APP_PATH_DIGOTA, specs_digota.Docker},
	"eshopmicroservices":     {utils.APP_PATH_ESHOPMICROSERVICES, specs_eshopmicroservices.Docker},
	"dsb_mediamicroservices": {utils.APP_PATH_DSB_MEDIAMICROSERVICES, specs_dsb_mediamicroservices.Docker},
	"sockshop":               {utils.APP_PATH_SOCKSHOP, specs_sockshop.Docker},
	"dsb_socialnetwork":      {utils.APP_PATH_DSB_SOCIALNETWORK, specs_dsb_socialnetwork.Docker},
	"dsb_hotel2":             {utils.APP_PATH_DSB_HOTEL2, specs_dsb_hotel2.Original},
	"trainticket":            {utils.APP_PATH_TRAIN_TICKET, specs_trainticket.Docker},
	"foobar":                 {utils.APP_PATH_FOO_BAR, specs_foobar.Docker},
	"foobar2":                {utils.APP_PATH_FOO_BAR2, specs_foobar2.Docker},
	"synthetic_app":          {utils.APP_PATH_SYNTHETIC_APP, specs_synthetic_app.Docker},
	"synthetic_appA":         {utils.APP_PATH_SYNTHETIC_APPA, specs_synthetic_appA.Docker},
	"synthetic_appB":         {utils.APP_PATH_SYNTHETIC_APPB, specs_synthetic_appB.Docker},
	"synthetic_app1":         {utils.APP_PATH_SYNTHETIC_APP1, specs_synthetic_app1.Docker},
	"synthetic_app2":         {utils.APP_PATH_SYNTHETIC_APP2, specs_synthetic_app2.Docker},
	"synthetic_app3":         {utils.APP_PATH_SYNTHETIC_APP3, specs_synthetic_app3.Docker},
	"synthetic_app4":         {utils.APP_PATH_SYNTHETIC_APP4, specs_synthetic_app4.Docker},
	"synthetic_app5":         {utils.APP_PATH_SYNTHETIC_APP5, specs_synthetic_app5.Docker},
	"synthetic_app6":         {utils.APP_PATH_SYNTHETIC_APP6, specs_synthetic_app6.Docker},
	"synthetic_app7":         {utils.APP_PATH_SYNTHETIC_APP7, specs_synthetic_app7.Docker},
	"synthetic_app8":         {utils.APP_PATH_SYNTHETIC_APP8, specs_synthetic_app8.Docker},
}

func loadAppSpec(app string) cmdbuilder.SpecOption {
	if info, ok := APPS_INFO[app]; ok {
		return info.BlueprintSpec
	}
	logrus.Fatalf("unknown application %s", app)
	return cmdbuilder.SpecOption{}
}
