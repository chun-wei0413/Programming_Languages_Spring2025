import unittest
from tf_06 import filter_chars_and_normalize, remove_stop_words_curry

class TestTF06(unittest.TestCase):
    def test_filter_chars_and_normalize(self):
        sentence = "This is a test. This test is only a test."
        expected = "this is a test this test is only a test "

        result = filter_chars_and_normalize(sentence)

        self.assertEqual(expected, result)

    def test_currying(self):
        # 測試是否能正確執行柯里化的 remove_stop_words
        stop_words = ["this", "is", "a"]
        remove_stop_words = remove_stop_words_curry(stop_words)

        input_words = ["this", "is", "a", "test", "sentence"]
        expected_output = ["test", "sentence"]

        self.assertEqual(expected_output, remove_stop_words(input_words))

if __name__ == "__main__":
    unittest.main()