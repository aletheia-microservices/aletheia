//go:build !synthetic_small && synthetic_medium && !synthetic_large

package blueprint

import (
	specs_synthetic_app4 "github.com/blueprint-uservices/blueprint/examples/synthetic_app4/wiring/specs"
	"github.com/blueprint-uservices/blueprint/plugins/cmdbuilder"
	"github.com/sirupsen/logrus"

	"analyzer/pkg/utils"
)

var APPS_INFO = map[string]AppInfo{
	"synthetic_app4":         {utils.APP_PATH_SYNTHETIC_APP4, specs_synthetic_app4.Docker},
}

func loadAppSpec(app string) cmdbuilder.SpecOption {
	logrus.WithField("app", app).Infof("loading synthetic app spec")
	if info, ok := APPS_INFO[app]; ok {
		return info.BlueprintSpec
	}
	logrus.Fatalf("unknown application %s", app)
	return cmdbuilder.SpecOption{}
}
