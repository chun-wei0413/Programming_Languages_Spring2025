package main

import (
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
)

// 測試 FileHandler
func TestFileHandler(t *testing.T) {
	// 建立臨時測試檔案
	content := "這是測試檔案的內容。\n它有一些單詞要計數。"
	tmpFile, err := os.CreateTemp("", "test_file.txt")
	if err != nil {
		t.Fatalf("無法建立臨時檔案: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	
	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("無法寫入臨時檔案: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("無法關閉臨時檔案: %v", err)
	}
	
	// 測試 ReadContent 方法
	fileHandler := NewFileHandler(tmpFile.Name())
	readContent, err := fileHandler.ReadContent()
	if err != nil {
		t.Fatalf("ReadContent 失敗: %v", err)
	}
	if readContent != content {
		t.Errorf("ReadContent 內容不匹配: 得到 %q, 預期 %q", readContent, content)
	}
	
	// 測試不存在的檔案
	nonExistentHandler := NewFileHandler("non_existent_file.txt")
	_, err = nonExistentHandler.ReadContent()
	if err == nil {
		t.Errorf("預期讀取不存在檔案時會出錯，但沒有")
	}
}

// 測試 WordProcessor
func TestWordProcessor(t *testing.T) {
	wordProcessor := NewWordProcessor()
	
	// 測試 Normalize 方法
	testCases := []struct {
		input    string
		expected string
	}{
		{"Hello World!", "hello world!"},
		{"TESTING", "testing"},
		{"MiXeD CaSe", "mixed case"},
		{"", ""},
	}
	
	for _, tc := range testCases {
		normalized := wordProcessor.Normalize(tc.input)
		if normalized != tc.expected {
			t.Errorf("Normalize 結果不匹配: 輸入 %q, 得到 %q, 預期 %q", tc.input, normalized, tc.expected)
		}
	}
	
	// 測試 SplitWords 方法
	splitTestCases := []struct {
		input    string
		expected []string
	}{
		{"Hello, World! This is a test 123.", []string{"hello", "world", "this", "is", "a", "test"}},
		{"word1, word2; word3.", []string{"word1", "word2", "word3"}},
		{"123, 456, 789", []string{}}, // 只有數字，應該回傳空陣列
		{"", []string{}},
	}
	
	for _, tc := range splitTestCases {
		words := wordProcessor.SplitWords(tc.input)
		if !reflect.DeepEqual(words, tc.expected) {
			t.Errorf("SplitWords 結果不匹配: 輸入 %q, 得到 %v, 預期 %v", tc.input, words, tc.expected)
		}
	}
}

// 測試 StopWordManager
func TestStopWordManager(t *testing.T) {
	// 建立臨時停用詞檔案
	stopWords := "a,the,and,of,to"
	tmpFile, err := os.CreateTemp("", "stop_words.txt")
	if err != nil {
		t.Fatalf("無法建立臨時停用詞檔案: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	
	if _, err := tmpFile.WriteString(stopWords); err != nil {
		t.Fatalf("無法寫入臨時停用詞檔案: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("無法關閉臨時停用詞檔案: %v", err)
	}
	
	// 複製停用詞檔案到當前目錄
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("無法獲取當前目錄: %v", err)
	}
	currentStopWordsPath := currentDir + "/stop_words.txt"
	stopWordsContent, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("無法讀取臨時停用詞檔案: %v", err)
	}
	
	// 保存原始檔案（如果存在）
	originalExists := false
	originalContent := []byte{}
	if _, err := os.Stat(currentStopWordsPath); err == nil {
		originalExists = true
		originalContent, err = os.ReadFile(currentStopWordsPath)
		if err != nil {
			t.Fatalf("無法讀取原始停用詞檔案: %v", err)
		}
	}
	
	// 寫入測試停用詞檔案
	err = os.WriteFile(currentStopWordsPath, stopWordsContent, 0644)
	if err != nil {
		t.Fatalf("無法在當前目錄寫入停用詞檔案: %v", err)
	}
	
	// 測試後恢復原始檔案或刪除測試檔案
	defer func() {
		if originalExists {
			os.WriteFile(currentStopWordsPath, originalContent, 0644)
		} else {
			os.Remove(currentStopWordsPath)
		}
	}()
	
	// 測試 IsStopWord 方法
	stopWordManager := NewStopWordManager()
	
	// 測試檔案中的停用詞
	if !stopWordManager.IsStopWord("a") {
		t.Errorf("IsStopWord 無法識別 'a' 為停用詞")
	}
	if !stopWordManager.IsStopWord("the") {
		t.Errorf("IsStopWord 無法識別 'the' 為停用詞")
	}
	
	// 測試為單字母詞
	if !stopWordManager.IsStopWord("x") {
		t.Errorf("IsStopWord 無法識別 'x' 為停用詞")
	}
	if !stopWordManager.IsStopWord("z") {
		t.Errorf("IsStopWord 無法識別 'z' 為停用詞")
	}
	
	// 測試非停用詞
	if stopWordManager.IsStopWord("hello") {
		t.Errorf("IsStopWord 錯誤地識別 'hello' 為停用詞")
	}
	if stopWordManager.IsStopWord("world") {
		t.Errorf("IsStopWord 錯誤地識別 'world' 為停用詞")
	}
}

// 測試 FrequencyCounter
func TestFrequencyCounter(t *testing.T) {
	// 建立帶有某些停用詞的停用詞管理器
	stopWordManager := &StopWordManager{
		stopWords: map[string]struct{}{
			"a": {}, "the": {}, "and": {}, "of": {}, "to": {},
		},
	}
	for c := 'a'; c <= 'z'; c++ {
		stopWordManager.stopWords[string(c)] = struct{}{}
	}
	
	// 測試 CountWords 方法
	freqCounter := NewFrequencyCounter()
	words := []string{"hello", "world", "hello", "test", "world", "world", "test", "the", "a", "of"}
	freqCounter.CountWords(words, stopWordManager)
	
	// 確認預期的計數結果
	expectedCounts := map[string]int{
		"hello": 2,
		"world": 3,
		"test":  2,
	}
	
	for word, expectedCount := range expectedCounts {
		actualCount, exists := freqCounter.counts[word]
		if !exists {
			t.Errorf("單詞 %q 應該存在於計數中但沒有", word)
		} else if actualCount != expectedCount {
			t.Errorf("單詞 %q 的計數不正確: 得到 %d, 預期 %d", word, actualCount, expectedCount)
		}
	}
	
	// 測試停用詞不應被計數
	for stopWord := range stopWordManager.stopWords {
		if _, exists := freqCounter.counts[stopWord]; exists && len(stopWord) > 1 {
			t.Errorf("停用詞 %q 不應被計數但存在於計數中", stopWord)
		}
	}
	
	// 測試 GetTopN 方法
	expectedTopResults := [][2]string{
		{"world", "3"},
		{"hello", "2"},
		{"test", "2"},
	}
	
	topResults := freqCounter.GetTopN(3)
	if !reflect.DeepEqual(topResults, expectedTopResults) {
		t.Errorf("GetTopN 結果不匹配: 得到 %v, 預期 %v", topResults, expectedTopResults)
	}
	
	// 測試 N 大於結果數量的情況
	largeNResults := freqCounter.GetTopN(10)
	if len(largeNResults) != 3 {
		t.Errorf("當 N 大於結果數量時，GetTopN 應返回所有結果: 得到 %d 項，預期 3 項", len(largeNResults))
	}
	
	// 測試空結果的情況
	emptyCounter := NewFrequencyCounter()
	emptyResults := emptyCounter.GetTopN(5)
	if len(emptyResults) != 0 {
		t.Errorf("空計數器的 GetTopN 應返回空陣列，得到 %v", emptyResults)
	}
}

// 測試 WordFrequencyController
func TestWordFrequencyController(t *testing.T) {
	// 建立帶有一些內容的臨時檔案
	content := "Hello world. This is a test. Hello again, world!"
	tmpFile, err := os.CreateTemp("", "test_wfc.txt")
	if err != nil {
		t.Fatalf("無法建立臨時檔案: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	
	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("無法寫入臨時檔案: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("無法關閉臨時檔案: %v", err)
	}
	
	// 建立臨時停用詞檔案
	stopWords := "a,the,and,of,to,is"
	stopWordsFile, err := os.CreateTemp("", "stop_words.txt")
	if err != nil {
		t.Fatalf("無法建立臨時停用詞檔案: %v", err)
	}
	defer os.Remove(stopWordsFile.Name())
	
	if _, err := stopWordsFile.WriteString(stopWords); err != nil {
		t.Fatalf("無法寫入臨時停用詞檔案: %v", err)
	}
	if err := stopWordsFile.Close(); err != nil {
		t.Fatalf("無法關閉臨時停用詞檔案: %v", err)
	}
	
	// 複製停用詞檔案到當前目錄
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("無法獲取當前目錄: %v", err)
	}
	currentStopWordsPath := currentDir + "/stop_words.txt"
	stopWordsContent, err := os.ReadFile(stopWordsFile.Name())
	if err != nil {
		t.Fatalf("無法讀取臨時停用詞檔案: %v", err)
	}
	
	// 保存原始檔案（如果存在）
	originalExists := false
	originalContent := []byte{}
	if _, err := os.Stat(currentStopWordsPath); err == nil {
		originalExists = true
		originalContent, err = os.ReadFile(currentStopWordsPath)
		if err != nil {
			t.Fatalf("無法讀取原始停用詞檔案: %v", err)
		}
	}
	
	// 寫入測試停用詞檔案
	err = os.WriteFile(currentStopWordsPath, stopWordsContent, 0644)
	if err != nil {
		t.Fatalf("無法在當前目錄寫入停用詞檔案: %v", err)
	}
	
	// 測試後恢復原始檔案或刪除測試檔案
	defer func() {
		if originalExists {
			os.WriteFile(currentStopWordsPath, originalContent, 0644)
		} else {
			os.Remove(currentStopWordsPath)
		}
	}()
	
	// 暫時重定向 stdout 以捕獲輸出
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	
	// 執行控制器
	controller := NewWordFrequencyController(tmpFile.Name())
	controller.Run()
	
	// 恢復 stdout 並讀取捕獲的輸出
	w.Close()
	os.Stdout = oldStdout
	var buf strings.Builder
	_, err = io.Copy(&buf, r)
	if err != nil {
		t.Fatalf("無法捕獲 stdout: %v", err)
	}
	
	// 檢查輸出 - 順序可能不同，所以只檢查是否存在
	output := buf.String()
	expectedWords := []string{"hello", "world", "test"}
	for _, word := range expectedWords {
		if !strings.Contains(output, word) {
			t.Errorf("在輸出中找不到單詞 %q", word)
		}
	}
	
	// 測試錯誤處理
	errorController := NewWordFrequencyController("non_existent_file.txt")
	err = errorController.Run()
	if err == nil {
		t.Errorf("對於不存在的檔案，Run 方法應返回錯誤，但沒有")
	}
}

// 測試多個元件的整合
func TestIntegration(t *testing.T) {
	// 建立帶有一些內容的臨時檔案
	content := "Hello, World! This is a test. Hello again, world! The test is working. And it should count words properly."
	tmpFile, err := os.CreateTemp("", "integration_test.txt")
	if err != nil {
		t.Fatalf("無法建立臨時檔案: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	
	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("無法寫入臨時檔案: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("無法關閉臨時檔案: %v", err)
	}
	
	// 建立臨時停用詞檔案
	stopWords := "a,the,and,of,to,is,it"
	stopWordsFile, err := os.CreateTemp("", "stop_words.txt")
	if err != nil {
		t.Fatalf("無法建立臨時停用詞檔案: %v", err)
	}
	defer os.Remove(stopWordsFile.Name())
	
	if _, err := stopWordsFile.WriteString(stopWords); err != nil {
		t.Fatalf("無法寫入臨時停用詞檔案: %v", err)
	}
	if err := stopWordsFile.Close(); err != nil {
		t.Fatalf("無法關閉臨時停用詞檔案: %v", err)
	}
	
	// 複製停用詞檔案到當前目錄
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("無法獲取當前目錄: %v", err)
	}
	currentStopWordsPath := currentDir + "/stop_words.txt"
	stopWordsContent, err := os.ReadFile(stopWordsFile.Name())
	if err != nil {
		t.Fatalf("無法讀取臨時停用詞檔案: %v", err)
	}
	
	// 保存原始檔案（如果存在）
	originalExists := false
	originalContent := []byte{}
	if _, err := os.Stat(currentStopWordsPath); err == nil {
		originalExists = true
		originalContent, err = os.ReadFile(currentStopWordsPath)
		if err != nil {
			t.Fatalf("無法讀取原始停用詞檔案: %v", err)
		}
	}
	
	// 寫入測試停用詞檔案
	err = os.WriteFile(currentStopWordsPath, stopWordsContent, 0644)
	if err != nil {
		t.Fatalf("無法在當前目錄寫入停用詞檔案: %v", err)
	}
	
	// 測試後恢復原始檔案或刪除測試檔案
	defer func() {
		if originalExists {
			os.WriteFile(currentStopWordsPath, originalContent, 0644)
		} else {
			os.Remove(currentStopWordsPath)
		}
	}()
	
	// 手動執行各元件以驗證其整合
	fileHandler := NewFileHandler(tmpFile.Name())
	wordProcessor := NewWordProcessor()
	stopWordManager := NewStopWordManager()
	freqCounter := NewFrequencyCounter()
	
	content, err = fileHandler.ReadContent()
	if err != nil {
		t.Fatalf("FileHandler.ReadContent 失敗: %v", err)
	}
	
	normalized := wordProcessor.Normalize(content)
	words := wordProcessor.SplitWords(normalized)
	
	freqCounter.CountWords(words, stopWordManager)
	topWords := freqCounter.GetTopN(5)
	
	// 驗證預期的單詞是否在頂部單詞中
	expectedWords := map[string]bool{
		"hello": true, "world": true, "test": true,
	}
	
	for _, pair := range topWords {
		word := pair[0]
		if _, found := expectedWords[word]; found {
			delete(expectedWords, word)
		}
	}
	
	if len(expectedWords) > 0 {
		t.Errorf("某些預期的單詞在頂部結果中未找到: %v", expectedWords)
	}
}