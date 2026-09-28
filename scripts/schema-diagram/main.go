// Draws the data schema of an analyzed app from output/{app}/schema.json: one box per database with
// its tables and their fields, and one edge per foreign key, from the referencing field to the
// referenced field.
//
// usage (from the aletheia directory): go run ./scripts/schema-diagram [-format png] [-dpi 300] [app ...]
//
//	app ...  apps to draw (default: every app with an output/{app}/schema.json)
//
// Writes output/{app}/diagrams/schema.dot and, when graphviz is installed, renders it to output/{app}/diagrams/schema.{format}.
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
	DATABASE_FILL    = "#f1f9f2"
	DATABASE_BORDER  = "#26803a"
	TABLE_HEADER     = "#d3f2d9"
	TABLE_BORDER     = "#26803a"
	KEY_COLOR        = "#6b6b6b"
	LEGEND_BORDER    = "#8a8a8a"
	FOREIGN_KEY_EDGE = "#555555"
	MANDATORY_EDGE   = "#d62728"
	TRANSITIVE_EDGE  = "#2a6fd6"
)

type Schema struct {
	Name        string   `json:"name"`
	Fields      []string `json:"fields"`
	Constraints []string `json:"constraints"`
}

type Database struct {
	Name    string   `json:"name"`
	Schemas []Schema `json:"schemas"`
}

// Table is a schema of a database (e.g., SQL table, NoSQL collection, queue topic, cache key)
type Table struct {
	ID     string // "{db}.{name}", the prefix of the path of every field in the table
	DB     string
	Name   string
	Fields []string                   // paths relative to the table (e.g., "Items[*].ID"), sorted
	Ports  map[string]string          // relative path => port of its row in the table label
	Keys   map[string]map[string]bool // relative path => key markers (PK, UQ, FK)
}

// Endpoint is a table or one of its fields
type Endpoint struct {
	Table *Table
	Field string // relative path, or "" for the table itself
}

type ForeignKey struct {
	From       Endpoint
	To         Endpoint
	Mandatory  bool
	Transitive bool
}

type SchemaGraph struct {
	Tables         []*Table // sorted by ID
	ForeignKeys    []ForeignKey
	EmptyDatabases []string // databases without tables, which are not drawn
}

var (
	// e.g. "FOREIGN_KEY assurance_db.assurance.OrderID REFERENCES order_db.order.ID [MANDATORY] [T]"
	foreignKeyRegex = regexp.MustCompile(`^FOREIGN_KEY (\S+) REFERENCES (\S+)(.*)$`)
	// e.g. "PRIMARY KEY (movie_id_db.movie._id)" or "UNIQUE (movie_id_db.movie.Title)"
	keyRegex = regexp.MustCompile(`^(PRIMARY KEY|UNIQUE) \((.*)\)$`)
)

var keyMarkers = map[string]string{
	"PRIMARY KEY": "PK",
	"UNIQUE":      "UQ",
}

// keyMarkersOrder is the order in which key markers are shown next to a field
var keyMarkersOrder = []string{"PK", "UQ", "FK"}

// resolve splits a path into the table it belongs to and the field relative to that table
// e.g. "order_db.order.Items[*].ID" => (order_db.order, "Items[*].ID")
// e.g. "delivery_queue.notification" => (delivery_queue.notification, "")
func (graph *SchemaGraph) resolve(path string) (Endpoint, bool) {
	var best *Table
	for _, table := range graph.Tables {
		if path != table.ID && !strings.HasPrefix(path, table.ID+".") {
			continue
		}
		// the longest match wins, in case a table name contains dots (e.g., tables "a" and "a.b")
		if best == nil || len(table.ID) > len(best.ID) {
			best = table
		}
	}
	if best == nil {
		return Endpoint{}, false
	}
	return Endpoint{Table: best, Field: strings.TrimPrefix(strings.TrimPrefix(path, best.ID), ".")}, true
}

// path returns the full path of the endpoint (e.g., "order_db.order.ID")
func (endpoint Endpoint) path() string {
	if endpoint.Field == "" {
		return endpoint.Table.ID
	}
	return endpoint.Table.ID + "." + endpoint.Field
}

func (table *Table) addField(field string) {
	if field != "" && !slices.Contains(table.Fields, field) {
		table.Fields = append(table.Fields, field)
	}
}

func (table *Table) addKey(field string, marker string) {
	if field == "" {
		return
	}
	table.addField(field)
	if table.Keys[field] == nil {
		table.Keys[field] = make(map[string]bool)
	}
	table.Keys[field][marker] = true
}

func buildSchemaGraph(appName string, databases map[string]Database) *SchemaGraph {
	graph := &SchemaGraph{}
	for dbName, database := range databases {
		if len(database.Schemas) == 0 {
			graph.EmptyDatabases = append(graph.EmptyDatabases, dbName)
			continue
		}
		for _, schema := range database.Schemas {
			table := &Table{
				ID:    dbName + "." + schema.Name,
				DB:    dbName,
				Name:  schema.Name,
				Ports: make(map[string]string),
				Keys:  make(map[string]map[string]bool),
			}
			for _, path := range schema.Fields {
				// the list of fields also includes the table itself (e.g., "order_db.order")
				if strings.HasPrefix(path, table.ID+".") {
					table.addField(strings.TrimPrefix(path, table.ID+"."))
				}
			}
			graph.Tables = append(graph.Tables, table)
		}
	}
	sort.Strings(graph.EmptyDatabases)
	sort.Slice(graph.Tables, func(i, j int) bool {
		return graph.Tables[i].ID < graph.Tables[j].ID
	})

	for _, database := range databases {
		for _, schema := range database.Schemas {
			for _, constraint := range schema.Constraints {
				graph.addConstraint(appName, constraint)
			}
		}
	}
	sort.Slice(graph.ForeignKeys, func(i, j int) bool {
		a, b := graph.ForeignKeys[i], graph.ForeignKeys[j]
		if a.From.path() != b.From.path() {
			return a.From.path() < b.From.path()
		}
		return a.To.path() < b.To.path()
	})

	for _, table := range graph.Tables {
		sort.Strings(table.Fields)
		for i, field := range table.Fields {
			table.Ports[field] = fmt.Sprintf("f%d", i)
		}
	}
	return graph
}

func (graph *SchemaGraph) addConstraint(appName string, constraint string) {
	if m := foreignKeyRegex.FindStringSubmatch(constraint); m != nil {
		from, okFrom := graph.resolve(m[1])
		to, okTo := graph.resolve(m[2])
		if !okFrom || !okTo {
			fmt.Fprintf(os.Stderr, "[WARN] %s: skipping constraint with unknown table: %s\n", appName, constraint)
			return
		}
		// constraints can reference fields missing from the list of fields of the schema, so draw them anyway
		from.Table.addKey(from.Field, "FK")
		to.Table.addField(to.Field)
		graph.ForeignKeys = append(graph.ForeignKeys, ForeignKey{
			From:       from,
			To:         to,
			Mandatory:  strings.Contains(m[3], "[MANDATORY]"),
			Transitive: strings.Contains(m[3], "[T]"),
		})
		return
	}
	if m := keyRegex.FindStringSubmatch(constraint); m != nil {
		for _, path := range strings.Split(m[2], ", ") {
			endpoint, ok := graph.resolve(path)
			if !ok {
				fmt.Fprintf(os.Stderr, "[WARN] %s: skipping constraint with unknown table: %s\n", appName, constraint)
				return
			}
			endpoint.Table.addKey(endpoint.Field, keyMarkers[m[1]])
		}
		return
	}
	fmt.Fprintf(os.Stderr, "[WARN] %s: skipping unknown constraint: %s\n", appName, constraint)
}

func tableLabel(table *Table) string {
	var b strings.Builder
	fmt.Fprintf(&b, `<TABLE BORDER="1" CELLBORDER="0" CELLSPACING="0" CELLPADDING="4" COLOR=%q BGCOLOR="white">`, TABLE_BORDER)
	fmt.Fprintf(&b, `<TR><TD PORT="t" BGCOLOR=%q><B>%s</B></TD></TR>`, TABLE_HEADER, html.EscapeString(table.Name))
	if len(table.Fields) == 0 {
		b.WriteString(`<TR><TD><I><FONT COLOR="#8a8a8a">no fields</FONT></I></TD></TR>`)
	}
	for _, field := range table.Fields {
		var markers []string
		for _, marker := range keyMarkersOrder {
			if table.Keys[field][marker] {
				markers = append(markers, marker)
			}
		}
		keys := ""
		if len(markers) > 0 {
			keys = " " + keysLabel(strings.Join(markers, " "))
		}
		fmt.Fprintf(&b, `<TR><TD PORT=%q ALIGN="LEFT">%s%s</TD></TR>`, table.Ports[field], html.EscapeString(field), keys)
	}
	b.WriteString(`</TABLE>`)
	return b.String()
}

// keysLabel formats key markers (e.g., "PK FK") as shown next to a field
func keysLabel(markers string) string {
	return fmt.Sprintf(`<FONT POINT-SIZE="9" COLOR=%q><B>%s</B></FONT>`, KEY_COLOR, markers)
}

type LegendItem struct {
	Sample string // e.g., a colored line (see lineSample) or a key marker
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

func (endpoint Endpoint) node() string {
	if endpoint.Field == "" {
		return fmt.Sprintf("%q:t", endpoint.Table.ID)
	}
	return fmt.Sprintf("%q:%s", endpoint.Table.ID, endpoint.Table.Ports[endpoint.Field])
}

// backEdges returns the indexes of the foreign keys that close a cycle between tables, found with a depth-first
// search in the order of the tables and foreign keys (self-references are not back edges)
func (graph *SchemaGraph) backEdges() map[int]bool {
	outgoing := make(map[*Table][]int) // table => indexes of its foreign keys to other tables
	for i, fk := range graph.ForeignKeys {
		if fk.From.Table != fk.To.Table {
			outgoing[fk.From.Table] = append(outgoing[fk.From.Table], i)
		}
	}
	const (
		unvisited = iota
		visiting
		visited
	)
	state := make(map[*Table]int)
	back := make(map[int]bool)
	var visit func(table *Table)
	visit = func(table *Table) {
		state[table] = visiting
		for _, i := range outgoing[table] {
			switch to := graph.ForeignKeys[i].To.Table; state[to] {
			case unvisited:
				visit(to)
			case visiting:
				back[i] = true
			}
		}
		state[table] = visited
	}
	for _, table := range graph.Tables {
		if state[table] == unvisited {
			visit(table)
		}
	}
	return back
}

// lastRankTable returns a table at the end of a longest chain of foreign keys (with the back edges
// reversed), which dot draws in the last column, or nil if there are no foreign keys between tables
func (graph *SchemaGraph) lastRankTable(back map[int]bool) *Table {
	next := make(map[*Table][]*Table)
	indegree := make(map[*Table]int)
	for i, fk := range graph.ForeignKeys {
		from, to := fk.From.Table, fk.To.Table
		if from == to {
			continue
		}
		if back[i] {
			from, to = to, from
		}
		next[from] = append(next[from], to)
		indegree[to]++
	}
	// longest path from any table with no incoming foreign keys, in topological order
	rank := make(map[*Table]int)
	var queue []*Table
	for _, table := range graph.Tables {
		if indegree[table] == 0 {
			queue = append(queue, table)
		}
	}
	for len(queue) > 0 {
		table := queue[0]
		queue = queue[1:]
		for _, to := range next[table] {
			rank[to] = max(rank[to], rank[table]+1)
			indegree[to]--
			if indegree[to] == 0 {
				queue = append(queue, to)
			}
		}
	}
	var last *Table
	for _, table := range graph.Tables {
		if (len(next[table]) > 0 || rank[table] > 0) && (last == nil || rank[table] > rank[last]) {
			last = table
		}
	}
	return last
}

func buildDot(appName string, graph *SchemaGraph) string {
	legend := legendLabel([]LegendItem{
		{lineSample(MANDATORY_EDGE), "mandatory foreign key"},
		{lineSample(FOREIGN_KEY_EDGE), "foreign key, not mandatory"},
		{lineSample(TRANSITIVE_EDGE), "transitive foreign key"},
		{keysLabel("PK"), "primary key"},
		{keysLabel("UQ"), "unique"},
		{keysLabel("FK"), "foreign key field"},
	})

	var b strings.Builder
	fmt.Fprintf(&b, "digraph %q {\n", appName)
	// straight segments (polyline) are easier to follow through dense areas than curves, and mclimit=10
	// makes dot work harder at reducing edge crossings
	fmt.Fprintf(&b, "  graph [rankdir=LR, splines=polyline, mclimit=10, nodesep=0.6, ranksep=1.4, fontname=\"Helvetica\", fontsize=20, labelloc=t, label=<<B>%s</B><BR/><FONT POINT-SIZE=\"16\">data schema</FONT>>];\n", html.EscapeString(appName))
	fmt.Fprintf(&b, "  node [shape=plain, fontname=\"Helvetica\", fontsize=11];\n")
	fmt.Fprintf(&b, "  edge [arrowsize=0.9, penwidth=1.3];\n")

	// the legend is the label of a borderless cluster around the whole diagram, which places it in the bottom-right corner
	fmt.Fprintf(&b, "\n  subgraph cluster_legend {\n")
	fmt.Fprintf(&b, "    pencolor=transparent; labelloc=b; labeljust=r; fontsize=11; label=<%s>;\n", legend)

	// one cluster per database with its tables (tables are sorted by ID, so tables of the same database are adjacent)
	// nested clusters inherit the attributes of cluster_legend, so override them all
	for i, table := range graph.Tables {
		if i == 0 || graph.Tables[i-1].DB != table.DB {
			fmt.Fprintf(&b, "\n    subgraph %q {\n", "cluster_db_"+table.DB)
			fmt.Fprintf(&b, "      label=%q; labelloc=t; labeljust=c; fontsize=13; style=\"rounded,filled\"; fillcolor=%q; pencolor=%q;\n", table.DB, DATABASE_FILL, DATABASE_BORDER)
		}
		fmt.Fprintf(&b, "      %q [label=<%s>];\n", table.ID, tableLabel(table))
		if i == len(graph.Tables)-1 || graph.Tables[i+1].DB != table.DB {
			fmt.Fprintf(&b, "    }\n")
		}
	}
	fmt.Fprintf(&b, "  }\n")

	back := graph.backEdges()

	// tables without foreign keys go in their own column after the last one, through invisible edges,
	// so that they don't take space between the tables with foreign keys
	if last := graph.lastRankTable(back); last != nil {
		connected := make(map[*Table]bool)
		for _, fk := range graph.ForeignKeys {
			connected[fk.From.Table], connected[fk.To.Table] = true, true
		}
		fmt.Fprintf(&b, "\n  // tables without foreign keys\n")
		for _, table := range graph.Tables {
			if !connected[table] {
				fmt.Fprintf(&b, "  %q -> %q [style=invis, weight=0];\n", last.ID, table.ID)
			}
		}
	}

	// every edge leaves the right side (e) of the referencing field and enters the left side (w) of the referenced field
	fmt.Fprintf(&b, "\n  // foreign keys\n")
	for i, fk := range graph.ForeignKeys {
		color := FOREIGN_KEY_EDGE
		if fk.Mandatory {
			color = MANDATORY_EDGE
		}
		// transitive takes precedence, so mandatory transitive foreign keys are also blue
		if fk.Transitive {
			color = TRANSITIVE_EDGE
		}
		switch {
		case fk.From.Table == fk.To.Table:
			// self-references loop on the right side of the table
			fmt.Fprintf(&b, "  %s:e -> %s:e [color=%q];\n", fk.From.node(), fk.To.node(), color)
		case back[i]:
			// foreign keys that close a cycle are written reversed with dir=back, so that dot never draws an
			// edge right to left, while the arrow still points at the referenced field
			fmt.Fprintf(&b, "  %s:e -> %s:w [color=%q, dir=back];\n", fk.To.node(), fk.From.node(), color)
		default:
			fmt.Fprintf(&b, "  %s:e -> %s:w [color=%q];\n", fk.From.node(), fk.To.node(), color)
		}
	}

	fmt.Fprintf(&b, "}\n")
	return b.String()
}

func drawSchema(appName string, format string, dpi int, dotPath string) {
	appDir := filepath.Join(DEFAULT_OUTPUT_DIR, appName)
	schemaFile := filepath.Join(appDir, "schema.json")

	data, err := os.ReadFile(schemaFile)
	if err != nil {
		log.Fatalf("[ERROR] error reading %s: %v", schemaFile, err)
	}
	var databases map[string]Database
	if err := json.Unmarshal(data, &databases); err != nil {
		log.Fatalf("[ERROR] error parsing %s: %v", schemaFile, err)
	}

	graph := buildSchemaGraph(appName, databases)
	if len(graph.EmptyDatabases) > 0 {
		fmt.Printf("[INFO] %s: not drawing databases without tables: %s\n", appName, strings.Join(graph.EmptyDatabases, ", "))
	}

	diagramsDir := filepath.Join(appDir, DIAGRAMS_DIR)
	if err := os.MkdirAll(diagramsDir, 0755); err != nil {
		log.Fatalf("[ERROR] error creating %s: %v", diagramsDir, err)
	}

	dotFile := filepath.Join(diagramsDir, "schema.dot")
	if err := os.WriteFile(dotFile, []byte(buildDot(appName, graph)), 0644); err != nil {
		log.Fatalf("[ERROR] error writing %s: %v", dotFile, err)
	}
	fmt.Printf("[INFO] generated %s\n", dotFile)

	if dotPath == "" {
		return
	}
	imageFile := filepath.Join(diagramsDir, "schema."+format)
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
		// no apps specified — draw every app with a schema.json in DEFAULT_OUTPUT_DIR
		matches, err := filepath.Glob(filepath.Join(DEFAULT_OUTPUT_DIR, "*", "schema.json"))
		if err != nil {
			log.Fatalf("[ERROR] error listing %s: %v", DEFAULT_OUTPUT_DIR, err)
		}
		if len(matches) == 0 {
			log.Fatalf("[ERROR] no schema.json files found in %s/*/", DEFAULT_OUTPUT_DIR)
		}
		for _, match := range matches {
			apps = append(apps, filepath.Base(filepath.Dir(match)))
		}
	}

	for _, appName := range apps {
		drawSchema(appName, *format, *dpi, dotPath)
	}
}
