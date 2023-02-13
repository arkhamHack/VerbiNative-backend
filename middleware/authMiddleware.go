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
