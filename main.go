package main

import (
	"github.com/arkhamHack/VerbiNative-backend/middleware"
	"github.com/arkhamHack/VerbiNative-backend/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	// client,contxt,cancel,err:=connect("http://localhost:27017")
	// if err!=nil{
	// 	panic(err)
	// }

	// defer close(client,contxt,cancel)
	// ping(client,contxt)
	router := gin.New()
	router.Use(gin.Logger())
	//configs.ConnectDB()
	routes.UserRoute(router)
	router.Use(cors.Default())
	router.Use(middleware.Authentication())
	// router.Use(middleware.CORSMiddleware())
	// router.GET("/", func(ctx *gin.Context) {
	// 	ctx.JSON(200, gin.H{
	// 		"data": "GIN",
	// 	})
	// })
	router.GET("/api-1", func(c *gin.Context) {

		c.JSON(200, gin.H{"success": "Access granted for api-1"})

	})

	// API-1
	router.GET("/api-2", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "Access granted for api-2"})
	})

	router.Run("localhost:8080")
}
