package utils

import (
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

type ProgramCloseCallback func()

// ListenForProgramClose sets a listener for program exits so we can exit gracefully.
func ListenForProgramClose(cb func()) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cb()
		os.Exit(0)
	}()
}

// RunCommandWithSpinner runs a command with a spinner instead
// of streaming the output.
func RunCommandWithSpinner(command string, args []string, terminalSpinner TerminalSpinner) error {
	// Start the spinner.
	terminalSpinner.Create()

	// Run the command.
	err := exec.Command(command, args...).Run()
	if err != nil {
		terminalSpinner.Fail()
		return err
	}

	// Stop the spinner and return nil.
	terminalSpinner.Stop()
	return nil
}

// RunCommand runs a command and return the output.
func RunCommand(command string, args []string) (string, error) {
	output, err := exec.Command(command, args...).Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

// RunCommandWithOutput runs a command and streams the output to
// the terminal.
func RunCommandWithOutput(command string, args []string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	return err
}
