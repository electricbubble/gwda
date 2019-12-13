package gwda

import (
	"net/url"
)

type Element struct {
	elementURL *url.URL
}

func (e *Element) Click() error {
	wdaResp, err := internalPost("Click", urlJoin(e.elementURL, "click"), nil)
	if err != nil {
		return err
	}
	return wdaResp.getErrMsg()
}

func (e *Element) Rect() (sJson string, err error) {
	wdaResp, err := internalGet("GetRect", urlJoin(e.elementURL, "rect"))
	if err != nil {
		return "", err
	}
	return wdaResp.getValue().String(), nil
}

func (e *Element) Enabled() (isEnabled bool, err error) {
	wdaResp, err := internalGet("GetRect", urlJoin(e.elementURL, "enabled"))
	if err != nil {
		return false, err
	}
	return wdaResp.getValue().Bool(), nil
}
