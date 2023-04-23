package websockets

import "github.com/gin-gonic/gin"

func WebSockRoute(router *gin.Engine) {
	router.GET("/ws/chat/:chatroomId", WebSocketConnection())
}
