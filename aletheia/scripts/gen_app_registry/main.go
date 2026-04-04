package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"gopkg.in/yaml.v2"
)

//go:embed apps.go.template
var rawTemplate string

var fileTemplate = template.Must(template.New("apps").Parse(rawTemplate))

const (
	DEFAULT_CONFIG_FILE = "apps.yaml"
	DEFAULT_CONFIG_DIR  = "configs"
	DEFAULT_OUTPUT_DIR  = "pkg/frameworks/blueprint/apps"
	DEFAULT_BUILD_TAG   = "!eval"
	DEFAULT_GO_MOD      = "go.mod"
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
	BuildTag string      `yaml:"build_tag"`
	Apps     []AppConfig `yaml:"apps"`
}

type specEntry struct {
	Alias      string
	ImportPath string
	CaseName   string
	Method     string
}

type templateData struct {
	BuildTag   string
	ConfigFile string
	Specs      []specEntry
}

// importAlias returns "specs_<name>" with hyphens replaced by underscores.
func importAlias(name string) string {
	return "specs_" + strings.ReplaceAll(name, "-", "_")
}

// deriveSpecMethod derives the Blueprint spec method name from spec_name by
// stripping the "<app_name>_" prefix and title-casing the remainder
// e.g. name="digota",    specName="digota_docker"     => "Docker"
// e.g. name="dsb_hotel2",specName="dsb_hotel2_original" => "Original"
func deriveSpecMethod(appName, specName string) string {
	suffix := strings.TrimPrefix(specName, appName+"_")
	if suffix == specName || suffix == "" {
		log.Fatalf("cannot derive method from spec_name %q and app name %q", specName, appName)
	}
	return strings.ToUpper(suffix[:1]) + suffix[1:]
}

// extractModulesFromAppRoot derives module paths and local replacement dir from app_root:
// wiringModule   = app_root + "/wiring"
// workflowModule = app_root + "/workflow"
// relativeBase      = "blueprint/examples/<last segment of app_root>"
func extractModulesFromAppRoot(appRoot string) (wiringModule, workflowModule, relativeBase string) {
	wiringModule = appRoot + "/wiring"
	workflowModule = appRoot + "/workflow"
	relativeBase = "../blueprint/examples/" + path.Base(appRoot)
	return
}

func updateGoMod(gomodPath string, apps []AppConfig) {
	data, err := os.ReadFile(gomodPath)
	if err != nil {
		log.Fatalf("[ERROR] error reading %s: %v", gomodPath, err)
	}
	content := string(data)

	// avoid duplicates
	has := func(s string) bool { return strings.Contains(content, s) }

	var additions strings.Builder
	addLine := func(line string) {
		additions.WriteString(line)
		additions.WriteString("\n")
		content += line + "\n" // keep content in sync so duplicate checks work within the same run
	}

	for _, app := range apps {
		if app.AppRoot == "" {
			continue
		}
		wiring, workflow, relativeBase := extractModulesFromAppRoot(app.AppRoot)
		appRoot := app.AppRoot

		if !has(wiring) {
			addLine(fmt.Sprintf("require %s v0.0.0", wiring))
		}
		if !has(workflow) {
			addLine(fmt.Sprintf("require %s v0.0.0 // indirect", workflow))
		}
		if !has("replace " + appRoot + " ") {
			addLine(fmt.Sprintf("replace %s => %s", appRoot, relativeBase))
		}
		if !has("replace " + workflow) {
			addLine(fmt.Sprintf("replace %s => %s/workflow", workflow, relativeBase))
		}
		if !has("replace " + wiring) {
			addLine(fmt.Sprintf("replace %s => %s/wiring", wiring, relativeBase))
		}
	}

	if additions.Len() == 0 {
		fmt.Printf("[INFO] go.mod already up to date\n")
		return
	}

	f, err := os.OpenFile(gomodPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("[ERROR] error opening %s: %v", gomodPath, err)
	}
	defer f.Close()

	if _, err := f.WriteString("\n" + additions.String()); err != nil {
		log.Fatalf("[ERROR] error writing to %s: %v", gomodPath, err)
	}
	fmt.Printf("[INFO] updated %s\n", gomodPath)
}


func deriveOutputFile(configFile string) string {
	base := path.Base(configFile)
	stem := strings.TrimSuffix(base, path.Ext(base))
	return path.Join(DEFAULT_OUTPUT_DIR, stem+".go")
}

func processConfigFile(configPath string) {
	outputFile := deriveOutputFile(configPath)

	fmt.Printf("[INFO] reading config file %s...\n", configPath)

	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("[ERROR] error reading config: %v", err)
	}
	var cfg AppsConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("[ERROR] error parsing config: %v", err)
	}

	if cfg.BuildTag == "" {
		cfg.BuildTag = DEFAULT_BUILD_TAG
	}

	// generate apps go file
	var specs []specEntry
	for _, app := range cfg.Apps {
		if app.SpecPath == "" || app.SpecName == "" {
			continue
		}
		specs = append(specs, specEntry{
			Alias:      importAlias(app.Name),
			ImportPath: app.SpecPath,
			CaseName:   app.SpecName,
			Method:     deriveSpecMethod(app.Name, app.SpecName),
		})
	}
	if len(specs) == 0 {
		log.Fatalf("[ERROR] no apps with spec_path and spec_name found in %s", configPath)
	}

	f, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("[ERROR] error creating output file: %v", err)
	}
	defer f.Close()

	if err := fileTemplate.Execute(f, templateData{
		BuildTag:   cfg.BuildTag,
		ConfigFile: configPath,
		Specs:      specs,
	}); err != nil {
		log.Fatalf("[ERROR] error rendering template: %v", err)
	}

	fmt.Printf("[INFO] generated %s\n", outputFile)

	// update go mod file
	updateGoMod(DEFAULT_GO_MOD, cfg.Apps)
}

func main() {
	configFile := flag.String("config", "", "path to apps YAML config (relative to configs/); if omitted, all YAML files in configs/ are processed")
	flag.Parse()

	if *configFile != "" {
		processConfigFile(path.Join(DEFAULT_CONFIG_DIR, *configFile))
		return
	}

	// no config specified — process all YAML files in DEFAULT_CONFIG_DIR
	matches, err := filepath.Glob(filepath.Join(DEFAULT_CONFIG_DIR, "*.yaml"))
	if err != nil {
		log.Fatalf("[ERROR] error listing config dir: %v", err)
	}
	if len(matches) == 0 {
		log.Fatalf("[ERROR] no YAML files found in %s", DEFAULT_CONFIG_DIR)
	}
	for _, configPath := range matches {
		processConfigFile(configPath)
	}
}
