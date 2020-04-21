package gwda

import (
	"testing"
)

var deviceURL = "http://localhost:8100"
var bundleId = "com.apple.Preferences"

func TestClient_NewSession(t *testing.T) {
	type sTmp struct {
		shouldUseTestManagerForVisibilityDetection *bool
	}
	// var wdaTrue *bool = true
	// t.Log(sTmp{shouldUseTestManagerForVisibilityDetection: *bool(true))})
	// t.Log(sTmp{shouldUseTestManagerForVisibilityDetection: sql.NullBool{Bool: true}})
	// return
	c, err := NewClient(deviceURL)
	if err != nil {
		t.Fatal(err)
	}
	bundleId = "com.taobao.taobao4iphone"
	// bundleId = "com.sgv.peanutsubwaywifi"
	Debug = true
	// _, err = c.NewSession()
	_, err = c.NewSession(bundleId)
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_LaunchUnattached(t *testing.T) {
	c, err := NewClient(deviceURL)
	if err != nil {
		t.Fatal(err)
	}
	Debug = true
	err = c.LaunchUnattachedApp(bundleId)
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_Status(t *testing.T) {
	Debug = true
	c, err := NewClient(deviceURL)
	if err != nil {
		t.Fatal(err)
	}
	wdaResp, err := c.Status()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(wdaResp)
}

func TestClient_Home(t *testing.T) {
	c, err := NewClient(deviceURL)
	if err != nil {
		t.Fatal(err)
	}
	Debug = true
	err = c.HomeScreen()
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_Locked(t *testing.T) {
	c, err := NewClient(deviceURL)
	if err != nil {
		t.Fatal(err)
	}
	Debug = true
	isLocked, err := c.IsLocked()
	if err != nil {
		t.Fatal(err)
	}
	if isLocked {
		t.Log("锁屏界面")
	} else {
		t.Log("非锁屏界面")
	}
}

func TestClient_Unlock(t *testing.T) {
	c, err := NewClient(deviceURL)
	if err != nil {
		t.Fatal(err)
	}
	Debug = true
	err = c.Unlock()
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_Lock(t *testing.T) {
	c, err := NewClient(deviceURL)
	if err != nil {
		t.Fatal(err)
	}
	Debug = true
	err = c.Lock()
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_Source(t *testing.T) {
	c, err := NewClient(deviceURL)
	if err != nil {
		t.Fatal(err)
	}
	Debug = true
	s, err := c.Source()
	// s, err := c.Source(true)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(s)

	t.Log(c.Source(true))
}

func TestClient_AccessibleSource(t *testing.T) {
	c, err := NewClient(deviceURL)
	if err != nil {
		t.Fatal(err)
	}
	Debug = true
	s, err := c.AccessibleSource()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(s)
}

func TestClient_ActiveAppInfo(t *testing.T) {
	c, err := NewClient(deviceURL)
	if err != nil {
		t.Fatal(err)
	}
	Debug = true
	info, err := c.ActiveAppInfo()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(info)
}

func TestTmp(t *testing.T) {
	// body := make(map[string]interface{})
	// body["capabilities"] = map[string]string{"bundleId": "com.netease.cloudmusic"}
	// bsResp, err := internalPost("tmp", deviceURL+"/session", body)
	// bsResp, err := internalGet("test", deviceURL+"/session/4713B32E-4F89-42CC-9118-3DB4C3A18A75")

	// body := newWdaBody().set("bundleId", "com.netease.cloudmusic").set("shouldWaitForQuiescence", false)
	// bodyCap := newWdaBody().set("capabilities", body)
	// bsBody, err := json.Marshal(body)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// t.Log(string(bsBody))
	// bsBody, _ = json.Marshal(bodyCap)
	// t.Log(string(bsBody))
	// return

	Debug = true
	c, err := NewClient(deviceURL)
	if err != nil {
		t.Fatal(err)
	}
	// err = c.HealthCheck()
	// t.Log("#@#", err)
	// wdaResp, err := internalGet("healthcheck", urlJoin(c.deviceURL, "/wda/healthcheck", ))
	body := newWdaBody()
	_ = body
	body.setBundleID("com.apple.Preferences")
	wdaResp, err := internalPost("#TEMP", urlJoin(c.deviceURL, "/wda/apps/launchUnattached"), body)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// t.Log(wdaResp)
	// s, err := c.NewSession("com.netease.cloudmusic")
	// s, err := c.NewSession("com.apple.mobilesafari")
	// if err != nil {
	// 	t.Fatal(err)
	// }
	//
	// // bsJson, err := s.Tap(230, 130)
	// // bsJson, err = s.Tap(210, 290)
	// // c.tttTmp()
	// _ = s
	// t.Log("client:", c.sessionID)
	// t.Log("session:", s.sessionID)
	// t.Log(c.ActiveAppInfo())
	// s.tttTmp()
	// // bsJson, err := s.ActiveAppInfo()
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// t.Log(string(bsJson))
	t.Log(wdaResp)
}
