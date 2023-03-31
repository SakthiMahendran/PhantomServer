package statuslogger

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	ct "github.com/daviddengcn/go-colortext"
)

// StatusLogger is a type representing a logger that can log messages to the console and a file.
type StatusLogger struct {
	logger      *log.Logger // The logger that will be used to log messages to the console and file
	loggerMutex sync.Mutex  // Mutex to ensure thread-safe access to the logger
}

// NewStatusLogger creates a new StatusLogger instance.
func NewStatusLogger() StatusLogger {
	sl := StatusLogger{}

	// Get a filename for the log file using the current date.
	logFileName := getLogFileName()

	// Open the log file, creating it if it doesn't exist, and append to it if it does.
	logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		panic(err)
	}

	// Create a new log writer to write to the console and file.
	lw := newLogWriter(logFile)

	// Create a new logger that writes to the log writer.
	sl.logger = log.New(&lw, "", log.Ldate|log.Ltime)
	sl.loggerMutex = sync.Mutex{}

	// Log a message indicating that logging has started.
	sl.LogInfo("Logging Started.")
	// Log a message indicating that logs will be written to the console and file.
	sl.LogInfo("Logs Will be Written to StdOut and ", logFileName)

	return sl
}

// getLogFileName generates a log file name using the current date.
func getLogFileName() string {
	t := time.Now()

	day := t.Day()
	month := int(t.Month())
	year := t.Year()

	fileName := fmt.Sprint(day, "_", month, "_", year, "_log.txt")

	return fileName
}

// LogErr logs an error message to the console and file.
func (sl *StatusLogger) LogErr(a ...any) {
	logMsg := fmt.Sprint(a...)

	sl.loggerMutex.Lock()

	// Set the console text color to red to indicate an error.
	ct.Foreground(ct.Red, false)

	// Set the logger prefix to "ERROR: " and log the message.
	sl.logger.SetPrefix("ERROR: ")
	sl.logger.Println(logMsg)

	// Reset the console text color.
	ct.ResetColor()

	sl.loggerMutex.Unlock()
}

// LogInfo logs an informational message to the console and file.
func (sl *StatusLogger) LogInfo(a ...any) {
	logMsg := fmt.Sprint(a...)

	sl.loggerMutex.Lock()

	// Set the console text color to green to indicate an informational message.
	ct.Foreground(ct.Green, false)

	// Set the logger prefix to "INFO: " and log the message.
	sl.logger.SetPrefix("INFO: ")
	sl.logger.Println(logMsg)

	// Reset the console text color.
	ct.ResetColor()

	sl.loggerMutex.Unlock()
}

// NewLine logs a new line to the console and file.
func (sl *StatusLogger) NewLine() {
	sl.loggerMutex.Lock()

	// Clear the logger prefix and flags, log a new line, and reset the flags.
	sl.logger.SetPrefix("")
	sl.logger.SetFlags(0)
	sl.logger.Println()
	sl.logger.SetFlags(log.Ldate | log.Ltime)

	sl.loggerMutex.Unlock()
}
