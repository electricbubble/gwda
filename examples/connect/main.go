package main

import (
	"github.com/electricbubble/gwda"
	"log"
)

func main() {
	// var urlPrefix = "http://localhost:8100"
	// 该函数或许还需要 `iproxy 8100 8100` 先进行设备端口转发
	// driver, err := gwda.NewDriver(nil, urlPrefix)

	// 通过 USB 直连设备
	driver, err := gwda.NewUSBDriver(nil)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(driver.IsWdaHealthy())
}
