package controllers

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rafinhacuri/SanchezDNS/db"
	"github.com/rafinhacuri/SanchezDNS/models"
	"github.com/rafinhacuri/SanchezDNS/passwords"
	"go.mongodb.org/mongo-driver/bson"
)

func InsertUser(ctx *gin.Context) {
	var request models.UserRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(400, gin.H{"message": "failed to bind JSON"})
		return
	}

	if err := request.ValidateRequest(); err != nil {
		ctx.JSON(400, gin.H{"message": err.Error()})
		return
	}

	passwordHash, err := passwords.BCrypt(request.Password)
	if err != nil {
		ctx.JSON(500, gin.H{"message": "Failed to hash password"})
		return
	}

	user := &models.User{
		Email:     request.Email,
		Password:  passwordHash,
		Level:     request.Level,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	ctxReq, cancel := context.WithTimeout(ctx.Request.Context(), 5*time.Second)
	defer cancel()

	countUsers, err := db.Database.Collection("users").CountDocuments(ctxReq, bson.M{})
	if err != nil {
		ctx.JSON(500, gin.H{"message": err.Error()})
		return
	}

	if countUsers == 0 {
		user.Level = "admin"
	} else {
		user.Level = "user"
	}

	count, err := db.Database.Collection("users").CountDocuments(ctxReq, bson.M{"email": user.Email})
	if err != nil {
		ctx.JSON(500, gin.H{"message": err.Error()})
		return
	}
	if count > 0 {
		ctx.JSON(400, gin.H{"message": "User with this email already exists"})
		return
	}

	if _, err := db.Database.Collection("users").InsertOne(ctxReq, user); err != nil {
		ctx.JSON(500, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(201, gin.H{"message": "User created successfully"})
}

func ChangePassword(ctx *gin.Context) {
	var request models.ChangePasswordRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(400, gin.H{"message": "failed to bind JSON"})
		return
	}

	if err := request.ValidateChangePasswordRequest(); err != nil {
		ctx.JSON(400, gin.H{"message": err.Error()})
		return
	}

	username, _ := ctx.Get("username")

	if username != request.Email {
		ctx.JSON(403, gin.H{"message": "You can only change your own password"})
		return
	}

	passwordHash, err := passwords.BCrypt(request.Password)
	if err != nil {
		ctx.JSON(500, gin.H{"message": "Failed to hash password"})
		return
	}

	ctxReq, cancel := context.WithTimeout(ctx.Request.Context(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"password":  passwordHash,
			"updatedAt": time.Now(),
		},
	}

	result, err := db.Database.Collection("users").UpdateOne(ctxReq, bson.M{"email": request.Email}, update)
	if err != nil {
		ctx.JSON(500, gin.H{"message": err.Error()})
		return
	}
	if result.MatchedCount == 0 {
		ctx.JSON(404, gin.H{"message": "User not found"})
		return
	}

	ctx.JSON(200, gin.H{"message": "Password changed successfully"})
}
