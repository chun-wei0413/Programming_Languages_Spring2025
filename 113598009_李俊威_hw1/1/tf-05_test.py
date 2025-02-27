import unittest
# from tf_05 import filter_chars_and_normalize
import tf_05

class TestTF05(unittest.TestCase):

    def setUp(self):
        tf_05.data = []
        tf_05.words = []
        tf_05.word_freqs = []

    def test_filter_chars_and_normalize(self):
        sentence = "This is a test. This test is only a test."
        expected = "this is a test  this test is only a test "
        tf_05.data = list(sentence)
        tf_05.filter_chars_and_normalize()
        self.assertEqual(expected, ''.join(tf_05.data))

    def test_scan(self):
        sentence = "this is a test  this test is only a test "
        expected = ['this', 'is', 'a', 'test', 'this', 'test', 'is', 'only', 'a', 'test']
        tf_05.data = list(sentence)
        tf_05.scan()
        self.assertEqual(expected, tf_05.words)

    def test_remove_stop_words(self):
        tf_05.words = ['this', 'is', 'a', 'test', 'this', 'test', 'is', 'only', 'a', 'test']
        tf_05.remove_stop_words()
        self.assertEqual(['test', 'test', 'test'], tf_05.words)

    def test_frequencies(self):
        tf_05.words = ['test', 'test', 'test']
        tf_05.frequencies()
        self.assertEqual([['test', 3]], tf_05.word_freqs)

if __name__ == '__main__':
    unittest.main()