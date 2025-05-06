package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
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

	printBanner()
	cfg, err := config.NewConfig()

	if err != nil {
		return
	}
	fmt.Println(cfg.Port)
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
			close(watcher.DoneChannel)
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
