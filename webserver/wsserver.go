package webserver

import (
	"net/http"

	"github.com/SakthiMahendran/PhantomServer/filelistener"
	"github.com/SakthiMahendran/PhantomServer/statuslogger"
	"github.com/gorilla/websocket"
)

// WsServer is a WebSocket server that reloads the webpage if there are changes made to resource files served by the HTTP server.
type WsServer struct {
	mfListener filelistener.MultiFileListener // MultiFileListener for listening to changes in resource files.
	con        *websocket.Conn                // WebSocket connection for reloading webpage.
	logger     *statuslogger.StatusLogger     // StatusLogger for logging.
}

// NewWsServer creates a new WsServer.
func NewWsServer(sl *statuslogger.StatusLogger) WsServer {
	ws := WsServer{
		mfListener: filelistener.NewMultiFileListener(), // Initialize MultiFileListener.
		logger:     sl,                                  // Set StatusLogger.
	}
	return ws
}

// Start upgrades a HTTP connection to a WebSocket connection and listens for changes in resource files.
func (ws *WsServer) Start(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection to WebSocket connection.
	up := websocket.Upgrader{
		ReadBufferSize:  0,  // No data is going to be read.
		WriteBufferSize: 32, // Set write buffer size.
	}

	con, err := up.Upgrade(w, r, nil)
	if err != nil {
		// If there is an error, log it and return.
		ws.logger.LogErr(err.Error())
		return
	}

	ws.con = con // Store the connection object pointer in shared memory.

	// Get the listen channel.
	lc := ws.mfListener.GetListenChan()

	// Listen for changes in resource files and reload the webpage.
	go func(listenChan <-chan struct{}) {
		<-listenChan // Wait for signal.
		ws.mfListener.Reset()
		for {
			err := ws.Reload() // Reload the webpage.
			if err != nil {    // If there is an error, try again.
				continue
			}
		}
	}(lc)
}

// Reload sends a reload message to the JavaScript client in the webpage, and then closes the WebSocket connection.
func (ws *WsServer) Reload() error {
	if ws.con != nil {
		// If ws.con is not nil, send a reload message to the JavaScript client.
		err := ws.con.WriteMessage(websocket.TextMessage, []byte("reload"))
		if err == nil {
			ws.con.Close() // Close the connection.
		}
		return err // Return any error.
	} else {
		return nil // Return nil if ws.con is nil.
	}
}

// AddFileListener adds a new FileListener to the MultiFileListener.
func (ws *WsServer) AddFileListener(filePath string) {
	// Create a new FileListener and add it to the MultiFileListener.
	fl := filelistener.NewFileListener(filePath)
	ws.mfListener.Add(&fl)
}
