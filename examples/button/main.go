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

	// err = driver.Homescreen()
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	err = driver.PressButton(gwda.DeviceButtonHome)
	// err = driver.PressButton(gwda.DeviceButtonVolumeUp)
	// err = driver.PressButton(gwda.DeviceButtonVolumeDown)
	if err != nil {
		log.Fatalln(err)
	}
}
