package main

import (
	"github.com/electricbubble/gwda"
	"log"
)

func main() {
	client, err := gwda.NewClient("http://localhost:8100")
	checkErr("连接设备", err)

	session, err := client.NewSession()
	checkErr("创建 session", err)

	source, err := client.AccessibleSource()
	// source, err = session.AccessibleSource()
	checkErr("查看当前可见的全部元素", err)
	log.Println(source)

	// 默认返回 xml 格式
	tree, err := session.Source()
	// tree, err = client.Source()
	checkErr("查看当前全部元素(包括不可见)", err)
	log.Println(tree)

	// 仅 xml 支持
	sTree, err := session.Source(gwda.NewWDASourceOption().SetExcludedAttributes([]string{"enabled", "visible", "type"}))
	checkErr("排除指定属性", err)
	log.Println(sTree)
	// 默认返回的 xml
	// <XCUIElementTypeOther type="XCUIElementTypeOther" name="程序坞" label="程序坞" enabled="true" visible="true" x="0" y="575" width="375" height="92">
	// 排除指定属性后
	// <XCUIElementTypeOther name="程序坞" label="程序坞" x="0" y="575" width="375" height="92">

	sJson, err := client.Source(gwda.NewWDASourceOption().SetFormatAsJson())
	checkErr("指定返回格式为 json", err)
	log.Println(sJson)

	sDesc, err := session.Source(gwda.NewWDASourceOption().SetFormatAsDescription())
	checkErr("指定返回格式为 description", err)
	log.Println(sDesc)

}

func checkErr(msg string, err error) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
