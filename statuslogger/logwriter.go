package statuslogger

import "os"

type logWriter struct {
	logFile *os.File // holds the file to which logs will be written
}

// creates a new logWriter and returns it
func newLogWriter(file *os.File) logWriter {
	lw := logWriter{}
	lw.logFile = file

	return lw
}

// Writes the given bytes to standard output and the log file
// Returns the number of bytes written and any error encountered
func (lw *logWriter) Write(p []byte) (n int, err error) {
	n1, err1 := os.Stdout.Write(p) // write the bytes to standard output and save the number of bytes written and any error

	// if an error was encountered, return the number of bytes written and the error
	if err != nil {
		return n1, err1
	}

	// write the bytes to the log file and save the number of bytes written and any error
	n2, err2 := lw.logFile.Write(p)

	// if an error was encountered, return the number of bytes written and the error
	if err != nil {
		return n2, err2
	}

	// return the number of bytes written and nil
	return n1, nil
}
