# Golang-wda

使用 Golang 实现 [appium/WebDriverAgent](https://github.com/appium/WebDriverAgent) 的客户端库

参考 [facebook-wda](https://github.com/openatx/facebook-wda)

## 安装

> 必须先安装好 `WDA`，安装步骤可参考 [ATX 文档 - iOS 真机如何安装 WebDriverAgent](https://testerhome.com/topics/7220) 或者
> [WebDriverAgent 安装](http://leixipaopao.com/posts/0005-wda-appium-installing/)

```shell script
go get -u github.com/electricbubble/gwda
```

## 使用
```go
package main

import (
	"github.com/electricbubble/gwda"
	"log"
	"os"
	"path/filepath"
	"time"
)

func main() {
	client, err := gwda.NewClient("http://localhost:8100")
	checkErr("连接设备", err)

	err = client.Lock()
	checkErr("触发锁屏", err)

	isLocked, err := client.IsLocked()
	checkErr("判断是否处于屏幕锁定状态", err)

	if isLocked {
		err = client.Unlock()
		checkErr("触发解锁", err)
	}

	err = client.Homescreen()
	checkErr("切换到主屏幕", err)

	userHomeDir, _ := os.UserHomeDir()
	err = client.ScreenshotToDisk(filepath.Join(userHomeDir, "Desktop", "homescreen.png"))
	checkErr("截图并保存", err)

	deviceInfo, err := client.DeviceInfo()
	checkErr("获取设备信息", err)
	log.Println("Name:", deviceInfo.Name)
	log.Println("IsSimulator:", deviceInfo.IsSimulator)

	session, err := client.NewSession()
	checkErr("创建 session", err)

	// defer session.DeleteSession()

	windowSize, err := session.WindowSize()
	checkErr("获取当前应用的大小", err)
	// 实际获取的是当前 App 的大小，但当前 App 是 主屏幕 时，通常得到的就是当前设备的屏幕大小
	log.Println("UIKit Size (Points):", windowSize.Width, "x", windowSize.Height)

	scale, err := session.Scale()
	checkErr("获取 UIKit Scale factor", err)
	log.Println("UIKit Scale factor:", scale)
	log.Println("Native Resolution (Pixels):", float64(windowSize.Width)*scale, "x", float64(windowSize.Height)*scale)

	statusBarSize, err := session.StatusBarSize()
	checkErr("获取 status bar 的大小", err)
	log.Println("Status bar size:", statusBarSize.Width, "x", statusBarSize.Height)

	batteryInfo, err := session.BatteryInfo()
	checkErr("获取🔋电量信息", err)
	switch batteryInfo.State {
	case gwda.WDABatteryUnplugged:
		log.Println("State:", batteryInfo.State)
	case gwda.WDABatteryCharging:
		if batteryInfo.Level == 1 {
			log.Println("State:", gwda.WDABatteryFull)
		} else {
			log.Println("State:", batteryInfo.State)
		}
	case gwda.WDABatteryFull:
		log.Println("State:", batteryInfo.State)
	}
	log.Printf("Level: %.00f%%\n", batteryInfo.Level*100)

	bundleId := "com.apple.Preferences"

	appRunState, err := session.AppState(bundleId)
	checkErr("获取指定 App 的运行状态", err)
	switch appRunState {
	case gwda.WDAAppNotRunning:
		log.Println("该 App 未运行, 开始打开 App:", bundleId)
		err = session.AppLaunch(bundleId)
		checkErr("启动指定 App", err)
	case gwda.WDAAppRunningBack:
		log.Println("该 App 正后台运行中, 开始切换到前台运行:", bundleId)
		err = session.AppActivate(bundleId)
		checkErr("切换指定 App 到前台运行", err)
	case gwda.WDAAppRunningFront:
		log.Println("该 App 正前台运行中, 开始关闭 App:", bundleId)
		err = session.AppTerminate(bundleId)
		checkErr("关闭指定 App", err)

		log.Println("重新打开 App:", bundleId)
		err = session.AppLaunch(bundleId)
		checkErr("再启动指定 App", err)
	}

	log.Println("使当前 App 退回 主屏幕, 并至少等待 3s 后(默认等待时间)再切换到前台")
	err = session.AppDeactivate()
	checkErr("使当前 App 退回 主屏幕, 并至少等待 3s 后(默认等待时间)再切换到前台", err)

	activeAppInfo, err := session.ActiveAppInfo()
	checkErr("获取当前 App 的信息", err)
	log.Println("当前 App 的 PID:", activeAppInfo.Pid)

	err = session.SwipeUp()
	checkErr("向上👆滑动", err)

	err = session.Tap(20, 1)
	checkErr("点击指定坐标点", err)

	time.Sleep(time.Second * 1)

	elemSearch, err := session.FindElement(gwda.WDALocator{ClassName: gwda.WDAElementType{SearchField: true}})
	checkErr("找到 搜索输入框", err)

	err = elemSearch.Click()
	checkErr("点击 搜索输入框", err)

	err = session.SendKeys("辅助功能\n", 1)
	checkErr("通过 session 输入文本", err)

	err = elemSearch.Clear()
	checkErr("清空 搜索输入框", err)

	err = elemSearch.SendKeys("音乐-" + gwda.WDATextBackspaceSequence + "\n")
	checkErr("输入文本", err)

	imgSearch, format, err := elemSearch.ScreenshotToImage()
	checkErr("截图元素并保存为 image.Image", err)
	log.Println("该元素图片的格式:", format)
	log.Println("该元素图片的大小(像素):", imgSearch.Bounds().Size())

	elemCancel, err := session.FindElement(gwda.WDALocator{Predicate: "type == 'XCUIElementTypeButton' && name == '取消'"})
	checkErr("找到 取消 按钮", err)

	rectCancel, err := elemCancel.Rect()
	checkErr("获取 取消 按钮的坐标和大小", err)
	log.Println(rectCancel)

	err = elemCancel.Click()
	checkErr("点击 取消 按钮", err)
}

func checkErr(msg string, err error) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}

```
> 以上代码仅使用了 iPhone X (13.4.1) 和 iPhone 6s (11.4.1) 进行了测试。

## TODO

待补充更多 Examples
