package pdfinject

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
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

func (s *ShellCommand) RunInPathWithStd(dir string,  stdInputFile string, stdOutputFile string, args ...string) error {

	// Create the command.
	var stderr bytes.Buffer
	cmd := exec.Command(s.cmdName, args...)
	cmd.Stderr = &stderr
	cmd.Dir = dir


	inputBytes, err := ioutil.ReadFile(stdInputFile)
	if err != nil {
		return err
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	// Start the command and wait for it to exit.
	err = cmd.Start()
	if err != nil {
		return fmt.Errorf(strings.TrimSpace(stderr.String()))
	}

	_, err = io.WriteString(stdin, string(inputBytes))
	if err != nil {
		return err
	}

	outFile, err := os.Create(stdOutputFile)
	// handle err
	defer outFile.Close()
	_, err = io.Copy(outFile, stdout)
	if err != nil {
		return err
	}

	return nil
}
