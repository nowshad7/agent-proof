#!/usr/bin/env python3
"""Run the unittest suite and emit a JUnit XML report.

This is a tiny, dependency-free runner so the example produces real JUnit
output that AgentProof can consume. It discovers tests in the current
directory and writes ``report.xml``.

The ``classname::name`` identities it emits (for example
``test_slugify.SlugifyTests::test_unicode``) are what the contract references.
"""

import sys
import time
import unittest
from xml.sax.saxutils import escape, quoteattr


def main() -> int:
    loader = unittest.TestLoader()
    discovered = loader.discover(start_dir=".", pattern="test_*.py")

    # Snapshot the test objects before running: TestSuite.run() nulls out its
    # references afterwards to free memory, so the list must be captured first.
    tests = list(iter_tests(discovered))
    suite = unittest.TestSuite(tests)

    result = unittest.TestResult()
    start = time.time()
    suite.run(result)
    duration = time.time() - start

    # Index failures/errors/skips by the test object for lookup.
    failed = {t: msg for t, msg in result.failures}
    errored = {t: msg for t, msg in result.errors}
    skipped = {t: msg for t, msg in result.skipped}

    cases = []
    for test in tests:
        classname = f"{test.__class__.__module__}.{test.__class__.__name__}"
        name = test._testMethodName
        if test in errored:
            body = f"<error message={quoteattr('error')}>{escape(errored[test])}</error>"
        elif test in failed:
            body = f"<failure message={quoteattr('assertion failed')}>{escape(failed[test])}</failure>"
        elif test in skipped:
            body = f"<skipped message={quoteattr(skipped[test])}/>"
        else:
            body = ""
        cases.append(
            f'  <testcase classname={quoteattr(classname)} name={quoteattr(name)}>{body}</testcase>'
        )

    total = result.testsRun
    n_fail = len(result.failures)
    n_err = len(result.errors)
    n_skip = len(result.skipped)

    xml = (
        '<?xml version="1.0" encoding="utf-8"?>\n'
        f'<testsuites tests="{total}" failures="{n_fail}" errors="{n_err}" skipped="{n_skip}">\n'
        f'  <testsuite name="slugify" tests="{total}" failures="{n_fail}" '
        f'errors="{n_err}" skipped="{n_skip}" time="{duration:.3f}">\n'
        + "\n".join(cases)
        + "\n  </testsuite>\n</testsuites>\n"
    )

    with open("report.xml", "w", encoding="utf-8") as fh:
        fh.write(xml)

    # Mirror the suite outcome in the process exit code.
    return 0 if result.wasSuccessful() else 1


def iter_tests(suite):
    for item in suite:
        if isinstance(item, unittest.TestSuite):
            yield from iter_tests(item)
        else:
            yield item


if __name__ == "__main__":
    sys.exit(main())
