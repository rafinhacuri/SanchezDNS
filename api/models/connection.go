package models

import (
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Connection struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string             `bson:"name" json:"name" binding:"required"`
	Host     string             `bson:"host" json:"host" binding:"required"`
	ApiKey   string             `bson:"apiKey" json:"apiKey" binding:"required"`
	ServerId string             `bson:"serverId" json:"serverId" binding:"required"`
	Users    []string           `bson:"users" json:"users" binding:"required"`
	CreateAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdateAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type ConnectionRequest struct {
	Name     string   `bson:"name" json:"name"`
	Host     string   `bson:"host" json:"host"`
	ApiKey   string   `bson:"apiKey" json:"apiKey"`
	ServerId string   `bson:"serverId" json:"serverId"`
	Users    []string `bson:"users" json:"users"`
}

func (u *ConnectionRequest) ValidateConnectionRequest() error {
	if strings.TrimSpace(u.Name) == "" {
		return errors.New("the field 'name' is required")
	}
	if strings.TrimSpace(u.Host) == "" {
		return errors.New("the field 'host' is required")
	}
	if strings.TrimSpace(u.ApiKey) == "" {
		return errors.New("the field 'apiKey' is required")
	}
	if len(u.ApiKey) < 10 {
		return errors.New("the field 'apiKey' must be at least 10 characters long")
	}
	if strings.TrimSpace(u.ServerId) == "" {
		u.ServerId = "localhost"
	}

	return nil
}
