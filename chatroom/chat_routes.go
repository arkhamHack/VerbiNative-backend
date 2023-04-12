package chatroom

import (
	"github.com/gin-gonic/gin"
)

func ChatRoutes(router *gin.Engine) {
	//	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

	router.GET("/chat/:userId", GetAllChats())
	router.POST("/chat/:userId/", CreateChatroom())
	router.PATCH("/chat/:userId/join/:chatroomId", JoinChat())
	router.GET("/chat/:chatroomId", GetChat())
	router.GET("/ws/:userId/:chatroomId", ChatMessenger())
	//router.PATCH("/:username/chats/:chatroom_id")
}
