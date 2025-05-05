package runner

import (
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	config "github.com/SXsid/gomon/interal/Config"
	"github.com/fatih/color"
)

// Runner manages the execution of a configured command
type Runner struct {
	config  *config.Config
	started time.Time
	cmd     *exec.Cmd
}

// NewRunner creates a new runner instance
func NewRunner(cfg *config.Config) *Runner {
	return &Runner{
		config: cfg,
	}
}

// Run starts the configured command
func (rn *Runner) Run() error {
	// If we already have a running process, stop it first
	if rn.cmd != nil && rn.cmd.Process != nil {
		rn.Stop()
	}

	commands := strings.Fields(rn.config.RunCMD)
	if len(commands) == 0 {
		color.Red("Please enter a valid run command")
		return nil
	}

	// Create the command
	rn.cmd = exec.Command(commands[0], commands[1:]...)
	rn.cmd.Stdout = os.Stdout
	rn.cmd.Stderr = os.Stderr

	if err := rn.cmd.Start(); err != nil {
		color.Red("Failed to start application: %v", err)
		return err
	}

	rn.started = time.Now()
	color.Green("Started application (PID: %d)", rn.cmd.Process.Pid)

	// Monitor the process in a goroutine
	go func() {
		err := rn.cmd.Wait()
		// Process has exited; check if it was terminated by the runner or naturally
		if rn.cmd != nil {
			if err != nil {
				// Only log non-nil errors if they're not due to a signal we sent
				if exitErr, ok := err.(*exec.ExitError); !ok || exitErr.ExitCode() != -1 {
					color.Yellow("Application exited with error: %v after %v", err, time.Since(rn.started))
				} else {
					color.Yellow("Application terminated after %v", time.Since(rn.started))
				}
			} else {
				color.Yellow("Application exited normally after %v", time.Since(rn.started))
			}
		}
	}()

	return nil
}

// Stop gracefully terminates the running process
func (rn *Runner) Stop() {
	if rn.cmd == nil || rn.cmd.Process == nil {
		return
	}

	process := rn.cmd.Process
	defer func() {
		rn.cmd = nil
	}()

	// Try graceful shutdown with SIGTERM
	err := process.Signal(syscall.SIGTERM)
	if err != nil {
		// Process already gone
		color.Yellow("Process already terminated")
		return
	}

	// Give the process time to shut down gracefully
	gracePeriod := 500 * time.Millisecond
	timeout := time.After(gracePeriod)

	done := make(chan struct{})
	go func() {
		rn.cmd.Wait()
		close(done)
	}()

	// Wait for either the process to exit or the timeout
	select {
	case <-done:
		color.Yellow("Stopped application (PID: %d) gracefully after running for %v", process.Pid, time.Since(rn.started))
	case <-timeout:
		// Process didn't exit gracefully, force kill it
		if err := process.Kill(); err != nil {
			color.Red("Failed to kill process: %v", err)
		} else {
			color.Yellow("Forcibly killed application (PID: %d) after running for %v", process.Pid, time.Since(rn.started))
		}
	}
}
