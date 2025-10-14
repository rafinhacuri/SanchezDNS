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
		ctx.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := user.Validate(); err != nil {
		slog.Error("Validation error", "error", err)
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctxReq, cancel := context.WithTimeout(ctx.Request.Context(), 5*time.Second)
	defer cancel()

	token, err := user.Login(ctxReq)
	if err != nil {
		slog.Error("Login error", "error", err)
		ctx.JSON(401, gin.H{"error": err.Error()})
		return
	}

	ctx.Header("Authorization", "Bearer "+token)
	ctx.Header("Access-Control-Expose-Headers", "Authorization")

	ctx.SetCookie("session", token, 86400, "/", "", false, true)

	ctx.JSON(200, gin.H{"message": "Login successful"})
}
