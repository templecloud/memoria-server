package logging

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

// ConfigureDefaultLogger initialises the global logger.
func ConfigureDefaultLogger() {

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		// TimestampFormat: "02-01-2006 15:04:05",
		DisableColors: false,
	})
	
	// log.SetFormatter(&log.JSONFormatter{
	// 	TimestampFormat: "02-01-2006 15:04:05",
	// })

	// log.SetReportCaller(true) // Adds overhead - debug only.

	logFile := "logfile.log"
	f, err := os.OpenFile(logFile, os.O_WRONLY | os.O_APPEND | os.O_CREATE, 0644)
    if err != nil {
        fmt.Println(err)
    }else{
        log.SetOutput(f)
    }
}
