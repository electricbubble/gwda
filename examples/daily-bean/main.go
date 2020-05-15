package main

import (
	"github.com/electricbubble/gwda"
	"log"
	"time"
)

func main() {
	client, err := gwda.NewClient("http://localhost:8100")
	checkErr(err)

	gwda.DefaultWaitTimeout = time.Second * 20

	session, err := client.NewSession()
	checkErr(err)

	bundleId := "com.360buy.jdmobile"
	// 打开App 京东
	checkErr(session.AppLaunch(bundleId))

	elemMine, err := waitForElement(session, gwda.WDALocator{Name: "我的"})
	checkErr(err, "找到底部导航栏按钮 '我的'")

	value, err := elemMine.Value()
	checkErr(err)
	if value != "1" {
		// 当前页不是 "我的", 点击进入
		checkErr(elemMine.Click(), "点击底部导航栏按钮 '我的'")
	}

	findAndClick(session, gwda.WDALocator{Name: "京豆"}, "按钮 '京豆'")

	_, err = waitForElement(session, gwda.WDALocator{Name: "京豆收支明细"})
	checkErr(err, "当前页 '我的京豆'")

	elemSignIn, err := waitForElement(session, gwda.WDALocator{Predicate: "rect.x == 82 && rect.y == 154 && rect.width == 211 && rect.height == 85"})
	checkErr(err, "找到进入签到页的按钮")

	lblSignIn, err := elemSignIn.Label()
	checkErr(err, "获取标签名 进入签到页的按钮")

	if lblSignIn == "已签到" {
		log.Println("已签到, 返回到 '我的'")
		findAndClick(session, gwda.WDALocator{Name: "返回按钮"}, "返回按钮")
		return
	}
	// 进入签到页
	checkErr(elemSignIn.Click())

	findAndClick(session, gwda.WDALocator{Name: "签到领京豆"}, "按钮 '签到领京豆'")

	// 签到后会跳转到 签到日历
	_, err = waitForElement(session, gwda.WDALocator{Predicate: "name BEGINSWITH '签到成功，'"})
	checkErr(err, "等待签到成功")

	findAndClick(session, gwda.WDALocator{Predicate: "rect.x == 0 && rect.y == 44 && rect.width == 44 && rect.height == 44"}, "签到日历的 返回按钮")

	elemBackLv2, err := waitForElement(session, gwda.WDALocator{Name: "返回"})
	checkErr(err, "找到签到页的 返回按钮")
	rectBackLv2, err := elemBackLv2.Rect()
	checkErr(err, "获取签到页 返回按钮 的坐标")

	// 判断第二层的返回按钮是否可见
	isShowBackLv2 := func(s *gwda.Session) (bool, error) {
		isDisplayed, fErr := elemBackLv2.IsDisplayed()
		// 如果查看 可见性 报错，则直接跳出判断，结束 `session.Wait`
		if fErr != nil {
			return false, fErr
		}
		if isDisplayed {
			return true, nil
		} else {
			// 如果 返回按钮 不可见，则可能是出现了一个提示性弹窗，点击可使其消失
			fErr := session.Tap(rectBackLv2.X+rectBackLv2.Width+1, rectBackLv2.Y+rectBackLv2.Height+1)
			// 如果查看 点击 报错，则直接跳出判断，结束 `session.Wait`
			return false, fErr
		}
	}
	checkErr(session.Wait(isShowBackLv2), "等待签到页的按钮可见 返回按钮")
	checkErr(elemBackLv2.Click(), "点击签到页的 返回按钮")

	findAndClick(session, gwda.WDALocator{Name: "返回按钮"}, "'我的京豆'的 返回按钮")
}

func findAndClick(session *gwda.Session, locator gwda.WDALocator, msg string) {
	elem, err := waitForElement(session, locator)
	checkErr(err, "找到 "+msg)
	checkErr(elem.Click(), "点击 "+msg)
}

func waitForElement(session *gwda.Session, locator gwda.WDALocator) (element *gwda.Element, err error) {
	var fErr error
	exists := func(s *gwda.Session) (bool, error) {
		element, fErr = s.FindElement(locator)
		if fErr == nil {
			return true, nil
		}
		return false, nil
	}

	if err = session.Wait(exists); err != nil {
		// 当前的判断条件，这里只可能是超时
		return nil, err
	}
	return element, fErr
}

func checkErr(err error, msg ...string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
