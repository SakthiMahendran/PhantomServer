package main

import (
	"bufio"
	"os"
	"runtime"
	"strings"

	"github.com/SakthiMahendran/PhantomServer/statuslogger"
	"github.com/SakthiMahendran/PhantomServer/webserver"
)

var reader = bufio.NewReader(os.Stdin)

func init() {
	if runtime.GOMAXPROCS(-1) < 4 {
		runtime.GOMAXPROCS(4)
	}

}

func main() {
	sl := statuslogger.NewStatusLogger()
	hs := webserver.NewHttpServer(&sl)

	var cmd string
	cmdExe := NewCmdExe(&hs, &sl)

	sl.LogInfo("Ready to use...")

	for {
		sl.NewLine()
		os.Stdout.WriteString("Enter your command: ")
		cmd = readln()

		if cmd != "" {
			cmdExe.exe(cmd)
		}
	}
}

func readln() string {
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)

	return line
}
