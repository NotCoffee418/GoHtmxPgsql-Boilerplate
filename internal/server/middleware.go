package server

import (
	"github.com/gin-gonic/gin"
)

func internalServerErrorHandlingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Process all errors
		for _, e := range c.Errors {
			if e.Type == gin.ErrorTypePrivate {
				switch e.Err.Error() {
				case "Internal Server Error":
					c.String(500, "Internal Server Error")
				case "Forbidden":
					c.String(403, "Forbidden")
				case "Bad Request":
					c.String(400, "Bad Request")
				}
				return
			}
		}
	}
}
