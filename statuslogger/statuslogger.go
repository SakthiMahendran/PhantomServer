package statuslogger

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	ct "github.com/daviddengcn/go-colortext"
)

type StatusLogger struct {
	logger      *log.Logger
	loggerMutex sync.Mutex
}

func NewStatusLogger() StatusLogger {
	sl := StatusLogger{}

	logFileName := getLogFileName()

	logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		panic(err)
	}

	lw := newLogWriter(logFile)

	sl.logger = log.New(&lw, "", log.Ldate|log.Ltime)
	sl.loggerMutex = sync.Mutex{}

	sl.LogInfo("Logging Started.")
	sl.LogInfo("Logs Will be Written to StdOut and ", logFileName)

	return sl
}

func getLogFileName() string {
	t := time.Now()

	day := t.Day()
	month := int(t.Month())
	year := t.Year()

	fileName := fmt.Sprint(day, "_", month, "_", year, "_log.txt")

	return fileName
}

func (sl *StatusLogger) LogErr(a ...any) {
	logMsg := fmt.Sprint(a...)

	sl.loggerMutex.Lock()

	ct.Foreground(ct.Red, false)

	sl.logger.SetPrefix("ERROR: ")
	sl.logger.Println(logMsg)

	ct.ResetColor()

	sl.loggerMutex.Unlock()
}

func (sl *StatusLogger) LogInfo(a ...any) {
	logMsg := fmt.Sprint(a...)

	sl.loggerMutex.Lock()

	ct.Foreground(ct.Green, false)

	sl.logger.SetPrefix("INFO: ")
	sl.logger.Println(logMsg)

	ct.ResetColor()

	sl.loggerMutex.Unlock()
}

func (sl *StatusLogger) NewLine() {
	sl.loggerMutex.Lock()

	sl.logger.SetPrefix("")
	sl.logger.SetFlags(0)
	sl.logger.Println()
	sl.logger.SetFlags(log.Ldate | log.Ltime)

	sl.loggerMutex.Unlock()
}
