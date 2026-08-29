package collector

import (
	"fmt"

	"github.com/StackExchange/wmi"
)

func CollectBattery() *BatteryInfo {
	type wmiBattery struct {
		EstimatedChargeRemaining *uint16
		BatteryStatus            *uint16
		EstimatedRunTime         *uint32
	}

	var batteries []wmiBattery
	err := wmi.Query("SELECT EstimatedChargeRemaining, BatteryStatus, EstimatedRunTime FROM Win32_Battery", &batteries)
	if err != nil || len(batteries) == 0 {
		return nil
	}

	b := batteries[0]
	if b.EstimatedChargeRemaining == nil || b.BatteryStatus == nil {
		return nil
	}

	statusMap := map[uint16]string{
		1: "Discharging",
		2: "Unknown",
		3: "Fully Charged",
		4: "Low",
		5: "Critical",
		6: "Charging",
		7: "Charging/High",
		8: "Charging/Low",
		9: "Charging/Critical",
		10: "Undefined",
		11: "Partially Charged",
	}

	statusStr := "Unknown"
	if s, ok := statusMap[*b.BatteryStatus]; ok {
		statusStr = s
	}

	timeLeft := ""
	// EstimatedRunTime is in minutes. 71582788 is the value when charging or unknown runtime
	if b.EstimatedRunTime != nil && *b.EstimatedRunTime > 0 && *b.EstimatedRunTime < 71582788 {
		hours := *b.EstimatedRunTime / 60
		mins := *b.EstimatedRunTime % 60
		if hours > 0 {
			timeLeft = fmt.Sprintf("%dh %dm", hours, mins)
		} else {
			timeLeft = fmt.Sprintf("%dm", mins)
		}
	}

	return &BatteryInfo{
		Percentage: int(*b.EstimatedChargeRemaining),
		Status:     statusStr,
		TimeLeft:   timeLeft,
	}
}
