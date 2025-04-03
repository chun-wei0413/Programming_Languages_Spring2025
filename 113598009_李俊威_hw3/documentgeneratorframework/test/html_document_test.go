package test

import (
	"documentgeneratorframework/documentgenerator"
	"testing"
)

func TestHTMLDocumentGenerate(t *testing.T) {
	doc := documentgenerator.NewHTMLDocument()
	result := doc.Generate()
	expected := "Saving HTML document: <div><html><body>This is raw HTML data.</body></html></div>"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestHTMLDocumentPrepareData(t *testing.T) {
	doc := documentgenerator.NewHTMLDocument()
	result := doc.PrepareData()
	expected := "<html><body>This is raw HTML data.</body></html>"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestHTMLDocumentFormatContent(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"test data", "<div>test data</div>"},
		{"", "<div></div>"},
		{"<p>html</p>", "<div><p>html</p></div>"},
	}
	doc := documentgenerator.NewHTMLDocument()
	for _, tt := range tests {
		result := doc.FormatContent(tt.input)
		if result != tt.expected {
			t.Errorf("Input %q: Expected %q, got %q", tt.input, tt.expected, result)
		}
	}
}

func TestHTMLDocumentSave(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"formatted content", "Saving HTML document: formatted content"},
		{"", "Saving HTML document: "},
		{"<div>test</div>", "Saving HTML document: <div>test</div>"},
	}
	doc := documentgenerator.NewHTMLDocument()
	for _, tt := range tests {
		result := doc.Save(tt.input)
		if result != tt.expected {
			t.Errorf("Input %q: Expected %q, got %q", tt.input, tt.expected, result)
		}
	}
}
