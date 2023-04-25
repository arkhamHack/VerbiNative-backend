package messages

import (
	"time"
)

type Msg struct {
	// Chatroom_id string    `json:"chatroom_id,omitempty" validate:"required"`
	Created_by  string    `json:"created_by" validate:"required"`
	Username    string    `json:"username" validate:"required"`
	Text        string    `json:"text,omitempty" validate:"required"`
	Timestamp   time.Time `json:"timestamp,omitempty" validate:"required"`
	Translation string    `json:"translation,omitempty"`
	MsgId       string    `json:"msgId" validate:"required"`
	Member_id   string    `json:"member_id,omitempty"`
}
