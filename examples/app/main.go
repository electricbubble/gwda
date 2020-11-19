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

	var bundleId = "com.apple.Preferences"

	err = driver.AppLaunchUnattached(bundleId)
	if err != nil {
		log.Fatalln(err)
	}

	err = driver.AppDeactivate(2)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = driver.AppTerminate(bundleId)
	if err != nil {
		log.Fatalln(err)
	}

	err = driver.AppActivate(bundleId)
	if err != nil {
		log.Fatalln(err)
	}

	// 重置当前 App 的 相机📷 权限
	// err = driver.AppAuthReset(gwda.ProtectedResourceCamera)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
}
