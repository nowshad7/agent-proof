//go:build unix

package proof

import (
	"os/exec"
	"syscall"
)

// configureProcGroup places the child in its own process group so that a
// timeout can cancel the whole tree.
func configureProcGroup(cmd *exec.Cmd) {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	cmd.SysProcAttr.Setpgid = true
}

// killProcGroup sends SIGKILL to the child's process group.
func killProcGroup(cmd *exec.Cmd) {
	if cmd.Process == nil {
		return
	}
	// Negative pid targets the process group.
	_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
}
