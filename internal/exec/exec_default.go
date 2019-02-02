// +build !windows

package exec

import (
	"io"
	"os"
	"os/exec"

	// Its very sad but i've gotta use this old ass library for pty's
	// Hopefully i can figure out the syscalls later and do it myself
	"github.com/kr/pty"
	"golang.org/x/crypto/ssh/terminal"
)

// Spawn will spawn cmd over readwriter (full pty not supported yet on windows)
func Spawn(readwriter io.ReadWriter, cmd *exec.Cmd) (err error) {
	f, err := pty.Start(cmd)
	if err != nil {
		return err
	}

	// Set stdin in raw mode.
	oldState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}

	// Remove raw mode after
	defer terminal.Restore(int(os.Stdin.Fd()), oldState)

	// Copy stdin to the pty and the pty to stdout
	done := make(chan error)
	go func() {
		_, err = io.Copy(f, readwriter)
		done <- err
	}()

	go func() {
		_, err = io.Copy(readwriter, f)
		done <- err
	}()

	// wait for one of them to finish then return possible error
	return <-done
}
