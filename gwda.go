package gwda

import (
	"bytes"
	"encoding/json"
	"fmt"
	goUSBMux "github.com/electricbubble/go-usbmuxd-device"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"time"

	"github.com/tidwall/gjson"
)

func init() {
	httpProxy := os.Getenv("http_proxy")
	if httpProxy != "" {
		if proxyURL, err := url.Parse(httpProxy); err == nil {
			http.DefaultTransport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}
		}
	}
}

var wdaDebugFlag = false

var usbHTTPClient = make(map[string]*http.Client)

var DefaultWaitTimeout = time.Second * 60
var DefaultWaitInterval = time.Millisecond * 250

var wdaHeader = map[string]string{
	"Content-Type": "application/json;charset=UTF-8",
	"accept":       "application/json",
}

// urlJoin fix `path.Join`
func urlJoin(endpoint *url.URL, elem string, isWdaFirst ...bool) string {
	tmp, _ := url.Parse(endpoint.String())
	if len(isWdaFirst) != 0 && isWdaFirst[0] {
		tmp.Path = path.Join(endpoint.Path, "wda", elem)
	} else {
		tmp.Path = path.Join(endpoint.Path, elem)
	}
	return tmp.String()
}

func executeGet(actionName, url string) (wdaResp wdaResponse, err error) {
	return executeHTTP(actionName, http.MethodGet, url, nil)
}

func executePost(actionName, url string, body wdaBody) (wdaResp wdaResponse, err error) {
	return executeHTTP(actionName, http.MethodPost, url, body)
}

func executeDelete(actionName, url string) (wdaResp wdaResponse, err error) {
	return executeHTTP(actionName, http.MethodDelete, url, nil)
}

func executeHTTP(actionName, method, sURL string, body wdaBody) (wdaResp wdaResponse, err error) {
	var req *http.Request
	var reqBody io.Reader = nil
	var bsBody []byte
	if body != nil {
		if bsBody, err = json.Marshal(body); err != nil {
			return nil, fmt.Errorf("%s: invalid request body %w", actionName, err)
		}
		reqBody = bytes.NewBuffer(bsBody)
	}

	req, _ = http.NewRequest(method, sURL, reqBody)
	for k, v := range wdaHeader {
		req.Header.Set(k, v)
	}

	httpClient := http.DefaultClient

	filteredURL, _ := url.Parse(sURL)
	if filteredURL.Port() == "" && len(filteredURL.Host) == 40 {
		udid := filteredURL.Host
		filteredURL.Host = "__UDID__"
		if tmpClient, ok := usbHTTPClient[udid]; !ok {
			// much better for debugging
			return nil, fmt.Errorf("no http client: %s", sURL)
			// return nil, fmt.Errorf("no http client: %s", filteredURL.String())
		} else {
			httpClient = tmpClient
		}
	}

	debugLog(fmt.Sprintf("--> %s %s %s\n%s", method, filteredURL.String(), actionName, bsBody))

	start := time.Now()
	var resp *http.Response
	resp, err = httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to send request %w", actionName, err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	wdaResp, err = ioutil.ReadAll(resp.Body)

	if actionName == "Screenshot" {
		debugLog(fmt.Sprintf("<-- %s %s %d %s %s 'too long, don't display'\n",
			method, filteredURL.String(), resp.StatusCode, time.Now().Sub(start), actionName))
	} else {
		debugLog(fmt.Sprintf("<-- %s %s %d %s %s\n%s\n", method, filteredURL.String(), resp.StatusCode, time.Now().Sub(start), actionName, wdaResp))
	}

	if err != nil {
		return nil, fmt.Errorf("%s: failed to read response %w", actionName, err)
	}

	err = wdaResp.getErrMsg()
	return
}

type wdaBody map[string]interface{}

func newWdaBody() wdaBody {
	return make(wdaBody)
}
func (wb wdaBody) set(k string, v interface{}) (body wdaBody) {
	wb[k] = v
	return wb
}
func (wb wdaBody) setBundleID(bundleId string) (body wdaBody) {
	return wb.set("bundleId", bundleId)
}
func (wb wdaBody) setXY(x, y interface{}) (body wdaBody) {
	return wb.set("x", x).set("y", y)
}

func (wb wdaBody) setAppLaunchOption(opt WDAAppLaunchOption) (body wdaBody) {
	for k, v := range opt {
		wb.set(k, v)
	}
	return wb
}

type wdaResponse []byte

func (wdaResp wdaResponse) String() string {
	return string(wdaResp)
}
func (wdaResp wdaResponse) getByPath(path string) gjson.Result {
	return gjson.GetBytes(wdaResp, path)
}

func (wdaResp wdaResponse) getValue() gjson.Result {
	return gjson.GetBytes(wdaResp, "value")
}

func (wdaResp wdaResponse) getErrMsg() error {
	// {
	//  "value" : {
	//    "error" : "unknown error",
	//    "message" : "Error Domain=com.facebook.WebDriverAgent Code=1 \"Timed out while waiting until the screen gets unlocked\" UserInfo={NSLocalizedDescription=Timed out while waiting until the screen gets unlocked}",
	//    "traceback" : ""
	//  },
	//  "sessionId" : "215BB5C5-B189-496F-83B7-37CBBB2DC54E"
	// }
	wdaErrType := wdaResp.getByPath("value.error").String()
	// if wdaErrType == "" && wdaResp.getValue().Type == gjson.Null {
	if wdaErrType == "" {
		return nil
	}
	wdaErrMsg := wdaResp.getByPath("value.message").String()
	errText := wdaErrMsg
	// 获取 NSLocalizedDescription 的值
	re := regexp.MustCompile(`{.+?=(.+?)}`)
	subMatch := re.FindStringSubmatch(wdaErrMsg)
	if len(subMatch) == 2 {
		errText = subMatch[1]
	}
	return fmt.Errorf("%s: %s", wdaErrType, errText)
}

func WDADebug(b ...bool) {
	if len(b) == 0 {
		b = []bool{true}
	}
	wdaDebugFlag = b[0]

	if len(b) == 2 {
		goUSBMux.Debug(b[1])
	}
}

func debugLog(msg string) {
	if !wdaDebugFlag {
		return
	}
	log.Println("[DEBUG] " + msg)
}
