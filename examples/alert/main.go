package main

import (
	"github.com/electricbubble/gwda"
	"log"
	"time"
)

func alertsMonitor() {
	// 这里为了测试弹窗，在 设置 里还原了 "还原位置与隐私" 重新打开 app 就会重新弹窗定位或者通知的弹窗
	// bundleId := "com.apple.weather"
	bundleId := "com.taobao.taobao4iphone"

	// 	⬇️ 弹窗自动处理 ⬇️

	// 在连接设备的时候，追加一个 "true" 参数，用于在连接设备时，预设“允许/拒绝”的按钮选择器
	// 允许: **/XCUIElementTypeButton[`label IN {'允许','好','仅在使用应用期间','暂不'}`]
	// 拒绝: **/XCUIElementTypeButton[`label IN {'不允许','暂不'}`]
	client, err := gwda.NewClient("http://localhost:8100", true)
	checkErr("连接设备", err)

	// 也可以选择不在连接设备时去预设选择器（默认值只针对了中文）
	// client, err := gwda.NewClient("http://localhost:8100")
	// checkErr("连接设备", err)
	// 通过下面两个函数也可以在创建 session 之前设置好想要设置的弹窗按钮选择器
	// client.SetAcceptAlertButtonSelector("**/XCUIElementTypeButton[`label IN {'允许','好','仅在使用应用期间','暂不'}`]")
	// client.SetDismissAlertButtonSelector("**/XCUIElementTypeButton[`label IN {'不允许','暂不'}`]")
	// ⚠️ 必须使用 `ClassChain` 来定位弹窗按钮选择器

	// 创建 session 时，设置当 Alert 出现时的默认处理行为（ Accept/Dismiss ）
	// gwda.WDASessionAlertActionAccept
	session, err := client.NewSession(
		// 可在创建 session 指定需要打开的 App （创建后会自动打开，并默认 SetShouldWaitForQuiescence(true)）
		// gwda.NewWDASessionCapability(bundleId).
		gwda.NewWDASessionCapability().
			SetDefaultAlertAction(gwda.WDASessionAlertActionAccept))
	checkErr("创建 session", err)

	// 将开始 2s 一次的弹窗监控

	// _ = session.AppLaunch(bundleId, gwda.NewWDAAppLaunchOption().SetShouldWaitForQuiescence(false))
	// SetShouldWaitForQuiescence(true)
	_ = session.AppLaunch(bundleId)

	time.Sleep(time.Second * 4)

	// 实际上 弹窗监控是可以跨 session 的，所以，在没有 DeleteSession() 之前可以自动处理其他的 App 弹窗
	bundleId = "com.apple.AppStore"
	_ = session.AppLaunch(bundleId, gwda.NewWDAAppLaunchOption().SetShouldWaitForQuiescence(false))
	// _ = session.AppLaunch(bundleId)

	defer func() {
		time.Sleep(time.Second * 12)
		_ = session.AppTerminate(bundleId)
		// ⚠️ 当设置自动弹窗处理的时候，请务必使用 session.DeleteSession() 来让 `WDA` 内部去关闭弹窗的监控
		// 删除 session 同时会关闭 `gwda.NewWDASessionCapability(bundleId)` 指定的 App
		_ = session.DeleteSession()
	}()
}
func main() {
	// 	⬇️ 弹窗自动处理 ⬇️
	// alertsMonitor()
	// return

	client, err := gwda.NewClient("http://localhost:8100", true)
	checkErr("连接设备", err)
	session, err := client.NewSession()
	checkErr("创建 session", err)

	// 	⬇️ 手动处理使用以下相关函数 ⬇️

	// 获取弹窗的内容
	alertText, err := session.AlertText()
	checkErr("弹窗内容", err)
	log.Println(alertText)

	// 获取弹窗的全部按钮文本
	alertButtons, err := session.AlertButtons()
	checkErr("当前弹窗的全部按钮文本", err)
	log.Println(alertButtons)

	// 点击可指定“名称”的“yes”按钮
	// 不指定则默认使用 预设/自定义的选择器
	err = session.AlertAccept()
	checkErr("accept", err)

	// 点击可指定“名称”的“no”按钮
	// 不指定则默认使用 预设/自定义的选择器
	err = session.AlertDismiss()
	checkErr("dismiss", err)

	// 在弹窗里的输入框输入内容
	err = session.AlertSendKeys("text")
	checkErr("弹窗内输入指定文本", err)
}

func checkErr(msg string, err error) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
