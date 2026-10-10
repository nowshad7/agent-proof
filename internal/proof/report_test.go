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
