package users

import (
	"github.com/gin-gonic/gin"
)

func UserRoute(router *gin.Engine) {
	// router.POST("/user", controllers.CreateUser())
	router.GET("/user/:username", GetUser())
	router.PUT("/user/:userId", EditUser())
	router.DELETE("/user/:userId", DeleteUser())
	router.POST("/user/signup", Signup())
	router.POST("/user/login", Login())
}
