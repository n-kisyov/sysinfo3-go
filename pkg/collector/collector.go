package collector

import (
	"fmt"
	"sync"
	"time"
)

type SystemSnapshot struct {
	Timestamp time.Time     `json:"timestamp" yaml:"timestamp"`
	Host      HostInfo      `json:"host" yaml:"host"`
	Memory    MemoryInfo    `json:"memory" yaml:"memory"`
	CPU       CPUInfo       `json:"cpu" yaml:"cpu"`
	Disks     []DiskInfo    `json:"disks" yaml:"disks"`
	PhysDisks []PhysicalDiskInfo `json:"physical_disks" yaml:"physical_disks"`
	Network   []NetInterface `json:"network" yaml:"network"`
	Battery   *BatteryInfo   `json:"battery,omitempty" yaml:"battery,omitempty"`
	Processes []ProcessInfo `json:"processes" yaml:"processes"`
	GPU       []GPUInfo     `json:"gpu" yaml:"gpu"`
	BIOS      BIOSInfo      `json:"bios" yaml:"bios"`
}

type HostInfo struct {
	Hostname        string `json:"hostname" yaml:"hostname"`
	OS              string `json:"os" yaml:"os"`
	Platform        string `json:"platform" yaml:"platform"`
	PlatformVersion string `json:"platform_version" yaml:"platform_version"`
	KernelVersion   string `json:"kernel_version" yaml:"kernel_version"`
	Uptime          uint64 `json:"uptime_seconds" yaml:"uptime_seconds"`
	UptimeDisplay   string `json:"uptime" yaml:"uptime"`
}

type MemoryInfo struct {
	Total       uint64  `json:"total_bytes" yaml:"total_bytes"`
	TotalH      string  `json:"total" yaml:"total"`
	Used        uint64  `json:"used_bytes" yaml:"used_bytes"`
	UsedH       string  `json:"used" yaml:"used"`
	Available   uint64  `json:"available_bytes" yaml:"available_bytes"`
	AvailableH  string  `json:"available" yaml:"available"`
	UsedPercent float64 `json:"used_percent" yaml:"used_percent"`
	SwapTotal   uint64  `json:"swap_total_bytes" yaml:"swap_total_bytes"`
	SwapTotalH  string  `json:"swap_total" yaml:"swap_total"`
	SwapUsed        uint64  `json:"swap_used_bytes" yaml:"swap_used_bytes"`
	SwapUsedH       string  `json:"swap_used" yaml:"swap_used"`
	SwapUsedPercent float64 `json:"swap_used_percent" yaml:"swap_used_percent"`
}

type CPUInfo struct {
	Model         string    `json:"model" yaml:"model"`
	PhysicalCores int       `json:"physical_cores" yaml:"physical_cores"`
	LogicalCores  int       `json:"logical_cores" yaml:"logical_cores"`
	UsagePercent  float64   `json:"usage_percent" yaml:"usage_percent"`
	PerCoreUsage  []float64 `json:"per_core_usage" yaml:"per_core_usage"`
	Temperature   string    `json:"temperature" yaml:"temperature"`
}

type DiskInfo struct {
	DriveLetter   string  `json:"drive_letter" yaml:"drive_letter"`
	MountPoint    string  `json:"mount_point" yaml:"mount_point"`
	FSType        string  `json:"fs_type" yaml:"fs_type"`
	Total         uint64  `json:"total_bytes" yaml:"total_bytes"`
	TotalH        string  `json:"total" yaml:"total"`
	Used          uint64  `json:"used_bytes" yaml:"used_bytes"`
	UsedH         string  `json:"used" yaml:"used"`
	Free          uint64  `json:"free_bytes" yaml:"free_bytes"`
	FreeH         string  `json:"free" yaml:"free"`
	UsedPercent   float64 `json:"used_percent" yaml:"used_percent"`
}

type PhysicalDiskInfo struct {
	Name             string `json:"name" yaml:"name"`
	Model            string `json:"model" yaml:"model"`
	Type             string `json:"type" yaml:"type"`
	Size             uint64 `json:"size_bytes" yaml:"size_bytes"`
	SizeH            string `json:"size" yaml:"size"`
	ReadBytes        uint64 `json:"read_bytes" yaml:"read_bytes"`
	WriteBytes       uint64 `json:"write_bytes" yaml:"write_bytes"`
	ReadBytesH       string `json:"read_bytes_h" yaml:"read_bytes_h"`
	WriteBytesH      string `json:"write_bytes_h" yaml:"write_bytes_h"`
	ReadBytesPerSec  uint64 `json:"read_bytes_per_sec" yaml:"read_bytes_per_sec"`
	WriteBytesPerSec uint64 `json:"write_bytes_per_sec" yaml:"write_bytes_per_sec"`
	ReadBytesPerSecH string `json:"read_bytes_per_sec_h" yaml:"read_bytes_per_sec_h"`
	WriteBytesPerSecH string `json:"write_bytes_per_sec_h" yaml:"write_bytes_per_sec_h"`
}

type NetInterface struct {
	Name       string   `json:"name" yaml:"name"`
	MAC        string   `json:"mac" yaml:"mac"`
	Addresses  []string `json:"addresses" yaml:"addresses"`
	MTU        int      `json:"mtu" yaml:"mtu"`
	Flags      []string `json:"flags" yaml:"flags"`
	BytesSent        uint64   `json:"bytes_sent" yaml:"bytes_sent"`
	BytesRecv        uint64   `json:"bytes_recv" yaml:"bytes_recv"`
	BytesSentH       string   `json:"bytes_sent_h" yaml:"bytes_sent_h"`
	BytesRecvH       string   `json:"bytes_recv_h" yaml:"bytes_recv_h"`
	BytesSentPerSec  uint64   `json:"bytes_sent_per_sec" yaml:"bytes_sent_per_sec"`
	BytesRecvPerSec  uint64   `json:"bytes_recv_per_sec" yaml:"bytes_recv_per_sec"`
	BytesSentPerSecH string   `json:"bytes_sent_per_sec_h" yaml:"bytes_sent_per_sec_h"`
	BytesRecvPerSecH string   `json:"bytes_recv_per_sec_h" yaml:"bytes_recv_per_sec_h"`
}

type BatteryInfo struct {
	Percentage int    `json:"percentage" yaml:"percentage"`
	Status     string `json:"status" yaml:"status"`
	TimeLeft   string `json:"time_left,omitempty" yaml:"time_left,omitempty"`
}

type ProcessInfo struct {
	PID         int32   `json:"pid" yaml:"pid"`
	Name        string  `json:"name" yaml:"name"`
	CPUPercent  float64 `json:"cpu_percent" yaml:"cpu_percent"`
	MemoryBytes uint64  `json:"memory_bytes" yaml:"memory_bytes"`
	MemoryH     string  `json:"memory" yaml:"memory"`
}

type GPUInfo struct {
	Name   string `json:"name" yaml:"name"`
	Driver string `json:"driver" yaml:"driver"`
	VRAM   uint64 `json:"vram_bytes" yaml:"vram_bytes"`
	VRAMH  string `json:"vram" yaml:"vram"`
}

type BIOSInfo struct {
	Vendor       string `json:"vendor" yaml:"vendor"`
	Version      string `json:"version" yaml:"version"`
	Date         string `json:"date" yaml:"date"`
	Manufacturer string `json:"manufacturer" yaml:"manufacturer"`
	Model        string `json:"model" yaml:"model"`
	SerialNumber string `json:"serial_number,omitempty" yaml:"serial_number,omitempty"`
}

func FormatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func FormatUptime(seconds uint64) string {
	if seconds == 0 {
		return "0m"
	}
	days := seconds / 86400
	hours := (seconds % 86400) / 3600
	minutes := (seconds % 3600) / 60
	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}

func CollectAll(static *SystemSnapshot) *SystemSnapshot {
	var hostInfo HostInfo
	var memoryInfo MemoryInfo
	var cpuInfo CPUInfo
	var disks []DiskInfo
	var physDisks []PhysicalDiskInfo
	var network []NetInterface
	var battery *BatteryInfo
	var processes []ProcessInfo
	var gpus []GPUInfo
	var bios BIOSInfo

	cpuCh := make(chan CPUInfo, 1)
	go func() {
		cpuCh <- CollectCPU()
	}()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() { defer wg.Done(); hostInfo = CollectHost() }()
	wg.Add(1)
	go func() { defer wg.Done(); memoryInfo = CollectMemory() }()
	wg.Add(1)
	go func() { defer wg.Done(); disks, physDisks = CollectDisks() }()
	wg.Add(1)
	go func() { defer wg.Done(); network = CollectNetwork() }()
	wg.Add(1)
	go func() { defer wg.Done(); battery = CollectBattery() }()
	wg.Add(1)
	go func() { defer wg.Done(); processes = CollectProcesses() }()

	collectStatic := static == nil
	if collectStatic {
		wg.Add(1)
		go func() { defer wg.Done(); gpus = CollectGPU() }()
		wg.Add(1)
		go func() { defer wg.Done(); bios = CollectBIOS() }()
	} else {
		gpus = static.GPU
		bios = static.BIOS
	}

	wg.Wait()
	cpuInfo = <-cpuCh

	return &SystemSnapshot{
		Timestamp: time.Now(),
		Host:      hostInfo,
		Memory:    memoryInfo,
		CPU:       cpuInfo,
		Disks:     disks,
		PhysDisks: physDisks,
		Network:   network,
		Battery:   battery,
		Processes: processes,
		GPU:       gpus,
		BIOS:      bios,
	}
}
