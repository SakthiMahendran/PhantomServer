package main

import (
	"fmt"
	"strings"

	"github.com/SakthiMahendran/PhantomServer/statuslogger"
	"github.com/SakthiMahendran/PhantomServer/webserver"
)

// NewCmdExe creates a new CmdExe instance with the provided HttpServer and StatusLogger instances.
func NewCmdExe(hs *webserver.HttpServer, sl *statuslogger.StatusLogger) CmdExe {
	ce := CmdExe{}
	ce.httpServer = hs
	ce.logger = sl

	return ce
}

// CmdExe represents a command executor that can execute various commands on the HttpServer instance.
type CmdExe struct {
	httpServer *webserver.HttpServer
	logger     *statuslogger.StatusLogger
}

// exe executes the given command on the HttpServer instance.
func (ce *CmdExe) exe(cmd string) {
	// Remove extra spaces from the command and split it into an array.
	cmdArr := ce.removeSpace(strings.Split(cmd, " "))

	// Extract the key (the first element) and the number of arguments from the command array.
	key := strings.TrimSpace(cmdArr[0])
	argsLen := len(cmdArr) - 1

	// Log the command execution.
	ce.logger.LogInfo("Executing Command ", "\"", cmd, "\"")

	// Execute the command based on the key.
	switch key {
	case "setmain":
		// Set the path for the main HTML file.
		if argsLen == 1 {
			main := strings.TrimSpace(cmdArr[1])
			ce.httpServer.SetMainHtml(main)
			break
		} else {
			ce.logger.LogErr("setmain takes 1 argument but ", argsLen, " provided.")
			break
		}
	case "setfavicon":
		// Set the path for the page icon.
		if argsLen == 1 {
			fav := strings.TrimSpace(cmdArr[1])
			ce.httpServer.SetFavIcon(fav)
			break
		} else {
			ce.logger.LogErr("setfavicon takes 1 arguments but ", argsLen, " provided.")
			break
		}
	case "link":
		// Link a resource with a request path.
		if argsLen == 2 {
			req := strings.TrimSpace(cmdArr[1])
			res := strings.TrimSpace(cmdArr[2])
			ce.httpServer.LinkRes(req, res)
			break
		} else {
			ce.logger.LogErr("link takes 2 arguments but ", argsLen, " provided.")
			break
		}
	case "start":
		// Start the Http server.
		if argsLen == 1 {
			port := strings.TrimSpace(cmdArr[1])
			ce.httpServer.SetPort(port)
			ce.httpServer.Start()
			break
		} else {
			ce.httpServer.Start()
			break
		}
	case "help":
		// Print a list of available commands.
		fmt.Println("setmain   	-> Sets the path for MainHtml file (syntax: setmain MainHtmlPath).")
		fmt.Println("setfavicon	-> Sets the path for PageIcon 	   (syntax: setfavicon FavIconPath).")
		fmt.Println("start 	   	-> Starts the Http server at given port else in default port(80) (synatx: start PortNumber or start).")
		fmt.Println("link 	   	-> Links a resourse with request (syntax: link ReqPath ResPath).")
		fmt.Println("help      	-> Gives info about the available commands.")
	default:
		// default is executed when the user enters an invalid command
		ce.logger.LogErr(cmd, " is not a valid Command.")
		ce.logger.LogInfo("type \"help\" for Help.")
	}

}

// removeSpace removes any leading or trailing whitespace from the elements of cmdArr
func (*CmdExe) removeSpace(cmdArr []string) []string {
	resultArr := make([]string, 0)

	for _, elem := range cmdArr {
		if strings.TrimSpace(elem) != "" {
			resultArr = append(resultArr, elem)
		}
	}

	return resultArr
}
