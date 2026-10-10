package proof

import (
	"strings"
	"testing"
)

func TestParseGoTestJSONValid(t *testing.T) {
	in := strings.Join([]string{
		`{"Action":"run","Package":"ex/p","Test":"TestA"}`,
		`{"Action":"output","Package":"ex/p","Test":"TestA","Output":"ok"}`,
		`{"Action":"pass","Package":"ex/p","Test":"TestA"}`,
		`{"Action":"run","Package":"ex/p","Test":"TestB"}`,
		`{"Action":"run","Package":"ex/p","Test":"TestB/sub"}`,
		`{"Action":"pass","Package":"ex/p","Test":"TestB/sub"}`,
		`{"Action":"fail","Package":"ex/p","Test":"TestB"}`,
		`{"Action":"skip","Package":"ex/p","Test":"TestC"}`,
		`{"Action":"fail","Package":"ex/p"}`,
	}, "\n")
	ev, err := ParseGoTestJSON(strings.NewReader(in))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := map[string]TestStatus{
		"ex/p::TestA":     TestPass,
		"ex/p::TestB":     TestFail,
		"ex/p::TestB/sub": TestPass,
		"ex/p::TestC":     TestSkip,
	}
	if len(ev.Tests) != len(want) {
		t.Fatalf("got %d tests, want %d: %+v", len(ev.Tests), len(want), ev.Tests)
	}
	for k, v := range want {
		if ev.Tests[k] != v {
			t.Errorf("%s: got %q want %q", k, ev.Tests[k], v)
		}
	}
}

func TestParseGoTestJSONMalformed(t *testing.T) {
	if _, err := ParseGoTestJSON(strings.NewReader("not json")); err == nil {
		t.Fatal("expected error for malformed line")
	}
}

func TestParseGoTestJSONEmpty(t *testing.T) {
	if _, err := ParseGoTestJSON(strings.NewReader("   \n")); err == nil {
		t.Fatal("expected error for empty stream")
	}
}

func TestParseGoTestJSONIncompleteTest(t *testing.T) {
	in := strings.Join([]string{
		`{"Action":"run","Package":"ex/p","Test":"TestA"}`,
		`{"Action":"pass","Package":"ex/p"}`,
	}, "\n")
	if _, err := ParseGoTestJSON(strings.NewReader(in)); err == nil {
		t.Fatal("expected error for test without terminal event")
	}
}

func TestParseGoTestJSONIncompletePackage(t *testing.T) {
	in := strings.Join([]string{
		`{"Action":"run","Package":"ex/p","Test":"TestA"}`,
		`{"Action":"pass","Package":"ex/p","Test":"TestA"}`,
	}, "\n")
	if _, err := ParseGoTestJSON(strings.NewReader(in)); err == nil {
		t.Fatal("expected error for package without terminal event")
	}
}

func TestParseGoTestJSONDuplicateIdentity(t *testing.T) {
	in := strings.Join([]string{
		`{"Action":"run","Package":"ex/p","Test":"TestA"}`,
		`{"Action":"pass","Package":"ex/p","Test":"TestA"}`,
		`{"Action":"pass","Package":"ex/p","Test":"TestA"}`,
		`{"Action":"pass","Package":"ex/p"}`,
	}, "\n")
	if _, err := ParseGoTestJSON(strings.NewReader(in)); err == nil {
		t.Fatal("expected error for duplicate identity")
	}
}

func TestParseGoTestJSONNoTests(t *testing.T) {
	in := `{"Action":"pass","Package":"ex/p"}`
	if _, err := ParseGoTestJSON(strings.NewReader(in)); err == nil {
		t.Fatal("expected error for zero test results")
	}
}
