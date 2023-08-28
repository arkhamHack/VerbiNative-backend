package main

import (
	"log"
	"os"

	"github.com/arkhamHack/VerbiNative-backend/users"
	"github.com/arkhamHack/VerbiNative-backend/websockets"

	"github.com/arkhamHack/VerbiNative-backend/chatroom"
	"github.com/arkhamHack/VerbiNative-backend/middleware"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"github.com/go-redis/redis"
)

// func init() {
// 	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
// 	defer rdb.Close()
// 	rdb.SAdd(controllers.Channels_key, "general", "random")
// }

var rdb *redis.Client

func main() {
	router := gin.Default()
	sessionKey := os.Getenv("SECRET_SESSION_KEY")
	if sessionKey == "" {
		log.Fatal("error: set SECRET_SESSION_KEY to a secret string and try again")
	}

	router.Use(gin.Logger())
	router.Use(middleware.CORSMiddleware())
	store := cookie.NewStore([]byte(sessionKey))

	router.Use(sessions.Sessions("verbinative-user-session", store))

	//router.Use(middleware.RedisMiddleware())

	users.UserRoute(router)
	authRoutes := router.Group("/user")
	authRoutes.Use(middleware.Authentication())

	chatroom.ChatRoutes(router)
	router.Group("/chat")
	// chatRouter.Use(middleware.RedisMiddleware(rdb))
	websockets.WebSockRoute(router)
	router.Group("/ws")
	
	port:=os.Getenv("PORT")
	
	address := fmt.Sprintf(":%s", port)
	router.Run(address)

}
