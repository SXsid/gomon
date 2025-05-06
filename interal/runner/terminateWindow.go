package runner

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

func terminateWindowProcess(pid int) error {
	cmd := exec.Command("taskkill", "/PID", strconv.Itoa(pid))
	return cmd.Run()
}

func forceKillWindowProcess(pid int) error {
	cmd := exec.Command("taskkill", "/F", "/PID", strconv.Itoa(pid))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("taskkill failed for PID %d: %v\nOutput: %s", pid, err, output)
	}
	return nil
}

// Improved port cleanup function that takes a port string directly
func forcePortCleanup(port string) error {
	color.Yellow("Cleaning up port %s...", port)

	// First try to find processes using the port
	findCmd := exec.Command("cmd", "/C", fmt.Sprintf("netstat -ano | findstr :%s | findstr LISTENING", port))
	output, err := findCmd.CombinedOutput()

	if err == nil && len(output) > 0 {
		// Found processes using the port, terminate them
		lines := strings.Split(string(output), "\r\n")
		for _, line := range lines {
			if line == "" {
				continue
			}

			fields := strings.Fields(line)
			if len(fields) >= 5 {
				pidStr := fields[4]
				pid, err := strconv.Atoi(pidStr)
				if err != nil {
					color.Yellow("Invalid PID: %s", pidStr)
					continue
				}

				color.Yellow("Killing process %d using port %s", pid, port)
				killCmd := exec.Command("taskkill", "/F", "/PID", pidStr)
				if err := killCmd.Run(); err != nil {
					color.Red("Failed to kill PID %d: %v", pid, err)
				} else {
					color.Green("Successfully killed process with PID %d", pid)
				}
			}
		}
	}

	portNum, err := strconv.Atoi(port)
	if err != nil {
		return fmt.Errorf("invalid port number: %s", port)
	}

	commands := []string{
		// Delete connections for specific local port
		fmt.Sprintf("netsh int ipv4 delete tcpconnection localport=%d", portNum),
		// Alternative syntax sometimes needed
		fmt.Sprintf("netsh int ipv4 delete tcpconnection 0.0.0.0:%d", portNum),
		fmt.Sprintf("netsh int ipv4 delete tcpconnection 127.0.0.1:%d", portNum),
		// Try IPv6 as well
		fmt.Sprintf("netsh int ipv6 delete tcpconnection localport=%d", portNum),
	}

	for _, cmdStr := range commands {
		cmd := exec.Command("cmd", "/C", cmdStr)
		cmd.Run()
	}

	checkCmd := exec.Command("cmd", "/C", fmt.Sprintf("netstat -ano | findstr :%s | findstr LISTENING", port))
	checkOutput, _ := checkCmd.CombinedOutput()

	if len(checkOutput) > 0 {
		color.Yellow("Warning: Port %s might still be in use after cleanup attempts", port)
	} else {
		color.Green("Port %s successfully freed", port)
	}

	return nil
}
