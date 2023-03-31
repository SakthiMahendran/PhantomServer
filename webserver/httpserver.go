package webserver

import (
	"net/http"

	"github.com/SakthiMahendran/PhantomServer/statuslogger"
)

// Defines the WebSocket request path
const WS_REQUEST_PATH string = "/sakthi/mahendran/2005/ws"

// Defines the HTTP server struct
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

// Defines a new HTTP server at the given port
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

// Starts the server if it is not already started
// Throws error if already started
func (hs *HttpServer) Start() {
	hs.logger.LogInfo("Starting server at port ", hs.port, ".")

	if hs.running {
		hs.logger.LogErr("Server already started can not start again.")
		return
	}

	if hs.mainHtmlPath == "" {
		hs.logger.LogErr("Can't start server without a MainHtml file.")
		return
	}

	hs.running = true

	http.HandleFunc(WS_REQUEST_PATH, hs.wsServer.Start)
	http.HandleFunc("/", hs.requestHandler)

	go http.ListenAndServe(":"+hs.port, nil)

	hs.logger.LogInfo("Server started.")
	hs.logger.LogInfo("Opening webpage in browser.")

	err := hs.util.openBrowser("http://localhost:" + hs.port + "/")

	if err != nil {
		hs.logger.LogErr(err)
	} else {
		hs.logger.LogInfo("Webpage opened in browser.")
	}

}

// Sets the server port
func (hs *HttpServer) SetPort(port string) {
	hs.logger.LogInfo("Setting ", port, " as Server Port.")

	if !hs.util.validPort(port) {
		hs.logger.LogErr(port, " is not a valid port.")
		return
	}

	if hs.running {
		hs.port = port
		hs.logger.LogErr("Server already started can not change the port.")
		return
	}

	hs.logger.LogInfo("Server port is set to ", port)
}

// Sets the favicon file path
func (hs *HttpServer) SetFavIcon(favIconPath string) {
	hs.logger.LogInfo("Setting ", favIconPath, " as FavIconPath.")

	if !hs.util.validPath(favIconPath) {
		hs.logger.LogErr(favIconPath, " is not a valid path.")
		return
	}

	hs.favIconPath = favIconPath
	hs.requestMap["/favicon.ico"] = favIconPath

	hs.wsServer.Reload()

	hs.logger.LogInfo(favIconPath, " is set as FavIconPath.")
}

// Sets the main html file path and validates if it is a valid html file path
func (hs *HttpServer) SetMainHtml(mainHtmlPath string) {
	hs.logger.LogInfo("Setting ", mainHtmlPath, " as MainHtml file.")

	// validate if the mainHtmlPath is a valid file path
	if !hs.util.validPath(mainHtmlPath) {
		hs.logger.LogErr(mainHtmlPath, " is not a valid path.")
		return
	}

	// validate if the mainHtmlPath is a html file
	if !hs.util.hasHtml(mainHtmlPath) {
		hs.logger.LogErr(mainHtmlPath, " is not a Html file.")
		return
	}

	// set the mainHtmlPath as the main html file path
	hs.mainHtmlPath = mainHtmlPath
	hs.requestMap["/"] = mainHtmlPath

	// reset the mfListener of websocket server and reload the server
	hs.wsServer.mfListener.Reset()
	hs.wsServer.Reload()

	hs.logger.LogInfo(mainHtmlPath, " was seted as MainHtml file.")
}

// Links request_url_path with file_path
// This function is used to map a URL path to a file path, so that if a request contains the mapped URL path, the server will respond with the file from the mapped file path
func (hs *HttpServer) LinkRes(reqst, resPath string) {
	// validate if the resPath is a valid file path
	if hs.util.validPath(resPath) {
		// add the mapping to the request map
		hs.requestMap[reqst] = resPath
		// reload the websocket server
		hs.wsServer.Reload()
		hs.logger.LogInfo("Linked.")
	} else {
		hs.logger.LogErr(resPath, " is not a ValidPath.")
	}
}

// Handles the incoming request and responds with the appropriate file for the requested URL path from the request map
func (hs *HttpServer) requestHandler(w http.ResponseWriter, r *http.Request) {
	// log the incoming request
	hs.logger.NewLine()
	hs.logger.LogInfo("Request: ", r.URL.Path)

	if filePath, ok := hs.requestMap[r.URL.Path]; ok {
		//Disable Browser cache.
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate") // HTTP 1.1.
		w.Header().Set("Pragma", "no-cache")                                   // HTTP 1.0.
		w.Header().Set("Expires", "0")                                         // Proxies.

		if filePath == hs.mainHtmlPath {
			// if the requested file is the main html file, inject javascript code and respond with the modified file
			hs.respondMainHtml(w, r)
		} else if filePath == hs.favIconPath {
			// if the requested file is the favicon, respond with the favicon
			hs.respondFavIcon(w, r)
		} else {
			// if the requested file is any other file, respond with the requested file
			hs.respondResFile(w, r, filePath)
		}

		// log the response
		hs.logger.LogInfo("Response: ", filePath)
		// start listening for file changes on the file path in a new goroutine
		go hs.wsServer.AddFileListener(filePath)

	} else {
		// if the requested file is not in the request map, respond with a 404 Not Found status
		w.WriteHeader(http.StatusNotFound)
		hs.logger.LogErr("Response: resource not found (404)")
	}
}

// responds the main html file with injected javascript code
func (hs *HttpServer) respondMainHtml(w http.ResponseWriter, r *http.Request) {
	injected, err := hs.util.injectCode(hs.mainHtmlPath)

	if err != nil {
		hs.logger.LogErr(err.Error())
		return
	}

	w.Write(injected)
}

// responds the Resource File
func (hs *HttpServer) respondResFile(w http.ResponseWriter, r *http.Request, resFilePath string) {
	http.ServeFile(w, r, resFilePath)
}

// responds the FavIcon
func (hs *HttpServer) respondFavIcon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, hs.favIconPath)
}
