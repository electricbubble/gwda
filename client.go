package gwda

import (
	"encoding/json"
	"errors"
	"net/url"
)

type Client struct {
	deviceURL *url.URL
	// sessionID string
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
	var wdaResp wdaResponse = []byte(sJson)
	if !wdaResp.getByPath("value.ready").Bool() {
		return nil, errors.New("device is not ready")
	}
	// if c.sessionID, err = wdaResp.getSessionID(); err != nil {
	// 	return nil, err
	// }
	return c, nil
}

// type WDAAppCapabilities map[string]interface{}
//
// func NewWDAAppCapabilities() WDAAppCapabilities {
// 	wdaAppCapabilities := make(WDAAppCapabilities)
// 	return wdaAppCapabilities
// }

// NewSession Creates and saves new session for application
func (c *Client) NewSession(bundleId ...string) (s *Session, err error) {
	// TODO
	//  shouldUseTestManagerForVisibilityDetection
	//  shouldUseCompactResponses
	//  elementResponseAttributes
	//  maxTypingFrequency
	//  shouldUseSingletonTestManager
	//  eventloopIdleDelaySec
	//  app
	//  arguments
	//  environment
	//  defaultAlertAction	// defaultAlertAction:@"accept"];	defaultAlertAction:@"dismiss"];
	// ⬆️ ["defaultAlertAction"] Creates and saves new session for application with default alert handling behaviour
	capabilities := newWdaBody() // .set("shouldWaitForQuiescence", false)
	// capabilities.set("defaultAlertAction", "accept")
	// capabilities.set("defaultAlertAction", "dismiss")
	if len(bundleId) != 0 {
		capabilities.setBundleID(bundleId[0])
		capabilities.set("shouldWaitForQuiescence", true)
	}
	body := newWdaBody().set("capabilities", newWdaBody().set("alwaysMatch", capabilities))
	var wdaResp wdaResponse
	if wdaResp, err = internalPost("New Session", urlJoin(c.deviceURL, "session"), body); err != nil {
		return nil, err
	}
	s = &Session{}
	sid := ""
	if sid = wdaResp.getByPath("sessionId").String(); sid == "" {
		return nil, errors.New("not find sessionId")
	}
	// s.bundleID = bundleId
	// c.deviceURL 已在新建时校验过, 理论上此处不再出现错误
	s.sessionURL, _ = url.Parse(urlJoin(c.deviceURL, "session", sid))
	// if err = s.Launch(bundleId); err != nil {
	// 	return nil, err
	// }
	return s, nil
}

// LaunchUnattachedApp
// Launch the app with the specified bundle ID
func (c *Client) LaunchUnattachedApp(bundleId string) (err error) {
	body := newWdaBody().setBundleID(bundleId)
	_, err = internalPost("LaunchUnattached", urlJoin(c.deviceURL, "/wda/apps/launchUnattached"), body)
	return
}

// Status Checking service status
func (c *Client) Status() (sJson string, err error) {
	var wdaResp wdaResponse
	if wdaResp, err = internalGet("Status", urlJoin(c.deviceURL, "status")); err != nil {
		return "", err
	}
	return wdaResp.String(), nil
}

// HomeScreen Forces the device under test to switch to the home screen
// 1. pressButton
// 2. WaitUntilApplicationBoardIsVisible
func (c *Client) HomeScreen() (err error) {
	_, err = internalPost("Homescreen", urlJoin(c.deviceURL, "wda", "homescreen"), nil)
	return
}

// HealthCheck
// Checks health of XCTest by:
// 1) Querying application for some elements,
// 2) Triggering some device events.
//
// !!! Health check might modify simulator state so it should only be called in-between testing sessions
func (c *Client) HealthCheck() (err error) {
	_, err = internalGet("HealthCheck", urlJoin(c.deviceURL, "wda", "healthcheck"))
	return
}

func isLocked(baseUrl *url.URL) (isLocked bool, err error) {
	var wdaResp wdaResponse
	if wdaResp, err = internalGet("Locked", urlJoin(baseUrl, "wda", "locked")); err != nil {
		return false, err
	}
	return wdaResp.getValue().Bool(), nil
}

// Locked whether the screen is locked
func (c *Client) IsLocked() (bool, error) {
	return isLocked(c.deviceURL)
}

func unlock(baseUrl *url.URL) (err error) {
	_, err = internalPost("Unlock", urlJoin(baseUrl, "wda", "unlock"), nil)
	return
}

// Unlock unlock screen
func (c *Client) Unlock() (err error) {
	return unlock(c.deviceURL)
}

func lock(baseUrl *url.URL) (err error) {
	_, err = internalPost("Lock", urlJoin(baseUrl, "wda", "lock"), nil)
	return
}

// Lock
func (c *Client) Lock() (err error) {
	return lock(c.deviceURL)
}

// TODO Screenshot
// func (c *Client) Screenshot() {}

func source(baseUrl *url.URL, formattedAsJson ...bool) (s string, err error) {
	tmp, _ := url.Parse(baseUrl.String())
	if len(formattedAsJson) != 0 && formattedAsJson[0] {
		q := tmp.Query()
		q.Set("format", "json")
		tmp.RawQuery = q.Encode()
	}
	var wdaResp wdaResponse
	if wdaResp, err = internalGet("Source", urlJoin(tmp, "source")); err != nil {
		return "", err
	}
	return wdaResp.getValue().String(), nil
}

// Source
//
// Source aka tree
func (c *Client) Source(formattedAsJson ...bool) (s string, err error) {
	return source(c.deviceURL, formattedAsJson...)
}

func accessibleSource(baseUrl *url.URL) (sJson string, err error) {
	var wdaResp wdaResponse
	if wdaResp, err = internalGet("AccessibleSource", urlJoin(baseUrl, "wda", "accessibleSource")); err != nil {
		return "", err
	}
	return wdaResp.getValue().String(), nil
}

// AccessibleSource
func (c *Client) AccessibleSource() (sJson string, err error) {
	return accessibleSource(c.deviceURL)
}

type WDAActiveAppInfo struct {
	ProcessArguments struct {
		Env  interface{}   `json:"env"`
		Args []interface{} `json:"args"`
	} `json:"processArguments"`
	Name     string `json:"name"`
	Pid      int    `json:"pid"`
	BundleID string `json:"bundleId"`
	_String  string
}

func (aai WDAActiveAppInfo) String() string {
	return aai._String
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
	if wdaResp, err = internalGet("ActiveAppInfo", urlJoin(baseUrl, "wda", "activeAppInfo")); err != nil {
		return
	}

	wdaActiveAppInfo._String = wdaResp.getValue().String()
	err = json.Unmarshal([]byte(wdaActiveAppInfo._String), &wdaActiveAppInfo)
	// err = json.Unmarshal(wdaResp.getValue2Bytes(), &wdaActiveAppInfo)
	return
}

// ActiveAppInfo Constructor used to get current active application
func (c *Client) ActiveAppInfo() (wdaActiveAppInfo WDAActiveAppInfo, err error) {
	return activeAppInfo(c.deviceURL)
}

// func (c *Client) tttTmp() {
// 	bsJson, err := internalGet("tttTmp", urlJoin(c.deviceURL, "/wd/hub/source"))
// 	fmt.Println(err, string(bsJson))
// }
