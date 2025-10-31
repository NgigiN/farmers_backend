package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func Validate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Bind handlers for global validation
		c.Next()
	}
}

func ValidateStruct(obj interface{}) error {
	return validate.Struct(obj)
}
