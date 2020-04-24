package gwda

import (
	"testing"
)

func TestElement_Click(t *testing.T) {
	c, err := NewClient(deviceURL)
	if err != nil {
		t.Fatal(err)
	}
	bundleId := "com.apple.Preferences"
	_ = bundleId
	s, err := c.NewSession()
	if err != nil {
		t.Fatal(err)
	}
	Debug = true
	element, err := s.FindElement("partial link text", "label=通知")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(element)

	err = element.Click()
	if err != nil {
		t.Fatal(err)
	}
}

func TestElement_Rect(t *testing.T) {
	c, err := NewClient(deviceURL)
	if err != nil {
		t.Fatal(err)
	}
	s, err := c.NewSession()
	if err != nil {
		t.Fatal(err)
	}
	Debug = true
	element, err := s.FindElement("partial link text", "label=通知")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(element)

	rect, err := element.Rect()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(rect)
}

func TestElement_Enabled(t *testing.T) {
	c, err := NewClient(deviceURL)
	if err != nil {
		t.Fatal(err)
	}
	bundleId := "com.apple.Preferences"
	_ = bundleId
	s, err := c.NewSession()
	if err != nil {
		t.Fatal(err)
	}
	Debug = true
	element, err := s.FindElement("partial link text", "label=通知")
	if err != nil {
		t.Fatal(err)
	}

	enabled, err := element.Enabled()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(enabled)
}
