package export

import (
	"strings"
	"testing"
)

func TestUnicode(t *testing.T) {
	out, err := ToCSV([][]string{{"café", "日本語", "naïve"}})
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"café", "日本語", "naïve"} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected %q to be preserved in %q", want, out)
		}
	}
}

func TestQuoting(t *testing.T) {
	out, err := ToCSV([][]string{{"a,b", "c\nd"}})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, `"a,b"`) {
		t.Fatalf("expected comma field to be quoted: %q", out)
	}
}
