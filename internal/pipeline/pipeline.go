// Package pipeline runs Aletheia's analysis of a registered Blueprint application: SSA taint
// propagation, abstract call graph, schema building and pattern detection
package pipeline

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/aletheia-microservices/aletheia/internal/analysis/service-level/ssagraph"
	"github.com/aletheia-microservices/aletheia/internal/analysis/service-level/ssagraph/parser"
	"github.com/aletheia-microservices/aletheia/internal/analysis/service-level/ssagraph/registry"
	"github.com/aletheia-microservices/aletheia/internal/analysis/service-level/ssagraph/tainter"
	"github.com/aletheia-microservices/aletheia/internal/analysis/system-level/abstractgraph"
	abstractgraphparser "github.com/aletheia-microservices/aletheia/internal/analysis/system-level/abstractgraph/parser"
	"github.com/aletheia-microservices/aletheia/internal/analysis/system-level/detection"
	"github.com/aletheia-microservices/aletheia/internal/analysis/system-level/detection/constraints/foreignkeycascade"
	"github.com/aletheia-microservices/aletheia/internal/analysis/system-level/detection/constraints/foreignkeyconcurrency"
	"github.com/aletheia-microservices/aletheia/internal/analysis/system-level/detection/constraints/keycoordination"
	"github.com/aletheia-microservices/aletheia/internal/analysis/system-level/detection/constraints/uniquenessconcurrency"
	"github.com/aletheia-microservices/aletheia/internal/app"
	appparser "github.com/aletheia-microservices/aletheia/internal/app/parser"
	"github.com/aletheia-microservices/aletheia/internal/config"
	"github.com/aletheia-microservices/aletheia/internal/utils"
)

// Options configures a run of the pipeline
type Options struct {
	// App is the name of the application in registry/apps.yaml
	App string
	// InitOnly only loads the app's Blueprint wiring and returns without analyzing it
	InitOnly bool
	// Synthetic marks the app as synthetic
	Synthetic bool
	// InputRefs reads extra foreign keys from input/{app}/*.yaml
	InputRefs bool
	// DetectionConfig is the path of a YAML file listing warnings to suppress (ignored if empty)
	DetectionConfig string
	// WriteOutputs saves the SSA code, app and schema JSON and detector results to output/{app}/
	WriteOutputs bool
	// Eval skips the intermediate outputs (only relevant with WriteOutputs)
	Eval bool
	// Debug also saves the SSA graphs and the abstract call graph as .dot files
	// (only relevant with WriteOutputs and without Eval)
	Debug bool
	// KeepGraphs keeps the SSA graphs in Result.FuncGraphs instead of releasing them once the
	// abstract call graph is built
	KeepGraphs bool
}

// Timings holds the time elapsed from the start of the pipeline until the end of each stage,
// except for SSATainting, Schema and Detection, which only measure their own stage
type Timings struct {
	Blueprint   time.Duration
	SSAParsing  time.Duration
	SSATainting time.Duration
	Parsing     time.Duration
	Schema      time.Duration
	Detection   time.Duration
	Total       time.Duration
}

// Result holds the state of every stage of the pipeline
type Result struct {
	App *app.App
	// FuncGraphs is nil unless Options.KeepGraphs is set
	FuncGraphs map[string]*ssagraph.SSAGraph
	AbsGraph   *abstractgraph.AbstractCallGraph
	Detectors  []detection.Detector
	// Results holds the results of each detector indexed by its type string
	// (e.g., "foreign-key-cascade")
	Results map[string]string
	// Summaries holds the formatted results of each detector, compared against the previous
	// output (only set with Options.WriteOutputs)
	Summaries []string
	Timings   Timings
}

// Run runs every stage of the pipeline for opts.App
func Run(opts Options) (*Result, error) {
	// the detection config is global state, so reset it for every run
	detection.Config = detection.InputConfig{}
	if opts.DetectionConfig != "" {
		detection.LoadInputConfig(opts.App, opts.DetectionConfig)
	}

	start := time.Now()
	log := logrus.WithField("app", opts.App)
	res := &Result{}

	// ------------ PART 1
	log.WithField("synthetic", opts.Synthetic).Infof("[1/12] initializing program")

	apppath := utils.GetAppRootPackagePath(opts.App)
	a := app.NewApp(opts.App)
	res.App = a
	appparser.Init(a, opts.Synthetic)

	if opts.InitOnly {
		return res, nil
	}

	res.Timings.Blueprint = time.Since(start)

	writeIntermediate := opts.WriteOutputs && !opts.Eval
	writeDebug := writeIntermediate && opts.Debug

	if opts.WriteOutputs {
		dirs := []string{fmt.Sprintf("output/%s", opts.App)}
		if opts.Debug {
			dirs = append(dirs,
				fmt.Sprintf("output/%s/ssagraphs/tainted", opts.App),
				fmt.Sprintf("output/%s/ssagraphs/untainted", opts.App),
				fmt.Sprintf("output/%s/abstractcallgraph", opts.App),
				fmt.Sprintf("output/%s/ssa", opts.App),
			)
		}
		for _, dir := range dirs {
			if err := os.MkdirAll(dir, os.ModePerm); err != nil {
				return nil, err
			}
		}
	}

	// ------------ PART 2
	log.Infof("[2/12] building program")

	prog, pkgs, err := utils.BuildProgram(apppath)
	if err != nil {
		return nil, fmt.Errorf("building program for %s: %w", opts.App, err)
	}

	appparser.InitServiceFields(a, pkgs)
	appparser.ParseSQLSchemaFromUserFile(a)
	appparser.ParseNoSQLSchemaFromUserFile(a)
	if opts.InputRefs {
		log.Infof("reading input refs...")
		appparser.ParseUserInputReferences(a)
	}

	funcGraphs := make(map[string]*ssagraph.SSAGraph)

	// ------------ PART 3
	log.Infof("[3/12] running SSA analysis")
	for _, pkg := range pkgs {
		if opts.WriteOutputs {
			parser.RunSSAAnalysis(a, prog, pkg, funcGraphs)
		} else {
			parser.RunSSAAnalysisTo(io.Discard, a, prog, pkg, funcGraphs)
		}
	}

	var graphsLst []*ssagraph.SSAGraph
	for _, graph := range funcGraphs {
		graph.Sort()
		graphsLst = append(graphsLst, graph)
	}

	// ------------ PART 4
	log.Infof("[4/12] registering fields")
	registry.RegisterFields(a, graphsLst)

	res.Timings.SSAParsing = time.Since(start)
	startSSATainting := time.Now()

	if writeIntermediate {
		a.WriteAppToJSON()
		if writeDebug {
			for fn, graph := range funcGraphs {
				graph.WriteToDOTFile(opts.App, fn, false)
			}
		}
	}

	// ------------ PART 5
	log.Infof("[5/12] running SSA tainter for single graphs")
	for _, graph := range funcGraphs {
		tainter.RunTainter(graph)
	}

	res.Timings.SSATainting = time.Since(startSSATainting)

	// ------------ PART 6
	log.Infof("[6/12] combining SSA graphs")
	for _, graph := range funcGraphs {
		tainter.Combine(graph, funcGraphs)
	}

	if writeDebug {
		written := make(map[string]bool)
		for fn, graph := range funcGraphs {
			graph.WriteToDOTFile(opts.App, fn, true)
			written[fn] = true
		}
		for fn, graph := range funcGraphs {
			for _, toGraph := range graph.GetAllCombinedGraphs() {
				newFn := fn + "." + toGraph.GetMethodName()
				if !written[newFn] {
					toGraph.WriteToDOTFile(opts.App, newFn, true)
				}
				written[newFn] = true
			}
		}
	}

	// ------------ PART 7
	log.Infof("[7/12] creating new abstract call graph")
	absgraph := abstractgraph.NewAbstractCallGraph(a)
	for _, entrypoint := range a.GetEntrypointsShortPaths() {
		abstractgraphparser.Parse(absgraph, entrypoint, true, funcGraphs)
	}
	res.AbsGraph = absgraph

	// ------------ PART 8
	if opts.KeepGraphs {
		res.FuncGraphs = funcGraphs
	} else {
		log.Infof("[8/12] releasing memory associated with ssa graph")
		graphsLst = nil
		pkgs = nil
		for fn, graph := range funcGraphs {
			if graph != nil {
				graph.Release()
			}
			delete(funcGraphs, fn)
		}
		funcGraphs = nil
	}

	res.Timings.Parsing = time.Since(start)

	res.Detectors = []detection.Detector{
		keycoordination.NewDetector(keycoordination.DETECTION_TYPE_PRIMARY_KEY),
		keycoordination.NewDetector(keycoordination.DETECTION_TYPE_FOREIGN_KEY),
		foreignkeycascade.NewDetector(),
		foreignkeyconcurrency.NewDetector(),
		uniquenessconcurrency.NewDetector(),
	}
	iterator := detection.NewIterator(a, absgraph, res.Detectors...)

	startSchema := time.Now()
	// ------------ PART 9
	log.Infof("[9/12] starting schema builder")
	iterator.Run(detection.PHASE_1_SCHEMA_BUILDER)
	// ------------ PART 10
	if config.Global.DualPassSchemaBuilding {
		log.Infof("[10/12] starting schema builder (read only)")
		iterator.Run(detection.PHASE_1_SCHEMA_BUILDER_READ_ONLY)
	} else {
		log.Infof("[10/12] skipping schema builder (read only)...")
	}

	res.Timings.Schema = time.Since(startSchema)

	// ------------ PART 11
	log.Infof("[11/12] starting pattern detection")
	startDetection := time.Now()

	// phase 2: one pass for all detectors
	iterator.Run(detection.PHASE_2_PATTERN_DETECTOR)

	res.Timings.Total = time.Since(start)
	res.Timings.Detection = time.Since(startDetection)

	if writeIntermediate {
		if writeDebug {
			// phase 0: dummy pass to generate dot files with taints for debugging
			iterator.Run(detection.PHASE_0_DEBUG)
			absgraph.WriteToDOTFile(opts.App, true)
			absgraph.WriteToDOTFile(opts.App, false)
		}

		a.WriteAppToJSON()
		a.WriteSchemaToJSON()
	}

	// ------------ PART 12
	log.Infof("[12/12] saving results")
	if opts.WriteOutputs {
		res.Summaries = detection.SaveResults(a, res.Detectors...)
	} else {
		for _, detector := range res.Detectors {
			detector.ComputeResults(a)
		}
	}
	res.Results = make(map[string]string, len(res.Detectors))
	for _, detector := range res.Detectors {
		res.Results[detector.GetTypeString()] = detector.GetResults()
	}

	if opts.WriteOutputs {
		a.WriteAppToJSON()
		a.WriteSchemaToJSON()
	}

	return res, nil
}
