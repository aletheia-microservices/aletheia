//go:build !synthetic_small && !synthetic_medium && !synthetic_large

package blueprint

import (
	specs_digota "github.com/blueprint-uservices/blueprint/examples/digota/wiring/specs"
	specs_dotnet_eshop "github.com/blueprint-uservices/blueprint/examples/dotnet_eshop/wiring/specs"
	specs_dsb_hotel2 "github.com/blueprint-uservices/blueprint/examples/dsb_hotel2/wiring/specs"
	specs_dsb_mediamicroservices "github.com/blueprint-uservices/blueprint/examples/dsb_mediamicroservices/wiring/specs"
	specs_dsb_socialnetwork "github.com/blueprint-uservices/blueprint/examples/dsb_socialnetwork/wiring/specs"
	specs_eshopmicroservices "github.com/blueprint-uservices/blueprint/examples/eshopmicroservices/wiring/specs"
	specs_foobar "github.com/blueprint-uservices/blueprint/examples/foobar/wiring/specs"
	specs_foobar2 "github.com/blueprint-uservices/blueprint/examples/foobar2/wiring/specs"
	specs_postnotification "github.com/blueprint-uservices/blueprint/examples/postnotification/wiring/specs"
	specs_sockshop "github.com/blueprint-uservices/blueprint/examples/sockshop/wiring/specs"
	specs_syntheticapp "github.com/blueprint-uservices/blueprint/examples/synthetic_app/wiring/specs"
	specs_syntheticapp3 "github.com/blueprint-uservices/blueprint/examples/synthetic_app3/wiring/specs"
	specs_trainticket "github.com/blueprint-uservices/blueprint/examples/trainticket/wiring/specs"
	"github.com/blueprint-uservices/blueprint/plugins/cmdbuilder"
	"github.com/sirupsen/logrus"

	"analyzer/pkg/utils"
)

var APPS_INFO = map[string]AppInfo{
	"synthetic_app":          {utils.APP_PATH_SYNTHETIC_APP, specs_syntheticapp.Docker},
	"synthetic_app3":         {utils.APP_PATH_SYNTHETIC_APP3, specs_syntheticapp3.Docker},
	"postnotification":       {utils.APP_PATH_POSTNOTIFICATION, specs_postnotification.Docker},
	"digota":                 {utils.APP_PATH_DIGOTA, specs_digota.Docker},
	"eshopmicroservices":     {utils.APP_PATH_ESHOPMICROSERVICES, specs_eshopmicroservices.Docker},
	"dsb_mediamicroservices": {utils.APP_PATH_DSB_MEDIAMICROSERVICES, specs_dsb_mediamicroservices.Docker},
	"dotnet_eshop":           {utils.APP_PATH_DOTNET_ESHOP, specs_dotnet_eshop.Docker},
	"sockshop":               {utils.APP_PATH_SOCKSHOP, specs_sockshop.Docker},
	"dsb_socialnetwork":      {utils.APP_PATH_DSB_SOCIALNETWORK, specs_dsb_socialnetwork.Docker},
	"dsb_hotel2":             {utils.APP_PATH_DSB_HOTEL2, specs_dsb_hotel2.Original},
	"trainticket":            {utils.APP_PATH_TRAIN_TICKET, specs_trainticket.Docker},
	"foobar":                 {utils.APP_PATH_FOO_BAR, specs_foobar.Docker},
	"foobar2":                {utils.APP_PATH_FOO_BAR2, specs_foobar2.Docker},
}

func loadAppSpec(app string) cmdbuilder.SpecOption {
	logrus.WithField("app", app).Infof("loading app spec")
	if info, ok := APPS_INFO[app]; ok {
		return info.BlueprintSpec
	}
	logrus.Fatalf("unknown application %s", app)
	return cmdbuilder.SpecOption{}
}
