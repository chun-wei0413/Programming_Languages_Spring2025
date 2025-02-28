#!/usr/bin/env python
import unittest
import os
import tempfile


# Import the functions from the main program
from tf_06 import read_file, filter_chars_and_normalize, scan, remove_stop_words, frequencies, sort

class TestCurriedWordFrequency(unittest.TestCase):
    
    def setUp(self):
        # Create a temporary test file with some content
        self.text_fd, self.text_path = tempfile.mkstemp()
        with open(self.text_path, 'w') as f:
            f.write("This is a test file. It contains some test words, test words that repeat.")
        
        # Create a temporary stop words file
        self.stop_fd, self.stop_path = tempfile.mkstemp()
        with open(self.stop_path, 'w') as f:
            f.write("a,is,it,this,some")
    
    def tearDown(self):
        # Clean up the temporary files
        os.close(self.text_fd)
        os.remove(self.text_path)
        os.close(self.stop_fd)
        os.remove(self.stop_path)
    
    def test_read_file(self):
        content = read_file(self.text_path)
        self.assertTrue(len(content) > 0)
        self.assertTrue("test file" in content)
    
    def test_filter_chars_and_normalize(self):
        text = "Hello, World! This is a TEST."
        filtered = filter_chars_and_normalize(text)
        self.assertEqual(filtered, "hello world this is a test ")
    
    def test_scan(self):
        text = "hello world this is a test"
        words = scan(text)
        self.assertEqual(words, ["hello", "world", "this", "is", "a", "test"])
    
    def test_remove_stop_words_curried(self):
        # Test the curried function
        remove_stop_words_func = remove_stop_words(self.stop_path)
        words = ["this", "is", "a", "test", "file", "with", "test", "words"]
        filtered_words = remove_stop_words_func(words)
        
        # Should remove 'this', 'is', 'a', and single letters
        expected = ["test", "file", "with", "test", "words"]
        self.assertEqual(filtered_words, expected)
    
    def test_frequencies(self):
        words = ["test", "file", "test", "words", "test"]
        freqs = frequencies(words)
        self.assertEqual(freqs, {"test": 3, "file": 1, "words": 1})
    
    def test_sort(self):
        freqs = {"test": 3, "file": 1, "words": 1}
        sorted_freqs = sort(freqs)
        self.assertEqual(sorted_freqs, [("test", 3), ("file", 1), ("words", 1)])
    
    def test_full_pipeline(self):
        # Test the full pipeline with the curried function
        text = read_file(self.text_path)
        normalized = filter_chars_and_normalize(text)
        words = scan(normalized)
        remove_stop_words_func = remove_stop_words(self.stop_path)
        filtered_words = remove_stop_words_func(words)
        word_freqs = frequencies(filtered_words)
        sorted_freqs = sort(word_freqs)
        
        # Verify expected results
        self.assertTrue(("test", 3) in sorted_freqs)
        self.assertTrue(("words", 2) in sorted_freqs)
        
        # Check that stop words were removed
        for word in ["a", "is", "it", "this", "some"]:
            self.assertFalse(any(w[0] == word for w in sorted_freqs))

if __name__ == "__main__":
    # python -m unittest tf_06_test.py
    unittest.main()