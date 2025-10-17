package controllers

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/rafinhacuri/SanchezDNS/db"
	"github.com/rafinhacuri/SanchezDNS/models"
	"github.com/rafinhacuri/SanchezDNS/utils"
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

	encryptedKey, err := utils.Encrypt(request.ApiKey)
	if err != nil {
		ctx.JSON(500, gin.H{"message": err.Error()})
		return
	}

	connection := &models.Connection{
		Name:      request.Name,
		Host:      request.Host,
		ApiKey:    encryptedKey,
		ServerId:  request.ServerId,
		Users:     request.Users,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	ctxReq, cancel := context.WithTimeout(ctx.Request.Context(), 5*time.Second)
	defer cancel()

	base := utils.NormalizeBase(connection.Host)

	httpc := resty.New().SetBaseURL(base).SetHeader("X-API-Key", request.ApiKey).SetHeader("Accept", "application/json").SetTimeout(6 * time.Second).SetRetryCount(2)

	resp, erro := httpc.R().SetContext(ctxReq).Get(fmt.Sprintf("/api/v1/servers/%s/statistics", connection.ServerId))
	if erro != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": fmt.Sprintf("failed to reach PowerDNS: %v", erro)})
		return
	}

	if resp.StatusCode() != http.StatusOK {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": fmt.Sprintf("unexpected response from PowerDNS: %s", resp.Status())})
		return
	}

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

	username := ctx.GetString("username")

	log := &models.Log{
		HostServer:   connection.Host,
		Zone:         "",
		IdConnection: "",
		Username:     username,
		Action:       "create_connection",
		Details:      fmt.Sprintf("User %s created connection %s", username, connection.Name),
		CreatedAt:    time.Now(),
	}

	_, _ = db.Database.Collection("logs").InsertOne(ctxReq, log)

	ctx.JSON(201, gin.H{"message": "connection created successfully"})
}

func AddUser(ctx *gin.Context) {
	if !ctx.GetBool("admin") {
		ctx.JSON(403, gin.H{"error": "forbidden"})
		return
	}

	var request models.AddUserRequest

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
		"$set":  bson.M{"updatedAt": time.Now()},
	})
	if err != nil {
		ctx.JSON(500, gin.H{"message": "failed to add user to connection"})
		return
	}

	username := ctx.GetString("username")

	log := &models.Log{
		HostServer:   connection.Host,
		IdConnection: request.Connection,
		Zone:         "",
		Username:     username,
		Action:       "add_user_to_connection",
		Details:      fmt.Sprintf("User %s added to connection %s", request.Email, connection.Name),
		CreatedAt:    time.Now(),
	}

	_, _ = db.Database.Collection("logs").InsertOne(ctxReq, log)

	ctx.JSON(200, gin.H{"message": "user added to connection successfully"})
}

func RemoveUser(ctx *gin.Context) {
	if !ctx.GetBool("admin") {
		ctx.JSON(403, gin.H{"error": "forbidden"})
		return
	}

	var request models.AddUserRequest

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

	if !slices.Contains(connection.Users, request.Email) {
		ctx.JSON(409, gin.H{"message": "user not associated with connection"})
		return
	}

	_, err = db.Database.Collection("connections").UpdateOne(ctxReq, bson.M{"_id": connectionID}, bson.M{
		"$pull": bson.M{"users": request.Email},
		"$set":  bson.M{"updatedAt": time.Now()},
	})
	if err != nil {
		ctx.JSON(500, gin.H{"message": "failed to remove user from connection"})
		return
	}

	username := ctx.GetString("username")

	log := &models.Log{
		HostServer:   connection.Host,
		Zone:         "",
		IdConnection: request.Connection,
		Username:     username,
		Action:       "remove_user_from_connection",
		Details:      fmt.Sprintf("User %s removed from connection %s", request.Email, connection.Name),
		CreatedAt:    time.Now(),
	}

	_, _ = db.Database.Collection("logs").InsertOne(ctxReq, log)

	ctx.JSON(200, gin.H{"message": "user removed from connection successfully"})
}

func GetConnection(ctx *gin.Context) {
	primitiveId := ctx.Query("connection")

	username := ctx.GetString("username")
	isAdmin := ctx.GetBool("admin")

	if primitiveId == "" {
		ctx.JSON(400, gin.H{"message": "to get a specific connection, use /connections/:id"})
		return
	}

	var connection models.Connection

	id, err := primitive.ObjectIDFromHex(primitiveId)

	if err != nil {
		ctx.JSON(400, gin.H{"message": "invalid connection ID"})
		return
	}

	err = db.Database.Collection("connections").FindOne(ctx.Request.Context(), bson.M{"_id": id}).Decode(&connection)
	if err != nil {
		ctx.JSON(404, gin.H{"message": "connection not found"})
		return
	}

	if !isAdmin && !slices.Contains(connection.Users, username) {
		ctx.JSON(403, gin.H{"message": "forbidden"})
		return
	}

	connection.ApiKey = ""
	connection.Users = nil

	ctx.JSON(200, connection)
}

func EditConnection(ctx *gin.Context) {
	primitiveId := ctx.Query("connection")

	username := ctx.GetString("username")
	isAdmin := ctx.GetBool("admin")

	if primitiveId == "" {
		ctx.JSON(400, gin.H{"message": "to get a specific connection, use /connections/:id"})
		return
	}

	var connection models.Connection

	id, err := primitive.ObjectIDFromHex(primitiveId)

	if err != nil {
		ctx.JSON(400, gin.H{"message": "invalid connection ID"})
		return
	}

	err = db.Database.Collection("connections").FindOne(ctx.Request.Context(), bson.M{"_id": id}).Decode(&connection)
	if err != nil {
		ctx.JSON(404, gin.H{"message": "connection not found"})
		return
	}

	if !isAdmin && !slices.Contains(connection.Users, username) {
		ctx.JSON(403, gin.H{"message": "forbidden"})
		return
	}

	var request models.ConnectingEditRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(400, gin.H{"message": "failed to bind JSON"})
		return
	}

	if err := request.ValidateConnectionEdit(); err != nil {
		ctx.JSON(400, gin.H{"message": err.Error()})
		return
	}

	update := bson.M{
		"name":      request.Name,
		"host":      request.Host,
		"serverId":  request.ServerId,
		"updatedAt": time.Now(),
	}

	ctxReq, cancel := context.WithTimeout(ctx.Request.Context(), 5*time.Second)
	defer cancel()

	_, err = db.Database.Collection("connections").UpdateOne(ctxReq, bson.M{"_id": id}, bson.M{"$set": update})
	if err != nil {
		ctx.JSON(500, gin.H{"message": err.Error()})
		return
	}

	changes := []string{}
	if connection.Name != request.Name {
		changes = append(changes, fmt.Sprintf("name: %s -> %s", connection.Name, request.Name))
	}
	if connection.Host != request.Host {
		changes = append(changes, fmt.Sprintf("host: %s -> %s", connection.Host, request.Host))
	}
	if connection.ServerId != request.ServerId {
		changes = append(changes, fmt.Sprintf("serverId: %s -> %s", connection.ServerId, request.ServerId))
	}

	changeDetails := "no fields changed"
	if len(changes) > 0 {
		changeDetails = strings.Join(changes, ", ")
	}

	log := &models.Log{
		HostServer:   connection.Host,
		IdConnection: primitiveId,
		Zone:         "",
		Username:     username,
		Action:       "edit_connection",
		Details:      fmt.Sprintf("User %s edited connection %s: %s", username, connection.Name, changeDetails),
		CreatedAt:    time.Now(),
	}

	_, _ = db.Database.Collection("logs").InsertOne(ctxReq, log)

	ctx.JSON(200, gin.H{"message": "connection updated successfully"})
}

func DeleteConnection(ctx *gin.Context) {
	primitiveId := ctx.Query("connection")

	username := ctx.GetString("username")
	isAdmin := ctx.GetBool("admin")

	if primitiveId == "" {
		ctx.JSON(400, gin.H{"message": "to get a specific connection, use /connections/:id"})
		return
	}

	var connection models.Connection

	id, err := primitive.ObjectIDFromHex(primitiveId)

	if err != nil {
		ctx.JSON(400, gin.H{"message": "invalid connection ID"})
		return
	}

	err = db.Database.Collection("connections").FindOne(ctx.Request.Context(), bson.M{"_id": id}).Decode(&connection)
	if err != nil {
		ctx.JSON(404, gin.H{"message": "connection not found"})
		return
	}

	if !isAdmin && !slices.Contains(connection.Users, username) {
		ctx.JSON(403, gin.H{"message": "forbidden"})
		return
	}

	_, err = db.Database.Collection("connections").DeleteOne(ctx.Request.Context(), bson.M{"_id": id})
	if err != nil {
		ctx.JSON(500, gin.H{"message": err.Error()})
		return
	}

	log := &models.Log{
		HostServer:   connection.Host,
		Zone:         "",
		IdConnection: primitiveId,
		Username:     username,
		Action:       "delete_connection",
		Details:      fmt.Sprintf("User %s deleted connection %s", username, connection.Name),
		CreatedAt:    time.Now(),
	}

	_, _ = db.Database.Collection("logs").InsertOne(ctx.Request.Context(), log)

	ctx.JSON(200, gin.H{"message": "connection deleted successfully"})
}
