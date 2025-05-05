package builder

import (
	"os"
	"os/exec"
	"strings"
	"time"

	config "github.com/SXsid/gomon/interal/Config"
	"github.com/fatih/color"
)

type Builder struct {
	config *config.Config
}

type BuildResult struct {
	Success bool
	Error   error
	Output  string
}

func NewBuilder(Cfg *config.Config) *Builder {
	return &Builder{
		config: Cfg,
	}
}

func (bd *Builder) Build() BuildResult {
	startTimer := time.Now()
	if strings.Contains(bd.config.BuildCMD, "./temp/") {
		os.MkdirAll("./temp", 0755)
	}
	commandSlice := strings.Fields(bd.config.BuildCMD)
	if len(commandSlice) == 0 {
		return res(false, nil, "please enter a valid build command")
	}
	cmd := exec.Command(commandSlice[0], commandSlice[1:]...)
	//exe the command
	output, err := cmd.CombinedOutput()
	buildTime := time.Since(startTimer)
	outputStr := string(output)
	if err != nil {
		color.Red("Build failed in %v", buildTime)
		return res(false, err, outputStr)
	}

	color.Green("builded succeded in %v", buildTime)
	return res(true, nil, outputStr)
}

func res(success bool, err error, output string) BuildResult {
	return BuildResult{
		Success: success,
		Error:   err,
		Output:  output,
	}
}
