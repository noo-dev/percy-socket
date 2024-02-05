package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

var (
	pongWaitTimeout = 10 * time.Second

	pingInterval = (pongWaitTimeout * 9) / 10
)

type ClientList map[*Client]bool

type Client struct {
	connection *websocket.Conn
	wsManager  *SocketManager

	egress chan Event
}

func NewClient(connection *websocket.Conn, wsManager *SocketManager) *Client {
	return &Client{
		connection: connection,
		wsManager:  wsManager,
		egress:     make(chan Event),
	}
}

func (c *Client) readMessages() {
	defer func() {
		c.wsManager.removeClient(c)
	}()

	if err := c.connection.SetReadDeadline(time.Now().Add(pongWaitTimeout)); err != nil {
		log.Println(err)
		return
	}

	c.connection.SetPongHandler(c.pongHandler)

	for {
		_, payload, err := c.connection.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			}

			break
		}

		var request Event

		if err := json.Unmarshal(payload, &request); err != nil {
			log.Printf("error marshalling event: %v", err)
			break
		}

		if err := c.wsManager.routeEvent(request, c); err != nil {
			log.Println("error handling message: ", err)
		}

	}
}

func (c *Client) writeMessages() {
	defer func() {
		c.wsManager.removeClient(c)
	}()

	ticker := time.NewTicker(pingInterval)

	for {
		select {
		case message, ok := <-c.egress:
			if !ok {
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("connection closed: ", err)
				}
				return
			}

			data, err := json.Marshal(message)
			if err != nil {
				log.Println(err)
				return
			}

			err = c.connection.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Printf("Failed to send message: %v", err)
			}
			log.Println("message sent")
		case <-ticker.C:
			log.Println("ping to client")

			// send a ping to the client
			if err := c.connection.WriteMessage(websocket.PingMessage, []byte(``)); err != nil {
				log.Println("writemsg err: ", err)
				return
			}
		}

	}
}

func (c *Client) pongHandler(pongStr string) error {
	log.Println("pong")
	return c.connection.SetReadDeadline(time.Now().Add(pongWaitTimeout))
}
