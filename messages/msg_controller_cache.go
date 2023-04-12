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

// func FirebaseConn(user_id string, txt string, translation string, timestamp time.Time) err.Error {

// 	opt := option.WithCredentialsFile("/home/kraken/Downloads/verbinative-firebase-adminsdk-pljdj-9ce1acc0d9.json")
// 	app, err := firebase.NewApp(context.Background(), nil, opt)
// 	if err != nil {
// 		return nil, fmt.Errorf("error initializing app: %v", err)
// 	}
// 	client, err := app.FireStore(context.Background())
// 	if err != nil {
// 		return nil, fmt.Errorf("error while fetching Firestore client ", err)
// 	}
// 	//query := client.Collection("messages").OrderBy("Timestamp", firestore.Descending).Limit(10)
// 	msg := models.Firebase_Msg{
// 		Created_by:  user_id,
// 		Text:        txt,
// 		Translation: translation,
// 		Timestamp:   timestamp,
// 		//MessageId:
// 	}
// 	ref, _, err := client.Collection("messages").Add(context.Background(), msg)
// 	if err != nil {
// 		log.Fatal("error adding message to Firestore: %v\n", err)
// 	}

// 	fmt.Printf("Added message with ID: %v\n", docRef.MsgId)
// }

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
