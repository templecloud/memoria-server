package logging

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

// ConfigureDefault configures the logging with the default values.
func ConfigureDefault() {
	Configure(NewDefaultConfig())
}

// Configure configures the logging.
func Configure(config *Config) {
	// Ensure Config
	if config == nil {
		config = NewDefaultConfig()
	}

	cfg := config.DefaultLogger
	// Handle default Logger type.
	if cfg.Type == JSONLogType {
		log.SetFormatter(&log.JSONFormatter{
			TimestampFormat: cfg.TimestampFormat,
		})
	} else {
		log.SetFormatter(&log.TextFormatter{
			TimestampFormat: cfg.TimestampFormat,
			FullTimestamp:   !cfg.DisableFullTimestamp,
			DisableColors:   cfg.DisableColors,
		})
	}

	// Handle default Logger type.
	log.SetReportCaller(cfg.EnableFunctionNameLogging) // Adds overhead - debug only.

	// Configure file logging.
	if cfg.LogFile != nil {
		fCfg := cfg.LogFile
		logFile := fmt.Sprintf("%s/%s", ensure(fCfg.Path), fCfg.Filename)
		f, err := os.OpenFile(logFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			fmt.Println(err)
			panic(err)
		} else {
			log.SetOutput(f)
		}
	}
}

// Safely ensure the 'log path directory' exists; else panic.
func ensure(path string) string {
	node, err := os.Stat(path)
	if os.IsNotExist(err) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			panic(err)
		}
	}
	if node.Mode().IsRegular() {
		panic(fmt.Sprintf("Cannot initialise log. File exists %s", path))
	}
	return path
}

