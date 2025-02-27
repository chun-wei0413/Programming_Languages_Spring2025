#!/usr/bin/env python
import sys, re, operator, string

#
# The functions
#

def read_file(path_to_file):
    """
    Takes a path to a file and returns the entire
    contents of the file as a string
    """
    with open(path_to_file) as f:
        data = f.read()
    return data

def filter_chars_and_normalize(str_data):
    """
    Takes a string and returns a copy with all nonalphanumeric 
    chars replaced by white space
    """
    pattern = re.compile(r'[\W_]+')
    return pattern.sub(' ', str_data).lower()

def scan(str_data):
    """
    Takes a string and scans for words, returning
    a list of words.
    """
    return str_data.split()

def remove_stop_words_curry(stop_words):
    """ 
    Takes a list of stop words and returns a function that
    removes stop words from any given word list.
    """
    def remove(word_list):
        return [w for w in word_list if w not in stop_words]
    return remove

def frequencies(word_list):
    """
    Takes a list of words and returns a dictionary associating
    words with frequencies of occurrence
    """
    word_freqs = {}
    for w in word_list:
        if w in word_freqs:
            word_freqs[w] += 1
        else:
            word_freqs[w] = 1
    return word_freqs

def sort(word_freq):
    """
    Takes a dictionary of words and their frequencies
    and returns a list of pairs where the entries are
    sorted by frequency 
    """
    return sorted(word_freq.items(), key=operator.itemgetter(1), reverse=True)

def print_all(word_freqs):
    """
    Takes a list of pairs where the entries are sorted by frequency and print them recursively.
    """
    if len(word_freqs) > 0:
        print(word_freqs[0][0], '-', word_freqs[0][1])
        print_all(word_freqs[1:])

#
# The main function
#
if __name__ == "__main__":
    # 讀取 stop words 並創建柯里化函式
    with open('../stop_words.txt') as f:
        stop_words = f.read().split(',')
    stop_words.extend(list(string.ascii_lowercase))

    # 取得柯里化版本的函式
    remove_stop_words = remove_stop_words_curry(stop_words)

    # 執行 pipeline
    print_all(sort(frequencies(remove_stop_words(scan(filter_chars_and_normalize(read_file(sys.argv[1]))))))[0:25])
