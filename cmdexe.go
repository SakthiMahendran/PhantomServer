package main

import (
	"fmt"
	"strconv"
	"strings"
	"webserver"
)

func NewCmdExe() CmdExe {
	return CmdExe{}
}

type CmdExe struct {
}

func (ce *CmdExe) exe(cmd string, hs *webserver.HttpServer) {
	cmdArr := strings.Split(cmd, " ")
	key := strings.TrimSpace(cmdArr[0])
	cmdArrLen := len(cmdArr)

	switch key {
	case "setmain":
		if cmdArrLen == 2 {
			main := strings.TrimSpace(cmdArr[1])
			hs.SetMainHtml(main)
			break
		} else {
			fmt.Println("setmain takes 1 arguments but ", cmdArrLen, " provided.")
		}
	case "setfavicon":
		if cmdArrLen == 2 {
			fav := strings.TrimSpace(cmdArr[1])
			hs.SetFavIcon(fav)
			break
		} else {
			fmt.Println("setfavicon takes 1 arguments but ", cmdArrLen, " provided.")
		}
	case "link":
		if cmdArrLen == 3 {
			req := strings.TrimSpace(cmdArr[1])
			res := strings.TrimSpace(cmdArr[2])
			hs.LinkRes(req, res)
			break
		} else {
			fmt.Println("link takes 2 arguments but ", cmdArrLen, " provided.")
		}
	case "start":
		if cmdArrLen == 2 {
			port, err := strconv.Atoi(strings.TrimSpace(cmdArr[1]))

			if err == nil {
				hs.SetPort(port)
				hs.Start()
			} else {
				fmt.Println(err.Error())
			}
		} else {
			hs.Start()
		}
		break
	case "help":
		fmt.Println("setmain   	-> Sets the path for MainHtml file (syntax: setmain MainHtmlPath).")
		fmt.Println("setfavicon	-> Sets the path for PageIcon 	   (syntax: setfavicon FavIconPath).")
		fmt.Println("start 	   	-> Starts the Http server at given port else in default port(80) (synatx: start PortNumber or start).")
		fmt.Println("link 	   	-> Links a resourse with request (syntax: link ResPath ReqPath).")
		fmt.Println("help      	-> Gives info about the available commands.")
		break
	default:
		fmt.Println(cmd, " is not a valid Command.")
		fmt.Println("type \"help\" for Help.")
		break
	}

}
