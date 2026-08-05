package collector

import (
	"github.com/StackExchange/wmi"
)

const unavailable = "<unavailable>"

func CollectBIOS() BIOSInfo {
	type wmiBIOS struct {
		Manufacturer string
		Version      string
		ReleaseDate  string
		SerialNumber string
	}

	var biosList []wmiBIOS
	err := wmi.Query("SELECT Manufacturer, Version, ReleaseDate, SerialNumber FROM Win32_BIOS", &biosList)
	if err != nil || len(biosList) == 0 {
		return BIOSInfo{
			Vendor:  unavailable,
			Version: unavailable,
		}
	}

	b := biosList[0]
	dateStr := ""
	if len(b.ReleaseDate) >= 8 {
		dateStr = b.ReleaseDate[0:4] + "-" + b.ReleaseDate[4:6] + "-" + b.ReleaseDate[6:8]
	}

	type wmiBoard struct {
		Manufacturer string
		Product      string
		SerialNumber string
	}
	var boards []wmiBoard
	wmi.Query("SELECT Manufacturer, Product, SerialNumber FROM Win32_BaseBoard", &boards)

	var boardMfg, boardModel, boardSerial string
	if len(boards) > 0 {
		boardMfg = boards[0].Manufacturer
		boardModel = boards[0].Product
		boardSerial = boards[0].SerialNumber
	}

	return BIOSInfo{
		Vendor:       b.Manufacturer,
		Version:      b.Version,
		Date:         dateStr,
		Manufacturer: boardMfg,
		Model:        boardModel,
		SerialNumber: boardSerial,
	}
}
