package main

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// DataStorageManager models the contents of the file
type DataStorageManager struct {
	data string
}

func NewDataStorageManager(path string) *DataStorageManager {
	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}
	// Replace non-word characters with space and convert to lowercase
	pattern := regexp.MustCompile(`[\W_]+`)
	data := pattern.ReplaceAllString(string(content), " ")
	return &DataStorageManager{data: strings.ToLower(data)}
}

func (d *DataStorageManager) Words() []string {
	words := strings.Fields(d.data)
	var validWords []string
	for _, word := range words {
		if len(word) > 0 && regexp.MustCompile(`^[a-z]+$`).MatchString(word) {
			validWords = append(validWords, word)
		}
	}
	return validWords
}

// StopWordManager models the stop word filter
type StopWordManager struct {
	stopWords map[string]struct{}
}

func NewStopWordManager() *StopWordManager {
	sw := &StopWordManager{stopWords: make(map[string]struct{})}
	// Load stop words from file
	content, err := os.ReadFile("./stop_words.txt")
	if err == nil {
		words := strings.Split(string(content), ",")
		for _, word := range words {
			trimmedWord := strings.TrimSpace(word)
			if len(trimmedWord) > 0 {
				sw.stopWords[trimmedWord] = struct{}{}
			}
		}
	} else {
		fmt.Println("Warning: Could not load stop_words.txt:", err)
	}
	// Add single-letter words
	for c := 'a'; c <= 'z'; c++ {
		sw.stopWords[string(c)] = struct{}{}
	}
	return sw
}

func (s *StopWordManager) IsStopWord(word string) bool {
	_, exists := s.stopWords[word]
	return exists
}

type WordFrequencyManager struct {
	wordFreqs map[string]int
}

func NewWordFrequencyManager() *WordFrequencyManager {
	return &WordFrequencyManager{wordFreqs: make(map[string]int)}
}

func (w *WordFrequencyManager) IncrementCount(word string) {
	if _, exists := w.wordFreqs[word]; exists {
		w.wordFreqs[word]++
	} else {
		w.wordFreqs[word] = 1
	}
}

func (w *WordFrequencyManager) Sorted() [][2]string {
	pairs := make([][2]string, 0, len(w.wordFreqs))
	for word, count := range w.wordFreqs {
		pairs = append(pairs, [2]string{word, fmt.Sprintf("%d", count)})
	}
	sort.Slice(pairs, func(i, j int) bool {
		countI, _ := strconv.Atoi(pairs[i][1])
		countJ, _ := strconv.Atoi(pairs[j][1])
		if countI != countJ {
			return countI > countJ
		}
		return pairs[i][0] < pairs[j][0]
	})
	return pairs
}

// WordFrequencyController orchestrates the word frequency counting process
type WordFrequencyController struct {
	storageManager   *DataStorageManager
	stopWordManager  *StopWordManager
	wordFreqManager  *WordFrequencyManager
}

func NewWordFrequencyController(path string) *WordFrequencyController {
	return &WordFrequencyController{
		storageManager:  NewDataStorageManager(path),
		stopWordManager: NewStopWordManager(),
		wordFreqManager: NewWordFrequencyManager(),
	}
}

func (w *WordFrequencyController) Run() {
	for _, word := range w.storageManager.Words() {
		if !w.stopWordManager.IsStopWord(word) {
			w.wordFreqManager.IncrementCount(word)
		}
	}
	wordFreqs := w.wordFreqManager.Sorted()
	for i, pair := range wordFreqs {
		if i >= 25 {
			break
		}
		fmt.Printf("%s - %s\n", pair[0], pair[1])
	}
}

func main() {
	//請使用以下指令 by Frank
	//go run hw2.go ./pride-and-prejudice.txt
	if len(os.Args) != 2 {
		fmt.Println("Usage: program <filename>")
		os.Exit(1)
	}
	controller := NewWordFrequencyController(os.Args[1])
	controller.Run()
}