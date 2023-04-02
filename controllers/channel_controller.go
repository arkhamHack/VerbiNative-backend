package controllers

import (
	"errors"
	"fmt"
	"log"

	"github.com/go-redis/redis"
)

const (
	Usr_keys        = "users"
	Usr_channel_fmt = "user:%s:channels"
	Channels_key    = "channels"
)

// var rdb *redis.Client = configs.RedisConn()
type Usr_chat struct {
	// Email            string `json:"email"`
	Username         string `json:"username"`
	Channels_handler *redis.PubSub
	StopListenerChan chan struct{}
	Listening        bool
	MessageChan      chan redis.Message
	Created_by       string `json:"created_by,omitempty"`
}
type Chatroom struct {
}

func (usr *Usr_chat) Subscribe(rdb *redis.Client, channel string) error {
	usr_channels_key := fmt.Sprintf(Usr_channel_fmt, usr.Username)
	if rdb.SIsMember(usr_channels_key, channel).Val() {
		return nil
	}

	if err := rdb.SAdd(usr_channels_key, channel).Err(); err != nil {
		return err
	}
	return usr.connect(rdb)
}

func (usr *Usr_chat) Unsubscribe(rdb *redis.Client, channel string) error {
	usr_channel_key := fmt.Sprintf(Usr_channel_fmt, usr.Username)
	if !rdb.SIsMember(usr_channel_key, channel).Val() {
		return nil
	}
	if err := rdb.SRem(usr_channel_key, channel).Err(); err != nil {
		return err
	}
	return usr.connect(rdb)
}

func Connect(rdb *redis.Client, username string) (*Usr_chat, error) {
	log.Println("called to add user")
	if _, err := rdb.SAdd(Usr_keys, username).Result(); err != nil {
		return nil, err
	}
	usr := &Usr_chat{
		Username:         username,
		StopListenerChan: make(chan struct{}),
		MessageChan:      make(chan redis.Message),
	}

	if err := usr.connect(rdb); err != nil {
		return nil, err
	}
	return usr, nil
}

func (u *Usr_chat) connect(rdb *redis.Client) error {
	var c []string
	c1, err := rdb.SMembers(Channels_key).Result()
	if err != nil {
		return err
	}
	c = append(c, c1...)
	c2, err := rdb.SMembers(fmt.Sprintf(Usr_channel_fmt, u.Username)).Result()
	if err != nil {
		return err
	}
	c = append(c, c2...)
	if len(c) == 0 {
		fmt.Println("no channels to connect ")
		return nil
	}
	if u.Channels_handler != nil {
		if err := u.Channels_handler.Unsubscribe(); err != nil {
			return err
		}
		if err := u.Channels_handler.Close(); err != nil {
			return err
		}
	}
	if u.Listening {
		u.StopListenerChan <- struct{}{}
	}
	return u.doConnect(rdb, c...)
}

func (usr *Usr_chat) doConnect(rdb *redis.Client, channels ...string) error {
	pub_sub := rdb.Subscribe(channels...)
	usr.Channels_handler = pub_sub
	go func() {
		usr.Listening = true
		fmt.Println("starting listener  for  user: ", usr.Username, "on channel: ")
		for {
			select {
			case msg, ok := <-pub_sub.Channel():
				if !ok {
					return
				}
				usr.MessageChan <- *msg
			case <-usr.StopListenerChan:
				fmt.Println("Stopping  listener  for user: ", usr.Username)
				return
			}
		}
	}()
	return nil
}

func (usr *Usr_chat) Disconnect() error {
	if usr.Channels_handler != nil {
		if err := usr.Channels_handler.Unsubscribe(); err != nil {
			return err
		}
		if err := usr.Channels_handler.Close(); err != nil {
			return err
		}
	}
	if usr.Listening {
		usr.StopListenerChan <- struct{}{}
	}
	close(usr.MessageChan)
	return nil

}

func Chat(rdb *redis.Client, channel string, content string) error {
	return rdb.Publish(channel, content).Err()
}

func List(rdb *redis.Client) ([]string, error) {
	return rdb.SMembers(Usr_keys).Result()
}

func GetChanList(rdb *redis.Client, username string) ([]string, error) {
	if !rdb.SIsMember(Usr_keys, username).Val() {
		return nil, errors.New("user doesn't exist")
	}
	var c []string
	c1, err := rdb.SMembers(Channels_key).Result()
	if err != nil {
		return nil, err
	}
	c = append(c, c1...)
	c2, err := rdb.SMembers(fmt.Sprintf(Usr_channel_fmt, username)).Result()
	if err != nil {
		return nil, err
	}
	c = append(c, c2...)
	return c, nil
}

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
