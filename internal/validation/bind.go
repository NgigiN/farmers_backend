package validation

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Sanitizer interface {
	Sanitize()
}

func BindAndValidate(c *gin.Context, req any) error {
	if err := c.ShouldBindJSON(req); err != nil {
		return err
	}
	if sanitizer, ok := req.(Sanitizer); ok {
		sanitizer.Sanitize()
	}
	return ValidateStruct(req)
}

func RespondBindingError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{"error": FormatError(err)})
}

func RespondError(c *gin.Context, err error) {
	if IsClientError(err) {
		c.JSON(http.StatusBadRequest, gin.H{"error": FormatError(err)})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": "An unexpected error occurred"})
}

func RespondNotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, gin.H{"error": message})
}