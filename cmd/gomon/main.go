package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	config "github.com/SXsid/gomon/interal/Config"
	"github.com/SXsid/gomon/interal/builder"
	"github.com/SXsid/gomon/interal/runner"
	"github.com/SXsid/gomon/interal/watcher"
	"github.com/fatih/color"
)

func main() {
	//r
	// fomon <paht of main.go file>
	args := os.Args
	if len(args) != 2 || args[1] == "--help" || args[1] == "-h" {
		fmt.Print(`
gomon - Hot reload for Go applications

Usage:
  gomon <path-to-main.go>

Example:
  gomon ./main.go

This will watch your Go files, rebuild the binary, and restart on changes.
`)
		return
	}
	printBanner()
	cfg := config.NewConfig()
	if runtime.GOOS == "windows" {
		cfg.BuildCMD = fmt.Sprintf("go build -o .\\temp\\gomonexe.exe %s", args[1])
		cfg.RunCMD = ".\\temp\\gomonexe.exe"
	} else {
		cfg.BuildCMD = fmt.Sprintf("go build -o ./temp/gomonexe %s", args[1])
		cfg.RunCMD = "./temp/gomonexe"
	}

	appRunner := runner.NewRunner(cfg)

	appBuilder := builder.NewBuilder(cfg)
	eventChan := make(chan struct{})
	signalCahn := make(chan os.Signal, 1)
	signal.Notify(signalCahn, syscall.SIGINT, syscall.SIGTERM)
	watcher, err := watcher.NewWatcher(eventChan, cfg)
	if err != nil {
		log.Fatalf("Failed to create file watcher: %v", err)
	}
	if err := watcher.Start(); err != nil {
		log.Fatalf("Failed to start file watcher: %v", err)
	}
	color.Yellow("Performing initial build...")
	//first build
	buildResult := appBuilder.Build()
	startServer(buildResult, appRunner)

	for {
		select {
		//file has changed
		case <-eventChan:
			color.Yellow("\n🔄 File changes detected, rebuilding...")
			buildResult := appBuilder.Build()
			startServer(buildResult, appRunner)

		case sig := <-signalCahn:
			color.Yellow("\n⛔ Received signal %v, shutting down...", sig)
			appRunner.Stop()
			watcher.DoneChannel <- struct{}{}
			return

		}

	}

}

func startServer(buildResult builder.BuildResult, appRunner *runner.Runner) {
	if buildResult.Success {
		appRunner.Stop()
		if err := appRunner.Run(); err != nil {
			log.Printf("Failed to restart application: %v", err)
		}
	} else {
		color.Red("Build failed:")
		fmt.Println(buildResult.Output)
	}
}

func printBanner() {
	banner := `
  ____       __  __              
 / ___| ___ |  \/  | ___  _ __  
| |  _ / _ \| |\/| |/ _ \| '_ \ 
| |_| | (_) | |  | | (_) | | | |
 \____|\___/|_|  |_|\___/|_| |_|
                                           
`
	color.Cyan(banner)
	color.Cyan("Hot reload for Go applications")
	color.Cyan("----------------------------------")
}
