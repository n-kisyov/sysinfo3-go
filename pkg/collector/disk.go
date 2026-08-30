package collector

import (
	"strings"
	"time"

	"github.com/StackExchange/wmi"
	"github.com/shirou/gopsutil/v3/disk"
)

type wmiPhysicalDisk struct {
	Name      string
	Model     string
	MediaType string
	Size      uint64
}

var (
	lastDiskCounters map[string]disk.IOCountersStat
	lastDiskTime     time.Time
)

func CollectDisks() ([]DiskInfo, []PhysicalDiskInfo) {
	partitions, err := disk.Partitions(false)
	if err != nil {
		return nil, nil
	}

	var physicalDisks []wmiPhysicalDisk
	wmi.Query("SELECT Name, Model, MediaType, Size FROM Win32_DiskDrive", &physicalDisks)

	now := time.Now()
	var timeDiffSec float64
	if !lastDiskTime.IsZero() {
		timeDiffSec = now.Sub(lastDiskTime).Seconds()
	}

	ioCounters, _ := disk.IOCounters() // Map of disk names to counters

	var physResult []PhysicalDiskInfo
	for _, p := range physicalDisks {
		mediaType := strings.TrimSpace(p.MediaType)
		if mediaType == "" {
			mediaType = "Unknown"
		}
		
		phys := PhysicalDiskInfo{
			Name:  strings.TrimSpace(p.Name),
			Model: strings.TrimSpace(p.Model),
			Type:  mediaType,
			Size:  p.Size,
			SizeH: FormatBytes(p.Size),
		}

		// Try to match WMI Name (e.g., \\.\PHYSICALDRIVE0) with gopsutil name (e.g., PhysicalDrive0)
		ioName := strings.ReplaceAll(phys.Name, "\\\\.\\PHYSICALDRIVE", "PhysicalDrive")
		
		if ioCounters != nil {
			if io, ok := ioCounters[ioName]; ok {
				phys.ReadBytes = io.ReadBytes
				phys.WriteBytes = io.WriteBytes
				phys.ReadBytesH = FormatBytes(io.ReadBytes)
				phys.WriteBytesH = FormatBytes(io.WriteBytes)

				if timeDiffSec > 0 && lastDiskCounters != nil {
					if last, ok := lastDiskCounters[ioName]; ok {
						readDiff := io.ReadBytes - last.ReadBytes
						writeDiff := io.WriteBytes - last.WriteBytes
						phys.ReadBytesPerSec = uint64(float64(readDiff) / timeDiffSec)
						phys.WriteBytesPerSec = uint64(float64(writeDiff) / timeDiffSec)
						phys.ReadBytesPerSecH = FormatBytes(phys.ReadBytesPerSec) + "/s"
						phys.WriteBytesPerSecH = FormatBytes(phys.WriteBytesPerSec) + "/s"
					}
				}
			}
		}

		physResult = append(physResult, phys)
	}

	lastDiskCounters = ioCounters
	lastDiskTime = now

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
