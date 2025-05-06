package runner

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func terminateWindowProcess(pid int) error {
	cmd := exec.Command("taskkill", "/PID", strconv.Itoa(pid))

	return cmd.Run()

}

func forceKillWindowProcess(pid int) error {
	cmd := exec.Command("taskkill", "/F", "PID", strconv.Itoa(pid))
	return cmd.Run()
}

//free the port forcefully

func forcePortCleanup(rn *Runner) error {
	//find the processs using our port

	cmd := exec.Command("cmd", "/C", fmt.Sprintf("netsat -ano |findstr :%s | findstr LISTENING", rn.config.Port))

	output, err := cmd.CombinedOutput()
	if err != nil {
		if len(output) == 0 {
			return nil
		}
		return fmt.Errorf("command failed: %v", err)

	}
	processes := strings.Split(string(output), "\n")
	for _, process := range processes {
		values := strings.Fields(process)
		if len(values) >= 5 {
			pid := values[4]
			pidInt, err := strconv.Atoi(pid)
			if err != nil {
				continue
			}
			forceKillWindowProcess(pidInt)
		}
	}

	return nil

}
