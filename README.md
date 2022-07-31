# Golang-WDA

[![go doc](https://godoc.org/github.com/electricbubble/gwda?status.svg)](https://pkg.go.dev/github.com/electricbubble/gwda?tab=doc#pkg-index)
[![go report](https://goreportcard.com/badge/github.com/electricbubble/gwda)](https://goreportcard.com/report/github.com/electricbubble/gwda)
[![license](https://img.shields.io/github/license/electricbubble/gwda)](https://github.com/electricbubble/gwda/blob/master/LICENSE)

[appium/WebDriverAgent](https://github.com/appium/WebDriverAgent) Client Library in Golang

> `Android` can use [electricbubble/guia2](https://github.com/electricbubble/guia2)

English | [🇨🇳中文](README_CN.md)

## Installation

> First, install WebDriverAgent for iOS devices

```shell script
go get github.com/electricbubble/gwda
```

## QuickStart

#### [Connection Device](examples/connect/main.go)

```go
package main

import (
	"github.com/electricbubble/gwda"
	"log"
)

func main() {
	// var urlPrefix = "http://localhost:8100"
	// The function may also require 'iproxy 8100 8100' to forward the device port first
	// driver, _ := gwda.NewDriver(nil, urlPrefix)

	// Connect devices via USB
	driver, _ := gwda.NewUSBDriver(nil)

	log.Println(driver.IsWdaHealthy())
}

```

#### [Touch](examples/touch/main.go)

```go
package main

import (
	"github.com/electricbubble/gwda"
)

func main() {
	driver, _ := gwda.NewUSBDriver(nil)

	x, y := 50, 256

	driver.Tap(x, y)

	driver.DoubleTap(x, y)

	driver.TouchAndHold(x, y)

	fromX, fromY, toX, toY := 50, 256, 100, 256

	driver.Drag(fromX, fromY, toX, toY)

	driver.Swipe(fromX, fromY, toX, toY)

	// requires hardware support: 3D Touch 
	// driver.ForceTouch(x, y, 0.8)
}

```

> [AssistiveTouch](examples/touch/main.go) `driver.PerformW3CActions` `driver.PerformAppiumTouchActions`

#### [App Actions](examples/app/main.go)

```go
package main

import (
	"github.com/electricbubble/gwda"
)

func main() {
	driver, _ := gwda.NewUSBDriver(nil)

	var bundleId = "com.apple.Preferences"

	driver.AppLaunchUnattached(bundleId)

	driver.AppDeactivate(2)

	driver.AppTerminate(bundleId)

	driver.AppActivate(bundleId)

	// Resets the 📷 camera authorization status of the current application
	// driver.AppAuthReset(gwda.ProtectedResourceCamera)
}

```

#### [Keyboard](examples/keyboard/main.go)

```go
package main

import (
	"github.com/electricbubble/gwda"
)

func main() {
	driver, _ := gwda.NewUSBDriver(nil)

	driver.SendKeys("hello")
}

```

> [specified element](examples/keyboard/main.go) `element.SendKeys`

#### [Siri](examples/siri/main.go)

```go
package main

import (
	"github.com/electricbubble/gwda"
)

func main() {
	driver, _ := gwda.NewUSBDriver(nil)

	driver.SiriActivate("What's the weather like today")
}

```

#### [Alert](examples/alert/main.go)

```go
package main

import (
	"github.com/electricbubble/gwda"
	"log"
)

func main() {
	driver, _ := gwda.NewUSBDriver(nil)

	text, _ := driver.AlertText()
	log.Println(text)

	alertButtons, _ := driver.AlertButtons()
	log.Println(alertButtons)

	driver.AlertAccept()
	// driver.AlertDismiss()

	// driver.SendKeys("ah")
}

```

#### [Device information](examples/info/main.go)

```go
package main

import (
	"github.com/electricbubble/gwda"
	"log"
)

func main() {
	driver, _ := gwda.NewUSBDriver(nil)

	deviceInfo, _ := driver.DeviceInfo()
	log.Println(deviceInfo.Name)

	batteryInfo, _ := driver.BatteryInfo()
	log.Println(batteryInfo.State)

	windowSize, _ := driver.WindowSize()
	log.Println(windowSize)

	location, err := driver.Location()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(location)

	// screen, _ := driver.Screen()
	// log.Println(screen)
}

```

#### [Hardware button](examples/button/main.go)

```go
package main

import (
	"github.com/electricbubble/gwda"
)

func main() {
	driver, _ := gwda.NewUSBDriver(nil)

	// driver.Homescreen()

	driver.PressButton(gwda.DeviceButtonHome)
	driver.PressButton(gwda.DeviceButtonVolumeUp)
	driver.PressButton(gwda.DeviceButtonVolumeDown)
}

```

#### [Screenshot](examples/screenshot/main.go)

```go
package main

import (
	"github.com/electricbubble/gwda"
	"image"
)

func main() {
	driver, _ := gwda.NewUSBDriver(nil)

	screenshot, _ := driver.Screenshot()

	img, format, _ := image.Decode(screenshot)
	_, _ = img, format
}

```

#### [Debug](examples/debug/main.go)

```go
package main

import (
	"fmt"
	"github.com/electricbubble/gwda"
)

func main() {
	driver, _ := gwda.NewUSBDriver(nil)

	source, _ := driver.Source()
	fmt.Println(source)

	// fmt.Println(driver.AccessibleSource())

	// gwda.SetDebug(true)
}

```

## Extensions

| |About|
|---|---|
|[electricbubble/gwda-ext-opencv](https://github.com/electricbubble/gwda-ext-opencv)|Operate with pictures|

## Alternatives

| |About|
|---|---|
|[openatx/facebook-wda](https://github.com/openatx/facebook-wda)|Facebook WebDriverAgent Python Client Library (not official)|

## Thanks

Thank you [JetBrains](https://www.jetbrains.com/?from=gwda) for providing free open source licenses
