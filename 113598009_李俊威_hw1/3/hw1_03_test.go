package main

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

// go test -v
// TestReadFile tests the file reading functionality
func TestReadFile(t *testing.T) {
	// Create a temporary test file
	content := "Hello, world!"
	tmpFile, err := os.CreateTemp("", "test_file*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	
	if _, err := tmpFile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()
	
	// Test reading the file
	result, err := readFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Error reading file: %v", err)
	}
	
	if result != content {
		t.Errorf("Expected content: %s, got: %s", content, result)
	}
}

// TestFilterCharsAndNormalize tests character filtering and normalization
func TestFilterCharsAndNormalize(t *testing.T) {
	input := "Hello, World! This is a test@123."
	result := filterCharsAndNormalize(input)
	
	// Check if output is lowercase and special characters are replaced with spaces
	if !strings.Contains(result, "hello") || !strings.Contains(result, "world") {
		t.Errorf("Expected filtered text to contain 'hello' and 'world', got: %s", result)
	}
	
	if strings.Contains(result, "!") || strings.Contains(result, "@") || strings.Contains(result, ",") {
		t.Errorf("Expected special characters to be removed, got: %s", result)
	}
}

// TestScan tests word scanning functionality
func TestScan(t *testing.T) {
	input := "hello world test"
	result := scan(input)
	expected := []string{"hello", "world", "test"}
	
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got: %v", expected, result)
	}
}

// TestRemoveStopWords tests stop word removal
func TestRemoveStopWords(t *testing.T) {
	// Create a temporary stop words file
	stopWordsContent := "a,is,the,this"
	stopFile, err := os.CreateTemp("", "stop_words*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(stopFile.Name())
	
	if _, err := stopFile.Write([]byte(stopWordsContent)); err != nil {
		t.Fatal(err)
	}
	stopFile.Close()
	
	// Test removing stop words
	input := []string{"hello", "the", "world", "a", "test", "this", "is"}
	expected := []string{"hello", "world", "test"}
	
	result, err := removeStopWords(input)
	if err != nil {
		t.Fatalf("Error removing stop words: %v", err)
	}
	
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got: %v", expected, result)
	}
	
	// Test that single letters are removed
	input = []string{"hello", "a", "b", "c", "test"}
	expected = []string{"hello", "test"}
	
	result, err = removeStopWords(input)
	if err != nil {
		t.Fatalf("Error removing stop words: %v", err)
	}
	
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got: %v", expected, result)
	}
}

// TestFrequencies tests word counting
func TestFrequencies(t *testing.T) {
	input := []string{"hello", "world", "hello", "test"}
	result := frequencies(input)
	
	if result["hello"] != 2 {
		t.Errorf("Expected frequency of 'hello' to be 2, got: %d", result["hello"])
	}
	
	if result["world"] != 1 {
		t.Errorf("Expected frequency of 'world' to be 1, got: %d", result["world"])
	}
	
	if result["test"] != 1 {
		t.Errorf("Expected frequency of 'test' to be 1, got: %d", result["test"])
	}
}

// TestSortFrequencies tests sorting by frequency
func TestSortFrequencies(t *testing.T) {
	input := map[string]int{
		"hello": 3,
		"world": 2,
		"test":  5,
	}
	
	result := sortFrequencies(input)
	
	// First word should be "test" with count 5
	if result[0].Word != "test" || result[0].Count != 5 {
		t.Errorf("Expected highest frequency word to be 'test' with count 5, got: %s with count %d", 
			result[0].Word, result[0].Count)
	}
	
	// Last word should be "world" with count 2
	if result[2].Word != "world" || result[2].Count != 2 {
		t.Errorf("Expected lowest frequency word to be 'world' with count 2, got: %s with count %d", 
			result[2].Word, result[2].Count)
	}
}