package gwda

import (
	"fmt"

	giDevice "github.com/electricbubble/gidevice"
)

const (
	defaultPort      = 8100
	defaultMjpegPort = 9100
)

type Device struct {
	deviceID     int
	serialNumber string
	Port         int
	MjpegPort    int

	d giDevice.Device
}

type DeviceOption func(d *Device)

func WithSerialNumber(serialNumber string) DeviceOption {
	return func(d *Device) {
		d.serialNumber = serialNumber
	}
}

func WithPort(port int) DeviceOption {
	return func(d *Device) {
		d.Port = port
	}
}

func WithMjpegPort(port int) DeviceOption {
	return func(d *Device) {
		d.MjpegPort = port
	}
}

func NewDevice(options ...DeviceOption) (device *Device, err error) {
	var usbmux giDevice.Usbmux
	if usbmux, err = giDevice.NewUsbmux(); err != nil {
		return nil, fmt.Errorf("init usbmux failed: %v", err)
	}

	var deviceList []giDevice.Device
	if deviceList, err = usbmux.Devices(); err != nil {
		return nil, fmt.Errorf("get attached devices failed: %v", err)
	}

	device = &Device{
		Port:      defaultPort,
		MjpegPort: defaultMjpegPort,
	}
	for _, option := range options {
		option(device)
	}

	serialNumber := device.serialNumber
	for _, d := range deviceList {
		// find device by serial number if specified
		if serialNumber != "" && d.Properties().SerialNumber != serialNumber {
			continue
		}

		device.deviceID = d.Properties().DeviceID
		device.serialNumber = d.Properties().SerialNumber
		device.d = d
		return device, nil
	}

	return nil, fmt.Errorf("device %s not found", device.serialNumber)
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
		devices[i].Port = defaultPort
		devices[i].MjpegPort = defaultMjpegPort
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
