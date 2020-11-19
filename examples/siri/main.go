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

	err = driver.SiriActivate("What's the weather like today")
	if err != nil {
		log.Fatalln(err)
	}

	// It doesn't actually work
	// driver.SiriOpenUrl("Prefs:root=Bluetooth")
}
