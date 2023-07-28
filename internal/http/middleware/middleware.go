package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestIdMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestId := uuid.New().String()
		c.Set("requestId", requestId) // set requestId in context
		c.Next()
	}
}
