package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rafinhacuri/SanchezDNS/utils"
)

func Authenticate(ctx *gin.Context) {
	token, err := ctx.Cookie("session")
	if err != nil {
		ctx.AbortWithStatusJSON(401, gin.H{"message": "Unauthorized"})
		return
	}

	email, adm, err := utils.JWTValidate(token)
	if err != nil {
		ctx.AbortWithStatusJSON(401, gin.H{"message": "Invalid token"})
		return
	}

	ctx.Set("email", email)
	ctx.Set("admin", adm)
	ctx.Next()
}
