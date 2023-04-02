package api

import (
	"encoding/json"
	"net/http"

	"github.com/arkhamHack/VerbiNative-backend/controllers"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

func Errhandle(err error, c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
}

func UserChannelHandler(c *gin.Context) {
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	username := c.Param("user")
	list, err := controllers.GetChanList(rdb, username)
	if err != nil {
		Errhandle(err, c)
		return
	}
	c.JSON(http.StatusOK, list)
}

// func UserHandler(w http.ResponseWriter, r *http.Request, rdb *redis.Client) {
// 	list, err := controllers.List(rdb)
// 	if err != nil {
// 		Errhandle(err, w)
// 		return
// 	}
// 	err = json.NewEncoder(w).Encode(list)
// 	if err != nil {
// 		Errhandle(err, w)
// 		return
// 	}
// }

func UserHandler(c *gin.Context) {

	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	list, err := controllers.List(rdb)
	if err != nil {
		Errhandle(err, c)
		return
	}
	err = json.NewEncoder(c.Writer).Encode(list)
	if err != nil {
		Errhandle(err, c)
		return
	}
}
