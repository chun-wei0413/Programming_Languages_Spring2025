package test

import (
	"documentgeneratorframework/documentgenerator"
	"testing"
)

func TestTextDocumentGenerate(t *testing.T) {
	doc := documentgenerator.NewTextDocument()
	result := doc.Generate()
	expected := "Saving text document: Formatted Text: This is the raw text data."
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestTextDocumentPrepareData(t *testing.T) {
	doc := documentgenerator.NewTextDocument()
	result := doc.PrepareData()
	expected := "This is the raw text data."
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestTextDocumentFormatContent(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"test data", "Formatted Text: test data"},
		{"", "Formatted Text: "},
		{"123", "Formatted Text: 123"},
	}
	doc := documentgenerator.NewTextDocument()
	for _, tt := range tests {
		result := doc.FormatContent(tt.input)
		if result != tt.expected {
			t.Errorf("Input %q: Expected %q, got %q", tt.input, tt.expected, result)
		}
	}
}

func TestTextDocumentSave(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"formatted content", "Saving text document: formatted content"},
		{"", "Saving text document: "},
		{"<test>", "Saving text document: <test>"},
	}
	doc := documentgenerator.NewTextDocument()
	for _, tt := range tests {
		result := doc.Save(tt.input)
		if result != tt.expected {
			t.Errorf("Input %q: Expected %q, got %q", tt.input, tt.expected, result)
		}
	}
}
