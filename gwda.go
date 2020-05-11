package gwda

import (
	"bytes"
	"encoding/json"
	"fmt"
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

var Debug = false
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

func internalGet(actionName, url string) (wdaResp wdaResponse, err error) {
	return internalDo(actionName, http.MethodGet, url, nil)
}

func internalPost(actionName, url string, body wdaBody) (wdaResp wdaResponse, err error) {
	return internalDo(actionName, http.MethodPost, url, body)
}

func internalDelete(actionName, url string) (wdaResp wdaResponse, err error) {
	return internalDo(actionName, http.MethodDelete, url, nil)
}

func internalDo(actionName, method, url string, body wdaBody) (wdaResp wdaResponse, err error) {
	var req *http.Request
	// 忽略 err 是因为在新建 Client 的时候已经校验了 URL 所以除此之外，应该不会出现其他错误
	var bsBody []byte
	if body != nil {
		// body 已经通过 `newWdaBody` 进行初始化和修改，理论上也不存在 err
		bsBody, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("%s: invalid request body %s", actionName, err.Error())
		}
		req, _ = http.NewRequest(method, url, bytes.NewBuffer(bsBody))
	} else {
		req, _ = http.NewRequest(method, url, nil)
	}
	for k, v := range wdaHeader {
		req.Header.Set(k, v)
	}
	start := time.Now()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to send request %s", actionName, err.Error())
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	output := ""
	if Debug {
		output = fmt.Sprintf("[DEBUG]↩︎\nActionName: %s\nMethod: %s\nURL: %s\n", actionName, method, req.URL.String())
		if body != nil {
			output += fmt.Sprintf("Body: %s\n", bsBody)
		}
		output += fmt.Sprintf("Duration: %s\n", time.Since(start).String())
	}
	wdaResp, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		if Debug {
			log.Println(output)
		}
		return nil, fmt.Errorf("%s: failed to read response %s", actionName, err.Error())
	}
	if Debug {
		if actionName == "Screenshot" {
			output += fmt.Sprintf("Response: too long, don't display (%s)\n", url)
		} else {
			output += fmt.Sprintf("Response: %s\n", wdaResp)
		}
		log.Println(output)
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

// func (wb wdaBody) setTextToType(text string) (body wdaBody) {
// 	return wb.set("value", strings.Split(text, ""))
// }

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

// func (wdaResp wdaResponse) getValue2Bytes() []byte {
// 	return []byte(wdaResp.getValue().Raw)
// }

// func (wdaResp wdaResponse) isReady() bool {
// 	return gjson.GetBytes(wdaResp, "value.ready").Bool()
// }

// func (wdaResp wdaResponse) getElementID() gjson.Result {
// 	return wdaResp.getValue().get
// }

// func (wdaResp wdaResponse) getSessionID() (sessionID string, err error) {
// 	sessionID = wdaResp.getByPath("sessionId").String()
// 	if sessionID == "" {
// 		err = errors.New("not find sessionId")
// 	}
// 	return
// }

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
