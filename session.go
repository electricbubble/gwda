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

// Launch
func (s *Session) Launch(bundleId string) error {
	body := newWdaBody().setBundleID(bundleId).set("shouldWaitForQuiescence", false)
	wdaResp, err := internalPost("Launch", urlJoin(s.sessionURL, "wda", "apps", "launch"), body)
	if err != nil {
		return err
	}
	return wdaResp.getErrMsg()
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
	wdaResp, err := internalGet("DeviceInfo", urlJoin(s.sessionURL, "wda", "device", "info"))
	if err != nil {
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
	wdaResp, err := internalGet("BatteryInfo", urlJoin(s.sessionURL, "wda", "batteryInfo"))
	if err != nil {
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
	wdaResp, err := internalGet("WindowSize", urlJoin(s.sessionURL, "window", "size"))
	if err != nil {
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
	wdaResp, err := internalGet("Screen", urlJoin(s.sessionURL, "wda", "screen"))
	if err != nil {
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
	// wdaResp, err := internalGet("Scale", urlJoin(s.sessionURL, "wda", "screen"))
	// if err != nil {
	// 	return 0, err
	// }
	// return wdaResp.getValue().Get("scale").Float(), nil
	screen, err := s.Screen()
	return screen.Scale, err
}

// StatusBarSize
//
// {
//    "width": 375,
//    "height": 44
// }
func (s *Session) StatusBarSize() (wdaStatusBarSize WDASize, err error) {
	// wdaResp, err := internalGet("StatusBarSize", urlJoin(s.sessionURL, "wda", "screen"))
	// if err != nil {
	// 	return "", err
	// }
	// return wdaResp.getValue().Get("statusBarSize").String(), nil
	screen, err := s.Screen()
	return screen.StatusBarSize, err
}

// ActiveAppInfo
//
// {
//    "processArguments": {
//        "env": {
//        },
//        "args": [
//        ]
//    },
//    "name": "",
//    "pid": 57,
//    "bundleId": "com.apple.springboard"
// }
func (s *Session) ActiveAppInfo() (wdaActiveAppInfo WDAActiveAppInfo, err error) {
	wdaResp, err := internalGet("ActiveAppInfo", urlJoin(s.sessionURL, "wda", "activeAppInfo"))
	if err != nil {
		return WDAActiveAppInfo{}, err
	}
	wdaActiveAppInfo._String = wdaResp.getValue().String()
	err = json.Unmarshal([]byte(wdaActiveAppInfo._String), &wdaActiveAppInfo)
	// err = json.Unmarshal(wdaResp.getValue2Bytes(), &wdaActiveAppInfo)
	return
}

// Tap
func (s *Session) Tap(x, y int) error {
	body := newWdaBody().setXY(x, y)
	wdaResp, err := internalPost("Tap", urlJoin(s.sessionURL, "wda", "tap", "0"), body)
	if err != nil {
		return err
	}
	return wdaResp.getErrMsg()
}

// DoubleTap
func (s *Session) DoubleTap(x, y int) error {
	body := newWdaBody().setXY(x, y)
	wdaResp, err := internalPost("DoubleTap", urlJoin(s.sessionURL, "wda", "doubleTap"), body)
	if err != nil {
		return err
	}
	return wdaResp.getErrMsg()
}

// TouchAndHold
func (s *Session) TouchAndHold(x, y int, duration ...float32) error {
	body := newWdaBody().setXY(x, y)
	if len(duration) == 0 {
		body.set("duration", 1.0)
	} else {
		body.set("duration", duration[0])
	}
	wdaResp, err := internalPost("TouchAndHold", urlJoin(s.sessionURL, "wda", "touchAndHold"), body)
	if err != nil {
		return err
	}
	return wdaResp.getErrMsg()
}

// AppTerminate Close the application by bundleId
func (s *Session) AppTerminate(bundleId string) error {
	body := newWdaBody().setBundleID(bundleId)
	wdaResp, err := internalPost("AppTerminate", urlJoin(s.sessionURL, "wda", "apps", "terminate"), body)
	if err != nil {
		return err
	}
	// "value" : true,
	// "value" : false,
	return wdaResp.getErrMsg()
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

// Locked
func (s *Session) Locked() (isLocked bool, err error) {
	wdaResp, err := internalGet("Locked", urlJoin(s.sessionURL, "wda", "locked"))
	if err != nil {
		return false, err
	}
	return wdaResp.getValue().Bool(), nil
}

// Unlock
func (s *Session) Unlock() (err error) {
	wdaResp, err := internalPost("Unlock", urlJoin(s.sessionURL, "wda", "unlock"), nil)
	if err != nil {
		return err
	}
	return wdaResp.getErrMsg()
}

// Lock
func (s *Session) Lock() (err error) {
	wdaResp, err := internalPost("Lock", urlJoin(s.sessionURL, "wda", "lock"), nil)
	if err != nil {
		return err
	}
	return wdaResp.getErrMsg()
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

const (
	WDAContentTypePlaintext = "plaintext"
	WDAContentTypeImage     = "image"
	WDAContentTypeUrl       = "url"
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

// SetPasteboard
func (s *Session) SetPasteboard(contentType, encodedContent string) (err error) {
	body := newWdaBody()
	body.set("contentType", contentType)
	body.set("content", encodedContent)

	wdaResp, err := internalPost("SetPasteboard", urlJoin(s.sessionURL, "/wda/setPasteboard"), body)
	if err != nil {
		return err
	}
	return wdaResp.getErrMsg()
}

// Source
//
// Source aka tree
func (s *Session) Source(formattedAsJson ...bool) (sTree string, err error) {
	tmp, _ := url.Parse(s.sessionURL.String())
	if len(formattedAsJson) != 0 && formattedAsJson[0] {
		q := tmp.Query()
		q.Set("format", "json")
		tmp.RawQuery = q.Encode()
	}
	wdaResp, err := internalGet("Source", urlJoin(tmp, "source"))
	if err != nil {
		return "", err
	}
	return wdaResp.getValue().String(), err
}

// AccessibleSource
func (s *Session) AccessibleSource() (sJson string, err error) {
	wdaResp, err := internalGet("AccessibleSource", urlJoin(s.sessionURL, "wda", "accessibleSource"))
	if err != nil {
		return "", err
	}
	return wdaResp.getValue().String(), err
}

// TODO DELETE	/session/{session id}
// TODO /timeouts
// TODO /screenshot
// TODO /wda/deactivateApp
// TODO /wda/setPasteboard
// TODO /wda/getPasteboard
// TODO /wda/pressButton
// TODO /wda/siri/activate
// TODO /wda/apps/launchUnattachedlaunch
// TODO /wda/keyboard/dismiss 失效？
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
	// body.set("duration", 7.1)
	_ = body
	body.set("contentType", "plaintext")
	content := "abcd123"
	vContent := base64.StdEncoding.EncodeToString([]byte(content))
	body.set("content", vContent)

	body = newWdaBody()
	body.set("contentType", "image")
	open, _ := os.Open("/Users/hero/Documents/leixipaopao/IMG_5246.JPG")
	all, _ := ioutil.ReadAll(open)
	vContent = base64.StdEncoding.EncodeToString(all)
	body.set("content", vContent)

	body = newWdaBody()
	body.set("contentType", "url")
	vContent = base64.URLEncoding.EncodeToString([]byte("http://baidu.com"))
	body.set("content", vContent)

	wdaResp, err := internalPost("#TEMP", urlJoin(s.sessionURL, "/wda/setPasteboard"), body)
	// body["bundleId"] = bundleId
	// bsJson, err := s.AppState("com.netease.cloudmusic")
	fmt.Println(err, wdaResp)
}
