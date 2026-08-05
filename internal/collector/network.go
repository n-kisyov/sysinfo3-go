package collector

import (
	"github.com/shirou/gopsutil/v3/net"
)

func CollectNetwork() []NetInterface {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil
	}

	var result []NetInterface
	for _, iface := range ifaces {
		addrs := make([]string, len(iface.Addrs))
		for i, a := range iface.Addrs {
			addrs[i] = a.Addr
		}

		result = append(result, NetInterface{
			Name:      iface.Name,
			MAC:       iface.HardwareAddr,
			Addresses: addrs,
			MTU:       iface.MTU,
			Flags:     iface.Flags,
		})
	}
	return result
}
