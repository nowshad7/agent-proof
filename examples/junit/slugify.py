"""A small, deterministic slug helper used to demonstrate JUnit evidence."""

import re
import unicodedata


def slugify(text: str) -> str:
    """Return a lowercase, hyphenated ASCII slug for text.

    Unicode letters are transliterated to their closest ASCII form; runs of
    non-alphanumeric characters collapse to a single hyphen.
    """
    normalized = unicodedata.normalize("NFKD", text)
    ascii_text = normalized.encode("ascii", "ignore").decode("ascii")
    ascii_text = ascii_text.lower()
    ascii_text = re.sub(r"[^a-z0-9]+", "-", ascii_text)
    return ascii_text.strip("-")
