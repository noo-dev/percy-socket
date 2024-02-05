package main

import (
	"github.com/gorilla/websocket"
	"log"
)

type ClientList map[*Client]bool

type Client struct {
	connection *websocket.Conn
	wsManager  *SocketManager

	egress chan []byte
}

func NewClient(connection *websocket.Conn, wsManager *SocketManager) *Client {
	return &Client{
		connection: connection,
		wsManager:  wsManager,
		egress:     make(chan []byte),
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

		for wsClient := range c.wsManager.clients {
			wsClient.egress <- payload
		}

		log.Println("MESSAGE TYPE: ", messageType)
		log.Println("PAYLOAD: ", string(payload))
	}
}

func (c *Client) writeMessages() {
	defer func() {
		c.wsManager.removeClient(c)
	}()

	for {
		select {
		case message, ok := <-c.egress:
			if !ok {
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("connection closed: ", err)
				}
				return
			}

			err := c.connection.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Printf("Failed to send message: %v", err)
			}
			log.Println("message sent")
		}

	}
}
