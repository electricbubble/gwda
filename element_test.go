package gwda

import (
	"math"
	"testing"
)

func setupElement(t *testing.T, by BySelector) WebElement {
	setup(t)
	element, err := driver.FindElement(by)
	if err != nil {
		t.Fatal(err)
	}
	return element
}

func Test_remoteWE_Click(t *testing.T) {
	element := setupElement(t, BySelector{LinkText: NewElementAttribute().WithLabel("设置")})

	err := element.Click()
	if err != nil {
		t.Fatal(err)
	}
}

func Test_remoteWE_SendKeys(t *testing.T) {
	element := setupElement(t, BySelector{ClassName: ElementType{SearchField: true}})

	err := element.SendKeys("App Store")
	// err := element.SendKeys("App Store", 3)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_remoteWE_Clear(t *testing.T) {
	element := setupElement(t, BySelector{ClassName: ElementType{SearchField: true}})

	err := element.Clear()
	if err != nil {
		t.Fatal(err)
	}
}

func Test_remoteWE_Tap(t *testing.T) {
	element := setupElement(t, BySelector{Name: "touchableView"})

	err := element.Tap(10, 20)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_remoteWE_DoubleTap(t *testing.T) {
	element := setupElement(t, BySelector{Name: "touchableView"})

	err := element.DoubleTap()
	if err != nil {
		t.Fatal(err)
	}
}

func Test_remoteWE_TouchAndHold(t *testing.T) {
	element := setupElement(t, BySelector{Name: "touchableView"})

	err := element.TouchAndHold(-1)
	// err := element.TouchAndHold(5)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_remoteWE_TwoFingerTap(t *testing.T) {
	element := setupElement(t, BySelector{Name: "touchableView"})

	err := element.TwoFingerTap()
	if err != nil {
		t.Fatal(err)
	}
}

func Test_remoteWE_TapWithNumberOfTaps(t *testing.T) {
	element := setupElement(t, BySelector{Name: "touchableView"})

	err := element.TapWithNumberOfTaps(3, 3)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_remoteWE_ForceTouch(t *testing.T) {
	element := setupElement(t, BySelector{Name: "touchableView"})

	// err := element.ForceTouch(1, -1)
	err := element.ForceTouchFloat(10, 20, 1, -1)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_remoteWE_Drag(t *testing.T) {
	element := setupElement(t, BySelector{Name: "touchableView"})

	// err := element.Drag(10, 20, 10, 300, -1)
	err := element.Swipe(10, 20, 10, 300)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_remoteWE_SwipeDirection(t *testing.T) {
	element := setupElement(t, BySelector{Name: "touchableView"})

	// err := element.SwipeDirection(DirectionUp, -1)
	err := element.SwipeDirection(DirectionDown, 120)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_remoteWE_Pinch(t *testing.T) {
	element := setupElement(t, BySelector{Name: "touchableView"})

	// zoom in
	// err := element.Pinch(2,10)
	// zoom out
	err := element.Pinch(0.9, -4.5)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_remoteWE_PinchToZoomOutByW3CAction(t *testing.T) {
	element := setupElement(t, BySelector{Name: "touchableView"})

	err := element.PinchToZoomOutByW3CAction(15)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_remoteWE_Rotate(t *testing.T) {
	element := setupElement(t, BySelector{Name: "touchableView"})

	// 90 CW
	// err := element.Rotate(math.Pi / 2)
	// 180 CCW
	err := element.Rotate(math.Pi * -2)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_remoteWE_PickerWheelSelect(t *testing.T) {
	element := setupElement(t, BySelector{ClassName: ElementType{PickerWheel: true}})

	err := element.PickerWheelSelect(PickerWheelOrderNext, 3)
	if err != nil {
		t.Fatal(err)
	}
	err = element.PickerWheelSelect(PickerWheelOrderPrevious)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_remoteWE_scroll(t *testing.T) {
	element := setupElement(t, BySelector{ClassName: ElementType{Table: true}})

	var err error
	// err = element.ScrollElementByName("电池")
	// err = element.ScrollElementByPredicate("type == 'XCUIElementTypeCell' AND name LIKE 'Safari*'")
	err = element.ScrollDirection(DirectionDown, 0.8)

	// element, err = driver.FindElement(BySelector{PartialLinkText: NewElementAttribute().WithLabel("Safari")})
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// err = element.ScrollToVisible()

	if err != nil {
		t.Fatal(err)
	}
}

func Test_remoteWE_FindElement(t *testing.T) {
	element := setupElement(t, BySelector{ClassName: ElementType{Table: true}})

	SetDebug(true)

	var err error
	element, err = element.FindElement(BySelector{PartialLinkText: NewElementAttribute().WithLabel("Safari")})
	if err != nil {
		t.Fatal(err)
	}

	SetDebug(false)
	err = element.Click()
	if err != nil {
		t.Fatal(err)
	}
}

func Test_remoteWE_FindElements(t *testing.T) {
	element := setupElement(t, BySelector{ClassName: ElementType{Table: true}})

	SetDebug(true)

	elements, err := element.FindElements(BySelector{ClassName: ElementType{Cell: true}})
	if err != nil {
		t.Fatal(err)
	}

	SetDebug(false)
	err = elements[0].Click()
	if err != nil {
		t.Fatal(err)
	}
}

func Test_remoteWE_FindVisibleCells(t *testing.T) {
	element := setupElement(t, BySelector{ClassName: ElementType{Table: true}})

	SetDebug(true)

	cells, err := element.FindVisibleCells()
	if err != nil {
		t.Fatal(err)
	}

	SetDebug(false)
	err = cells[0].Click()
	if err != nil {
		t.Fatal(err)
	}
}

func Test_remoteWE_Rect(t *testing.T) {
	element := setupElement(t, BySelector{ClassName: ElementType{Switch: true}})

	rect, err := element.Rect()
	if err != nil {
		t.Fatal(err)
	}
	location, err := element.Location()
	if err != nil {
		t.Fatal(err)
	}
	size, err := element.Size()
	if err != nil {
		t.Fatal(err)
	}
	_, _, _ = rect, location, size
	t.Log(rect, location, size)
}

func Test_remoteWE_Text(t *testing.T) {
	element := setupElement(t, BySelector{ClassName: ElementType{Switch: true}})

	text, err := element.Text()
	if err != nil {
		t.Fatal(err)
	}
	_ = text
	// t.Log(text)
}

func Test_remoteWE_Type(t *testing.T) {
	element := setupElement(t, BySelector{ClassName: ElementType{Switch: true}})

	elemType, err := element.Type()
	if err != nil {
		t.Fatal(err)
	}
	_ = elemType
	// t.Log(elemType)
}

func Test_remoteWE_IsEnabled(t *testing.T) {
	element := setupElement(t, BySelector{ClassName: ElementType{Switch: true}})

	enabled, err := element.IsEnabled()
	if err != nil {
		t.Fatal(err)
	}
	_ = enabled
	// t.Log(enabled)
}

func Test_remoteWE_IsDisplayed(t *testing.T) {
	element := setupElement(t, BySelector{PartialLinkText: NewElementAttribute().WithLabel("Safari")})

	displayed, err := element.IsDisplayed()
	if err != nil {
		t.Fatal(err)
	}
	_ = displayed
	// t.Log(displayed)
}

func Test_remoteWE_IsSelected(t *testing.T) {
	element := setupElement(t, BySelector{ClassName: ElementType{Switch: true}})
	// element := setupElement(t, BySelector{Name: "添加到主屏幕"})
	// element := setupElement(t, BySelector{Name: "仅App资源库"})

	selected, err := element.IsSelected()
	if err != nil {
		t.Fatal(err)
	}
	_ = selected
	// t.Log(selected)
}

func Test_remoteWE_IsAccessible(t *testing.T) {
	element := setupElement(t, BySelector{ClassName: ElementType{Switch: true}})

	accessible, err := element.IsAccessible()
	if err != nil {
		t.Fatal(err)
	}
	_ = accessible
	// t.Log(accessible)
}

func Test_remoteWE_IsAccessibilityContainer(t *testing.T) {
	// element := setupElement(t, BySelector{ClassName: ElementType{Switch: true}})
	element := setupElement(t, BySelector{ClassName: ElementType{Table: true}})

	isAccessibilityContainer, err := element.IsAccessibilityContainer()
	if err != nil {
		t.Fatal(err)
	}
	_ = isAccessibilityContainer
	// t.Log(isAccessibilityContainer)
}

func Test_remoteWE_GetAttribute(t *testing.T) {
	element := setupElement(t, BySelector{ClassName: ElementType{StaticText: true}})

	value, err := element.GetAttribute(NewElementAttribute().WithValue(""))
	if err != nil {
		t.Fatal(err)
	}
	_ = value
	// t.Log(value)
}

func Test_remoteWE_Screenshot(t *testing.T) {
	element := setupElement(t, BySelector{ClassName: ElementType{TextView: true}})

	screenshot, err := element.Screenshot()
	if err != nil {
		t.Fatal(err)
	}
	_ = screenshot

	// img, format, err := image.Decode(screenshot)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// userHomeDir, _ := os.UserHomeDir()
	// file, err := os.Create(userHomeDir + "/Desktop/e1." + format)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// defer func() { _ = file.Close() }()
	// switch format {
	// case "png":
	// 	err = png.Encode(file, img)
	// case "jpeg":
	// 	err = jpeg.Encode(file, img, nil)
	// }
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// t.Log(file.Name())
}
