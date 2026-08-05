package collector

import (
	"github.com/shirou/gopsutil/v3/host"
)

func CollectHost() HostInfo {
	info, err := host.Info()
	if err != nil {
		return HostInfo{
			Hostname:      "<error>",
			OS:            "<error>",
			UptimeDisplay: "0m",
		}
	}

	uptime, _ := host.Uptime()

	osName := info.Platform
	if osName == "" {
		osName = info.OS
	}

	return HostInfo{
		Hostname:        info.Hostname,
		OS:              osName,
		Platform:        info.Platform,
		PlatformVersion: info.PlatformVersion,
		KernelVersion:   info.KernelVersion,
		Uptime:          uptime,
		UptimeDisplay:   FormatUptime(uptime),
	}
}
