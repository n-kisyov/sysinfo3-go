package collector

import (
	"github.com/shirou/gopsutil/v3/net"
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
		}

		result = append(result, ni)
	}
	return result
}
