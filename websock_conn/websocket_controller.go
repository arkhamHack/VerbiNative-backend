package websockets

import (
	"context"
	"errors"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var (
	mut       sync.Mutex
	userspool WebSocketClientsPool
)

func WebSocketConnection() gin.HandlerFunc {
	return func(c *gin.Context) {
		//chatroomId := c.Param("chatroomId")
		//chatroom:=c.Param("chatroomId")
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		//session, err := users.CookieStorage().Get(c.Request, "verbinative-user-session")
		if err != nil {
			log.Fatalln(err)
		}
		//uid := session.Values["userId"]

		go StartClient(c, ws, "123")
	}
}
func (c *webSocketClient) Launch(ctx context.Context) {
	c.ws.SetReadLimit(MaxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(PongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(PongWait)); return nil })
	c.once.Do(func() { go c.launchSupport(ctx) })
}

func (c *webSocketClient) launchSupport(ctx context.Context) {
	var wg sync.WaitGroup
	cancellation, cancel := context.WithCancel(ctx)
	defer func() {
		cancel()
		c.Send(websocket.CloseMessage)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		c.Read(cancellation)
		cancel()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		c.Ping(cancellation)
		cancel()
	}()
	wg.Wait()
	c.done <- struct{}{}
}
func (c *webSocketClient) Read(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			var msg WebSocketMessages
			err := c.ws.ReadJSON(&msg)
			if err != nil {
				c.HandleError(err)
				return
			}
			c.msgs <- msg
		}
	}
}
func (c *webSocketClient) Ping(ctx context.Context) {
	timer := time.NewTicker(PingInterval)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			c.Send(websocket.PingMessage)
		case <-ctx.Done():
			return
		}
	}
}
func (c *webSocketClient) Id() string {
	return c.id
}
func (c *webSocketClient) HandleError(err error) {
	if _, ok := err.(*websocket.CloseError); ok {
		return
	}
	if errors.Is(err, websocket.ErrCloseSent) {
		return
	}
	c.err <- err
}
func (c *webSocketClient) Close() error {
	close(c.msgs)
	return c.ws.Close()
}

func (c *webSocketClient) Listen() <-chan WebSocketMessages {
	return c.msgs
}

func (c *webSocketClient) Done() <-chan interface{} {
	return c.done
}

func (c *webSocketClient) Error() <-chan error {
	return c.err
}
func (c *webSocketClient) Write(m WebSocketMessages) error {
	c.mutex.Lock()
	defer c.mutex.Lock()
	c.ws.SetWriteDeadline(time.Now().Add(WriteWait))
	return c.ws.WriteJSON(m)
}

func (c *webSocketClient) Send(msg_type int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.ws.SetWriteDeadline(time.Now().Add(WriteWait))
	if err := c.ws.WriteMessage(msg_type, nil); err != nil {
		c.HandleError(err)
	}
}
