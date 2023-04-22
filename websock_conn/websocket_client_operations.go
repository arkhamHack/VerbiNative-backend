package websockets

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/arkhamHack/VerbiNative-backend/messages"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var m sync.Mutex

func NewWebSocketClient(ws *websocket.Conn, userid string) WebSocketClient {
	return &webSocketClient{
		id:   userid,
		ws:   ws,
		msgs: make(chan WebSocketMessages),
		err:  make(chan error),
		done: make(chan interface{}),
	}
}

func StartClient(ctx context.Context, ws *websocket.Conn, userid string) {
	usr := NewWebSocketClient(ws, userid)
	mut.Lock()
	users = append(users, usr)
	mut.Unlock()
	defer func() {
		if err := recover(); err != nil {
			log.Printf("error:%v", err)
		}
		mut.Lock()
		defer mut.Unlock()
		users = Except(users, func(item WebSocketClient) bool {
			return item.Id() == usr.Id()
		})
		usr.Close()
	}()
	usr.Launch(ctx)
	MemberJoin(users, usr)
	for {
		select {
		case msg, ok := <-usr.Listen():
			if !ok {
				return
			} else {
				switch msg.Type {
				case "MESSAGE":
					NewMessage(users, usr, msg.Content.Text)
				default:
					log.Printf("unknown message type: %s", msg.Type)
					return

				}
			}
		case err := <-usr.Error():
			log.Printf("web socket error: %v", err)
		default:
		case <-usr.Done():
			MemberLeave(users, usr)
			return
		}
	}
}

func NewMessage(users WebSocketClientsPool, usr WebSocketClient, text string) {
	broadcast(users, usr, WebSocketMessages{
		Type: "MESSAGE",
		Content: messages.Msg{
			Created_by: usr.Id(),
			Text:       text,
			Timestamp:  time.Now(),
			MsgId:      uuid.NewString(),
		},
	})
}

func MemberLeave(users WebSocketClientsPool, usr WebSocketClient) {
	broadcast(users, usr, WebSocketMessages{
		Type: "MEMBER_LEAVE",
		Content: messages.Msg{
			Created_by: usr.Id(),
		},
	})
}

func MemberJoin(users WebSocketClientsPool, usr WebSocketClient) {
	broadcast(users, usr, WebSocketMessages{
		Type: "MEMBER_JOIN",
		Content: messages.Msg{
			Created_by: usr.Id(),
		},
	})
}

func broadcast(users WebSocketClientsPool, usr WebSocketClient, msg WebSocketMessages) {
	ForEach(Except(users, func(item WebSocketClient) bool {
		return item.Id() == usr.Id()
	}), func(item WebSocketClient) {
		item.Write(msg)
	},
	)

}
