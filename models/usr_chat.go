package models

import "github.com/go-redis/redis"

type Usr_chat struct {
	email            string
	channels_handler *redis.PubSub
	stopListenerChan chan struct{}
	listening        bool
	MessageChan      chan redis.Message
}
