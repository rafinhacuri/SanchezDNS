package middleware

import "github.com/gin-gonic/gin"

func AdminOnly(c *gin.Context) {
	level := c.GetString("level")

	if level != "admin" {
		c.AbortWithStatusJSON(403, gin.H{"message": "Forbidden"})

		return
	}

	c.Next()
}
