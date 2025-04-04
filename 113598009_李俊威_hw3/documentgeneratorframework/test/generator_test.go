package test

import (
	"documentgeneratorframework/documentgenerator"
	"testing"
)

func TestBaseGeneratorWithTextDocument(t *testing.T) {
	textDoc := documentgenerator.NewTextDocument()
	// 這裡inject textDoc到BaseGenerator是為了測試base的Generate()
	base := &documentgenerator.BaseGenerator{textDoc}

	result := base.Generate()
	expected := "Saving text document: Formatted Text: This is the raw text data."
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestBaseGeneratorWithHTMLDocument(t *testing.T) {
	htmlDoc := documentgenerator.NewHTMLDocument()
	// 這裡同理
	base := &documentgenerator.BaseGenerator{htmlDoc}

	result := base.Generate()
	expected := "Saving HTML document: <div><html><body>This is raw HTML data.</body></html></div>"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestBaseGeneratorNilDoc(t *testing.T) {
	base := &documentgenerator.BaseGenerator{nil}
	//原本call nil的Generate會觸發panic但我在BaseGenerator的Generate裡面加了nil判斷
	//所以會return錯誤訊息
	result := base.Generate()
	expected := "Error: no document generator provided"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}
