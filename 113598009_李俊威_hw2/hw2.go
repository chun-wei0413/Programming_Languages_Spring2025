package main

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// FileHandler encapsulates file operations
type FileHandler struct {
	path string
}

func NewFileHandler(path string) *FileHandler {
	return &FileHandler{path: path}
}

func (f *FileHandler) ReadContent() (string, error) {
	content, err := os.ReadFile(f.path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// WordProcessor encapsulates word processing operations
type WordProcessor struct{}

func NewWordProcessor() *WordProcessor {
	return &WordProcessor{}
}

func (w *WordProcessor) Normalize(text string) string {
	return strings.ToLower(text)
}

func (w *WordProcessor) SplitWords(text string) []string {
	// Use regexp to keep only letters
	re := regexp.MustCompile(`[^a-zA-Z]+`)
	replaced := re.ReplaceAllString(text, " ")
	words := strings.Fields(replaced)
	// Filter out empty or invalid words
	var validWords []string
	for _, word := range words {
		if len(word) > 0 && regexp.MustCompile(`^[a-z]+$`).MatchString(word) {
			validWords = append(validWords, word)
		}
	}
	return validWords
}

// StopWordManager encapsulates stop word filtering
type StopWordManager struct {
	stopWords map[string]struct{}
}

func NewStopWordManager() *StopWordManager {
	sw := &StopWordManager{stopWords: make(map[string]struct{})}

	// Load stop words from file in current directory
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

	// Add single-letter words and common titles
	for c := 'a'; c <= 'z'; c++ {
		sw.stopWords[string(c)] = struct{}{}
	}
	
	return sw
}

func (s *StopWordManager) IsStopWord(word string) bool {
	_, exists := s.stopWords[word]
	return exists
}

// FrequencyCounter encapsulates frequency counting operations
type FrequencyCounter struct {
	counts map[string]int
}

func NewFrequencyCounter() *FrequencyCounter {
	return &FrequencyCounter{counts: make(map[string]int)}
}

func (f *FrequencyCounter) CountWords(words []string, stopManager *StopWordManager) {
	for _, word := range words {
		if len(word) > 1 && !stopManager.IsStopWord(word) {
			f.counts[word]++
		}
	}
}

func (f *FrequencyCounter) GetTopN(n int) [][2]string {
	// Convert map to slice of pairs
	pairs := make([][2]string, 0, len(f.counts))
	for word, count := range f.counts {
		pairs = append(pairs, [2]string{word, strconv.Itoa(count)})
	}

	// Sort by count (descending) and then by word (ascending)
	sort.Slice(pairs, func(i, j int) bool {
		countI, _ := strconv.Atoi(pairs[i][1])
		countJ, _ := strconv.Atoi(pairs[j][1])
		if countI != countJ {
			return countI > countJ // Descending by count
		}
		return pairs[i][0] < pairs[j][0] // Ascending by word if counts are equal
	})

	// Return top N or all if less than N
	if len(pairs) < n {
		return pairs
	}
	return pairs[:n]
}

// WordFrequencyController orchestrates the word frequency counting process
type WordFrequencyController struct {
	fileHandler   *FileHandler
	wordProcessor *WordProcessor
	stopManager   *StopWordManager
	freqCounter   *FrequencyCounter
}

func NewWordFrequencyController(path string) *WordFrequencyController {
	return &WordFrequencyController{
		fileHandler:   NewFileHandler(path),
		wordProcessor: NewWordProcessor(),
		stopManager:   NewStopWordManager(),
		freqCounter:   NewFrequencyCounter(),
	}
}

func (w *WordFrequencyController) Run() error {
	// Read file content
	content, err := w.fileHandler.ReadContent()
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	// Process words
	normalized := w.wordProcessor.Normalize(content)
	words := w.wordProcessor.SplitWords(normalized)
	w.freqCounter.CountWords(words, w.stopManager)

	// Get and print top 25
	top25 := w.freqCounter.GetTopN(25)
	for _, pair := range top25 {
		fmt.Printf("%s - %s\n", pair[0], pair[1])
	}

	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: program <filename>")
		os.Exit(1)
	}

	controller := NewWordFrequencyController(os.Args[1])
	if err := controller.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}