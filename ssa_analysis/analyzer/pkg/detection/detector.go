package detection

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"analyzer/pkg/abstractgraph"
	"analyzer/pkg/app"
)

type Detector interface {
	GetResults() string
	GetSummary() string
	SetSummary(summary string)
	ComputeResults()
	GetTypeString() string

	OnNewRun(app *app.App)
	OnEndRun(app *app.App)
	OnNewRequest(node *abstractgraph.AbstractNode)
	OnEndRequest(app *app.App)
	OnNewNode(app *app.App, node *abstractgraph.AbstractNode)
	OnEndNode(app *app.App, node *abstractgraph.AbstractNode)

	// database calls
	OnRead(app *app.App, edge *abstractgraph.AbstractEdge)
	OnWrite(app *app.App, edge *abstractgraph.AbstractEdge)
	OnUpdate(app *app.App, edge *abstractgraph.AbstractEdge)
	OnDelete(app *app.App, edge *abstractgraph.AbstractEdge)
}

const TEXT_BOLD_LIGHT_YELLOW = "\033[1;38;5;179m"
const TEXT_BOLD_LIGHT_RED = "\033[1;31m"
const TEXT_RESET_COLOR = "\033[0m"
const TEXT_BOLD_LIGHT_BLUE = "\033[1;38;5;68m"
const TEXT_BOLD_LIGHT_GREEN = "\033[1;32m"

func SaveResults(app *app.App, detector Detector) string {
	detector.ComputeResults()
	results := detector.GetResults()
	analysisTypeString := detector.GetTypeString()
	analysisPrefix := strings.ToUpper(analysisTypeString)

	// ensure the path for the results file exists
	path := fmt.Sprintf("output/%s/analysis/%s.txt", app.GetName(), detector.GetTypeString())
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		log.Fatalf("[%s] error creating directory %s: %s", analysisPrefix, dir, err.Error())
	}

	// read previous file (if it exists) and check if results have changed
	var previousContent []byte
	if _, err := os.Stat(path); err == nil {
		previousContent, err = os.ReadFile(path)
		if err != nil {
			log.Fatalf("[%s] error reading existing file %s: %s", analysisPrefix, path, err.Error())
		}
	}

	color := TEXT_BOLD_LIGHT_BLUE
	if string(previousContent) == results {
		return color + "\t\t\t\t\t\t\t (unmodified) \n" + results + TEXT_RESET_COLOR + "\n\n"
	}

	// if content changed but the analysis summary is the same then the result is printed in yellow
	// otherwise (i.e., changes in both content and analysis summary), result is printed in red
	color = TEXT_BOLD_LIGHT_RED

	// if we have a new file or the content has changed then update the results
	err = os.WriteFile(path, []byte(results), 0644)
	if err != nil {
		log.Fatalf("[%s] error writing data to %s: %s", analysisPrefix, path, err.Error())
	}

	str := fmt.Sprintf("%s\t\t\t\t\t\t\t   (modified)\n%s%s\n\n", color, results, TEXT_RESET_COLOR)

	return str
}
