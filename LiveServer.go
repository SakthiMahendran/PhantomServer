package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"
	"webserver"
)

var reader = bufio.NewReader(os.Stdin)

func main() {

	setMinThreads()

	hs := webserver.NewHttpServer()
	var cmd string
	cmdExe := NewCmdExe()

	fmt.Println("Ready to use...")

	for {
		fmt.Println()
		fmt.Print("Enter your command: ")
		cmd = readln()

		if cmd != "" {
			cmdExe.exe(cmd, &hs)
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
