package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var (
	websocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type SocketManager struct {
}

func NewSocketManager() *SocketManager {
	return &SocketManager{}
}

func (m *SocketManager) serveWS(w http.ResponseWriter, r *http.Request) {
	log.Println("new connection")

	// upgrade regular http connection into websocket
	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	conn.Close()
}
