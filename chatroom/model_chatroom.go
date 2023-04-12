package chatroom

import (
	"github.com/arkhamHack/VerbiNative-backend/messages"
	"github.com/arkhamHack/VerbiNative-backend/users"

	"github.com/gorilla/websocket"
)

type Chatroom struct {
	Name           string                   `json:"chat_name" validate:"required"`
	Chatroom_id    string                   `json:"chatroom_id" validate:"required"`
	Messages       []messages.Msg           `json:"messages,omitempty" validate:"required"`
	Chatroom_users []users.User             `json:"users" validate:"required"`
	Broadcast      chan messages.Msg        `json:"broadcast,omitempty"`
	Users          map[*websocket.Conn]bool `json:"messages,omitempty" validate:"required"`
	Register       chan *websocket.Conn     `json:"register,omitempty" validate:"required"`
	Unregister     chan *websocket.Conn
}
