package main

import (
	"log"

	"github.com/electricbubble/gwda"
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

	// send keys with specified frequency
	err = driver.SendKeys("world", gwda.WithFrequency(30))
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
