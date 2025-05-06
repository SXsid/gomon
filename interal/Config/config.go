package config

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/fatih/color"

	"github.com/SXsid/gomon/interal/util"
)

type Config struct {
	WatchDir string

	BuildCMD string

	RunCMD string

	IgnorePattern []string
	Port          string
	IsWindows     bool
}

//constructor

func NewConfig() (*Config, error) {

	config := &Config{
		WatchDir: ".",
		IgnorePattern: []string{
			"temp", "temp/*",
			".git", ".git/*",
			"node_modules", "node_modules/*",
			"vendor", "vendor/*",
			"*.exe", "*.tmp", "*.log",
		},

		BuildCMD:  "",
		RunCMD:    "",
		Port:      "",
		IsWindows: false,
	}

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `

Usage:
  gomon --main <path-to-main.go> [--watch <directory>] [--port <port>]

Flags:
`)
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, `Example:
  gomon --main ./main.go --port 8080`)
	}
	watchDir := flag.String("watch", config.WatchDir, "Directory to watch for changes (default: current directory)")
	fileLocation := flag.String("main", "", "Path to the main Go file (e.g., ./cmd/server/main.go)")
	port := flag.String("port", config.Port, "Port number if the app runs an HTTP server (e.g., 8080)")

	flag.Parse()
	if *watchDir != config.WatchDir {
		config.WatchDir = *watchDir
	}

	if *fileLocation == "" {
		fmt.Println("Please provide the path to the main Go file using --main")
		flag.Usage()
		return nil, fmt.Errorf("no file path was provided")
	}
	config.IsWindows = runtime.GOOS == "windows"

	if config.IsWindows {
		config.BuildCMD = fmt.Sprintf("go build -mod=mod -o .\\temp\\gomon.exe %s", *fileLocation)

		config.RunCMD = ".\\temp\\gomon.exe"
	} else {

		config.BuildCMD = fmt.Sprintf("go build -mod=mod -o ./temp/gomon.exe %s", *fileLocation)

		config.RunCMD = "./temp/gomon.exe"
	}

	if *port != "" {
		config.Port = *port
	}
	if util.IsAWebserver(*fileLocation) && config.Port == "" {
		color.Yellow("YOUR ARE LIKELY RUNNING AN WEBSERVER WITHOUT POROVIDING THE PORT")

	}

	return config, nil
}
