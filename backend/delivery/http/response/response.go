package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Meta struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

func JSON(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, gin.H{
		"success": statusCode >= 200 && statusCode < 300,
		"data":    data,
	})
}

func JSONWithMeta(c *gin.Context, statusCode int, data interface{}, meta *Meta) {
	c.JSON(statusCode, gin.H{
		"success": statusCode >= 200 && statusCode < 300,
		"data":    data,
		"meta":    meta,
	})
}

func Error(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{
		"success": false,
		"error":   message,
	})
}

func Unauthorized(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"success": false,
		"error":   message,
	})
}
