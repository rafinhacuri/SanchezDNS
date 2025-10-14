package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/rafinhacuri/SanchezDNS/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetUsers(ctx *gin.Context) {
	if !ctx.GetBool("admin") {
		ctx.JSON(403, gin.H{"error": "forbidden"})
		return
	}

	username := ctx.GetString("username")

	opts := options.Find().SetProjection(bson.M{"_id": 1, "email": 1})
	cursor, err := db.Database.Collection("users").Find(ctx.Request.Context(), bson.M{}, opts)
	if err != nil {
		ctx.JSON(500, gin.H{"message": "failed to fetch users"})
		return
	}

	var users []bson.M
	if err := cursor.All(ctx.Request.Context(), &users); err != nil {
		ctx.JSON(500, gin.H{"message": "failed to parse users"})
		return
	}

	for i, user := range users {
		if user["email"] == username {
			users = append(users[:i], users[i+1:]...)
			break
		}
	}

	ctx.JSON(200, users)
}
