package webserver

import (
	"filelistener"
	"net/http"
	"statuslogger"
	"time"

	"github.com/gorilla/websocket"
)

func NewWsServer(sl *statuslogger.StatusLogger) WsServer {
	ws := WsServer{}
	ws.mfl = filelistener.NewMultiFileListener()
	ws.logger = sl

	return ws
}

type WsServer struct {
	mfl    filelistener.MultiFileListener
	con    *websocket.Conn
	logger *statuslogger.StatusLogger
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

	lc := ws.mfl.GetListenChan()

	go func(w http.ResponseWriter, r *http.Request, listenChan <-chan struct{}) {
		for range listenChan {
			err := ws.con.WriteMessage(websocket.TextMessage, []byte("reload"))

			if err != nil {
				ws.logger.LogErr(err.Error())
			}

			ws.con.Close()
			break
		}
	}(w, r, lc)
}

func (ws *WsServer) AddFileListener(filePath string, rate time.Duration) {
	fl := filelistener.NewFileListener(filePath, rate)
	ws.mfl.Add(fl)
}
