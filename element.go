package gwda

import (
	"encoding/json"
	"net/url"
)

type Element struct {
	elementURL *url.URL
}

type Position struct {
	Y int `json:"y"`
	X int `json:"x"`
}

type WDARect struct {
	// Y int `json:"y"`
	// X int `json:"x"`
	Position
	WDASize
}

func (e *Element) Click() error {
	wdaResp, err := internalPost("Click", urlJoin(e.elementURL, "click"), nil)
	if err != nil {
		return err
	}
	return wdaResp.getErrMsg()
}

func (e *Element) Rect() (wdaRect WDARect, err error) {
	wdaResp, err := internalGet("Rect", urlJoin(e.elementURL, "rect"))
	if err != nil {
		return WDARect{}, err
	}
	wdaRect._String = wdaResp.getValue().String()
	err = json.Unmarshal([]byte(wdaRect._String), &wdaRect)
	return
}

func (e *Element) Enabled() (isEnabled bool, err error) {
	wdaResp, err := internalGet("Enabled", urlJoin(e.elementURL, "enabled"))
	if err != nil {
		return false, err
	}
	return wdaResp.getValue().Bool(), nil
}
