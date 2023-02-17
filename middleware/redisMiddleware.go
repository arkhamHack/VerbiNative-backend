package middleware

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

func RedisMiddleware() gin.HandlerFunc {
	fmt.Println("Redis Connected")
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	_, err := rdb.Ping().Result()
	if err != nil {
		log.Fatal(err)
	}
	return func(c *gin.Context) {
		// add the Redis client to the context
		c.Set("redisClient", rdb)

		// call the next middleware or handler function
		c.Next()
	}
}
