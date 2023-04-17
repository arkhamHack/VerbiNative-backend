package messages

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/arkhamHack/VerbiNative-backend/responses"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func MsgfetchCache(ctx context.Context, chatroomId string, rdb *redis.Client) ([]*Msg, error) {
	msg_ids, err := rdb.ZRangeByScore(ctx, "chat:"+chatroomId, &redis.ZRangeBy{
		Min: "-inf",
		Max: "+inf",
	}).Result()
	if err != nil {
		return nil, err
	}
	var msgs []*Msg
	for _, msg_id := range msg_ids {
		msg_bytes, err := rdb.Get(ctx, msg_id).Bytes()
		if err != nil {
			return nil, err
		}
		var msg Msg
		err = json.Unmarshal(msg_bytes, &msg)
		if err != nil {
			return nil, err
		}
		msgs = append(msgs, &msg)
	}

	return msgs, nil

}

func storeInCache(ctx context.Context, msg *Msg, rdb *redis.Client) error {
	msg_bytes, err := json.Marshal(msg)
	if err != nil {
		return err

	}
	// ctx,cancel:=context.WithTimeout(context.Background,time.Second*10)
	// defer cancel()
	err = rdb.Set(ctx, msg.MsgId, msg_bytes, 10*time.Minute).Err()
	//err=rdb.ZAdd(ctx,msg.MessageId,&redis.Z{Score:float64(msg.Timestamp.UnixNano()),Member:msg_bytes}).Err()

	if err != nil {
		return err
	}
	return nil
}

func CreateMessage(c *gin.Context, rdb *redis.Client) {
	var msg Msg
	err := c.BindJSON(&msg)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"err": err}})
	}
	err = storeInCache(c, &msg, rdb)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"err": err}})
		return
	}
	c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": msg.Text}})
}
