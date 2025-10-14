package controllers

import "github.com/gin-gonic/gin"

func Logout(ctx *gin.Context) {
	ctx.SetCookie("session", "", -1, "/", "", false, true)
	ctx.JSON(200, gin.H{"message": "Logout successful"})
}
