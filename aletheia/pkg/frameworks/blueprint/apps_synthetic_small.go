//go:build synthetic_small && !synthetic_medium && !synthetic_large

package blueprint

import (
	specs_synthetic_app1 "github.com/blueprint-uservices/blueprint/examples/synthetic_app1/wiring/specs"
	specs_synthetic_app2 "github.com/blueprint-uservices/blueprint/examples/synthetic_app2/wiring/specs"
	specs_synthetic_app3 "github.com/blueprint-uservices/blueprint/examples/synthetic_app3/wiring/specs"
	"github.com/blueprint-uservices/blueprint/plugins/cmdbuilder"
	"github.com/sirupsen/logrus"

	"analyzer/pkg/utils"
)

var APPS_INFO = map[string]AppInfo{
	"synthetic_app1":         {utils.APP_PATH_SYNTHETIC_APP1, specs_synthetic_app1.Docker},
	"synthetic_app2":         {utils.APP_PATH_SYNTHETIC_APP2, specs_synthetic_app2.Docker},
	"synthetic_app3":         {utils.APP_PATH_SYNTHETIC_APP3, specs_synthetic_app3.Docker},
}

func loadAppSpec(app string) cmdbuilder.SpecOption {
	logrus.WithField("app", app).Infof("loading synthetic app spec")
	if info, ok := APPS_INFO[app]; ok {
		return info.BlueprintSpec
	}
	logrus.Fatalf("unknown application %s", app)
	return cmdbuilder.SpecOption{}
}
