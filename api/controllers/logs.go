package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rafinhacuri/SanchezDNS/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetLogs(ctx *gin.Context) {
	username := ctx.GetString("username")
	isAdmin := ctx.GetBool("admin")

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	search, _ := ctx.GetQuery("filter")

	skip := (page - 1) * limit

	filter := bson.M{}

	var orFilters []bson.M
	if search != "" {
		orFilters = []bson.M{
			{"username": bson.M{"$regex": search, "$options": "i"}},
			{"action": bson.M{"$regex": search, "$options": "i"}},
			{"details": bson.M{"$regex": search, "$options": "i"}},
			{"zone": bson.M{"$regex": search, "$options": "i"}},
			{"hostServer": bson.M{"$regex": search, "$options": "i"}},
		}
	}

	if !isAdmin {
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

		connectionFilter := bson.M{"idConnection": bson.M{"$in": connectionIds}}
		if len(orFilters) > 0 {
			filter = bson.M{
				"$and": []bson.M{
					connectionFilter,
					{"$or": orFilters},
				},
			}
		} else {
			filter = connectionFilter
		}
	} else {
		if len(orFilters) > 0 {
			filter = bson.M{"$or": orFilters}
		}
	}

	total, _ := db.Database.Collection("logs").CountDocuments(ctx.Request.Context(), filter)

	opts := options.Find().
		SetSort(bson.M{"createdAt": -1}).
		SetSkip(int64(skip)).
		SetLimit(int64(limit))

	cursor, err := db.Database.Collection("logs").Find(ctx.Request.Context(), filter, opts)
	if err != nil {
		ctx.JSON(500, gin.H{"message": "failed to fetch logs"})
		return
	}

	var logs []bson.M
	if err := cursor.All(ctx.Request.Context(), &logs); err != nil {
		ctx.JSON(500, gin.H{"message": "failed to parse logs"})
		return
	}

	ctx.JSON(200, gin.H{
		"data":  logs,
		"total": total,
	})
}
