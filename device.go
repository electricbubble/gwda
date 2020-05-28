package gwda

import goUSBMux "github.com/electricbubble/go-usbmuxd-device"

type Device struct {
	deviceID                         int
	serialNumber                     string
	WDAPort                          int
	MjpegPort                        int
	IsInitializesAlertButtonSelector bool
}

func DeviceList() (devices []Device, err error) {
	var deviceList []goUSBMux.USBDevice
	if deviceList, err = goUSBMux.NewUSBHub().DeviceList(); err != nil {
		return nil, err
	}
	devices = make([]Device, len(deviceList))

	for i := range devices {
		devices[i].deviceID = deviceList[i].DeviceID
		devices[i].serialNumber = deviceList[i].SerialNumber
		devices[i].WDAPort = 8100
		devices[i].MjpegPort = 9100
		devices[i].IsInitializesAlertButtonSelector = true
	}

	return
}

func (d Device) DeviceID() int {
	return d.deviceID
}

func (d Device) SerialNumber() string {
	return d.serialNumber
}
