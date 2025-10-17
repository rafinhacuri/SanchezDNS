package controllers

import (
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/rafinhacuri/SanchezDNS/db"
	"github.com/rafinhacuri/SanchezDNS/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func permission(ctx *gin.Context) (bool, *models.Connection) {
	primitiveId := ctx.Query("connection")

	username := ctx.GetString("username")
	isAdmin := ctx.GetBool("admin")

	var connection models.Connection

	if primitiveId == "" {
		ctx.JSON(400, gin.H{"message": "connection ID is required"})
		return false, nil
	}

	id, err := primitive.ObjectIDFromHex(primitiveId)

	if err != nil {
		ctx.JSON(400, gin.H{"message": "invalid connection ID"})
		return false, nil
	}

	err = db.Database.Collection("connections").FindOne(ctx.Request.Context(), bson.M{"_id": id}).Decode(&connection)
	if err != nil {
		ctx.JSON(404, gin.H{"message": "connection not found"})
		return false, nil
	}

	if !isAdmin && !slices.Contains(connection.Users, username) {
		ctx.JSON(403, gin.H{"message": "forbidden"})
		return false, nil
	}

	return true, &connection
}
