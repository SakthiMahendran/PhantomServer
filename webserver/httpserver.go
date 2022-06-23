package webserver

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

//Defines the path for the javascript file which contains the code to be injected
const INJECTABLE_CODE_PATH string = "injectable_code/injectable.html"

//Defines the WS(Web Socket) Request_Path
const WS_REQUEST_PATH string = "/sakthi/mahendran/2005/ws"

//Defines the
const RATE time.Duration = time.Millisecond * 250

//Makes a new HttpServer at the given port
func NewHttpServer() HttpServer {
	hs := HttpServer{}
	hs.wss = NewWsServer()
	hs.port = 80
	hs.requestMap = make(map[string]string)
	hs.running = false

	return hs
}

//HttpServer struct
type HttpServer struct {
	requestMap   map[string]string
	port         int
	running      bool
	favIconPath  string
	mainHtmlPath string
	wss          WsServer
}

//Starts the server if it is not already started
//Throws error if already Started
func (hs *HttpServer) Start() {
	fmt.Println("Starting server at port ", hs.port, ".")

	if hs.running {
		fmt.Println("Server already started.")
		return
	}

	hs.running = true

	http.HandleFunc(WS_REQUEST_PATH, hs.wss.Start)
	http.HandleFunc("/favicon.ico", hs.respondFavIcon)
	http.HandleFunc("/", hs.requestHandler)

	go http.ListenAndServe(fmt.Sprint(":", hs.port), nil)

	fmt.Println("Server started.")
}
func (hs *HttpServer) SetPort(port int) {
	if !hs.running {
		hs.port = port
		fmt.Println("Server port is seted to ", port)
	} else {
		fmt.Println("Server already started can not change the port.")
	}
}

//Sets the FavIcon file path
func (hs *HttpServer) SetFavIcon(favIconPath string) {
	hs.favIconPath = favIconPath
	fmt.Println(favIconPath, " is seted as FavIconPath.")
}

//Sets the main html file path
//returns error if the givenFilePath does not contain a html file
func (hs *HttpServer) SetMainHtml(mainHtmlPath string) {
	fmt.Println("Setting ", mainHtmlPath, " as MainHtml file.")

	if !hs.validPath(mainHtmlPath) {
		fmt.Println(mainHtmlPath, " is not a valid path.")
		return
	}

	if hs.hasHtml(mainHtmlPath) {
		hs.mainHtmlPath = mainHtmlPath
		hs.requestMap["/"] = mainHtmlPath

		fmt.Println(mainHtmlPath, " was seted as MainHtml file.")
	} else {
		fmt.Println(mainHtmlPath, " is not a Html file.")
	}
}

//Links request_url_path with file_path
//So that if the request contains the given url_path then file from the given file_path is responded
func (hs *HttpServer) LinkRes(rqst, res string) {
	if hs.validPath(res) {
		hs.requestMap[rqst] = res
		fmt.Println("Linked.")
	} else {
		fmt.Println(res, " is not a ValidPath")
	}
}

//Indicates whether the server is running or not
func (hs *HttpServer) IsRunning() bool {
	return hs.running
}

//Gives the port number of the server
func (hs *HttpServer) GetPort() int {
	return hs.port
}

//Handles the incoming request
//Responds with the appropriate file for the request_url_path from the rqstMap
func (hs *HttpServer) requestHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println()
	fmt.Println("request: ", r.URL.Path)

	if filePath, ok := hs.requestMap[r.URL.Path]; ok {

		go hs.wss.AddFileListener(filePath, RATE)
		fmt.Println("response: ", filePath)

		if filePath == hs.mainHtmlPath {
			hs.respondMainHtml(w, r)

		} else if filePath == hs.favIconPath {
			hs.respondFavIcon(w, r)

		} else {
			hs.respondResFile(w, r, filePath)

		}
	} else {
		fmt.Println("response: Page not found (404)")
	}
}

//responds the main html file with injected javascript code
func (hs *HttpServer) respondMainHtml(w http.ResponseWriter, r *http.Request) {
	injected, err := hs.injectCode(hs.mainHtmlPath)

	if err != nil {
		fmt.Printf(err.Error())
	}

	fmt.Fprint(w, string(injected))
}

//responds the Resource File
func (hs *HttpServer) respondResFile(w http.ResponseWriter, r *http.Request, resFilePath string) {
	http.ServeFile(w, r, resFilePath)
}

//responds the FavIcon
func (hs *HttpServer) respondFavIcon(w http.ResponseWriter, r *http.Request) {
	if hs.favIconPath != "" {
		http.ServeFile(w, r, hs.favIconPath)
	}
}

//injects the code from INJECTABLE_CODE_PATH
func (hs *HttpServer) injectCode(htmlFilePath string) ([]byte, error) {
	file, err := ioutil.ReadFile(htmlFilePath)
	if err != nil {
		return nil, err
	}

	code, err := ioutil.ReadFile(INJECTABLE_CODE_PATH)
	if err != nil {
		return nil, err
	}

	return append(file, []byte(code)...), nil
}

//Checks whether the given filepath contains html file
func (hs *HttpServer) hasHtml(filePath string) bool {
	return strings.HasSuffix(filePath, ".html")
}

func (hs *HttpServer) validPath(filePath string) bool {
	if _, err := os.Stat(filePath); err == nil {
		return true
	}

	return false
}
