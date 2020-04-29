package gwda

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
)

type Session struct {
	sessionURL *url.URL
}

func newSession(deviceURL *url.URL, sid string) (s *Session) {
	s = new(Session)
	s.sessionURL, _ = url.Parse(deviceURL.String() + "/session/" + sid)
	return
}

type WDASessionInfo struct {
	Capabilities struct {
		CFBundleIdentifier string `json:"CFBundleIdentifier"`
		BrowserName        string `json:"browserName"`
		Device             string `json:"device"`
		SdkVersion         string `json:"sdkVersion"`
	} `json:"capabilities"`
	SessionID string `json:"sessionId"`
	_string   string
}

func (si WDASessionInfo) String() string {
	return si._string
}

// GetActiveSession
//
// get current session information
//
// {
//    "sessionId" : "8BF16568-832F-4A14-A137-FD0CA566FC64",
//    "capabilities" : {
//      "device" : "iphone",
//      "browserName" : "设置",
//      "sdkVersion" : "11.4.1",
//      "CFBundleIdentifier" : "com.apple.Preferences"
//    }
// }
func (s *Session) GetActiveSession() (wdaSessionInfo WDASessionInfo, err error) {
	var wdaResp wdaResponse
	if wdaResp, err = internalGet("GetActiveSession", urlJoin(s.sessionURL, "")); err != nil {
		return WDASessionInfo{}, err
	}

	wdaSessionInfo._string = wdaResp.getValue().String()
	err = json.Unmarshal([]byte(wdaSessionInfo._string), &wdaSessionInfo)
	return
}

// DeleteSession
//
// kill session (and App) associated with that request
//
//	1. alertsMonitor disable
//	2. testedApplicationBundleId terminate
func (s *Session) DeleteSession() (err error) {
	_, err = internalDelete("DeleteSession", s.sessionURL.String())
	return
}

// launch application configuration
type WDAAppLaunchOption wdaBody

func NewWDAAppLaunchOption() WDAAppLaunchOption {
	return make(WDAAppLaunchOption)
}

// SetShouldWaitForQuiescence
//
// It allows to turn on/off waiting for application quiescence, while performing queries.
func (alo WDAAppLaunchOption) SetShouldWaitForQuiescence(b bool) WDAAppLaunchOption {
	return WDAAppLaunchOption(wdaBody(alo).set("shouldWaitForQuiescence", b))
}

// SetArguments
//
// The optional array of application command line arguments. The arguments are going to be applied if the application was not running before.
func (alo WDAAppLaunchOption) SetArguments(args []string) WDAAppLaunchOption {
	return WDAAppLaunchOption(wdaBody(alo).set("arguments", args))
}

// SetEnvironment
//
// The optional dictionary of environment variables for the application, which is going to be executed. The environment variables are going to be applied if the application was not running before.
func (alo WDAAppLaunchOption) SetEnvironment(env map[string]string) WDAAppLaunchOption {
	return WDAAppLaunchOption(wdaBody(alo).set("environment", env))
}

// AppLaunch
//
// Launch an application with given bundle identifier in scope of current session.
// !This method is only available since Xcode9 SDK
//
// Default wait for quiescence
//
//	1. registerApplicationWithBundleId
//	2. launch OR activate
func (s *Session) AppLaunch(bundleId string, opt ...WDAAppLaunchOption) (err error) {
	// TODO BundleId is required 如果是不存在的 bundleId 会导致 wda 内部报错导致接下来的操作都无法接收处理
	if len(opt) == 0 {
		opt = []WDAAppLaunchOption{NewWDAAppLaunchOption().SetShouldWaitForQuiescence(true)}
	}
	body := newWdaBody().setBundleID(bundleId)
	body.setAppLaunchOption(opt[0])
	_, err = internalPost("AppLaunch", urlJoin(s.sessionURL, "/wda/apps/launch"), body)
	return
}

type WDADeviceInfo struct {
	TimeZone           string `json:"timeZone"`
	CurrentLocale      string `json:"currentLocale"`
	Model              string `json:"model"`
	UUID               string `json:"uuid"`
	UserInterfaceIdiom int    `json:"userInterfaceIdiom"`
	UserInterfaceStyle string `json:"userInterfaceStyle"`
	Name               string `json:"name"`
	IsSimulator        bool   `json:"isSimulator"`
	_string            string
}

func (di WDADeviceInfo) String() string {
	return di._string
}

// DeviceInfo
//
// {
//    "timeZone" : "Asia\/Shanghai",
//    "currentLocale" : "zh_CN",
//    "model" : "iPhone",
//    "uuid" : "x-x-x-x-x",
//    "userInterfaceIdiom" : 0,
//    "userInterfaceStyle" : "unsupported",
//    "name" : "TEST’s iPhone",
//    "isSimulator" : false
//  }
func (s *Session) DeviceInfo() (wdaDeviceInfo WDADeviceInfo, err error) {
	return deviceInfo(s.sessionURL)
}

type WDABatteryInfo struct {
	Level   float64         `json:"level"` // Battery level in range [0.0, 1.0], where 1.0 means 100% charge.
	State   WDABatteryState `json:"state"` // Battery state ( 1: on battery, discharging; 2: plugged in, less than 100%, 3: plugged in, at 100% )
	_string string
}

func (bi WDABatteryInfo) String() string {
	return bi._string
}

type WDABatteryState int

const (
	_                                   = iota
	WDABatteryUnplugged WDABatteryState = iota // on battery, discharging
	WDABatteryCharging                         // plugged in, less than 100%
	WDABatteryFull                             // plugged in, at 100%
)

func (v WDABatteryState) String() string {
	switch v {
	case WDABatteryUnplugged:
		return "On battery, discharging"
	case WDABatteryCharging:
		return "Plugged in, less than 100%"
	case WDABatteryFull:
		return "Plugged in, at 100%"
	default:
		return "UNKNOWN"
	}
}

// BatteryInfo
//
// {
//    "level": 0.92000001668930054,
//    "state": 2
// }
//
// level - Battery level in range [0.0, 1.0], where 1.0 means 100% charge.
// state - Battery state. The following values are possible:
// UIDeviceBatteryStateUnplugged = 1  // on battery, discharging
// UIDeviceBatteryStateCharging = 2   // plugged in, less than 100%
// UIDeviceBatteryStateFull = 3       // plugged in, at 100%
func (s *Session) BatteryInfo() (wdaBatteryInfo WDABatteryInfo, err error) {
	var wdaResp wdaResponse
	if wdaResp, err = internalGet("BatteryInfo", urlJoin(s.sessionURL, "/wda/batteryInfo")); err != nil {
		return
	}

	wdaBatteryInfo._string = wdaResp.getValue().String()
	err = json.Unmarshal([]byte(wdaBatteryInfo._string), &wdaBatteryInfo)
	// err = json.Unmarshal(wdaResp.getValue2Bytes(), &wdaBatteryInfo)
	return
}

// WindowSize
//
// {
//    "width": 812,
//    "height": 375
// }
func (s *Session) WindowSize() (wdaSize WDASize, err error) {
	var wdaResp wdaResponse
	if wdaResp, err = internalGet("WindowSize", urlJoin(s.sessionURL, "/window/size")); err != nil {
		return
	}

	wdaSize._string = wdaResp.getValue().String()
	err = json.Unmarshal([]byte(wdaSize._string), &wdaSize)
	// err = json.Unmarshal(wdaResp.getValue2Bytes(), &wdaSize)
	return
}

type WDASize struct {
	Width   int `json:"width"`
	Height  int `json:"height"`
	_string string
}

func (s WDASize) String() string {
	return s._string
}

type WDAScreen struct {
	StatusBarSize WDASize `json:"statusBarSize"`
	Scale         float64 `json:"scale"`
	_string       string
}

func (s WDAScreen) String() string {
	return s._string
}

// Screen
//
// {
//    "statusBarSize": {
//        "width": 375,
//        "height": 44
//    },
//    "scale": 3
// }
func (s *Session) Screen() (wdaScreen WDAScreen, err error) {
	var wdaResp wdaResponse
	if wdaResp, err = internalGet("Screen", urlJoin(s.sessionURL, "/wda/screen")); err != nil {
		return
	}

	wdaScreen.StatusBarSize._string = wdaResp.getValue().Get("statusBarSize").String()
	wdaScreen._string = wdaResp.getValue().String()
	err = json.Unmarshal([]byte(wdaScreen._string), &wdaScreen)
	// err = json.Unmarshal(wdaResp.getValue2Bytes(), &wdaScreen)
	return
}

// Scale
func (s *Session) Scale() (scale float64, err error) {
	screen, err := s.Screen()
	return screen.Scale, err
}

// StatusBarSize
func (s *Session) StatusBarSize() (wdaStatusBarSize WDASize, err error) {
	screen, err := s.Screen()
	return screen.StatusBarSize, err
}

// ActiveAppInfo
//
// get current active application
func (s *Session) ActiveAppInfo() (wdaActiveAppInfo WDAActiveAppInfo, err error) {
	return activeAppInfo(s.sessionURL)
}

// ActiveAppsList
//
// use multitasking on iPad
//
// [
//    {
//      "pid" : 3573,
//      "bundleId" : "com.apple.DocumentsApp"
//    },
//    {
//      "pid" : 3311,
//      "bundleId" : "com.apple.reminders"
//    }
//  ]
func (s *Session) ActiveAppsList() (appsList []WDAAppBaseInfo, err error) {
	var wdaResp wdaResponse
	if wdaResp, err = internalGet("ActiveAppsList", urlJoin(s.sessionURL, "/wda/apps/list")); err != nil {
		return nil, err
	}
	appsList = make([]WDAAppBaseInfo, 0)
	err = json.Unmarshal([]byte(wdaResp.getValue().String()), &appsList)
	return
}

func tap(baseUrl *url.URL, x, y interface{}, elemUID ...string) (err error) {
	body := newWdaBody().set("x", x).set("y", y)
	// [FBRoute POST:@"/wda/tap/:uuid"]
	tmpPath := "/wda/tap"
	if len(elemUID) == 0 {
		tmpPath += "/0"
	} else {
		tmpPath += "/" + elemUID[0]
	}
	_, err = internalPost("Tap", urlJoin(baseUrl, tmpPath), body)
	return
}

// Tap
func (s *Session) Tap(x, y int) error {
	return tap(s.sessionURL, x, y)
}

// TapFloat
func (s *Session) TapFloat(x, y float64) error {
	return tap(s.sessionURL, x, y)
}

// WDACoordinate
func (s *Session) TapCoordinate(wdaCoordinate WDACoordinate) error {
	return tap(s.sessionURL, wdaCoordinate.X, wdaCoordinate.Y)
}

// doubleTap
//
// [FBRoute POST:@"/wda/doubleTap"]
// [FBRoute POST:@"/wda/element/:uuid/doubleTap"]
func doubleTap(baseUrl *url.URL, x, y interface{}, elemPrefixPath ...string) (err error) {
	body := newWdaBody()
	tmpPath := "/doubleTap"
	if len(elemPrefixPath) == 0 {
		body.set("x", x).set("y", y)
	} else {
		tmpPath = elemPrefixPath[0] + tmpPath
	}
	_, err = internalPost("DoubleTap", urlJoin(baseUrl, tmpPath, true), body)
	return
}

// DoubleTap
//
// double tap coordinate
func (s *Session) DoubleTap(x, y int) (err error) {
	return doubleTap(s.sessionURL, x, y)
}

func (s *Session) DoubleTapFloat(x, y float64) (err error) {
	return doubleTap(s.sessionURL, x, y)
}

// touchAndHold
//
// [FBRoute POST:@"/wda/touchAndHold"]
// [FBRoute POST:@"/wda/element/:uuid/touchAndHold"]
func touchAndHold(baseUrl *url.URL, x, y, duration interface{}, elemPrefixPath ...string) (err error) {
	body := newWdaBody().set("duration", duration)
	tmpPath := "/touchAndHold"
	if len(elemPrefixPath) == 0 {
		body.set("x", x).set("y", y)
	} else {
		tmpPath = elemPrefixPath[0] + tmpPath
	}
	_, err = internalPost("TouchAndHold", urlJoin(baseUrl, tmpPath, true), body)
	return
}

// TouchAndHold
//
// touch and hold coordinate
func (s *Session) TouchAndHold(x, y int, duration ...int) (err error) {
	if len(duration) == 0 {
		duration = []int{1}
	}
	return touchAndHold(s.sessionURL, x, y, duration[0])
}

func (s *Session) TouchAndHoldFloat(x, y float64, duration ...float64) (err error) {
	if len(duration) == 0 {
		duration = []float64{1.0}
	}
	return touchAndHold(s.sessionURL, x, y, duration[0])
}

// drag
//
// [FBRoute POST:@"/wda/dragfromtoforduration"]
// [FBRoute POST:@"/wda/element/:uuid/dragfromtoforduration"]
func drag(baseUrl *url.URL, fromX, fromY, toX, toY, pressForDuration interface{}, elemPrefixPath ...string) (err error) {
	body := newWdaBody().set("duration", pressForDuration)
	body.set("fromX", fromX).set("fromY", fromY)
	body.set("toX", toX).set("toY", toY)
	tmpPath := "/dragfromtoforduration"
	if len(elemPrefixPath) != 0 {
		tmpPath = elemPrefixPath[0] + tmpPath
	}
	_, err = internalPost("Drag", urlJoin(baseUrl, tmpPath, true), body)
	return
}

// Drag
//
// Clicks and holds for a specified duration (generally long enough to start a drag operation) then drags to the other coordinate.
func (s *Session) Drag(fromX, fromY, toX, toY int, pressForDuration ...int) (err error) {
	if len(pressForDuration) == 0 {
		pressForDuration = []int{1}
	}
	return drag(s.sessionURL, fromX, fromY, toX, toY, pressForDuration[0])
}

func (s *Session) DragFloat(fromX, fromY, toX, toY float64, pressForDuration ...float64) (err error) {
	if len(pressForDuration) == 0 {
		pressForDuration = []float64{1}
	}
	return drag(s.sessionURL, fromX, fromY, toX, toY, pressForDuration[0])
}

func (s *Session) Swipe(fromX, fromY, toX, toY int) (err error) {
	return drag(s.sessionURL, fromX, fromY, toX, toY, 0)
}

func (s *Session) SwipeFloat(fromX, fromY, toX, toY float64) (err error) {
	return drag(s.sessionURL, fromX, fromY, toX, toY, 0)
}

func (s *Session) SwipeCoordinate(fromCoordinate, toCoordinate WDACoordinate) (err error) {
	return drag(s.sessionURL, fromCoordinate.X, fromCoordinate.Y, toCoordinate.X, toCoordinate.Y, 0)
}

// SwipeUp
func (s *Session) SwipeUp() (err error) {
	var fromCoordinate, toCoordinate WDACoordinate
	if windowSize, err := s.WindowSize(); err != nil {
		return err
	} else {
		center := WDACoordinate{X: windowSize.Width / 2, Y: windowSize.Height / 2}
		fromCoordinate, toCoordinate = center, center
	}
	fromCoordinate.Y += 100
	toCoordinate.Y -= 100
	return s.SwipeCoordinate(fromCoordinate, toCoordinate)
}

// SwipeDown
func (s *Session) SwipeDown() (err error) {
	var fromCoordinate, toCoordinate WDACoordinate
	if windowSize, err := s.WindowSize(); err != nil {
		return err
	} else {
		center := WDACoordinate{X: windowSize.Width / 2, Y: windowSize.Height / 2}
		fromCoordinate, toCoordinate = center, center
	}
	fromCoordinate.Y -= 100
	toCoordinate.Y += 100
	return s.SwipeCoordinate(fromCoordinate, toCoordinate)
}

// SwipeLeft
func (s *Session) SwipeLeft() (err error) {
	var fromCoordinate, toCoordinate WDACoordinate
	if windowSize, err := s.WindowSize(); err != nil {
		return err
	} else {
		center := WDACoordinate{X: windowSize.Width / 2, Y: windowSize.Height / 2}
		fromCoordinate, toCoordinate = center, center
	}
	fromCoordinate.X += 100
	toCoordinate.X -= 100
	return s.SwipeCoordinate(fromCoordinate, toCoordinate)
}

// SwipeRight
func (s *Session) SwipeRight() (err error) {
	var fromCoordinate, toCoordinate WDACoordinate
	if windowSize, err := s.WindowSize(); err != nil {
		return err
	} else {
		center := WDACoordinate{X: windowSize.Width / 2, Y: windowSize.Height / 2}
		fromCoordinate, toCoordinate = center, center
	}
	fromCoordinate.X -= 100
	toCoordinate.X += 100
	return s.SwipeCoordinate(fromCoordinate, toCoordinate)
}

// AppTerminate
//
// Close the application by bundleId
//
//	1. unregisterApplicationWithBundleId
func (s *Session) AppTerminate(bundleId string) (err error) {
	body := newWdaBody().setBundleID(bundleId)
	_, err = internalPost("AppTerminate", urlJoin(s.sessionURL, "/wda/apps/terminate"), body)
	// "value" : true,
	// "value" : false,
	return
}

type WDAAppRunState int

const (
	WDAAppNotRunning WDAAppRunState = 1 << iota
	WDAAppRunningBack
	WDAAppRunningFront
)

func (v WDAAppRunState) String() string {
	switch v {
	case WDAAppNotRunning:
		return "Not Running"
	case WDAAppRunningBack:
		return "Running (Back)"
	case WDAAppRunningFront:
		return "Running (Front)"
	default:
		return "UNKNOWN"
	}
}

// AppState
//
// Get the state of the particular application in scope of the current session.
// !This method is only returning reliable results since Xcode9 SDK
func (s *Session) AppState(bundleId string) (appRunState WDAAppRunState, err error) {
	body := newWdaBody().setBundleID(bundleId)
	var wdaResp wdaResponse
	if wdaResp, err = internalPost("AppState", urlJoin(s.sessionURL, "/wda/apps/state"), body); err != nil {
		return -1, err
	}
	return WDAAppRunState(wdaResp.getValue().Int()), nil
}

const WDATextBackspaceDeleteSequence = "\u0008\u007F"

// sendKeys
func sendKeys(url string, text string, typingFrequency ...int) (err error) {
	body := newWdaBody().set("value", strings.Split(text, ""))
	if len(typingFrequency) != 0 {
		body.set("frequency", typingFrequency[0])
	}
	_, err = internalPost("SendKeys", url, body)
	return
}

// SendKeys
//
// static NSUInteger FBMaxTypingFrequency = 60;
func (s *Session) SendKeys(text string, typingFrequency ...int) error {
	return sendKeys(urlJoin(s.sessionURL, "/wda/keys"), text, typingFrequency...)
}

func findUidOfElement(baseUrl *url.URL, wdaLocator WDALocator) (elemUID string, err error) {
	using, value := wdaLocator.getUsingAndValue()
	body := newWdaBody().set("using", using).set("value", value)
	var wdaResp wdaResponse
	if wdaResp, err = internalPost("FindElement", urlJoin(baseUrl, "/element"), body); err != nil {
		return "", err
	}
	return wdaResp.getValue().Get("ELEMENT").String(), nil
}

// FindElement
func (s *Session) FindElement(wdaLocator WDALocator) (element *Element, err error) {
	var elemUID string
	if elemUID, err = findUidOfElement(s.sessionURL, wdaLocator); err != nil {
		return nil, err
	}
	return newElement(s.sessionURL, elemUID), nil
}

func findUidOfElements(baseUrl *url.URL, wdaLocator WDALocator) (elemUIDs []string, err error) {
	using, value := wdaLocator.getUsingAndValue()
	body := newWdaBody().set("using", using).set("value", value)
	var wdaResp wdaResponse
	if wdaResp, err = internalPost("FindElements", urlJoin(baseUrl, "/elements"), body); err != nil {
		return nil, err
	}
	results := wdaResp.getValue().Array()
	if len(results) == 0 {
		return nil, errors.New(fmt.Sprintf("no such element: unable to find an element using '%s', value '%s'", using, value))
	}
	elemUIDs = make([]string, len(results))
	for i := range elemUIDs {
		elemUIDs[i] = results[i].Get("ELEMENT").String()
	}
	return
}

// FindElements
func (s *Session) FindElements(wdaLocator WDALocator) (elements []*Element, err error) {
	var elemUIDs []string
	if elemUIDs, err = findUidOfElements(s.sessionURL, wdaLocator); err != nil {
		return nil, err
	}
	elements = make([]*Element, len(elemUIDs))
	for i := range elements {
		elements[i] = newElement(s.sessionURL, elemUIDs[i])
	}
	return
}

// ActiveElement
//
// returns the currently active element
//
// [NSPredicate predicateWithFormat:@"hasKeyboardFocus == YES"]
func (s *Session) ActiveElement() (element *Element, err error) {
	var wdaResp wdaResponse
	if wdaResp, err = internalGet("ActiveElement", urlJoin(s.sessionURL, "/element/active")); err != nil {
		return nil, err
	}
	element = newElement(s.sessionURL, wdaResp.getValue().Get("ELEMENT").String())
	return
}

// IsLocked
//
// Checks if the screen is locked or not.
func (s *Session) IsLocked() (bool, error) {
	return isLocked(s.sessionURL)
}

// Unlock
//
// Forces the device under test to unlock.
// An immediate return will happen if the device is already unlocked and an error is going to be thrown if the screen has not been unlocked after the timeout.
func (s *Session) Unlock() (err error) {
	return unlock(s.sessionURL)
}

// Lock
//
// Forces the device under test to switch to the lock screen.
// An immediate return will happen if the device is already locked and an error is going to be thrown if the screen has not been locked after the timeout.
func (s *Session) Lock() (err error) {
	return lock(s.sessionURL)
}

// AppActivate
//
// Activate the application by restoring it from the background.
// Nothing will happen if the application is already in foreground.
// This method is only supported since Xcode9.
func (s *Session) AppActivate(bundleId string) (err error) {
	body := newWdaBody().setBundleID(bundleId)
	_, err = internalPost("AppActivate", urlJoin(s.sessionURL, "/wda/apps/activate"), body)
	return
}

// AppDeactivate
//
// Deactivates application for given time and then activate it again
func (s *Session) AppDeactivate(seconds ...float64) (err error) {
	body := newWdaBody()
	if len(seconds) != 0 {
		body.set("duration", seconds[0])
	}
	wdaResp, err := internalPost("AppDeactivate", urlJoin(s.sessionURL, "/wda/deactivateApp"), body)
	if err != nil {
		return err
	}
	return wdaResp.getErrMsg()
}

type WDAContentType string

const (
	WDAContentTypePlaintext WDAContentType = "plaintext"
	WDAContentTypeImage     WDAContentType = "image"
	WDAContentTypeUrl       WDAContentType = "url"
)

// SetPasteboardForType
func (s *Session) SetPasteboardForPlaintext(content string) (err error) {
	encodedContent := base64.StdEncoding.EncodeToString([]byte(content))
	return s.SetPasteboard(WDAContentTypePlaintext, encodedContent)
}

// SetPasteboardForImage
func (s *Session) SetPasteboardForImage(filename string) (err error) {
	imgFile, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer imgFile.Close()
	all, err := ioutil.ReadAll(imgFile)
	if err != nil {
		return err
	}
	encodedContent := base64.StdEncoding.EncodeToString(all)
	return s.SetPasteboard(WDAContentTypeImage, encodedContent)
}

// SetPasteboardForUrl
func (s *Session) SetPasteboardForUrl(url string) (err error) {
	encodedContent := base64.URLEncoding.EncodeToString([]byte(url))
	return s.SetPasteboard(WDAContentTypeUrl, encodedContent)
}

// SetPasteboard Sets data to the general pasteboard
func (s *Session) SetPasteboard(contentType WDAContentType, encodedContent string) (err error) {
	body := newWdaBody()
	body.set("contentType", contentType)
	body.set("content", encodedContent)

	_, err = internalPost("SetPasteboard", urlJoin(s.sessionURL, "/wda/setPasteboard"), body)
	return
}

type WDADeviceButtonName string

const (
	WDADeviceButtonHome       WDADeviceButtonName = "home"
	WDADeviceButtonVolumeUp   WDADeviceButtonName = "volumeUp"
	WDADeviceButtonVolumeDown WDADeviceButtonName = "volumeDown"
)

func (s *Session) PressHomeButton() (err error) {
	return s.PressButton(WDADeviceButtonHome)
}

func (s *Session) PressVolumeUpButton() (err error) {
	return s.PressButton(WDADeviceButtonVolumeUp)
}

func (s *Session) PressVolumeDownButton() (err error) {
	return s.PressButton(WDADeviceButtonVolumeDown)
}

// PressButton Presses the corresponding hardware button on the device
//
// !!! not a synchronous action
func (s *Session) PressButton(wdaDeviceButton WDADeviceButtonName) (err error) {
	body := newWdaBody().set("name", wdaDeviceButton)
	_, err = internalPost("PressButton", urlJoin(s.sessionURL, "/wda/pressButton"), body)
	return
}

// SiriActivate
//
// Activates Siri service voice recognition with the given text to parse
func (s *Session) SiriActivate(text string) (err error) {
	body := newWdaBody().set("text", text)
	_, err = internalPost("SiriActivate", urlJoin(s.sessionURL, "/wda/siri/activate"), body)
	return
}

// SiriOpenURL Open {%@}
// It doesn't actually work, right?
func (s *Session) SiriOpenURL(url string) (err error) {
	body := newWdaBody().set("url", url)
	_, err = internalPost("SiriOpenURL", urlJoin(s.sessionURL, "/url"), body)
	return
}

// Screenshot
//
// OR takes a screenshot of the specified element
func (s *Session) Screenshot(elemUID ...string) (raw *bytes.Buffer, err error) {
	return screenshot(s.sessionURL, elemUID...)
}

// ScreenshotToDisk
func (s *Session) ScreenshotToDisk(filename string, elemUID ...string) (err error) {
	return screenshotToDisk(s.sessionURL, filename, elemUID...)
}

// ScreenshotToImage
func (s *Session) ScreenshotToImage(elemUID ...string) (img image.Image, format string, err error) {
	return screenshotToImage(s.sessionURL, elemUID...)
}

// Source
func (s *Session) Source(srcOpt ...WDASourceOption) (sTree string, err error) {
	return source(s.sessionURL, srcOpt...)
}

// AccessibleSource
//
// Return application elements accessibility tree
//
// ignore all elements except for the main window for accessibility tree
func (s *Session) AccessibleSource() (sJson string, err error) {
	return accessibleSource(s.sessionURL)
}

func (s *Session) GetAppiumSettings() (sJson string, err error) {
	var wdaResp wdaResponse
	if wdaResp, err = internalGet("GetAppiumSettings", urlJoin(s.sessionURL, "/appium/settings")); err != nil {
		return "", err
	}
	return wdaResp.getValue().String(), nil
}

func (s *Session) SetAppiumSetting(key string, value interface{}) (sJson string, err error) {
	return s.SetAppiumSettings(map[string]interface{}{key: value})
}

func (s *Session) SetAppiumSettings(settings map[string]interface{}) (sJson string, err error) {
	body := newWdaBody().set("settings", settings)
	var wdaResp wdaResponse
	if wdaResp, err = internalPost("SetAppiumSettings", urlJoin(s.sessionURL, "/appium/settings"), body); err != nil {
		return "", err
	}
	return wdaResp.getValue().String(), nil
}

// It's not working
// /timeouts
// /wda/keyboard/dismiss
// /wda/getPasteboard

func (s *Session) tttTmp() {
	body := newWdaBody()
	_ = body

	// [NSPredicate predicateWithFormat:@"hasKeyboardFocus == YES"]
	wdaResp, err := internalGet("###############", urlJoin(s.sessionURL, "/element/active"))
	fmt.Println(err, wdaResp)
}
