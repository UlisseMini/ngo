// exec is used for spawning processes with -e
// will spawn using a pty if it can

package exec

import (
	"os/exec"
	"strings"
)

// Parse splits s into a list and creates an exec.Cmd using it.
func Parse(s string) *exec.Cmd {
	split := strings.Split(s, " ")
	return exec.Command(split[0], split[1:]...)
}
