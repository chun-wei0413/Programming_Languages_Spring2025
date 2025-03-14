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

	content := "The cat, and the dog! 123 Cat HERE"
	filePath := setupTestFile(t, content)
	defer cleanupTestFile(t, filePath)

	dsm := NewDataStorageManager(filePath)
	words := dsm.Words()

	expected := []string{"the", "cat", "and", "the", "dog", "cat", "here"}
	if !reflect.DeepEqual(expected, words) {
		t.Errorf("Expected words %v, but got %v", expected, words)
	}
}

func TestStopWordManager_IsStopWord(t *testing.T) {
	//為了不要耦合stop_words.txt我寫了一個stopWordsContent
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
		if test.expected != result {
			t.Errorf("For word %q, expected %v, but got %v", test.word, test.expected, result)
		}
	}
}

func TestWordFrequencyManager_IncrementCount(t *testing.T) {
	wfm := NewWordFrequencyManager()

	wfm.IncrementCount("cat")
	wfm.IncrementCount("dog")
	wfm.IncrementCount("cat")

	if wfm.wordFreqs["cat"] != 2 {
		t.Errorf("Expected frequency of 'cat' to be 2, but got %d", wfm.wordFreqs["cat"])
	}
	if wfm.wordFreqs["dog"] != 1 {
		t.Errorf("Expected frequency of 'dog' to be 1, but got %d", wfm.wordFreqs["dog"])
	}
}

func TestWordFrequencyManager_Sorted(t *testing.T) {
	wfm := NewWordFrequencyManager()

	// 先增加一些詞頻
	wfm.IncrementCount("pig")
	wfm.IncrementCount("bird")
	wfm.IncrementCount("pig")

	sorted := wfm.Sorted()
	expected := [][2]string{
		{"pig", "2"},
		{"bird", "1"},
	}

	if !reflect.DeepEqual(expected, sorted) {
		t.Errorf("Expected sorted result %v, but got %v", expected, sorted)
	}
}

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
	if expectedOutput != output.String() {
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
	if expectedOutput != output.String() {
		t.Errorf("Expected output %q, but got %q", expectedOutput, output.String())
	}
}