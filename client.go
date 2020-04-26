package gwda

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"net/url"
	"strings"
)

type Client struct {
	deviceURL *url.URL
}

func NewClient(deviceURL string) (c *Client, err error) {
	c = &Client{}
	if c.deviceURL, err = url.Parse(deviceURL); err != nil {
		return nil, err
	}
	sJson, err := c.Status()
	if err != nil {
		return nil, err
	}
	// ready: always true
	if !wdaResponse(sJson).getByPath("value.ready").Bool() {
		return nil, errors.New("device is not ready")
	}

	settings := newWdaBody().set("acceptAlertButtonSelector", _acceptAlertButtonSelector).set("dismissAlertButtonSelector", _dismissAlertButtonSelector)
	c.setAppiumSettings(settings, wdaResponse(sJson).getByPath("sessionId").String())
	return c, nil
}

func (c *Client) setAppiumSettings(settings map[string]interface{}, sid ...string) {
	if len(sid) == 0 {
		status, err := c.Status()
		if err != nil {
			log.Printf("[ERROR]↩︎\n[setAppiumSettings] failed to get status %s\n", err.Error())
			return
		}
		sid = []string{wdaResponse(status).getValue().String()}
	}
	if _, err := newSession(c.deviceURL, sid[0]).SetAppiumSettings(settings); err != nil {
		// TODO return err ?
		//  [settings objectForKey:ACTIVE_APP_DETECTION_POINT]
		//  [settings objectForKey:SCREENSHOT_ORIENTATION]
		log.Printf("[ERROR]↩︎\n[setAppiumSettings] failed to set AppiumSettings %s\n", err.Error())
	}
}

const _acceptAlertButtonSelector = "**/XCUIElementTypeButton[`label == '允许' OR label == '好' OR label == '仅在使用应用期间' OR label == '暂不'`]"
const _dismissAlertButtonSelector = "**/XCUIElementTypeButton[`label == '不允许' OR label == '暂不'`]"

// SetAcceptAlertButtonSelector
//
// Sets custom class chain locators for accept/dismiss alert buttons location.
//
// This might be useful if the default buttons detection algorithm fails to determine alert buttons properly when defaultAlertAction is set.
//
// Default: **/XCUIElementTypeButton[`label == '允许' OR label == '好' OR label == '仅在使用应用期间' OR label == '暂不'`]
func (c *Client) SetAcceptAlertButtonSelector(classChainSelector string) {
	c.setAppiumSettings(map[string]interface{}{"acceptAlertButtonSelector": classChainSelector})
}

// SetDismissAlertButtonSelector
//
// Default: **/XCUIElementTypeButton[`label == '不允许' OR label == '暂不'`]
func (c *Client) SetDismissAlertButtonSelector(classChainSelector string) {
	c.setAppiumSettings(map[string]interface{}{"dismissAlertButtonSelector": classChainSelector})
}

type WDASessionCapability wdaBody

// NewWDASessionCapability
//
// Default wait for quiescence
func NewWDASessionCapability(bundleId ...string) WDASessionCapability {
	sCapabilities := make(WDASessionCapability)
	if len(bundleId) != 0 {
		wdaBody(sCapabilities).setBundleID(bundleId[0])
		sCapabilities.SetAppLaunchOption(NewWDAAppLaunchOption().SetShouldWaitForQuiescence(true))
	}
	return sCapabilities
}

type WDASessionDefaultAlertAction string

const (
	WDASessionAlertActionAccept  WDASessionDefaultAlertAction = "accept"
	WDASessionAlertActionDismiss WDASessionDefaultAlertAction = "dismiss"
)

// SetDefaultAlertAction
//
// Creates and saves new session for application with default alert handling behaviour
//
// Default is disabled
func (sc WDASessionCapability) SetDefaultAlertAction(sAlertAction WDASessionDefaultAlertAction) WDASessionCapability {
	return WDASessionCapability(wdaBody(sc).set("defaultAlertAction", sAlertAction))
}

// SetAppLaunchOption
func (sc WDASessionCapability) SetAppLaunchOption(opt WDAAppLaunchOption) WDASessionCapability {
	return WDASessionCapability(wdaBody(sc).setAppLaunchOption(opt))
}

// SetShouldUseTestManagerForVisibilityDetection
//
// Default is `false`
// static BOOL FBShouldUseTestManagerForVisibilityDetection = NO;
func (sc WDASessionCapability) SetShouldUseTestManagerForVisibilityDetection(b bool) WDASessionCapability {
	return WDASessionCapability(wdaBody(sc).set("shouldUseTestManagerForVisibilityDetection", b))
}

// SetShouldUseCompactResponses
//
// Default is `true`
// static BOOL FBShouldUseCompactResponses = YES;
func (sc WDASessionCapability) SetShouldUseCompactResponses(b bool) WDASessionCapability {
	return WDASessionCapability(wdaBody(sc).set("shouldUseCompactResponses", b))
}

// SetElementResponseAttributes
//
// Default is `"type,label"`
// static NSString *FBElementResponseAttributes = @"type,label";
func (sc WDASessionCapability) SetElementResponseAttributes(s string) WDASessionCapability {
	return WDASessionCapability(wdaBody(sc).set("elementResponseAttributes", s))
}

// SetMaxTypingFrequency
//
// Default is `60`
// static NSUInteger FBMaxTypingFrequency = 60;
func (sc WDASessionCapability) SetMaxTypingFrequency(n int) WDASessionCapability {
	return WDASessionCapability(wdaBody(sc).set("maxTypingFrequency", n))
}

// SetShouldUseSingletonTestManager
//
// Default is `true`
// static BOOL FBShouldUseSingletonTestManager = YES;
func (sc WDASessionCapability) SetShouldUseSingletonTestManager(b bool) WDASessionCapability {
	return WDASessionCapability(wdaBody(sc).set("shouldUseSingletonTestManager", b))
}

// SetEventloopIdleDelaySec
//
// Once the methods were swizzled they stay like that since the only change in the implementation is the thread sleep, which is skipped on setting it to zero.
//
// <= 0 disableEventLoopDelay
//
// Default is `0`
// static NSTimeInterval eventloopIdleDelay = 0;
func (sc WDASessionCapability) SetEventloopIdleDelaySec(seconds int) WDASessionCapability {
	return WDASessionCapability(wdaBody(sc).set("eventloopIdleDelaySec", seconds))
}

// NewSession
//
// Creates and saves new session for application
func (c *Client) NewSession(capabilities ...WDASessionCapability) (s *Session, err error) {
	// TODO BundleId is required 如果是不存在的 bundleId 会导致 wda 内部报错导致接下来的操作都无法接收处理
	body := newWdaBody()
	if len(capabilities) != 0 {
		body.set("capabilities", newWdaBody().set("alwaysMatch", capabilities[0]))
	} else {
		body.set("capabilities", newWdaBody()) // .set("alwaysMatch", nil))
	}
	var wdaResp wdaResponse
	if wdaResp, err = internalPost("NewSession", urlJoin(c.deviceURL, "/session"), body); err != nil {
		return nil, err
	}
	if sid := wdaResp.getByPath("sessionId").String(); sid == "" {
		return nil, errors.New("not find sessionId")
	} else {
		// c.deviceURL 已在新建时校验过, 理论上此处不再出现错误
		s = newSession(c.deviceURL, sid)
	}
	return s, nil
}

// AppLaunchUnattached
//
// Launch the app with the specified bundle ID
//
// shouldWaitForQuiescence: false
func (c *Client) AppLaunchUnattached(bundleId string) (err error) {
	body := newWdaBody().setBundleID(bundleId)
	_, err = internalPost("AppLaunchUnattached", urlJoin(c.deviceURL, "/wda/apps/launchUnattached"), body)
	return
}

// Status
//
// Checking service status
func (c *Client) Status() (sJson string, err error) {
	var wdaResp wdaResponse
	if wdaResp, err = internalGet("Status", urlJoin(c.deviceURL, "/status")); err != nil {
		return "", err
	}
	return wdaResp.String(), nil
}

// Homescreen
//
// Forces the device under test to switch to the home screen
//
// 1. pressButton
// 2. WaitUntilApplicationBoardIsVisible
func (c *Client) Homescreen() (err error) {
	_, err = internalPost("Homescreen", urlJoin(c.deviceURL, "/wda/homescreen"), nil)
	return
}

// HealthCheck
//
// Checks health of XCTest by:
// 1) Querying application for some elements,
// 2) Triggering some device events.
//
// !!! Health check might modify simulator state so it should only be called in-between testing sessions
func (c *Client) HealthCheck() (err error) {
	_, err = internalGet("HealthCheck", urlJoin(c.deviceURL, "/wda/healthcheck"))
	return
}

func isLocked(baseUrl *url.URL) (isLocked bool, err error) {
	var wdaResp wdaResponse
	if wdaResp, err = internalGet("Locked", urlJoin(baseUrl, "/wda/locked")); err != nil {
		return false, err
	}
	return wdaResp.getValue().Bool(), nil
}

// IsLocked
//
// Checks if the screen is locked or not.
func (c *Client) IsLocked() (bool, error) {
	return isLocked(c.deviceURL)
}

func unlock(baseUrl *url.URL) (err error) {
	_, err = internalPost("Unlock", urlJoin(baseUrl, "/wda/unlock"), nil)
	return
}

// Unlock
//
// Forces the device under test to unlock.
// An immediate return will happen if the device is already unlocked and an error is going to be thrown if the screen has not been unlocked after the timeout.
func (c *Client) Unlock() (err error) {
	return unlock(c.deviceURL)
}

func lock(baseUrl *url.URL) (err error) {
	_, err = internalPost("Lock", urlJoin(baseUrl, "/wda/lock"), nil)
	return
}

// Lock
//
// Forces the device under test to switch to the lock screen.
// An immediate return will happen if the device is already locked and an error is going to be thrown if the screen has not been locked after the timeout.
func (c *Client) Lock() (err error) {
	return lock(c.deviceURL)
}

// screenshot
//
// [FBRoute GET:@"/screenshot"]					format: png
// [FBRoute GET:@"/element/:uuid/screenshot"]	format: jpeg
// [FBRoute GET:@"/screenshot/:uuid"]			format: jpeg
func screenshot(baseUrl *url.URL, elemUID ...string) (raw *bytes.Buffer, err error) {
	var wdaResp wdaResponse

	tmpPath := "/screenshot"
	if len(elemUID) != 0 && elemUID[0] != "" {
		tmpPath += "/" + elemUID[0]
	}

	if wdaResp, err = internalGet("Screenshot", urlJoin(baseUrl, tmpPath)); err != nil {
		return nil, err
	}

	if decodeString, err := base64.StdEncoding.DecodeString(wdaResp.getValue().String()); err != nil {
		return nil, err
	} else {
		raw = bytes.NewBuffer(decodeString)
		return raw, nil
	}
}

func screenshotToDisk(baseUrl *url.URL, filename string, elemUID ...string) (err error) {
	var raw *bytes.Buffer
	if raw, err = screenshot(baseUrl, elemUID...); err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, raw.Bytes(), 0666)
	return
}

func screenshotToPng(baseUrl *url.URL) (img image.Image, err error) {
	if raw, err := screenshot(baseUrl); err != nil {
		return nil, err
	} else {
		img, err = png.Decode(raw)
		return img, err
	}
}

func screenshotToJpeg(baseUrl *url.URL, elemUID string) (img image.Image, err error) {
	if raw, err := screenshot(baseUrl, elemUID); err != nil {
		return nil, err
	} else {
		img, err = jpeg.Decode(raw)
		return img, err
	}
}

// Screenshot
func (c *Client) Screenshot() (raw *bytes.Buffer, err error) {
	return screenshot(c.deviceURL)
}

// ScreenshotToDiskAsJpeg
func (c *Client) ScreenshotToDiskAsPng(filename string) (err error) {
	return screenshotToDisk(c.deviceURL, filename)
}

// ScreenshotToJpeg
func (c *Client) ScreenshotToPng() (img image.Image, err error) {
	return screenshotToPng(c.deviceURL)
}

type WDASourceOption wdaBody

// NewWDASourceOption
//
// Default: "format": "xml"
func NewWDASourceOption() WDASourceOption {
	return make(WDASourceOption)
}

func (so WDASourceOption) SetFormatAsJson() WDASourceOption {
	return WDASourceOption(wdaBody(so).set("format", "json"))
}

func (so WDASourceOption) SetFormatAsXml() WDASourceOption {
	return WDASourceOption(wdaBody(so).set("format", "xml"))
}

func (so WDASourceOption) SetFormatAsDescription() WDASourceOption {
	return WDASourceOption(wdaBody(so).set("format", "description"))
}

// SetExcludedAttributes
//
// only `xml` supported.
func (so WDASourceOption) SetExcludedAttributes(excludedAttributes []string) WDASourceOption {
	if vFormat, ok := so["format"]; ok && vFormat != "xml" {
		return so
	}
	return WDASourceOption(wdaBody(so).set("excluded_attributes", strings.Join(excludedAttributes, ",")))
}

// source
func source(baseUrl *url.URL, srcOpt ...WDASourceOption) (s string, err error) {
	tmp, _ := url.Parse(baseUrl.String())
	if len(srcOpt) != 0 {
		q := tmp.Query()
		if vFormat, ok := srcOpt[0]["format"]; ok {
			q.Set("format", vFormat.(string))
		}
		if vEAttr, ok := srcOpt[0]["excluded_attributes"]; ok {
			q.Set("excluded_attributes", vEAttr.(string))
		}
		tmp.RawQuery = q.Encode()
	}
	var wdaResp wdaResponse
	if wdaResp, err = internalGet("Source", urlJoin(tmp, "/source")); err != nil {
		return "", err
	}
	return wdaResp.getValue().String(), nil
}

// Source
func (c *Client) Source(srcOpt ...WDASourceOption) (s string, err error) {
	return source(c.deviceURL, srcOpt...)
}

func accessibleSource(baseUrl *url.URL) (sJson string, err error) {
	var wdaResp wdaResponse
	if wdaResp, err = internalGet("AccessibleSource", urlJoin(baseUrl, "/wda/accessibleSource")); err != nil {
		return "", err
	}
	return wdaResp.getValue().String(), nil
}

// AccessibleSource
//
// Return application elements accessibility tree
//
// ignore all elements except for the main window for accessibility tree
func (c *Client) AccessibleSource() (sJson string, err error) {
	return accessibleSource(c.deviceURL)
}

type WDAActiveAppInfo struct {
	ProcessArguments struct {
		Env  interface{}   `json:"env"`
		Args []interface{} `json:"args"`
	} `json:"processArguments"`
	Name string `json:"name"`
	WDAAppBaseInfo
	_string string
}

func (aai WDAActiveAppInfo) String() string {
	return aai._string
}

type WDAAppBaseInfo struct {
	Pid      int    `json:"pid"`
	BundleID string `json:"bundleId"`
}

// activeAppInfo
//
// {
//    "processArguments": {
//        "env": {},
//        "args": []
//    },
//    "name": "",
//    "pid": 57,
//    "bundleId": "com.apple.springboard"
// }
func activeAppInfo(baseUrl *url.URL) (wdaActiveAppInfo WDAActiveAppInfo, err error) {
	var wdaResp wdaResponse
	if wdaResp, err = internalGet("ActiveAppInfo", urlJoin(baseUrl, "/wda/activeAppInfo")); err != nil {
		return WDAActiveAppInfo{}, err
	}

	wdaActiveAppInfo._string = wdaResp.getValue().String()
	err = json.Unmarshal([]byte(wdaActiveAppInfo._string), &wdaActiveAppInfo)
	// err = json.Unmarshal(wdaResp.getValue2Bytes(), &wdaActiveAppInfo)
	return
}

// ActiveAppInfo
//
// get current active application
func (c *Client) ActiveAppInfo() (wdaActiveAppInfo WDAActiveAppInfo, err error) {
	return activeAppInfo(c.deviceURL)
}

func (c *Client) IsWdaHealth() (isHealth bool, err error) {
	var wdaResp wdaResponse
	if wdaResp, err = internalGet("IsWdaHealth", urlJoin(c.deviceURL, "/health")); err != nil {
		return false, err
	}
	if wdaResp.String() != "I-AM-ALIVE" {
		return false, nil
	}
	return true, nil
}

func (c *Client) WdaShutdown() (err error) {
	_, err = internalGet("WdaShutdown", urlJoin(c.deviceURL, "/wda/shutdown"))
	return
}

func (c *Client) tttTmp() {
	wdaResp, err := internalGet("tttTmp", urlJoin(c.deviceURL, "/health"))
	fmt.Println(err, wdaResp)
}
