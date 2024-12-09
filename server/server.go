package server

import (
	"github.com/gorilla/websocket"
	"github.com/jonnywei/yi_tunnel/common"
	"log"
	"net/http"
	"strconv"
)

type WebSocketServer struct {
	Config common.ServerConfig
}

func (s *WebSocketServer) serveWs(w http.ResponseWriter, r *http.Request) {

	upgrader := websocket.Upgrader{
		WriteBufferSize: 8192,
		ReadBufferSize:  8192,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		return
	}

	tunnel := NewWebSocketServerTunnel(ws, s.Config.Secret)

	for {
		mt, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("read error", err)
			break
		}

		log.Printf("recv from tunnel msg length = %d", len(message))
		if mt == websocket.TextMessage {
			log.Println("receive text msg: ", string(message))
			continue
		}
		err = tunnel.OnMessage(message)
		if err != nil {
			log.Println("message process :", err)
			break
		}
	}

	tunnel.OnClose()
}

func (s *WebSocketServer) Listen() {

	http.HandleFunc(s.Config.Path, s.serveWs)
	wsAddr := s.Config.Listen + ":" + strconv.Itoa(s.Config.Port)
	log.Println("websocket listen on " + wsAddr)
	http.ListenAndServe(wsAddr, nil)
}
