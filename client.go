package main

import "github.com/gorilla/websocket"

type ClientList map[*Client]bool

type Client struct {
	connection *websocket.Conn
	wsManager  *SocketManager
}

func NewClient(connection *websocket.Conn, wsManager *SocketManager) *Client {
	return &Client{
		connection: connection,
		wsManager:  wsManager,
	}
}
