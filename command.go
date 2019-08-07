package pdfinject

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type ShellCommand struct {
	cmdName string
}

func NewShellCommand(name string) ShellCommand {
	return ShellCommand{
		name,
	}
}

// runCommandInPath runs a command and waits for it to exit.
// The working directory is also set.
// The stderr error message is returned on error.
func (s *ShellCommand) RunInPath(dir string, args ...string) error {
	// Create the command.
	var stderr bytes.Buffer
	cmd := exec.Command(s.cmdName, args...)
	cmd.Stderr = &stderr
	cmd.Dir = dir

	// Start the command and wait for it to exit.
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf(strings.TrimSpace(stderr.String()))
	}

	return nil
}
