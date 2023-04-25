package chatroom

import (
	"github.com/arkhamHack/VerbiNative-backend/messages"
)

type Chatroom struct {
	Name        string          `json:"name" validate:"required"`
	Chatroom_id string          `json:"chatroom_id" validate:"required"`
	Messages    []messages.Msg  `json:"messages,omitempty" validate:"required"`
	Msg_len     int             `json:"msg_len,omitempty" `
	User        []User_Identify `json:"user" validate:"required"`

	// User        []map[string]string `json:"user_ids" validate:"required"` //map from userid to username
	// Type        int                 `json:"type" validate:"required"`     //type defines type of chatroom: 0->Public Chatroom(official),1->Dm,2->Public Chatroom(unofficial)
}
type User_Identify struct {
	User_id   string `json:"user_id" validate:"required"`
	Usernames string `json:"username" validate:"required"`
}
