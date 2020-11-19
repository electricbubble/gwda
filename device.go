package gwda

import (
	"fmt"
	goUSBMux "github.com/electricbubble/go-usbmuxd-device"
)

type Device struct {
	deviceID     int
	serialNumber string
	Port         int
	MjpegPort    int
}

func DeviceList() (devices []Device, err error) {
	var deviceList []goUSBMux.USBDevice
	if deviceList, err = goUSBMux.NewUSBHub().DeviceList(); err != nil {
		return nil, fmt.Errorf("device list: %w", err)
	}
	devices = make([]Device, len(deviceList))

	for i := range devices {
		devices[i].deviceID = deviceList[i].DeviceID
		devices[i].serialNumber = deviceList[i].SerialNumber
		devices[i].Port = 8100
		devices[i].MjpegPort = 9100
	}

	return
}

func (d Device) DeviceID() int {
	return d.deviceID
}

func (d Device) SerialNumber() string {
	return d.serialNumber
}
