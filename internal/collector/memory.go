package collector

import (
	"github.com/shirou/gopsutil/v3/mem"
)

func CollectMemory() MemoryInfo {
	vmem, err := mem.VirtualMemory()
	if err != nil {
		return MemoryInfo{}
	}

	info := MemoryInfo{
		Total:       vmem.Total,
		TotalH:      FormatBytes(vmem.Total),
		Used:        vmem.Used,
		UsedH:       FormatBytes(vmem.Used),
		Available:   vmem.Available,
		AvailableH:  FormatBytes(vmem.Available),
		UsedPercent: vmem.UsedPercent,
	}

	swap, err := mem.SwapMemory()
	if err == nil && swap != nil {
		info.SwapTotal = swap.Total
		info.SwapTotalH = FormatBytes(swap.Total)
		info.SwapUsed = swap.Used
		info.SwapUsedH = FormatBytes(swap.Used)
		if swap.Total > 0 {
			info.SwapUsedPercent = float64(swap.Used) / float64(swap.Total) * 100.0
		}
	}

	return info
}
