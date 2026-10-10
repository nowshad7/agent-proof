package proof

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"path"
	"strings"
)

// Parser identifies how a check's output is interpreted.
const (
	ParserCommand = "command"
	ParserGo      = "go"
	ParserJUnit   = "junit"
)

// Contract is the reviewable acceptance policy. It maps explicit requirements
// to the commands or exact test identities that supply their evidence.
type Contract struct {
	Version      int           `json:"version"`
	Checks       []Check       `json:"checks"`
	Requirements []Requirement `json:"requirements"`
}

// Check is a single executable step whose output becomes evidence.
type Check struct {
	ID             string   `json:"id"`
	Command        []string `json:"command"`
	Parser         string   `json:"parser"`
	TimeoutSeconds int      `json:"timeout_seconds"`
	// ReportPath is the relative path to a JUnit XML report produced by the
	// command. Required for the junit parser, forbidden otherwise.
	ReportPath string `json:"report_path,omitempty"`
}

// Requirement is an acceptance statement pointing at a check and, optionally,
// at an exact test identity within that check's evidence.
type Requirement struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Check       string `json:"check"`
	Test        string `json:"test,omitempty"`
}

// ConfigHash is the SHA-256 of the exact policy bytes, so whitespace changes
// also change the hash.
func ConfigHash(raw []byte) string {
	sum := sha256.Sum256(raw)
	return "sha256:" + hex.EncodeToString(sum[:])
}

// ParseContract strictly decodes a contract from raw JSON bytes. Unknown
// fields and trailing content are rejected. The returned contract is
// validated before it is returned.
func ParseContract(raw []byte) (*Contract, error) {
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()

	var c Contract
	if err := dec.Decode(&c); err != nil {
		return nil, fmt.Errorf("invalid contract JSON: %w", err)
	}
	// Reject trailing content after the first JSON value.
	if _, err := dec.Token(); err != io.EOF {
		return nil, fmt.Errorf("invalid contract: unexpected trailing content after JSON object")
	}
	if err := c.Validate(); err != nil {
		return nil, err
	}
	return &c, nil
}

// Validate enforces the contract rules documented in the configuration
// reference. A failing contract must stop verification before any command
// runs.
func (c *Contract) Validate() error {
	if c.Version != 1 {
		return fmt.Errorf("unsupported contract version %d: only version 1 is supported", c.Version)
	}
	if len(c.Checks) == 0 {
		return fmt.Errorf("contract must declare at least one check")
	}
	if len(c.Requirements) == 0 {
		return fmt.Errorf("contract must declare at least one requirement")
	}

	checkByID := make(map[string]*Check, len(c.Checks))
	for i := range c.Checks {
		ch := &c.Checks[i]
		if strings.TrimSpace(ch.ID) == "" {
			return fmt.Errorf("check %d: id must be nonempty", i)
		}
		if _, dup := checkByID[ch.ID]; dup {
			return fmt.Errorf("duplicate check id %q", ch.ID)
		}
		checkByID[ch.ID] = ch

		if len(ch.Command) == 0 {
			return fmt.Errorf("check %q: command must be a nonempty argument array", ch.ID)
		}
		for j, arg := range ch.Command {
			if arg == "" {
				return fmt.Errorf("check %q: command argument %d is empty", ch.ID, j)
			}
		}
		switch ch.Parser {
		case ParserCommand, ParserGo, ParserJUnit:
		default:
			return fmt.Errorf("check %q: parser must be one of %q, %q or %q", ch.ID, ParserCommand, ParserGo, ParserJUnit)
		}
		if ch.TimeoutSeconds < 1 || ch.TimeoutSeconds > 3600 {
			return fmt.Errorf("check %q: timeout_seconds must be between 1 and 3600", ch.ID)
		}
		if ch.Parser == ParserJUnit {
			if err := validateReportPath(ch.ReportPath); err != nil {
				return fmt.Errorf("check %q: %w", ch.ID, err)
			}
		} else if ch.ReportPath != "" {
			return fmt.Errorf("check %q: report_path is only valid for the %q parser", ch.ID, ParserJUnit)
		}
	}

	reqIDs := make(map[string]struct{}, len(c.Requirements))
	for i := range c.Requirements {
		r := &c.Requirements[i]
		if strings.TrimSpace(r.ID) == "" {
			return fmt.Errorf("requirement %d: id must be nonempty", i)
		}
		if _, dup := reqIDs[r.ID]; dup {
			return fmt.Errorf("duplicate requirement id %q", r.ID)
		}
		reqIDs[r.ID] = struct{}{}

		ch, ok := checkByID[r.Check]
		if !ok {
			return fmt.Errorf("requirement %q references unknown check %q", r.ID, r.Check)
		}
		// A command-only parser cannot satisfy a named-test requirement.
		if r.Test != "" && ch.Parser == ParserCommand {
			return fmt.Errorf("requirement %q names a test but check %q uses the %q parser, which has no test identities", r.ID, r.Check, ParserCommand)
		}
	}
	return nil
}

// validateReportPath enforces a safe, relative JUnit report path.
func validateReportPath(p string) error {
	if strings.TrimSpace(p) == "" {
		return fmt.Errorf("report_path is required for the junit parser")
	}
	if path.IsAbs(p) || strings.HasPrefix(p, "/") {
		return fmt.Errorf("report_path %q must be relative", p)
	}
	clean := path.Clean(p)
	if clean == ".." || strings.HasPrefix(clean, "../") || clean == "." {
		return fmt.Errorf("report_path %q must stay within the repository", p)
	}
	for _, seg := range strings.Split(clean, "/") {
		if seg == ".." {
			return fmt.Errorf("report_path %q must not traverse parent directories", p)
		}
	}
	return nil
}

// CheckByID returns the check with the given id, or nil.
func (c *Contract) CheckByID(id string) *Check {
	for i := range c.Checks {
		if c.Checks[i].ID == id {
			return &c.Checks[i]
		}
	}
	return nil
}
