package proof

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// Version is the AgentProof build version. It is overridden at release time.
var Version = "0.1.0"

// Options configures a verification run.
type Options struct {
	RepoDir    string
	Revision   string
	ConfigPath string // optional path to a separately reviewed policy
	Runner     string
	Now        func() time.Time
}

// Verify resolves the commit, exports a snapshot, runs the contract's checks,
// parses evidence, aggregates requirements and returns a receipt. It does not
// write any files.
func Verify(ctx context.Context, opts Options) (*Receipt, error) {
	if opts.Now == nil {
		opts.Now = time.Now
	}
	repoDir := opts.RepoDir
	if repoDir == "" {
		repoDir = "."
	}

	commit, err := ResolveCommit(ctx, repoDir, opts.Revision)
	if err != nil {
		return nil, err
	}

	snapshot, err := ExtractSnapshot(ctx, repoDir, commit)
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(snapshot)

	// Load the contract: an explicit config path takes priority, otherwise
	// the committed snapshot's agent-proof.json.
	var rawConfig []byte
	if opts.ConfigPath != "" {
		rawConfig, err = os.ReadFile(opts.ConfigPath)
		if err != nil {
			return nil, fmt.Errorf("reading config %q: %w", opts.ConfigPath, err)
		}
	} else {
		configPath := filepath.Join(snapshot, "agent-proof.json")
		rawConfig, err = os.ReadFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("reading agent-proof.json from the committed snapshot: %w", err)
		}
	}

	contract, err := ParseContract(rawConfig)
	if err != nil {
		return nil, err
	}

	receipt := &Receipt{
		ReceiptVersion:    ReceiptSchemaVersion,
		AgentProofVersion: Version,
		Commit:            commit,
		ConfigHash:        ConfigHash(rawConfig),
		StartedAt:         opts.Now().UTC().Format(time.RFC3339),
		Runner:            opts.Runner,
		Warnings: []string{
			"Trusted local execution is not a sandbox: checks run with the caller's privileges.",
			"This receipt is unsigned and editable; it is not an authenticated attestation.",
			"AgentProof verifies the declared evidence workflow, not the honesty or completeness of the tests.",
		},
	}

	// Run each check and collect evidence.
	evidence := make(map[string]*Evidence, len(contract.Checks))
	for i := range contract.Checks {
		ch := &contract.Checks[i]
		cr, ev := runAndParse(ctx, snapshot, ch)
		receipt.Checks = append(receipt.Checks, cr)
		if ev != nil {
			evidence[ch.ID] = ev
		}
	}

	// Aggregate requirements.
	checkStatus := make(map[string]string, len(receipt.Checks))
	for _, cr := range receipt.Checks {
		checkStatus[cr.ID] = cr.Status
	}
	for i := range contract.Requirements {
		r := &contract.Requirements[i]
		receipt.Requirements = append(receipt.Requirements, evaluateRequirement(r, checkStatus[r.Check], evidence[r.Check]))
	}

	receipt.Status = overallStatus(receipt)
	return receipt, nil
}

// runAndParse executes a check, parses evidence where applicable and returns a
// CheckResult plus any parsed evidence.
func runAndParse(ctx context.Context, snapshot string, ch *Check) (CheckResult, *Evidence) {
	cr := CheckResult{
		ID:       ch.ID,
		Command:  ch.Command,
		Parser:   ch.Parser,
		ExitCode: -1,
	}

	// For junit, remove any stale tracked report before running so that old
	// XML cannot be mistaken for fresh evidence.
	var reportTarget string
	if ch.Parser == ParserJUnit {
		var perr error
		reportTarget, perr = safeJoin(mustAbs(snapshot), ch.ReportPath)
		if perr != nil {
			cr.Status = StatusError
			cr.Diagnostic = perr.Error()
			return cr, nil
		}
		if err := removeStaleReport(reportTarget); err != nil {
			cr.Status = StatusError
			cr.Diagnostic = err.Error()
			return cr, nil
		}
	}

	exec := RunCheck(ctx, snapshot, ch)
	cr.ExitCode = exec.ExitCode
	cr.DurationMS = exec.Duration.Milliseconds()
	cr.Stdout = string(exec.Stdout)
	cr.Stderr = string(exec.Stderr)

	if exec.TimedOut {
		cr.Status = StatusError
		cr.Diagnostic = fmt.Sprintf("command timed out after %d seconds", ch.TimeoutSeconds)
		return cr, nil
	}
	if exec.Err != nil {
		cr.Status = StatusError
		cr.Diagnostic = exec.Err.Error()
		return cr, nil
	}

	switch ch.Parser {
	case ParserCommand:
		if exec.ExitCode == 0 {
			cr.Status = StatusPass
		} else {
			cr.Status = StatusFail
		}
		return cr, nil

	case ParserGo:
		ev, err := ParseGoTestJSON(bytes.NewReader(exec.Stdout))
		if err != nil {
			cr.Status = StatusError
			cr.Diagnostic = err.Error()
			return cr, nil
		}
		cr.Tests = sortedTests(ev)
		cr.Status = checkStatusFromEvidence(exec.ExitCode, ev)
		return cr, ev

	case ParserJUnit:
		data, err := os.ReadFile(reportTarget)
		if err != nil {
			cr.Status = StatusError
			cr.Diagnostic = fmt.Sprintf("reading junit report %q: %v", ch.ReportPath, err)
			return cr, nil
		}
		if err := rejectSymlink(reportTarget); err != nil {
			cr.Status = StatusError
			cr.Diagnostic = err.Error()
			return cr, nil
		}
		ev, err := ParseJUnitXML(bytes.NewReader(data))
		if err != nil {
			cr.Status = StatusError
			cr.Diagnostic = err.Error()
			return cr, nil
		}
		cr.Tests = sortedTests(ev)
		cr.Status = checkStatusFromEvidence(exec.ExitCode, ev)
		return cr, ev
	}

	cr.Status = StatusError
	cr.Diagnostic = "unknown parser"
	return cr, nil
}

// checkStatusFromEvidence derives a check's status from its exit code and
// parsed tests. A nonzero exit, or any failing test, fails the check.
func checkStatusFromEvidence(exitCode int, ev *Evidence) string {
	if exitCode != 0 {
		return StatusFail
	}
	for _, st := range ev.Tests {
		if st == TestFail {
			return StatusFail
		}
	}
	return StatusPass
}

// evaluateRequirement computes a requirement's resulting state.
func evaluateRequirement(r *Requirement, checkStatus string, ev *Evidence) RequirementResult {
	res := RequirementResult{
		ID:          r.ID,
		Description: r.Description,
		Check:       r.Check,
		Test:        r.Test,
	}

	if checkStatus == StatusError {
		res.State = StatusError
		res.Detail = "the referenced check reported an operational error"
		return res
	}

	if r.Test == "" {
		// Command-level requirement: tied directly to the check status.
		if checkStatus == StatusPass {
			res.State = StatusPass
		} else {
			res.State = StatusFail
			res.Detail = "the referenced check did not pass"
		}
		return res
	}

	// Named-test requirement.
	if ev == nil {
		res.State = StatusUnverified
		res.Detail = "no evidence was available for the referenced check"
		return res
	}
	st, ok := ev.Tests[r.Test]
	if !ok {
		res.State = StatusUnverified
		res.Detail = fmt.Sprintf("named test %q was not found in the evidence", r.Test)
		return res
	}
	switch st {
	case TestSkip:
		res.State = StatusSkipped
		res.Detail = "the named test was reported as skipped"
	case TestFail:
		res.State = StatusFail
		res.Detail = "the named test failed"
	case TestPass:
		if checkStatus == StatusFail {
			// A named test cannot rescue a failing parent check.
			res.State = StatusFail
			res.Detail = "the named test passed but its parent check failed"
		} else {
			res.State = StatusPass
		}
	}
	return res
}

// overallStatus folds every check and requirement into the gate. An error
// takes precedence; otherwise any non-passing check or requirement fails.
func overallStatus(r *Receipt) string {
	for _, cr := range r.Checks {
		if cr.Status == StatusError {
			return StatusError
		}
	}
	for _, rr := range r.Requirements {
		if rr.State == StatusError {
			return StatusError
		}
	}
	for _, cr := range r.Checks {
		if cr.Status != StatusPass {
			return StatusFail
		}
	}
	for _, rr := range r.Requirements {
		if rr.State != StatusPass {
			return StatusFail
		}
	}
	return StatusPass
}

func sortedTests(ev *Evidence) []TestResult {
	out := make([]TestResult, 0, len(ev.Tests))
	for name, st := range ev.Tests {
		out = append(out, TestResult{Name: name, Status: string(st)})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

func removeStaleReport(path string) error {
	err := os.Remove(path)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("removing stale report %q: %w", path, err)
	}
	return nil
}

func rejectSymlink(path string) error {
	info, err := os.Lstat(path)
	if err != nil {
		return fmt.Errorf("inspecting report %q: %w", path, err)
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("report %q is a symlink, which is rejected", path)
	}
	return nil
}

func mustAbs(p string) string {
	abs, err := filepath.Abs(p)
	if err != nil {
		return p
	}
	return abs
}
