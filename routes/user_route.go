package routes

import (
	"github.com/arkhamHack/VerbiNative-backend/controllers"
	"github.com/gin-gonic/gin"
)

func UserRoute(router *gin.Engine) {
	// router.POST("/user", controllers.CreateUser())
	router.GET("/user/:username", controllers.GetUser())
	router.PUT("/user/:userId", controllers.EditUser())
	router.DELETE("/user/:userId", controllers.DeleteUser())
	router.POST("/user/signup", controllers.Signup())
	router.POST("/user/login", controllers.Login())
}
