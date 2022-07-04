package webserver

import (
	"net/http"
	"statuslogger"
	"time"
)

//Defines the WS(Web Socket) Request_Path
const WS_REQUEST_PATH string = "/sakthi/mahendran/2005/ws"

//Defines the
const RATE time.Duration = time.Millisecond * 250

//Makes a new HttpServer at the given port
func NewHttpServer(sl *statuslogger.StatusLogger) HttpServer {

	hs := HttpServer{}
	hs.wsServer = NewWsServer(sl)
	hs.port = "80"
	hs.requestMap = make(map[string]string)
	hs.running = false
	hs.logger = sl
	hs.util = utility{}

	return hs
}

//HttpServer struct
type HttpServer struct {
	requestMap   map[string]string
	port         string
	running      bool
	favIconPath  string
	mainHtmlPath string
	wsServer     WsServer
	logger       *statuslogger.StatusLogger
	util         utility
}

//Starts the server if it is not already started
//Throws error if already Started
func (hs *HttpServer) Start() {
	hs.logger.LogInfo("Starting server at port ", hs.port, ".")

	if hs.running {
		hs.logger.LogErr("Server already started can not start again.")
		return
	}

	hs.running = true

	http.HandleFunc(WS_REQUEST_PATH, hs.wsServer.Start)
	http.HandleFunc("/", hs.requestHandler)

	go http.ListenAndServe(":"+hs.port, nil)

	hs.logger.LogInfo("Server started.")

	hs.logger.LogInfo("Opening webpage in browser.")
	err := hs.util.openbrowser("http://localhost/")

	if err != nil {
		hs.logger.LogErr(err)
	} else {
		hs.logger.LogInfo("Webpage opened in browser.")
	}

}

func (hs *HttpServer) SetPort(port string) {
	hs.logger.LogInfo("Setting ", port, " as Server Port")

	if !hs.running {
		hs.port = port
		hs.logger.LogInfo("Server port is seted to ", port)
	} else {
		hs.logger.LogErr("Server already started can not change the port.")
	}
}

//Sets the FavIcon file path
func (hs *HttpServer) SetFavIcon(favIconPath string) {
	hs.logger.LogInfo("Setting ", favIconPath, " as FavIconPath.")

	if hs.util.validPath(favIconPath) {
		hs.logger.LogErr(favIconPath, " is not a valid path.")
		return
	}

	hs.favIconPath = favIconPath
	hs.requestMap["/favicon.ico"] = favIconPath

	hs.logger.LogInfo(favIconPath, " is seted as FavIconPath.")
}

//Sets the main html file path
//returns error if the givenFilePath does not contain a html file
func (hs *HttpServer) SetMainHtml(mainHtmlPath string) {
	hs.logger.LogInfo("Setting ", mainHtmlPath, " as MainHtml file.")

	if hs.util.validPath(mainHtmlPath) {
		hs.logger.LogErr(mainHtmlPath, " is not a valid path.")
		return
	}

	if hs.util.hasHtml(mainHtmlPath) {
		hs.logger.LogErr(mainHtmlPath, " is not a Html file.")
		return
	}

	hs.mainHtmlPath = mainHtmlPath
	hs.requestMap["/"] = mainHtmlPath

	hs.logger.LogInfo(mainHtmlPath, " was seted as MainHtml file.")
}

//Links request_url_path with file_path
//So that if the request contains the given url_path then file from the given file_path is responded
func (hs *HttpServer) LinkRes(rqst, res string) {
	if hs.util.validPath(res) {
		hs.requestMap[rqst] = res
		hs.logger.LogInfo("Linked.")
	} else {
		hs.logger.LogErr(res, " is not a ValidPath")
	}
}

//Indicates whether the server is running or not
func (hs *HttpServer) IsRunning() bool {
	return hs.running
}

//Gives the port number of the server
func (hs *HttpServer) GetPort() string {
	return hs.port
}

//Handles the incoming request
//Responds with the appropriate file for the request_url_path from the rqstMap
func (hs *HttpServer) requestHandler(w http.ResponseWriter, r *http.Request) {

	hs.logger.NewLine()
	hs.logger.LogInfo("request: ", r.URL.Path)

	if filePath, ok := hs.requestMap[r.URL.Path]; ok {

		go hs.wsServer.AddFileListener(filePath, RATE)

		hs.logger.LogInfo("response: ", filePath)

		if filePath == hs.mainHtmlPath {
			hs.respondMainHtml(w, r)
		} else if filePath == hs.favIconPath {
			hs.respondFavIcon(w, r)

		} else {
			hs.respondResFile(w, r, filePath)
		}

	} else {
		w.WriteHeader(http.StatusNotFound)
		hs.logger.LogErr("response: resource not found (404)")
	}
}

//responds the main html file with injected javascript code
func (hs *HttpServer) respondMainHtml(w http.ResponseWriter, r *http.Request) {
	injected, err := hs.util.injectCode(hs.mainHtmlPath)

	if err != nil {
		hs.logger.LogErr(err.Error())
	}

	w.Write(injected)
}

//responds the Resource File
func (hs *HttpServer) respondResFile(w http.ResponseWriter, r *http.Request, resFilePath string) {
	http.ServeFile(w, r, resFilePath)
}

//responds the FavIcon
func (hs *HttpServer) respondFavIcon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, hs.favIconPath)
}
