package chatroom

import (
	"github.com/gin-gonic/gin"
)

func ChatRoutes(router *gin.Engine) {
	//	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

	router.GET("/chat/user/:userId/", GetAllChats())
	router.POST("/chat/create/", CreateChatroom())
	router.PATCH("/chat/join/:chatroomId/", JoinChat())
	router.GET("/chat/:chatroomId", GetChat())
	router.GET("/chat/ws/:userId/:chatroomId", ChatMessenger())
	//router.PATCH("/:username/chats/:chatroom_id")
}
