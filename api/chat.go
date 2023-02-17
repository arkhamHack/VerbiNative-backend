package api

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/arkhamHack/VerbiNative-backend/controllers"
	"github.com/arkhamHack/VerbiNative-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
)

var upgrader websocket.Upgrader

// var clients = make(map[*websocket.Conn]bool)
// var broadcaster = make(chan models.Message)

//	var upgrader = websocket.Upgrader{
//		CheckOrigin: func(r *http.Request) bool {
//			return true
//		},
//	}
const (
	cmd_sub = iota
	cmd_unsub
	cmd_chat
)

var connected_users = make(map[string]*controllers.Usr_chat)

func H(rdb *redis.Client, fn func(http.ResponseWriter, *http.Request, *redis.Client)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, rdb)
	}
}

// func H(rdb *redis.Client, handler gin.HandlerFunc) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		handler(c)
// 	}
// }

func webSockErrHandle(err error, ws *websocket.Conn) {
	_ = ws.WriteJSON(models.Message{Err: err.Error()})
}

func ChatHandler(c *gin.Context) {
	log.Println("Func Chat Called")
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	web_sock, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		webSockErrHandle(err, web_sock)
		return
	}
	err = Connect(c, web_sock, rdb)
	if err != nil {
		webSockErrHandle(err, web_sock)
		return
	}
	close_chan := Disconnect(c, rdb, web_sock)
	ChanMsg(web_sock, c)
loop:
	for {
		select {
		case <-close_chan:
			break loop
		default:
			UserMsg(c, rdb, web_sock)
		}
	}
}

// func ChatHandler(w http.ResponseWriter, r *http.Request, rdb *redis.Client) {
// 	web_sock, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		webSockErrHandle(err, web_sock)
// 		return
// 	}
// 	err = Connect(r, web_sock, rdb)
// 	if err != nil {
// 		webSockErrHandle(err, web_sock)
// 		return
// 	}
// 	close_chan := Disconnect(r, rdb, web_sock)
// 	ChanMsg(web_sock, r)
// loop:
// 	for {
// 		select {
// 		case <-close_chan:
// 			break loop
// 		default:
// 			UserMsg(r, rdb, web_sock)
// 		}
// 	}
// }

// func Connect(r *http.Request, web_sock *websocket.Conn, rdb *redis.Client) error {

//		username := r.URL.Query()["username"][0]
//		fmt.Println("connnected from ", web_sock.RemoteAddr(), "user:", username)
//		usr, err := controllers.Connect(rdb, username)
//		if err != nil {
//			return err
//		}
//		connected_users[username] = usr
//		return nil
//	}
func Connect(c *gin.Context, web_sock *websocket.Conn, rdb *redis.Client) error {
	log.Println("Conection attempted")
	username := c.Param("username")
	fmt.Println("connected from", web_sock.RemoteAddr(), "user:", username)
	usr, err := controllers.Connect(rdb, username)
	if err != nil {
		return err
	}
	connected_users[username] = usr
	return nil
}

func Disconnect(c *gin.Context, rdb *redis.Client, ws *websocket.Conn) chan struct{} {
	log.Println("Disconnect attempted")

	chan_close := make(chan struct{})
	username := c.Param("username")
	ws.SetCloseHandler(func(code int, text string) error {
		fmt.Println("connection closed for user: ", username)
		u := connected_users[username]
		if err := u.Disconnect(); err != nil {
			return err
		}
		delete(connected_users, username)
		close(chan_close)
		return nil
	})
	return chan_close
}

func UserMsg(c *gin.Context, rdb *redis.Client, ws *websocket.Conn) {
	log.Println("Readin bitch")
	var msg models.Message
	if err := ws.ReadJSON(&msg); err != nil {
		webSockErrHandle(err, ws)
		return
	}
	username := c.Param("username")

	u := connected_users[username]
	switch msg.Command {
	case cmd_sub:
		if err := u.Subscribe(rdb, msg.Channel); err != nil {
			webSockErrHandle(err, ws)

		}
	case cmd_unsub:
		if err := u.Unsubscribe(rdb, msg.Channel); err != nil {
			webSockErrHandle(err, ws)
		}
	case cmd_chat:
		if err := controllers.Chat(rdb, msg.Channel, msg.Text); err != nil {
			webSockErrHandle(err, ws)
		}
	}

}
func lang_conv(text string) string {
	langpair := "auto-en" // source and target languages

	url := fmt.Sprintf("https://script.google.com/macros/s/AKfycbzWv40r3aHy79g9T-6PrvkfipALo_4UgW-oF1y2827JWfa4qlKPlSg0cFK96ybRjB6Qog/exec", text, langpair)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("X-RapidAPI-Key", "SIGN-UP-FOR-KEY")
	req.Header.Add("X-RapidAPI-Host", "petapro-translate-v1.p.rapidapi.com")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	return string(body)
}

func ChanMsg(ws *websocket.Conn, c *gin.Context) {
	username := c.Param("username")
	usr := connected_users[username]
	go func() {
		for msgch := range usr.MessageChan {
			msg := models.Message{
				Text:    msgch.Payload,
				Channel: msgch.Channel,
				//Translation: lang_conv(msgch.Payload),
			}
			if err := ws.WriteJSON(msg); err != nil {
				fmt.Println(err)
			}
		}

	}()
}
