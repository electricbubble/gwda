package gwda

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

var urlPrefix = "http://localhost:8100"
var bundleId = "com.apple.Preferences"
var driver WebDriver

func setup(t *testing.T) {
	var err error
	driver, err = NewUSBDriver(nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestViaUSB(t *testing.T) {
	devices, err := DeviceList()
	if err != nil {
		t.Fatal(err)
	}

	drivers := make([]WebDriver, 0, len(devices))

	for _, dev := range devices {
		d, err := NewUSBDriver(nil, dev)
		if err != nil {
			t.Errorf("%s: %s", dev.SerialNumber(), err)
			continue
		}
		drivers = append(drivers, d)
	}

	for _, d := range drivers {
		t.Log(d.Status())
	}
}

func TestNewDevice(t *testing.T) {
	device, _ := NewDevice()
	if device != nil {
		t.Log(device)
	}

	device, _ = NewDevice(WithSerialNumber("xxxx"))
	if device != nil {
		t.Log(device)
	}

	device, _ = NewDevice(WithPort(8700), WithMjpegPort(8800))
	if device != nil {
		t.Log(device)
	}

	device, _ = NewDevice(WithSerialNumber("xxxx"), WithPort(8700), WithMjpegPort(8800))
	if device != nil {
		t.Log(device)
	}
}

func TestNewDriver(t *testing.T) {
	var err error
	driver, err = NewDriver(nil, urlPrefix)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewUSBDriver(t *testing.T) {
	setup(t)

	// t.Log(driver.IsWdaHealthy())
}

func Test_remoteWD_NewSession(t *testing.T) {
	setup(t)

	// sessionInfo, err := driver.NewSession(nil)
	sessionInfo, err := driver.NewSession(
		NewCapabilities().WithAppLaunchOption(
			NewAppLaunchOption().WithBundleId(bundleId).WithArguments([]string{"-AppleLanguages", "(Russian)"}),
		),
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(sessionInfo.SessionId) == 0 {
		t.Fatal(sessionInfo)
	}
}

func Test_remoteWD_ActiveSession(t *testing.T) {
	setup(t)

	sessionInfo, err := driver.ActiveSession()
	if err != nil {
		t.Fatal(err)
	}
	if len(sessionInfo.SessionId) == 0 {
		t.Fatal(sessionInfo)
	}
}

func Test_remoteWD_DeleteSession(t *testing.T) {
	setup(t)

	err := driver.DeleteSession()
	if err != nil {
		t.Fatal(err)
	}
}

func Test_remoteWD_HealthCheck(t *testing.T) {
	setup(t)

	err := driver.HealthCheck()
	if err != nil {
		t.Fatal(err)
	}
}

func Test_remoteWD_GetAppiumSettings(t *testing.T) {
	setup(t)

	settings, err := driver.GetAppiumSettings()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(settings)
}

func Test_remoteWD_SetAppiumSettings(t *testing.T) {
	setup(t)

	const _acceptAlertButtonSelector = "**/XCUIElementTypeButton[`label IN {'允许','好','仅在使用应用期间','暂不'}`]"
	const _dismissAlertButtonSelector = "**/XCUIElementTypeButton[`label IN {'不允许','暂不'}`]"

	key := "acceptAlertButtonSelector"
	value := _acceptAlertButtonSelector

	// settings, err := driver.SetAppiumSettings(map[string]interface{}{"dismissAlertButtonSelector": "暂不"})
	settings, err := driver.SetAppiumSettings(map[string]interface{}{key: value})
	if err != nil {
		t.Fatal(err)
	}
	if settings[key] != value {
		t.Fatal(settings[key])
	}
}

func Test_remoteWD_IsWdaHealthy(t *testing.T) {
	setup(t)

	healthy, err := driver.IsWdaHealthy()
	if err != nil {
		t.Fatal(err)
	}
	if healthy == false {
		t.Fatal("healthy =", healthy)
	}
}

// func Test_remoteWD_WdaShutdown(t *testing.T) {
// 	setup(t)
//
// 	if err := driver.WdaShutdown(); err != nil {
// 		t.Fatal(err)
// 	}
// }

func Test_remoteWD_Status(t *testing.T) {
	setup(t)

	status, err := driver.Status()
	if err != nil {
		t.Fatal(err)
	}
	if status.Ready == false {
		t.Fatal("deviceStatus =", status)
	}
}

func Test_remoteWD_DeviceInfo(t *testing.T) {
	setup(t)

	SetDebug(true)

	info, err := driver.DeviceInfo()
	if err != nil {
		t.Fatal(err)
	}
	if len(info.Model) == 0 {
		t.Fatal(info)
	}
}

func Test_remoteWD_BatteryInfo(t *testing.T) {
	setup(t)

	batteryInfo, err := driver.BatteryInfo()
	if err != nil {
		t.Fatal()
	}
	t.Log(batteryInfo)
}

func Test_remoteWD_WindowSize(t *testing.T) {
	setup(t)

	size, err := driver.WindowSize()
	if err != nil {
		t.Fatal()
	}
	t.Log(size)
}

func Test_remoteWD_Screen(t *testing.T) {
	setup(t)

	screen, err := driver.Screen()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(screen)
}

func Test_remoteWD_ActiveAppInfo(t *testing.T) {
	setup(t)

	appInfo, err := driver.ActiveAppInfo()
	if err != nil {
		t.Fatal(err)
	}
	if len(appInfo.BundleId) == 0 {
		t.Fatal(appInfo)
	}
	t.Log(appInfo)
}

func Test_remoteWD_ActiveAppsList(t *testing.T) {
	setup(t)

	appsList, err := driver.ActiveAppsList()
	if err != nil {
		t.Fatal(err)
	}
	if len(appsList) == 0 {
		t.Fatal(appsList)
	}
	t.Log(appsList)
}

func Test_remoteWD_AppState(t *testing.T) {
	setup(t)

	runState, err := driver.AppState(bundleId)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(runState)
}

func Test_remoteWD_IsLocked(t *testing.T) {
	setup(t)

	locked, err := driver.IsLocked()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(locked)
}

func Test_remoteWD_Unlock(t *testing.T) {
	setup(t)

	err := driver.Unlock()
	if err != nil {
		t.Fatal(err)
	}
}

func Test_remoteWD_Lock(t *testing.T) {
	setup(t)

	err := driver.Lock()
	if err != nil {
		t.Fatal(err)
	}
}

func Test_remoteWD_AlertText(t *testing.T) {
	setup(t)

	text, err := driver.AlertText()
	if err != nil {
		t.Fatal(err)
	}
	_ = text
	t.Log(text)
}

func Test_remoteWD_AlertButtons(t *testing.T) {
	setup(t)

	btnLabels, err := driver.AlertButtons()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(btnLabels)
}

func Test_remoteWD_AlertAccept(t *testing.T) {
	// Test_remoteWD_AppAuthReset(t)
	// return

	setup(t)

	err := driver.AlertAccept()
	// err := driver.AlertAccept("")
	// err := driver.AlertAccept("好")
	if err != nil {
		t.Fatal(err)
	}
}

func Test_remoteWD_AlertDismiss(t *testing.T) {
	// Test_remoteWD_AppAuthReset(t)
	// return

	setup(t)

	err := driver.AlertDismiss()
	// err := driver.AlertDismiss("")
	// err := driver.AlertDismiss("不允许")
	if err != nil {
		t.Fatal(err)
	}
}
func Test_remoteWD_AlertSendKeys(t *testing.T) {
	setup(t)

	err := driver.AlertSendKeys("todo")
	if err != nil {
		t.Fatal(err)
	}
}

func Test_remoteWD_Homescreen(t *testing.T) {
	setup(t)

	err := driver.Homescreen()
	if err != nil {
		t.Fatal(err)
	}
}

func Test_remoteWD_AppLaunch(t *testing.T) {
	setup(t)

	// SetDebug(true)

	// bundleId = "com.hustlzp.xcz"
	// bundleId = "com.github.stormbreaker.prod"
	// bundleId = "com.360buy.jdmobile"
	// bundleId = "com.zhihu.ios"
	// bundleId = "com.tencent.xin"
	// bundleId = "com.jsmcc.ZP7267A6ES"
	err := driver.AppLaunch(bundleId)
	// err := driver.AppLaunch(bundleId, NewAppLaunchOption().WithShouldWaitForQuiescence(true))
	// err := driver.AppLaunch(bundleId, NewAppLaunchOption().WithArguments([]string{"-AppleLanguages", "(Russian)"}))
	if err != nil {
		t.Fatal(err)
	}
}

func Test_remoteWD_AppLaunchUnattached(t *testing.T) {
	setup(t)

	err := driver.AppLaunchUnattached(bundleId)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_remoteWD_AppTerminate(t *testing.T) {
	setup(t)

	_, err := driver.AppTerminate(bundleId)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_remoteWD_AppActivate(t *testing.T) {
	setup(t)

	err := driver.AppActivate(bundleId)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_remoteWD_AppDeactivate(t *testing.T) {
	setup(t)

	err := driver.AppDeactivate(2)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_remoteWD_AppAuthReset(t *testing.T) {
	setup(t)

	err := driver.AppAuthReset(ProtectedResourceCamera)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_remoteWD_Tap(t *testing.T) {
	setup(t)

	err := driver.Tap(200, 300)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_remoteWD_DoubleTap(t *testing.T) {
	setup(t)

	err := driver.DoubleTap(200, 300)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_remoteWD_TouchAndHold(t *testing.T) {
	setup(t)

	// err := driver.TouchAndHold(200, 300)
	err := driver.TouchAndHold(200, 300, -1)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_remoteWD_Drag(t *testing.T) {
	setup(t)

	// err := driver.Drag(200, 300, 200, 500, -1)
	err := driver.Swipe(200, 300, 200, 500)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_remoteWD_ForceTouch(t *testing.T) {
	setup(t)

	SetDebug(true)
	err := driver.ForceTouch(256, 400, 0.8, -1)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_remoteWD_SetPasteboard(t *testing.T) {
	setup(t)

	// err := driver.SetPasteboard(PasteboardTypePlaintext, "gwda")
	err := driver.SetPasteboard(PasteboardTypeUrl, "Clock-stopwatch://")
	// userHomeDir, _ := os.UserHomeDir()
	// bytesImg, _ := ioutil.ReadFile(userHomeDir + "/Pictures/IMG_0806.jpg")
	// err := driver.SetPasteboard(PasteboardTypeImage, string(bytesImg))
	if err != nil {
		t.Fatal(err)
	}
}

func Test_remoteWD_GetPasteboard(t *testing.T) {
	setup(t)

	var buffer *bytes.Buffer
	var err error

	buffer, err = driver.GetPasteboard(PasteboardTypePlaintext)
	// buffer, err = driver.GetPasteboard(PasteboardTypeUrl)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(buffer.String())

	// buffer, err = driver.GetPasteboard(PasteboardTypeImage)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// userHomeDir, _ := os.UserHomeDir()
	// if err = ioutil.WriteFile(userHomeDir+"/Desktop/p1.png", buffer.Bytes(), 0600); err != nil {
	// 	t.Error(err)
	// }
}

func Test_remoteWD_SendKeys(t *testing.T) {
	setup(t)

	err := driver.SendKeys("App Store")
	// err := driver.SendKeys("App Store", 3)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_remoteWD_PressButton(t *testing.T) {
	setup(t)

	err := driver.PressButton(DeviceButtonVolumeUp)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 1)
	err = driver.PressButton(DeviceButtonVolumeDown)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 1)
	err = driver.PressButton(DeviceButtonHome)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_remoteWD_SiriActivate(t *testing.T) {
	setup(t)

	err := driver.SiriActivate("What's the weather like today")
	if err != nil {
		t.Fatal(err)
	}
}

func Test_remoteWD_SiriOpenUrl(t *testing.T) {
	setup(t)

	err := driver.SiriOpenUrl("Prefs:root=Bluetooth")
	// err := driver.SiriOpenUrl("Prefs:root=WIFI![]()")
	if err != nil {
		t.Fatal(err)
	}
}

func Test_remoteWD_Orientation(t *testing.T) {
	setup(t)

	orientation, err := driver.Orientation()
	if err != nil {
		t.Fatal(err)
	}
	if orientation == "" {
		t.Fatal(orientation)
	}
}

func Test_remoteWD_SetOrientation(t *testing.T) {
	setup(t)

	var err error
	err = driver.SetOrientation(OrientationLandscapeLeft)
	err = driver.SetOrientation(OrientationLandscapeRight)
	err = driver.SetOrientation(OrientationPortrait)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_remoteWD_Rotation(t *testing.T) {
	setup(t)

	rotation, err := driver.Rotation()
	if err != nil {
		t.Fatal()
	}
	t.Log(rotation)
}

func Test_remoteWD_SetRotation(t *testing.T) {
	setup(t)

	err := driver.SetRotation(Rotation{X: 0, Y: 0, Z: 270})
	if err != nil {
		t.Fatal(err)
	}
}

func Test_remoteWD_PerformW3CActions(t *testing.T) {
	// setup(t)
	// actions := NewW3CActions().SendKeys("App Store")

	element := setupElement(t, BySelector{Name: "touchableView"})
	actions := NewW3CActions().FingerAction(
		NewFingerAction().
			Move(NewFingerMove().WithXY(-15, -85).WithOrigin(element)).
			Down().
			Pause(0.25).
			Move(NewFingerMove().WithOrigin(element)).
			Pause(0.25).
			Up(),
		NewFingerAction().
			Move(NewFingerMove().WithXY(15, 85).WithOrigin(element)).
			Down().
			Pause(0.25).
			Move(NewFingerMove().WithOrigin(element)).
			Pause(0.25).
			Up(),
	)
	err := driver.PerformW3CActions(actions)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_remoteWD_PerformAppiumTouchActions(t *testing.T) {
	element := setupElement(t, BySelector{Name: "touchableView"})

	actions := NewTouchActions().
		Press(NewTouchActionPress().WithElement(element).WithXY(100, 150).WithPressure(0.2)).
		Wait(0.2).
		MoveTo(NewTouchActionMoveTo().WithXY(300, 150)).
		Wait(0.2).
		MoveTo(NewTouchActionMoveTo().WithElement(element)).
		Wait(0.2).
		MoveTo(NewTouchActionMoveTo().WithElement(element).WithXY(300, 400)).
		Release()

	err := driver.PerformAppiumTouchActions(actions)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_remoteWD_ActiveElement(t *testing.T) {
	setup(t)

	element, err := driver.ActiveElement()
	if err != nil {
		t.Fatal(err)
	}
	_ = element
	// t.Log(element)
}

func Test_remoteWD_FindElement(t *testing.T) {
	setup(t)

	element, err := driver.FindElement(BySelector{Name: "设置"})
	if err != nil {
		t.Fatal(err)
	}
	_ = element
	// t.Log(element)
}

func Test_remoteWD_FindElements(t *testing.T) {
	setup(t)

	elements, err := driver.FindElements(BySelector{ClassName: ElementType{Icon: true}})
	if err != nil {
		t.Fatal(err)
	}
	_ = elements
	t.Log(elements)
}

func Test_remoteWD_Screenshot(t *testing.T) {
	setup(t)

	screenshot, err := driver.Screenshot()
	if err != nil {
		t.Fatal(err)
	}
	_ = screenshot

	// img, format, err := image.Decode(screenshot)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// userHomeDir, _ := os.UserHomeDir()
	// file, err := os.Create(userHomeDir + "/Desktop/s1." + format)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// defer func() { _ = file.Close() }()
	// switch format {
	// case "png":
	// 	err = png.Encode(file, img)
	// case "jpeg":
	// 	err = jpeg.Encode(file, img, nil)
	// }
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// t.Log(file.Name())
}

func Test_remoteWD_Source(t *testing.T) {
	setup(t)

	var source string
	var err error

	SetDebug(true)

	// source, err = driver.Source()
	// if err != nil {
	// 	t.Fatal(err)
	// }

	source, err = driver.Source(NewSourceOption().WithScope("AppiumAUT"))
	if err != nil {
		t.Fatal(err)
	}

	// source, err = driver.Source(NewSourceOption().WithFormatAsJson())
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// source, err = driver.Source(NewSourceOption().WithFormatAsDescription())
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// source, err = driver.Source(NewSourceOption().WithFormatAsXml().WithExcludedAttributes([]string{"label", "type", "index"}))
	// if err != nil {
	// 	t.Fatal(err)
	// }

	_ = source
	fmt.Println(source)
}

func Test_remoteWD_AccessibleSource(t *testing.T) {
	setup(t)

	source, err := driver.AccessibleSource()
	if err != nil {
		t.Fatal(err)
	}
	_ = source
	fmt.Println(source)
}

func Test_remoteWD_Wait(t *testing.T) {
	setup(t)

	var element WebElement
	var err error

	by := BySelector{Name: "通知"}
	// driver.AppLaunch()
	exists := func(d WebDriver) (bool, error) {
		element, err = d.FindElement(by)
		if err == nil {
			return true, nil
		}
		return false, nil
	}
	_ = exists
	_ = element

	err = driver.AppLaunchUnattached(bundleId)
	if err != nil {
		t.Fatal(err)
	}
	// element, err = driver.FindElement(by)
	err = driver.WaitWithTimeoutAndInterval(exists, time.Second*10, time.Millisecond*10)
	if err != nil {
		t.Fatal(err)
	}

	// t.Log(element.Rect())
}

func Test_remoteWD_Location(t *testing.T) {
	setup(t)

	location, err := driver.Location()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(location)
}

func Test_remoteWD_KeyboardDismiss(t *testing.T) {
	setup(t)
	SetDebug(true)

	err := driver.KeyboardDismiss()
	if err != nil {
		t.Fatal(err)
	}
}

func Test_remoteWD_ExpectNotification(t *testing.T) {
	setup(t)
	SetDebug(true)

	// bundleId = "com.apple.shortcuts"
	// err := driver.ExpectNotification("shortcuts", NotificationTypePlain, 10)
	// if err != nil {
	// 	t.Fatal(err)
	// }
}

func Test_remoteWD_IOHIDEvent(t *testing.T) {
	setup(t)
	SetDebug(true)

	err := driver.IOHIDEvent(EventPageIDConsumer, EventUsageIDCsmrVolumeDown)
	if err != nil {
		t.Fatal(err)
	}
}
