package detector

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"analyzer/pkg/abstractgraph"
	"analyzer/pkg/app"
	"analyzer/pkg/logger"
)

const TEXT_BOLD_LIGHT_YELLOW = "\033[1;38;5;179m"
const TEXT_BOLD_LIGHT_RED = "\033[1;31m"
const TEXT_RESET_COLOR = "\033[0m"
const TEXT_BOLD_LIGHT_BLUE = "\033[1;38;5;75m"
const TEXT_BOLD_LIGHT_GREEN = "\033[1;32m"

type Detector interface {
	GetResults() string
	GetSummary() string
	SetSummary(summary string)
	ComputeResults()
	GetAnalysisTypeString() string

	OnNewRun(app *app.App)
	OnEndRun(app *app.App)
	OnNewRequest(entryNode *abstractgraph.AbstractServiceCall)
	OnEndRequest(app *app.App)
	OnNewNode(app *app.App, node abstractgraph.AbstractNode)
	OnEndNode(app *app.App, node abstractgraph.AbstractNode)
	OnRead(app *app.App, dbCall *abstractgraph.AbstractDatabaseCall, lastServiceCallNode *abstractgraph.AbstractServiceCall, child_idx int)
	OnWrite(app *app.App, dbCall *abstractgraph.AbstractDatabaseCall, lastServiceCallNode *abstractgraph.AbstractServiceCall, child_idx int)
	OnUpdate(app *app.App, dbCall *abstractgraph.AbstractDatabaseCall, lastServiceCallNode *abstractgraph.AbstractServiceCall, child_idx int)
	OnDelete(app *app.App, dbCall *abstractgraph.AbstractDatabaseCall, lastServiceCallNode *abstractgraph.AbstractServiceCall, child_idx int)
}

func SaveResults(app *app.App, detector Detector) string {
	detector.ComputeResults()
	results := detector.GetResults()
	analysisTypeString := detector.GetAnalysisTypeString()
	analysisPrefix := strings.ToUpper(analysisTypeString)

	// ensure the path for the results file exists
	path := fmt.Sprintf("output/%s/analysis/%s.txt", app.Name, detector.GetAnalysisTypeString())
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		logger.Logger.Fatalf("[%s] error creating directory %s: %s", analysisPrefix, dir, err.Error())
	}

	// read previous file (if it exists) and check if results have changed
	var previousContent []byte
	if _, err := os.Stat(path); err == nil {
		previousContent, err = os.ReadFile(path)
		if err != nil {
			logger.Logger.Fatalf("[%s] error reading existing file %s: %s", analysisPrefix, path, err.Error())
		}
	}

	// extract analysis summary from the results - always in line 5, after ">>" and in between ( )
	var summary string
	newLines := strings.Split(results, "\n")
	if len(newLines) >= 5 {
		summary = strings.TrimSpace(strings.TrimPrefix(newLines[4], ">>"))
	}

	color := TEXT_BOLD_LIGHT_BLUE
	if string(previousContent) == results {
		detector.SetSummary(color + strings.ToUpper(detector.GetAnalysisTypeString()) + ": " + summary + TEXT_RESET_COLOR + "\n")
		return color + "(unmodified) \n" + results + TEXT_RESET_COLOR + "\n\n"
	}

	// if content changed but the analysis summary is the same then the result is printed in yellow
	// otherwise (i.e., changes in both content and analysis summary), result is printed in red
	color = TEXT_BOLD_LIGHT_RED
	previousLines := strings.Split(string(previousContent), "\n")
	if len(previousLines) >= 5 && len(newLines) >= 5 {
		previousHeader := strings.TrimSpace(strings.TrimPrefix(previousLines[4], ">>"))
		newHeader := strings.TrimSpace(strings.TrimPrefix(newLines[4], ">>"))

		if previousHeader == newHeader {
			color = TEXT_BOLD_LIGHT_YELLOW // Apply yellow only if summaries are identical
		}
	}

	// if we have a new file or the content has changed then update the results
	logger.Logger.Warnf("[%s] WARNING: detected modified results for %s", analysisPrefix, path)
	err = os.WriteFile(path, []byte(results), 0644)
	if err != nil {
		logger.Logger.Fatalf("[%s] error writing data to %s: %s", analysisPrefix, path, err.Error())
	}
	logger.Logger.Tracef("[%s] saved cascading detection results to %s", analysisPrefix, path)

	detector.SetSummary(color + strings.ToUpper(detector.GetAnalysisTypeString()) + ": " + summary + TEXT_RESET_COLOR + "\n")
	return fmt.Sprintf("%s(modified)\n%s%s\n\n", color, results, TEXT_RESET_COLOR)
}
