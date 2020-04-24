package gwda

import (
	"testing"
	"time"
)

var deviceURL = "http://localhost:8100"
var bundleId = "com.apple.Preferences"

func TestClient_NewSession(t *testing.T) {
	Debug = true
	c, err := NewClient(deviceURL)
	if err != nil {
		t.Fatal(err)
	}
	bundleId = "com.taobao.taobao4iphone"
	bundleId = "com.sgv.peanutsubwaywifi"
	// bundleId = "com.apple.camera"

	// _, err = c.NewSession()
	// 未解锁状态下指定 bundleId 会导致 wda 内部会报错导致接下来的操作都无法接收处理
	// s, err := c.NewSession(NewWDASessionCapability(bundleId))
	// s, err := c.NewSession(NewWDASessionCapability())
	s, err := c.NewSession(NewWDASessionCapability(bundleId).SetDefaultAlertAction(WDASessionAlertActionAccept))
	// s, err := c.NewSession(NewWDASessionCapability().SetDefaultAlertAction(WDASessionAlertActionAccept))
	// s, err := c.NewSession(NewWDASessionCapability())
	// s, err := c.NewSession()
	if err != nil {
		t.Fatal(err)
	}
	// 如果使用了弹窗监控，却没有删除 session，新建的 session 也使用了弹窗监控，将会导致弹窗监控全都失效(弹窗监控是全局的，脱离于 session)
	// defer s.Delete()
	defer func() {
		s.DeleteSession()
		time.Sleep(time.Second * 1)
	}()
	_ = s
	time.Sleep(time.Second * 30)
	// s.Delete()
	// s.AppLaunch(bundleId)
	// s.AppLaunch(bundleId)
	// s.AppLaunch("com.apple.DocumentsApp")
	// _, err = c.NewSession("com.apple.DocumentsApp")
}

func TestClient_LaunchUnattached(t *testing.T) {
	c, err := NewClient(deviceURL)
	if err != nil {
		t.Fatal(err)
	}
	Debug = true
	// 锁屏界面下，启动会失败 `LSApplicationWorkspace failed to launch app`
	// 但是，如果被打开的 App 正在运行中（前台或后台），则不会报错
	// 但也不会点亮屏幕
	err = c.AppLaunchUnattached(bundleId)
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
	status, err := c.Status()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(status)
}

func TestClient_Home(t *testing.T) {
	c, err := NewClient(deviceURL)
	if err != nil {
		t.Fatal(err)
	}
	Debug = true
	err = c.Homescreen()
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

func TestClient_Screenshot(t *testing.T) {
	c, err := NewClient(deviceURL)
	if err != nil {
		t.Fatal(err)
	}
	Debug = true
	_, err = c.Screenshot()
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ScreenshotToPng(t *testing.T) {
	c, err := NewClient(deviceURL)
	if err != nil {
		t.Fatal(err)
	}
	Debug = true
	toPng, err := c.ScreenshotToPng()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("图片大小", toPng.Bounds().Size())
}

func TestClient_ScreenshotToDiskAsPng(t *testing.T) {
	c, err := NewClient(deviceURL)
	if err != nil {
		t.Fatal(err)
	}
	Debug = true
	err = c.ScreenshotToDiskAsPng("/Users/hero/Desktop/3.png")
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
	s, err := c.Source() // xml
	// s, err := c.Source(NewWDASourceOption().SetFormatAsJson())
	// s, err := c.Source(NewWDASourceOption().SetFormatAsDescription())
	// s, err := c.Source(NewWDASourceOption().SetFormatAsJson().SetExcludedAttributes([]string{"enabled", "visible", "type"}))
	// s, err := c.Source(NewWDASourceOption().SetExcludedAttributes([]string{"enabled", "visible", "type"}))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(s)

	// s2, err := c.Source()
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// t.Log(s2)
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
	t.Log(info.BundleID)
	t.Log(info.Pid)
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
