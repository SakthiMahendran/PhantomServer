package statuslogger

import "os"

type logWriter struct {
	logFile *os.File
}

func newLogWriter(file *os.File) logWriter {
	lw := logWriter{}
	lw.logFile = file

	return lw
}

func (lw *logWriter) Write(p []byte) (n int, err error) {
	n1, err1 := os.Stdout.Write(p)

	if err != nil {
		return n1, err1
	}

	n2, err2 := lw.logFile.Write(p)

	if err != nil {
		return n2, err2
	}

	return n1, nil
}
