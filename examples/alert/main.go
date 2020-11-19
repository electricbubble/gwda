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

	text, err := driver.AlertText()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(text)

	alertButtons, err := driver.AlertButtons()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(alertButtons)

	err = driver.AlertAccept()
	// err = driver.AlertDismiss()
	if err != nil {
		log.Fatalln(err)
	}

	// driver.SendKeys("ah")
}
