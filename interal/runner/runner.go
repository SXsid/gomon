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

	//if we alredy hae a running process
	if rn.cmd != nil && rn.cmd.Process != nil {
		rn.Stop()
	}

	commands := strings.Fields(rn.config.RunCMD)
	if len(commands) == 0 {
		color.Red("Please enter a valid run command")
		return nil
	}
	//create the command
	rn.cmd = exec.Command(commands[0], commands[1:]...)
	rn.cmd.Stdout = os.Stdout
	rn.cmd.Stderr = os.Stderr
	if err := rn.cmd.Start(); err != nil {
		color.Red("Failed to start application: %v", err)
		return err
	}
	rn.started = time.Now()
	color.Green("Started application (PID: %d)", rn.cmd.Process.Pid)
	go func() {
		rn.cmd.Wait()
		//exec befroe stop completes
		if rn.cmd != nil {
			color.Yellow("Application exited after %v", time.Since(rn.started))
		}
	}()

	return nil

}

func (rn *Runner) Stop() {

	if rn.cmd == nil || rn.cmd.Process == nil {
		return
	}

	process := rn.cmd.Process

	//shutdown
	_ = process.Signal(syscall.SIGTERM)
	time.Sleep(100 * time.Millisecond)
	if process.Signal(syscall.Signal(0)) == nil {
		_ = process.Kill()
		 rn.cmd.Wait()
	}
	color.Yellow("Stopped application (PID: %d) after running for %v", process.Pid, time.Since(rn.started))
	rn.cmd = nil
}
