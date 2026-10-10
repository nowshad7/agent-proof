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
