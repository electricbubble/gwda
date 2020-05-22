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
	WDADebug = true
	c, err := NewClient(deviceURL)
	checkErr(t, err)
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
	checkErr(t, err)
	// 如果使用了弹窗监控，却没有删除 session，新建的 session 也使用了弹窗监控，将会导致弹窗监控全都失效(弹窗监控是全局的，脱离于 session)
	// defer s.Delete()
	defer func() {
		_ = s.DeleteSession()
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
	checkErr(t, err)
	WDADebug = true
	// 锁屏界面下，启动会失败 `LSApplicationWorkspace failed to launch app`
	// 但是，如果被打开的 App 正在运行中（前台或后台），则不会报错
	// 但也不会点亮屏幕
	err = c.AppLaunchUnattached(bundleId)
	checkErr(t, err)
}

func TestClient_Status(t *testing.T) {
	WDADebug = true
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	status, err := c.Status()
	checkErr(t, err)
	t.Log(status)
}

func TestClient_Homescreen(t *testing.T) {
	WDADebug = true
	c, err := NewClient(deviceURL)
	// c, err := NewClient(deviceURL, true)
	checkErr(t, err)
	err = c.Homescreen()
	checkErr(t, err)
}

func TestClient_AlertAccept(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	WDADebug = true
	// err = c.AlertAccept()
	err = c.AlertAccept("始终允许")
	checkErr(t, err)
}

func TestClient_AlertDismiss(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	WDADebug = true
	// err = c.AlertDismiss()
	err = c.AlertDismiss("不允许")
	checkErr(t, err)
}

func TestClient_AlertText(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	WDADebug = true
	text, err := c.AlertText()
	checkErr(t, err)
	t.Log(text)
}

func TestClient_IsLocked(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	WDADebug = true
	isLocked, err := c.IsLocked()
	checkErr(t, err)
	if isLocked {
		t.Log("锁屏界面")
	} else {
		t.Log("非锁屏界面")
	}
}

func TestClient_Unlock(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	WDADebug = true
	err = c.Unlock()
	checkErr(t, err)
}

func TestClient_Lock(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	WDADebug = true
	err = c.Lock()
	checkErr(t, err)
}

func TestClient_DeviceInfo(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	WDADebug = true
	deviceInfo, err := c.DeviceInfo()
	checkErr(t, err)
	t.Log(deviceInfo.Name)
	t.Log(deviceInfo.TimeZone)
}

func TestClient_Screenshot(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	WDADebug = true
	_, err = c.Screenshot()
	checkErr(t, err)
}

func TestClient_ScreenshotToImage(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	WDADebug = true
	// toPng, err := c.ScreenshotToPng()
	img, format, err := c.ScreenshotToImage()
	checkErr(t, err)
	t.Log("图片的格式:", format)
	t.Log("图片的大小:", img.Bounds().Size())
}

func TestClient_ScreenshotToDisk(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	WDADebug = true
	userHomeDir, _ := os.UserHomeDir()
	// err = c.ScreenshotToDiskAsPng(filepath.Join(userHomeDir, "Desktop", "c1.png"))
	err = c.ScreenshotToDisk(filepath.Join(userHomeDir, "Desktop", "c1.png"))
	checkErr(t, err)
}

func TestClient_Source(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	WDADebug = true
	s, err := c.Source() // xml
	// s, err := c.Source(NewWDASourceOption().SetFormatAsJson())
	// s, err := c.Source(NewWDASourceOption().SetFormatAsDescription())
	// s, err := c.Source(NewWDASourceOption().SetFormatAsJson().SetExcludedAttributes([]string{"enabled", "visible", "type"}))
	// s, err := c.Source(NewWDASourceOption().SetExcludedAttributes([]string{"enabled", "visible", "type"}))
	checkErr(t, err)
	t.Log(s)

	// s2, err := c.Source()
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// t.Log(s2)
}

func TestClient_AccessibleSource(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	WDADebug = true
	s, err := c.AccessibleSource()
	checkErr(t, err)
	t.Log(s)
}

func TestClient_ActiveAppInfo(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	WDADebug = true
	info, err := c.ActiveAppInfo()
	checkErr(t, err)
	t.Log(info)
	t.Log(info.BundleID)
	t.Log(info.Pid)
}

func TestClient_IsWdaHealth(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	WDADebug = true
	isWdaHealth, err := c.IsWdaHealth()
	checkErr(t, err)
	t.Log(isWdaHealth)
}

func TestClient_WdaShutdown(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	WDADebug = true
	err = c.WdaShutdown()
	checkErr(t, err)
	_, err = c.IsWdaHealth()
	if err == nil {
		t.Fatal("wda 关闭失败")
	}
	t.Log(err)
}

func TestTmp(t *testing.T) {
	c, err := NewClient(deviceURL)
	checkErr(t, err)
	WDADebug = true
	c.tttTmp()
}

func checkErr(t *testing.T, err error, msg ...string) {
	if err != nil {
		if len(msg) == 0 {
			t.Fatal(err)
		} else {
			t.Fatal(msg, err)
		}
	}
}
