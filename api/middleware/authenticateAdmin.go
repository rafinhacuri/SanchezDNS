package middleware

import (
	"github.com/gin-gonic/gin"
)

func AuthenticateAdmin(c *gin.Context) {
	admin, exists := c.Get("admin")
	if !exists || admin == nil {
		c.AbortWithStatusJSON(401, gin.H{"message": "Unauthorized"})
		return
	}

	if !admin.(bool) {
		c.AbortWithStatusJSON(403, gin.H{"message": "Forbidden"})
		return
	}

	c.Next()
}
