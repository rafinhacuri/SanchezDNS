package controllers

import (
	"context"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rafinhacuri/SanchezDNS/models"
)

func Auth(ctx *gin.Context) {
	var user models.Auth
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(400, gin.H{"message": "Invalid request payload"})
		return
	}

	if err := user.Validate(); err != nil {
		slog.Error("Validation error", "error", err)
		ctx.JSON(400, gin.H{"message": "Username or password incorrect"})
		return
	}

	ctxReq, cancel := context.WithTimeout(ctx.Request.Context(), 5*time.Second)
	defer cancel()

	token, isAdmin, err := user.Login(ctxReq)
	if err != nil {
		ctx.JSON(401, gin.H{"message": "Username or password incorrect"})
		return
	}

	ctx.JSON(200, gin.H{"message": "Login successful", "token": token, "isAdmin": isAdmin})
}
