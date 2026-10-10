package proof

import "testing"

func TestParseContractValid(t *testing.T) {
	raw := []byte(`{
      "version": 1,
      "checks": [{"id":"unit","command":["go","test","-json","./..."],"parser":"go","timeout_seconds":120}],
      "requirements": [{"id":"r1","description":"d","check":"unit","test":"pkg::TestX"}]
    }`)
	c, err := ParseContract(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Version != 1 || len(c.Checks) != 1 || len(c.Requirements) != 1 {
		t.Fatalf("unexpected contract: %+v", c)
	}
}

func TestParseContractRejectsUnknownFields(t *testing.T) {
	raw := []byte(`{"version":1,"checks":[{"id":"u","command":["x"],"parser":"command","timeout_seconds":1}],"requirements":[{"id":"r","check":"u"}],"extra":true}`)
	if _, err := ParseContract(raw); err == nil {
		t.Fatal("expected error for unknown field")
	}
}

func TestParseContractRejectsTrailingContent(t *testing.T) {
	raw := []byte(`{"version":1,"checks":[{"id":"u","command":["x"],"parser":"command","timeout_seconds":1}],"requirements":[{"id":"r","check":"u"}]} extra`)
	if _, err := ParseContract(raw); err == nil {
		t.Fatal("expected error for trailing content")
	}
}

func TestContractValidation(t *testing.T) {
	cases := []struct {
		name string
		c    Contract
		ok   bool
	}{
		{"wrong version", Contract{Version: 2, Checks: []Check{{ID: "u", Command: []string{"x"}, Parser: "command", TimeoutSeconds: 1}}, Requirements: []Requirement{{ID: "r", Check: "u"}}}, false},
		{"no checks", Contract{Version: 1, Requirements: []Requirement{{ID: "r", Check: "u"}}}, false},
		{"no requirements", Contract{Version: 1, Checks: []Check{{ID: "u", Command: []string{"x"}, Parser: "command", TimeoutSeconds: 1}}}, false},
		{"duplicate check id", Contract{Version: 1, Checks: []Check{{ID: "u", Command: []string{"x"}, Parser: "command", TimeoutSeconds: 1}, {ID: "u", Command: []string{"y"}, Parser: "command", TimeoutSeconds: 1}}, Requirements: []Requirement{{ID: "r", Check: "u"}}}, false},
		{"empty command", Contract{Version: 1, Checks: []Check{{ID: "u", Command: []string{}, Parser: "command", TimeoutSeconds: 1}}, Requirements: []Requirement{{ID: "r", Check: "u"}}}, false},
		{"bad parser", Contract{Version: 1, Checks: []Check{{ID: "u", Command: []string{"x"}, Parser: "zzz", TimeoutSeconds: 1}}, Requirements: []Requirement{{ID: "r", Check: "u"}}}, false},
		{"timeout too large", Contract{Version: 1, Checks: []Check{{ID: "u", Command: []string{"x"}, Parser: "command", TimeoutSeconds: 4000}}, Requirements: []Requirement{{ID: "r", Check: "u"}}}, false},
		{"timeout zero", Contract{Version: 1, Checks: []Check{{ID: "u", Command: []string{"x"}, Parser: "command", TimeoutSeconds: 0}}, Requirements: []Requirement{{ID: "r", Check: "u"}}}, false},
		{"requirement unknown check", Contract{Version: 1, Checks: []Check{{ID: "u", Command: []string{"x"}, Parser: "command", TimeoutSeconds: 1}}, Requirements: []Requirement{{ID: "r", Check: "missing"}}}, false},
		{"command parser cannot name test", Contract{Version: 1, Checks: []Check{{ID: "u", Command: []string{"x"}, Parser: "command", TimeoutSeconds: 1}}, Requirements: []Requirement{{ID: "r", Check: "u", Test: "pkg::T"}}}, false},
		{"junit requires report path", Contract{Version: 1, Checks: []Check{{ID: "u", Command: []string{"x"}, Parser: "junit", TimeoutSeconds: 1}}, Requirements: []Requirement{{ID: "r", Check: "u"}}}, false},
		{"junit absolute report path", Contract{Version: 1, Checks: []Check{{ID: "u", Command: []string{"x"}, Parser: "junit", TimeoutSeconds: 1, ReportPath: "/etc/passwd"}}, Requirements: []Requirement{{ID: "r", Check: "u"}}}, false},
		{"junit traversal report path", Contract{Version: 1, Checks: []Check{{ID: "u", Command: []string{"x"}, Parser: "junit", TimeoutSeconds: 1, ReportPath: "../out.xml"}}, Requirements: []Requirement{{ID: "r", Check: "u"}}}, false},
		{"report path on non-junit", Contract{Version: 1, Checks: []Check{{ID: "u", Command: []string{"x"}, Parser: "command", TimeoutSeconds: 1, ReportPath: "out.xml"}}, Requirements: []Requirement{{ID: "r", Check: "u"}}}, false},
		{"valid", Contract{Version: 1, Checks: []Check{{ID: "u", Command: []string{"x"}, Parser: "junit", TimeoutSeconds: 10, ReportPath: "build/report.xml"}}, Requirements: []Requirement{{ID: "r", Check: "u", Test: "C::t"}}}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.c.Validate()
			if tc.ok && err != nil {
				t.Fatalf("expected valid, got %v", err)
			}
			if !tc.ok && err == nil {
				t.Fatal("expected invalid, got nil")
			}
		})
	}
}

func TestConfigHashChangesWithWhitespace(t *testing.T) {
	a := ConfigHash([]byte(`{"a":1}`))
	b := ConfigHash([]byte(`{"a":1} `))
	if a == b {
		t.Fatal("hash should change with whitespace")
	}
}
