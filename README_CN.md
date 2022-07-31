# Golang-WDA
[![go doc](https://godoc.org/github.com/electricbubble/gwda?status.svg)](https://pkg.go.dev/github.com/electricbubble/gwda?tab=doc#pkg-index)
[![go report](https://goreportcard.com/badge/github.com/electricbubble/gwda)](https://goreportcard.com/report/github.com/electricbubble/gwda)
[![license](https://img.shields.io/github/license/electricbubble/gwda)](https://github.com/electricbubble/gwda/blob/master/LICENSE)

使用 Golang 实现 [appium/WebDriverAgent](https://github.com/appium/WebDriverAgent) 的客户端库

参考 facebook-wda (python): [https://github.com/openatx/facebook-wda](https://github.com/openatx/facebook-wda)

## 扩展库

- [electricbubble/gwda-ext-opencv](https://github.com/electricbubble/gwda-ext-opencv) 直接通过指定图片进行操作

> 如果使用 `Android` 设备, 可查看 [electricbubble/guia2](https://github.com/electricbubble/guia2)

## 安装

> 必须先安装好 `WDA`，安装步骤可参考 [ATX 文档 - iOS 真机如何安装 WebDriverAgent](https://testerhome.com/topics/7220) 或者
> [WebDriverAgent 安装](http://leixipaopao.com/posts/0005-wda-appium-installing/)

```shell script
go get github.com/electricbubble/gwda
```

## 快速上手

#### [连接设备](examples/connect/main.go)

```go
package main

import (
	"github.com/electricbubble/gwda"
	"log"
)

func main() {
	// var urlPrefix = "http://localhost:8100"
	// 该函数或许还需要 `iproxy 8100 8100` 先进行设备端口转发
	// driver, _ := gwda.NewDriver(nil, urlPrefix)

	// 通过 USB 直连设备
	driver, _ := gwda.NewUSBDriver(nil)

	log.Println(driver.IsWdaHealthy())
}

```

#### [手势操作](examples/touch/main.go)

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

	// 需要 3D Touch 硬件支持
	// driver.ForceTouch(x, y, 0.8)
}

```

> [自定义手势](examples/touch/main.go) `driver.PerformW3CActions` `driver.PerformAppiumTouchActions`

#### [App 操作](examples/app/main.go)

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

	// 重置当前 App 的 相机📷 权限
	// driver.AppAuthReset(gwda.ProtectedResourceCamera)
}

```

#### [键盘输入](examples/keyboard/main.go)

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

> [指定元素的输入](examples/keyboard/main.go) `element.SendKeys`


#### [Siri 操作](examples/siri/main.go)

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

#### [弹窗操作](examples/alert/main.go)

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

#### [基本设备信息](examples/info/main.go)

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

#### [按键操作](examples/button/main.go)

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

#### [截图](examples/screenshot/main.go)

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

#### [调试函数](examples/debug/main.go)

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

## Thanks

Thank you [JetBrains](https://www.jetbrains.com/?from=gwda) for providing free open source licenses
