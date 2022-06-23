package webserver

import (
	"filelistener"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

func NewWsServer() WsServer {
	ws := WsServer{}
	ws.mfl = filelistener.NewMultiFileListener()

	return ws
}

type WsServer struct {
	mfl filelistener.MultiFileListener
	con *websocket.Conn
}

func (ws *WsServer) Start(w http.ResponseWriter, r *http.Request) {
	up := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	con, _ := up.Upgrade(w, r, nil)
	ws.con = con

	lc := ws.mfl.GetListenChan()

	go func(w http.ResponseWriter, r *http.Request, listenChan <-chan struct{}) {
		for range listenChan {
			ws.con.WriteMessage(websocket.TextMessage, []byte("reload"))
			ws.con.Close()
		}
	}(w, r, lc)
}

func (ws *WsServer) AddFileListener(filePath string, rate time.Duration) {
	fl := filelistener.NewFileListener(filePath, rate)
	ws.mfl.Add(fl)
}
