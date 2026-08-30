package collector

import (
	"time"

	"github.com/shirou/gopsutil/v3/net"
)

var (
	lastNetworkCounters map[string]net.IOCountersStat
	lastNetworkTime     time.Time
)

func CollectNetwork() []NetInterface {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil
	}

	// Build a map of I/O counters keyed by interface name.
	ioMap := make(map[string]net.IOCountersStat)
	counters, err := net.IOCounters(true) // true = per-interface
	if err == nil {
		for _, c := range counters {
			ioMap[c.Name] = c
		}
	}

	var result []NetInterface
	now := time.Now()
	var timeDiffSec float64
	if !lastNetworkTime.IsZero() {
		timeDiffSec = now.Sub(lastNetworkTime).Seconds()
	}

	for _, iface := range ifaces {
		addrs := make([]string, len(iface.Addrs))
		for i, a := range iface.Addrs {
			addrs[i] = a.Addr
		}

		ni := NetInterface{
			Name:      iface.Name,
			MAC:       iface.HardwareAddr,
			Addresses: addrs,
			MTU:       iface.MTU,
			Flags:     iface.Flags,
		}

		if io, ok := ioMap[iface.Name]; ok {
			ni.BytesSent = io.BytesSent
			ni.BytesRecv = io.BytesRecv
			ni.BytesSentH = FormatBytes(io.BytesSent)
			ni.BytesRecvH = FormatBytes(io.BytesRecv)

			if timeDiffSec > 0 && lastNetworkCounters != nil {
				if last, ok := lastNetworkCounters[iface.Name]; ok {
					sentDiff := io.BytesSent - last.BytesSent
					recvDiff := io.BytesRecv - last.BytesRecv
					ni.BytesSentPerSec = uint64(float64(sentDiff) / timeDiffSec)
					ni.BytesRecvPerSec = uint64(float64(recvDiff) / timeDiffSec)
					ni.BytesSentPerSecH = FormatBytes(ni.BytesSentPerSec) + "/s"
					ni.BytesRecvPerSecH = FormatBytes(ni.BytesRecvPerSec) + "/s"
				}
			}
		}

		result = append(result, ni)
	}

	lastNetworkCounters = ioMap
	lastNetworkTime = now

	return result
}
