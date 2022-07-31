package gwda

import (
	"fmt"

	giDevice "github.com/electricbubble/gidevice"
)

type Device struct {
	deviceID     int
	serialNumber string
	Port         int
	MjpegPort    int

	d giDevice.Device
}

func DeviceList() (devices []Device, err error) {
	var usbmux giDevice.Usbmux
	if usbmux, err = giDevice.NewUsbmux(); err != nil {
		return nil, fmt.Errorf("usbmuxd: %w", err)
	}

	var deviceList []giDevice.Device
	if deviceList, err = usbmux.Devices(); err != nil {
		return nil, fmt.Errorf("device list: %w", err)
	}

	devices = make([]Device, len(deviceList))

	for i := range devices {
		devices[i].deviceID = deviceList[i].Properties().DeviceID
		devices[i].serialNumber = deviceList[i].Properties().SerialNumber
		devices[i].Port = 8100
		devices[i].MjpegPort = 9100
		devices[i].d = deviceList[i]
	}

	return
}

func (d Device) DeviceID() int {
	return d.deviceID
}

func (d Device) SerialNumber() string {
	return d.serialNumber
}

func (d Device) GIDevice() giDevice.Device {
	return d.d
}
