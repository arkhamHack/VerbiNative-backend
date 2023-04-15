package video

type Video struct {
	User_ids []string `json:"user_ids"validate:"required"`
}

type ChatroomVideo struct {
	Chatroom_id string   `json:"chatroom_id"validate:"required"`
	User_ids    []string `json:"user_ids"validate:"required"`
}
