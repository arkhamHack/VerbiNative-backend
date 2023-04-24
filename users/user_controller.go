package users

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/arkhamHack/VerbiNative-backend/configs"
	"github.com/arkhamHack/VerbiNative-backend/helpers"
	"github.com/arkhamHack/VerbiNative-backend/responses"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(user_pwd string, provided_pwd string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(provided_pwd), []byte(user_pwd))
	check := true
	msg := ""
	if err != nil {
		msg = fmt.Sprintf("login or passowrd is incorrect")
		check = false
	}
	return check, msg
}

var UserCollec *mongo.Collection = configs.GetCollec(configs.Mongo_DB, "users")

var validate = validator.New()

// func CookieStorage() *sessions.CookieStore {
// 	SecretKey := []byte(os.Getenv("SECRET_SESSION_KEY"))
// 	cookieStore := sessions.NewCookieStore(SecretKey)

// 	return cookieStore
// }

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		var user User
		defer cancel()
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		validate_err := validate.Struct(&user)
		if validate_err != nil {
			c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "validation error", Data: map[string]interface{}{"data": validate_err.Error()}})
			return
		}
		count, err := UserCollec.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error: email already registered", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		count2, err := UserCollec.CountDocuments(ctx, bson.M{"username": user.Username})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error: username already exists", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if count > 0 || count2 > 0 {
			c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error: this email or username already exists", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		password := HashPassword(user.Password)
		new_usr := User{
			Id:          primitive.NewObjectID(),
			Username:    user.Username,
			Email:       user.Email,
			Region:      user.Region,
			Language:    user.Language,
			Language_id: user.Language_id,
			Password:    password,
			User_id:     uuid.New().String(),
		}

		new_usr.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		token, refresh_token, _ := helpers.GenerateAllTokens(new_usr.Username, new_usr.Email, new_usr.User_id, new_usr.Region)
		new_usr.Token = token
		new_usr.Refresh_token = refresh_token
		fin, err := UserCollec.InsertOne(ctx, new_usr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})

			return
		}
		result := bson.M{
			"user_id":     new_usr.User_id,
			"_id":         fin.InsertedID,
			"email":       new_usr.Email,
			"username":    new_usr.Username,
			"language":    new_usr.Language,
			"language_id": new_usr.Language_id,
			"region":      new_usr.Region,
		}
		c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": result}})
	}
}
func MyHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set("key", "value")
		session.Save()

		c.String(http.StatusOK, "Value stored in session")
		user := session.Get("key")
		if user == nil {
			log.Println("Invalid session token")
			return
		}
		log.Println(user)
	}
}
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set("hi", "123")
		session.Save()
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		// userId := c.Param("userId")
		var user User
		var usr_found User
		defer cancel()
		// objId, _ := primitive.ObjectIDFromHex(userId)
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		err := UserCollec.FindOne(ctx, bson.M{"email": user.Email}).Decode(&usr_found)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		pwdValid, msg := VerifyPassword(user.Password, usr_found.Password)
		defer cancel()
		if pwdValid != true {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": msg}})
			return
		}
		token, refreshToken, _ := helpers.GenerateAllTokens(usr_found.Email, usr_found.Username, usr_found.User_id, usr_found.Region)
		helpers.UpdateAllTokens(token, refreshToken, usr_found.User_id)
		uid_init := usr_found.User_id
		session.Set("verbinative-userid", uid_init)
		err_sess := session.Save()
		uid := session.Get("verbinative-userid")
		log.Println("\nUser id:", uid)
		if err_sess != nil {
			// handle the error
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": usr_found}})
	}
}
func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		// secret_key := os.Getenv("SECRET_KEY")

		// token_val := c.GetHeader("Authorization")
		// if token_val == "" {
		// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
		// 	return
		// }
		// token, err := jwt.Parse(strings.Replace(token_val, "Bearer ", "", 1), func(token *jwt.Token) (interface{}, error) {
		// 	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		// 		return nil, fmt.Errorf("invalid token")
		// 	}
		// 	return []byte(secret_key), nil
		// })
		// if err != nil {
		// 	c.JSON(http.StatusUnauthorized, responses.UserResponse{Status: http.StatusUnauthorized, Message: "error: unauthorized", Data: map[string]interface{}{"data": err.Error()}})
		// 	return
		// }
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		username := c.Query("username")
		defer cancel()
		// objId, _ := primitive.ObjectIDFromHex(userId)
		// claims, ok := token.Claims.(jwt.MapClaims)
		// if !ok || claims["username"].(string) != username {
		// c.JSON(http.StatusUnauthorized, responses.UserResponse{Status: http.StatusUnauthorized, Message: "error: unauthorized", Data: map[string]interface{}{"data": err.Error()}})
		// return
		// }
		regex_pattern := bson.M{"$regex": primitive.Regex{Pattern: username, Options: "i"}}
		var matching_users []User
		filter := bson.M{"username": regex_pattern}
		coll, err := UserCollec.Find(ctx, filter)
		for coll.Next(ctx) {
			var potential_user User
			err := coll.Decode(&potential_user)
			if err != nil {
				c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err}})
				return
			}
			matching_users = append(matching_users, potential_user)
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": matching_users}})
	}
}

func EditUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		userId := c.Param("userId")
		var userUpdate map[string]interface{}
		defer cancel()
		if err := c.BindJSON(&userUpdate); err != nil {
			c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error: could'nt bind to JSON", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		// if validate_err := validate.Struct(&user); validate_err != nil {
		// 	c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"validation error": validate_err.Error()}})
		// 	return
		// }
		delete(userUpdate, "email")
		delete(userUpdate, "password")
		delete(userUpdate, "username")
		// update := bson.M{"region": user.Region, "language": user.Language, "language_id": user.Language_id, "password": user.Password}
		update := bson.M{"$set": userUpdate}
		fin, err := UserCollec.UpdateOne(ctx, bson.M{"user_id": userId}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error:couldn't find user or couldn't update", Data: map[string]interface{}{"error": err.Error()}})
			return
		}
		var updated_user User

		if fin.MatchedCount == 1 {
			err := UserCollec.FindOne(ctx, bson.M{"user_id": userId}).Decode(&updated_user)
			if err != nil {
				c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
		}
		c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": updated_user}})

	}

}

func DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		userId := c.Param("userId")

		defer cancel()
		objId, _ := primitive.ObjectIDFromHex(userId)
		fin, err := UserCollec.DeleteOne(ctx, bson.M{"_id": objId})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"error": err.Error()}})
			return
		}
		if fin.DeletedCount < 1 {
			c.JSON(http.StatusNotFound, responses.UserResponse{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "User with specified ID not found"}})
			return
		}
		c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "User deleted successfully."}})

	}
}

func GetByRegion() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		region := c.Param("region")
		var user []User
		defer cancel()
		pipeline := bson.A{
			bson.M{
				"$match": bson.M{
					"region": bson.M{
						"$in": bson.A{region},
					},
				},
			},
		}
		coll, err := UserCollec.Aggregate(ctx, pipeline)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		err = coll.All(ctx, &user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": user}})
	}
}

func GetUserDetails() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		userid := c.Param("userId")
		defer cancel()
		filter := bson.M{"user_id": userid}
		coll, err := UserCollec.Find(ctx, filter)

		defer coll.Close(context.Background())
		var result bson.M

		for coll.Next(context.Background()) {
			err := coll.Decode(&result)
			if err != nil {
				return
			}
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err}})
			return
		}
		c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": result}})

	}
}
