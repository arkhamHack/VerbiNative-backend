package main

import (
	"github.com/arkhamHack/VerbiNative-backend/controllers"
	"github.com/arkhamHack/VerbiNative-backend/middleware"
	"github.com/arkhamHack/VerbiNative-backend/routes"

	"github.com/gin-gonic/gin"

	"github.com/go-redis/redis"
)

func init() {
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	defer rdb.Close()
	rdb.SAdd(controllers.Channels_key, "general", "random")
}

var rdb *redis.Client

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
	//routes.UserRoute(router)
	// corsMiddleware := cors.New(cors.Options{
	// 	AllowedOrigins:   []string{"http://localhost:3000"},
	// 	AllowCredentials: true,
	// 	Debug:            true,
	// })
	//router.Use(middleware.Authentication())
	// router.Use(corsMiddleware.Handler)

	//router.Use(middleware.CORSMiddleware())
	// router.GET("/", func(ctx *gin.Context) {
	// 	ctx.JSON(200, gin.H{
	// 		"data": "GIN",
	// 	})
	// })
	router.Use(middleware.RedisMiddleware())
	routes.ChatRoutes(router)
	router.GET("/api-1", func(c *gin.Context) {

		c.JSON(200, gin.H{"success": "Access granted for api-1"})

	})

	// API-1
	router.GET("/api-2", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "Access granted for api-2"})
	})

	router.Run("localhost:8080")
	// rdb = redis.NewClient(&redis.Options{Addr: "localhost:6379"})

	// r := mux.NewRouter()

	// r.Path("/chat").Methods("GET").HandlerFunc(api.H(rdb, api.ChatHandler))
	// r.Path("/user/{user}/channels").Methods("GET").HandlerFunc(api.H(rdb, api.UserChannelHandler))
	// r.Path("/users").Methods("GET").HandlerFunc(api.H(rdb, api.UserHandler))

	// port := ":" + os.Getenv("PORT")
	// if port == ":" {
	// 	port = ":8080"
	// }
	// fmt.Println("chat service started on port", port)
	// log.Fatal(http.ListenAndServe(port, r))
}
