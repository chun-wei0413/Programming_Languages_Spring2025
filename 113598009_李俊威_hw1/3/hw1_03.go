package main

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
)

// readFile reads the contents of a file and returns it as a string
func readFile(pathToFile string) (string, error) {
	content, err := os.ReadFile(pathToFile)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// filterCharsAndNormalize replaces non-alphanumeric characters with spaces and converts to lowercase
func filterCharsAndNormalize(strData string) string {
	re := regexp.MustCompile(`[\W_]+`)
	return strings.ToLower(re.ReplaceAllString(strData, " "))
}

// scan splits a string into words
func scan(strData string) []string {
	return strings.Fields(strData)
}

// removeStopWords removes common stop words and single letters from a word list
func removeStopWords(wordList []string) ([]string, error) {

	stopWordsPath := "../../stop_words.txt"
	// Read stop words file
	content, err := os.ReadFile(stopWordsPath)
	if err != nil {
		return nil, err
	}
	
	stopWords := make(map[string]bool)
	for _, word := range strings.Split(string(content), ",") {
		stopWords[word] = true
	}
	
	// Add single letter words as stop words
	for c := 'a'; c <= 'z'; c++ {
		stopWords[string(c)] = true
	}
	
	// Filter out stop words
	filteredWords := make([]string, 0)
	for _, word := range wordList {
		if !stopWords[word] {
			filteredWords = append(filteredWords, word)
		}
	}
	
	return filteredWords, nil
}

// frequencies counts word occurrences and returns a map of word to count
func frequencies(wordList []string) map[string]int {
	wordFreqs := make(map[string]int)
	for _, word := range wordList {
		wordFreqs[word]++
	}
	return wordFreqs
}

// sortFrequencies sorts word frequencies in descending order
type wordFreq struct {
	Word  string
	Count int
}

func sortFrequencies(wordFreqs map[string]int) []wordFreq {
	pairs := make([]wordFreq, 0, len(wordFreqs))
	for word, count := range wordFreqs {
		pairs = append(pairs, wordFreq{Word: word, Count: count})
	}
	
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Count > pairs[j].Count
	})
	
	return pairs
}

// printAll prints word frequencies recursively
// func printAll(wordFreqs []wordFreq, index int) {
// 	if index < len(wordFreqs) {
// 		fmt.Printf("%s - %d\n", wordFreqs[index].Word, wordFreqs[index].Count)
// 		printAll(wordFreqs, index+1)
// 	}
// }


// go run hw1_03.go ../../pride-and-prejudice.txt
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run hw1_03.go <file_path>")
		os.Exit(1)
	}

	filePath := os.Args[1]

	//這種寫法有點不太直觀
	sorted := then(then(thenError(then(then(pipeOf(readFile(filePath)),filterCharsAndNormalize),scan),removeStopWords),frequencies),sortFrequencies)

	if sorted.err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", sorted.err)
		os.Exit(1)
	}

	limit := 25
	if len(sorted.value) < limit {
		limit = len(sorted.value)
	}
	for i := 0; i < limit; i++ {
		fmt.Printf("%s - %d\n", sorted.value[i].Word, sorted.value[i].Count)
	}
}


// **泛型函數擴展 pipeline**
type pipe[T any] struct {
	value T
	err   error
}

// 創建一個 pipeline
func pipeOf[T any](value T, err error) pipe[T] {
	return pipe[T]{value, err}
}

func then[T, U any](p pipe[T], f func(T) U) pipe[U] {
	if p.err != nil {
		var zero U
		return pipe[U]{zero, p.err} // 回傳零值
	}
	return pipeOf(f(p.value), nil)
}

func thenError[T, U any](p pipe[T], f func(T) (U, error)) pipe[U] {
	if p.err != nil {
		var zero U
		return pipe[U]{zero, p.err}
	}
	v, err := f(p.value)
	return pipeOf(v, err)
}