package proof

import (
	"strings"
	"testing"
)

func TestParseJUnitSuites(t *testing.T) {
	in := `<?xml version="1.0"?>
<testsuites>
  <testsuite name="outer">
    <testcase classname="pkg.Mod" name="test_pass"/>
    <testcase classname="pkg.Mod" name="test_fail"><failure message="boom"/></testcase>
    <testcase classname="pkg.Mod" name="test_err"><error message="x"/></testcase>
    <testcase classname="pkg.Mod" name="test_skip"><skipped/></testcase>
    <testsuite name="nested">
      <testcase name="test_noclass"/>
    </testsuite>
  </testsuite>
</testsuites>`
	ev, err := ParseJUnitXML(strings.NewReader(in))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := map[string]TestStatus{
		"pkg.Mod::test_pass": TestPass,
		"pkg.Mod::test_fail": TestFail,
		"pkg.Mod::test_err":  TestFail,
		"pkg.Mod::test_skip": TestSkip,
		"::test_noclass":     TestPass,
	}
	if len(ev.Tests) != len(want) {
		t.Fatalf("got %d, want %d: %+v", len(ev.Tests), len(want), ev.Tests)
	}
	for k, v := range want {
		if ev.Tests[k] != v {
			t.Errorf("%s: got %q want %q", k, ev.Tests[k], v)
		}
	}
}

func TestParseJUnitBareSuite(t *testing.T) {
	in := `<testsuite name="s"><testcase classname="C" name="t"/></testsuite>`
	ev, err := ParseJUnitXML(strings.NewReader(in))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ev.Tests["C::t"] != TestPass {
		t.Fatalf("unexpected: %+v", ev.Tests)
	}
}

func TestParseJUnitMalformed(t *testing.T) {
	if _, err := ParseJUnitXML(strings.NewReader("<testsuite><bad")); err == nil {
		t.Fatal("expected error for malformed xml")
	}
}

func TestParseJUnitEmpty(t *testing.T) {
	if _, err := ParseJUnitXML(strings.NewReader("   ")); err == nil {
		t.Fatal("expected error for empty report")
	}
}

func TestParseJUnitMissingName(t *testing.T) {
	in := `<testsuite><testcase classname="C"/></testsuite>`
	if _, err := ParseJUnitXML(strings.NewReader(in)); err == nil {
		t.Fatal("expected error for missing name")
	}
}

func TestParseJUnitDuplicateIdentity(t *testing.T) {
	in := `<testsuite><testcase classname="C" name="t"/><testcase classname="C" name="t"/></testsuite>`
	if _, err := ParseJUnitXML(strings.NewReader(in)); err == nil {
		t.Fatal("expected error for duplicate identity")
	}
}

func TestParseJUnitNoCases(t *testing.T) {
	in := `<testsuite name="s"></testsuite>`
	if _, err := ParseJUnitXML(strings.NewReader(in)); err == nil {
		t.Fatal("expected error for no test cases")
	}
}
