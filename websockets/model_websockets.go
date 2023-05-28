package websockets

import (
	"context"
	"sync"
	"time"

	"github.com/arkhamHack/VerbiNative-backend/messages"
	"github.com/gorilla/websocket"
)

const (
	MaxMessageSize = 1024
	WriteWait      = 10 * time.Second
	PingInterval   = 10 * time.Second
	PongWait       = 20 * time.Second
)

type webSocketClient struct {
	id          string
	chatroom_id string
	ws          *websocket.Conn
	msgs        chan WebSocketMessages
	err         chan error
	done        chan interface{}
	mutex       sync.Mutex
	once        sync.Once
}

type WebSocketClient interface {
	Id() string
	ChatroomId() string
	Launch(ctx context.Context)
	Write(m WebSocketMessages) error
	Close() error
	Listen() <-chan WebSocketMessages
	Done() <-chan interface{}
	Error() <-chan error
}
type WebSocketMessages struct {
	Type    string
	Content messages.Msg
}
type WebSocketClientsPool []WebSocketClient
