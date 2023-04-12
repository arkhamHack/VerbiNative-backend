package api

import (
	"encoding/json"
	"net/http"

	"github.com/arkhamHack/VerbiNative-backend/controllers"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

func Errhandle(err error, c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
}

func UserChannelHandler(c *gin.Context) {
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	username := c.Param("user")
	list, err := controllers.GetChanList(rdb, username)
	if err != nil {
		Errhandle(err, c)
		return
	}
	c.JSON(http.StatusOK, list)
}

func UserHandler(c *gin.Context) {

	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	list, err := controllers.List(rdb)
	if err != nil {
		Errhandle(err, c)
		return
	}
	err = json.NewEncoder(c.Writer).Encode(list)
	if err != nil {
		Errhandle(err, c)
		return
	}
}

// func ChatroomHandler(c *gin.Context) {
// 	upgrader := websocket.Upgrader{WriteBufferSize: 1024, ReadBufferSize: 1024, CheckOrigin: func(r *http.Request) bool {
// 		return true
// 	}}
// 	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
// 	chatroomId := c.Param("chatroom_id")
// 	go ChatroomHandler.HandleChatrooms(ws, chatroomId)

// }
