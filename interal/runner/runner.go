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

type Runner struct {
	config  *config.Config
	started time.Time
	cmd     *exec.Cmd
}

func NewRunner(cfg *config.Config) *Runner {
	return &Runner{
		config: cfg,
	}
}

func (rn *Runner) Run() error {
	if rn.cmd != nil && rn.cmd.Process != nil {
		rn.Stop()

		// Wait for process to fully terminate before continuing
		if rn.config.IsWindows {
			time.Sleep(300 * time.Millisecond)
			if rn.config.Port != "" {

				if err := forcePortCleanup(rn.config.Port); err != nil {
					color.Yellow("Warning during port cleanup: %v", err)
				}
			}
		}
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

	if rn.config.IsWindows && rn.config.Port != "" {
		rn.cmd.Env = append(os.Environ(), rn.cmd.Env...)
	}

	if err := rn.cmd.Start(); err != nil {
		color.Red("Failed to start application: %v", err)
		return err
	}

	rn.started = time.Now()
	color.Green("Started application (PID: %d)", rn.cmd.Process.Pid)

	// Monitor the process in a goroutine
	go func() {
		err := rn.cmd.Wait()

		if rn.cmd != nil {
			if err != nil {
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
	pid := process.Pid

	port := ""
	if rn.config.IsWindows && rn.config.Port != "" {
		port = rn.config.Port
	}

	defer func() {
		rn.cmd = nil
	}()

	var err error

	if rn.config.IsWindows {
		err = terminateWindowProcess(pid)
	} else {
		err = process.Signal(syscall.SIGTERM)
	}

	if err != nil {

		color.Yellow("Process already terminated")

		if port != "" {
			forcePortCleanup(port)
		}
		return
	}

	gracePeriod := 500 * time.Millisecond
	timeout := time.After(gracePeriod)

	done := make(chan struct{})
	go func() {
		rn.cmd.Wait()
		close(done)
	}()

	select {
	case <-done:
		color.Yellow("Stopped application (PID: %d) gracefully after running for %v", pid, time.Since(rn.started))

		if port != "" {
			time.Sleep(100 * time.Millisecond)
			forcePortCleanup(port)
		}
	case <-timeout:

		var err error
		if rn.config.IsWindows {
			err = forceKillWindowProcess(pid)
		} else {
			err = process.Kill()
		}

		if err != nil {
			color.Red("Failed to kill process: %v", err)
		} else {
			color.Yellow("Forcibly killed application (PID: %d) after running for %v", pid, time.Since(rn.started))

			if port != "" {
				time.Sleep(300 * time.Millisecond)
				forcePortCleanup(port)
			}
		}
	}
}
