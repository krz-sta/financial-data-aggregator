package middleware

import "github.com/gin-gonic/gin"

func AuthMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader = c.GetHeader("Authorization")

	}
}
