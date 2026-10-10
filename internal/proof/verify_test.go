package proof

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// initRepo creates a git repository in a temp dir with the given files and a
// single commit, and returns the repo directory.
func initRepo(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	run := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(),
			"GIT_AUTHOR_NAME=test", "GIT_AUTHOR_EMAIL=test@example.com",
			"GIT_COMMITTER_NAME=test", "GIT_COMMITTER_EMAIL=test@example.com",
		)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	run("init", "-q")
	run("config", "user.email", "test@example.com")
	run("config", "user.name", "test")
	for name, content := range files {
		full := filepath.Join(dir, name)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	run("add", "-A")
	run("commit", "-q", "-m", "initial")
	return dir
}

// TestVerifySnapshot exercises the full flow against a real git archive,
// which includes git's pax_global_header metadata entry. This is the
// regression guard for safe tar metadata handling.
func TestVerifySnapshot(t *testing.T) {
	repo := initRepo(t, map[string]string{
		"agent-proof.json": `{
          "version": 1,
          "checks": [{"id":"present","command":["sh","-c","test -f marker.txt"],"parser":"command","timeout_seconds":10}],
          "requirements": [{"id":"file-present","description":"marker exists","check":"present"}]
        }`,
		"marker.txt": "hello\n",
	})

	r, err := Verify(context.Background(), Options{RepoDir: repo})
	if err != nil {
		t.Fatalf("verify failed: %v", err)
	}
	if r.Status != StatusPass {
		t.Fatalf("expected PASS, got %s (%+v)", r.Status, r.Requirements)
	}
	if len(r.Commit) != 40 {
		t.Fatalf("expected full SHA, got %q", r.Commit)
	}
}

// TestVerifyExcludesDirty confirms uncommitted changes do not affect the
// snapshot.
func TestVerifyExcludesDirty(t *testing.T) {
	repo := initRepo(t, map[string]string{
		"agent-proof.json": `{
          "version": 1,
          "checks": [{"id":"absent","command":["sh","-c","test ! -f dirty.txt"],"parser":"command","timeout_seconds":10}],
          "requirements": [{"id":"no-dirty","description":"dirty excluded","check":"absent"}]
        }`,
	})
	// Add an untracked working-tree file that must not appear in the snapshot.
	if err := os.WriteFile(filepath.Join(repo, "dirty.txt"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	r, err := Verify(context.Background(), Options{RepoDir: repo})
	if err != nil {
		t.Fatalf("verify failed: %v", err)
	}
	if r.Status != StatusPass {
		t.Fatalf("expected PASS (dirty excluded), got %s", r.Status)
	}
}

// TestVerifyFailClosedUnverified confirms a missing named test becomes
// UNVERIFIED and fails the gate.
func TestVerifyFailClosedUnverified(t *testing.T) {
	repo := initRepo(t, map[string]string{
		"agent-proof.json": `{
          "version": 1,
          "checks": [{"id":"j","command":["sh","-c","cat > report.xml <<'EOF'\n<testsuite name=\"s\"><testcase classname=\"C\" name=\"present\"/></testsuite>\nEOF"],"parser":"junit","timeout_seconds":10,"report_path":"report.xml"}],
          "requirements": [{"id":"has-missing","description":"needs a test that is not emitted","check":"j","test":"C::missing"}]
        }`,
	})
	r, err := Verify(context.Background(), Options{RepoDir: repo})
	if err != nil {
		t.Fatalf("verify failed: %v", err)
	}
	if r.Status != StatusFail {
		t.Fatalf("expected FAIL, got %s", r.Status)
	}
	if r.Requirements[0].State != StatusUnverified {
		t.Fatalf("expected UNVERIFIED, got %s", r.Requirements[0].State)
	}
}

// TestVerifyStaleReportRemoved confirms a committed stale report cannot pass
// when the command fails to regenerate it.
func TestVerifyStaleReportRemoved(t *testing.T) {
	repo := initRepo(t, map[string]string{
		"agent-proof.json": `{
          "version": 1,
          "checks": [{"id":"j","command":["sh","-c","exit 0"],"parser":"junit","timeout_seconds":10,"report_path":"report.xml"}],
          "requirements": [{"id":"r","description":"d","check":"j","test":"C::t"}]
        }`,
		// A stale, tracked report that claims a passing test.
		"report.xml": `<testsuite name="s"><testcase classname="C" name="t"/></testsuite>`,
	})
	r, err := Verify(context.Background(), Options{RepoDir: repo})
	if err != nil {
		t.Fatalf("verify failed: %v", err)
	}
	// The command does not regenerate the report, so reading it must error.
	if r.Status != StatusError {
		t.Fatalf("expected ERROR (stale report removed), got %s: %s", r.Status, r.Checks[0].Diagnostic)
	}
}

func TestVerifyGoAdapterEndToEnd(t *testing.T) {
	repo := initRepo(t, map[string]string{
		"go.mod": "module example.com/proj\n\ngo 1.23\n",
		"thing_test.go": `package proj

import "testing"

func TestUnicode(t *testing.T) {
	if "café" == "" {
		t.Fatal("unreachable")
	}
}
`,
		"agent-proof.json": `{
          "version": 1,
          "checks": [{"id":"unit","command":["go","test","-json","./..."],"parser":"go","timeout_seconds":120}],
          "requirements": [{"id":"unicode","description":"unicode preserved","check":"unit","test":"example.com/proj::TestUnicode"}]
        }`,
	})
	r, err := Verify(context.Background(), Options{RepoDir: repo})
	if err != nil {
		t.Fatalf("verify failed: %v", err)
	}
	if r.Status != StatusPass {
		t.Fatalf("expected PASS, got %s (%+v / %+v)", r.Status, r.Requirements, r.Checks[0].Diagnostic)
	}
}
