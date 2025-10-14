package controllers

import (
	"context"
	"slices"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rafinhacuri/SanchezDNS/db"
	"github.com/rafinhacuri/SanchezDNS/models"
	"github.com/rafinhacuri/SanchezDNS/passwords"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func InsertConnection(ctx *gin.Context) {
	if !ctx.GetBool("admin") {
		ctx.JSON(403, gin.H{"error": "forbidden"})
		return
	}

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

func AddUserToConnection(ctx *gin.Context) {
	if !ctx.GetBool("admin") {
		ctx.JSON(403, gin.H{"error": "forbidden"})
		return
	}

	var request struct {
		Email      string `json:"email"`
		Connection string `json:"connection"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(400, gin.H{"message": "failed to bind JSON"})
		return
	}

	if request.Email == "" || request.Connection == "" {
		ctx.JSON(400, gin.H{"message": "email and connection are required"})
		return
	}

	ctxReq, cancel := context.WithTimeout(ctx.Request.Context(), 5*time.Second)
	defer cancel()

	var user models.User
	err := db.Database.Collection("users").FindOne(ctxReq, bson.M{"email": request.Email}).Decode(&user)
	if err != nil {
		ctx.JSON(404, gin.H{"message": "user not found"})
		return
	}

	var connection models.Connection

	connectionID, err := primitive.ObjectIDFromHex(request.Connection)
	if err != nil {
		ctx.JSON(400, gin.H{"message": "invalid connection ID"})
		return
	}

	err = db.Database.Collection("connections").FindOne(ctxReq, bson.M{"_id": connectionID}).Decode(&connection)
	if err != nil {
		ctx.JSON(404, gin.H{"message": "connection not found"})
		return
	}

	if slices.Contains(connection.Users, request.Email) {
		ctx.JSON(409, gin.H{"message": "user already added to connection"})
		return
	}

	_, err = db.Database.Collection("connections").UpdateOne(ctxReq, bson.M{"_id": connectionID}, bson.M{
		"$push": bson.M{"users": request.Email},
		"$set":  bson.M{"updateAt": time.Now()},
	})
	if err != nil {
		ctx.JSON(500, gin.H{"message": "failed to add user to connection"})
		return
	}

	ctx.JSON(200, gin.H{"message": "user added to connection successfully"})
}
