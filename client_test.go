package gwda

import (
	"testing"
)

var deviceURL = "http://localhost:8100"

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
	isLocked, err := c.Locked()
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
	sJson, err := c.ActiveAppInfo()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(sJson)
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
	wdaResp, err := internalGet("AppList", urlJoin(c.deviceURL, "/wda/apps/list", ))
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
