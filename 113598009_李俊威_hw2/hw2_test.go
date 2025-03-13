package main

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

// 模擬檔案內容的輔助函數
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

// 清理臨時檔案
func cleanupTestFile(t *testing.T, path string) {
	err := os.Remove(path)
	if err != nil {
		t.Fatalf("Failed to remove temp file: %v", err)
	}
}

// 測試 DataStorageManager
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
	// 模擬 stop_words.txt 檔案
	stopWordsContent := "the, and, is"
	stopFile := setupTestFile(t, stopWordsContent)
	defer cleanupTestFile(t, stopFile)

	// 創建 StopWordManager，並模擬檔案內容
	swm := &StopWordManager{stopWords: make(map[string]struct{})}
	content, _ := os.ReadFile(stopFile)
	words := strings.Split(string(content), ",")
	for _, word := range words {
		trimmedWord := strings.TrimSpace(word)
		if len(trimmedWord) > 0 {
			swm.stopWords[trimmedWord] = struct{}{}
		}
	}
	// 添加單字母詞
	for c := 'a'; c <= 'z'; c++ {
		swm.stopWords[string(c)] = struct{}{}
	}

	tests := []struct {
		word     string
		expected bool
	}{
		{"the", true},    // 在停用詞列表中
		{"and", true},    // 在停用詞列表中
		{"is", true},     // 在停用詞列表中
		{"a", true},      // 單字母停用詞
		{"cat", false},   // 不在停用詞列表中
		{"dog", false},   // 不在停用詞列表中
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

	// 測試 Sorted
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
	// 模擬輸入檔案
	content := "The cat and the dog. Cat is here!"
	filePath := setupTestFile(t, content)
	defer cleanupTestFile(t, filePath)

	// 模擬 stop_words.txt
	stopWordsContent := "the, and, is"
	stopFile := setupTestFile(t, stopWordsContent)
	defer cleanupTestFile(t, stopFile)

	// 創建 StopWordManager 並手動填充停用詞
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

	// 重定向輸出以檢查結果
	originalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// 手動創建 controller，避免依賴 NewStopWordManager
	controller := &WordFrequencyController{
		storageManager:  NewDataStorageManager(filePath),
		stopWordManager: swm,
		wordFreqManager: NewWordFrequencyManager(),
	}
	controller.Run()

	w.Close()
	os.Stdout = originalStdout

	// 讀取輸出
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

// 測試主程式流程（模擬 main）
func TestMainFlow(t *testing.T) {
	// 模擬命令列參數
	originalArgs := os.Args
	os.Args = []string{"program", "testfile.txt"}
	defer func() { os.Args = originalArgs }()

	// 模擬輸入檔案
	content := "The cat and the dog."
	filePath := setupTestFile(t, content)
	defer cleanupTestFile(t, filePath)
	os.Args[1] = filePath

	// 模擬 stop_words.txt
	stopWordsContent := "the, and"
	stopFile := setupTestFile(t, stopWordsContent)
	defer cleanupTestFile(t, stopFile)

	// 重定向輸出
	originalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// 執行 main
	main()

	w.Close()
	os.Stdout = originalStdout

	// 讀取輸出
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