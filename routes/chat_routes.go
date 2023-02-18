package routes

import (
	"github.com/arkhamHack/VerbiNative-backend/api"
	"github.com/gin-gonic/gin"
)

func ChatRoutes(router *gin.Engine) {
	//	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

	router.GET("/chat/:username", api.ChatHandler)
	router.GET("/chat/:username/channels", api.UserChannelHandler)
	router.GET("/chat/users", api.UserHandler)

}
