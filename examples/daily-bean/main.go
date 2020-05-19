package main

import (
	"fmt"
	. "github.com/electricbubble/gwda"
	"log"
	"time"
)

func main() {
	client, err := NewClient("http://localhost:8100")
	checkErr(err)

	DefaultWaitTimeout = time.Second * 20

	session, err := client.NewSession()
	checkErr(err)

	bundleId := "com.360buy.jdmobile"
	// 打开App 京东
	checkErr(session.AppLaunch(bundleId))

	elemMine, err := waitForElement(session, WDALocator{Name: "我的"})
	checkErr(err, "找到底部导航栏按钮 '我的'")

	value, err := elemMine.Value()
	checkErr(err)
	if value != "1" {
		// 当前页不是 "我的", 点击进入
		checkErr(elemMine.Click(), "点击底部导航栏按钮 '我的'")
	}

	findAndClick(session, WDALocator{Name: "京豆"}, "按钮 '京豆'")

	_, err = waitForElement(session, WDALocator{Name: "京豆收支明细"})
	checkErr(err, "当前页 '我的京豆'")

	elemSignIn, err := waitForElement(session, WDALocator{Predicate: "rect.x == 82 && rect.y == 154 && rect.width == 211 && rect.height == 85"})
	checkErr(err, "找到进入签到页的按钮")

	lblSignIn, err := elemSignIn.Label()
	checkErr(err, "获取标签名 进入签到页的按钮")

	if lblSignIn == "已签到" {
		log.Println("已签到, 返回到 '我的'")
		findAndClick(session, WDALocator{Name: "返回按钮"}, "返回按钮")
		return
	}
	// 进入签到页
	checkErr(elemSignIn.Click())

	findAndClick(session, WDALocator{Name: "签到领京豆"}, "按钮 '签到领京豆'")

	// 签到后会跳转到 签到日历
	_, err = waitForElement(session, WDALocator{Predicate: "name BEGINSWITH '签到成功，'"})
	checkErr(err, "等待签到成功")

	findAndClick(session, WDALocator{Predicate: "rect.x == 0 && rect.y == 44 && rect.width == 44 && rect.height == 44"}, "签到日历的 返回按钮")

	locatorBackLv2 := WDALocator{Name: "返回"}
	findAndClick(session, locatorBackLv2, "签到页的 返回按钮")
	// 点击后再查找这个 返回按钮
	elemBackLv2, err := session.FindElement(locatorBackLv2)
	if err == nil {
		// 如果找到了，意味着出现了每日仅提示两次的一个弹窗
		rectBackLv2, err := elemBackLv2.Rect()
		checkErr(err, "获取签到页 返回按钮 的坐标")
		// 通过点击使弹窗消失
		checkErr(session.Tap(rectBackLv2.X+rectBackLv2.Width+1, rectBackLv2.Y+rectBackLv2.Height+1))

		isVisible := func(s *Session) (bool, error) {
			isDisplayed, err2 := elemBackLv2.IsDisplayed()
			if err2 != nil {
				return false, err2
			}
			if isDisplayed {
				return true, nil
			}
			return false, nil
		}
		_ = session.WaitWithTimeoutAndInterval(isVisible, 3, 0.1)
		checkErr(elemBackLv2.Click(), "点击签到页的 返回按钮")
	}

	findAndClick(session, WDALocator{Name: "返回按钮"}, "'我的京豆'的 返回按钮")
}

func findAndClick(session *Session, locator WDALocator, msg string) {
	elem, err := waitForElement(session, locator)
	checkErr(err, "找到 "+msg)
	checkErr(elem.Click(), "点击 "+msg)
}

func waitForElement(session *Session, locator WDALocator) (element *Element, err error) {
	var fErr error
	exists := func(s *Session) (bool, error) {
		element, fErr = s.FindElement(locator)
		if fErr == nil {
			return true, nil
		}
		// 如果直接返回 err 将直接终止 `session.Wait`
		return false, nil
	}

	if err = session.Wait(exists); err != nil {
		// 当前的判断条件，这里只可能是超时
		return nil, fmt.Errorf("%s: %w", err.Error(), fErr)
	}
	return element, fErr
}

func checkErr(err error, msg ...string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
