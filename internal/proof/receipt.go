package proof

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

// ReceiptSchemaVersion is the version of the receipt document format.
const ReceiptSchemaVersion = 1

// ParseReceipt decodes a receipt from JSON bytes, rejecting unknown fields and
// trailing content so that a malformed receipt cannot be silently accepted.
func ParseReceipt(raw []byte) (*Receipt, error) {
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	var r Receipt
	if err := dec.Decode(&r); err != nil {
		return nil, fmt.Errorf("invalid receipt JSON: %w", err)
	}
	if _, err := dec.Token(); err != io.EOF {
		return nil, fmt.Errorf("invalid receipt: unexpected trailing content")
	}
	if r.ReceiptVersion != ReceiptSchemaVersion {
		return nil, fmt.Errorf("unsupported receipt version %d", r.ReceiptVersion)
	}
	return &r, nil
}

// Receipt is the portable verification record. It is unsigned and proves only
// the declared checks; it is not an authenticated attestation.
type Receipt struct {
	ReceiptVersion    int                 `json:"receipt_version"`
	AgentProofVersion string              `json:"agentproof_version"`
	Commit            string              `json:"commit"`
	ConfigHash        string              `json:"config_hash"`
	StartedAt         string              `json:"started_at"`
	Runner            string              `json:"runner"`
	Status            string              `json:"status"`
	Checks            []CheckResult       `json:"checks"`
	Requirements      []RequirementResult `json:"requirements"`
	Warnings          []string            `json:"warnings"`
}

// CheckResult is the recorded outcome of one check.
type CheckResult struct {
	ID         string       `json:"id"`
	Command    []string     `json:"command"`
	Parser     string       `json:"parser"`
	ExitCode   int          `json:"exit_code"`
	DurationMS int64        `json:"duration_ms"`
	Status     string       `json:"status"`
	Tests      []TestResult `json:"tests,omitempty"`
	Stdout     string       `json:"stdout,omitempty"`
	Stderr     string       `json:"stderr,omitempty"`
	Diagnostic string       `json:"diagnostic,omitempty"`
}

// RequirementResult is the recorded outcome of one requirement.
type RequirementResult struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Check       string `json:"check"`
	Test        string `json:"test,omitempty"`
	State       string `json:"state"`
	Detail      string `json:"detail,omitempty"`
}
