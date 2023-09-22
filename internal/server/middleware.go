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
				// handle errors here, like logging or sending a generic error response
				c.String(500, "Internal Server Error")
				return
			}
		}
	}
}
