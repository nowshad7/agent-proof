//go:build !unix

package proof

import "os/exec"

// configureProcGroup is a no-op on platforms without Unix process groups.
// Descendant-process cancellation is weaker here; this limitation is
// documented in the execution model.
func configureProcGroup(cmd *exec.Cmd) {}

// killProcGroup relies on the context cancellation already applied by
// exec.CommandContext, which kills the direct child process only.
func killProcGroup(cmd *exec.Cmd) {}
