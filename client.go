package gwda

import (
	"errors"
	"net/url"
)

type Client struct {
	deviceURL *url.URL
	// sessionID string
}

func NewClient(deviceURL string) (c *Client, err error) {
	c = &Client{}
	c.deviceURL, err = url.Parse(deviceURL)
	if err != nil {
		return nil, err
	}
	sJson, err := c.Status()
	if err != nil {
		return nil, err
	}
	var wdaResp wdaResponse = []byte(sJson)
	if !wdaResp.isReady() {
		return nil, errors.New("device is not ready")
	}
	// if c.sessionID, err = wdaResp.getSessionID(); err != nil {
	// 	return nil, err
	// }
	return c, nil
}

func (c *Client) NewSession(bundleId string) (s *Session, err error) {
	// body := make(map[string]interface{})
	// // capabilities := map[string]interface{}{"alwaysMatch": map[string]string{"bundleId": bundleId}}
	// // body["capabilities"] = capabilities
	// body["capabilities"] = map[string]interface{}{
	// 	"bundleId":                bundleId,
	// 	"shouldWaitForQuiescence": false,
	// }
	// "arguments":
	// "environment":
	// body["desiredCapabilities"] = map[string]string{"bundleId": bundleId}
	capabilities := newWdaBody().set("bundleId", bundleId).set("shouldWaitForQuiescence", false)
	body := newWdaBody().set("capabilities", capabilities)
	wdaResp, err := internalPost("create session", urlJoin(c.deviceURL, "session"), body)
	if err != nil {
		return nil, err
	}
	if err = wdaResp.getErrMsg(); err != nil {
		return nil, err
	}
	s = &Session{}
	sid := ""
	if sid, err = wdaResp.getSessionID(); err != nil {
		return nil, err
	}
	// s.bundleID = bundleId
	// c.deviceURL 已在新建时校验过, 理论上此处不再出现错误
	s.sessionURL, _ = url.Parse(urlJoin(c.deviceURL, "session", sid))
	// if err = s.Launch(bundleId); err != nil {
	// 	return nil, err
	// }
	return s, nil
}

// Status Checking service status
func (c *Client) Status() (sJson string, err error) {
	wdaResp, err := internalGet("检查服务状态", urlJoin(c.deviceURL, "status"))
	if err != nil {
		return "", err
	}
	return wdaResp.String(), nil
}

// HomeScreen Go to home screen
func (c *Client) HomeScreen() (err error) {
	wdaResp, err := internalPost("前往 Home 界面", urlJoin(c.deviceURL, "wda", "homescreen"), nil)
	if err != nil {
		return err
	}
	// value.error	"unknown error"
	// value.message	"Error Domain=com.facebook.WebDriverAgent Code=1 \"Timeout waiting until SpringBoard is visible\" UserInfo={NSLocalizedDescription=Timeout waiting until SpringBoard is visible}",
	return wdaResp.getErrMsg()
}

// HealthCheck
func (c *Client) HealthCheck() (err error) {
	wdaResp, err := internalGet("Hit healthcheck", urlJoin(c.deviceURL, "wda", "healthcheck"))
	if err != nil {
		return err
	}
	return wdaResp.getErrMsg()
}

// Locked
func (c *Client) Locked() (isLocked bool, err error) {
	wdaResp, err := internalGet("判断锁屏界面", urlJoin(c.deviceURL, "wda", "locked"))
	if err != nil {
		return false, err
	}
	return wdaResp.getValue().Bool(), nil
}

// Unlock
func (c *Client) Unlock() (err error) {
	wdaResp, err := internalPost("触发解锁", urlJoin(c.deviceURL, "wda", "unlock"), nil)
	if err != nil {
		return err
	}
	return wdaResp.getErrMsg()
}

// Lock
func (c *Client) Lock() (err error) {
	wdaResp, err := internalPost("触发锁屏", urlJoin(c.deviceURL, "wda", "lock"), nil)
	if err != nil {
		return err
	}
	return wdaResp.getErrMsg()
}

// TODO Screenshot
// func (c *Client) Screenshot() {}

// Source
//
// Source aka tree
func (c *Client) Source(formattedAsJson ...bool) (s string, err error) {
	tmp, _ := url.Parse(c.deviceURL.String())
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
func (c *Client) AccessibleSource() (sJson string, err error) {
	wdaResp, err := internalGet("Source aka tree", urlJoin(c.deviceURL, "wda", "accessibleSource"))
	if err != nil {
		return "", err
	}
	return wdaResp.getValue().String(), err
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
func (c *Client) ActiveAppInfo() (sJson string, err error) {
	wdaResp, err := internalGet("ActiveAppInfo", urlJoin(c.deviceURL, "wda", "activeAppInfo"))
	if err != nil {
		return "", err
	}
	return wdaResp.getValue().String(), nil
}

// TODO launchUnattached

// func (c *Client) tttTmp() {
// 	bsJson, err := internalGet("tttTmp", urlJoin(c.deviceURL, "/wd/hub/source"))
// 	fmt.Println(err, string(bsJson))
// }
