package main

import (
	"github.com/electricbubble/gwda"
	"log"
	"time"
)

func main() {
	client, err := gwda.NewClient("http://localhost:8100")
	checkErr("连接设备", err)
	session, err := client.NewSession()
	checkErr("创建 session", err)
	windowSize, err := session.WindowSize()
	checkErr("获取当前应用的大小", err)

	setup(session)

	center := gwda.WDACoordinate{X: windowSize.Width / 2, Y: windowSize.Height / 2}

	err = session.Drag(
		center.X-100, center.Y-100,
		center.X+100, center.Y-100)
	// 可指定按住的时间，默认 1s
	// err = session.Drag(center.X-100, center.Y-100, center.X+100, center.Y-100, 3)
	checkErr("拖拽手势", err)

	setup(session, 3)

	err = session.SwipeLeft()
	checkErr("向左👈滑动", err)

	setup(session, 3)

	// 默认按压 1s
	err = session.ForceTouchCoordinate(center, 3.5)
	checkErr("指定压力值触发 3D Touch", err)

	setup(session, 3)

	topLeft := gwda.WDACoordinate{X: center.X - 100, Y: center.Y - 100}
	topRight := gwda.WDACoordinate{X: center.X + 100, Y: topLeft.Y}
	lowerLeft := gwda.WDACoordinate{X: topLeft.X, Y: center.Y + 100}
	lowerRight := gwda.WDACoordinate{X: topRight.X, Y: center.Y + 100}

	_ = lowerLeft

	element, err := session.FindElement(gwda.WDALocator{Name: "自定手势作用区域"})
	checkErr("自定手势作用区域", err)

	touchActions := gwda.NewWDATouchActions().
		// 同时设置元素和坐标，坐标是元素的坐标位置
		Press(gwda.NewWDATouchActionOptionPress().SetElement(element).SetXYCoordinate(topLeft).SetPressure(0.8)).
		// LongPress(gwda.NewWDATouchActionOptionLongPress().SetElement(element).SetXY(topLeft.X, topLeft.Y)).
		Wait(0.2).
		// 只设置了坐标，则是当前屏幕的坐标位置
		MoveTo(gwda.NewWDATouchActionOptionMoveTo().SetXYCoordinate(topRight)).
		Wait(0.2).
		// 如果只设置了元素，则默认坐标为 元素的中心
		MoveTo(gwda.NewWDATouchActionOptionMoveTo().SetElement(element)).
		Wait(0.2).
		MoveTo(gwda.NewWDATouchActionOptionMoveTo().SetElement(element).SetXYCoordinate(lowerRight)).
		Release()
	err = session.PerformTouchActions(touchActions)
	checkErr("z 手势", err)

	setup(session, 3)

	actions := gwda.NewWDAActions(2).
		Swipe(center.X-100, center.Y-100, center.X, center.Y).
		Swipe(center.X+100, center.Y+100, center.X, center.Y)
	// 如果设置了元素，则坐标是从元素中心点开始的相对坐标
	// actions = gwda.NewWDAActions().
	// 	Swipe(0-100, 0-100, 0, 0, element).
	// 	Swipe(100, 100, 0, 0, element)
	err = session.PerformActions(actions)
	checkErr("缩放 手势", err)

	setup(session, 3)

	actions = gwda.NewWDAActions().
		DoubleTap(center.X+60, center.Y).
		Swipe(center.X, center.Y-100, center.X, center.Y+100)
	err = session.PerformActions(actions)
	checkErr("组合手势，下滑并双击", err)
}

func setup(session *gwda.Session, duration ...time.Duration) {
	if len(duration) != 0 {
		time.Sleep(time.Second * duration[0])
	}
	bundleId := "com.apple.Preferences"

	appRunState, err := session.AppState(bundleId)
	checkErr("获取指定 App 的运行状态", err)
	switch appRunState {
	case gwda.WDAAppNotRunning:
		log.Println("该 App 未运行, 开始打开 App:", bundleId)
		err = session.AppLaunch(bundleId)
		checkErr("启动指定 App", err)
	case gwda.WDAAppRunningFront:
		if activeNavBarName(session) == "新建手势" {
			findAndClick(session, gwda.WDALocator{LinkText: gwda.NewWDAElementAttribute().SetLabel("取消")}, "新建手势 取消")
			time.Sleep(time.Second * 1)
			findAndClick(session, gwda.WDALocator{PartialLinkText: gwda.NewWDAElementAttribute().SetLabel("创建新手势")}, "创建新手势…")
			time.Sleep(time.Second * 1)
			return
		} else {
			restartApp(session, bundleId)
		}
	default:
		restartApp(session, bundleId)
	}

	err = session.SwipeDown()
	checkErr("向下👇滑动", err)

	elemSearch, err := session.FindElement(gwda.WDALocator{ClassName: gwda.WDAElementType{SearchField: true}})
	checkErr("找到 搜索输入框", err)

	// targetName := "切换控制"

	err = elemSearch.SendKeys("切换控制" + "\n")
	checkErr("输入文本", err)

	elemSearchRet, err := session.FindElement(gwda.WDALocator{Predicate: "type in {'XCUIElementTypeTable', 'XCUIElementTypeCollectionView'} && visible == true"})
	checkErr("找到 搜索结果列表框", err)

	findAndClick(elemSearchRet, gwda.WDALocator{ClassName: gwda.WDAElementType{Cell: true}}, "第一个搜索结果")

	// 获取当前导航栏的 name 属性值
	navBarName := activeNavBarName(session)

	if navBarName != "切换控制" {
		findAndClick(session, gwda.WDALocator{LinkText: gwda.NewWDAElementAttribute().SetLabel("切换控制")}, "切换控制")
	}

	isSwitched := func(s *gwda.Session) (bool, error) {
		if activeNavBarName(s) == "切换控制" {
			return true, nil
		}
		return false, nil
	}
	checkErr("等待列表切换", session.WaitWithTimeoutAndInterval(isSwitched, 10, 0.1))

	elemList, err := session.FindElement(gwda.WDALocator{ClassName: gwda.WDAElementType{Table: true}})
	checkErr("找到当前列表 切换控制", err)

	// targetItem := "已存储的手势"

	err = elemList.ScrollElementByPredicate("type == 'XCUIElementTypeCell' && name == '已存储的手势'")
	checkErr("滚动找到 已存储的手势", err)

	findAndClick(session, gwda.WDALocator{Name: "已存储的手势"}, "已存储的手势")

	findAndClick(session, gwda.WDALocator{PartialLinkText: gwda.NewWDAElementAttribute().SetLabel("创建新手势")}, "创建新手势…")
}

func activeNavBarName(session *gwda.Session) string {
	navBar, err := session.FindElement(gwda.WDALocator{ClassName: gwda.WDAElementType{NavigationBar: true}})
	checkErr("找到当前页导航栏", err)

	attrName, err := navBar.GetAttribute(gwda.NewWDAElementAttribute().SetName(""))
	checkErr("获取导航栏 name 属性值", err)
	// log.Println("当前导航栏 name 属性值:", attrName)
	return attrName
}

func findAndClick(scope interface{}, locator gwda.WDALocator, msg string) {
	var elem *gwda.Element
	var err error
	switch scope := scope.(type) {
	case *gwda.Session:
		elem, err = scope.FindElement(locator)
	case *gwda.Element:
		elem, err = scope.FindElement(locator)
	}
	checkErr("找到 "+msg, err)
	err = elem.Click()
	checkErr("点击 "+msg, err)
}

func restartApp(session *gwda.Session, bundleId string) {
	log.Println("重新启动 App:", bundleId)
	err := session.AppTerminate(bundleId)
	checkErr("关闭指定 App", err)
	err = session.AppLaunch(bundleId)
	checkErr("再启动指定 App", err)
}

func checkErr(msg string, err error) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
