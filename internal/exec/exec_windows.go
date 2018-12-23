// +build windows

package exec

// Spawn will spawn cmd over conn (full pty not supported yet on windows)
func Spawn(conn net.Conn, cmd *exec.Cmd) error {
	cmd.Stdout = conn
	cmd.Stderr = conn
	cmd.Stdin = conn

	cmd.Run()

	return nil
}
