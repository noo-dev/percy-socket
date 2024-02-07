package main

import (
	"encoding/json"
	"time"
)

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type EventHandler func(event Event, c *Client) error

const (
	EventSendMessage    = "send_message"
	EventNewMessage     = "new_message"
	EventChangeChatroom = "change_chatroom"
)

type SendMessageEvent struct {
	Message string `json:"message"`
	From    string `json:"from"`
}

type NewMessageEvent struct {
	SendMessageEvent
	SentTime time.Time `json:"sent_time"`
}

type ChangeRoomEvent struct {
	Name string `json:"name"`
}
