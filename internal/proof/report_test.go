package proof

import (
	"strings"
	"testing"
)

func sampleReceipt() *Receipt {
	return &Receipt{
		ReceiptVersion:    ReceiptSchemaVersion,
		AgentProofVersion: "0.1.0-test",
		Commit:            "0123456789012345678901234567890123456789",
		ConfigHash:        "sha256:abc",
		StartedAt:         "2026-10-09T00:00:00Z",
		Status:            StatusFail,
		Checks: []CheckResult{{
			ID: "unit", Command: []string{"go", "test"}, Parser: "go",
			Status: StatusFail, ExitCode: 1,
			Tests:  []TestResult{{Name: "p::T<script>", Status: "fail"}},
			Stderr: "<b>danger</b>",
		}},
		Requirements: []RequirementResult{{
			ID: "r<x>", Description: "d", Check: "unit", Test: "p::T<script>", State: StatusFail,
		}},
		Warnings: []string{"unsigned"},
	}
}

func TestRenderHTMLEscapes(t *testing.T) {
	out := string(RenderHTML(sampleReceipt()))
	if strings.Contains(out, "<script>") {
		t.Fatal("unescaped script tag in HTML report")
	}
	if strings.Contains(out, "<b>danger</b>") {
		t.Fatal("unescaped stderr in HTML report")
	}
	if !strings.Contains(out, "Content-Security-Policy") {
		t.Fatal("missing CSP")
	}
}

func TestRenderJSONRoundTrip(t *testing.T) {
	b, err := RenderJSON(sampleReceipt())
	if err != nil {
		t.Fatal(err)
	}
	r, err := ParseReceipt(b)
	if err != nil {
		t.Fatalf("round trip failed: %v", err)
	}
	if r.Status != StatusFail || len(r.Checks) != 1 {
		t.Fatalf("unexpected: %+v", r)
	}
}

func TestRenderMarkdownContainsStatus(t *testing.T) {
	out := string(RenderMarkdown(sampleReceipt()))
	if !strings.Contains(out, "FAIL") || !strings.Contains(out, "## Requirements") {
		t.Fatal("markdown missing expected content")
	}
}

func TestRenderGitHubComment(t *testing.T) {
	out := string(RenderGitHubComment(sampleReceipt()))
	if !strings.HasPrefix(out, GitHubCommentMarker) {
		t.Fatal("gh comment must start with the sticky marker")
	}
	if !strings.Contains(out, "AgentProof — FAIL") {
		t.Fatalf("gh comment missing status header: %s", out)
	}
	// Commit is shortened to 12 chars.
	if !strings.Contains(out, "commit 012345678901`") {
		t.Fatalf("gh comment should show short commit: %s", out)
	}
	// Table-breaking characters in evidence must be escaped, not raw.
	if strings.Contains(out, "p::T<script>|") {
		t.Fatal("unescaped pipe in gh comment table")
	}
	if !strings.Contains(out, "<details><summary>Checks (1)</summary>") {
		t.Fatal("gh comment missing collapsible checks section")
	}
}

func TestRenderBadgePass(t *testing.T) {
	r := sampleReceipt()
	r.Status = StatusPass
	r.Requirements = []RequirementResult{
		{ID: "a", State: StatusPass},
		{ID: "b", State: StatusPass},
	}
	out, err := RenderBadge(r)
	if err != nil {
		t.Fatal(err)
	}
	s := string(out)
	if !strings.Contains(s, `"schemaVersion":1`) {
		t.Fatalf("badge missing shields schema: %s", s)
	}
	if !strings.Contains(s, `"label":"agent-proof"`) {
		t.Fatalf("badge missing label: %s", s)
	}
	if !strings.Contains(s, `"message":"2/2 passed"`) {
		t.Fatalf("badge message wrong: %s", s)
	}
	if !strings.Contains(s, `"color":"brightgreen"`) {
		t.Fatalf("passing badge must be green: %s", s)
	}
}

// TestRenderBadgeFailsClosed is the critical property: no status other than an
// explicit PASS may ever be coloured green, including an unknown or empty one.
// A green badge is a public claim of verification; it must never appear unless
// the gate actually passed.
func TestRenderBadgeFailsClosed(t *testing.T) {
	for _, status := range []string{
		StatusFail, StatusUnverified, StatusError, StatusSkipped,
		"", "WAT", "pass", "Pass", "PASSED", "OK",
	} {
		r := sampleReceipt()
		r.Status = status
		r.Requirements = []RequirementResult{{ID: "a", State: status}}
		out, err := RenderBadge(r)
		if err != nil {
			t.Fatalf("status %q: %v", status, err)
		}
		s := string(out)
		for _, green := range []string{"brightgreen", "green", "success"} {
			if strings.Contains(s, `"color":"`+green+`"`) {
				t.Fatalf("status %q produced a green badge (fail-closed violated): %s", status, s)
			}
		}
	}
}

func TestRenderBadgeErrorAndEmpty(t *testing.T) {
	// An ERROR status reports the status word, not a misleading count.
	r := sampleReceipt()
	r.Status = StatusError
	r.Requirements = []RequirementResult{{ID: "a", State: StatusPass}}
	out, _ := RenderBadge(r)
	if !strings.Contains(string(out), `"message":"error"`) {
		t.Fatalf("error badge should report error, got: %s", out)
	}
	if !strings.Contains(string(out), `"color":"red"`) {
		t.Fatalf("error badge must be red, got: %s", out)
	}

	// No requirements and an empty status must read "unknown", never a pass.
	r2 := &Receipt{ReceiptVersion: ReceiptSchemaVersion}
	out2, _ := RenderBadge(r2)
	if !strings.Contains(string(out2), `"message":"unknown"`) {
		t.Fatalf("empty receipt badge should be unknown, got: %s", out2)
	}
}

func TestStatusEmoji(t *testing.T) {
	cases := map[string]string{
		StatusPass:       "✅",
		StatusFail:       "❌",
		StatusError:      "❌",
		StatusUnverified: "⚠️",
		StatusSkipped:    "⚠️",
	}
	for status, want := range cases {
		if got := statusEmoji(status); got != want {
			t.Errorf("statusEmoji(%s) = %q, want %q", status, got, want)
		}
	}
}
