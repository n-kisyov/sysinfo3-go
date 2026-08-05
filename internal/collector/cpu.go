package collector

import (
	"fmt"
	"runtime"
	"time"

	"github.com/StackExchange/wmi"
	"github.com/shirou/gopsutil/v3/cpu"
)

type wmiCPU struct {
	Name                      string
	NumberOfCores             *uint32
	NumberOfLogicalProcessors *uint32
}

func CollectCPU() CPUInfo {
	var physicalCores, logicalCores int
	var model string

	var wmiCPUs []wmiCPU
	if err := wmi.Query("SELECT Name, NumberOfCores, NumberOfLogicalProcessors FROM Win32_Processor", &wmiCPUs); err == nil {
		for _, c := range wmiCPUs {
			model = c.Name
			if c.NumberOfCores != nil {
				physicalCores = int(*c.NumberOfCores)
			}
			if c.NumberOfLogicalProcessors != nil {
				logicalCores = int(*c.NumberOfLogicalProcessors)
			}
			break
		}
	}

	if model == "" {
		cpuInfos, err := cpu.Info()
		if err == nil && len(cpuInfos) > 0 {
			model = cpuInfos[0].ModelName
		}
	}
	if physicalCores == 0 {
		physicalCores = runtime.NumCPU()
	}
	if logicalCores == 0 {
		logicalCores = runtime.NumCPU()
	}

	usage, _ := cpu.Percent(1*time.Second, false)
	var usagePercent float64
	if len(usage) > 0 {
		usagePercent = usage[0]
	}

	var perCoreUsage []float64
	perCoreUsage, _ = cpu.Percent(1*time.Second, true)

	temperature := collectCPUTemperature()

	return CPUInfo{
		Model:         model,
		PhysicalCores: physicalCores,
		LogicalCores:  logicalCores,
		UsagePercent:  usagePercent,
		PerCoreUsage:  perCoreUsage,
		Temperature:   temperature,
	}
}

func collectCPUTemperature() string {
	type wmiTemp struct {
		CurrentReading *uint32
	}
	var probes []wmiTemp
	err := wmi.Query("SELECT * FROM Win32_TemperatureProbe", &probes)
	if err != nil || len(probes) == 0 {
		return "<unavailable>"
	}
	for _, p := range probes {
		if p.CurrentReading != nil {
			t := float64(*p.CurrentReading) / 10.0
			return fmt.Sprintf("%.1f °C", t)
		}
	}
	return "<unavailable>"
}
