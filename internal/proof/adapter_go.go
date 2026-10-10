package proof

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

// goEvent is one record from `go test -json` output.
type goEvent struct {
	Action  string `json:"Action"`
	Package string `json:"Package"`
	Test    string `json:"Test"`
}

// ParseGoTestJSON reads a complete `go test -json` stream and returns the
// evidence it contains. It fails closed: malformed lines, duplicate test
// identities, incomplete test or package events, and zero results are all
// surfaced so that missing evidence cannot be mistaken for success.
//
// Test identities use the form "package::TestName" and subtest paths are
// preserved (for example "pkg::TestThing/case").
func ParseGoTestJSON(r io.Reader) (*Evidence, error) {
	type testKey struct {
		pkg  string
		test string
	}
	// terminal records the final outcome seen for a started test.
	terminal := make(map[testKey]TestStatus)
	// started records tests that emitted a run action.
	started := make(map[testKey]bool)
	// seenTerminal guards against duplicate terminal events for one identity.
	seenTerminal := make(map[testKey]bool)
	// pkgStarted / pkgDone track package-level completion.
	pkgStarted := make(map[string]bool)
	pkgDone := make(map[string]bool)

	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 0, 64*1024), 8*1024*1024)
	var sawLine bool
	for sc.Scan() {
		line := bytes.TrimSpace(sc.Bytes())
		if len(line) == 0 {
			continue
		}
		sawLine = true
		var ev goEvent
		// go test may add fields over time; only reject lines that are not
		// valid JSON objects at all.
		if err := json.Unmarshal(line, &ev); err != nil {
			return nil, fmt.Errorf("malformed go test json line: %w", err)
		}
		if ev.Package == "" {
			return nil, fmt.Errorf("malformed go test json: event without package")
		}

		if ev.Test == "" {
			// Package-level event.
			switch ev.Action {
			case "start", "run":
				pkgStarted[ev.Package] = true
			case "pass", "fail", "skip":
				pkgStarted[ev.Package] = true
				pkgDone[ev.Package] = true
			}
			continue
		}

		key := testKey{pkg: ev.Package, test: ev.Test}
		switch ev.Action {
		case "run":
			started[key] = true
		case "pass", "fail", "skip":
			started[key] = true
			if seenTerminal[key] {
				return nil, fmt.Errorf("duplicate go test identity %q: ambiguous repeated result", identity(ev.Package, ev.Test))
			}
			seenTerminal[key] = true
			switch ev.Action {
			case "pass":
				terminal[key] = TestPass
			case "fail":
				terminal[key] = TestFail
			case "skip":
				terminal[key] = TestSkip
			}
		}
	}
	if err := sc.Err(); err != nil {
		if err == bufio.ErrTooLong {
			return nil, fmt.Errorf("go test json line exceeded buffer limit")
		}
		return nil, fmt.Errorf("reading go test json: %w", err)
	}
	if !sawLine {
		return nil, fmt.Errorf("go test produced no json events")
	}

	// Every started test must have a terminal event.
	for key := range started {
		if _, ok := terminal[key]; !ok {
			return nil, fmt.Errorf("incomplete go test json: test %q started but never finished", identity(key.pkg, key.test))
		}
	}
	// Every package that emitted any test must have completed.
	pkgsWithTests := make(map[string]bool)
	for key := range started {
		pkgsWithTests[key.pkg] = true
	}
	for pkg := range pkgsWithTests {
		if !pkgDone[pkg] {
			return nil, fmt.Errorf("incomplete go test json: package %q never reported a final result", pkg)
		}
	}

	tests := make(map[string]TestStatus, len(terminal))
	for key, st := range terminal {
		tests[identity(key.pkg, key.test)] = st
	}
	if len(tests) == 0 {
		return nil, fmt.Errorf("go test json contained no test results")
	}
	return &Evidence{Tests: tests}, nil
}

func identity(pkg, test string) string {
	return pkg + "::" + test
}
