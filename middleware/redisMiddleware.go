package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

func RedisMiddleware(rdb *redis.Client) gin.HandlerFunc {

	return func(c *gin.Context) {
		// add the Redis client to the context
		c.Set("redisClient", rdb)

		// call the next middleware or handler function
		c.Next()
	}
}
