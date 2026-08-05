package collector

import (
	"github.com/shirou/gopsutil/v3/mem"
)

func CollectMemory() MemoryInfo {
	vmem, err := mem.VirtualMemory()
	if err != nil {
		return MemoryInfo{}
	}

	swap, _ := mem.SwapMemory()

	return MemoryInfo{
		Total:       vmem.Total,
		TotalH:      FormatBytes(vmem.Total),
		Used:        vmem.Used,
		UsedH:       FormatBytes(vmem.Used),
		Available:   vmem.Available,
		AvailableH:  FormatBytes(vmem.Available),
		UsedPercent: vmem.UsedPercent,
		SwapTotal:   swap.Total,
		SwapTotalH:  FormatBytes(swap.Total),
		SwapUsed:    swap.Used,
		SwapUsedH:   FormatBytes(swap.Used),
	}
}
