package main

import (
	"github.com/gorilla/websocket"
	"log"
)

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

func (c *Client) readMessages() {
	defer func() {

		c.wsManager.removeClient(c)
	}()

	for {
		messageType, payload, err := c.connection.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			}

			break
		}

		log.Println("MESSAGE TYPE: ", messageType)
		log.Println("PAYLOAD: ", string(payload))
	}
}
