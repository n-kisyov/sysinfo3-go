package collector

import (
	"sort"

	"github.com/shirou/gopsutil/v3/process"
)

func CollectProcesses() []ProcessInfo {
	procs, err := process.Processes()
	if err != nil {
		return nil
	}

	var results []ProcessInfo
	for _, p := range procs {
		name, err := p.Name()
		if err != nil || name == "" {
			continue
		}

		cpu, _ := p.CPUPercent()
		memInfo, err := p.MemoryInfo()
		var mem uint64
		if err == nil && memInfo != nil {
			mem = memInfo.RSS
		}

		results = append(results, ProcessInfo{
			PID:         p.Pid,
			Name:        name,
			CPUPercent:  cpu,
			MemoryBytes: mem,
			MemoryH:     FormatBytes(mem),
		})
	}

	// Sort by memory usage as the primary metric, since CPU tracking needs time to measure properly
	sort.Slice(results, func(i, j int) bool {
		if results[i].MemoryBytes == results[j].MemoryBytes {
			return results[i].CPUPercent > results[j].CPUPercent
		}
		return results[i].MemoryBytes > results[j].MemoryBytes
	})

	limit := 5
	if len(results) < limit {
		limit = len(results)
	}
	return results[:limit]
}
