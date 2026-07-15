package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"time"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		xRequestID := uuid.New().String()

		c.Set("requestid", xRequestID)

		fmt.Printf("[GIN-debug] %s [%s] - \"%s %s\"\n", time.Now().Format(time.RFC3339), xRequestID, c.Request.Method, c.Request.URL.Path)

		c.Next()
	}
}
