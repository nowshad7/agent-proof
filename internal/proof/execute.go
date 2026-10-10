package proof

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// maxStreamBytes bounds captured stdout and stderr independently.
const maxStreamBytes = 1 << 20 // 1 MiB each

// ExecResult is the outcome of running one check command.
type ExecResult struct {
	ExitCode int
	Stdout   []byte
	Stderr   []byte
	Duration time.Duration
	// TimedOut is true when the timeout cancelled the command.
	TimedOut bool
	// Err is a non-acceptance operational failure: start failure or output
	// overflow. A nonzero exit is not reported here.
	Err error
}

// RunCheck executes a check's command inside dir with the configured timeout.
// It supplies a minimal, non-inherited environment with a fresh HOME,
// temporary directory and Go caches, while inheriting PATH so installed
// executables remain discoverable.
func RunCheck(ctx context.Context, dir string, ch *Check) *ExecResult {
	res := &ExecResult{ExitCode: -1}

	runCtx, cancel := context.WithTimeout(ctx, time.Duration(ch.TimeoutSeconds)*time.Second)
	defer cancel()

	cmd := exec.CommandContext(runCtx, ch.Command[0], ch.Command[1:]...)
	cmd.Dir = dir

	env, cleanup, err := buildEnv(dir)
	if err != nil {
		res.Err = err
		return res
	}
	defer cleanup()
	cmd.Env = env

	stdout := &cappedBuffer{limit: maxStreamBytes}
	stderr := &cappedBuffer{limit: maxStreamBytes}
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	configureProcGroup(cmd)

	start := time.Now()
	startErr := cmd.Start()
	if startErr != nil {
		res.Err = fmt.Errorf("starting command: %w", startErr)
		return res
	}

	waitErr := cmd.Wait()
	res.Duration = time.Since(start)
	res.Stdout = stdout.Bytes()
	res.Stderr = stderr.Bytes()

	// On timeout, make a best effort to cancel the whole process group.
	if runCtx.Err() == context.DeadlineExceeded {
		res.TimedOut = true
		killProcGroup(cmd)
		return res
	}

	if stdout.overflowed || stderr.overflowed {
		res.Err = fmt.Errorf("command output exceeded %d bytes", maxStreamBytes)
		return res
	}

	if waitErr != nil {
		if exitErr, ok := waitErr.(*exec.ExitError); ok {
			res.ExitCode = exitErr.ExitCode()
			return res
		}
		res.Err = fmt.Errorf("command failed: %w", waitErr)
		return res
	}
	res.ExitCode = 0
	return res
}

// buildEnv returns a minimal environment plus a cleanup for the scratch
// directories it creates.
func buildEnv(dir string) ([]string, func(), error) {
	home, err := os.MkdirTemp("", "agentproof-home-")
	if err != nil {
		return nil, func() {}, fmt.Errorf("creating scratch HOME: %w", err)
	}
	tmp, err := os.MkdirTemp("", "agentproof-tmp-")
	if err != nil {
		_ = os.RemoveAll(home)
		return nil, func() {}, fmt.Errorf("creating scratch TMPDIR: %w", err)
	}
	cleanup := func() {
		_ = os.RemoveAll(home)
		_ = os.RemoveAll(tmp)
	}

	env := []string{
		"HOME=" + home,
		"TMPDIR=" + tmp,
		"TEMP=" + tmp,
		"TMP=" + tmp,
		"GOCACHE=" + filepath.Join(home, "go-build"),
		"GOPATH=" + filepath.Join(home, "go"),
		"GOMODCACHE=" + filepath.Join(home, "go", "pkg", "mod"),
		"GOFLAGS=-mod=mod",
	}
	// Inherit PATH only, so installed toolchains remain discoverable without
	// forwarding the whole caller environment.
	if path := os.Getenv("PATH"); path != "" {
		env = append(env, "PATH="+path)
	}
	// Preserve a proxy configuration if the host requires one for module
	// downloads; this does not forward secrets.
	for _, k := range []string{"HTTPS_PROXY", "HTTP_PROXY", "NO_PROXY", "https_proxy", "http_proxy", "no_proxy", "SSL_CERT_FILE", "SSL_CERT_DIR"} {
		if v := os.Getenv(k); v != "" {
			env = append(env, k+"="+v)
		}
	}
	return env, cleanup, nil
}

// cappedBuffer records up to limit bytes and flags overflow.
type cappedBuffer struct {
	buf        bytes.Buffer
	limit      int
	overflowed bool
}

func (c *cappedBuffer) Write(p []byte) (int, error) {
	if c.overflowed {
		return len(p), nil
	}
	remaining := c.limit - c.buf.Len()
	if remaining <= 0 {
		c.overflowed = true
		return len(p), nil
	}
	if len(p) > remaining {
		c.buf.Write(p[:remaining])
		c.overflowed = true
		return len(p), nil
	}
	return c.buf.Write(p)
}

func (c *cappedBuffer) Bytes() []byte { return c.buf.Bytes() }
