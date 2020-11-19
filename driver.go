package gwda

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	goUSBMux "github.com/electricbubble/go-usbmuxd-device"
	"net"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

// NewDriver creates new remote client, this will also start a new session.
func NewDriver(capabilities Capabilities, urlPrefix string) (driver WebDriver, err error) {
	wd := new(remoteWD)
	if wd.urlPrefix, err = url.Parse(urlPrefix); err != nil {
		return nil, err
	}
	var sessionInfo SessionInfo
	if sessionInfo, err = wd.NewSession(capabilities); err != nil {
		return nil, err
	}
	wd.sessionId = sessionInfo.SessionId
	return wd, nil
}

func NewUSBDriver(capabilities Capabilities, device ...Device) (driver WebDriver, err error) {
	if len(device) == 0 {
		if device, err = DeviceList(); err != nil {
			return nil, err
		}
	}
	dev := device[0]

	wd := new(remoteWD)
	wd.viaUSB = true

	var conn net.Conn
	if conn, err = goUSBMux.NewUSBHub().CreateConnect(dev.DeviceID(), dev.Port); err != nil {
		return nil, fmt.Errorf("create connection: %w", err)
	}
	wd.httpClient = convertToHTTPClient(conn)

	if conn, err = goUSBMux.NewUSBHub().CreateConnect(dev.DeviceID(), dev.MjpegPort); err != nil {
		return nil, fmt.Errorf("usb %w", err)
	}
	wd.mjpegClient = convertToHTTPClient(conn)

	if wd.urlPrefix, err = url.Parse("http://" + dev.serialNumber); err != nil {
		return nil, err
	}
	_, err = wd.NewSession(capabilities)
	return wd, err
}

var _ WebDriver = (*remoteWD)(nil)

func (wd *remoteWD) _requestURL(tmpURL *url.URL, elem ...string) string {
	var tmp *url.URL
	if tmpURL == nil {
		tmpURL = wd.urlPrefix
	}
	tmp, _ = url.Parse(tmpURL.String())
	tmp.Path = path.Join(append([]string{tmpURL.Path}, elem...)...)
	return tmp.String()
}

func (wd *remoteWD) _usbHTTPClient() []*http.Client {
	if wd.viaUSB {
		return []*http.Client{wd.httpClient}
	}
	return nil
}

func (wd *remoteWD) executeGet(pathElem ...string) (rawResp rawResponse, err error) {
	return executeHTTP(http.MethodGet, wd._requestURL(nil, pathElem...), nil, wd._usbHTTPClient()...)
}

func (wd *remoteWD) executePost(data interface{}, pathElem ...string) (rawResp rawResponse, err error) {
	var bsJSON []byte = nil
	if data != nil {
		if bsJSON, err = json.Marshal(data); err != nil {
			return nil, err
		}
	}
	return executeHTTP(http.MethodPost, wd._requestURL(nil, pathElem...), bsJSON, wd._usbHTTPClient()...)
}

func (wd *remoteWD) executeDelete(pathElem ...string) (rawResp rawResponse, err error) {
	return executeHTTP(http.MethodDelete, wd._requestURL(nil, pathElem...), nil, wd._usbHTTPClient()...)
}

func (wd *remoteWD) GetMjpegHTTPClient() *http.Client {
	return wd.mjpegClient
}

type remoteWD struct {
	urlPrefix *url.URL
	sessionId string

	viaUSB                  bool
	httpClient, mjpegClient *http.Client
}

func (wd *remoteWD) NewSession(capabilities Capabilities) (sessionInfo SessionInfo, err error) {
	// [[FBRoute POST:@"/session"].withoutSession respondWithTarget:self action:@selector(handleCreateSession:)]
	data := make(map[string]interface{})
	if capabilities == nil || len(capabilities) == 0 {
		data["capabilities"] = make(map[string]interface{})
	} else {
		data["capabilities"] = map[string]interface{}{"alwaysMatch": capabilities}
	}

	var rawResp rawResponse
	if rawResp, err = wd.executePost(data, "/session"); err != nil {
		return SessionInfo{}, err
	}
	if sessionInfo, err = rawResp.valueConvertToSessionInfo(); err != nil {
		return SessionInfo{}, err
	}
	wd.sessionId = sessionInfo.SessionId
	return
}

func (wd *remoteWD) ActiveSession() (sessionInfo SessionInfo, err error) {
	// [[FBRoute GET:@""] respondWithTarget:self action:@selector(handleGetActiveSession:)]
	var rawResp rawResponse
	if rawResp, err = wd.executeGet("/session", wd.sessionId); err != nil {
		return SessionInfo{}, err
	}
	if sessionInfo, err = rawResp.valueConvertToSessionInfo(); err != nil {
		return SessionInfo{}, err
	}
	return
}

func (wd *remoteWD) DeleteSession() (err error) {
	// [[FBRoute DELETE:@""] respondWithTarget:self action:@selector(handleDeleteSession:)]
	_, err = wd.executeDelete("/session", wd.sessionId)
	return
}

func (wd *remoteWD) Status() (deviceStatus DeviceStatus, err error) {
	// [[FBRoute GET:@"/status"].withoutSession respondWithTarget:self action:@selector(handleGetStatus:)]
	var rawResp rawResponse
	if rawResp, err = wd.executeGet("/status"); err != nil {
		return DeviceStatus{}, err
	}
	var reply = new(struct{ Value struct{ DeviceStatus } })
	if err = json.Unmarshal(rawResp, reply); err != nil {
		return DeviceStatus{}, err
	}
	deviceStatus = reply.Value.DeviceStatus
	return
}

func (wd *remoteWD) DeviceInfo() (deviceInfo DeviceInfo, err error) {
	// [[FBRoute GET:@"/wda/device/info"] respondWithTarget:self action:@selector(handleGetDeviceInfo:)]
	// [[FBRoute GET:@"/wda/device/info"].withoutSession
	var rawResp rawResponse
	if rawResp, err = wd.executeGet("/session", wd.sessionId, "/wda/device/info"); err != nil {
		return DeviceInfo{}, err
	}
	var reply = new(struct{ Value struct{ DeviceInfo } })
	if err = json.Unmarshal(rawResp, reply); err != nil {
		return DeviceInfo{}, err
	}
	deviceInfo = reply.Value.DeviceInfo
	return
}

func (wd *remoteWD) BatteryInfo() (batteryInfo BatteryInfo, err error) {
	// [[FBRoute GET:@"/wda/batteryInfo"] respondWithTarget:self action:@selector(handleGetBatteryInfo:)]
	var rawResp rawResponse
	if rawResp, err = wd.executeGet("/session", wd.sessionId, "/wda/batteryInfo"); err != nil {
		return BatteryInfo{}, err
	}
	var reply = new(struct{ Value struct{ BatteryInfo } })
	if err = json.Unmarshal(rawResp, reply); err != nil {
		return BatteryInfo{}, err
	}
	batteryInfo = reply.Value.BatteryInfo
	return
}

func (wd *remoteWD) WindowSize() (size Size, err error) {
	// [[FBRoute GET:@"/window/size"] respondWithTarget:self action:@selector(handleGetWindowSize:)]
	var rawResp rawResponse
	if rawResp, err = wd.executeGet("/session", wd.sessionId, "/window/size"); err != nil {
		return Size{}, err
	}
	var reply = new(struct{ Value struct{ Size } })
	if err = json.Unmarshal(rawResp, reply); err != nil {
		return Size{}, err
	}
	size = reply.Value.Size
	return
}

func (wd *remoteWD) Screen() (screen Screen, err error) {
	// [[FBRoute GET:@"/wda/screen"] respondWithTarget:self action:@selector(handleGetScreen:)]
	var rawResp rawResponse
	if rawResp, err = wd.executeGet("/session", wd.sessionId, "/wda/screen"); err != nil {
		return Screen{}, err
	}
	var reply = new(struct{ Value struct{ Screen } })
	if err = json.Unmarshal(rawResp, reply); err != nil {
		return Screen{}, err
	}
	screen = reply.Value.Screen
	return
}

func (wd *remoteWD) ActiveAppInfo() (info AppInfo, err error) {
	// [[FBRoute GET:@"/wda/activeAppInfo"] respondWithTarget:self action:@selector(handleActiveAppInfo:)]
	// [[FBRoute GET:@"/wda/activeAppInfo"].withoutSession
	var rawResp rawResponse
	if rawResp, err = wd.executeGet("/session", wd.sessionId, "/wda/activeAppInfo"); err != nil {
		return AppInfo{}, err
	}
	var reply = new(struct{ Value struct{ AppInfo } })
	if err = json.Unmarshal(rawResp, reply); err != nil {
		return AppInfo{}, err
	}
	info = reply.Value.AppInfo
	return
}

func (wd *remoteWD) ActiveAppsList() (appsList []AppBaseInfo, err error) {
	// [[FBRoute GET:@"/wda/apps/list"] respondWithTarget:self action:@selector(handleGetActiveAppsList:)]
	var rawResp rawResponse
	if rawResp, err = wd.executeGet("/session", wd.sessionId, "/wda/apps/list"); err != nil {
		return nil, err
	}
	var reply = new(struct{ Value []AppBaseInfo })
	if err = json.Unmarshal(rawResp, reply); err != nil {
		return nil, err
	}
	appsList = reply.Value
	return
}

func (wd *remoteWD) AppState(bundleId string) (runState AppState, err error) {
	// [[FBRoute POST:@"/wda/apps/state"] respondWithTarget:self action:@selector(handleSessionAppState:)]
	data := map[string]interface{}{"bundleId": bundleId}
	var rawResp rawResponse
	if rawResp, err = wd.executePost(data, "/session", wd.sessionId, "/wda/apps/state"); err != nil {
		return 0, err
	}
	var reply = new(struct{ Value AppState })
	if err = json.Unmarshal(rawResp, reply); err != nil {
		return 0, err
	}
	runState = reply.Value
	_ = rawResp
	return
}

func (wd *remoteWD) IsLocked() (locked bool, err error) {
	// [[FBRoute GET:@"/wda/locked"] respondWithTarget:self action:@selector(handleIsLocked:)]
	// [[FBRoute GET:@"/wda/locked"].withoutSession
	var rawResp rawResponse
	if rawResp, err = wd.executeGet("/session", wd.sessionId, "/wda/locked"); err != nil {
		return false, err
	}
	if locked, err = rawResp.valueConvertToBool(); err != nil {
		return false, err
	}
	return
}

func (wd *remoteWD) Unlock() (err error) {
	// [[FBRoute POST:@"/wda/unlock"] respondWithTarget:self action:@selector(handleUnlock:)]
	// [[FBRoute POST:@"/wda/unlock"].withoutSession
	_, err = wd.executePost(nil, "/session", wd.sessionId, "/wda/unlock")
	return
}

func (wd *remoteWD) Lock() (err error) {
	// [[FBRoute POST:@"/wda/lock"] respondWithTarget:self action:@selector(handleLock:)]
	// [[FBRoute POST:@"/wda/lock"].withoutSession
	_, err = wd.executePost(nil, "/session", wd.sessionId, "/wda/lock")
	return
}

func (wd *remoteWD) Homescreen() (err error) {
	// [[FBRoute POST:@"/wda/homescreen"].withoutSession respondWithTarget:self action:@selector(handleHomescreenCommand:)]
	_, err = wd.executePost(nil, "/wda/homescreen")
	return
}

func (wd *remoteWD) AlertText() (text string, err error) {
	// [[FBRoute GET:@"/alert/text"] respondWithTarget:self action:@selector(handleAlertGetTextCommand:)]
	// [[FBRoute GET:@"/alert/text"].withoutSession
	var rawResp rawResponse
	if rawResp, err = wd.executeGet("/session", wd.sessionId, "/alert/text"); err != nil {
		return "", err
	}
	if text, err = rawResp.valueConvertToString(); err != nil {
		return "", err
	}
	return
}

func (wd *remoteWD) AlertButtons() (btnLabels []string, err error) {
	// [[FBRoute GET:@"/wda/alert/buttons"] respondWithTarget:self action:@selector(handleGetAlertButtonsCommand:)]
	var rawResp rawResponse
	if rawResp, err = wd.executeGet("/session", wd.sessionId, "/wda/alert/buttons"); err != nil {
		return nil, err
	}
	var reply = new(struct{ Value []string })
	if err = json.Unmarshal(rawResp, reply); err != nil {
		return nil, err
	}
	btnLabels = reply.Value
	return
}

func (wd *remoteWD) AlertAccept(label ...string) (err error) {
	// [[FBRoute POST:@"/alert/accept"] respondWithTarget:self action:@selector(handleAlertAcceptCommand:)]
	// [[FBRoute POST:@"/alert/accept"].withoutSession
	data := make(map[string]interface{})
	if len(label) != 0 && label[0] != "" {
		data["name"] = label[0]
	}
	_, err = wd.executePost(data, "/alert/accept")
	return
}

func (wd *remoteWD) AlertDismiss(label ...string) (err error) {
	// [[FBRoute POST:@"/alert/dismiss"] respondWithTarget:self action:@selector(handleAlertDismissCommand:)]
	// [[FBRoute POST:@"/alert/dismiss"].withoutSession
	data := make(map[string]interface{})
	if len(label) != 0 && label[0] != "" {
		data["name"] = label[0]
	}
	_, err = wd.executePost(data, "/alert/dismiss")
	return
}

func (wd *remoteWD) AlertSendKeys(text string) (err error) {
	// [[FBRoute POST:@"/alert/text"] respondWithTarget:self action:@selector(handleAlertSetTextCommand:)]
	data := map[string]interface{}{"value": strings.Split(text, "")}
	_, err = wd.executePost(data, "/session", wd.sessionId, "/alert/text")
	return
}

func (wd *remoteWD) AppLaunch(bundleId string, launchOpt ...AppLaunchOption) (err error) {
	// [[FBRoute POST:@"/wda/apps/launch"] respondWithTarget:self action:@selector(handleSessionAppLaunch:)]
	data := make(map[string]interface{})
	if len(launchOpt) != 0 {
		data = launchOpt[0]
	}
	data["bundleId"] = bundleId
	_, err = wd.executePost(data, "/session", wd.sessionId, "/wda/apps/launch")
	return
}

func (wd *remoteWD) AppLaunchUnattached(bundleId string) (err error) {
	// [[FBRoute POST:@"/wda/apps/launchUnattached"].withoutSession respondWithTarget:self action:@selector(handleLaunchUnattachedApp:)]
	data := map[string]interface{}{"bundleId": bundleId}
	_, err = wd.executePost(data, "/wda/apps/launchUnattached")
	return
}

func (wd *remoteWD) AppTerminate(bundleId string) (successful bool, err error) {
	// [[FBRoute POST:@"/wda/apps/terminate"] respondWithTarget:self action:@selector(handleSessionAppTerminate:)]
	data := map[string]interface{}{"bundleId": bundleId}
	var rawResp rawResponse
	if rawResp, err = wd.executePost(data, "/session", wd.sessionId, "/wda/apps/terminate"); err != nil {
		return false, err
	}
	if successful, err = rawResp.valueConvertToBool(); err != nil {
		return false, err
	}
	return
}

func (wd *remoteWD) AppActivate(bundleId string) (err error) {
	// [[FBRoute POST:@"/wda/apps/activate"] respondWithTarget:self action:@selector(handleSessionAppActivate:)]
	data := map[string]interface{}{"bundleId": bundleId}
	_, err = wd.executePost(data, "/session", wd.sessionId, "/wda/apps/activate")
	return
}

func (wd *remoteWD) AppDeactivate(second float64) (err error) {
	// [[FBRoute POST:@"/wda/deactivateApp"] respondWithTarget:self action:@selector(handleDeactivateAppCommand:)]
	if second < 3 {
		second = 3.0
	}
	data := map[string]interface{}{"duration": second}
	_, err = wd.executePost(data, "/session", wd.sessionId, "/wda/deactivateApp")
	return
}

func (wd *remoteWD) AppAuthReset(resource ProtectedResource) (err error) {
	// [[FBRoute POST:@"/wda/resetAppAuth"] respondWithTarget:self action:@selector(handleResetAppAuth:)]
	data := map[string]interface{}{"resource": resource}
	_, err = wd.executePost(data, "/session", wd.sessionId, "/wda/resetAppAuth")
	return
}

func (wd *remoteWD) Tap(x, y int) error {
	return wd.TapFloat(float64(x), float64(y))
}

func (wd *remoteWD) TapFloat(x, y float64) (err error) {
	// [[FBRoute POST:@"/wda/tap/:uuid"] respondWithTarget:self action:@selector(handleTap:)]
	data := map[string]interface{}{
		"x": x,
		"y": y,
	}
	_, err = wd.executePost(data, "/session", wd.sessionId, "/wda/tap/0")
	return
}

func (wd *remoteWD) DoubleTap(x, y int) error {
	return wd.DoubleTapFloat(float64(x), float64(y))
}

func (wd *remoteWD) DoubleTapFloat(x, y float64) (err error) {
	// [[FBRoute POST:@"/wda/doubleTap"] respondWithTarget:self action:@selector(handleDoubleTapCoordinate:)]
	data := map[string]interface{}{
		"x": x,
		"y": y,
	}
	_, err = wd.executePost(data, "/session", wd.sessionId, "/wda/doubleTap")
	return
}

func (wd *remoteWD) TouchAndHold(x, y int, second ...float64) error {
	return wd.TouchAndHoldFloat(float64(x), float64(y), second...)
}

func (wd *remoteWD) TouchAndHoldFloat(x, y float64, second ...float64) (err error) {
	// [[FBRoute POST:@"/wda/touchAndHold"] respondWithTarget:self action:@selector(handleTouchAndHoldCoordinate:)]
	data := map[string]interface{}{
		"x": x,
		"y": y,
	}
	if len(second) == 0 || second[0] <= 0 {
		second = []float64{1.0}
	}
	data["duration"] = second[0]
	_, err = wd.executePost(data, "/session", wd.sessionId, "/wda/touchAndHold")
	return
}

func (wd *remoteWD) Drag(fromX, fromY, toX, toY int, pressForDuration ...float64) error {
	return wd.DragFloat(float64(fromX), float64(fromY), float64(toX), float64(toY), pressForDuration...)
}

func (wd *remoteWD) DragFloat(fromX, fromY, toX, toY float64, pressForDuration ...float64) (err error) {
	// [[FBRoute POST:@"/wda/dragfromtoforduration"] respondWithTarget:self action:@selector(handleDragCoordinate:)]
	data := map[string]interface{}{
		"fromX": fromX,
		"fromY": fromY,
		"toX":   toX,
		"toY":   toY,
	}
	if len(pressForDuration) == 0 || pressForDuration[0] < 0 {
		pressForDuration = []float64{1.0}
	}
	data["duration"] = pressForDuration[0]
	_, err = wd.executePost(data, "/session", wd.sessionId, "/wda/dragfromtoforduration")
	return
}

func (wd *remoteWD) Swipe(fromX, fromY, toX, toY int) error {
	return wd.SwipeFloat(float64(fromX), float64(fromY), float64(toX), float64(toY))
}

func (wd *remoteWD) SwipeFloat(fromX, fromY, toX, toY float64) error {
	return wd.DragFloat(fromX, fromY, toX, toY, 0)
}

func (wd *remoteWD) ForceTouch(x, y int, pressure float64, second ...float64) error {
	return wd.ForceTouchFloat(float64(x), float64(y), pressure, second...)
}

func (wd *remoteWD) ForceTouchFloat(x, y, pressure float64, second ...float64) error {
	if len(second) == 0 || second[0] <= 0 {
		second = []float64{1.0}
	}
	actions := NewTouchActions().
		Press(
			NewTouchActionPress().WithXYFloat(x, y).WithPressure(pressure)).
		Wait(second[0]).
		Release()
	return wd.PerformAppiumTouchActions(actions)
}

func (wd *remoteWD) PerformW3CActions(actions *W3CActions) (err error) {
	// [[FBRoute POST:@"/actions"] respondWithTarget:self action:@selector(handlePerformW3CTouchActions:)]
	data := map[string]interface{}{"actions": actions}
	_, err = wd.executePost(data, "/session", wd.sessionId, "/actions")
	return
}

func (wd *remoteWD) PerformAppiumTouchActions(touchActs *TouchActions) (err error) {
	// [[FBRoute POST:@"/wda/touch/perform"] respondWithTarget:self action:@selector(handlePerformAppiumTouchActions:)]
	// [[FBRoute POST:@"/wda/touch/multi/perform"]
	data := map[string]interface{}{"actions": touchActs}
	_, err = wd.executePost(data, "/session", wd.sessionId, "/wda/touch/multi/perform")
	return
}

func (wd *remoteWD) SetPasteboard(contentType PasteboardType, content string) (err error) {
	// [[FBRoute POST:@"/wda/setPasteboard"] respondWithTarget:self action:@selector(handleSetPasteboard:)]
	data := map[string]interface{}{
		"contentType": contentType,
		"content":     base64.StdEncoding.EncodeToString([]byte(content)),
	}
	_, err = wd.executePost(data, "/session", wd.sessionId, "/wda/setPasteboard")
	return
}

func (wd *remoteWD) GetPasteboard(contentType PasteboardType) (raw *bytes.Buffer, err error) {
	// [[FBRoute POST:@"/wda/getPasteboard"] respondWithTarget:self action:@selector(handleGetPasteboard:)]
	data := map[string]interface{}{"contentType": contentType}
	var rawResp rawResponse
	if rawResp, err = wd.executePost(data, "/session", wd.sessionId, "/wda/getPasteboard"); err != nil {
		return nil, err
	}
	if raw, err = rawResp.valueDecodeAsBase64(); err != nil {
		return nil, err
	}
	return
}

func (wd *remoteWD) SendKeys(text string, frequency ...int) (err error) {
	// [[FBRoute POST:@"/wda/keys"] respondWithTarget:self action:@selector(handleKeys:)]
	data := map[string]interface{}{"value": strings.Split(text, "")}
	if len(frequency) == 0 || frequency[0] <= 0 {
		frequency = []int{60}
	}
	data["frequency"] = frequency[0]
	_, err = wd.executePost(data, "/session", wd.sessionId, "/wda/keys")
	return
}

func (wd *remoteWD) PressButton(devBtn DeviceButton) (err error) {
	// [[FBRoute POST:@"/wda/pressButton"] respondWithTarget:self action:@selector(handlePressButtonCommand:)]
	data := map[string]interface{}{"name": devBtn}
	_, err = wd.executePost(data, "/session", wd.sessionId, "/wda/pressButton")
	return
}

func (wd *remoteWD) SiriActivate(text string) (err error) {
	// [[FBRoute POST:@"/wda/siri/activate"] respondWithTarget:self action:@selector(handleActivateSiri:)]
	data := map[string]interface{}{"text": text}
	_, err = wd.executePost(data, "/session", wd.sessionId, "/wda/siri/activate")
	return
}

func (wd *remoteWD) SiriOpenUrl(url string) (err error) {
	// [[FBRoute POST:@"/url"] respondWithTarget:self action:@selector(handleOpenURL:)]
	data := map[string]interface{}{"url": url}
	_, err = wd.executePost(data, "/session", wd.sessionId, "/url")
	return
}

func (wd *remoteWD) Orientation() (orientation Orientation, err error) {
	// [[FBRoute GET:@"/orientation"] respondWithTarget:self action:@selector(handleGetOrientation:)]
	var rawResp rawResponse
	if rawResp, err = wd.executeGet("/session", wd.sessionId, "/orientation"); err != nil {
		return "", err
	}
	var reply = new(struct{ Value Orientation })
	if err = json.Unmarshal(rawResp, reply); err != nil {
		return "", err
	}
	orientation = reply.Value
	return
}

func (wd *remoteWD) SetOrientation(orientation Orientation) (err error) {
	// [[FBRoute POST:@"/orientation"] respondWithTarget:self action:@selector(handleSetOrientation:)]
	data := map[string]interface{}{"orientation": orientation}
	_, err = wd.executePost(data, "/session", wd.sessionId, "/orientation")
	return
}

func (wd *remoteWD) Rotation() (rotation Rotation, err error) {
	// [[FBRoute GET:@"/rotation"] respondWithTarget:self action:@selector(handleGetRotation:)]
	var rawResp rawResponse
	if rawResp, err = wd.executeGet("/session", wd.sessionId, "/rotation"); err != nil {
		return Rotation{}, err
	}
	var reply = new(struct{ Value Rotation })
	if err = json.Unmarshal(rawResp, reply); err != nil {
		return Rotation{}, err
	}
	rotation = reply.Value
	return
}

func (wd *remoteWD) SetRotation(rotation Rotation) (err error) {
	// [[FBRoute POST:@"/rotation"] respondWithTarget:self action:@selector(handleSetRotation:)]
	_, err = wd.executePost(rotation, "/session", wd.sessionId, "/rotation")
	return
}

func (wd *remoteWD) MatchTouchID(isMatch bool) (err error) {
	// [FBRoute POST:@"/wda/touch_id"]
	data := map[string]interface{}{"match": isMatch}
	_, err = wd.executePost(data, "/session", wd.sessionId, "/wda/touch_id")
	return
}

func (wd *remoteWD) ActiveElement() (element WebElement, err error) {
	// [[FBRoute GET:@"/element/active"] respondWithTarget:self action:@selector(handleGetActiveElement:)]
	var rawResp rawResponse
	if rawResp, err = wd.executeGet("/session", wd.sessionId, "/element/active"); err != nil {
		return nil, err
	}
	var elementID string
	if elementID, err = rawResp.valueConvertToElementID(); err != nil {
		return nil, err
	}
	element = &remoteWE{parent: wd, id: elementID}
	return
}

func (wd *remoteWD) FindElement(by BySelector) (element WebElement, err error) {
	// [[FBRoute POST:@"/element"] respondWithTarget:self action:@selector(handleFindElement:)]
	using, value := by.getUsingAndValue()
	data := map[string]interface{}{
		"using": using,
		"value": value,
	}
	var rawResp rawResponse
	if rawResp, err = wd.executePost(data, "/session", wd.sessionId, "/element"); err != nil {
		return nil, err
	}
	var elementID string
	if elementID, err = rawResp.valueConvertToElementID(); err != nil {
		if errors.Is(err, errNoSuchElement) {
			return nil, fmt.Errorf("%w: unable to find an element using '%s', value '%s'", err, using, value)
		}
		return nil, err
	}
	element = &remoteWE{parent: wd, id: elementID}
	return
}

func (wd *remoteWD) FindElements(by BySelector) (elements []WebElement, err error) {
	// [[FBRoute POST:@"/elements"] respondWithTarget:self action:@selector(handleFindElements:)]
	using, value := by.getUsingAndValue()
	data := map[string]interface{}{
		"using": using,
		"value": value,
	}
	var rawResp rawResponse
	if rawResp, err = wd.executePost(data, "/session", wd.sessionId, "/elements"); err != nil {
		return nil, err
	}
	var elementIDs []string
	if elementIDs, err = rawResp.valueConvertToElementIDs(); err != nil {
		if errors.Is(err, errNoSuchElement) {
			return nil, fmt.Errorf("%w: unable to find an element using '%s', value '%s'", err, using, value)
		}
		return nil, err
	}
	elements = make([]WebElement, len(elementIDs))
	for i := range elementIDs {
		elements[i] = &remoteWE{parent: wd, id: elementIDs[i]}
	}
	return
}

func (wd *remoteWD) Screenshot() (raw *bytes.Buffer, err error) {
	// [[FBRoute GET:@"/screenshot"] respondWithTarget:self action:@selector(handleGetScreenshot:)]
	// [[FBRoute GET:@"/screenshot"].withoutSession respondWithTarget:self action:@selector(handleGetScreenshot:)]
	var rawResp rawResponse
	if rawResp, err = wd.executeGet("/session", wd.sessionId, "/screenshot"); err != nil {
		return nil, err
	}

	if raw, err = rawResp.valueDecodeAsBase64(); err != nil {
		return nil, err
	}
	return
}

func (wd *remoteWD) Source(srcOpt ...SourceOption) (source string, err error) {
	// [[FBRoute GET:@"/source"] respondWithTarget:self action:@selector(handleGetSourceCommand:)]
	// [[FBRoute GET:@"/source"].withoutSession
	tmp, _ := url.Parse(wd._requestURL(nil, "/session", wd.sessionId))
	toJsonRaw := false
	if len(srcOpt) != 0 {
		q := tmp.Query()
		if vFormat, ok := srcOpt[0]["format"]; ok {
			q.Set("format", vFormat.(string))
			if vFormat.(string) == "json" {
				toJsonRaw = true
			}
		}
		if vAttr, ok := srcOpt[0]["excluded_attributes"]; ok {
			q.Set("excluded_attributes", vAttr.(string))
		}
		tmp.RawQuery = q.Encode()
	}

	var rawResp rawResponse
	if rawResp, err = executeHTTP(http.MethodGet, wd._requestURL(tmp, "/source"), nil, wd._usbHTTPClient()...); err != nil {
		return "", nil
	}
	if toJsonRaw {
		var jr json.RawMessage
		if jr, err = rawResp.valueConvertToJsonRawMessage(); err != nil {
			return "", err
		}
		return string(jr), nil
	}
	if source, err = rawResp.valueConvertToString(); err != nil {
		return "", err
	}
	return
}

func (wd *remoteWD) AccessibleSource() (source string, err error) {
	// [[FBRoute GET:@"/wda/accessibleSource"] respondWithTarget:self action:@selector(handleGetAccessibleSourceCommand:)]
	// [[FBRoute GET:@"/wda/accessibleSource"].withoutSession
	var rawResp rawResponse
	if rawResp, err = wd.executeGet("/session", wd.sessionId, "/wda/accessibleSource"); err != nil {
		return "", err
	}
	var jr json.RawMessage
	if jr, err = rawResp.valueConvertToJsonRawMessage(); err != nil {
		return "", err
	}
	source = string(jr)
	return
}

func (wd *remoteWD) HealthCheck() (err error) {
	// [[FBRoute GET:@"/wda/healthcheck"].withoutSession respondWithTarget:self action:@selector(handleGetHealthCheck:)]
	_, err = wd.executeGet("/wda/healthcheck")
	return
}

func (wd *remoteWD) GetAppiumSettings() (settings map[string]interface{}, err error) {
	// [[FBRoute GET:@"/appium/settings"] respondWithTarget:self action:@selector(handleGetSettings:)]
	var rawResp rawResponse
	if rawResp, err = wd.executeGet("/session", wd.sessionId, "/appium/settings"); err != nil {
		return nil, err
	}
	var reply = new(struct{ Value map[string]interface{} })
	if err = json.Unmarshal(rawResp, reply); err != nil {
		return nil, err
	}
	settings = reply.Value
	return
}

func (wd *remoteWD) SetAppiumSettings(settings map[string]interface{}) (ret map[string]interface{}, err error) {
	// [[FBRoute POST:@"/appium/settings"] respondWithTarget:self action:@selector(handleSetSettings:)]
	data := map[string]interface{}{"settings": settings}
	var rawResp rawResponse
	if rawResp, err = wd.executePost(data, "/session", wd.sessionId, "/appium/settings"); err != nil {
		return nil, err
	}
	var reply = new(struct{ Value map[string]interface{} })
	if err = json.Unmarshal(rawResp, reply); err != nil {
		return nil, err
	}
	ret = reply.Value
	return
}

func (wd *remoteWD) IsWdaHealthy() (healthy bool, err error) {
	var rawResp rawResponse
	if rawResp, err = wd.executeGet("/health"); err != nil {
		return false, err
	}
	if string(rawResp) != "I-AM-ALIVE" {
		return false, nil
	}
	return true, nil
}

func (wd *remoteWD) WdaShutdown() (err error) {
	_, err = wd.executeGet("/wda/shutdown")
	return
}

func (wd *remoteWD) WaitWithTimeoutAndInterval(condition Condition, timeout, interval time.Duration) error {
	startTime := time.Now()
	for {
		done, err := condition(wd)
		if err != nil {
			return err
		}
		if done {
			return nil
		}

		if elapsed := time.Since(startTime); elapsed > timeout {
			return fmt.Errorf("timeout after %v", elapsed)
		}
		time.Sleep(interval)
	}
}

func (wd *remoteWD) WaitWithTimeout(condition Condition, timeout time.Duration) error {
	return wd.WaitWithTimeoutAndInterval(condition, timeout, DefaultWaitInterval)
}

func (wd *remoteWD) Wait(condition Condition) error {
	return wd.WaitWithTimeoutAndInterval(condition, DefaultWaitTimeout, DefaultWaitInterval)
}
