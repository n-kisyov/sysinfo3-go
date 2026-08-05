package collector

import (
	"github.com/StackExchange/wmi"
)

func CollectGPU() []GPUInfo {
	type wmiGPU struct {
		Name          string
		DriverVersion string
		AdapterRAM    *uint64
	}

	var gpus []wmiGPU
	err := wmi.Query("SELECT Name, DriverVersion, AdapterRAM FROM Win32_VideoController WHERE Name IS NOT NULL", &gpus)
	if err != nil {
		return nil
	}

	var result []GPUInfo
	for _, g := range gpus {
		var vram uint64
		var vramH string
		if g.AdapterRAM != nil {
			vram = *g.AdapterRAM
			vramH = FormatBytes(vram)
		} else {
			vramH = "<unknown>"
		}
		result = append(result, GPUInfo{
			Name:   g.Name,
			Driver: g.DriverVersion,
			VRAM:   vram,
			VRAMH:  vramH,
		})
	}
	return result
}
