package collector

import (
	"github.com/StackExchange/wmi"
	"github.com/shirou/gopsutil/v3/disk"
)

type wmiPhysicalDisk struct {
	Model  string
	MediaType string
	Size   uint64
}

func CollectDisks() []DiskInfo {
	partitions, err := disk.Partitions(false)
	if err != nil {
		return nil
	}

	var physicalDisks []wmiPhysicalDisk
	wmi.Query("SELECT Model, MediaType, Size FROM Win32_DiskDrive", &physicalDisks)

	var result []DiskInfo
	for _, p := range partitions {
		usage, err := disk.Usage(p.Mountpoint)
		if err != nil {
			continue
		}

		di := DiskInfo{
			DriveLetter: p.Device,
			MountPoint:  p.Mountpoint,
			FSType:      p.Fstype,
			Total:       usage.Total,
			TotalH:      FormatBytes(usage.Total),
			Used:        usage.Used,
			UsedH:       FormatBytes(usage.Used),
			Free:        usage.Free,
			FreeH:       FormatBytes(usage.Free),
			UsedPercent: usage.UsedPercent,
		}

		// Only attach physical info when there's exactly one physical disk,
		// since we cannot reliably map partitions to physical disks without
		// the full WMI association chain.
		if len(physicalDisks) == 1 {
			di.PhysicalModel = physicalDisks[0].Model
			di.PhysicalType = physicalDisks[0].MediaType
			if physicalDisks[0].Size > 0 {
				di.PhysicalSizeH = FormatBytes(physicalDisks[0].Size)
			}
		}

		result = append(result, di)
	}

	return result
}
