package proof

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

// maxJUnitBytes bounds JUnit report input.
const maxJUnitBytes = 8 << 20 // 8 MiB

// junitSuites models both a top-level <testsuites> wrapper and a single
// <testsuite>; nested suites are supported.
type junitSuites struct {
	XMLName xml.Name     `xml:"testsuites"`
	Suites  []junitSuite `xml:"testsuite"`
}

type junitSuite struct {
	XMLName xml.Name     `xml:"testsuite"`
	Suites  []junitSuite `xml:"testsuite"`
	Cases   []junitCase  `xml:"testcase"`
}

type junitCase struct {
	Name      string        `xml:"name,attr"`
	ClassName string        `xml:"classname,attr"`
	Failures  []junitDetail `xml:"failure"`
	Errors    []junitDetail `xml:"error"`
	Skipped   []junitDetail `xml:"skipped"`
}

type junitDetail struct {
	Message string `xml:"message,attr"`
}

// ParseJUnitXML reads a JUnit XML report and returns the evidence it
// contains. It fails closed: malformed XML, missing case names and duplicate
// identities are errors. Identities use "classname::name"; an absent
// classname yields "::name".
func ParseJUnitXML(r io.Reader) (*Evidence, error) {
	limited := io.LimitReader(r, maxJUnitBytes+1)
	data, err := io.ReadAll(limited)
	if err != nil {
		return nil, fmt.Errorf("reading junit report: %w", err)
	}
	if len(data) > maxJUnitBytes {
		return nil, fmt.Errorf("junit report exceeds %d bytes", maxJUnitBytes)
	}
	if len(strings.TrimSpace(string(data))) == 0 {
		return nil, fmt.Errorf("junit report is empty")
	}

	var cases []junitCase

	// Try <testsuites> first, then a bare <testsuite>.
	var suites junitSuites
	if err := xml.Unmarshal(data, &suites); err == nil && suites.XMLName.Local == "testsuites" {
		for _, s := range suites.Suites {
			collectCases(s, &cases)
		}
	} else {
		var single junitSuite
		if err := xml.Unmarshal(data, &single); err != nil {
			return nil, fmt.Errorf("malformed junit xml: %w", err)
		}
		if single.XMLName.Local != "testsuite" {
			return nil, fmt.Errorf("malformed junit xml: expected <testsuites> or <testsuite> root, got <%s>", single.XMLName.Local)
		}
		collectCases(single, &cases)
	}

	tests := make(map[string]TestStatus, len(cases))
	for _, c := range cases {
		name := strings.TrimSpace(c.Name)
		if name == "" {
			return nil, fmt.Errorf("junit report contains a testcase without a name")
		}
		id := strings.TrimSpace(c.ClassName) + "::" + name
		if _, dup := tests[id]; dup {
			return nil, fmt.Errorf("duplicate junit test identity %q", id)
		}
		tests[id] = junitCaseStatus(c)
	}
	if len(tests) == 0 {
		return nil, fmt.Errorf("junit report contained no test cases")
	}
	return &Evidence{Tests: tests}, nil
}

func collectCases(s junitSuite, out *[]junitCase) {
	*out = append(*out, s.Cases...)
	for _, child := range s.Suites {
		collectCases(child, out)
	}
}

func junitCaseStatus(c junitCase) TestStatus {
	// Failure and error take precedence over skipped.
	if len(c.Failures) > 0 || len(c.Errors) > 0 {
		return TestFail
	}
	if len(c.Skipped) > 0 {
		return TestSkip
	}
	return TestPass
}
