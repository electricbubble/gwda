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
	"strconv"
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
	// BundleId is required 如果是不存在的 bundleId 会导致 wda 内部报错导致接下来的操作都无法接收处理
	if len(opt) == 0 {
		opt = []WDAAppLaunchOption{NewWDAAppLaunchOption().SetShouldWaitForQuiescence(true)}
	}
	body := newWdaBody().setBundleID(bundleId)
	body.setAppLaunchOption(opt[0])
	_, err = internalPost("AppLaunch", urlJoin(s.sessionURL, "/wda/apps/launch"), body)
	return
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
	if len(seconds) == 0 {
		seconds = []float64{3.0}
	}
	body := newWdaBody().set("duration", seconds[0])
	wdaResp, err := internalPost("AppDeactivate", urlJoin(s.sessionURL, "/wda/deactivateApp"), body)
	if err != nil {
		return err
	}
	return wdaResp.getErrMsg()
}

const (
	WDATextBackspaceSequence = "\u0008"
	WDATextDeleteSequence    = "\u007F"
)

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

func tap(baseUrl *url.URL, x, y interface{}, elemUID ...string) (err error) {
	body := newWdaBody().setXY(x, y)
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

// TapCoordinate
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
		body.setXY(x, y)
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
		body.setXY(x, y)
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

func (s *Session) _getCenterCoordinates() (c WDACoordinate, err error) {
	if windowSize, err := s.WindowSize(); err != nil {
		return WDACoordinate{}, err
	} else {
		c = WDACoordinate{X: windowSize.Width / 2, Y: windowSize.Height / 2}
	}
	return
}

// SwipeUp
func (s *Session) SwipeUp() (err error) {
	var fromCoordinate, toCoordinate WDACoordinate
	if c, err := s._getCenterCoordinates(); err != nil {
		return err
	} else {
		fromCoordinate, toCoordinate = c, c
	}
	fromCoordinate.Y += 100
	toCoordinate.Y -= 100
	return s.SwipeCoordinate(fromCoordinate, toCoordinate)
}

// SwipeDown
func (s *Session) SwipeDown() (err error) {
	var fromCoordinate, toCoordinate WDACoordinate
	if c, err := s._getCenterCoordinates(); err != nil {
		return err
	} else {
		fromCoordinate, toCoordinate = c, c
	}
	fromCoordinate.Y -= 100
	toCoordinate.Y += 100
	return s.SwipeCoordinate(fromCoordinate, toCoordinate)
}

// SwipeLeft
func (s *Session) SwipeLeft() (err error) {
	var fromCoordinate, toCoordinate WDACoordinate
	if c, err := s._getCenterCoordinates(); err != nil {
		return err
	} else {
		fromCoordinate, toCoordinate = c, c
	}
	fromCoordinate.X += 100
	toCoordinate.X -= 100
	return s.SwipeCoordinate(fromCoordinate, toCoordinate)
}

// SwipeRight
func (s *Session) SwipeRight() (err error) {
	var fromCoordinate, toCoordinate WDACoordinate
	if c, err := s._getCenterCoordinates(); err != nil {
		return err
	} else {
		fromCoordinate, toCoordinate = c, c
	}
	fromCoordinate.X -= 100
	toCoordinate.X += 100
	return s.SwipeCoordinate(fromCoordinate, toCoordinate)
}

type WDAContentType string

const (
	WDAContentTypePlaintext WDAContentType = "plaintext"
	WDAContentTypeImage     WDAContentType = "image"
	WDAContentTypeUrl       WDAContentType = "url"
)

// SetPasteboard Sets data to the general pasteboard
func (s *Session) SetPasteboard(contentType WDAContentType, content string) (err error) {
	body := newWdaBody()
	body.set("contentType", contentType)
	body.set("content", base64.StdEncoding.EncodeToString([]byte(content)))

	_, err = internalPost("SetPasteboard", urlJoin(s.sessionURL, "/wda/setPasteboard"), body)
	return
}

// SetPasteboardForType
func (s *Session) SetPasteboardForPlaintext(content string) (err error) {
	return s.SetPasteboard(WDAContentTypePlaintext, content)
}

// SetPasteboardForImageFromFile
func (s *Session) SetPasteboardForImageFromFile(filename string) (err error) {
	var content []byte
	if content, err = ioutil.ReadFile(filename); err != nil {
		return err
	}
	return s.SetPasteboard(WDAContentTypeImage, string(content))
}

// SetPasteboardForUrl
func (s *Session) SetPasteboardForUrl(url string) (err error) {
	return s.SetPasteboard(WDAContentTypeUrl, url)
}

// GetPasteboard
//
// It might work when `WebDriverAgentRunner` is in foreground on real devices.
// https://github.com/appium/WebDriverAgent/issues/330
func (s *Session) GetPasteboard(contentType WDAContentType) (raw *bytes.Buffer, err error) {
	var wdaResp wdaResponse
	body := newWdaBody().set("contentType", contentType)
	// [FBRoute POST:@"/wda/getPasteboard"]
	if wdaResp, err = internalPost("GetPasteboard", urlJoin(s.sessionURL, "/wda/getPasteboard"), body); err != nil {
		return nil, err
	}
	if decodeString, err := base64.StdEncoding.DecodeString(wdaResp.getValue().String()); err != nil {
		return nil, err
	} else {
		raw = bytes.NewBuffer(decodeString)
		return raw, nil
	}
}

func (s *Session) GetPasteboardForPlaintext() (content string, err error) {
	var raw *bytes.Buffer
	if raw, err = s.GetPasteboard(WDAContentTypePlaintext); err != nil {
		return "", err
	}
	content = raw.String()
	return
}

func (s *Session) GetPasteboardForUrl() (content string, err error) {
	var raw *bytes.Buffer
	if raw, err = s.GetPasteboard(WDAContentTypeUrl); err != nil {
		return "", err
	}
	content = raw.String()
	return
}

func (s *Session) GetPasteboardForImage() (img image.Image, format string, err error) {
	var raw *bytes.Buffer
	if raw, err = s.GetPasteboard(WDAContentTypeImage); err != nil {
		return nil, "", err
	}
	return image.Decode(raw)
}

func (s *Session) GetPasteboardForImageToDisk(filename string) (err error) {
	var raw *bytes.Buffer
	if raw, err = s.GetPasteboard(WDAContentTypeImage); err != nil {
		return err
	}
	return ioutil.WriteFile(filename, raw.Bytes(), 0666)
}

type WDADeviceButtonName string

const (
	WDADeviceButtonHome       WDADeviceButtonName = "home"
	WDADeviceButtonVolumeUp   WDADeviceButtonName = "volumeUp"
	WDADeviceButtonVolumeDown WDADeviceButtonName = "volumeDown"
)

// PressButton
//
// Presses the corresponding hardware button on the device
// !!! not a synchronous action
func (s *Session) PressButton(wdaDeviceButton WDADeviceButtonName) (err error) {
	body := newWdaBody().set("name", wdaDeviceButton)
	_, err = internalPost("PressButton", urlJoin(s.sessionURL, "/wda/pressButton"), body)
	return
}

func (s *Session) PressHomeButton() (err error) {
	return s.PressButton(WDADeviceButtonHome)
}

func (s *Session) PressVolumeUpButton() (err error) {
	return s.PressButton(WDADeviceButtonVolumeUp)
}

func (s *Session) PressVolumeDownButton() (err error) {
	return s.PressButton(WDADeviceButtonVolumeDown)
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

func (s *Session) AlertSendKeys(text string) (err error) {
	// [FBRoute POST:@"/alert/text"]
	return sendKeys(urlJoin(s.sessionURL, "/alert/text"), text)
}

func (s *Session) AlertAccept(label ...string) (err error) {
	return alertAccept(s.sessionURL, label...)
}

func (s *Session) AlertDismiss(label ...string) (err error) {
	return alertDismiss(s.sessionURL, label...)
}

func (s *Session) AlertText() (text string, err error) {
	return alertText(s.sessionURL)
}

func (s *Session) AlertButtons() (buttons []string, err error) {
	var wdaResp wdaResponse
	// [FBRoute GET:@"/wda/alert/buttons"]
	if wdaResp, err = internalGet("AlertButtons", urlJoin(s.sessionURL, "/wda/alert/buttons")); err != nil {
		return nil, err
	}
	results := wdaResp.getValue().Array()
	buttons = make([]string, len(results))
	for i := range buttons {
		buttons[i] = results[i].String()
	}
	return
}

type WDAOrientation string

const (
	WDAOrientationPortrait           WDAOrientation = "PORTRAIT"                                   // Device oriented vertically, home button on the bottom
	WDAOrientationPortraitUpsideDown WDAOrientation = "UIA_DEVICE_ORIENTATION_PORTRAIT_UPSIDEDOWN" // Device oriented vertically, home button on the top
	WDAOrientationLandscapeLeft      WDAOrientation = "LANDSCAPE"                                  // Device oriented horizontally, home button on the right
	WDAOrientationLandscapeRight     WDAOrientation = "UIA_DEVICE_ORIENTATION_LANDSCAPERIGHT"      // Device oriented horizontally, home button on the left
)

func (v WDAOrientation) String() string {
	switch v {
	case WDAOrientationPortrait:
		return "Device oriented vertically, home button on the bottom"
	case WDAOrientationPortraitUpsideDown:
		return "Device oriented vertically, home button on the top"
	case WDAOrientationLandscapeLeft:
		return "Device oriented horizontally, home button on the right"
	case WDAOrientationLandscapeRight:
		return "Device oriented horizontally, home button on the left"
	default:
		return "UNKNOWN"
	}
}

func (s *Session) Orientation() (orientation WDAOrientation, err error) {
	var wdaResp wdaResponse
	// [FBRoute GET:@"/orientation"]
	if wdaResp, err = internalGet("Orientation", urlJoin(s.sessionURL, "/orientation")); err != nil {
		return "", err
	}
	return WDAOrientation(wdaResp.getValue().String()), nil
}

func (s *Session) SetOrientation(orientation WDAOrientation) (err error) {
	body := newWdaBody().set("orientation", orientation)
	// [FBRoute POST:@"/orientation"]
	_, err = internalPost("SetOrientation", urlJoin(s.sessionURL, "/orientation"), body)
	return
}

type WDARotation struct {
	X       int `json:"x"`
	Y       int `json:"y"`
	Z       int `json:"z"`
	_string string
}

func (r WDARotation) String() string {
	return r._string
}

func (s *Session) Rotation() (wdaRotation WDARotation, err error) {
	var wdaResp wdaResponse
	// [FBRoute GET:@"/rotation"]
	if wdaResp, err = internalGet("Rotation", urlJoin(s.sessionURL, "/rotation")); err != nil {
		return WDARotation{}, err
	}
	wdaRotation._string = wdaResp.getValue().String()
	err = json.Unmarshal([]byte(wdaRotation._string), &wdaRotation)
	return
}

func (s *Session) SetRotation(wdaRotation WDARotation) (err error) {
	body := newWdaBody()
	body.setXY(wdaRotation.X, wdaRotation.Y)
	body.set("z", wdaRotation.Z)
	// [FBRoute POST:@"/rotation"]
	_, err = internalPost("SetRotation", urlJoin(s.sessionURL, "/rotation"), body)
	return
}

//  █████  ██████  ██████  ██ ██    ██ ███    ███     ████████  ██████  ██    ██  ██████ ██   ██  █████   ██████ ████████ ██  ██████  ███    ██ ███████
// ██   ██ ██   ██ ██   ██ ██ ██    ██ ████  ████        ██    ██    ██ ██    ██ ██      ██   ██ ██   ██ ██         ██    ██ ██    ██ ████   ██ ██
// ███████ ██████  ██████  ██ ██    ██ ██ ████ ██        ██    ██    ██ ██    ██ ██      ███████ ███████ ██         ██    ██ ██    ██ ██ ██  ██ ███████
// ██   ██ ██      ██      ██ ██    ██ ██  ██  ██        ██    ██    ██ ██    ██ ██      ██   ██ ██   ██ ██         ██    ██ ██    ██ ██  ██ ██      ██
// ██   ██ ██      ██      ██  ██████  ██      ██        ██     ██████   ██████   ██████ ██   ██ ██   ██  ██████    ██    ██  ██████  ██   ████ ███████

type WDATouchActions []wdaBody

func NewWDATouchActions(cap ...int) *WDATouchActions {
	if len(cap) == 0 || cap[0] <= 0 {
		cap = []int{8}
	}
	tmp := make(WDATouchActions, 0, cap[0])
	return &tmp
}

// PerformTouchActions
//
// fb_performAppiumTouchActions
func (s *Session) PerformTouchActions(touchActions *WDATouchActions) (err error) {
	body := newWdaBody().set("actions", touchActions)
	// [FBRoute POST:@"/wda/touch/perform"]
	// [FBRoute POST:@"/wda/touch/multi/perform"]
	_, err = internalPost("PerformTouchActions", urlJoin(s.sessionURL, "/wda/touch/multi/perform"), body)
	return
}

type WDATouchActionOptionTap wdaBody

func NewWDATouchActionOptionTap() WDATouchActionOptionTap {
	return make(WDATouchActionOptionTap)
}
func (tao WDATouchActionOptionTap) SetXY(x, y int) WDATouchActionOptionTap {
	return WDATouchActionOptionTap(wdaBody(tao).setXY(x, y))
}
func (tao WDATouchActionOptionTap) SetXYFloat(x, y float64) WDATouchActionOptionTap {
	return WDATouchActionOptionTap(wdaBody(tao).setXY(x, y))
}
func (tao WDATouchActionOptionTap) SetElement(element *Element) WDATouchActionOptionTap {
	return WDATouchActionOptionTap(wdaBody(tao).set("element", element.UID))
}
func (tao WDATouchActionOptionTap) SetCount(count int) WDATouchActionOptionTap {
	return WDATouchActionOptionTap(wdaBody(tao).set("count", count))
}
func (ta *WDATouchActions) Tap(optTap WDATouchActionOptionTap) *WDATouchActions {
	tmp := newWdaBody().set("action", "tap")
	tmp.set("options", optTap)
	*ta = append(*ta, tmp)
	return ta
}

type WDATouchActionOptionLongPress wdaBody

func NewWDATouchActionOptionLongPress() WDATouchActionOptionLongPress {
	return make(WDATouchActionOptionLongPress)
}
func (tao WDATouchActionOptionLongPress) SetXY(x, y int) WDATouchActionOptionLongPress {
	return WDATouchActionOptionLongPress(wdaBody(tao).setXY(x, y))
}
func (tao WDATouchActionOptionLongPress) SetXYFloat(x, y float64) WDATouchActionOptionLongPress {
	return WDATouchActionOptionLongPress(wdaBody(tao).setXY(x, y))
}
func (tao WDATouchActionOptionLongPress) SetElement(element *Element) WDATouchActionOptionLongPress {
	return WDATouchActionOptionLongPress(wdaBody(tao).set("element", element.UID))
}
func (ta *WDATouchActions) LongPress(optLongPress WDATouchActionOptionLongPress) *WDATouchActions {
	tmp := newWdaBody().set("action", "longPress")
	tmp.set("options", optLongPress)
	*ta = append(*ta, tmp)
	return ta
}

type WDATouchActionOptionPress wdaBody

func NewWDATouchActionOptionPress() WDATouchActionOptionPress {
	return make(WDATouchActionOptionPress)
}
func (tao WDATouchActionOptionPress) SetXY(x, y int) WDATouchActionOptionPress {
	return WDATouchActionOptionPress(wdaBody(tao).setXY(x, y))
}
func (tao WDATouchActionOptionPress) SetXYFloat(x, y float64) WDATouchActionOptionPress {
	return WDATouchActionOptionPress(wdaBody(tao).setXY(x, y))
}
func (tao WDATouchActionOptionPress) SetElement(element *Element) WDATouchActionOptionPress {
	return WDATouchActionOptionPress(wdaBody(tao).set("element", element.UID))
}
func (tao WDATouchActionOptionPress) SetPressure(pressure float64) WDATouchActionOptionPress {
	return WDATouchActionOptionPress(wdaBody(tao).set("pressure", pressure))
}
func (ta *WDATouchActions) Press(optPress WDATouchActionOptionPress) *WDATouchActions {
	tmp := newWdaBody().set("action", "press")
	tmp.set("options", optPress)
	*ta = append(*ta, tmp)
	return ta
}

func (ta *WDATouchActions) Release() *WDATouchActions {
	tmp := newWdaBody().set("action", "release")
	*ta = append(*ta, tmp)
	return ta
}

type WDATouchActionOptionMoveTo wdaBody

func NewWDATouchActionOptionMoveTo() WDATouchActionOptionMoveTo {
	return make(WDATouchActionOptionMoveTo)
}
func (tao WDATouchActionOptionMoveTo) SetXY(x, y int) WDATouchActionOptionMoveTo {
	return WDATouchActionOptionMoveTo(wdaBody(tao).setXY(x, y))
}
func (tao WDATouchActionOptionMoveTo) SetXYFloat(x, y float64) WDATouchActionOptionMoveTo {
	return WDATouchActionOptionMoveTo(wdaBody(tao).setXY(x, y))
}
func (tao WDATouchActionOptionMoveTo) SetElement(element *Element) WDATouchActionOptionMoveTo {
	return WDATouchActionOptionMoveTo(wdaBody(tao).set("element", element.UID))
}
func (ta *WDATouchActions) MoveTo(optMoveTo WDATouchActionOptionMoveTo) *WDATouchActions {
	tmp := newWdaBody().set("action", "moveTo")
	tmp.set("options", optMoveTo)
	*ta = append(*ta, tmp)
	return ta
}

func (ta *WDATouchActions) Wait(duration ...float64) *WDATouchActions {
	tmp := newWdaBody().set("action", "wait")
	if len(duration) == 0 {
		duration = []float64{1.0}
	}
	tmp.set("options", newWdaBody().set("ms", duration[0]*1000))
	*ta = append(*ta, tmp)
	return ta
}

func (ta *WDATouchActions) Cancel() *WDATouchActions {
	tmp := newWdaBody().set("action", "cancel")
	*ta = append(*ta, tmp)
	return ta
}

// ██     ██ ██████   ██████      █████   ██████ ████████ ██  ██████  ███    ██ ███████
// ██     ██      ██ ██          ██   ██ ██         ██    ██ ██    ██ ████   ██ ██
// ██  █  ██  █████  ██          ███████ ██         ██    ██ ██    ██ ██ ██  ██ ███████
// ██ ███ ██      ██ ██          ██   ██ ██         ██    ██ ██    ██ ██  ██ ██      ██
//  ███ ███  ██████   ██████     ██   ██  ██████    ██    ██  ██████  ██   ████ ███████

func performActions(baseUrl *url.URL, actions *WDAActions) (err error) {
	body := newWdaBody().set("actions", actions)
	// [FBRoute POST:@"/actions"]
	_, err = internalPost("PerformActions", urlJoin(baseUrl, "/actions"), body)
	return
}

type WDAActions []wdaBody

func NewWDAActions(cap ...int) *WDAActions {
	if len(cap) == 0 || cap[0] <= 0 {
		cap = []int{8}
	}
	tmp := make(WDAActions, 0, cap[0])
	return &tmp
}

// PerformActions
//
// fb_performW3CActions
func (s *Session) PerformActions(actions *WDAActions) (err error) {
	return performActions(s.sessionURL, actions)
}

type WDAActionOptionFinger []wdaBody

func NewWDAActionOptionFinger(cap ...int) *WDAActionOptionFinger {
	if len(cap) == 0 || cap[0] <= 0 {
		cap = []int{8}
	}
	tmp := make(WDAActionOptionFinger, 0, cap[0])
	return &tmp
}

type WDAActionOptionFingerMove wdaBody

func NewWWDAActionOptionFingerMove() WDAActionOptionFingerMove {
	return WDAActionOptionFingerMove(newWdaBody().set("type", "pointerMove"))
}
func (ofm WDAActionOptionFingerMove) _setXY(x, y interface{}) WDAActionOptionFingerMove {
	return WDAActionOptionFingerMove(wdaBody(ofm).setXY(x, y))
}
func (ofm WDAActionOptionFingerMove) SetXY(x, y int) WDAActionOptionFingerMove {
	return ofm._setXY(x, y)
}
func (ofm WDAActionOptionFingerMove) SetXYFloat(x, y float64) WDAActionOptionFingerMove {
	return ofm._setXY(x, y)
}
func (ofm WDAActionOptionFingerMove) SetOrigin(element *Element) WDAActionOptionFingerMove {
	return WDAActionOptionFingerMove(wdaBody(ofm).set("origin", element.UID))
}
func (ofm WDAActionOptionFingerMove) SetDuration(duration float64) WDAActionOptionFingerMove {
	return WDAActionOptionFingerMove(wdaBody(ofm).set("duration", duration))
}

func (aof *WDAActionOptionFinger) Move(ofm WDAActionOptionFingerMove) *WDAActionOptionFinger {
	*aof = append(*aof, wdaBody(ofm))
	return aof
}
func (aof *WDAActionOptionFinger) Down() *WDAActionOptionFinger {
	*aof = append(*aof, newWdaBody().set("type", "pointerDown"))
	return aof
}
func (aof *WDAActionOptionFinger) Up() *WDAActionOptionFinger {
	*aof = append(*aof, newWdaBody().set("type", "pointerUp"))
	return aof
}
func (aof *WDAActionOptionFinger) Pause(duration ...float64) *WDAActionOptionFinger {
	if len(duration) == 0 || duration[0] < 0 {
		duration = []float64{0.5}
	}
	*aof = append(*aof, newWdaBody().set("type", "pause").set("duration", duration[0]*1000))
	return aof
}

func (act *WDAActions) _newTypeForKeyboard() wdaBody {
	pointer := newWdaBody().set("type", "key")
	pointer.set("id", "keyboard"+strconv.FormatInt(int64(len(*act)+1), 10))
	return pointer
}

func (act *WDAActions) SendKeys(text string) *WDAActions {
	keyboard := act._newTypeForKeyboard()
	ss := strings.Split(text, "")
	actOptKey := make([]wdaBody, 0, len(ss)+1)

	for i := range ss {
		actOptKey = append(actOptKey,
			newWdaBody().set("type", "keyDown").set("value", ss[i]),
			newWdaBody().set("type", "keyUp").set("value", ss[i]))
	}

	keyboard.set("actions", actOptKey)
	*act = append(*act, keyboard)
	return act
}

func (act *WDAActions) _newTypeForFinger() wdaBody {
	pointer := newWdaBody().set("type", "pointer")
	pointer.set("id", "finger"+strconv.FormatInt(int64(len(*act)+1), 10))
	pointer.set("parameters", newWdaBody().set("pointerType", "touch"))
	return pointer
}

func (act *WDAActions) FingerActionOption(actOptFinger *WDAActionOptionFinger) *WDAActions {
	pointer := act._newTypeForFinger()
	pointer.set("actions", *actOptFinger)
	*act = append(*act, pointer)
	return act
}

func (act *WDAActions) Tap(x, y int, element ...*Element) *WDAActions {
	optMove := NewWWDAActionOptionFingerMove().SetXY(x, y)
	if len(element) != 0 {
		optMove.SetOrigin(element[0])
	}
	actOptFinger := NewWDAActionOptionFinger().
		Move(optMove).
		Down().
		Pause(0.1).
		Up()

	return act.FingerActionOption(actOptFinger)
}

func (act *WDAActions) DoubleTap(x, y int, element ...*Element) *WDAActions {
	optMove := NewWWDAActionOptionFingerMove().SetXY(x, y)
	if len(element) != 0 {
		optMove.SetOrigin(element[0])
	}
	actOptFinger := NewWDAActionOptionFinger().
		Move(optMove).
		Down().
		Pause(0.1).
		Up().
		Pause(0.04).
		Down().
		Pause(0.1).
		Up()

	return act.FingerActionOption(actOptFinger)
}

func (act *WDAActions) Press(x, y int, duration float64, element ...*Element) *WDAActions {
	optMove := NewWWDAActionOptionFingerMove().SetXY(x, y)
	if len(element) != 0 {
		optMove.SetOrigin(element[0])
	}
	actOptFinger := NewWDAActionOptionFinger().
		Move(optMove).
		Down().
		Pause(duration).
		Up()

	return act.FingerActionOption(actOptFinger)
}

func (act *WDAActions) _swipe(fromX, fromY, toX, toY interface{}, element ...*Element) *WDAActions {
	optMoveFrom := NewWWDAActionOptionFingerMove()._setXY(fromX, fromY)
	optMoveTo := NewWWDAActionOptionFingerMove()._setXY(toX, toY)
	if len(element) != 0 {
		optMoveFrom.SetOrigin(element[0])
		optMoveTo.SetOrigin(element[0])
	}
	actOptFinger := NewWDAActionOptionFinger().
		Move(optMoveFrom).
		Down().
		Pause(0.25).
		Move(optMoveTo).
		Pause(0.25).
		Up()

	return act.FingerActionOption(actOptFinger)
}

func (act *WDAActions) Swipe(fromX, fromY, toX, toY int, element ...*Element) *WDAActions {
	return act._swipe(fromX, fromY, toX, toY, element...)
}

func (act *WDAActions) SwipeFloat(fromX, fromY, toX, toY float64, element ...*Element) *WDAActions {
	return act._swipe(fromX, fromY, toX, toY, element...)
}

func (act *WDAActions) SwipeCoordinate(fromCoordinate, toCoordinate WDACoordinate, element ...*Element) *WDAActions {
	return act._swipe(fromCoordinate.X, fromCoordinate.Y, toCoordinate.X, toCoordinate.Y, element...)
}

// MatchTouchID
//
// Matches or mismatches TouchID request
func (s *Session) MatchTouchID(isMatch bool) (bool, error) {
	body := newWdaBody().set("match", isMatch)
	// [FBRoute POST:@"/wda/touch_id"]
	wdaResp, err := internalPost("MatchTouchID", urlJoin(s.sessionURL, "/wda/touch_id"), body)
	return wdaResp.getValue().Bool(), err
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
func (s *Session) ActiveAppsList() (appsList []WDAAppBaseInfo, err error) {
	var wdaResp wdaResponse
	if wdaResp, err = internalGet("ActiveAppsList", urlJoin(s.sessionURL, "/wda/apps/list")); err != nil {
		return nil, err
	}
	appsList = make([]WDAAppBaseInfo, 0)
	err = json.Unmarshal([]byte(wdaResp.getValue().String()), &appsList)
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
//	level - Battery level in range [0.0, 1.0], where 1.0 means 100% charge.
//	state - Battery state. The following values are possible:
//	UIDeviceBatteryStateUnplugged = 1  // on battery, discharging
//	UIDeviceBatteryStateCharging = 2   // plugged in, less than 100%
//	UIDeviceBatteryStateFull = 3       // plugged in, at 100%
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
// CGRect frame = request.session.activeApplication.wdFrame;
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
// /wda/touch_id

func (s *Session) tttTmp() {
	body := newWdaBody()
	_ = body

	// element, err := s.FindElement(WDALocator{ClassName: WDAElementType{ScrollView: true}})
	// element, err := s.FindElement(WDALocator{Name: "自定手势作用区域"})
	// _, _ = element, err

	body.set("match", true)

	// [FBRoute POST:@"/wda/touch_id"]
	wdaResp, err := internalPost("###############", urlJoin(s.sessionURL, "/wda/touch_id"), body)
	_, _ = err, wdaResp
	// fmt.Println(err, wdaResp)
}
