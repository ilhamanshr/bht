package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger returns a middleware that logs request details
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()

		log.Printf("[%s] %s %s | %d | %v",
			time.Now().Format("2006-01-02 15:04:05"),
			method,
			path,
			statusCode,
			latency,
		)

		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				log.Printf("[ERROR] %s", e.Error())
			}
		}
	}
}
