package collector

import (
	"strings"

	"github.com/StackExchange/wmi"
	"github.com/shirou/gopsutil/v3/disk"
)

type wmiPhysicalDisk struct {
	Model     string
	MediaType string
	Size      uint64
}

func CollectDisks() ([]DiskInfo, []PhysicalDiskInfo) {
	partitions, err := disk.Partitions(false)
	if err != nil {
		return nil, nil
	}

	var physicalDisks []wmiPhysicalDisk
	wmi.Query("SELECT Model, MediaType, Size FROM Win32_DiskDrive", &physicalDisks)

	var physResult []PhysicalDiskInfo
	for _, p := range physicalDisks {
		mediaType := strings.TrimSpace(p.MediaType)
		if mediaType == "" {
			mediaType = "Unknown"
		}
		physResult = append(physResult, PhysicalDiskInfo{
			Model: strings.TrimSpace(p.Model),
			Type:  mediaType,
			Size:  p.Size,
			SizeH: FormatBytes(p.Size),
		})
	}

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

		result = append(result, di)
	}

	return result, physResult
}
