package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	Id          primitive.ObjectID `bson:"_id"`
	Time_stamp  time.Time          `json:"time_stamp"`
	Language    string             `json:"language,omitempty" validate:"required"`
	Chat_id     string             `json:"chat_id"`
	Text        string             `json:"text,omitempty" validate:"required"`
	Translation string             `json:"translation,omitempty" validate:"required"`
}
