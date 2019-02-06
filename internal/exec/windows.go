// +build windows

package exec

import (
	"fmt"
	"io"
	"os/exec"
	"syscall"
)

// Exec execeutes a command over a stream (usually a net.Conn),
// It will execute cmd inside a pty if possible.
func Exec(cmd *exec.Cmd, conn io.ReadWriter) error {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}

	cmd.Stdout = conn
	cmd.Stderr = conn
	cmd.Stdin = conn

	if err := cmd.Start(); err != nil {
		fmt.Fprintf(conn, "error starting process: %v\n", err)
		return err
	}

	if _, err := cmd.Process.Wait(); err != nil {
		fmt.Fprintf(conn, "error waiting for process: %v\n", err)
		return err
	}

	return nil
}
