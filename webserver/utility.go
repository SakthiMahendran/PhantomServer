package webserver

import (
	"os"
	"os/exec"
	"runtime"
)

//Defines the path for the javascript file which contains the code to be injected
const INJECTABLE_CODE_PATH string = "injectable_code/injectable.html"

type utility struct {
}

//injects the code from INJECTABLE_CODE_PATH
func (*utility) injectCode(htmlFilePath string) ([]byte, error) {
	fileContent, err := os.ReadFile(htmlFilePath)

	if err != nil {
		return nil, err
	}

	code, err := os.ReadFile(INJECTABLE_CODE_PATH)
	if err != nil {
		return nil, err
	}

	return append(fileContent, code...), nil
}

//Checks whether the given filepath contains html file
func (*utility) hasHtml(filePath string) bool {
	return filePath[len(filePath)-5:] == ".html"
}

func (*utility) validPath(filePath string) bool {
	if _, err := os.Stat(filePath); err == nil {
		return true
	}

	return false
}

//Opens the URL in default browser.
func (*utility) openbrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}

	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
