package detection

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type IgnoreCascadeEntry struct {
	Database string `yaml:"database"`
	Entity   string `yaml:"entity"`
}

type InputConfig struct {
	App           string               `yaml:"app"`
	IgnoreCascade []IgnoreCascadeEntry `yaml:"ignore_cascade"`
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

	for _, entry := range Config.IgnoreCascade {
		logrus.Infof("loaded ignore cascade entry: database = %s, entity = %s\n", entry.Database, entry.Entity)
	}

	if appname != Config.App {
		logrus.Fatalf("missmatch between appname argument (%s) and app in detection config (%s)", appname, Config.App)
	}
}
