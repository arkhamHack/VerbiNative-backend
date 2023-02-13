package middleware

import (
	"net/http"

	"github.com/arkhamHack/VerbiNative-backend/helpers"
	"github.com/arkhamHack/VerbiNative-backend/responses"

	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": "No Authorization header provided"}})
			c.Abort()
			return
		}

		claims, err := helpers.ValidateToken(clientToken)
		if err != "" {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err}})
			c.Abort()
			return
		}
		c.Set("email", claims.Email)
		c.Set("username", claims.Username)
		c.Set("region", claims.Region)
		c.Set("uid", claims.Uid)
		c.Next()
	}
}
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")

		// Skip preflight request (OPTIONS)
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		// Continue to the next middleware or endpoint handler

		c.Next()
	}
}
