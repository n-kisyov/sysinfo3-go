package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
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
		fmt.Fprintf(os.Stderr, "  sysinfo3-go --category cpu,memory  Only CPU and Memory\n")
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

	opts := output.Options{
		Basic:    *basic,
		Verbose:  *verbose,
		NoColor:  *noColor,
		Category: *category,
	}

	switch {
	case *jsonOut:
		snap := collectAll(nil)
		output.RenderJSON(snap)
		return
	case *csvOut:
		snap := collectAll(nil)
		output.RenderCSV(snap)
		return
	case *yamlOut:
		snap := collectAll(nil)
		output.RenderYAML(snap)
		return
	case *watch:
		runWatch(*interval, opts)
		return
	default:
		snap := collectAll(nil)
		output.RenderTerminal(snap, opts)
	}
}

func collectAll(static *collector.SystemSnapshot) *collector.SystemSnapshot {
	var hostInfo collector.HostInfo
	var memoryInfo collector.MemoryInfo
	var cpuInfo collector.CPUInfo
	var disks []collector.DiskInfo
	var network []collector.NetInterface
	var gpus []collector.GPUInfo
	var bios collector.BIOSInfo

	cpuCh := make(chan collector.CPUInfo, 1)
	go func() {
		cpuCh <- collector.CollectCPU()
	}()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() { defer wg.Done(); hostInfo = collector.CollectHost() }()
	wg.Add(1)
	go func() { defer wg.Done(); memoryInfo = collector.CollectMemory() }()
	wg.Add(1)
	go func() { defer wg.Done(); disks = collector.CollectDisks() }()
	wg.Add(1)
	go func() { defer wg.Done(); network = collector.CollectNetwork() }()

	collectStatic := static == nil
	if collectStatic {
		wg.Add(1)
		go func() { defer wg.Done(); gpus = collector.CollectGPU() }()
		wg.Add(1)
		go func() { defer wg.Done(); bios = collector.CollectBIOS() }()
	} else {
		gpus = static.GPU
		bios = static.BIOS
	}

	wg.Wait()
	cpuInfo = <-cpuCh

	return &collector.SystemSnapshot{
		Timestamp: time.Now(),
		Host:      hostInfo,
		Memory:    memoryInfo,
		CPU:       cpuInfo,
		Disks:     disks,
		Network:   network,
		GPU:       gpus,
		BIOS:      bios,
	}
}

func runWatch(intervalSec int, opts output.Options) {
	if intervalSec < 1 {
		intervalSec = 1
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	static := collectAll(nil)
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
			dynamic := collectAll(static)
			clearScreen()
			output.RenderTerminal(dynamic, opts)
		}
	}
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
	fmt.Print("\033[?25l")
}
