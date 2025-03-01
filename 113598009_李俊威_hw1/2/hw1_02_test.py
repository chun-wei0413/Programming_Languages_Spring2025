import unittest
# from hw1_02 import filter_chars_and_normalize
import hw1_02

class TestTF05(unittest.TestCase):

    def setUp(self):
        # Reset global variables to ensure a clean test environment
        hw1_02.data = []
        hw1_02.words = []
        hw1_02.word_freqs = []

    def test_read_file_not_idempotent(self):
        """
        read_file is not idempotent because each execution appends
        the content to data.
        """
        expected_content = "existing data" + "This is a test. This test is only a test." * 2
        test_filename = "test.txt"
        
        # Create a temporary test file
        with open(test_filename, "w") as f:
            f.write("This is a test. This test is only a test.")

        try:
            hw1_02.data = list("existing data")
            hw1_02.read_file(test_filename)
            hw1_02.read_file(test_filename)
            
            # Should contain original data plus twice-read content
            self.assertIn("existing data", ''.join(hw1_02.data))
            self.assertEqual(expected_content, "".join(hw1_02.data))
        finally:
            # Clean up the file after the test
            import os
            os.remove(test_filename)

    def test_filter_chars_and_normalize_idempotent(self):
        """
        filter_chars_and_normalize is idempotent because multiple executions 
        produce the same result.
        """
        sentence = "This is a test. This test is only a test."
        expected = "this is a test  this test is only a test "
        hw1_02.data = list(sentence)
        hw1_02.filter_chars_and_normalize()
        hw1_02.filter_chars_and_normalize()
        self.assertEqual(expected, ''.join(hw1_02.data))

    def test_scan_not_idempotent(self):
        """
        scan is not idempotent because it appends words to words list on each execution.
        """
        hw1_02.data = list("hello world hello world")
        hw1_02.scan()
        hw1_02.scan()
        
        # Expected word count should double
        self.assertEqual(4, hw1_02.words.count("hello"))
        self.assertEqual(4, hw1_02.words.count("world"))

    def test_remove_stop_words_idempotent(self):
        """
        remove_stop_words is idempotent because executing it multiple times
        does not affect the result.
        """
        hw1_02.words = ['test', 'test', 'test']
        hw1_02.remove_stop_words()
        hw1_02.remove_stop_words()
        self.assertEqual(['test', 'test', 'test'], hw1_02.words)

    def test_frequencies_not_idempotent(self):
        """
        frequencies is not idempotent because it updates word_freqs each time
        based on words without resetting, leading to accumulated counts.
        """
        hw1_02.words = ["test", "test", "test"]
        hw1_02.frequencies()
        hw1_02.frequencies()
        
        # Expected test count should be doubled
        self.assertEqual([["test", 6]], hw1_02.word_freqs)
    
    def test_sort_idempotent(self):
        """
        sort is idempotent because sorting multiple times does not change the result.
        """
        hw1_02.word_freqs = [['test', 3], ['hello', 5], ['world', 2]]
        hw1_02.sort()
        expected = [['hello', 5], ['test', 3], ['world', 2]]
        self.assertEqual(expected, hw1_02.word_freqs)
        
        hw1_02.sort()
        self.assertEqual(expected, hw1_02.word_freqs)

if __name__ == '__main__':
    #python -m unittest hw1_02_test.py
    unittest.main()
