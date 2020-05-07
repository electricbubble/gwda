package gwda

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"math"
	"net/url"
	"path"
	"reflect"
	"strconv"
)

type Element struct {
	endpoint *url.URL
	UID      string
}

func newElement(endpoint *url.URL, elemUID string) (elem *Element) {
	elem = new(Element)
	elem.endpoint = endpoint
	elem.UID = elemUID
	return
}

// /element/:uuid
func (e *Element) _withFormat(elem ...string) string {
	return path.Join(append([]string{"element", e.UID}, elem...)...)
}

// http://ip:port/session/:uuid/element/:uuid
func (e *Element) _withFormatToUrl(elem ...string) *url.URL {
	tmp, _ := url.Parse(urlJoin(e.endpoint, e._withFormat()))
	path.Join(append([]string{"element", e.UID}, elem...)...)
	return tmp
}

func (e *Element) Tap(x, y int) error {
	return tap(e.endpoint, x, y, e.UID)
}

func (e *Element) TapFloat(x, y float64) error {
	return tap(e.endpoint, x, y, e.UID)
}

// DoubleTap
//
// Sends a double tap event to a hittable point computed for the element.
func (e *Element) DoubleTap() error {
	return doubleTap(e.endpoint, -1, -1, e._withFormat())
}

// TwoFingerTap
//
// Sends a two finger tap event to a hittable point computed for the element.
func (e *Element) TwoFingerTap() (err error) {
	// [FBRoute POST:@"/wda/element/:uuid/twoFingerTap"]
	_, err = internalPost("TwoFingerTap", urlJoin(e.endpoint, e._withFormat("/twoFingerTap"), true), nil)
	return
}

// TapWithNumberOfTaps
//
// Sends one or more taps with one of more touch points.
func (e *Element) TapWithNumberOfTaps(numberOfTaps, numberOfTouches int) (err error) {
	if numberOfTouches <= 0 {
		return errors.New("'numberOfTouches' must be greater than zero")
	}
	if numberOfTouches > 5 {
		return errors.New("'numberOfTouches' cannot be greater than 5")
	}
	if numberOfTaps <= 0 {
		return errors.New("'numberOfTaps' must be greater than zero")
	}
	if numberOfTaps > 10 {
		return errors.New("'numberOfTaps' cannot be greater than 10")
	}
	body := newWdaBody().set("numberOfTaps", numberOfTaps).set("numberOfTouches", numberOfTouches)
	// [FBRoute POST:@"/wda/element/:uuid/tapWithNumberOfTaps"]
	_, err = internalPost("TapWithNumberOfTaps", urlJoin(e.endpoint, e._withFormat("/tapWithNumberOfTaps"), true), body)
	return
}

// TouchAndHold
//
// Sends a long press gesture to a hittable point computed for the element, holding for the specified duration.
func (e *Element) TouchAndHold(duration ...int) (err error) {
	if len(duration) == 0 {
		duration = []int{1}
	}
	return touchAndHold(e.endpoint, -1, -1, duration[0], e._withFormat())
}

func (e *Element) TouchAndHoldFloat(duration ...float64) (err error) {
	if len(duration) == 0 {
		duration = []float64{1.0}
	}
	return touchAndHold(e.endpoint, -1, -1, duration[0], e._withFormat())
}

func (e *Element) _forceTouch(wdaCoordinate WDACoordinate, pressure float64, duration ...float64) (err error) {
	body := newWdaBody()
	if wdaCoordinate.X != -1 && wdaCoordinate.Y != -1 {
		body.setXY(wdaCoordinate.X, wdaCoordinate.Y)
	}
	body.set("pressure", pressure)
	if len(duration) == 0 {
		duration = []float64{1.0}
	}
	body.set("duration", duration[0])
	// [FBRoute POST:@"/wda/element/:uuid/forceTouch"]
	_, err = internalPost("ForceTouch", urlJoin(e.endpoint, e._withFormat("/forceTouch"), true), body)
	return
}

// ForceTouch
//
// 3D Touch
func (e *Element) ForceTouch(pressure float64, duration ...float64) (err error) {
	return e._forceTouch(WDACoordinate{X: -1, Y: -1}, pressure, duration...)
}

func (e *Element) ForceTouchCoordinate(wdaCoordinate WDACoordinate, pressure float64, duration ...float64) (err error) {
	return e._forceTouch(wdaCoordinate, pressure, duration...)
}

// WDATouchSensitivity [Light, Medium, Firm]
//  reveal content previews, actions,	[peek] 轻瞄
//  contextual menus				[pop] 突显
//  显示内容预览、操作和上下文菜单

// func (e *Element) ForceTouchPeek() (err error) {
// 	// return e.ForceTouch(0.5) // Light
// 	return e.ForceTouch(0.68) // Medium
// 	// return e.ForceTouch(0.68) // Firm
// }
//
// func (e *Element) ForceTouchPop() (err error) {
// 	// return e.ForceTouch(2.3, 1.2) // Light
// 	return e.ForceTouch(5, 0.5) // Light
//
// 	// return e.ForceTouch(0.xx) // Medium
// 	// return e.ForceTouch(0.xx) // Firm
// }

// Drag
//
// Clicks and holds for a specified duration (generally long enough to start a drag operation) then drags to the other coordinate.
func (e *Element) Drag(fromX, fromY, toX, toY int, pressForDuration ...int) (err error) {
	if len(pressForDuration) == 0 {
		pressForDuration = []int{1}
	}
	return drag(e.endpoint, fromX, fromY, toX, toY, pressForDuration[0], e._withFormat())
}

func (e *Element) DragFloat(fromX, fromY, toX, toY float64, pressForDuration ...float64) (err error) {
	if len(pressForDuration) == 0 {
		pressForDuration = []float64{1}
	}
	return drag(e.endpoint, fromX, fromY, toX, toY, pressForDuration[0], e._withFormat())
}

type WDASwipeDirection string

const (
	WDASwipeDirectionUp    WDASwipeDirection = "up"
	WDASwipeDirectionDown  WDASwipeDirection = "down"
	WDASwipeDirectionLeft  WDASwipeDirection = "left"
	WDASwipeDirectionRight WDASwipeDirection = "right"
)

// Swipe
//
//	element.frame.origin.x + [request.arguments[@"fromX"] doubleValue]
// 	element.frame.origin.y + [request.arguments[@"fromY"] doubleValue]
// 	element.frame.origin.x + [request.arguments[@"toX"] doubleValue]
//	element.frame.origin.y + [request.arguments[@"toY"] doubleValue]
func (e *Element) Swipe(fromX, fromY, toX, toY int) (err error) {
	return drag(e.endpoint, fromX, fromY, toX, toY, 0, e._withFormat())
}

func (e *Element) SwipeFloat(fromX, fromY, toX, toY float64) (err error) {
	return drag(e.endpoint, fromX, fromY, toX, toY, 0, e._withFormat())
}

// SwipeDirection
//
// Sends a swipe gesture in the specified direction.
func (e *Element) SwipeDirection(direction WDASwipeDirection) (err error) {
	body := newWdaBody().set("direction", direction)
	// [FBRoute POST:@"/wda/element/:uuid/swipe"]
	_, err = internalPost("SwipeDirection", urlJoin(e.endpoint, e._withFormat("/swipe"), true), body)
	return
}

// SwipeUp
//
// Sends a swipe-up gesture.
func (e *Element) SwipeUp() (err error) {
	return e.SwipeDirection(WDASwipeDirectionUp)
}

// SwipeDown
//
// Sends a swipe-down gesture.
func (e *Element) SwipeDown() (err error) {
	return e.SwipeDirection(WDASwipeDirectionDown)
}

// SwipeLeft
//
// Sends a swipe-left gesture.
func (e *Element) SwipeLeft() (err error) {
	return e.SwipeDirection(WDASwipeDirectionLeft)
}

// SwipeRight
//
// Sends a swipe-right gesture.
func (e *Element) SwipeRight() (err error) {
	return e.SwipeDirection(WDASwipeDirectionRight)
}

// Pinch
//
// Sends a pinching gesture with two touches.
//
// The system makes a best effort to synthesize the requested scale and velocity: absolute accuracy is not guaranteed.
// Some values may not be possible based on the size of the element's frame - these will result in test failures.
//
// @param scale
// The scale of the pinch gesture.  Use a scale between 0 and 1 to "pinch close" or zoom out and a scale greater than 1 to "pinch open" or zoom in.
//
// @param velocity
// The velocity of the pinch in scale factor per second.
func (e *Element) Pinch(scale, velocity float64) (err error) {
	if scale <= 0 {
		return errors.New("'scale' must be greater than zero")
	}
	if scale == 1 {
		return errors.New("'scale' must be greater or less than 1")
	}
	if scale < 1 && velocity > 0 {
		return errors.New("'velocity' must be less than zero when 'scale' is less than 1")
	}
	if scale > 1 && velocity <= 0 {
		return errors.New("'velocity' must be greater than zero when 'scale' is greater than 1")
	}
	body := newWdaBody().set("scale", scale).set("velocity", velocity)
	// [FBRoute POST:@"/wda/element/:uuid/pinch"]
	_, err = internalPost("Pinch", urlJoin(e.endpoint, e._withFormat("/pinch"), true), body)
	return
}

// PinchToZoomIn
//
// scale, velocity = 2, 10
func (e *Element) PinchToZoomIn() (err error) {
	return e.Pinch(2, 10)
}

// PinchToZoomOut
//
// !! may not work
//
// scale, velocity = 0.9, -4.5
func (e *Element) PinchToZoomOut() (err error) {
	return e.Pinch(0.9, -4.5)
}

func (e *Element) PinchToZoomOutByActions(scale ...float64) (err error) {
	if len(scale) == 0 {
		scale = []float64{1.0}
	} else if scale[0] > 23 {
		scale[0] = 23
	}
	var rect WDARect
	if rect, err = e.Rect(); err != nil {
		return err
	}
	r := scale[0] * 2 / 100.0
	leftX, leftY := float64(rect.Width)*r, float64(rect.Height)*r

	offsetX, offsetY := int(leftX), int(leftY)

	actions := NewWDAActions().FingerSwipe(0-offsetX, 0-offsetY, 0, 0, e).FingerSwipe(offsetX, offsetY, 0, 0, e)
	return performActions(e.endpoint, actions)
}

// Rotate
//
// Sends a rotation gesture with two touches.
//
// The system makes a best effort to synthesize the requested rotation and velocity: absolute accuracy is not guaranteed.
// Some values may not be possible based on the size of the element's frame - these will result in test failures.
//
// @param rotation
// The rotation of the gesture in radians.
//
// @param velocity
// The velocity of the rotation gesture in radians per second.
func (e *Element) Rotate(rotation float64, velocity ...float64) (err error) {
	if rotation > math.Pi*2 || rotation < math.Pi*-2 {
		return errors.New("'rotation' must not be more than 2π or less than -2π")
	}
	if len(velocity) == 0 || velocity[0] == 0 {
		velocity = []float64{rotation}
	}
	if rotation > 0 && velocity[0] < 0 || rotation < 0 && velocity[0] > 0 {
		return errors.New("'rotation' and 'velocity' must have the same sign")
	}
	body := newWdaBody().set("rotation", rotation).set("velocity", velocity[0])
	// [FBRoute POST:@"/wda/element/:uuid/rotate"]
	_, err = internalPost("Rotate", urlJoin(e.endpoint, e._withFormat("/rotate"), true), body)
	return
}

func (e *Element) _scroll(body wdaBody) (err error) {
	// [FBRoute POST:@"/wda/element/:uuid/scroll"]
	_, err = internalPost("Scroll", urlJoin(e.endpoint, e._withFormat("/scroll"), true), body)
	return
}

// ScrollElementByName
func (e *Element) ScrollElementByName(name string) (err error) {
	return e._scroll(newWdaBody().set("name", name))
}

func (e *Element) ScrollElementByPredicate(predicate string) (err error) {
	return e._scroll(newWdaBody().set("predicateString", predicate))
}

func (e *Element) ScrollToVisible() (err error) {
	return e._scroll(newWdaBody().set("toVisible", true))
}

func (e *Element) _scrollDirection(direction WDASwipeDirection, distance ...float64) (err error) {
	if len(distance) == 0 {
		distance = []float64{0.5}
	}
	body := newWdaBody().set("direction", direction).set("distance", distance[0])
	return e._scroll(body)
}

func (e *Element) ScrollUp(distance ...float64) (err error) {
	return e._scrollDirection(WDASwipeDirectionUp, distance...)
}

func (e *Element) ScrollDown(distance ...float64) (err error) {
	return e._scrollDirection(WDASwipeDirectionDown, distance...)
}

func (e *Element) ScrollLeft(distance ...float64) (err error) {
	return e._scrollDirection(WDASwipeDirectionLeft, distance...)
}

func (e *Element) ScrollRight(distance ...float64) (err error) {
	return e._scrollDirection(WDASwipeDirectionRight, distance...)
}

type WDAPickerWheelSelectOrder string

const (
	WDAPickerWheelSelectOrderNext     WDAPickerWheelSelectOrder = "next"
	WDAPickerWheelSelectOrderPrevious WDAPickerWheelSelectOrder = "previous"
)

func (e *Element) PickerWheelSelect(order WDAPickerWheelSelectOrder, offset ...int) (err error) {
	if len(offset) == 0 {
		offset = []int{2}
	} else if offset[0] <= 0 || offset[0] > 5 {
		return errors.New(fmt.Sprintf("'offset' value is expected to be in range (0, 5]. '%d' was given instead", offset[0]))
	}
	body := newWdaBody().set("order", order).set("offset", float64(offset[0])*0.1)
	// [FBRoute POST:@"/wda/pickerwheel/:uuid/select"]
	_, err = internalPost("PickerWheelSelect", urlJoin(e.endpoint, path.Join("/pickerwheel", e.UID, "/select"), true), body)
	return
}

func (e *Element) PickerWheelSelectNext(offset ...int) (err error) {
	if len(offset) == 0 {
		offset = []int{1}
	}
	return e.PickerWheelSelect(WDAPickerWheelSelectOrderNext, offset...)
}

func (e *Element) PickerWheelSelectPrevious(offset ...int) (err error) {
	if len(offset) == 0 {
		offset = []int{1}
	}
	return e.PickerWheelSelect(WDAPickerWheelSelectOrderPrevious, offset...)
}

func (e *Element) Click() (err error) {
	// [FBRoute POST:@"/element/:uuid/click"]
	_, err = internalPost("Click", urlJoin(e.endpoint, e._withFormat("/click")), nil)
	return
}

func (e *Element) SendKeys(text string, typingFrequency ...int) error {
	// [FBRoute POST:@"/element/:uuid/value"]
	return sendKeys(urlJoin(e.endpoint, e._withFormat("/value")), text, typingFrequency...)
}

func (e *Element) Clear() (err error) {
	// [FBRoute POST:@"/element/:uuid/clear"]
	_, err = internalPost("Clear", urlJoin(e.endpoint, e._withFormat("/clear")), nil)
	return
}

type WDACoordinate struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type WDARect struct {
	WDACoordinate
	WDASize
}

func (e *Element) Rect() (wdaRect WDARect, err error) {
	var wdaResp wdaResponse
	// [FBRoute GET:@"/element/:uuid/rect"]
	if wdaResp, err = internalGet("Rect", urlJoin(e.endpoint, e._withFormat("/rect"))); err != nil {
		return WDARect{}, err
	}
	wdaRect._string = wdaResp.getValue().String()
	err = json.Unmarshal([]byte(wdaRect._string), &wdaRect)
	return
}

func (e *Element) IsEnabled() (isEnabled bool, err error) {
	var wdaResp wdaResponse
	// [FBRoute GET:@"/element/:uuid/enabled"]
	if wdaResp, err = internalGet("IsEnabled", urlJoin(e.endpoint, e._withFormat("/enabled"))); err != nil {
		return false, err
	}
	return wdaResp.getValue().Bool(), nil
}

func (e *Element) IsDisplayed() (isDisplayed bool, err error) {
	var wdaResp wdaResponse
	// [FBRoute GET:@"/element/:uuid/displayed"]
	if wdaResp, err = internalGet("IsDisplayed", urlJoin(e.endpoint, e._withFormat("/displayed"))); err != nil {
		return false, err
	}
	return wdaResp.getValue().Bool(), nil
}

func (e *Element) IsSelected() (isSelected bool, err error) {
	var wdaResp wdaResponse
	// [FBRoute GET:@"/element/:uuid/selected"]
	if wdaResp, err = internalGet("IsSelected", urlJoin(e.endpoint, e._withFormat("/selected"))); err != nil {
		return false, err
	}
	return wdaResp.getValue().Bool(), nil
}

func (e *Element) IsAccessible() (isAccessible bool, err error) {
	var wdaResp wdaResponse
	// [FBRoute GET:@"/wda/element/:uuid/accessible"]
	if wdaResp, err = internalGet("IsAccessible", urlJoin(e.endpoint, e._withFormat("/accessible"), true)); err != nil {
		return false, err
	}
	return wdaResp.getValue().Bool(), nil
}

func (e *Element) IsAccessibilityContainer() (isAccessibilityContainer bool, err error) {
	var wdaResp wdaResponse
	// [FBRoute GET:@"/wda/element/:uuid/accessibilityContainer"]
	if wdaResp, err = internalGet("IsAccessibilityContainer", urlJoin(e.endpoint, e._withFormat("/accessibilityContainer"), true)); err != nil {
		return false, err
	}
	return wdaResp.getValue().Bool(), nil
}

// GetAttribute
//
// Returns value of given property specified in WebDriver Spec
// Check the FBElement protocol to get list of supported attributes.
// This method also supports shortcuts, like wdName == name, wdValue == value.
func (e *Element) GetAttribute(attr WDAElementAttribute) (value string, err error) {
	attrName := attr.getAttributeName()
	if attrName == "UNKNOWN" {
		return "", errors.New("'WDAElementAttribute' does not have 'Attribute Name'")
	}
	var wdaResp wdaResponse
	// [FBRoute GET:@"/element/:uuid/attribute/:name"]
	if wdaResp, err = internalGet("GetAttribute", urlJoin(e.endpoint, e._withFormat("/attribute", attrName))); err != nil {
		return "", err
	}
	return wdaResp.getValue().String(), nil
}

func (e *Element) Text() (text string, err error) {
	var wdaResp wdaResponse
	// [FBRoute GET:@"/element/:uuid/text"]
	if wdaResp, err = internalGet("Text", urlJoin(e.endpoint, e._withFormat("/text"))); err != nil {
		return "", err
	}
	return wdaResp.getValue().String(), nil
}

// Type
//
// Element's type ( WDAElementType )
func (e *Element) Type() (elemType string, err error) {
	var wdaResp wdaResponse
	// [FBRoute GET:@"/element/:uuid/name"]
	if wdaResp, err = internalGet("Type", urlJoin(e.endpoint, e._withFormat("/name"))); err != nil {
		return "", err
	}
	return wdaResp.getValue().String(), nil
}

// FindElement
func (e *Element) FindElement(wdaLocator WDALocator) (element *Element, err error) {
	var elemUID string
	// [FBRoute POST:@"/element/:uuid/element"]
	if elemUID, err = findUidOfElement(e._withFormatToUrl(), wdaLocator); err != nil {
		return nil, err
	}
	return newElement(e.endpoint, elemUID), nil
}

// FindElements
func (e *Element) FindElements(wdaLocator WDALocator) (elements []*Element, err error) {
	var elemUIDs []string
	// [FBRoute POST:@"/element/:uuid/elements"]
	if elemUIDs, err = findUidOfElements(e._withFormatToUrl(), wdaLocator); err != nil {
		return nil, err
	}
	elements = make([]*Element, len(elemUIDs))
	for i := range elements {
		elements[i] = newElement(e.endpoint, elemUIDs[i])
	}
	return
}

// FindVisibleCells
func (e *Element) FindVisibleCells() (elements []*Element, err error) {
	var wdaResp wdaResponse
	// [FBRoute GET:@"/wda/element/:uuid/getVisibleCells"]
	if wdaResp, err = internalGet("FindVisibleCells", urlJoin(e.endpoint, e._withFormat("/getVisibleCells"), true)); err != nil {
		return nil, err
	}
	results := wdaResp.getValue().Array()
	if len(results) == 0 {
		return nil, errors.New(fmt.Sprintf("no such element: unable to find a cell element in this element"))
	}
	elements = make([]*Element, len(results))
	for i := range elements {
		elements[i] = newElement(e.endpoint, results[i].Get("ELEMENT").String())
	}
	return
}

// W3C element screenshot
// [[FBRoute GET:@"/element/:uuid/screenshot"] respondWithTarget:self action:@selector(handleElementScreenshot:)],
// JSONWP element screenshot
// [[FBRoute GET:@"/screenshot/:uuid"] respondWithTarget:self action:@selector(handleElementScreenshot:)],

// Screenshot
func (e *Element) Screenshot() (raw *bytes.Buffer, err error) {
	// [FBRoute GET:@"/element/:uuid/screenshot"]
	return screenshot(e._withFormatToUrl())
}

// ScreenshotToDisk
func (e *Element) ScreenshotToDisk(filename string) (err error) {
	return screenshotToDisk(e._withFormatToUrl(), filename)
}

// ScreenshotToImage
func (e *Element) ScreenshotToImage() (img image.Image, format string, err error) {
	return screenshotToImage(e._withFormatToUrl())
}

// func addToRootWda(baseUrl *url.URL) {
// fmt.Println(url.Parse(urlJoin(baseUrl, "/accessible")))
// }

func (e *Element) tttTmp() {

	var err error
	// body := newWdaBody()
	// _ = body

	body := newWdaBody()

	// element, _ := e.FindElement(WDALocator{Predicate: "type == 'XCUIElementTypeCell' AND name == '音乐'"})
	// attribute, _ := element.GetAttribute(NewWDAElementAttribute().SetUID(""))
	// fmt.Println("###", attribute)

	// (0.0, 0.5]
	// 'offset' value is expected to be in range (0.0, 0.5]. '0.0' was given instead
	body.set("offset", 0.3)

	body.set("order", "next")
	body.set("order", "previous")

	var wdaResp wdaResponse
	// [FBRoute POST:@"/wda/pickerwheel/:uuid/select"]
	wdaResp, err = internalPost("###############", urlJoin(e.endpoint, path.Join("/pickerwheel", e.UID, "/select"), true), body)
	_ = wdaResp
	_ = err
	// fmt.Println(err, wdaResp)
}

type WDALocator struct {
	ClassName WDAElementType `json:"class name"`

	// isSearchByIdentifier
	Name            string `json:"name"`
	Id              string `json:"id"`
	AccessibilityId string `json:"accessibility id"`
	// isSearchByIdentifier

	// partialSearch
	LinkText        WDAElementAttribute `json:"link text"`
	PartialLinkText WDAElementAttribute `json:"partial link text"`
	// partialSearch

	Predicate string `json:"predicate string"`

	ClassChain string `json:"class chain"`

	XPath string `json:"xpath"`
}

func (wl WDALocator) getUsingAndValue() (using, value string) {
	vBy := reflect.ValueOf(wl)
	tBy := reflect.TypeOf(wl)
	for i := 0; i < vBy.NumField(); i++ {
		vi := vBy.Field(i).Interface()
		switch vi.(type) {
		case WDAElementType:
			value = vi.(WDAElementType).String()
		case string:
			value = vi.(string)
		case WDAElementAttribute:
			value = vi.(WDAElementAttribute).String()
		}
		if value != "" && value != "UNKNOWN" {
			using = tBy.Field(i).Tag.Get("json")
			return
		}
	}
	return
}

type WDAElementAttribute wdaBody

func (ea WDAElementAttribute) String() string {
	for k, v := range ea {
		switch v.(type) {
		case bool:
			return k + "=" + strconv.FormatBool(v.(bool))
		case string:
			return k + "=" + v.(string)
		default:
			return k + "=" + fmt.Sprintf("%v", v)
		}
	}
	return "UNKNOWN"
}

func (ea WDAElementAttribute) getAttributeName() string {
	for k := range ea {
		return k
	}
	return "UNKNOWN"
}

func NewWDAElementAttribute() WDAElementAttribute {
	return make(WDAElementAttribute)
}

// SetUID
//
// Element's unique identifier
func (ea WDAElementAttribute) SetUID(uid string) WDAElementAttribute {
	return WDAElementAttribute(wdaBody(ea).set("UID", uid))
}

// SetAccessibilityContainer
//
// Whether element is an accessibility container (contains children of any depth that are accessible)
func (ea WDAElementAttribute) SetAccessibilityContainer(b bool) WDAElementAttribute {
	return WDAElementAttribute(wdaBody(ea).set("accessibilityContainer", b))
}

// SetAccessible
//
// Whether element is accessible
func (ea WDAElementAttribute) SetAccessible(b bool) WDAElementAttribute {
	return WDAElementAttribute(wdaBody(ea).set("accessible", b))
}

// SetEnabled
//
// Whether element is enabled
func (ea WDAElementAttribute) SetEnabled(b bool) WDAElementAttribute {
	return WDAElementAttribute(wdaBody(ea).set("enabled", b))
}

// SetLabel
//
// Element's label
func (ea WDAElementAttribute) SetLabel(s string) WDAElementAttribute {
	return WDAElementAttribute(wdaBody(ea).set("label", s))
}

// SetName
//
// Element's name
func (ea WDAElementAttribute) SetName(s string) WDAElementAttribute {
	return WDAElementAttribute(wdaBody(ea).set("name", s))
}

// SetSelected
//
// Element's selected state
func (ea WDAElementAttribute) SetSelected(b bool) WDAElementAttribute {
	return WDAElementAttribute(wdaBody(ea).set("selected", b))
}

// SetType
//
// Element's type
func (ea WDAElementAttribute) SetType(elemType WDAElementType) WDAElementAttribute {
	return WDAElementAttribute(wdaBody(ea).set("type", elemType.String()))
}

// SetValue
//
// Element's value
func (ea WDAElementAttribute) SetValue(s string) WDAElementAttribute {
	return WDAElementAttribute(wdaBody(ea).set("value", s))
}

// SetVisible
//
// Whether element is visible
func (ea WDAElementAttribute) SetVisible(b bool) WDAElementAttribute {
	return WDAElementAttribute(wdaBody(ea).set("visible", b))
}

func (et WDAElementType) String() string {
	vBy := reflect.ValueOf(et)
	tBy := reflect.TypeOf(et)
	for i := 0; i < vBy.NumField(); i++ {
		if vBy.Field(i).Bool() {
			return tBy.Field(i).Tag.Get("json")
		}
	}
	return "UNKNOWN"
}

// WDAElementType
// !!! This mapping should be updated if there are changes after each new XCTest release"`
type WDAElementType struct {
	Any                bool `json:"XCUIElementTypeAny"`
	Other              bool `json:"XCUIElementTypeOther"`
	Application        bool `json:"XCUIElementTypeApplication"`
	Group              bool `json:"XCUIElementTypeGroup"`
	Window             bool `json:"XCUIElementTypeWindow"`
	Sheet              bool `json:"XCUIElementTypeSheet"`
	Drawer             bool `json:"XCUIElementTypeDrawer"`
	Alert              bool `json:"XCUIElementTypeAlert"`
	Dialog             bool `json:"XCUIElementTypeDialog"`
	Button             bool `json:"XCUIElementTypeButton"`
	RadioButton        bool `json:"XCUIElementTypeRadioButton"`
	RadioGroup         bool `json:"XCUIElementTypeRadioGroup"`
	CheckBox           bool `json:"XCUIElementTypeCheckBox"`
	DisclosureTriangle bool `json:"XCUIElementTypeDisclosureTriangle"`
	PopUpButton        bool `json:"XCUIElementTypePopUpButton"`
	ComboBox           bool `json:"XCUIElementTypeComboBox"`
	MenuButton         bool `json:"XCUIElementTypeMenuButton"`
	ToolbarButton      bool `json:"XCUIElementTypeToolbarButton"`
	Popover            bool `json:"XCUIElementTypePopover"`
	Keyboard           bool `json:"XCUIElementTypeKeyboard"`
	Key                bool `json:"XCUIElementTypeKey"`
	NavigationBar      bool `json:"XCUIElementTypeNavigationBar"`
	TabBar             bool `json:"XCUIElementTypeTabBar"`
	TabGroup           bool `json:"XCUIElementTypeTabGroup"`
	Toolbar            bool `json:"XCUIElementTypeToolbar"`
	StatusBar          bool `json:"XCUIElementTypeStatusBar"`
	Table              bool `json:"XCUIElementTypeTable"`
	TableRow           bool `json:"XCUIElementTypeTableRow"`
	TableColumn        bool `json:"XCUIElementTypeTableColumn"`
	Outline            bool `json:"XCUIElementTypeOutline"`
	OutlineRow         bool `json:"XCUIElementTypeOutlineRow"`
	Browser            bool `json:"XCUIElementTypeBrowser"`
	CollectionView     bool `json:"XCUIElementTypeCollectionView"`
	Slider             bool `json:"XCUIElementTypeSlider"`
	PageIndicator      bool `json:"XCUIElementTypePageIndicator"`
	ProgressIndicator  bool `json:"XCUIElementTypeProgressIndicator"`
	ActivityIndicator  bool `json:"XCUIElementTypeActivityIndicator"`
	SegmentedControl   bool `json:"XCUIElementTypeSegmentedControl"`
	Picker             bool `json:"XCUIElementTypePicker"`
	PickerWheel        bool `json:"XCUIElementTypePickerWheel"`
	Switch             bool `json:"XCUIElementTypeSwitch"`
	Toggle             bool `json:"XCUIElementTypeToggle"`
	Link               bool `json:"XCUIElementTypeLink"`
	Image              bool `json:"XCUIElementTypeImage"`
	Icon               bool `json:"XCUIElementTypeIcon"`
	SearchField        bool `json:"XCUIElementTypeSearchField"`
	ScrollView         bool `json:"XCUIElementTypeScrollView"`
	ScrollBar          bool `json:"XCUIElementTypeScrollBar"`
	StaticText         bool `json:"XCUIElementTypeStaticText"`
	TextField          bool `json:"XCUIElementTypeTextField"`
	SecureTextField    bool `json:"XCUIElementTypeSecureTextField"`
	DatePicker         bool `json:"XCUIElementTypeDatePicker"`
	TextView           bool `json:"XCUIElementTypeTextView"`
	Menu               bool `json:"XCUIElementTypeMenu"`
	MenuItem           bool `json:"XCUIElementTypeMenuItem"`
	MenuBar            bool `json:"XCUIElementTypeMenuBar"`
	MenuBarItem        bool `json:"XCUIElementTypeMenuBarItem"`
	Map                bool `json:"XCUIElementTypeMap"`
	WebView            bool `json:"XCUIElementTypeWebView"`
	IncrementArrow     bool `json:"XCUIElementTypeIncrementArrow"`
	DecrementArrow     bool `json:"XCUIElementTypeDecrementArrow"`
	Timeline           bool `json:"XCUIElementTypeTimeline"`
	RatingIndicator    bool `json:"XCUIElementTypeRatingIndicator"`
	ValueIndicator     bool `json:"XCUIElementTypeValueIndicator"`
	SplitGroup         bool `json:"XCUIElementTypeSplitGroup"`
	Splitter           bool `json:"XCUIElementTypeSplitter"`
	RelevanceIndicator bool `json:"XCUIElementTypeRelevanceIndicator"`
	ColorWell          bool `json:"XCUIElementTypeColorWell"`
	HelpTag            bool `json:"XCUIElementTypeHelpTag"`
	Matte              bool `json:"XCUIElementTypeMatte"`
	DockItem           bool `json:"XCUIElementTypeDockItem"`
	Ruler              bool `json:"XCUIElementTypeRuler"`
	RulerMarker        bool `json:"XCUIElementTypeRulerMarker"`
	Grid               bool `json:"XCUIElementTypeGrid"`
	LevelIndicator     bool `json:"XCUIElementTypeLevelIndicator"`
	Cell               bool `json:"XCUIElementTypeCell"`
	LayoutArea         bool `json:"XCUIElementTypeLayoutArea"`
	LayoutItem         bool `json:"XCUIElementTypeLayoutItem"`
	Handle             bool `json:"XCUIElementTypeHandle"`
	Stepper            bool `json:"XCUIElementTypeStepper"`
	Tab                bool `json:"XCUIElementTypeTab"`
	TouchBar           bool `json:"XCUIElementTypeTouchBar"`
	StatusItem         bool `json:"XCUIElementTypeStatusItem"`
}
