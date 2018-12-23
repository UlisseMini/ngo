// +build windows

package exec

import (
	"io"
	"os/exec"
)

// Spawn will spawn cmd over conn (full pty not supported yet on windows)
func Spawn(readwriter io.ReadWriter, cmd *exec.Cmd) error {
	cmd.Stdout = readwriter
	cmd.Stderr = readwriter
	cmd.Stdin = readwriter

	cmd.Run()

	return nil
}
