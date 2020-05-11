package gwda

import (
	"fmt"
	"testing"
)

func TestTmpGwda(t *testing.T) {
	// elemAttr := WDAElementAttribute{Label: "通知"}
	// fmt.Println(elemAttr)
	// fmt.Println(elemAttr.getMethodAndValue())

	// fmt.Println(NewWDAElementAttribute().SetUID("计算器"))
	fmt.Printf("%v\n", NewWDAElementAttribute().SetLabel("计算器"))

	// 以 struct 的形式可以更好的保证参数传入的正确性
	// 虽然只是为了获取字符串，但是形式参数使用字符串相对容易出错
	// type WDAElementType2 struct {
	// 	Alert string
	// }
	// _globalDefault := WDAElementType2{Alert: "XCUIElementTypeAlert"}
	// fmt.Println(_globalDefault.Alert)

	// 暂无想到可保持 NewXXX 使用风格的方式
	fmt.Println(WDAElementType{Alert: true})

	fmt.Println(WDAElementType{})
	fmt.Println(WDAElementAttribute{})
	fmt.Println(WDAElementAttribute{}.SetLabel("TestFlight"))

	fmt.Println()

	// using, value := WDALocator{ClassName: WDAElementType{Application: true}}.getUsingAndValue()
	// using, value := WDALocator{Name: "App Store"}.getUsingAndValue()
	// using, value := WDALocator{AccessibilityId: "设置"}.getUsingAndValue()
	using, value := WDALocator{LinkText: NewWDAElementAttribute().SetLabel("TestFlight")}.getUsingAndValue()
	fmt.Println("using:", using)
	fmt.Println("value:", value)

	fmt.Println()
}

// type WDAElementAttribute struct {
// 	UID                    string `json:"UID"`                    // Element's unique identifier
// 	AccessibilityContainer bool   `json:"accessibilityContainer"` // Whether element is an accessibility container (contains children of any depth that are accessible)
// 	Accessible             bool   `json:"accessible"`             // Whether element is accessible
// 	Enabled                bool   `json:"enabled"`                // Whether element is enabled
// 	Label                  string `json:"label"`                  // Element's label
// 	Name                   string `json:"name"`                   // Element's name
// 	Selected               bool   `json:"selected"`               // Element's selected state
// 	Type                   string `json:"type"`                   // Element's type
// 	Value                  string `json:"value"`                  // Element's value
// 	Visible                bool   `json:"visible"`                // Whether element is visible
// 	// Frame                  string `json:"wdFrame"`                    // Element's frame in CGRect format
// 	// Rect                   string `json:"wdRect"`                     // Element's frame in NSDictionary format
//
// 	// WdUID                    string `json:"wdUID"`
// 	// WdAccessibilityContainer bool `json:"wdAccessibilityContainer"`
// 	// WdAccessible             bool `json:"wdAccessible"`
// 	// WdEnabled                bool `json:"wdEnabled"`
// 	// WdLabel                  string `json:"wdLabel"`
// 	// WdName                   string `json:"wdName"`
// 	// WdSelected               bool `json:"wdSelected"`
// 	// WdType                   string `json:"wdType"`
// 	// WdValue                  string `json:"wdValue"`
// 	// WdVisible                bool `json:"wdVisible"`
// 	// WdFrame                  string `json:"wdFrame"`
// 	// WdRect                   string `json:"wdRect"`
// }
