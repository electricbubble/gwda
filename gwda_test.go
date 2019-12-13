package gwda

import (
	"fmt"
	"reflect"
	"testing"
)

type selector interface{}
type PartialLinkText struct {
}

type PredicateString struct {
	UID                      string `json:"UID"`
	AccessibilityContainer   string `json:"accessibilityContainer"`
	Accessible               string `json:"accessible"`
	Enabled                  bool   `json:"enabled"`
	Frame                    string `json:"frame"`
	Label                    string `json:"label"`
	Name                     string `json:"name"`
	Rect                     string `json:"rect"`
	Type                     string `json:"type"`
	Value                    string `json:"value"`
	Visible                  string `json:"visible"`
	WdAccessibilityContainer string `json:"wdAccessibilityContainer"`
	WdAccessible             string `json:"wdAccessible"`
	WdEnabled                bool   `json:"wdEnabled"`
	WdFrame                  string `json:"wdFrame"`
	WdLabel                  string `json:"wdLabel"`
	WdName                   string `json:"wdName"`
	WdRect                   string `json:"wdRect"`
	WdType                   string `json:"wdType"`
	WdUID                    string `json:"wdUID"`
	WdValue                  string `json:"wdValue"`
	WdVisible                bool   `json:"wdVisible"`
}

func TestTmpGwda(t *testing.T) {
	ps := PredicateString{Enabled: true}
	fmt.Println(reflect.TypeOf(ps))
	fmt.Println(reflect.ValueOf(ps))
	v := reflect.ValueOf(ps)
	fmt.Println(v.Kind())

	fmt.Println(v.NumField())
	fmt.Println()

	myType := reflect.TypeOf(ps)
	fmt.Println(myType.NumField())
	for i := 0; i < myType.NumField(); i++ {
		fmt.Println(myType.Field(i).Name, myType.Field(i).Tag.Get("json"))
	}
}
