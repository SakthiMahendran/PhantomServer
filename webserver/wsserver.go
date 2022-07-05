package webserver

import (
	"filelistener"
	"net/http"
	"statuslogger"

	"github.com/gorilla/websocket"
)

func NewWsServer(sl *statuslogger.StatusLogger) WsServer {
	ws := WsServer{}
	ws.mfl = filelistener.NewMultiFileListener()
	ws.logger = sl
	ws.connected = false

	return ws
}

type WsServer struct {
	mfl       filelistener.MultiFileListener
	con       *websocket.Conn
	logger    *statuslogger.StatusLogger
	connected bool
}

func (ws *WsServer) Start(w http.ResponseWriter, r *http.Request) {
	up := websocket.Upgrader{
		ReadBufferSize:  0,
		WriteBufferSize: 32,
	}

	con, err := up.Upgrade(w, r, nil)

	if err != nil {
		ws.logger.LogErr(err.Error())
		return
	}

	ws.con = con
	ws.connected = true

	lc := ws.mfl.GetListenChan()

	go func(w http.ResponseWriter, r *http.Request, listenChan <-chan struct{}) {
		<-listenChan

		err := ws.Reload()

		if err != nil {
			ws.logger.LogErr(err.Error())
		}

	}(w, r, lc)
}

func (ws *WsServer) Reload() error {
	if ws.connected {
		err := ws.con.WriteMessage(websocket.TextMessage, []byte("reload"))
		ws.con.Close()
		ws.connected = false
		return err
	} else {
		return nil
	}
}

func (ws *WsServer) AddFileListener(filePath string) {
	fl := filelistener.NewFileListener(filePath)
	ws.mfl.Add(fl)
}
