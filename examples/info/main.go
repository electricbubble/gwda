package main

import (
	"github.com/electricbubble/gwda"
	"log"
)

func main() {
	driver, err := gwda.NewUSBDriver(nil)
	if err != nil {
		log.Fatalln(err)
	}

	deviceInfo, err := driver.DeviceInfo()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(deviceInfo.Name)

	batteryInfo, err := driver.BatteryInfo()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(batteryInfo.State)

	windowSize, err := driver.WindowSize()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(windowSize)

	// screen, err := driver.Screen()
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// log.Println(screen)
}
