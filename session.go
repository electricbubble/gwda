package gwda

import (
	"errors"
	"fmt"
	"net/url"
)

type Session struct {
	sessionURL *url.URL
	// sessionID  string
	// bundleID   string
}

// Launch
func (s *Session) Launch(bundleId string) error {
	body := newWdaBody().setBundleID(bundleId)
	wdaResp, err := internalPost("Launch", urlJoin(s.sessionURL, "wda", "apps", "launch"), body)
	if err != nil {
		return err
	}
	return wdaResp.getErrMsg()
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
//    "name": "hero’s iPhone X"
// }
func (s *Session) DeviceInfo() (sJson string, err error) {
	wdaResp, err := internalGet("DeviceInfo", urlJoin(s.sessionURL, "wda", "device", "info"))
	if err != nil {
		return "", err
	}
	return wdaResp.getValue().String(), nil
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
func (s *Session) BatteryInfo() (sJson string, err error) {
	wdaResp, err := internalGet("电池信息", urlJoin(s.sessionURL, "wda", "batteryInfo"))
	if err != nil {
		return "", err
	}
	return wdaResp.getValue().String(), nil
}

// WindowSize
//
// {
//    "width": 812,
//    "height": 375
// }
func (s *Session) WindowSize() (sJson string, err error) {
	wdaResp, err := internalGet("WindowSize", urlJoin(s.sessionURL, "window", "size"))
	if err != nil {
		return "", err
	}
	return wdaResp.getValue().String(), nil
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
func (s *Session) Screen() (sJson string, err error) {
	wdaResp, err := internalGet("Screen", urlJoin(s.sessionURL, "wda", "screen"))
	if err != nil {
		return "", err
	}
	return wdaResp.getValue().String(), nil
}

// Scale
func (s *Session) Scale() (scale float64, err error) {
	wdaResp, err := internalGet("Scale", urlJoin(s.sessionURL, "wda", "screen"))
	if err != nil {
		return 0, err
	}
	return wdaResp.getValue().Get("scale").Float(), nil
}

// StatusBarSize
//
// {
//    "width": 375,
//    "height": 44
// }
func (s *Session) StatusBarSize() (sJson string, err error) {
	wdaResp, err := internalGet("StatusBarSize", urlJoin(s.sessionURL, "wda", "screen"))
	if err != nil {
		return "", err
	}
	return wdaResp.getValue().Get("statusBarSize").String(), nil
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
func (s *Session) ActiveAppInfo() (sJson string, err error) {
	wdaResp, err := internalGet("ActiveAppInfo", urlJoin(s.sessionURL, "wda", "activeAppInfo"))
	if err != nil {
		return "", err
	}
	return wdaResp.getValue().String(), nil
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

// AppState
//
// 1 未运行？
// 2 运行中（后台活动）
// 4 运行中（前台活动）
func (s *Session) AppState(bundleId string) (state int, err error) {
	body := newWdaBody().setBundleID(bundleId)
	wdaResp, err := internalPost("AppState", urlJoin(s.sessionURL, "wda", "apps", "state"), body)
	if err != nil {
		return -1, nil
	}
	return int(wdaResp.getValue().Int()), nil
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
		return
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
		return Element{}, errors.New("no such element")
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
	wdaResp, err := internalGet("判断锁屏界面", urlJoin(s.sessionURL, "wda", "locked"))
	if err != nil {
		return false, err
	}
	return wdaResp.getValue().Bool(), nil
}

// Unlock
func (s *Session) Unlock() (err error) {
	wdaResp, err := internalPost("触发解锁", urlJoin(s.sessionURL, "wda", "unlock"), nil)
	if err != nil {
		return err
	}
	return wdaResp.getErrMsg()
}

// Lock
func (s *Session) Lock() (err error) {
	wdaResp, err := internalPost("触发锁屏", urlJoin(s.sessionURL, "wda", "lock"), nil)
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
	wdaResp, err := internalGet("Source aka tree", urlJoin(tmp, "source"))
	if err != nil {
		return "", err
	}
	return wdaResp.getValue().String(), err
}

// AccessibleSource
func (s *Session) AccessibleSource() (sJson string, err error) {
	wdaResp, err := internalGet("Source aka tree", urlJoin(s.sessionURL, "wda", "accessibleSource"))
	if err != nil {
		return "", err
	}
	return wdaResp.getValue().String(), err
}

// TODO DELETE	/session/{session id}
// TODO /timeouts
// TODO /screenshot
// TODO /wda/apps/activate
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

	body := make(map[string]interface{})
	// body["bundleId"] = "com.netease.cloudmusic"
	// body["url"] = "baidu.com"
	// body["shouldWaitForQuiescence"] = true
	// body["x"] = 230
	// body["y"] = 130
	// body["duration"] = 1.0
	// body["value"] = strings.Split("中文测试1.0", "")
	// body["using"] = "link text"
	// body["value"] = "label=发现"
	body["url"] = "http://www.baidu.com"
	// bsJson, err := internalPost("tttTmp", urlJoin(s.sessionURL, "/wda/keyboard/dismiss"), nil)
	// bsJson, err := internalPost("tttTmp", urlJoin(s.sessionURL, "/elements"), body)
	// bsJson, err := internalPost("tttTmp", urlJoin(s.sessionURL, "/url"), body)
	// bsJson, err := internalGet("tttTmp", urlJoin(s.sessionURL, "/window/size"))
	bsJson, err := internalGet("tttTmp", urlJoin(s.sessionURL, "/wda/apps/list"))
	// body["bundleId"] = bundleId
	// bsJson, err := s.AppState("com.netease.cloudmusic")
	fmt.Println(err, string(bsJson))
}
