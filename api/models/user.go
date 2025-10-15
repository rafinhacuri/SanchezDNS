package models

import (
	"errors"
	"strings"
	"time"

	"github.com/rafinhacuri/SanchezDNS/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Email     string             `bson:"email" json:"email" binding:"required"`
	Password  string             `bson:"password" json:"password" binding:"required"`
	Level     string             `bson:"level" json:"level" binding:"required,oneof=admin user"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type UserRequest struct {
	Email    string `bson:"email" json:"email"`
	Password string `bson:"password" json:"password"`
	Level    string `bson:"level" json:"level"`
}

func (u *UserRequest) ValidateRequest() error {
	if strings.TrimSpace(u.Email) == "" {
		return errors.New("the field 'email' is required")
	}
	if strings.TrimSpace(u.Password) == "" {
		return errors.New("the field 'password' is required")
	}
	if strings.TrimSpace(u.Level) == "" {
		return errors.New("the field 'level' is required")
	}
	if u.Level != "admin" && u.Level != "user" {
		return errors.New("the field 'level' must be 'admin' or 'user'")
	}
	if err := utils.ValidateEmail(u.Email); err != nil {
		return errors.New("invalid email format")
	}

	if err := utils.ValidatePassword(u.Password); err != nil {
		return errors.New("invalid password format")
	}
	return nil
}

type ChangePasswordRequest struct {
	Email    string `bson:"email" json:"email"`
	Password string `bson:"password" json:"password"`
}

func (u *ChangePasswordRequest) ValidateChangePasswordRequest() error {
	if strings.TrimSpace(u.Email) == "" {
		return errors.New("the field 'email' is required")
	}
	if strings.TrimSpace(u.Password) == "" {
		return errors.New("the field 'password' is required")
	}
	if err := utils.ValidateEmail(u.Email); err != nil {
		return errors.New("invalid email format")
	}

	if err := utils.ValidatePassword(u.Password); err != nil {
		return errors.New("invalid password format")
	}
	return nil
}

type AddUserRequest struct {
	Email      string `bson:"email" json:"email"`
	Connection string `bson:"connection" json:"connection"`
}

func (u *AddUserRequest) ValidateAddUserRequest() error {
	if strings.TrimSpace(u.Email) == "" {
		return errors.New("the field 'email' is required")
	}
	if strings.TrimSpace(u.Connection) == "" {
		return errors.New("the field 'connection' is required")
	}
	if err := utils.ValidateEmail(u.Email); err != nil {
		return errors.New("invalid email format")
	}

	return nil
}
