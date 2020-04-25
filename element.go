package gwda

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
)

type Element struct {
	elementURL *url.URL
}

type WDAPosition struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type WDARect struct {
	WDAPosition
	WDASize
}

func (e *Element) Click() (err error) {
	_, err = internalPost("Click", urlJoin(e.elementURL, "/click"), nil)
	return
}

func (e *Element) Rect() (wdaRect WDARect, err error) {
	var wdaResp wdaResponse
	if wdaResp, err = internalGet("Rect", urlJoin(e.elementURL, "/rect")); err != nil {
		return WDARect{}, err
	}
	wdaRect._string = wdaResp.getValue().String()
	err = json.Unmarshal([]byte(wdaRect._string), &wdaRect)
	return
}

func (e *Element) IsEnabled() (isEnabled bool, err error) {
	var wdaResp wdaResponse
	if wdaResp, err = internalGet("IsEnabled", urlJoin(e.elementURL, "/enabled")); err != nil {
		return false, err
	}
	return wdaResp.getValue().Bool(), nil
}

func (e *Element) IsDisplayed() (isDisplayed bool, err error) {
	var wdaResp wdaResponse
	if wdaResp, err = internalGet("IsDisplayed", urlJoin(e.elementURL, "/displayed")); err != nil {
		return false, err
	}
	return wdaResp.getValue().Bool(), nil
}

func (e *Element) IsSelected() (isSelected bool, err error) {
	var wdaResp wdaResponse
	if wdaResp, err = internalGet("IsSelected", urlJoin(e.elementURL, "/selected")); err != nil {
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
	if wdaResp, err = internalGet("GetAttribute", urlJoin(e.elementURL, "/attribute", attrName)); err != nil {
		return "", err
	}
	return wdaResp.getValue().String(), nil
}

func (e *Element) Text() (text string, err error) {
	var wdaResp wdaResponse
	if wdaResp, err = internalGet("Text", urlJoin(e.elementURL, "/text")); err != nil {
		return "", err
	}
	return wdaResp.getValue().String(), nil
}

// Type
//
// Element's type ( WDAElementType )
func (e *Element) Type() (elemType string, err error) {
	var wdaResp wdaResponse
	if wdaResp, err = internalGet("Type", urlJoin(e.elementURL, "/name")); err != nil {
		return "", err
	}
	return wdaResp.getValue().String(), nil
}

func (e *Element) tttTmp() {
	// TODO [[FBRoute POST:@"/element/:uuid/element"] respondWithTarget:self action:@selector(handleFindSubElement:)],
	// TODO [[FBRoute POST:@"/element/:uuid/elements"] respondWithTarget:self action:@selector(handleFindSubElements:)],
	// TODO [[FBRoute GET:@"/wda/element/:uuid/getVisibleCells"] respondWithTarget:self action:@selector(handleFindVisibleCells:)],
	body := newWdaBody()
	_ = body

	// attrName := "type"
	// attrName = NewWDAElementAttribute().SetType(WDAElementType{}).GetAttributeName()
	// attrName = WDAElementType{T}
	// tmp, _ := url.Parse(urlJoin(e.elementURL, "/attribute"))
	// q := tmp.Query()
	// q.Set("name", attrName)
	// tmp.RawQuery = q.Encode()
	// wdaResp, err := internalGet("###############", tmp.String())

	wdaResp, err := internalGet("###############", urlJoin(e.elementURL, "/name"))
	fmt.Println(err, wdaResp)
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
	for k, _ := range ea {
		// switch v.(type) {
		// case bool:
		// 	return k + "=" + strconv.FormatBool(v.(bool))
		// case string:
		// 	return k + "=" + v.(string)
		// default:
		// 	return k + "=" + fmt.Sprintf("%v", v)
		// }
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
