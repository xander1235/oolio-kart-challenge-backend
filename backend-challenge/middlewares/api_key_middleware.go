package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"oolio.com/kart/configs"
	"oolio.com/kart/dtos/responses"
)

// APIKeyMiddleware checks if the API key is present in the request header
func APIKeyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("api_key")
		if apiKey != "" && apiKey == configs.APIKey {
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, responses.APIResponse{
				Code:    http.StatusUnauthorized,
				Type:    "error",
				Message: "Unauthorized",
			})
		}
	}
}
