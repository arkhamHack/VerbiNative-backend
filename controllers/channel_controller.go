package controllers

import (
	"errors"
	"fmt"
	"log"

	"github.com/go-redis/redis"
)

const (
	Usr_keys        = "users"
	Usr_channel_fmt = "user:%s:channels"
	Channels_key    = "channels"
)

// var rdb *redis.Client = configs.RedisConn()
type Usr_chat struct {
	// Email            string `json:"email"`
	Username         string `json:"username"`
	Channels_handler *redis.PubSub
	StopListenerChan chan struct{}
	Listening        bool
	MessageChan      chan redis.Message
	Created_by       string `json:"created_by,omitempty"`
}
type Chatroom struct {
}

func (usr *Usr_chat) Subscribe(rdb *redis.Client, channel string) error {
	usr_channels_key := fmt.Sprintf(Usr_channel_fmt, usr.Username)
	if rdb.SIsMember(usr_channels_key, channel).Val() {
		return nil
	}

	if err := rdb.SAdd(usr_channels_key, channel).Err(); err != nil {
		return err
	}
	return usr.connect(rdb)
}

func (usr *Usr_chat) Unsubscribe(rdb *redis.Client, channel string) error {
	usr_channel_key := fmt.Sprintf(Usr_channel_fmt, usr.Username)
	if !rdb.SIsMember(usr_channel_key, channel).Val() {
		return nil
	}
	if err := rdb.SRem(usr_channel_key, channel).Err(); err != nil {
		return err
	}
	return usr.connect(rdb)
}

func Connect(rdb *redis.Client, username string) (*Usr_chat, error) {
	log.Println("called to add user")
	if _, err := rdb.SAdd(Usr_keys, username).Result(); err != nil {
		return nil, err
	}
	usr := &Usr_chat{
		Username:         username,
		StopListenerChan: make(chan struct{}),
		MessageChan:      make(chan redis.Message),
	}

	if err := usr.connect(rdb); err != nil {
		return nil, err
	}
	return usr, nil
}

func (u *Usr_chat) connect(rdb *redis.Client) error {
	var c []string
	c1, err := rdb.SMembers(Channels_key).Result()
	if err != nil {
		return err
	}
	c = append(c, c1...)
	c2, err := rdb.SMembers(fmt.Sprintf(Usr_channel_fmt, u.Username)).Result()
	if err != nil {
		return err
	}
	c = append(c, c2...)
	if len(c) == 0 {
		fmt.Println("no channels to connect ")
		return nil
	}
	if u.Channels_handler != nil {
		if err := u.Channels_handler.Unsubscribe(); err != nil {
			return err
		}
		if err := u.Channels_handler.Close(); err != nil {
			return err
		}
	}
	if u.Listening {
		u.StopListenerChan <- struct{}{}
	}
	return u.doConnect(rdb, c...)
}

func (usr *Usr_chat) doConnect(rdb *redis.Client, channels ...string) error {
	pub_sub := rdb.Subscribe(channels...)
	usr.Channels_handler = pub_sub
	go func() {
		usr.Listening = true
		fmt.Println("starting listener  for  user: ", usr.Username, "on channel: ")
		for {
			select {
			case msg, ok := <-pub_sub.Channel():
				if !ok {
					return
				}
				usr.MessageChan <- *msg
			case <-usr.StopListenerChan:
				fmt.Println("Stopping  listener  for user: ", usr.Username)
				return
			}
		}
	}()
	return nil
}

func (usr *Usr_chat) Disconnect() error {
	if usr.Channels_handler != nil {
		if err := usr.Channels_handler.Unsubscribe(); err != nil {
			return err
		}
		if err := usr.Channels_handler.Close(); err != nil {
			return err
		}
	}
	if usr.Listening {
		usr.StopListenerChan <- struct{}{}
	}
	close(usr.MessageChan)
	return nil

}

func Chat(rdb *redis.Client, channel string, content string) error {
	return rdb.Publish(channel, content).Err()
}

func List(rdb *redis.Client) ([]string, error) {
	return rdb.SMembers(Usr_keys).Result()
}

func GetChanList(rdb *redis.Client, username string) ([]string, error) {
	if !rdb.SIsMember(Usr_keys, username).Val() {
		return nil, errors.New("user doesn't exist")
	}
	var c []string
	c1, err := rdb.SMembers(Channels_key).Result()
	if err != nil {
		return nil, err
	}
	c = append(c, c1...)
	c2, err := rdb.SMembers(fmt.Sprintf(Usr_channel_fmt, username)).Result()
	if err != nil {
		return nil, err
	}
	c = append(c, c2...)
	return c, nil
}
