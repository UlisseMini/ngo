// +build !windows

// exec is used for spawning processes with -e
// will spawn using a pty if it can
package exec

import (
	"io"
	"net"
	"os"
	"os/exec"

	// Its very sad but i've gotta use this old ass library for pty's
	// Hopefully i can figure out the syscalls later and do it myself
	"github.com/kr/pty"
	"golang.org/x/crypto/ssh/terminal"
)

// Spawn will spawn cmd over conn (full pty not supported yet on windows)
func Spawn(conn net.Conn, cmd *exec.Cmd) error {
	f, err := pty.Start(cmd)
	if err != nil {
		return err
	}

	defer f.Close()

	// Set stdin in raw mode.
	oldState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}

	// Remove raw mode after
	defer terminal.Restore(int(os.Stdin.Fd()), oldState)

	// Copy stdin to the pty and the pty to stdout.
	go io.Copy(f, conn)
	io.Copy(conn, f)

	return nil
}
