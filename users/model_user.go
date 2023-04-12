package users

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id            primitive.ObjectID `bson:"_id"`
	User_id       string             `json:"user_id"`
	Token         string             `json:"token"`
	Refresh_token string             `json:"refresh_token"`
	Username      string             `json:"username,omitempty" validate:"required"`
	Region        string             `json:"region,omitempty" validate:"required"`
	Language      string             `json:"language,omitempty",validate:"required"`
	Email         string             `json:"email,omitempty" validate:"required"`
	Password      string             `json:"password,omitempty"  validate:"required"`
	Created_at    time.Time          `json:"created_at"`
	// Updated_at    time.Time          `json:"updated_at"`
}
