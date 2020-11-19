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

	err = driver.SendKeys("hello")
	if err != nil {
		log.Fatalln(err)
	}

	// element, err := driver.ActiveElement()
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// err = element.SendKeys("little monkey", 5)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
}
