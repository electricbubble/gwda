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
	"fmt"
	"github.com/electricbubble/gwda"
	"log"
	"time"
)

func main() {
	var deviceURL = "http://localhost:8100"
	// gwda.Debug = true
	// 连接设备
	c, err := gwda.NewClient(deviceURL)
	if err != nil {
		log.Fatalln(err)
	}

	// 获取当前屏幕锁定状态
	locked, err := c.IsLocked()
	if err != nil {
		log.Fatalln(err)
	}
	if locked {
		// 触发解锁
		err = c.Unlock()
		if err != nil {
			log.Fatalln("解锁失败", err)
		}
	}

	cAppInfo, err := c.ActiveAppInfo()
	if err != nil {
		log.Fatalln("查看当前 APP 信息 失败", err)
	}
	fmt.Printf("当前 App 的 PID: %d\n当前 App 的 bundleId: %s\n", cAppInfo.Pid, cAppInfo.BundleID)

	bundleId := "com.apple.Preferences"
	// 创建 session，可选输入 bundle id，如指定，则启动并且等待 app 加载完毕
	// 弹窗自动处理
	// gwda.NewWDASessionCapability(bundleId).SetDefaultAlertAction(gwda.WDASessionAlertActionAccept)
	// gwda.NewWDASessionCapability(bundleId).SetDefaultAlertAction(gwda.WDASessionAlertActionDismiss)
	s, err := c.NewSession(gwda.NewWDASessionCapability(bundleId))
	if err != nil {
		log.Fatalln("创建 session 失败", err)
	}
	defer func() {
		time.Sleep(time.Second * 10)
		// 如果使用了弹窗自动处理，一定要执行删除，保证 弹窗监控 被禁用，避免 wda 内部错误
		s.DeleteSession()
	}()

	btyInfo, err := s.BatteryInfo()
	if err != nil {
		log.Fatalln("电量获取失败", err)
	}
	fmt.Printf("当前电量: %.0f%%\n", btyInfo.Level*100)
	switch btyInfo.State {
	case gwda.WDABatteryUnplugged:
		fmt.Println("未充电")
	case gwda.WDABatteryCharging:
		fmt.Println("充电中，电量少于 100%")
	case gwda.WDABatteryFull:
		fmt.Println("充电中，并且电量已满 100%")
	}

	wSize, err := s.WindowSize()
	if err != nil {
		log.Fatalln("获取 逻辑分辨率 失败", err)
	}
	fmt.Printf("逻辑分辨率：\t宽: %d\t高: %d\n", wSize.Width, wSize.Height)

	statusBarSize, err := s.StatusBarSize()
	if err != nil {
		log.Fatalln("获取 状态栏大小 失败", err)
	}
	fmt.Printf("状态栏(逻辑宽高) 宽: %d\t高: %d\n", statusBarSize.Width, statusBarSize.Height)
	scale, err := s.Scale()
	if err != nil {
		log.Fatalln("获取 缩放倍率 失败", err)
	}
	fmt.Printf("屏幕缩放倍率: %.2f\n", scale)
	// screenInfo, err := s.Screen()
	// if err != nil {
	// 	log.Fatalln("获取 状态栏大小 和 缩放倍率 失败", err)
	// }
	fmt.Printf("渲染后的屏幕分辨率 宽: %.2f\t高: %.2f\n", float32(wSize.Width)*scale, float32(wSize.Height)*scale)

	appRunState, err := s.AppState(bundleId)
	if err != nil {
		log.Fatalln("获取 App 运行状态 失败", err)
	}
	switch appRunState {
	case gwda.WDAAppRunningBack:
		fmt.Println("该 App 正后台运行中，开始切换到前台运行")
		_ = s.AppActivate(bundleId)
	case gwda.WDAAppNotRunning:
		fmt.Println("该 App 未运行，开始打开 App")
		_ = s.AppLaunch(bundleId)
	}

	element, err := s.FindElement(gwda.WDALocator{LinkText: gwda.NewWDAElementAttribute().SetLabel("通用")})
	if err != nil {
		log.Fatalln("查找元素失败", err)
	}

	err = element.Click()
	if err != nil {
		log.Fatalln("元素点击失败", err)
	}

	// gwda.Debug = false
	// source, _ := s.Source()
	// gwda.Debug = true

	elementAbout, err := s.FindElement(gwda.WDALocator{PartialLinkText: gwda.NewWDAElementAttribute().SetValue("关于本")})
	if err != nil {
		log.Fatalln("查找元素失败", err)
	}
	// fmt.Println(elementAbout.Rect())
	rect, err := elementAbout.Rect()
	if err != nil {
		log.Fatalln("元素 rect 获取失败", err)
	}
	fmt.Printf("该元素的坐标 x,y: (%d, %d)\t宽: %d, 高: %d\n", rect.X, rect.Y, rect.Width, rect.Height)

	// 点击指定坐标
	err = s.Tap(rect.X+1, rect.Y+1)
	if err != nil {
		log.Fatalln("点击失败", err)
	}

	// 可以当作字符串输出从 wda 收到的响应值
	// fmt.Println(rect)
	// {
	//    "y" : 99,
	//    "x" : 0,
	//    "width" : 375,
	//    "height" : 44
	//  }

	// fmt.Println(source)

}

```
> 以上代码仅使用了 iPhone X (13.4.1) 和 iPhone 6s (11.4.1) 进行了测试。

## TODO

暂时只模仿了 [facebook-wda](https://github.com/openatx/facebook-wda) 的部分功能

需要做的还有很多很多，列不出来🤕
