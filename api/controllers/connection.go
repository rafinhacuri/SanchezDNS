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

func InsertConnection(ctx *gin.Context) {
	var request models.ConnectionRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(400, gin.H{"message": "failed to bind JSON"})
		return
	}

	if err := request.ValidateConnectionRequest(); err != nil {
		ctx.JSON(400, gin.H{"message": err.Error()})
		return
	}

	encryptedKey, err := passwords.Encrypt(request.ApiKey)
	if err != nil {
		ctx.JSON(500, gin.H{"message": err.Error()})
		return
	}

	connection := &models.Connection{
		Name:     request.Name,
		Host:     request.Host,
		ApiKey:   encryptedKey,
		ServerId: request.ServerId,
		Users:    request.Users,
		CreateAt: time.Now(),
		UpdateAt: time.Now(),
	}

	ctxReq, cancel := context.WithTimeout(ctx.Request.Context(), 5*time.Second)
	defer cancel()

	count, err := db.Database.Collection("connections").CountDocuments(ctxReq, bson.M{"name": connection.Name})
	if err != nil {
		ctx.JSON(500, gin.H{"message": err.Error()})
		return
	}

	if count > 0 {
		ctx.JSON(409, gin.H{"message": "connection already exists"})
		return
	}

	_, err = db.Database.Collection("connections").InsertOne(ctxReq, connection)
	if err != nil {
		ctx.JSON(500, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(201, gin.H{"message": "connection created successfully"})
}
