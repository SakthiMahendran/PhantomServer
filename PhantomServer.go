package main

import (
	"bufio"
	"os"
	"runtime"
	"strings"

	"github.com/SakthiMahendran/PhantomServer/statuslogger"
	"github.com/SakthiMahendran/PhantomServer/webserver"
)

// initialize a buffered reader to read input from stdin
var reader = bufio.NewReader(os.Stdin)

// set the GOMAXPROCS to 4 if it's less than 4
func init() {
	if runtime.GOMAXPROCS(-1) < 4 {
		runtime.GOMAXPROCS(4)
	}
}

func main() {
	//initialize status logger
	sl := statuslogger.NewStatusLogger()

	//initialize http server with the status logger
	hs := webserver.NewHttpServer(&sl)

	//initialize a command executor with the http server and status logger
	cmdExe := NewCmdExe(&hs, &sl)

	//log that the server is ready to use
	sl.LogInfo("Ready to use...")

	//loop for reading commands from stdin and executing them
	for {
		//print a new line
		sl.NewLine()

		//prompt for user input
		os.Stdout.WriteString("Enter your command: ")
		cmd := readln()

		//execute the command if it's not empty
		if cmd != "" {
			cmdExe.exe(cmd)
		}
	}
}

// function to read a line from the buffered reader and return the trimmed string
func readln() string {
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)

	return line
}
