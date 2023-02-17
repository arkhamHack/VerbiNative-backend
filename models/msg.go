package models

type Message struct {
	//Time_stamp  time.Time          `json:"time_stamp,omitempty"`
	Language    string `json:"language,omitempty" validate:"required"`
	Channel     string `json:"channel,omitempty"`
	Text        string `json:"text,omitempty"`
	Translation string `json:"translation,omitempty"`
	Err         string `json:"err,omitempty"`
	Command     int    `json:"command,omitempty"`
}
