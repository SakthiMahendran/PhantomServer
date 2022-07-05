package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"statuslogger"
	"strings"
	"webserver"
)

var reader = bufio.NewReader(os.Stdin)

func main() {

	setMinThreads()

	sl := statuslogger.NewStatusLogger()
	hs := webserver.NewHttpServer(&sl)

	var cmd string
	cmdExe := NewCmdExe(&hs, &sl)

	sl.LogInfo("Ready to use...")

	for {
		sl.NewLine()
		fmt.Print("Enter your command: ")
		cmd = readln()

		if cmd != "" {
			sl.LogInfo("Executing Command ", "\"", cmd, "\"")
			cmdExe.exe(cmd)
		}
	}
}

func readln() string {
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)

	return line
}

func setMinThreads() {
	if runtime.GOMAXPROCS(-1) < 4 {
		runtime.GOMAXPROCS(4)
	}
}
