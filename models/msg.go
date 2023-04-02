package models

import "time"

type Message struct {
	//Time_stamp  time.Time          `json:"time_stamp,omitempty"`
	Language    string `json:"language,omitempty" validate:"required"`
	Channel     string `json:"channel,omitempty"`
	Text        string `json:"text,omitempty"`
	Translation string `json:"translation,omitempty"`
	Err         string `json:"err,omitempty"`
	Command     int    `json:"command,omitempty"`
	Username    string `json:"username,omitempty" validate:"required"`
}

type Msg struct {
	// Chatroom_id string    `json:"chatroom_id,omitempty" validate:"required"`
	Created_by  string    `json:"created_by,omitempty" validate:"required"`
	Text        string    `json:"text,omitempty" validate:"required"`
	Timestamp   time.Time `json:"timestamp,omitempty" validate:"required"`
	Translation string    `json:"translation,omitempty" validate:"required"`
	MsgId       string    `json:"msgId,omitempty" validate:"required"`
}

type Chatroom struct {
	Chatroom_id string `json:"chatroom_id,omitempty" validate:"required"`
	Messages    []Msg
}
