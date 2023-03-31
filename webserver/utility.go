package webserver

import (
	"os"
	"os/exec"
	"runtime"
	"strconv"
)

// This constant holds the JavaScript code to be injected into HTML pages for live reloading
var INJECTABLE_CODE = []byte(`
	<script>
		if ('WebSocket' in window) {
			var protocol = 'ws://';
			var address = protocol + window.location.host + '/sakthi/mahendran/2005/ws'
			var socket = new WebSocket(address);
		
			socket.onmessage = function (msg) {
				if (msg.data == 'reload') {
					window.location.reload();
				}
			}
	
		} else {
			window.alert("Browser does'nt support live reload. Please upgrade your browser (by PhantomServer)")
			console.log("Browser does'nt support live reload. Please upgrade your browser (by PhantomServer)")
		}
	</script>
`)

type utility struct {
}

// This method injects the INJECTABLE_CODE into the HTML file specified by htmlFilePath
func (*utility) injectCode(htmlFilePath string) ([]byte, error) {
	fileContent, err := os.ReadFile(htmlFilePath)

	if err != nil {
		return nil, err
	}

	return append(fileContent, INJECTABLE_CODE...), nil
}

// This method checks whether the file specified by filePath is an HTML file
func (*utility) hasHtml(filePath string) bool {
	return filePath[len(filePath)-5:] == ".html"
}

// This method checks whether the file specified by filePath exists
func (*utility) validPath(filePath string) bool {
	if _, err := os.Stat(filePath); err == nil {
		return true
	}

	return false
}

// This method checks whether the given port is valid
func (*utility) validPort(port string) bool {
	if _, err := strconv.Atoi(port); err == nil {
		return true
	}

	return false
}

// This method opens the specified URL in the default browser
func (*utility) openBrowser(url string) error {
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
