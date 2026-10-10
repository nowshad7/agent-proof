package proof

import (
	"context"
	"testing"
)

func TestRunCheckExitCodes(t *testing.T) {
	dir := t.TempDir()
	ctx := context.Background()

	pass := RunCheck(ctx, dir, &Check{ID: "p", Command: []string{"sh", "-c", "exit 0"}, Parser: ParserCommand, TimeoutSeconds: 10})
	if pass.Err != nil || pass.ExitCode != 0 {
		t.Fatalf("pass: err=%v exit=%d", pass.Err, pass.ExitCode)
	}

	fail := RunCheck(ctx, dir, &Check{ID: "f", Command: []string{"sh", "-c", "exit 3"}, Parser: ParserCommand, TimeoutSeconds: 10})
	if fail.Err != nil || fail.ExitCode != 3 {
		t.Fatalf("fail: err=%v exit=%d", fail.Err, fail.ExitCode)
	}
}

func TestRunCheckStartFailure(t *testing.T) {
	res := RunCheck(context.Background(), t.TempDir(), &Check{ID: "x", Command: []string{"this-binary-does-not-exist-agentproof"}, Parser: ParserCommand, TimeoutSeconds: 10})
	if res.Err == nil {
		t.Fatal("expected start failure error")
	}
}

func TestRunCheckTimeout(t *testing.T) {
	res := RunCheck(context.Background(), t.TempDir(), &Check{ID: "t", Command: []string{"sh", "-c", "sleep 5"}, Parser: ParserCommand, TimeoutSeconds: 1})
	if !res.TimedOut {
		t.Fatalf("expected timeout, got %+v", res)
	}
}

func TestRunCheckOutputOverflow(t *testing.T) {
	// Emit more than 1 MiB on stdout.
	res := RunCheck(context.Background(), t.TempDir(), &Check{ID: "o", Command: []string{"sh", "-c", "head -c 2000000 /dev/zero | tr '\\0' 'a'"}, Parser: ParserCommand, TimeoutSeconds: 20})
	if res.Err == nil {
		t.Fatalf("expected overflow error, got exit=%d", res.ExitCode)
	}
}
