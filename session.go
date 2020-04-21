package gwda

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
)

type Session struct {
	sessionURL *url.URL
	// sessionID  string
	// bundleID   string
}

// Delete kill session associated with that request
//
// 1. alertsMonitor disable
// 2. testedApplicationBundleId terminate
func (s *Session) Delete() (err error) {
	_, err = internalDelete("DeleteSession", s.sessionURL.String())
	return
}

// AppLaunch always wait for quiescence
//
// 1. registerApplicationWithBundleId
// 2. launch OR activate
func (s *Session) AppLaunch(bundleId string) (err error) {
	// TODO
	//  arguments
	//  environment
	body := newWdaBody().setBundleID(bundleId).set("shouldWaitForQuiescence", true)
	_, err = internalPost("Launch", urlJoin(s.sessionURL, "wda", "apps", "launch"), body)
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
	_String            string
}

func (di WDADeviceInfo) String() string {
	return di._String
}

// DeviceInfo
//
// {
//    "timeZone": "GMT+0800",
//    "currentLocale": "zh_CN",
//    "model": "iPhone",
//    "uuid": "x-x-x-x-x",
//    "userInterfaceIdiom": 0,
//    "isSimulator": false,
//    "name": "x’s iPhone X"
// }
func (s *Session) DeviceInfo() (wdaDeviceInfo WDADeviceInfo, err error) {
	var wdaResp wdaResponse
	if wdaResp, err = internalGet("DeviceInfo", urlJoin(s.sessionURL, "wda", "device", "info")); err != nil {
		return
	}

	wdaDeviceInfo._String = wdaResp.getValue().String()
	err = json.Unmarshal([]byte(wdaDeviceInfo._String), &wdaDeviceInfo)
	// err = json.Unmarshal(wdaResp.getValue2Bytes(), &wdaDeviceInfo)
	return
}

type WDABatteryInfo struct {
	Level   float32         `json:"level"` // Battery level in range [0.0, 1.0], where 1.0 means 100% charge.
	State   WDABatteryState `json:"state"` // Battery state ( 1: on battery, discharging; 2: plugged in, less than 100%, 3: plugged in, at 100% )
	_String string
}

func (bi WDABatteryInfo) String() string {
	return bi._String
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
	if wdaResp, err = internalGet("BatteryInfo", urlJoin(s.sessionURL, "wda", "batteryInfo")); err != nil {
		return
	}

	wdaBatteryInfo._String = wdaResp.getValue().String()
	err = json.Unmarshal([]byte(wdaBatteryInfo._String), &wdaBatteryInfo)
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
	if wdaResp, err = internalGet("WindowSize", urlJoin(s.sessionURL, "window", "size")); err != nil {
		return
	}

	wdaSize._String = wdaResp.getValue().String()
	err = json.Unmarshal([]byte(wdaSize._String), &wdaSize)
	// err = json.Unmarshal(wdaResp.getValue2Bytes(), &wdaSize)
	return
}

type WDASize struct {
	Width   int `json:"width"`
	Height  int `json:"height"`
	_String string
}

func (s WDASize) String() string {
	return s._String
}

type WDAScreen struct {
	StatusBarSize WDASize `json:"statusBarSize"`
	Scale         float32 `json:"scale"`
	_String       string
}

func (s WDAScreen) String() string {
	return s._String
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
	if wdaResp, err = internalGet("Screen", urlJoin(s.sessionURL, "wda", "screen")); err != nil {
		return
	}

	wdaScreen.StatusBarSize._String = wdaResp.getValue().Get("statusBarSize").String()
	wdaScreen._String = wdaResp.getValue().String()
	err = json.Unmarshal([]byte(wdaScreen._String), &wdaScreen)
	// err = json.Unmarshal(wdaResp.getValue2Bytes(), &wdaScreen)
	return
}

// Scale
func (s *Session) Scale() (scale float32, err error) {
	screen, err := s.Screen()
	return screen.Scale, err
}

// StatusBarSize
func (s *Session) StatusBarSize() (wdaStatusBarSize WDASize, err error) {
	screen, err := s.Screen()
	return screen.StatusBarSize, err
}

// ActiveAppInfo Constructor used to get current active application
func (s *Session) ActiveAppInfo() (wdaActiveAppInfo WDAActiveAppInfo, err error) {
	return activeAppInfo(s.sessionURL)
}

// Tap
// TODO tap
func (s *Session) Tap(x, y int) error {
	body := newWdaBody().setXY(x, y)
	wdaResp, err := internalPost("Tap", urlJoin(s.sessionURL, "wda", "tap", "0"), body)
	if err != nil {
		return err
	}
	return wdaResp.getErrMsg()
}

// DoubleTap double tap coordinate
func (s *Session) DoubleTap(x, y int) (err error) {
	body := newWdaBody().setXY(x, y)
	_, err = internalPost("DoubleTap", urlJoin(s.sessionURL, "wda", "doubleTap"), body)
	return
}

// TouchAndHold touch and hold coordinate
func (s *Session) TouchAndHold(x, y int, duration ...float32) (err error) {
	body := newWdaBody().setXY(x, y)
	if len(duration) == 0 {
		body.set("duration", 1.0)
	} else {
		body.set("duration", duration[0])
	}
	_, err = internalPost("TouchAndHold", urlJoin(s.sessionURL, "wda", "touchAndHold"), body)
	return
}

// AppTerminate Close the application by bundleId
//
// 1. unregisterApplicationWithBundleId
func (s *Session) AppTerminate(bundleId string) (err error) {
	body := newWdaBody().setBundleID(bundleId)
	_, err = internalPost("AppTerminate", urlJoin(s.sessionURL, "wda", "apps", "terminate"), body)
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
// 1 未运行？
// 2 运行中（后台活动）
// 4 运行中（前台活动）
func (s *Session) AppState(bundleId string) (appRunState WDAAppRunState, err error) {
	body := newWdaBody().setBundleID(bundleId)
	wdaResp, err := internalPost("AppState", urlJoin(s.sessionURL, "wda", "apps", "state"), body)
	if err != nil {
		return -1, err
	}
	return WDAAppRunState(wdaResp.getValue().Int()), nil
}

// SendKeys
// TODO 每个字符输入等待时间 5s, 输入失败不会报错
// TODO frequency
// `确定` 使用 `\n`
func (s *Session) SendKeys(text string) error {
	body := newWdaBody().setSendKeys(text)
	wdaResp, err := internalPost("SendKeys", urlJoin(s.sessionURL, "/wda/keys"), body)
	if err != nil {
		return err
	}
	return wdaResp.getErrMsg()
}

var ErrNoSuchElement = errors.New("no such element")

// TODO FindElements
func (s *Session) FindElements(using, value string) (elements []Element, err error) {
	body := newWdaBody().set("using", using).set("value", value)
	wdaResp, err := internalPost("FindElements", urlJoin(s.sessionURL, "elements"), body)
	if err != nil {
		return nil, err
	}
	// fmt.Println(wdaResp)
	// fmt.Println(wdaResp.getValue().IsArray())
	results := wdaResp.getValue().Array()
	if len(results) == 0 {
		return nil, ErrNoSuchElement
	}
	elements = make([]Element, 0, len(results))
	for _, res := range results {
		elementId := res.Get("ELEMENT").String()
		ele := Element{}
		ele.elementURL, _ = url.Parse(urlJoin(s.sessionURL, "element", elementId))
		elements = append(elements, ele)
		// fmt.Println("###", elementId)
	}
	return
}

// TODO FindElement
func (s *Session) FindElement(using, value string) (element Element, err error) {
	body := newWdaBody().set("using", using).set("value", value)
	wdaResp, err := internalPost("FindElements", urlJoin(s.sessionURL, "element"), body)
	if err != nil {
		return Element{}, err
	}
	// fmt.Println(wdaResp)
	// fmt.Println(wdaResp.getValue().IsArray())
	elementID := wdaResp.getValue().Get("ELEMENT").String()
	if elementID == "" {
		return Element{}, ErrNoSuchElement
	}
	element = Element{}
	element.elementURL, _ = url.Parse(urlJoin(s.sessionURL, "element", elementID))
	// element = make([]Element, 0, len(elementID))
	// for _, res := range elementID {
	// 	elementId := res.Get("ELEMENT").String()
	// 	ele := Element{}
	// 	ele.elementURL, _ = url.Parse(urlJoin(s.sessionURL, "element", elementId))
	// 	element = append(element, ele)
	// 	// fmt.Println("###", elementId)
	// }
	return
}

// IsLocked Whether the screen is locked
func (s *Session) IsLocked() (bool, error) {
	return isLocked(s.sessionURL)
}

// Unlock unlock screen
func (s *Session) Unlock() (err error) {
	return unlock(s.sessionURL)
}

// Lock
func (s *Session) Lock() (err error) {
	return lock(s.sessionURL)
}

// TODO /wda/apps/activate
func (s *Session) Activate(bundleId string) error {
	body := newWdaBody().setBundleID(bundleId)
	wdaResp, err := internalPost("Activate", urlJoin(s.sessionURL, "wda", "apps", "activate"), body)
	if err != nil {
		return err
	}
	return wdaResp.getErrMsg()
}

// DeactivateApp Deactivates application for given time
func (s *Session) DeactivateApp(seconds ...float32) (err error) {
	body := newWdaBody()
	if len(seconds) != 0 {
		body.set("seconds", seconds[0])
	}
	wdaResp, err := internalPost("DeactivateApp", urlJoin(s.sessionURL, "/wda/deactivateApp"), body)
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

// SiriActivate Activates Siri service voice recognition with the given text to parse
func (s *Session) SiriActivate(text string) (err error) {
	body := newWdaBody().set("text", text)
	_, err = internalPost("SiriActivate", urlJoin(s.sessionURL, "/wda/siri/activate"), body)
	return
}

// SiriOpenURL Open {%@}
// It doesn't actually work, right?
func (s *Session) SiriOpenURL(url string) (err error) {
	body := newWdaBody().set("url", url)
	_, err = internalPost("OpenURL", urlJoin(s.sessionURL, "/url"), body)
	return
}

// Source
//
// Source aka tree
func (s *Session) Source(formattedAsJson ...bool) (sTree string, err error) {
	return source(s.sessionURL, formattedAsJson...)
}

// AccessibleSource
func (s *Session) AccessibleSource() (sJson string, err error) {
	return accessibleSource(s.sessionURL)
}

// It's not working
// /timeouts
// /wda/keyboard/dismiss
// /wda/getPasteboard

// TODO DELETE	/session/{session id}
// TODO /screenshot
// TODO wdaResp, err := internalGet("AppList", urlJoin(s.sessionURL, "/wda/apps/list", ))	handleGetActiveAppsList	fb_activeAppsInfo

func (s *Session) tttTmp() {

	// wdaResp, err := internalGet("AppList", urlJoin(s.sessionURL, "/wda/apps/list", ))
	// fmt.Println(err, wdaResp)
	// 
	// return

	// body := make(map[string]interface{})
	// body["bundleId"] = "com.netease.cloudmusic"
	// body["url"] = "baidu.com"
	// body["shouldWaitForQuiescence"] = true
	// body["x"] = 230
	// body["y"] = 130
	// body["duration"] = 1.0
	// body["value"] = strings.Split("中文测试1.0", "")
	// body["using"] = "link text"
	// body["value"] = "label=发现"
	// body["url"] = "http://www.baidu.com"
	// bsJson, err := internalPost("tttTmp", urlJoin(s.sessionURL, "/wda/keyboard/dismiss"), nil)
	// bsJson, err := internalPost("tttTmp", urlJoin(s.sessionURL, "/elements"), body)
	// bsJson, err := internalPost("tttTmp", urlJoin(s.sessionURL, "/url"), body)
	// bsJson, err := internalGet("tttTmp", urlJoin(s.sessionURL, "/window/size"))
	// bsJson, err := internalGet("tttTmp", urlJoin(s.sessionURL, "/wda/apps/list"))
	body := newWdaBody()
	_ = body

	// body.set("bundleId", "com.netease.cloudmusic")
	body.set("url", "baidu.com")
	wdaResp, err := internalPost("#TEMP", urlJoin(s.sessionURL, "/url"), body)
	// body["bundleId"] = bundleId
	// bsJson, err := s.AppState("com.netease.cloudmusic")
	fmt.Println(err, wdaResp)
}
