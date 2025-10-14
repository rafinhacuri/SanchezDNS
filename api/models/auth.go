package models

import (
	"context"
	"errors"

	"github.com/rafinhacuri/SanchezDNS/db"
	"github.com/rafinhacuri/SanchezDNS/passwords"
	"github.com/rafinhacuri/SanchezDNS/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Auth struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (u *Auth) Validate() error {
	if u.Email == "" {
		return errors.New("the field 'email' is required")
	}
	if u.Password == "" {
		return errors.New("the field 'password' is required")
	}
	if err := utils.ValidatePassword(u.Password); err != nil {
		return errors.New("the field 'password' must be at least 6 characters long")
	}
	if err := utils.ValidateEmail(u.Email); err != nil {
		return errors.New("invalid email format")
	}

	return nil
}

func (u *Auth) Login(ctx context.Context) (token string, isAdmin bool, err error) {
	var user User
	if err := db.Database.Collection("users").FindOne(ctx, bson.M{"email": u.Email}).Decode(&user); err != nil {
		return "", false, errors.New("invalid email or password")
	}

	if !passwords.VerifyBCrypt(u.Password, user.Password) {
		return "", false, errors.New("invalid email or password")
	}

	token, err = utils.GenerateJWT(user.Email, user.Level == "admin")
	if err != nil {
		return "", false, err
	}

	return token, user.Level == "admin", nil
}
