package blueprint_apps

import (
	"os"

	"github.com/blueprint-uservices/blueprint/plugins/cmdbuilder"
	"gopkg.in/yaml.v2"

	"analyzer/pkg/utils"
)

type AppConfig struct {
	Name        string   `yaml:"name"`
	AppRoot     string   `yaml:"app_root"`
	PackagePath string   `yaml:"package_path"`
	SpecName    string   `yaml:"spec_name"`
	SpecPath    string   `yaml:"spec_path"`
	SQLTables   []string `yaml:"sql_tables"`
	NoSQLPath   string   `yaml:"nosql_path"`
}

type AppsConfig struct {
	Apps []AppConfig `yaml:"apps"`
}

type AppInfo struct {
	PackagePath   string
	BlueprintSpec cmdbuilder.SpecOption
}

var APPS_INFO = map[string]AppInfo{}
const BLUEPRINT_EXAMPLES_PKG_PREFIX    = "github.com/blueprint-uservices/blueprint/examples/"

func loadAppsConfig(filepath string) (*AppsConfig, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	var cfg AppsConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func initAppsFromConfig(cfg *AppsConfig) {
	for _, app := range cfg.Apps {
		if app.PackagePath != "" {
			var spec cmdbuilder.SpecOption
			if app.SpecName != "" {
				spec = resolveSpec(app.SpecName)
			}
			APPS_INFO[app.Name] = AppInfo{
				PackagePath:   app.PackagePath,
				BlueprintSpec: spec,
			}
			utils.RegisterApp(BLUEPRINT_EXAMPLES_PKG_PREFIX + app.PackagePath)
		}
		if app.NoSQLPath != "" {
			utils.APPS_NOSQL_SCHEMAS[app.Name] = app.NoSQLPath
		}
		if len(app.SQLTables) > 0 {
			utils.APPS_SQL_TABLES[app.Name] = app.SQLTables
		}
	}
}
