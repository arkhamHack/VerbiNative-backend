package chatroom

import (
	"github.com/arkhamHack/VerbiNative-backend/messages"
)

type Chatroom struct {
	Name        string         `json:"name" validate:"required"`
	Chatroom_id string         `json:"chatroom_id" validate:"required"`
	Messages    []messages.Msg `json:"messages,omitempty" validate:"required"`
	User_ids    []string       `json:"user_ids" validate:"required"`
	// User        []map[string]string `json:"user_ids" validate:"required"` //map from userid to username
	// Type        int                 `json:"type" validate:"required"`     //type defines type of chatroom: 0->Public Chatroom(official),1->Dm,2->Public Chatroom(unofficial)
}
