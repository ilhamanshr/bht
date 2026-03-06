package middleware

import (
	"bht-test/internal/domain"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// Recovery returns a middleware that recovers from any panics and writes a 500 if there was one.
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the error and stack trace
				slog.Error("Panic recovered",
					slog.Any("error", err),
					slog.String("stack", string(debug.Stack())),
				)

				// Respond with a 500 error using standardized format
				c.AbortWithStatusJSON(http.StatusInternalServerError, domain.NewErrorResponse("Internal server error"))
			}
		}()
		c.Next()
	}
}
