package main

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

func setupTestFile(t *testing.T, content string) string {
	f, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer f.Close()

	_, err = f.WriteString(content)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	return f.Name()
}

func cleanupTestFile(t *testing.T, path string) {
	err := os.Remove(path)
	if err != nil {
		t.Fatalf("Failed to remove temp file: %v", err)
	}
}

func TestDataStorageManager_Words(t *testing.T) {
	// 測試用檔案內容
	content := "The cat, and the dog! 123 Cat HERE"
	filePath := setupTestFile(t, content)
	defer cleanupTestFile(t, filePath)

	dsm := NewDataStorageManager(filePath)
	words := dsm.Words()

	expected := []string{"the", "cat", "and", "the", "dog", "cat", "here"}
	if !reflect.DeepEqual(words, expected) {
		t.Errorf("Expected words %v, but got %v", expected, words)
	}
}

// 測試 StopWordManager
func TestStopWordManager_IsStopWord(t *testing.T) {

	stopWordsContent := "the, and, is"
	stopFile := setupTestFile(t, stopWordsContent)
	defer cleanupTestFile(t, stopFile)

	swm := &StopWordManager{stopWords: make(map[string]struct{})}
	content, _ := os.ReadFile(stopFile)
	words := strings.Split(string(content), ",")
	for _, word := range words {
		trimmedWord := strings.TrimSpace(word)
		if len(trimmedWord) > 0 {
			swm.stopWords[trimmedWord] = struct{}{}
		}
	}

	for c := 'a'; c <= 'z'; c++ {
		swm.stopWords[string(c)] = struct{}{}
	}

	tests := []struct {
		word     string
		expected bool
	}{
		{"the", true},    
		{"and", true},    
		{"is", true},     
		{"a", true},      
		{"cat", false},   
		{"dog", false},   
	}

	for _, test := range tests {
		result := swm.IsStopWord(test.word)
		if result != test.expected {
			t.Errorf("For word %q, expected %v, but got %v", test.word, test.expected, result)
		}
	}
}

// 測試 WordFrequencyManager
func TestWordFrequencyManager_IncrementCountAndSorted(t *testing.T) {
	wfm := NewWordFrequencyManager()

	// 測試 IncrementCount
	wfm.IncrementCount("cat")
	wfm.IncrementCount("dog")
	wfm.IncrementCount("cat")

	if wfm.wordFreqs["cat"] != 2 {
		t.Errorf("Expected frequency of 'cat' to be 2, but got %d", wfm.wordFreqs["cat"])
	}
	if wfm.wordFreqs["dog"] != 1 {
		t.Errorf("Expected frequency of 'dog' to be 1, but got %d", wfm.wordFreqs["dog"])
	}

	sorted := wfm.Sorted()
	expected := [][2]string{
		{"cat", "2"},
		{"dog", "1"},
	}
	if !reflect.DeepEqual(sorted, expected) {
		t.Errorf("Expected sorted result %v, but got %v", expected, sorted)
	}
}

// 測試 WordFrequencyController
func TestWordFrequencyController_Run(t *testing.T) {
	content := "The cat and the dog. Cat is here!"
	filePath := setupTestFile(t, content)
	defer cleanupTestFile(t, filePath)

	stopWordsContent := "the, and, is"
	stopFile := setupTestFile(t, stopWordsContent)
	defer cleanupTestFile(t, stopFile)

	swm := &StopWordManager{stopWords: make(map[string]struct{})}
	stopContent, _ := os.ReadFile(stopFile)
	words := strings.Split(string(stopContent), ",")
	for _, word := range words {
		trimmedWord := strings.TrimSpace(word)
		if len(trimmedWord) > 0 {
			swm.stopWords[trimmedWord] = struct{}{}
		}
	}

	for c := 'a'; c <= 'z'; c++ {
		swm.stopWords[string(c)] = struct{}{}
	}

	originalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	controller := &WordFrequencyController{
		storageManager:  NewDataStorageManager(filePath),
		stopWordManager: swm,
		wordFreqManager: NewWordFrequencyManager(),
	}
	controller.Run()

	w.Close()
	os.Stdout = originalStdout

	var output strings.Builder
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf)
		if n > 0 {
			output.Write(buf[:n])
		}
		if err != nil {
			break
		}
	}
	r.Close()

	expectedOutput := "cat - 2\ndog - 1\nhere - 1\n"
	if output.String() != expectedOutput {
		t.Errorf("Expected output %q, but got %q", expectedOutput, output.String())
	}
}

func TestMainFlow(t *testing.T) {
	originalArgs := os.Args
	os.Args = []string{"program", "testfile.txt"}
	defer func() { os.Args = originalArgs }()

	content := "The cat and the dog."
	filePath := setupTestFile(t, content)
	defer cleanupTestFile(t, filePath)
	os.Args[1] = filePath

	stopWordsContent := "the, and"
	stopFile := setupTestFile(t, stopWordsContent)
	defer cleanupTestFile(t, stopFile)

	originalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	main()

	w.Close()
	os.Stdout = originalStdout

	var output strings.Builder
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf)
		if n > 0 {
			output.Write(buf[:n])
		}
		if err != nil {
			break
		}
	}
	r.Close()

	expectedOutput := "cat - 1\ndog - 1\n"
	if output.String() != expectedOutput {
		t.Errorf("Expected output %q, but got %q", expectedOutput, output.String())
	}
}