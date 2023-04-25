package chatroom

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/arkhamHack/VerbiNative-backend/configs"
	"github.com/arkhamHack/VerbiNative-backend/messages"
	"github.com/arkhamHack/VerbiNative-backend/responses"

	//	"github.com/arkhamHack/VerbiNative-backend/users"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var chatroomMutex sync.Mutex
var validate = validator.New()
var ChatCollec *mongo.Collection = configs.GetCollec(configs.Mongo_DB, "servers")

// initially on joining the user will be added to the servers of each major region(country based)->so those are default created chats
func CreateChatroom() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		var chatroom Chatroom
		defer cancel()
		if err := c.BindJSON(&chatroom); err != nil {
			c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		chat_stored := Chatroom{
			Chatroom_id: uuid.New().String(),
			Name:        chatroom.Name,
			User_ids:    chatroom.User_ids,
			Messages:    chatroom.Messages,
		}
		validate_err := validate.Struct(&chat_stored)
		if validate_err != nil {
			c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "validation error", Data: map[string]interface{}{"data": validate_err.Error()}})
			return
		}
		fin, err := ChatCollec.InsertOne(ctx, chat_stored)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		result := bson.M{
			"chatroom_id": chat_stored.Chatroom_id,
			"_id":         fin.InsertedID,
			"name":        chat_stored.Name,
		}
		c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": result}})

	}
}

func GetChat() gin.HandlerFunc {
	return func(c *gin.Context) {
		// session := sessions.Default(c)
		// uid := session.Get("verbinative-userid")
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		chat := c.Param("chatroomId")
		var chatr Chatroom
		defer cancel()
		filter := bson.M{"chatroom_id": chat}
		err := ChatCollec.FindOne(ctx, filter).Decode(&chatr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"chatroom": chatr}})

	}
}
func JoinChat() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		chat := c.Param("chatroomId")
		//user := c.Param("userId")
		type User struct {
			UserID string `json:"user_id"`
		}
		var user User
		var chatr Chatroom
		defer cancel()
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		if validate_err := validate.Struct(&user); validate_err != nil {
			c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"validation error": validate_err.Error()}})
			return
		}
		filter := bson.M{"chatroom_id": chat}
		update := bson.M{"$addToSet": bson.M{"user_ids": user.UserID}}
		err := ChatCollec.FindOneAndUpdate(ctx, filter, update).Decode(&chatr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"chatroom": chatr}})
	}
}

// get's all the chats associated with user
func GetAllChats() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		user := c.Param("userId")
		var chat []Chatroom
		defer cancel()
		pipeline := bson.A{
			bson.M{
				"$match": bson.M{
					"user_ids": bson.M{
						"$in": bson.A{user},
					},
				},
			},
		}
		coll, err := ChatCollec.Aggregate(ctx, pipeline)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		err = coll.All(ctx, &chat)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"chatroom": chat}})

	}
}

// user wants to leave chat
func LeaveChat() {}

//update chat records

func UpdateChat(chatroomId string, message messages.Msg) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	var chatr Chatroom
	defer cancel()

	filter := bson.M{"chatroom_id": chatroomId}
	update := bson.M{}
	update["$addToSet"] = bson.M{"messages": message}
	// update["$setOnInsert"] = bson.M{"messages": []interface{}{}}
	// opts := options.FindOneAndUpdate().SetUpsert(true)
	err := ChatCollec.FindOneAndUpdate(ctx, filter, update).Decode(&chatr)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println(chatr)

	return nil
}

func GetMessages() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		chatroom := c.Param("chatroomId")
		var chat Chatroom

		defer cancel()
		skip := 0
		limit := 0
		if s, err := strconv.Atoi(c.Query("skip")); err == nil {
			skip = s
		}
		if l, err := strconv.Atoi(c.Query("limit")); err == nil {
			limit = l
		}
		filter := bson.M{"chatroom_id": chatroom}
		err := ChatCollec.FindOne(ctx, filter).Decode(&chat)
		msg_val := len(chat.Messages)
		projection := bson.M{
			"messages": bson.M{

				"$slice": []interface{}{bson.M{"$reverseArray": "$messages"}, skip, limit}},
		}
		//bson.M{"$subtract": []interface{}{bson.M{"$size": "$messages"}, skip}}
		opts := options.FindOne().SetProjection(projection)
		err = ChatCollec.FindOne(ctx, filter, opts).Decode(&chat)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"couldn't find chatroom": err.Error()}})
			return
		}

		c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"messages": chat.Messages, "msg_len": msg_val}})

	}
}
