package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rafinhacuri/SanchezDNS/utils"
)

func CheckSession(ctx *gin.Context) {
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

	ctx.JSON(200, gin.H{"username": email, "isAdmin": adm})
}
