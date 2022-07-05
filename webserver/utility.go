package webserver

import (
	"os"
	"os/exec"
	"runtime"
)

//code to be injected
var INJECTABLE_CODE = []byte(`
<script>
	if ('WebSocket' in window) {
		var protocol = 'ws://';
		var address = protocol + window.location.host + '/sakthi/mahendran/2005/ws'
		var socket = new WebSocket(address);
	
		socket.onmessage = function (msg) {
			console.log(msg.data)
		    if (msg.data == 'reload') {
				window.location.reload();
			}
		}

	} else {
		window.alert("Browser does'nt support live reload. Please upgrade your browser (by LiveServer)")
		console.log("Browser does'nt support live reload. Please upgrade your browser (by LiveServer)")
	}
</script>
`)

type utility struct {
}

//injects the code from INJECTABLE_CODE_PATH
func (*utility) injectCode(htmlFilePath string) ([]byte, error) {
	fileContent, err := os.ReadFile(htmlFilePath)

	if err != nil {
		return nil, err
	}

	return append(fileContent, INJECTABLE_CODE...), nil
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
