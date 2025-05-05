package config

type Config struct {
	WatchDir string

	BuildCMD string

	RunCMD string

	IgnorePattern []string
}

//constructor

func NewConfig() *Config {

	config := &Config{
		WatchDir: ".",
		IgnorePattern: []string{
			"temp", "temp/*", // Ignore the temp directory and all its contents
			".git/*", "node_modules/*", "vendor/*",
			"*.exe", "*.tmp", "*.log",
		},
		BuildCMD: "",
		RunCMD:   "",
	}

	return config
}
