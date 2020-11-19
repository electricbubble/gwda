package main

import (
	"github.com/electricbubble/gwda"
	"log"
)

func main() {
	driver, err := gwda.NewUSBDriver(nil)
	if err != nil {
		log.Fatalln(err)
	}

	x, y := 50, 256

	err = driver.Tap(x, y)
	if err != nil {
		log.Fatalln(err)
	}

	err = driver.DoubleTap(x, y)
	if err != nil {
		log.Fatalln(err)
	}

	err = driver.TouchAndHold(x, y)
	if err != nil {
		log.Fatalln(err)
	}

	fromX, fromY, toX, toY := 50, 256, 100, 256

	err = driver.Drag(fromX, fromY, toX, toY)
	if err != nil {
		log.Fatalln(err)
	}

	err = driver.Swipe(fromX, fromY, toX, toY)
	if err != nil {
		log.Fatalln(err)
	}

	// 需要 3D Touch 硬件支持
	// err = driver.ForceTouch(x, y, 0.8)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// 多点触控
	// actions := gwda.NewW3CActions().FingerAction(
	// 	gwda.NewFingerAction().
	// 		Move(gwda.NewFingerMove().WithXY(50, 128)).
	// 		Down().
	// 		Pause(0.25).
	// 		Move(gwda.NewFingerMove().WithXY(200, 160)).
	// 		Pause(0.25).
	// 		Up(),
	// 	gwda.NewFingerAction().
	// 		Move(gwda.NewFingerMove().WithXY(300, 256)).
	// 		Down().
	// 		Pause(0.25).
	// 		Move(gwda.NewFingerMove().WithXY(200, 160)).
	// 		Pause(0.25).
	// 		Up(),
	// )
	// err = driver.PerformW3CActions(actions)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// 一笔画手势
	// touchActs := gwda.NewTouchActions().
	// 	Press(gwda.NewTouchActionPress().WithXY(50, 150)).
	// 	Wait(0.2).
	// 	MoveTo(gwda.NewTouchActionMoveTo().WithXY(300, 150)).
	// 	Wait(0.2).
	// 	MoveTo(gwda.NewTouchActionMoveTo().WithXY(50, 256)).
	// 	Wait(0.2).
	// 	MoveTo(gwda.NewTouchActionMoveTo().WithXY(300, 256)).
	// 	Release()
	//
	// err = driver.PerformAppiumTouchActions(touchActs)
	// if err != nil {
	// 	log.Fatal(err)
	// }
}
