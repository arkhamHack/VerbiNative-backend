package helpers

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/arkhamHack/VerbiNative-backend/configs"
	jwt "github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDetails struct {
	Email    string
	Username string
	Region   string
	Uid      string
	jwt.StandardClaims
}

var UserCollec *mongo.Collection = configs.GetCollec(configs.Mongo_DB, "users")
var secret_key = os.Getenv("SECRET_KEY")

func GenerateAllTokens(email string, username string, uid string, region string) (used_token string, used_refreshToken string, err error) {
	claims := &SignedDetails{
		Email:    email,
		Username: username,
		Uid:      uid,
		Region:   region,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}
	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret_key))
	refresh_token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(secret_key))
	if err != nil {
		log.Panic(err)
		return
	}
	return token, refresh_token, err

}

func ValidateToken(signed_token string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(
		signed_token, &SignedDetails{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret_key), nil
		},
	)

	if err != nil {
		msg = err.Error()
		return
	}
	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = err.Error()
		return
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = fmt.Sprintf("token has expired")
		msg = err.Error()
		return

	}
	return claims, msg
}
func UpdateAllTokens(signedToken string, signed_refreshToken string, user_id string) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	var update_obj primitive.D
	update_obj = append(update_obj, bson.E{"token", signedToken})
	update_obj = append(update_obj, bson.E{"refresh_token", signed_refreshToken})

	upsert := true
	filter := bson.M{"user_id": user_id}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}
	_, err := UserCollec.UpdateOne(
		ctx, filter, bson.D{
			{"$set", update_obj},
		},
		&opt,
	)
	defer cancel()
	if err != nil {
		log.Panic(err)
		return
	}
	return

}
