package test

import (
	"documentgeneratorframework/documentgenerator"
	"testing"
)

func TestBaseGeneratorWithTextDocument(t *testing.T) {
	textDoc := documentgenerator.NewTextDocument()
	base := &documentgenerator.BaseGenerator{Doc: textDoc}
	result := base.Generate()
	expected := "Saving text document: Formatted Text: This is the raw text data."
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestBaseGeneratorWithHTMLDocument(t *testing.T) {
	htmlDoc := documentgenerator.NewHTMLDocument()
	base := &documentgenerator.BaseGenerator{Doc: htmlDoc}
	result := base.Generate()
	expected := "Saving HTML document: <div><html><body>This is raw HTML data.</body></html></div>"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestBaseGeneratorNilDoc(t *testing.T) {
	base := &documentgenerator.BaseGenerator{Doc: nil}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic when doc is nil, but no panic occurred")
		}
	}()
	base.Generate() // 應該觸發 panic
}
