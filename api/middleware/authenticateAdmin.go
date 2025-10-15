package middleware

import (
	"github.com/gin-gonic/gin"
)

func AuthenticateAdmin(ctx *gin.Context) {
	if adm := ctx.GetBool("admin"); !adm {
		ctx.AbortWithStatusJSON(403, gin.H{"message": "Forbidden"})
		return
	}

	ctx.Next()
}
