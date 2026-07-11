package middleware

import (
	"farm-backend/internal/validation"
	"net/http"

	"github.com/gin-gonic/gin"
)

func MaxBodySize() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Body != nil && c.Request.ContentLength > validation.MaxRequestBody {
			c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{
				"error": "request body too large",
			})
			return
		}
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, validation.MaxRequestBody)
		c.Next()
	}
}