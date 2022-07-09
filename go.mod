module main

require statuslogger v1.0.0

replace statuslogger v1.0.0 => /statuslogger

require (
	filelistener v1.0.0 // indirect
	github.com/daviddengcn/go-colortext v1.0.0 // indirect
)

replace filelistener v1.0.0 => /filelistener

require (
	github.com/gorilla/websocket v1.5.0 // indirect
	webserver v1.0.0
)

replace webserver v1.0.0 => /webserver

go 1.18
