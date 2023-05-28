package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/arkhamHack/VerbiNative-backend/messages"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

// func RedisMiddleware() gin.HandlerFunc {
// 	rdb := redis.NewClient(&redis.Options{
// 		Addr:     "127.0.0.1:6379",
// 		Password: "", // no password set
// 		DB:       0,  // use default DB
// 	})
// 	pong, err := rdb.Ping().Result()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println("Redis Connected")
// 	fmt.Println(pong, err)

// 	return func(c *gin.Context) {
// 		// add the Redis client to the context
// 		c.Set("redisClient", rdb)

//			// call the next middleware or handler function
//			c.Next()
//		}
//	}
type RedisMsg struct {
	msg        messages.Msg
	chatroomId string
}

func RedisStore() gin.HandlerFunc {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{
			"localhost:7000",
			"localhost:7001",
		},
	})
	err := rdb.ForEachShard(ctx, func(ctx context.Context, shard *redis.Client) error {
		return shard.Ping(ctx).Err()

	})
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Redis Connected")
	}
	return func(c *gin.Context) {
		c.Set("redisClient", rdb)
		c.Next()
	}
}

func RedisPush(rdb *redis.ClusterClient, msg messages.Msg, chatroomId string) error {
	var strMsg RedisMsg
	strMsg.msg = msg
	strMsg.chatroomId = chatroomId
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	now := float64(time.Now().UnixNano())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err = rdb.ZAdd(ctx, chatroomId, &redis.Z{Score: now, Member: msgJSON}).Result()
	return err
}

func RedisPull(rdb *redis.ClusterClient, chatroomId string) ([]messages.Msg, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := rdb.ZRange(ctx, chatroomId, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	msgs := make([]messages.Msg, len(res))
	for i, msgJson := range res {
		err := json.Unmarshal([]byte(msgJson), &msgs[i])
		if err != nil {
			return nil, err
		}
	}
	return msgs, nil
}
