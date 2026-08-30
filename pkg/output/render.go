package output

import (
	"fmt"
	"strings"
	"sysinfo3-go/internal/collector"

	"github.com/fatih/color"
)

type Options struct {
	Basic    bool
	Verbose  bool
	NoColor  bool
	Category string
}

var (
	headerColor = color.New(color.FgCyan, color.Bold)
	labelColor  = color.New(color.FgWhite)
	valueColor  = color.New(color.FgHiWhite)
	greenColor  = color.New(color.FgGreen)
	yellowColor = color.New(color.FgYellow)
	redColor    = color.New(color.FgRed)
	warnColor   = yellowColor
)

func setNoColor(v bool) {
	color.NoColor = v
}

type section struct {
	name      string
	filterKey string // alias used by --category; falls back to lowercase name
	render    func(*collector.SystemSnapshot, Options)
}

var sections = []section{
	{"System", "system", renderHost},
	{"Memory", "memory", renderMemory},
	{"CPU", "cpu", renderCPU},
	{"Disks (Volumes)", "disks", renderDisks},
	{"Physical Disks", "physical disks", renderPhysDisks},
	{"GPU", "gpu", renderGPU},
	{"Network", "network", renderNetwork},
	{"Battery", "battery", renderBattery},
	{"Top Processes", "processes", renderProcesses},
	{"BIOS", "bios", renderBIOS},
}

func RenderTerminal(s *collector.SystemSnapshot, opts Options) {
	setNoColor(opts.NoColor)

	filter := make(map[string]bool)
	if opts.Category != "" {
		for _, cat := range strings.Split(opts.Category, ",") {
			filter[strings.ToLower(strings.TrimSpace(cat))] = true
		}
	}

	for _, sec := range sections {
		key := sec.filterKey
		if key == "" {
			key = strings.ToLower(sec.name)
		}
		if len(filter) > 0 && !filter[key] {
			continue
		}
		printSectionHeader(sec.name)
		sec.render(s, opts)
		fmt.Println()
	}
}

func printSectionHeader(name string) {
	headerColor.Printf("\n══ %s ", name)
	remaining := 50 - len(name) - 4
	if remaining > 0 {
		for i := 0; i < remaining; i++ {
			headerColor.Print("═")
		}
	}
	fmt.Println()
}

func lb(key string) string {
	return labelColor.Sprintf("  %-16s", key)
}



func pctColor(pct float64) *color.Color {
	switch {
	case pct >= 85:
		return redColor
	case pct >= 60:
		return yellowColor
	default:
		return greenColor
	}
}

func printBar(pct float64) {
	const width = 20
	filled := int(pct/100.0*float64(width) + 0.5)
	if filled > width {
		filled = width
	}
	if filled < 0 {
		filled = 0
	}

	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	c := pctColor(pct)
	fmt.Printf("    %s %s\n", c.Sprint(bar), c.Sprintf("%.1f%%", pct))
}

func renderHost(s *collector.SystemSnapshot, opts Options) {
	h := s.Host
	fmt.Printf("%s %s\n", lb("Hostname"), valueColor.Sprint(h.Hostname))
	fmt.Printf("%s %s\n", lb("OS"), valueColor.Sprintf("%s (%s)", h.OS, h.PlatformVersion))
	if !opts.Basic {
		fmt.Printf("%s %s\n", lb("Kernel"), valueColor.Sprint(h.KernelVersion))
	}
	fmt.Printf("%s %s\n", lb("Uptime"), valueColor.Sprint(h.UptimeDisplay))

	if opts.Verbose {
		fmt.Printf("%s %s\n", lb("Platform"), valueColor.Sprintf("%s %s", h.Platform, h.PlatformVersion))
		fmt.Printf("%s %s\n", lb("Uptime (raw)"), valueColor.Sprintf("%d s", h.Uptime))
	}
}

func renderMemory(s *collector.SystemSnapshot, opts Options) {
	m := s.Memory
	fmt.Printf("%s %s\n", lb("Total"), valueColor.Sprint(m.TotalH))
	fmt.Printf("%s %s\n", lb("Used"), pctColor(m.UsedPercent).Sprintf("%s (%.1f%%)", m.UsedH, m.UsedPercent))
	fmt.Printf("%s %s\n", lb("Available"), greenColor.Sprint(m.AvailableH))

	if !opts.Basic {
		fmt.Print(lb("Usage Bar"))
		fmt.Println()
		printBar(m.UsedPercent)
	}

	if (!opts.Basic && m.SwapTotal > 0) || opts.Verbose {
		fmt.Printf("%s %s\n", lb("Swap Total"), valueColor.Sprint(m.SwapTotalH))
		fmt.Printf("%s %s\n", lb("Swap Used"), pctColor(m.SwapUsedPercent).Sprintf("%s (%.1f%%)", m.SwapUsedH, m.SwapUsedPercent))
		if opts.Verbose {
			fmt.Print(lb("Swap Bar"))
			fmt.Println()
			printBar(m.SwapUsedPercent)
		}
	}
}

func renderCPU(s *collector.SystemSnapshot, opts Options) {
	cpu := s.CPU
	fmt.Printf("%s %s\n", lb("Model"), valueColor.Sprint(cpu.Model))
	fmt.Printf("%s %s\n", lb("Cores"), valueColor.Sprintf("%d physical / %d logical", cpu.PhysicalCores, cpu.LogicalCores))
	fmt.Printf("%s %s\n", lb("Usage"), pctColor(cpu.UsagePercent).Sprintf("%.1f%%", cpu.UsagePercent))

	if !opts.Basic {
		fmt.Print(lb("Usage Bar"))
		fmt.Println()
		printBar(cpu.UsagePercent)

		tempStr := cpu.Temperature
		if tempStr == "<unavailable>" || tempStr == "" {
			tempStr = warnColor.Sprint("<unavailable>")
		}
		fmt.Printf("%s %s\n", lb("Temperature"), tempStr)
	}

	if opts.Verbose && len(cpu.PerCoreUsage) > 0 {
		fmt.Printf("%s\n", lb("Core Loads"))
		for i, u := range cpu.PerCoreUsage {
			fmt.Printf("    %s %s\n", labelColor.Sprintf("Core %2d", i), pctColor(u).Sprintf("%.1f%%", u))
		}
	}
}

func renderDisks(s *collector.SystemSnapshot, opts Options) {
	for _, d := range s.Disks {
		if opts.Basic {
			fmt.Printf("%s %s\n", lb(d.DriveLetter),
				pctColor(d.UsedPercent).Sprintf("%s free (%.1f%% used)", d.FreeH, d.UsedPercent))
		} else {
			driveLabel := fmt.Sprintf("%s (%s)", d.DriveLetter, d.FSType)
			usageStr := pctColor(d.UsedPercent).Sprintf("%s total / %s free (%.1f%% used)", d.TotalH, d.FreeH, d.UsedPercent)
			fmt.Printf("%s %s\n", lb(driveLabel), usageStr)

			fmt.Print(lb("Usage Bar"))
			fmt.Println()
			printBar(d.UsedPercent)
		}

	}
}

func renderPhysDisks(s *collector.SystemSnapshot, opts Options) {
	if len(s.PhysDisks) == 0 {
		fmt.Printf("%s %s\n", lb("Status"), warnColor.Sprint("<no physical disks detected>"))
		return
	}
	for i, d := range s.PhysDisks {
		label := fmt.Sprintf("Disk %d", i)
		details := valueColor.Sprintf("%s [%s] %s", d.Model, d.Type, d.SizeH)
		fmt.Printf("%s %s\n", lb(label), details)
	}
}

func renderGPU(s *collector.SystemSnapshot, opts Options) {
	if len(s.GPU) == 0 {
		fmt.Printf("%s %s\n", lb("Status"), warnColor.Sprint("<none detected>"))
		return
	}
	for _, g := range s.GPU {
		fmt.Printf("%s %s\n", lb("Name"), valueColor.Sprint(g.Name))
		if !opts.Basic {
			fmt.Printf("%s %s\n", lb("VRAM"), valueColor.Sprint(g.VRAMH))
		}
		if opts.Verbose {
			fmt.Printf("%s %s\n", lb("Driver"), valueColor.Sprint(g.Driver))
		}
	}
}

func renderNetwork(s *collector.SystemSnapshot, opts Options) {
	for _, iface := range s.Network {
		isLoopback := strings.Contains(strings.ToLower(iface.Name), "loopback")
		hasAddrs := len(iface.Addresses) > 0

		// In non-verbose mode, skip loopback and interfaces with no addresses.
		if !opts.Verbose && (isLoopback || !hasAddrs) {
			continue
		}

		addrStr := "<no address>"
		if hasAddrs {
			if opts.Basic {
				addrStr = iface.Addresses[0]
			} else {
				addrStr = strings.Join(iface.Addresses, ", ")
			}
		}
		fmt.Printf("%s %s\n", lb(iface.Name), valueColor.Sprint(addrStr))
		if opts.Verbose {
			fmt.Printf("%s %s\n", lb("  MAC"), valueColor.Sprint(iface.MAC))
			fmt.Printf("%s %s\n", lb("  MTU"), valueColor.Sprint(fmt.Sprintf("%d", iface.MTU)))
			fmt.Printf("%s %s\n", lb("  Flags"), valueColor.Sprint(strings.Join(iface.Flags, ", ")))
			fmt.Printf("%s %s\n", lb("  Sent"), valueColor.Sprint(iface.BytesSentH))
			fmt.Printf("%s %s\n", lb("  Received"), valueColor.Sprint(iface.BytesRecvH))
		}
	}
}

func renderBattery(s *collector.SystemSnapshot, opts Options) {
	if s.Battery == nil {
		fmt.Printf("%s %s\n", lb("Status"), warnColor.Sprint("<no battery detected>"))
		return
	}
	b := s.Battery
	fmt.Printf("%s %s\n", lb("Status"), valueColor.Sprint(b.Status))
	fmt.Printf("%s %s\n", lb("Charge"), pctColor(float64(b.Percentage)).Sprintf("%d%%", b.Percentage))
	if b.TimeLeft != "" {
		fmt.Printf("%s %s\n", lb("Time Left"), valueColor.Sprint(b.TimeLeft))
	}
}

func renderProcesses(s *collector.SystemSnapshot, opts Options) {
	if len(s.Processes) == 0 {
		fmt.Printf("%s %s\n", lb("Status"), warnColor.Sprint("<no process info>"))
		return
	}
	
	for i, p := range s.Processes {
		label := fmt.Sprintf("#%d %s", i+1, p.Name)
		details := valueColor.Sprintf("PID: %d | Mem: %s | CPU: %.1f%%", p.PID, p.MemoryH, p.CPUPercent)
		fmt.Printf("%s %s\n", lb(label), details)
	}
}

func renderBIOS(s *collector.SystemSnapshot, opts Options) {
	b := s.BIOS
	if opts.Basic {
		if b.Manufacturer != "" && b.Model != "" {
			fmt.Printf("%s %s\n", lb("Board"), valueColor.Sprintf("%s %s", b.Manufacturer, b.Model))
		}
		return
	}
	fmt.Printf("%s %s\n", lb("Vendor"), valueColor.Sprint(b.Vendor))
	fmt.Printf("%s %s\n", lb("Version"), valueColor.Sprint(b.Version))
	if b.Date != "" {
		fmt.Printf("%s %s\n", lb("Date"), valueColor.Sprint(b.Date))
	}
	if b.Manufacturer != "" && b.Model != "" {
		fmt.Printf("%s %s\n", lb("Board"), valueColor.Sprintf("%s %s", b.Manufacturer, b.Model))
	}
	if opts.Verbose && b.SerialNumber != "" {
		fmt.Printf("%s %s\n", lb("Serial"), valueColor.Sprint(b.SerialNumber))
	}
}
