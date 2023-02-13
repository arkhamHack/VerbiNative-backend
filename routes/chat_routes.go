package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

var rdb *redis.Client

func ChatRoutes(router *gin.Engine) {
	router.GET("/chat", func(c *gin.Context) {
		// api.H(rdb, api.ChatWe)
	})
}
