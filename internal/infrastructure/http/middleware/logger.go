package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		log.Printf("[REQUEST START] %s %s - Client: %s",
			c.Request.Method,
			c.Request.URL.Path,
			c.ClientIP(),
		)

		c.Next()

		duration := time.Since(startTime)
		statusCode := c.Writer.Status()

		if statusCode >= 400 {
			log.Printf("[REQUEST ERROR] %s %s - Status: %d - Duration: %v - Client: %s",
				c.Request.Method,
				c.Request.URL.Path,
				statusCode,
				duration,
				c.ClientIP(),
			)
		} else {
			log.Printf("[REQUEST END] %s %s - Status: %d - Duration: %v - Client: %s",
				c.Request.Method,
				c.Request.URL.Path,
				statusCode,
				duration,
				c.ClientIP(),
			)
		}
	}
}
