package users

import (
	"github.com/gin-gonic/gin"
)

func UserRoute(router *gin.Engine) {
	// router.POST("/user", controllers.CreateUser())
	router.GET("/user/search", GetUser())
	router.GET("/user/region/:region", GetByRegion())
	//router.GET("/user/language/:language", GetByRegion())
	router.PUT("/user/:userId", EditUser())
	router.GET("/user/:userId", GetUserDetails())
	router.DELETE("/user/:userId", DeleteUser())
	router.POST("/user/signup", Signup())
	router.POST("/user/login", Login())
	//router.POST("/user/test", MyHandler())

}
