// Draws the service dependency graph of an analyzed app from output/{app}/app.json:
// services are blue nodes, databases are green nodes, and edges are service -> service
// and service -> database dependencies.
//
// usage (from the aletheia directory): go run ./scripts/app-diagram [-format png] [-dpi 300] [app ...]
//
//	app ...  apps to draw (default: every app with an output/{app}/app.json)
//
// Writes output/{app}/diagrams/app.dot and, when graphviz is installed, renders it to output/{app}/diagrams/app.{format}.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strings"
)

const (
	DEFAULT_OUTPUT_DIR = "output"
	DIAGRAMS_DIR       = "diagrams" // inside output/{app}/
)

// BITMAP_FORMATS are the graphviz output formats whose resolution is set with -dpi
// (vector formats such as svg and pdf stay sharp at any size)
var BITMAP_FORMATS = []string{"png", "jpg", "jpeg", "gif", "bmp", "tif", "tiff", "webp"}

const (
	SERVICE_FILL    = "#cfe2ff"
	SERVICE_BORDER  = "#1a5fb4"
	DATABASE_FILL   = "#d3f2d9"
	DATABASE_BORDER = "#26803a"
	LEGEND_BORDER   = "#8a8a8a"
	SERVICE_EDGE    = "#555555"
	DATABASE_EDGE   = "#26803a"
)

type Service struct {
	Name      string   `json:"name"`
	Fields    []string `json:"fields"`
	Services  []string `json:"services"`
	Databases []string `json:"databases"`
}

type App struct {
	Name      string    `json:"name"`
	Databases []string  `json:"databases"`
	Services  []Service `json:"services"`
}

// fieldInstanceRegex extracts the instance name of a field
// e.g. "db #0 (assurance_db)" => "assurance_db"
var fieldInstanceRegex = regexp.MustCompile(`\(([^)]+)\)\s*$`)

// serviceDatabases returns the databases used by a service, sorted by name.
// The "databases" list of a service in app.json is not populated yet, so we also
// match the instance name of each field against the databases of the app
func serviceDatabases(service Service, appDatabases map[string]bool) []string {
	set := make(map[string]bool)
	for _, db := range service.Databases {
		set[db] = true
	}
	for _, field := range service.Fields {
		if m := fieldInstanceRegex.FindStringSubmatch(field); m != nil && appDatabases[m[1]] {
			set[m[1]] = true
		}
	}
	dbs := make([]string, 0, len(set))
	for db := range set {
		dbs = append(dbs, db)
	}
	sort.Strings(dbs)
	return dbs
}

func serviceNode(name string) string  { return fmt.Sprintf("%q", "svc:"+name) }
func databaseNode(name string) string { return fmt.Sprintf("%q", "db:"+name) }

type LegendItem struct {
	Sample string // e.g., a colored line (see lineSample) or a node (see nodeSample)
	Text   string
}

// legendLabel returns a boxed legend with one row per item
func legendLabel(items []LegendItem) string {
	var b strings.Builder
	fmt.Fprintf(&b, `<TABLE BORDER="1" CELLBORDER="0" CELLSPACING="0" CELLPADDING="3" COLOR=%q BGCOLOR="white">`, LEGEND_BORDER)
	b.WriteString(`<TR><TD COLSPAN="2" ALIGN="LEFT"><B>Legend</B></TD></TR>`)
	for _, item := range items {
		fmt.Fprintf(&b, `<TR><TD WIDTH="30">%s</TD><TD ALIGN="LEFT">%s</TD></TR>`, item.Sample, item.Text)
	}
	b.WriteString(`</TABLE>`)
	return b.String()
}

// lineSample draws a short line in the given color, like an edge of the diagram
func lineSample(color string) string {
	return fmt.Sprintf(`<TABLE BORDER="0" CELLPADDING="0" CELLSPACING="0"><TR><TD BGCOLOR=%q WIDTH="24" HEIGHT="3" FIXEDSIZE="TRUE"></TD></TR></TABLE>`, color)
}

// nodeSample draws a small box with the given fill and border colors, like a node of the diagram
func nodeSample(fill string, border string, rounded bool) string {
	style := ""
	if rounded {
		style = ` STYLE="rounded"`
	}
	return fmt.Sprintf(`<TABLE BORDER="1" CELLBORDER="0"%s COLOR=%q BGCOLOR=%q CELLPADDING="0" CELLSPACING="0" FIXEDSIZE="TRUE" WIDTH="24" HEIGHT="13"><TR><TD></TD></TR></TABLE>`, style, border, fill)
}

func buildDot(app App) string {
	appDatabases := make(map[string]bool, len(app.Databases))
	for _, db := range app.Databases {
		appDatabases[db] = true
	}

	var b strings.Builder
	fmt.Fprintf(&b, "digraph %q {\n", app.Name)
	legend := legendLabel([]LegendItem{
		{nodeSample(SERVICE_FILL, SERVICE_BORDER, true), "service"},
		{nodeSample(DATABASE_FILL, DATABASE_BORDER, false), "database"},
		{lineSample(SERVICE_EDGE), "service &#8594; service dependency"},
		{lineSample(DATABASE_EDGE), "service &#8594; database dependency"},
	})
	fmt.Fprintf(&b, "  graph [rankdir=LR, splines=true, nodesep=0.3, ranksep=1.2, fontname=\"Helvetica\", fontsize=20, labelloc=t, label=<<B>%s</B><BR/><FONT POINT-SIZE=\"16\">service dependency graph</FONT>>];\n", html.EscapeString(app.Name))
	fmt.Fprintf(&b, "  node [fontname=\"Helvetica\", fontsize=12, style=filled, penwidth=1.5];\n")
	fmt.Fprintf(&b, "  edge [arrowsize=0.7, penwidth=1.1];\n\n")

	// declare every node up front so dependencies missing from the top-level lists are still drawn with the right style
	serviceNames := make(map[string]bool)
	databaseNames := make(map[string]bool)
	dbsByService := make(map[string][]string)
	for _, service := range app.Services {
		serviceNames[service.Name] = true
		for _, dep := range service.Services {
			serviceNames[dep] = true
		}
		dbs := serviceDatabases(service, appDatabases)
		dbsByService[service.Name] = dbs
		for _, db := range dbs {
			databaseNames[db] = true
		}
	}
	for _, db := range app.Databases {
		databaseNames[db] = true
	}

	// the legend is the label of a borderless cluster around all nodes, which places it in the bottom-right corner
	fmt.Fprintf(&b, "  subgraph cluster_legend {\n")
	fmt.Fprintf(&b, "  pencolor=transparent; labelloc=b; labeljust=r; fontsize=11; label=<%s>;\n\n", legend)

	fmt.Fprintf(&b, "  // services\n")
	fmt.Fprintf(&b, "  node [shape=box, style=\"filled,rounded\", fillcolor=%q, color=%q];\n", SERVICE_FILL, SERVICE_BORDER)
	for _, name := range sortedKeys(serviceNames) {
		fmt.Fprintf(&b, "  %s [label=%q];\n", serviceNode(name), name)
	}

	fmt.Fprintf(&b, "\n  // databases\n")
	fmt.Fprintf(&b, "  node [shape=cylinder, style=filled, fillcolor=%q, color=%q];\n", DATABASE_FILL, DATABASE_BORDER)
	for _, name := range sortedKeys(databaseNames) {
		fmt.Fprintf(&b, "  %s [label=%q];\n", databaseNode(name), name)
	}

	fmt.Fprintf(&b, "  }\n")

	fmt.Fprintf(&b, "\n  // service -> service\n")
	fmt.Fprintf(&b, "  edge [color=%q];\n", SERVICE_EDGE)
	for _, service := range app.Services {
		for _, dep := range service.Services {
			fmt.Fprintf(&b, "  %s -> %s;\n", serviceNode(service.Name), serviceNode(dep))
		}
	}

	fmt.Fprintf(&b, "\n  // service -> database\n")
	fmt.Fprintf(&b, "  edge [color=%q];\n", DATABASE_EDGE)
	for _, service := range app.Services {
		for _, db := range dbsByService[service.Name] {
			fmt.Fprintf(&b, "  %s -> %s;\n", serviceNode(service.Name), databaseNode(db))
		}
	}

	fmt.Fprintf(&b, "}\n")
	return b.String()
}

func sortedKeys(set map[string]bool) []string {
	keys := make([]string, 0, len(set))
	for k := range set {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func drawApp(appName string, format string, dpi int, dotPath string) {
	appDir := filepath.Join(DEFAULT_OUTPUT_DIR, appName)
	appFile := filepath.Join(appDir, "app.json")

	data, err := os.ReadFile(appFile)
	if err != nil {
		log.Fatalf("[ERROR] error reading %s: %v", appFile, err)
	}
	var app App
	if err := json.Unmarshal(data, &app); err != nil {
		log.Fatalf("[ERROR] error parsing %s: %v", appFile, err)
	}

	diagramsDir := filepath.Join(appDir, DIAGRAMS_DIR)
	if err := os.MkdirAll(diagramsDir, 0755); err != nil {
		log.Fatalf("[ERROR] error creating %s: %v", diagramsDir, err)
	}

	dotFile := filepath.Join(diagramsDir, "app.dot")
	if err := os.WriteFile(dotFile, []byte(buildDot(app)), 0644); err != nil {
		log.Fatalf("[ERROR] error writing %s: %v", dotFile, err)
	}
	fmt.Printf("[INFO] generated %s\n", dotFile)

	if dotPath == "" {
		return
	}
	imageFile := filepath.Join(diagramsDir, "app."+format)
	args := []string{"-T" + format, dotFile, "-o", imageFile}
	if slices.Contains(BITMAP_FORMATS, format) {
		args = append(args, fmt.Sprintf("-Gdpi=%d", dpi))
	}
	out, err := exec.Command(dotPath, args...).CombinedOutput()
	if err != nil {
		log.Fatalf("[ERROR] error rendering %s: %v\n%s", imageFile, err, out)
	}
	fmt.Printf("[INFO] generated %s\n", imageFile)
}

func ensureRunFromAletheia() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("[ERROR] %v", err)
	}
	if filepath.Base(cwd) != "aletheia" {
		fmt.Fprintf(os.Stderr, "[ERROR] this script must be run from the aletheia directory (current: %s)\n", cwd)
		os.Exit(1)
	}
}

func main() {
	ensureRunFromAletheia()

	format := flag.String("format", "png", "graphviz output format of the rendered diagram (e.g. png, svg, pdf)")
	dpi := flag.Int("dpi", 300, "resolution of bitmap formats (e.g. png) in dots per inch; graphviz defaults to 96")
	flag.Parse()

	dotPath, err := exec.LookPath("dot")
	if err != nil {
		fmt.Fprintf(os.Stderr, "[WARN] graphviz 'dot' not found in PATH, writing .dot files only (install with: brew install graphviz)\n")
		dotPath = ""
	}

	apps := flag.Args()
	if len(apps) == 0 {
		// no apps specified — draw every app with an app.json in DEFAULT_OUTPUT_DIR
		matches, err := filepath.Glob(filepath.Join(DEFAULT_OUTPUT_DIR, "*", "app.json"))
		if err != nil {
			log.Fatalf("[ERROR] error listing %s: %v", DEFAULT_OUTPUT_DIR, err)
		}
		if len(matches) == 0 {
			log.Fatalf("[ERROR] no app.json files found in %s/*/", DEFAULT_OUTPUT_DIR)
		}
		for _, match := range matches {
			apps = append(apps, filepath.Base(filepath.Dir(match)))
		}
	}

	for _, appName := range apps {
		drawApp(appName, *format, *dpi, dotPath)
	}
}
