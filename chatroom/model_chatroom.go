package chatroom

import (
	"github.com/arkhamHack/VerbiNative-backend/messages"

	"github.com/gorilla/websocket"
)

type Chatroom struct {
	Name        string         `json:"name" validate:"required"`
	Chatroom_id string         `json:"chatroom_id" validate:"required"`
	Messages    []messages.Msg `json:"messages,omitempty"`
	User_ids    []string       `json:"user_ids" validate:"required"`
	Broadcast   chan messages.Msg
	Register    chan *websocket.Conn
	Unregister  chan *websocket.Conn
}

type ChatroomSummary struct {
	Name        string         `json:"name" validate:"required"`
	Chatroom_id string         `json:"chatroom_id" validate:"required"`
	Messages    []messages.Msg `json:"messages,omitempty"`
	User_ids    []string       `json:"user_ids" validate:"required"`
}
