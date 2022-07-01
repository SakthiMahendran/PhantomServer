module webserver

require filelistener v1.0.0

replace filelistener v1.0.0 => /filelistener

require (
	statuslogger v1.0.0 // indirect
)

replace statuslogger v1.0.0 => /statuslogger


go 1.18
