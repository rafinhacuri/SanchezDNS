package middleware

import (
	"github.com/gin-gonic/gin"
)

func AuthenticateAdmin(ctx *gin.Context) {
	admin, exists := ctx.Get("admin")
	if !exists || admin == nil {
		ctx.AbortWithStatusJSON(401, gin.H{"message": "Unauthorized"})
		return
	}

	if !admin.(bool) {
		ctx.AbortWithStatusJSON(403, gin.H{"message": "Forbidden"})
		return
	}

	ctx.Next()
}
