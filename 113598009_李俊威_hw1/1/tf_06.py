#!/usr/bin/env python
import sys, re, operator, string

#
#python tf_06.py ../../pride-and-prejudice.txt ../../stop_words.txt
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

def remove_stop_words(stop_words_file):
    """ 
    Takes a path to the stop words file and returns a function that takes
    a list of words and returns a copy with all stop words removed.
    This is an example of currying.
    """
    def _remove_stop_words(word_list):
        with open(stop_words_file) as f:
            stop_words = f.read().split(',')
        # add single-letter words
        stop_words.extend(list(string.ascii_lowercase))
        return [w for w in word_list if not w in stop_words]
    return _remove_stop_words

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
    if(len(word_freqs) > 0):
        print(word_freqs[0][0], '-', word_freqs[0][1])
        print_all(word_freqs[1:])

#
# The main function
#
if __name__ == "__main__":
    if len(sys.argv) < 3:
        print("Usage: python3 tf_06.py <text_file> <stop_words_file>")
        sys.exit(1)
    # Apply each function in the pipeline, passing the result to the next function
    # Get the stop_words_file from sys.argv[2] and create the curried function
    remove_stop_words_with_file = remove_stop_words(sys.argv[2])
    # Use the curried function in the pipeline
    print_all(sort(frequencies(remove_stop_words_with_file(scan(filter_chars_and_normalize(read_file(sys.argv[1]))))))[0:25])