import unittest
from tf_06 import filter_chars_and_normalize

class TestTF06(unittest.TestCase):
    def test_filter_chars_and_normalize(self):
        sentence = "This is a test. This test is only a test."
        expected = "this is a test this test is only a test "

        result = filter_chars_and_normalize(sentence)

        self.assertEqual(expected, result)

    def test_currying(self):
        self.assertEqual(5, f(1, 2, 3))

        self.assertEqual(5, cf(1)(2)(3))

def f(x, y, z):
    return x * y + z

def cf(x):
    def g(y):
        def h(z):
            return x * y + z
        return h
    return g

if __name__ == "__main__":
    unittest.main()