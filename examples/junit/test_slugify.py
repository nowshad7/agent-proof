import unittest

from slugify import slugify


class SlugifyTests(unittest.TestCase):
    def test_basic(self):
        self.assertEqual(slugify("Hello World"), "hello-world")

    def test_unicode(self):
        self.assertEqual(slugify("Café au lait"), "cafe-au-lait")

    def test_collapses_separators(self):
        self.assertEqual(slugify("a  --  b"), "a-b")


if __name__ == "__main__":
    unittest.main()
