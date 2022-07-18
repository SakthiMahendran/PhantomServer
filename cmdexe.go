package main

import (
	"fmt"
	"strings"

	"github.com/SakthiMahendran/PhantomServer/statuslogger"
	"github.com/SakthiMahendran/PhantomServer/webserver"
)

func NewCmdExe(hs *webserver.HttpServer, sl *statuslogger.StatusLogger) CmdExe {
	ce := CmdExe{}
	ce.httpServer = hs
	ce.logger = sl

	return ce
}

type CmdExe struct {
	httpServer *webserver.HttpServer
	logger     *statuslogger.StatusLogger
}

func (ce *CmdExe) exe(cmd string) {
	cmdArr := ce.removeSpace(strings.Split(cmd, " "))
	key := strings.TrimSpace(cmdArr[0])
	argsLen := len(cmdArr) - 1

	ce.logger.LogInfo("Executing Command ", "\"", cmd, "\"")

	switch key {
	case "setmain":
		if argsLen == 1 {
			main := strings.TrimSpace(cmdArr[1])
			ce.httpServer.SetMainHtml(main)
			break
		} else {
			ce.logger.LogErr("setmain takes 1 argument but ", argsLen, " provided.")
			break
		}
	case "setfavicon":
		if argsLen == 1 {
			fav := strings.TrimSpace(cmdArr[1])
			ce.httpServer.SetFavIcon(fav)
			break
		} else {
			ce.logger.LogErr("setfavicon takes 1 arguments but ", argsLen, " provided.")
			break
		}
	case "link":
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
		fmt.Println("setmain   	-> Sets the path for MainHtml file (syntax: setmain MainHtmlPath).")
		fmt.Println("setfavicon	-> Sets the path for PageIcon 	   (syntax: setfavicon FavIconPath).")
		fmt.Println("start 	   	-> Starts the Http server at given port else in default port(80) (synatx: start PortNumber or start).")
		fmt.Println("link 	   	-> Links a resourse with request (syntax: link ResPath ReqPath).")
		fmt.Println("help      	-> Gives info about the available commands.")
		break
	default:
		ce.logger.LogErr(cmd, " is not a valid Command.")
		ce.logger.LogInfo("type \"help\" for Help.")
		break
	}

}

func (*CmdExe) removeSpace(cmdArr []string) []string {
	resultArr := make([]string, 0)

	for _, elem := range cmdArr {
		if strings.TrimSpace(elem) != "" {
			resultArr = append(resultArr, elem)
		}
	}

	return resultArr
}
