package gwda

import (
	"os"
	"path/filepath"
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

func TestClient_AppLaunchUnattached(t *testing.T) {
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

func TestClient_Homescreen(t *testing.T) {
	Debug = true
	c, err := NewClient(deviceURL)
	// c, err := NewClient(deviceURL, true)
	if err != nil {
		t.Fatal(err)
	}
	err = c.Homescreen()
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_AlertAccept(t *testing.T) {
	c, err := NewClient(deviceURL)
	if err != nil {
		t.Fatal(err)
	}
	Debug = true
	// err = c.AlertAccept()
	err = c.AlertAccept("始终允许")
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_AlertDismiss(t *testing.T) {
	c, err := NewClient(deviceURL)
	if err != nil {
		t.Fatal(err)
	}
	Debug = true
	// err = c.AlertDismiss()
	err = c.AlertDismiss("不允许")
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_AlertText(t *testing.T) {
	c, err := NewClient(deviceURL)
	if err != nil {
		t.Fatal(err)
	}
	Debug = true
	text, err := c.AlertText()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(text)
}

func TestClient_IsLocked(t *testing.T) {
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

func TestClient_DeviceInfo(t *testing.T) {
	c, err := NewClient(deviceURL)
	if err != nil {
		t.Fatal(err)
	}
	Debug = true
	deviceInfo, err := c.DeviceInfo()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(deviceInfo.Name)
	t.Log(deviceInfo.TimeZone)
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

func TestClient_ScreenshotToImage(t *testing.T) {
	c, err := NewClient(deviceURL)
	if err != nil {
		t.Fatal(err)
	}
	Debug = true
	// toPng, err := c.ScreenshotToPng()
	img, format, err := c.ScreenshotToImage()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("图片的格式:", format)
	t.Log("图片的大小:", img.Bounds().Size())
}

func TestClient_ScreenshotToDisk(t *testing.T) {
	c, err := NewClient(deviceURL)
	if err != nil {
		t.Fatal(err)
	}
	Debug = true
	userHomeDir, _ := os.UserHomeDir()
	// err = c.ScreenshotToDiskAsPng(filepath.Join(userHomeDir, "Desktop", "c1.png"))
	err = c.ScreenshotToDisk(filepath.Join(userHomeDir, "Desktop", "c1.png"))
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

func TestClient_IsWdaHealth(t *testing.T) {
	c, err := NewClient(deviceURL)
	if err != nil {
		t.Fatal(err)
	}
	Debug = true
	isWdaHealth, err := c.IsWdaHealth()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(isWdaHealth)
}

func TestClient_WdaShutdown(t *testing.T) {
	c, err := NewClient(deviceURL)
	if err != nil {
		t.Fatal(err)
	}
	Debug = true
	err = c.WdaShutdown()
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.IsWdaHealth()
	if err == nil {
		t.Fatal("wda 关闭失败")
	}
	t.Log(err)
}

func TestTmp(t *testing.T) {
	c, err := NewClient(deviceURL)
	if err != nil {
		t.Fatal(err)
	}
	Debug = true
	c.tttTmp()
}
