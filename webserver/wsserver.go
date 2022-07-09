package webserver

import (
	"filelistener"
	"net/http"
	"statuslogger"

	"github.com/gorilla/websocket"
)

//First read the code in "filelistener" package, description and then the code line by line with the comments to for a better understanding.

//Description
/*
	WsServer (WebSocket Server) will reload the webpage if there is any
	changes made to resource file that is served by the HttpServe.
*/

//Makes a new WsServer.
func NewWsServer(sl *statuslogger.StatusLogger) WsServer {
	ws := WsServer{}                             //Instantiation.
	ws.mfl = filelistener.NewMultiFileListener() //Setting the "MultiFileListener" (for listening changes in resource file).
	ws.logger = sl                               // Setting the statuslogger (for logging).

	return ws
}

type WsServer struct {
	mfl    filelistener.MultiFileListener //"MultiFileListener" for listening changes in resource file.
	con    *websocket.Conn                // WebSocket connection for reloading webpage.
	logger *statuslogger.StatusLogger     //"StatusLogger" for logging.
}

//UpGrades a Http connection into a WebSocket connection.
func (ws *WsServer) Start(w http.ResponseWriter, r *http.Request) {
	up := websocket.Upgrader{ // Defining buffer size.
		ReadBufferSize:  0, // No data is going to be readed.
		WriteBufferSize: 32,
	}

	con, err := up.Upgrade(w, r, nil) //Upgrading

	if err != nil { //Cheacking for error
		ws.logger.LogErr(err.Error()) //if error then log error
		return                        // and return
	}

	ws.con = con                 //Writting the connection obeject pointer into shared memory (else it will be cleared when it exits the scope)
	lc := ws.mfl.GetListenChan() //Getting the "listenChan" (Signal will come through this channel if any changes is made to the resource files)

	//Starting a new goroutine
	go func(w http.ResponseWriter, r *http.Request, listenChan <-chan struct{}) {
		<-listenChan // Waiting for signal

		err := ws.Reload() // Reloading the webpage

		if err != nil {
			//If error then log
			ws.logger.LogErr(err.Error())
		}

	}(w, r, lc)
}

func (ws *WsServer) Reload() error {
	if ws.con != nil { //"ws.con" should not be nil pointer
		err := ws.con.WriteMessage(websocket.TextMessage, []byte("reload")) //Send reload message for JavaScript client in the webpage (That will be injected by HttpServer)
		ws.con.Close()                                                      //Close the connection (Reloading makes the webpage to make another WebSocket request so close previous connection)
		return err                                                          // return if any error
	} else {
		return nil //Just return nil if "ws.con" is nil ptr
	}
}

//Adds a new FileListener
func (ws *WsServer) AddFileListener(filePath string) {
	fl := filelistener.NewFileListener(filePath) // Makes a new FileListener
	ws.mfl.Add(fl)                               //Add it to MultiFileListener
}
