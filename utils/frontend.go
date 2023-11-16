package utils

import "github.com/gin-gonic/gin"

func IsHTMXRequest(c *gin.Context) bool {
	return c.GetHeader("HX-Request") == "true"
}
