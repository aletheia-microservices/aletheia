// Command aletheia analyzes a Blueprint application registered in registry/apps.yaml for
// integrity violations across microservices and saves the results in output/<app>/
package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	blueprint_apps "github.com/aletheia-microservices/aletheia/internal/frameworks/blueprint/apps"
	"github.com/aletheia-microservices/aletheia/internal/pipeline"
)

const evalMetricsBase = "eval/metrics"

func printUsage() {
	out := flag.CommandLine.Output()
	fmt.Fprintf(out, "usage: aletheia [flags] <app>\n\n")
	fmt.Fprintf(out, "Analyzes a Blueprint application registered in registry/apps.yaml and\nsaves the results in output/<app>/. Run it from the repository root.\n\n")
	fmt.Fprintln(out, "flags:")
	flag.PrintDefaults()
	fmt.Fprintln(out, "\nregistered apps:")
	for _, name := range blueprint_apps.AppNames() {
		fmt.Fprintf(out, "  %s\n", name)
	}
}

func main() {
	opts := pipeline.Options{WriteOutputs: true}
	flag.BoolVar(&opts.InitOnly, "init", false, "only load the app's Blueprint wiring, then exit without analyzing it")
	flag.BoolVar(&opts.Eval, "eval", false, "evaluation mode: skip intermediate outputs, print timings and save them in "+evalMetricsBase+"/")
	flag.BoolVar(&opts.Synthetic, "synthetic", false, "mark the app as synthetic (only changes where -eval saves timings)")
	flag.BoolVar(&opts.InputRefs, "refs", false, "read extra foreign keys from input/<app>/*.yaml\n(one per line: FOREIGN_KEY db.table.field REFERENCES db.table.field)")
	flag.BoolVar(&opts.Debug, "debug", false, "also save the SSA graphs and the abstract call graph as .dot files, and print timings")
	flag.StringVar(&opts.DetectionConfig, "detection_config", "", "YAML file listing warnings to suppress (examples in config/)")
	flag.Usage = printUsage
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	opts.App = flag.Arg(0)
	if _, ok := blueprint_apps.APPS_INFO[opts.App]; !ok {
		fmt.Fprintf(os.Stderr, "unknown app %q\n\n", opts.App)
		flag.Usage()
		os.Exit(1)
	}

	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.TimeOnly,
	})
	logrus.SetLevel(logrus.InfoLevel)

	if opts.Eval {
		go func() {
			for {
				for _, r := range `-\|/` {
					fmt.Printf("\rRunning... %c", r)
					time.Sleep(1 * time.Second)
				}
			}
		}()
	}

	res, err := pipeline.Run(opts)
	if err != nil {
		logrus.Fatalf("error: %s", err.Error())
	}
	if opts.InitOnly {
		return
	}

	for _, summary := range res.Summaries {
		fmt.Println(summary)
	}

	t := res.Timings
	if opts.Eval || opts.Debug {
		fmt.Printf("Execution time (TOTAL):\t\t%.4f s\n", t.Total.Seconds())
		fmt.Printf("Execution time (BLUEPRINT):\t%.4f s\n", t.Blueprint.Seconds())
		fmt.Printf("Execution time (PARSING):\t%.4f s\n", t.Parsing.Seconds())
		fmt.Printf("Execution time (SSA PARS):\t%.4f s\n", t.SSAParsing.Seconds())
		fmt.Printf("Execution time (SSA TAIN):\t%.4f s\n", t.SSATainting.Seconds())
		fmt.Printf("Execution time (SCHEMA):\t%.4f s\n", t.Schema.Seconds())
		fmt.Printf("Execution time (DETECTION):\t%.4f s\n", t.Detection.Seconds())
	}

	if opts.Eval {
		times := analysisTimes{
			App:              res.App.GetName(),
			NumMicroservices: res.App.NumberOfMicroservices(),
			NumDatastores:    res.App.NumberOfDatastores(),
			NumCallGraphs:    res.AbsGraph.ComputeAndGetNumCallGraphs(),
			Blueprint:        t.Blueprint.Seconds(),
			Rpcs:             res.AbsGraph.GetRPCCount(),
			Total:            t.Total.Seconds(),
			Parsing:          t.Parsing.Seconds(),
			Schema:           t.Schema.Seconds(),
			Detection:        t.Detection.Seconds(),
		}
		if err := saveAnalysisTimes(times, opts.Synthetic); err != nil {
			logrus.Fatalf("error saving analysis times: %s", err.Error())
		}
	}
}

type analysisTimes struct {
	App              string  `yaml:"app"`
	NumMicroservices int     `yaml:"ms_count"`
	NumDatastores    int     `yaml:"ds_count"`
	NumCallGraphs    int     `yaml:"callgraphs"`
	Rpcs             int     `yaml:"rpcs"`
	Blueprint        float64 `yaml:"blueprint"`
	Total            float64 `yaml:"total_s"`
	Parsing          float64 `yaml:"parsing_s"`
	Schema           float64 `yaml:"schema_s"`
	Detection        float64 `yaml:"detection_s"`
}

// saveAnalysisTimes saves times to eval/metrics/{date}/{realistic|synthetic}/{app}_{unix}.yaml
func saveAnalysisTimes(times analysisTimes, synthetic bool) error {
	ts := time.Now().Unix()
	dir := path.Join(evalMetricsBase, time.Now().Format(time.DateOnly))
	if synthetic {
		dir += "/synthetic"
	} else {
		dir += "/realistic"
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	out, err := yaml.Marshal(times)
	if err != nil {
		return err
	}

	filepath := fmt.Sprintf("%s/%s_%d.yaml", dir, times.App, ts)
	return os.WriteFile(filepath, out, 0644)
}
