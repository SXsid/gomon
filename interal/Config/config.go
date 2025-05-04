package config

import "flag"

type Config struct {
	WatchDir string

	BuildCMD string

	RunCMD string

	IgnorePattern []string
}

//constructor

func NewConfig() *Config {

	config := &Config{
		WatchDir:      ".",
		BuildCMD:      "go build -o ./temp/gomonexe ./cmd/gomon",
		RunCMD:        "./temp/gomonexe",
		IgnorePattern: []string{"*/.git/*", "*/temp/*", "*/venodr/*"},
	}

	flag.StringVar(&config.WatchDir, "dir", config.WatchDir, "Directory to keep eye on")
	flag.StringVar(&config.BuildCMD, "buildCommand", config.BuildCMD, "command to build the applicaton")
	flag.StringVar(&config.RunCMD, "runCommand", config.RunCMD, "command to build the applicaton")

	return config
}
