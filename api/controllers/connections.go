package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/rafinhacuri/SanchezDNS/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetConnections(ctx *gin.Context) {
	username := ctx.GetString("username")
	isAdmin := ctx.GetBool("admin")

	var filter bson.M
	if isAdmin {
		filter = bson.M{}
	} else {
		filter = bson.M{"users": username}
	}

	opts := options.Find().SetProjection(bson.M{"_id": 1, "name": 1, "users": 1})

	cursor, err := db.Database.Collection("connections").Find(ctx.Request.Context(), filter, opts)
	if err != nil {
		ctx.JSON(500, gin.H{"message": "failed to fetch connections"})
		return
	}

	var connections []bson.M
	if err := cursor.All(ctx.Request.Context(), &connections); err != nil {
		ctx.JSON(500, gin.H{"message": "failed to parse connections"})
		return
	}

	ctx.JSON(200, connections)
}

func GetFullConnections(ctx *gin.Context) {
	if isAdmin := ctx.GetBool("admin"); !isAdmin {
		ctx.JSON(403, gin.H{"error": "forbidden"})
		return
	}

	opts := options.Find().SetProjection(bson.M{"apiKey": 0})

	cursor, err := db.Database.Collection("connections").Find(ctx.Request.Context(), bson.M{}, opts)
	if err != nil {
		ctx.JSON(500, gin.H{"message": "failed to fetch connections"})
		return
	}

	var connections []bson.M
	if err := cursor.All(ctx.Request.Context(), &connections); err != nil {
		ctx.JSON(500, gin.H{"message": "failed to parse connections"})
		return
	}

	ctx.JSON(200, connections)
}
