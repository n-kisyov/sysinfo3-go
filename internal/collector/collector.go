package collector

import (
	"fmt"
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
	SwapUsed    uint64  `json:"swap_used_bytes" yaml:"swap_used_bytes"`
	SwapUsedH   string  `json:"swap_used" yaml:"swap_used"`
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
	Model string `json:"model" yaml:"model"`
	Type  string `json:"type" yaml:"type"`
	Size  uint64 `json:"size_bytes" yaml:"size_bytes"`
	SizeH string `json:"size" yaml:"size"`
}

type NetInterface struct {
	Name      string   `json:"name" yaml:"name"`
	MAC       string   `json:"mac" yaml:"mac"`
	Addresses []string `json:"addresses" yaml:"addresses"`
	MTU       int      `json:"mtu" yaml:"mtu"`
	Flags     []string `json:"flags" yaml:"flags"`
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
