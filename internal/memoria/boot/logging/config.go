package logging

import (
	"fmt"
	"os"
)

const (
	// TextLogType defines a TextFormatter logger.
	TextLogType = "text"
	// JSONLogType defines a JSONFormatter logger.
	JSONLogType = "json"
)

const (
	// DefaultTimestampFormat defines the default time format of the timestamps.
	DefaultTimestampFormat = "0001-01-01T00:00:00Z"
	// DefaultLogFilename defines the default name of the log file.
	DefaultLogFilename = "log.log"
)

//-------------------------------------------------------------------------------------------------
// Models

// Config denotes the configuration of the logging.
type Config struct {
	DefaultLogger *LogConfig `json:"defaultLogger" binding:"optional"`
}

// LogConfig denotes the configuration of the loggers.
type LogConfig struct {
	Type                      string         `json:"type" binding:"required"`
	TimestampFormat           string         `json:"host" binding:"optional"`
	DisableFullTimestamp      bool           `json:"disableFullTimestamp" binding:"optional"`
	DisableColors             bool           `json:"disableColors" binding:"optional"`
	EnableFunctionNameLogging bool           `json:"enabledFunctionNameLogging" binding:"optional"`
	LogFile                   *LogFileConfig `json:"logFile" binding:"optional"`
}

// LogFileConfig denotes the configuration of the log file.
type LogFileConfig struct {
	Path      string `json:"path" binding:"required"`
	Filename string `json:"filename" binding:"required"`
}

//-------------------------------------------------------------------------------------------------
// Private

// Return a default logging configuration.
func newDefaultLoggingConfig() *Config {
	return &Config{
		DefaultLogger: &LogConfig{
			Type:                      TextLogType,
			TimestampFormat:           DefaultTimestampFormat,
			DisableFullTimestamp:      false,
			DisableColors:             false,
			EnableFunctionNameLogging: true,
			LogFile: &LogFileConfig{
				Path:      cwd(),
				Filename: DefaultLogFilename,
			},
		},
	}
}

// Panic if the configuration is invalid.
func (c *Config) validate() {
	if c.DefaultLogger == nil {
		panic(fmt.Sprintf("No default logger was defined."))
	}
	if c.DefaultLogger.Type == "" {
		panic(fmt.Sprintf("The type of the default logger was not defined."))
	}
	if c.DefaultLogger.Type != TextLogType && c.DefaultLogger.Type != JSONLogType {
		panic(fmt.Sprintf("The type of the default logger '%s' was invalid.", c.DefaultLogger.Type))
	}
}

// Safely get the 'current working directory'; else panic.
func cwd() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return cwd
}
