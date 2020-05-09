package main

import (
	"github.com/electricbubble/gwda"
	"log"
	"os"
	"path/filepath"
	"time"
)

func main() {
	// TODO 预设自动弹窗 `Selector` 的值
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

	time.Sleep(time.Second * 1)

	userHomeDir, _ := os.UserHomeDir()
	err = client.ScreenshotToDisk(filepath.Join(userHomeDir, "Desktop", "homescreen.png"))
	checkErr("截图并保存", err)

	deviceInfo, err := client.DeviceInfo()
	checkErr("获取设备信息", err)
	log.Println("Name:", deviceInfo.Name)
	log.Println("IsSimulator:", deviceInfo.IsSimulator)

	// TODO 自动弹窗、App 启动参数等设置
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
		// TODO app 启动参数设置
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

		log.Println("重新启动 App:", bundleId)
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

	err = session.SendKeys("音乐\n", 1)
	checkErr("通过 session 输入文本", err)

	err = elemSearch.Clear()
	checkErr("清空 搜索输入框", err)

	err = elemSearch.SendKeys("辅助功能-" + gwda.WDATextBackspaceSequence + "\n")
	checkErr("输入文本", err)

	imgSearch, format, err := elemSearch.ScreenshotToImage()
	checkErr("截图元素并保存为 image.Image", err)
	log.Println("搜索输入框 的截图图片格式:", format)
	log.Println("搜索输入框 的截图图片大小(像素):", imgSearch.Bounds().Size())

	elemSearchRet, err := session.FindElement(gwda.WDALocator{Predicate: "type in {'XCUIElementTypeTable', 'XCUIElementTypeCollectionView'} && visible == true"})
	checkErr("找到 搜索结果列表框", err)

	cellElemRets, err := elemSearchRet.FindVisibleCells()
	checkErr("找到全部 搜索结果", err)
	log.Printf("共找到 %d 个搜索结果\n", len(cellElemRets))

	elemCancel, err := session.FindElement(gwda.WDALocator{Predicate: "type == 'XCUIElementTypeButton' && name == '取消'"})
	checkErr("找到 取消 按钮", err)

	err = elemCancel.Click()
	checkErr("点击 取消 按钮", err)

	err = session.PressVolumeUpButton()
	checkErr("触发设备按键🔊音量⬆️", err)

	time.Sleep(time.Millisecond * 500)

	err = session.PressHomeButton()
	checkErr("触发设备按键 Home️", err)

	time.Sleep(time.Millisecond * 500)

	err = session.PressVolumeDownButton()
	checkErr("触发设备按键🔊音量⬇️", err)

	time.Sleep(time.Millisecond * 1500)
	err = session.SwipeLeft()
	checkErr("向左👈滑动", err)
	time.Sleep(time.Millisecond * 350)

	elemIcon, err := session.FindElement(gwda.WDALocator{ClassChain: "**/XCUIElementTypeIcon[`visible == true`]"})
	checkErr("找到 当前屏幕的第一个 App/文件夹", err)

	text, err := elemIcon.Text()
	checkErr("获取 当前屏幕第一个 App/文件夹 的文本内容", err)
	log.Println("当前屏幕第一个 App/文件夹 的文本内容:", text)

	rectIcon, err := elemIcon.Rect()
	checkErr("获取该 App/文件夹 的坐标和大小", err)
	log.Println("该 App/文件夹 的坐标和大小:", rectIcon)

	err = elemIcon.TouchAndHold(3)
	checkErr("按住并保持指定秒数 (默认1s)", err)

	time.Sleep(time.Millisecond * 150)
	err = session.PressHomeButton()
	checkErr("触发设备按键 Home️", err)
	time.Sleep(time.Millisecond * 150)

	err = session.ForceTouch(rectIcon.X+rectIcon.Width/2, rectIcon.Y+rectIcon.Height/2, 1, 0.5)
	checkErr("指定压力值, 触发 3D Touch, (默认保持 1s)", err)

	time.Sleep(time.Second * 3)
	err = session.PressHomeButton()
	checkErr("触发设备按键 Home️", err)
	time.Sleep(time.Millisecond * 150)

	orientation, err := session.Orientation()
	checkErr("获取当前设备方向", err)
	rotation, err := session.Rotation()
	checkErr("获取当前设备 Rotation", err)
	log.Println("Orientation:", orientation)
	log.Println("Rotation:", rotation)

	bundleId = "com.apple.calculator"
	err = session.AppLaunch(bundleId)
	checkErr("启动 App 计算器", err)

	switch orientation {
	case gwda.WDAOrientationPortrait:
		err = session.SetOrientation(gwda.WDAOrientationLandscapeLeft)
	default:
		err = session.SetRotation(gwda.WDARotation{X: 0, Y: 0, Z: 0})
	}
	checkErr("修改设备方向", err)

	err = session.SiriActivate("当前时间")
	checkErr("激活 Siri", err)
}

func checkErr(msg string, err error) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
