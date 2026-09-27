package detection

import (
	"os"
	"path/filepath"
	"testing"
)

// TestDetectionOutput compares the warnings of every detector (and the inferred constraints)
// against the expected output in tests/expected/{app}
func TestDetectionOutput(t *testing.T) {
	for _, appname := range analyzedApps {
		t.Run(appname, func(t *testing.T) {
			runExpectedOutputTests(t, appname, "", filepath.Join("tests", "expected", appname))
		})
	}
}

// TestDetectionOutputWithConfig is the same as TestDetectionOutput but suppresses
// warnings using the detection config files in config/{app}.yaml
func TestDetectionOutputWithConfig(t *testing.T) {
	for _, appname := range analyzedApps {
		configPath := filepath.Join("config", appname+".yaml")
		if _, err := os.Stat(configPath); err != nil {
			continue
		}
		t.Run(appname, func(t *testing.T) {
			runExpectedOutputTests(t, appname, configPath, filepath.Join("tests", "expected", appname+"_config"))
		})
	}
}
