package gwda

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSession_GetActiveSession(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	Debug = true
	s, err := c.NewSession()
	checkErr(t, err)
	sessionInfo, err := s.GetActiveSession()
	checkErr(t, err)
	if sessionInfo.SessionID == "" {
		t.Fatal("session id should not be empty")
	}

	t.Logf("\n当前 App 的 名称:\t\t%s\n"+
		"当前 App 的 BundleId:\t%s\n"+
		"当前设备的系统版本:\t\t%s\n",
		sessionInfo.Capabilities.BrowserName, sessionInfo.Capabilities.CFBundleIdentifier, sessionInfo.Capabilities.SdkVersion)
	// t.Log("当前 App 的 BundleId:", sessionInfo.Capabilities.CFBundleIdentifier)
}

func TestSession_DeleteSession(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	Debug = true
	s, err := c.NewSession()
	checkErr(t, err)
	err = s.DeleteSession()
	checkErr(t, err)

	// t.Log(c.ActiveAppInfo())
	err = s.SiriActivate("打开 微信")
	if err == nil {
		t.Fatal("It should not be nil")
	}
	if !strings.EqualFold(err.Error(), "invalid session id: Session does not exist") {
		t.Fatal(err)
	}

	// t.Log(s.SiriOpenURL("weixin://"))
}

func TestSession_DeviceInfo(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	Debug = true
	s, err := c.NewSession()
	checkErr(t, err)
	dInfo, err := s.DeviceInfo()
	checkErr(t, err)
	t.Log(dInfo)
	t.Log(dInfo.Name)
	t.Log(dInfo.CurrentLocale)
}

func TestSession_BatteryInfo(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	s, err := c.NewSession()
	checkErr(t, err)
	Debug = true
	batteryInfo, err := s.BatteryInfo()
	checkErr(t, err)
	t.Log(batteryInfo)
	t.Log(batteryInfo.Level)
	t.Log(batteryInfo.State)
}

func TestSession_WindowSize(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	s, err := c.NewSession()
	checkErr(t, err)
	Debug = true
	windowSize, err := s.WindowSize()
	checkErr(t, err)
	t.Log(windowSize)
}

func TestSession_Screen(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	s, err := c.NewSession()
	checkErr(t, err)
	Debug = true
	screen, err := s.Screen()
	checkErr(t, err)
	t.Log(screen)
}

func TestSession_Scale(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	s, err := c.NewSession()
	checkErr(t, err)
	Debug = true
	scale, err := s.Scale()
	checkErr(t, err)
	t.Log(scale)
}

func TestSession_StatusBarSize(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	s, err := c.NewSession()
	checkErr(t, err)
	Debug = true
	statusBarSize, err := s.StatusBarSize()
	checkErr(t, err)
	t.Log(statusBarSize)
	t.Log(statusBarSize.Width)
	t.Log(statusBarSize.Height)
}

func TestSession_ActiveAppInfo(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	s, err := c.NewSession()
	checkErr(t, err)
	appInfo, err := s.ActiveAppInfo()
	checkErr(t, err)
	t.Log(appInfo)
}

func TestSession_ActiveAppsList(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	s, err := c.NewSession()
	checkErr(t, err)
	Debug = true
	appsList, err := s.ActiveAppsList()
	checkErr(t, err)
	t.Log(appsList)
}

func TestSession_Tap(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	s, err := c.NewSession()
	checkErr(t, err)
	Debug = true
	err = s.Tap(230, 130)
	checkErr(t, err)
	// err = s.Tap(210, 290)
	// if err != nil {
	// 	t.Fatal(err)
	// }
}

func TestSession_DoubleTap(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	s, err := c.NewSession()
	checkErr(t, err)
	Debug = true
	err = s.DoubleTap(230, 130)
	checkErr(t, err)
}

func TestSession_TouchAndHold(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	s, err := c.NewSession()
	checkErr(t, err)
	Debug = true
	// err = s.TouchAndHold(210, 290)
	err = s.TouchAndHold(230, 130)
	// err = s.TouchAndHoldFloat(230, 130, 2.5)
	checkErr(t, err)
}

func TestSession_ForceTouch(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	s, err := c.NewSession()
	checkErr(t, err)
	Debug = true
	// err = s.ForceTouch(230, 130, 1)
	err = s.ForceTouch(230, 130, 1, 0.5)
	checkErr(t, err)
}

func TestSession_Drag(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	s, err := c.NewSession()
	checkErr(t, err)
	Debug = true
	// err = s.Drag(230, 130, 230, 480,2)
	// err = s.Drag(230, 130, 230, 480)
	err = s.Drag(230, 130, 230, 30)
	// err = s.Drag(230, 130, 130, 130)
	// err = s.Drag(230, 130, 330, 130)
	checkErr(t, err)
}

func TestSession_Swipe(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	s, err := c.NewSession()
	checkErr(t, err)
	Debug = true
	err = s.Swipe(230, 130, 230*2, 130*2)
	checkErr(t, err)
}

func TestSession_SwipeSwipeDirection(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	s, err := c.NewSession()
	checkErr(t, err)
	Debug = true
	// err = s.SwipeUp()
	err = s.SwipeDown()
	// err = s.SwipeLeft()
	// err = s.SwipeRight()
	checkErr(t, err)
}

func TestSession_AppLaunch(t *testing.T) {
	Debug = true
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	s, err := c.NewSession()
	checkErr(t, err)
	// bundleId = "com.apple.AppStore"
	// _ = s.AppTerminate(bundleId)
	// bool 类型初始值为 false，在设置启动操作选项需要主动设置 ShouldWaitForQuiescence 为 true （如果需要的话）
	launchOpt := NewWDAAppLaunchOption().SetShouldWaitForQuiescence(true).SetArguments([]string{"-AppleLanguages", "(Russian)"})
	// launchOpt.SetEnvironment(map[string]string{"DYLD_PRINT_STATISTICS": "YES"})
	_ = launchOpt
	// 未解锁状态下启动指定 bundleId 会导致 wda 内部会报错
	// 虽然点亮了屏幕，但是内部报错了 Enqueue Failure: Failed to launch com.apple.Preferences: 未能完成该操作。Unable to launch com.apple.Preferences because the device was not, or could not be, unlocked.
	// Unable to launch com.apple.Preferences because the device was not, or could not be, unlocked.
	// 如果一段时间内解锁，还是可以继续后续的操作
	err = s.AppLaunch(bundleId, launchOpt)
	// err = s.AppLaunch(bundleId)
	checkErr(t, err)
}

func TestSession_AppTerminate(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	bundleId := "com.apple.Preferences"
	s, err := c.NewSession()
	checkErr(t, err)
	Debug = true
	err = s.AppTerminate(bundleId)
	checkErr(t, err)
}

func TestSession_AppState(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	// bundleId := "com.apple.Preferences"
	s, err := c.NewSession()
	checkErr(t, err)
	Debug = true
	state, err := s.AppState(bundleId)
	checkErr(t, err)
	t.Log(state)
	t.Log("app 是否前台活动中", state == WDAAppRunningFront)
}

func TestSession_SendKeys(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	s, err := c.NewSession()
	checkErr(t, err)
	Debug = true
	err = s.SendKeys(bundleId + "\n")
	checkErr(t, err)

	file, _ := ioutil.ReadFile("/Users/hero/Documents/Workspace/Golang/gwda/examples/main.go")
	err = s.SendKeys(string(file), 30)
	checkErr(t, err)
}

func TestSession_FindElement(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	s, err := c.NewSession()
	checkErr(t, err)
	_ = s.AppLaunch(bundleId)
	Debug = true
	// **/XCUIElementTypeButton[`label == '允许' OR label == '好'`]
	// "XCUIElementTypeWindow/*/*[$type == 'XCUIElementTypeButton' AND label BEGINSWITH 'A'$]"

	element, err := s.FindElement(WDALocator{PartialLinkText: NewWDAElementAttribute().SetValue("中心")})
	// element, err := s.FindElement(WDALocator{LinkText: NewWDAElementAttribute().SetLabel("“ ”文件夹")})
	// element, err := s.FindElement(WDALocator{ClassName: WDAElementType{PageIndicator: true}})
	// element, err := s.FindElement(WDALocator{ClassChain: "**/XCUIElementTypeCell[`label == '关于本机' OR label == 'Siri信息播报'`]"})
	// element, err := s.FindElement(WDALocator{Predicate: "label = 'Alerts'"})
	// element, err := s.FindElement(WDALocator{Predicate: "type == 'XCUIElementTypeCell'"})
	// element, err := s.FindElement(WDALocator{Predicate: "type = 'XCUIElementTypeButton'"})
	// element, err := s.FindElement(WDALocator{Predicate: "selected == true"})
	// element, err := s.FindElement(WDALocator{Predicate: "type == 'XCUIElementTypeIcon'"})
	// element, err := s.FindElement(WDALocator{LinkText: NewWDAElementAttribute().SetType(WDAElementType{Icon: true})})
	// element, err := s.FindElement(WDALocator{ClassChain: "**/XCUIElementTypeCell[`label == '通知' OR label == '通知'`]"})

	checkErr(t, err)
	t.Log(element)

	t.Log(element.Rect())
	t.Log(element.Click())

	// if len(element) == 1 {
	// 	err := element[0].Click()
	// 	t.Log(err)
	// }
}

func TestSession_FindElements(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	s, err := c.NewSession()
	checkErr(t, err)
	Debug = true
	// elements, err := s.FindElements("partial link text", "label=看一看")
	// elements, err := s.FindElements("partial link text", "label=发现")
	// **/XCUIElementTypeButton[`label == '允许' OR label == '好'`]
	// "XCUIElementTypeWindow/*/*[$type == 'XCUIElementTypeButton' AND label BEGINSWITH 'A'$]"
	// elements, err := s.FindElements("class chain", "**/XCUIElementTypeButton[`label == '允许' OR label == '好' OR label == '仅在使用应用期间' OR label == '暂不'`]")
	// elements, err := s.FindElements("class chain", "**/XCUIElementTypeButton[`label == '允许' OR label == '好' OR label == '仅在使用应用期间' OR label == '暂不'`]")
	// elements, err := s.FindElements(WDALocator{Predicate: "label == 'Siri信息播报'"})
	// elements, err := s.FindElements(WDALocator{Predicate: "selected == true AND label == '通用'"})
	elements, err := s.FindElements(WDALocator{Predicate: "label == '通用'"})
	checkErr(t, err)
	t.Log(elements)
	t.Log(len(elements))

	t.Log(elements[0].Rect())
	t.Log(elements[0].Click())

	// if len(elements) == 1 {
	// 	err := elements[0].Click()
	// 	t.Log(err)
	// }
}

func TestSession_ActiveElement(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	s, err := c.NewSession()
	checkErr(t, err)
	Debug = true
	element, err := s.ActiveElement()
	checkErr(t, err)
	t.Log(element.Rect())
}

func TestSession_AlertSendKeys(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	s, err := c.NewSession()
	checkErr(t, err)
	Debug = true
	err = s.AlertSendKeys("test")
	checkErr(t, err)
}

func TestSession_AlertAccept(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	s, err := c.NewSession()
	checkErr(t, err)
	Debug = true
	// err = s.AlertAccept()
	// err = s.AlertAccept("好")
	err = s.AlertAccept("允许")
	checkErr(t, err)
}

func TestSession_AlertDismiss(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	s, err := c.NewSession()
	checkErr(t, err)
	Debug = true
	// err = s.AlertDismiss()
	err = s.AlertDismiss("不允许")
	checkErr(t, err)
}

func TestSession_AlertText(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	s, err := c.NewSession()
	checkErr(t, err)
	Debug = true
	text, err := s.AlertText()
	checkErr(t, err)
	t.Log(text)
}

func TestSession_AlertButtons(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	s, err := c.NewSession()
	checkErr(t, err)
	Debug = true
	buttons, err := s.AlertButtons()
	checkErr(t, err)
	t.Log(buttons)
}

func TestSession_Orientation(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	s, err := c.NewSession()
	checkErr(t, err)
	Debug = true
	orientation, err := s.Orientation()
	checkErr(t, err)
	t.Log(orientation)
}

func TestSession_SetOrientation(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	s, err := c.NewSession()
	checkErr(t, err)
	Debug = true
	err = s.SetOrientation(WDAOrientationPortraitUpsideDown)
	checkErr(t, err)
}

func TestSession_Rotation(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	s, err := c.NewSession()
	checkErr(t, err)
	Debug = true
	rotation, err := s.Rotation()
	checkErr(t, err)
	t.Log(rotation)
}

func TestSession_SetRotation(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	s, err := c.NewSession()
	checkErr(t, err)
	Debug = true
	err = s.SetRotation(WDARotation{X: 0, Y: 0, Z: 270})
	checkErr(t, err)
}

func TestSession_PerformAppiumTouchActions(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	s, err := c.NewSession()
	checkErr(t, err)
	element, err := s.FindElement(WDALocator{Name: "自定手势作用区域"})
	checkErr(t, err)
	Debug = true

	touchActions := NewWDATouchActions().
		Press(NewWDATouchActionOptionPress().SetElement(element).SetXY(200, 200).SetPressure(0.8)).
		// LongPress(NewWDATouchActionOptionLongPress().SetElement(element).SetXY(200, 200)).
		Wait(0.2).
		MoveTo(NewWDATouchActionOptionMoveTo().SetXY(300, 200)).
		Wait(0.2).
		MoveTo(NewWDATouchActionOptionMoveTo().SetElement(element)).
		Wait(0.2).
		MoveTo(NewWDATouchActionOptionMoveTo().SetElement(element).SetXY(300, 400)).
		Release()

	err = s.PerformTouchActions(touchActions)
	checkErr(t, err)
}

func TestSession_PerformActions(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	s, err := c.NewSession()
	checkErr(t, err)
	// element, err := s.FindElement(WDALocator{Name: "自定手势作用区域"})
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// actOptFingerLeft := NewWDAActionOptionFinger().
	// 	Move(NewWWDAActionOptionFingerMove().SetXY(-75, -185).SetOrigin(element)).
	// 	Down().
	// 	Pause(0.25).
	// 	Move(NewWWDAActionOptionFingerMove().SetOrigin(element)).
	// 	Pause(0.25).
	// 	Up()
	// actOptFingerRight := NewWDAActionOptionFinger().
	// 	Move(NewWWDAActionOptionFingerMove().SetXY(75, 185).SetOrigin(element)).
	// 	Down().
	// 	Pause(0.25).
	// 	Move(NewWWDAActionOptionFingerMove().SetOrigin(element)).
	// 	Pause(0.25).
	// 	Up()
	// _, _, _ = element, actOptFingerLeft, actOptFingerRight
	Debug = true

	// actions := NewWDAActions().Tap(80, 100)
	// actions := NewWDAActions().Tap(50, 0, element)
	// actions := NewWDAActions().Press(50, 0, 3, element)
	// actions := NewWDAActions().DoubleTap(0, 50, element)
	// actions := NewWDAActions().Swipe(-75, -185, 0, 0, element)
	// actions := NewWDAActions().Swipe(-75, -185, 0, 0, element).Swipe(75, 185, 0, 0, element)
	// actions := NewWDAActions().FingerActionOption(actOptFingerLeft).FingerActionOption(actOptFingerRight)
	actions := NewWDAActions().SendKeys("WebDriverAgent")

	err = s.PerformActions(actions)
	checkErr(t, err)
}

func TestSession_IsLocked(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	s, err := c.NewSession()
	checkErr(t, err)
	Debug = true
	locked, err := s.IsLocked()
	checkErr(t, err)
	if locked {
		t.Log("设备处于 屏幕锁定 界面")
	} else {
		t.Log("设备已屏幕解锁")
	}
}

func TestSession_Unlock(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	bundleId := "com.apple.Preferences"
	_ = bundleId
	s, err := c.NewSession()
	checkErr(t, err)
	Debug = true
	err = s.Unlock()
	checkErr(t, err)
}

func TestSession_Lock(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	bundleId := "com.apple.Preferences"
	_ = bundleId
	s, err := c.NewSession()
	checkErr(t, err)
	Debug = true
	err = s.Lock()
	checkErr(t, err)
}

func TestSession_AppActivate(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	// bundleId := "com.apple.Preferences"
	s, err := c.NewSession()
	checkErr(t, err)
	Debug = true
	err = s.AppActivate(bundleId)
	checkErr(t, err)
}

func TestSession_AppDeactivate(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	// bundleId := "com.apple.Preferences"
	_ = bundleId
	s, err := c.NewSession()
	checkErr(t, err)
	Debug = true
	err = s.AppDeactivate(10.6)
	checkErr(t, err)
}

func TestSession_SetPasteboardForPlaintext(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	bundleId := "com.apple.Preferences"
	_ = bundleId
	s, err := c.NewSession()
	checkErr(t, err)
	Debug = true
	err = s.SetPasteboardForPlaintext("abcd1234")
	checkErr(t, err)
}

func TestSession_SetPasteboardForImageFromFile(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	bundleId := "com.apple.Preferences"
	_ = bundleId
	s, err := c.NewSession()
	checkErr(t, err)
	Debug = true
	err = s.SetPasteboardForImageFromFile("/Users/hero/Documents/leixipaopao/IMG_5246.JPG")
	checkErr(t, err)
}

func TestSession_SetPasteboardForUrl(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	bundleId := "com.apple.Preferences"
	_ = bundleId
	s, err := c.NewSession()
	checkErr(t, err)
	Debug = true
	err = s.SetPasteboardForUrl("https://www.apple.com.cn")
	checkErr(t, err)
}

func TestSession_SetPasteboard(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	bundleId := "com.apple.Preferences"
	_ = bundleId
	s, err := c.NewSession()
	checkErr(t, err)
	Debug = true
	err = s.SetPasteboard(WDAContentTypeUrl, "https://www.google.com")
	checkErr(t, err)
}

func TestSession_GetPasteboard(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	s, err := c.NewSession()
	checkErr(t, err)
	Debug = true
	// raw, err := s.GetPasteboard(WDAContentTypePlaintext)
	// content, err := s.GetPasteboardForPlaintext()
	// url, err := s.GetPasteboardForUrl()
	// image, format, err := s.GetPasteboardForImage()
	userHomeDir, _ := os.UserHomeDir()
	err = s.GetPasteboardForImageToDisk(filepath.Join(userHomeDir, "Desktop", "s3.png"))
	checkErr(t, err)
	// t.Log(raw.String())
	// t.Log(content)
	// t.Log(url)
	// t.Log(format)
	// t.Log(image.Bounds().Size())
}

func TestSession_PressHomeButton(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	bundleId := "com.apple.Preferences"
	_ = bundleId
	s, err := c.NewSession()
	checkErr(t, err)
	Debug = true
	err = s.PressHomeButton()
	checkErr(t, err)
}

func TestSession_PressVolumeUpButton(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	bundleId := "com.apple.Preferences"
	_ = bundleId
	s, err := c.NewSession()
	checkErr(t, err)
	Debug = true
	err = s.PressVolumeUpButton()
	checkErr(t, err)
}

func TestSession_PressVolumeDownButton(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	bundleId := "com.apple.Preferences"
	_ = bundleId
	s, err := c.NewSession()
	checkErr(t, err)
	Debug = true
	err = s.PressVolumeDownButton()
	checkErr(t, err)
}

func TestSession_SiriActivate(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	s, err := c.NewSession()
	checkErr(t, err)
	Debug = true
	err = s.SiriActivate("查一下附近的餐厅")
	checkErr(t, err)
}

func TestSession_SiriOpenURL(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	s, err := c.NewSession()
	checkErr(t, err)
	Debug = true
	err = s.SiriOpenURL("https://apple.com")
	checkErr(t, err)
}

func TestSession_Screenshot(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	s, err := c.NewSession()
	checkErr(t, err)
	// defer s.DeleteSession()
	Debug = true
	_, err = s.Screenshot()
	checkErr(t, err)
	Debug = false
	element, err := s.FindElement(WDALocator{Predicate: "type == 'XCUIElementTypeCell' AND name == '通知'"})
	checkErr(t, err)

	Debug = true
	_, err = s.Screenshot(element)
	checkErr(t, err)
}

func TestSession_ScreenshotToDisk(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	s, err := c.NewSession()
	checkErr(t, err)
	// defer s.DeleteSession()
	Debug = true
	userHomeDir, _ := os.UserHomeDir()
	err = s.ScreenshotToDisk(filepath.Join(userHomeDir, "Desktop", "s1.png"))
	checkErr(t, err)
	Debug = false
	element, err := s.FindElement(WDALocator{Predicate: "type == 'XCUIElementTypeCell' AND name == '通知'"})
	checkErr(t, err)

	Debug = true
	err = s.ScreenshotToDisk(filepath.Join(userHomeDir, "Desktop", "s2.png"), element)
	checkErr(t, err)
}

func TestSession_ScreenshotToImage(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	s, err := c.NewSession()
	checkErr(t, err)
	// defer s.DeleteSession()
	Debug = true
	// toPng, err := s.ScreenshotToPng()
	img, format, err := s.ScreenshotToImage()
	checkErr(t, err)
	t.Log("图片的格式:", format)
	t.Log("图片的大小:", img.Bounds().Size())

	element, err := s.FindElement(WDALocator{Predicate: "type == 'XCUIElementTypeCell' AND name == '通知'"})
	checkErr(t, err)

	Debug = true
	img, format, err = s.ScreenshotToImage(element)
	checkErr(t, err)
	t.Log("元素图片的格式:", format)
	t.Log("元素图片的大小:", img.Bounds().Size())
	t.Log(element.Rect())
}

func TestSession_Source(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	s, err := c.NewSession()
	checkErr(t, err)
	Debug = true
	// sTree, err := s.Source()
	sTree, err := s.Source()
	checkErr(t, err)
	t.Log(sTree)
}

func TestSession_AccessibleSource(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	bundleId := "com.apple.Preferences"
	_ = bundleId
	s, err := c.NewSession()
	checkErr(t, err)
	Debug = true
	source, err := s.AccessibleSource()
	checkErr(t, err)
	t.Log(source)
}

func TestSession_GetAppiumSettings(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	// bundleId := "com.apple.Preferences"
	// _ = bundleId
	s, err := c.NewSession()
	checkErr(t, err)
	Debug = true
	settings, err := s.GetAppiumSettings()
	checkErr(t, err)
	t.Log(settings)
}

func TestTmpSession(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	s, err := c.NewSession()
	checkErr(t, err)
	Debug = true

	// err = s.AppLaunch("com.apple.calculator")
	// checkErr(t, err)
	// orientation, err := s.Orientation()
	// if orientation == WDAOrientationPortrait {
	// 	err = s.SetOrientation(WDAOrientationLandscapeLeft)
	// }
	//
	// userHomeDir, _ := os.UserHomeDir()
	// err = s.ScreenshotToDisk(filepath.Join(userHomeDir, "Desktop", "s4.png"))
	// checkErr(t, err)
	//
	// element, err := s.FindElement(WDALocator{Name: "("})
	// checkErr(t, err)
	//
	// err = element.ScreenshotToDisk(filepath.Join(userHomeDir, "Desktop", "e5.png"))
	// checkErr(t, err)
	//
	// err = s.ScreenshotToDisk(filepath.Join(userHomeDir, "Desktop", "e6.png"), element)
	// checkErr(t, err)
	//
	// return

	// s.SetOrientation(WDAOrientationLandscapeRight)
	s.tttTmp()
	// _ = s
}
