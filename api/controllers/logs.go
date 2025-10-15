package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/rafinhacuri/SanchezDNS/db"
	"go.mongodb.org/mongo-driver/bson"
)

func GetLogs(ctx *gin.Context) {
	username := ctx.GetString("username")
	isAdmin := ctx.GetBool("admin")

	if isAdmin {
		cursor, err := db.Database.Collection("logs").Find(ctx.Request.Context(), gin.H{})
		if err != nil {
			ctx.JSON(500, gin.H{"message": "failed to fetch logs"})
			return
		}

		var logs []bson.M
		if err := cursor.All(ctx.Request.Context(), &logs); err != nil {
			ctx.JSON(500, gin.H{"message": "failed to parse logs"})
			return
		}

		ctx.JSON(200, logs)
		return
	}

	cursorConnections, err := db.Database.Collection("connections").Find(ctx.Request.Context(), bson.M{"users": username})
	if err != nil {
		ctx.JSON(500, gin.H{"message": "failed to fetch connections"})
		return
	}

	var connections []bson.M
	if err := cursorConnections.All(ctx.Request.Context(), &connections); err != nil {
		ctx.JSON(500, gin.H{"message": "failed to parse connections"})
		return
	}

	var connectionIds []string
	for _, conn := range connections {
		if id, ok := conn["_id"].(string); ok {
			connectionIds = append(connectionIds, id)
		}
	}

	cursorLogs, err := db.Database.Collection("logs").Find(ctx.Request.Context(), bson.M{"idConnection": bson.M{"$in": connectionIds}})
	if err != nil {
		ctx.JSON(500, gin.H{"message": "failed to fetch logs"})
		return
	}

	var logs []bson.M
	if err := cursorLogs.All(ctx.Request.Context(), &logs); err != nil {
		ctx.JSON(500, gin.H{"message": "failed to parse logs"})
		return
	}

	ctx.JSON(200, logs)
}
