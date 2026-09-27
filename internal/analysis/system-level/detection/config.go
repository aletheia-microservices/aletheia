package detection

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type IgnoreCascadeEntry struct {
	Database        string `yaml:"database"`
	Entity          string `yaml:"entity"`
	TriggerDatabase string `yaml:"trigger_database"`
	TriggerEntity   string `yaml:"trigger_entity"`
}
type InputConfig struct {
	App               string               `yaml:"app"`
	IgnoreCascade     []IgnoreCascadeEntry `yaml:"ignore_cascade"`
	IgnoreForeignKeys []string             `yaml:"ignore_foreignkeys"`
	// linearizable databases; the others are eventually consistent
	LinearizableDatabases []string `yaml:"linearizable_databases"`
}

var Config InputConfig

func LoadInputConfig(appname string, path string) {
	logrus.Infof("loading detection config from %s\n", path)
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read detection config file: %v\n", err)
		os.Exit(1)
	}

	if err := yaml.Unmarshal(data, &Config); err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse detection config yaml: %v\n", err)
		os.Exit(1)
	}

	for _, entry := range Config.IgnoreForeignKeys {
		logrus.Infof("loaded ignore foreign key entry: %s\n", entry)
	}

	for _, entry := range Config.IgnoreCascade {
		logrus.Infof("loaded ignore cascade entry: database = %s, entity = %s, trigger_database = %s, trigger_entity = %s\n", entry.Database, entry.Entity, entry.TriggerDatabase, entry.TriggerEntity)
	}

	for _, database := range Config.LinearizableDatabases {
		logrus.Infof("loaded linearizable database: %s\n", database)
	}

	if appname != Config.App {
		logrus.Fatalf("missmatch between appname argument (%s) and app in detection config (%s)", appname, Config.App)
	}
}

// IsLinearizable reports whether the database is listed as linearizable; databases not listed are
// eventually consistent
func IsLinearizable(database string) bool {
	for _, entry := range Config.LinearizableDatabases {
		if entry == database {
			return true
		}
	}
	return false
}
