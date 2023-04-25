package websockets

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/arkhamHack/VerbiNative-backend/chatroom"
	"github.com/arkhamHack/VerbiNative-backend/configs"
	"github.com/arkhamHack/VerbiNative-backend/messages"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/mongo"
)

var m sync.Mutex
var ChatCollec *mongo.Collection = configs.GetCollec(configs.Mongo_DB, "servers")

func NewWebSocketClient(ws *websocket.Conn) WebSocketClient {
	return &webSocketClient{
		id:   uuid.NewString(),
		ws:   ws,
		msgs: make(chan WebSocketMessages),
		err:  make(chan error),
		done: make(chan interface{}),
	}
}

func StartClient(ctx context.Context, ws *websocket.Conn, chatroomId string) {
	usr := NewWebSocketClient(ws)
	mut.Lock()
	userspool = append(userspool, usr)
	mut.Unlock()
	defer func() {
		if err := recover(); err != nil {
			log.Printf("error:%v", err)
		}
		mut.Lock()
		defer mut.Unlock()
		userspool = Except(userspool, func(item WebSocketClient) bool {
			return item.Id() == usr.Id()
		})
		usr.Close()
	}()
	usr.Launch(ctx)
	MemberJoin(userspool, usr)
	for {
		select {
		case msg, ok := <-usr.Listen():
			if !ok {
				return
			} else {
				switch msg.Type {
				case "MESSAGE":
					NewMessage(userspool, usr, chatroomId, msg.Content.Created_by, msg.Content.Username, msg.Content.Text)
				default:
					log.Printf("unknown message type: %s", msg.Type)
					return

				}
			}
		case err := <-usr.Error():
			log.Printf("web socket error: %v", err)
		case <-usr.Done():
			MemberLeave(userspool, usr)
			return
		}
	}
}

func NewMessage(users WebSocketClientsPool, usr WebSocketClient, chatroomId string, userid string, username string, text string) {
	broadcast(users, usr, WebSocketMessages{
		Type: "MESSAGE",
		Content: messages.Msg{
			Created_by: userid,
			Username:   username,
			Text:       text,
			MsgId:      uuid.NewString(),
		},
	})
	err := chatroom.UpdateChat(chatroomId, messages.Msg{
		Created_by: userid,
		Username:   username,
		Text:       text,
		MsgId:      uuid.NewString(),
		Timestamp:  time.Now(),
	})
	if err != nil {
		log.Println("Error while uploading messages to database: ", err)
	}
}

func MemberLeave(users WebSocketClientsPool, usr WebSocketClient) {
	broadcast(users, usr, WebSocketMessages{
		Type: "MEMBER_LEAVE",
		Content: messages.Msg{
			Member_id: usr.Id(),
		},
	})
}

func MemberJoin(users WebSocketClientsPool, usr WebSocketClient) {
	broadcast(users, usr, WebSocketMessages{
		Type: "MEMBER_JOIN",
		Content: messages.Msg{
			Member_id: usr.Id(),
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
