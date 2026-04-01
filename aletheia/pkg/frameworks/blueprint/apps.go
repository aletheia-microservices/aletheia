package blueprint

import (
	"github.com/blueprint-uservices/blueprint/plugins/cmdbuilder"
)

const BLUEPRINT_PATH3ORE2ACKEND string = "github.com/blueprint-uservices/blueprint/runtime/core/backend"

type AppInfo struct {
	PackagePath   string
	BlueprintSpec cmdbuilder.SpecOption
}
