package proof

// Status values classify checks, requirements and the overall gate. Missing
// results are never silently treated as success.
const (
	// StatusPass means the declared command or test supplied passing evidence.
	StatusPass = "PASS"
	// StatusFail means a command returned nonzero or a test failed.
	StatusFail = "FAIL"
	// StatusUnverified means there was no usable test evidence, or a named
	// test was missing.
	StatusUnverified = "UNVERIFIED"
	// StatusError means execution, timeout, parsing or report handling failed.
	StatusError = "ERROR"
	// StatusSkipped means a required test was reported as skipped.
	StatusSkipped = "SKIPPED"
)

// TestStatus is the per-test outcome drawn from parsed evidence.
type TestStatus string

const (
	TestPass TestStatus = "pass"
	TestFail TestStatus = "fail"
	TestSkip TestStatus = "skip"
)

// Evidence is the parsed result of one check's output: a set of test
// identities mapped to their outcomes.
type Evidence struct {
	// Tests maps a test identity to its outcome. For the command parser this
	// is always empty.
	Tests map[string]TestStatus
}

// TestResult is a single test identity and outcome, used for stable,
// serialised reporting.
type TestResult struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}
