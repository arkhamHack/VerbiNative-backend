package chatroom

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/arkhamHack/VerbiNative-backend/configs"
	"github.com/arkhamHack/VerbiNative-backend/messages"
	"github.com/arkhamHack/VerbiNative-backend/responses"
	"github.com/gin-contrib/sessions"

	//	"github.com/arkhamHack/VerbiNative-backend/users"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var chatroomMutex sync.Mutex
var validate = validator.New()
var ChatCollec *mongo.Collection = configs.GetCollec(configs.Mongo_DB, "servers")
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type ChatroomHandler struct {
	Chatrooms map[string]*Chatroom
}

func NewChatroomManager() *ChatroomHandler {
	return &ChatroomHandler{
		Chatrooms: make(map[string]*Chatroom),
	}

}

// func ChatMessenger() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		fmt.Print("func called")
// 		chatroomId := c.Param("chatroomId")
// 		userId := c.Param("userId")
// 		chatroomHandler := &ChatroomHandler{}
// 		chatroomHandler.HandleChatrooms(c.Writer, c.Request, chatroomId, userId)
// 	}
// }

// func (ch *ChatroomHandler) HandleChatrooms(w http.ResponseWriter, r *http.Request, chatroomId string, userId string) {
// 	ws, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Println("Failed to upgrade to WebSocket:", err)
// 		return
// 	}
// 	if _, ok := ch.Chatrooms[chatroomId]; !ok {
// 		ch.Chatrooms[chatroomId] = &Chatroom{
// 			Chatroom_id: chatroomId,
// 			User_ids:    []string{},
// 			//Messages:       []messages.Msg{},
// 			Broadcast:  make(chan messages.Msg),
// 			Register:   make(chan *websocket.Conn),
// 			Unregister: make(chan *websocket.Conn),
// 		}

// 	}
// 	ch.Chatrooms[chatroomId].Register <- ws
// 	go func() {
// 		for {
// 			_, msg, err := ws.ReadMessage()
// 			if err != nil {
// 				log.Println("Couldn't read webdocket data")
// 				break
// 			}
// 			msg_data := messages.Msg{
// 				Created_by: userId,
// 				Text:       string(msg),
// 				Timestamp:  time.Now(),
// 				MsgId:      uuid.New().String(),
// 			}
// 			ch.Chatrooms[chatroomId].Broadcast <- msg_data
// 			err = UpdateChat(chatroomId, msg_data)
// 			if err != nil {
// 				log.Println("Failed to save message to MongoDB:", err)
// 				//c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"error": err.Error()}})

// 			}
// 		}
// 	}()
// 	go func() {
// 		for {
// 			_, _, err := ws.ReadMessage()
// 			if err != nil {
// 				ch.Chatrooms[chatroomId].Unregister <- ws
// 				break
// 			}
// 		}
// 	}()

// }

// initially on joining the user will be added to the servers of each major region(country based)->so those are default created chats
func CreateChatroom() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		var chatroom ChatroomSummary
		defer cancel()
		if err := c.BindJSON(&chatroom); err != nil {
			c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		chat_stored := ChatroomSummary{
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
		session := sessions.Default(c)
		uid := session.Get("verbinative-userid")
		log.Println("\nUser id called in chat:", uid)
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
		var chatr ChatroomSummary
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
		var chat []ChatroomSummary
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
	// if err := c.BindJSON(&chat); err != nil {
	// 	c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
	// 	return
	// }
	// if validate_err := validate.Struct(&chat); validate_err != nil {
	// 	c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"validation error": validate_err.Error()}})
	// 	return
	// }
	filter := bson.M{"chatroom_id": chatroomId}
	update := bson.M{"$addToSet": bson.M{"messages": message}}
	err := ChatCollec.FindOneAndUpdate(ctx, filter, update).Decode(&chatr)
	if err != nil {
		return err
	}
	fmt.Println(chatr)
	return nil
}

func UpdateMongo(chatroomId string, m messages.Msg) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	var message messages.Msg
	defer cancel()
	filter := bson.M{"chatroom_id": chatroomId}
	update := bson.M{"$addToSet": bson.M{"messages": "hi"}}
	err := ChatCollec.FindOneAndUpdate(ctx, filter, update).Decode(&message)
	if err != nil {
		return err
	}
	return nil
}
