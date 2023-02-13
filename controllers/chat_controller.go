package controllers

import (
	// "encoding/json"
	// "io"
	// "log"
	// "net/http"

	"github.com/arkhamHack/VerbiNative-backend/models"
	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var broadcaster = make(chan models.Message)

// var upgrader = websocket.Upgrader{
// 	CheckOrigin: func(r *http.Request) bool {
// 		return true
// 	},
// }

// //var rdb *redis.Client = configs.RedisConn()

// // func Connect(rdb *redis.Client, name string) (*models.Usr_chat, error) {
// // if _,err:=rdb.SAdd()
// // }
// func connection(w http.ResponseWriter, r *http.Request) {
// 	web_sock, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer web_sock.Close()
// 	clients[web_sock] = true
// 	// http.Handle("/", http.FileServer(http.Dir("./public")))
// 	// log.Print("Server started at localhost:4444")
// 	if rdb.Exists("chat_msg").Val() != 0 {
// 		sendPreviousMessages(web_sock)
// 	}
// 	for {
// 		var msg models.Message
// 		err := web_sock.ReadJSON(&msg)
// 		if err != nil {
// 			delete(clients, web_sock)
// 			break
// 		}
// 		broadcaster <- msg
// 	}

// }
// func HandleMessages() {
// 	for {
// 		msg := <-broadcaster
// 		RedisStore(msg)
// 		MsgClients(msg)
// 	}

// }

// func RedisStore(msg models.Message) {
// 	json, err := json.Marshal(msg)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	if err := rdb.RPush("chat_messages", json); err != nil {
// 		log.Fatal(err)
// 	}
// }

// func sendPreviousMessages(web_sock *websocket.Conn) {
// 	chat_msgs, err := rdb.LRange("chat_messages", 0, -1).Result()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	for _, chatMsg := range chat_msgs {
// 		var msg models.Message
// 		json.Unmarshal([]byte(chatMsg), &msg)
// 		MsgClient(web_sock, msg)
// 	}
// }
// func MsgClient(client *websocket.Conn, msg models.Message) {
// 	err := client.WriteJSON(msg)
// 	if err != nil && unsafeError(err) {
// 		log.Printf("error: %v", err)
// 		client.Close()
// 		delete(clients, client)
// 	}
// }

// func MsgClients(msg models.Message) {
// 	for client := range clients {
// 		MsgClient(client, msg)
// 	}
// }

// func unsafeError(err error) bool {
// 	return !websocket.IsCloseError(err, websocket.CloseGoingAway) && err != io.EOF
// }
