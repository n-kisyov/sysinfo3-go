package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"time"

	"sysinfo3-go/internal/collector"
	"sysinfo3-go/internal/output"
)

func main() {
	basic := flag.Bool("basic", false, "Show minimal system summary")
	verbose := flag.Bool("verbose", false, "Show all available details")
	jsonOut := flag.Bool("json", false, "Output as JSON")
	csvOut := flag.Bool("csv", false, "Output as CSV")
	yamlOut := flag.Bool("yaml", false, "Output as YAML")
	watch := flag.Bool("watch", false, "Continuously refresh system info")
	interval := flag.Int("interval", 2, "Refresh interval in seconds (with --watch)")
	noColor := flag.Bool("no-color", false, "Disable ANSI color output")
	category := flag.String("category", "", "Show only specified categories (comma-separated)")
	outputFile := flag.String("output", "", "Write structured output to file (use with --json/--csv/--yaml)")
	topN := flag.Int("top", 5, "Number of top processes to display")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "sysinfo3-go — Windows system information tool\n\n")
		fmt.Fprintf(os.Stderr, "Usage: sysinfo3-go [flags]\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  sysinfo3-go                        Full system overview\n")
		fmt.Fprintf(os.Stderr, "  sysinfo3-go --basic                Minimal summary\n")
		fmt.Fprintf(os.Stderr, "  sysinfo3-go --verbose              All details\n")
		fmt.Fprintf(os.Stderr, "  sysinfo3-go --watch --interval 3   Live refresh every 3s\n")
		fmt.Fprintf(os.Stderr, "  sysinfo3-go --json                 JSON output\n")
		fmt.Fprintf(os.Stderr, "  sysinfo3-go --json --output sys.json  Save to file\n")
		fmt.Fprintf(os.Stderr, "  sysinfo3-go --category cpu,memory  Only CPU and Memory\n")
		fmt.Fprintf(os.Stderr, "  sysinfo3-go --top 10               Show top 10 processes\n")
	}

	flag.Parse()

	structCount := 0
	if *jsonOut {
		structCount++
	}
	if *csvOut {
		structCount++
	}
	if *yamlOut {
		structCount++
	}
	if structCount > 1 {
		fmt.Fprintln(os.Stderr, "Error: --json, --csv, and --yaml are mutually exclusive")
		os.Exit(1)
	}

	if structCount > 0 && *watch {
		fmt.Fprintln(os.Stderr, "Error: --watch cannot be used with structured output formats")
		os.Exit(1)
	}

	if *basic && *verbose {
		fmt.Fprintln(os.Stderr, "Error: --basic and --verbose are mutually exclusive")
		os.Exit(1)
	}

	if *outputFile != "" && structCount == 0 {
		fmt.Fprintln(os.Stderr, "Error: --output requires --json, --csv, or --yaml")
		os.Exit(1)
	}

	if *topN < 1 {
		*topN = 1
	}
	collector.TopN = *topN

	opts := output.Options{
		Basic:    *basic,
		Verbose:  *verbose,
		NoColor:  *noColor,
		Category: *category,
	}

	switch {
	case *jsonOut:
		snap := collector.CollectAll(nil)
		w := getWriter(*outputFile)
		if f, ok := w.(*os.File); ok && f != os.Stdout {
			defer f.Close()
		}
		output.RenderJSON(snap, w)
		return
	case *csvOut:
		snap := collector.CollectAll(nil)
		w := getWriter(*outputFile)
		if f, ok := w.(*os.File); ok && f != os.Stdout {
			defer f.Close()
		}
		output.RenderCSV(snap, w)
		return
	case *yamlOut:
		snap := collector.CollectAll(nil)
		w := getWriter(*outputFile)
		if f, ok := w.(*os.File); ok && f != os.Stdout {
			defer f.Close()
		}
		output.RenderYAML(snap, w)
		return
	case *watch:
		runWatch(*interval, opts)
		return
	default:
		snap := collector.CollectAll(nil)
		output.RenderTerminal(snap, opts)
	}
}


func runWatch(intervalSec int, opts output.Options) {
	if intervalSec < 1 {
		intervalSec = 1
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	static := collector.CollectAll(nil)
	clearScreen()
	output.RenderTerminal(static, opts)

	ticker := time.NewTicker(time.Duration(intervalSec) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-sigCh:
			fmt.Print("\033[?25h")
			fmt.Println()
			return
		case <-ticker.C:
			dynamic := collector.CollectAll(static)
			clearScreen()
			output.RenderTerminal(dynamic, opts)
		}
	}
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
	fmt.Print("\033[?25l")
}

func getWriter(outputFile string) io.Writer {
	if outputFile == "" {
		return os.Stdout
	}
	f, err := os.Create(outputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
		os.Exit(1)
	}
	return f
}
