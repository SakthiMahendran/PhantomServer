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
		<-listenChan

		err := ws.Reload()

		if err != nil {
			ws.logger.LogErr(err.Error())
		}

		ws.con.Close()

	}(w, r, lc)
}

func (ws *WsServer) Reload() error {
	err := ws.con.WriteMessage(websocket.TextMessage, []byte("reload"))
	ws.con.Close()
	return err
}

func (ws *WsServer) AddFileListener(filePath string) {
	fl := filelistener.NewFileListener(filePath)
	ws.mfl.Add(fl)
}
