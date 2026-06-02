package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/rafinhacuri/SanchezDNS/api/mongo"
)

func newSID() (string, error) {
	b := make([]byte, 32)

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}

func normalize(email string) string {
	return strings.ToLower(strings.ReplaceAll(email, " ", ""))
}

func CreateSession(ctx context.Context, email, ip, os, browser, location string) (string, error) {
	sid, err := newSID()
	if err != nil {
		return "", err
	}

	email = normalize(email)

	if ip == "" {
		ip = "Unknown"
	}

	if os == "" {
		os = "Unknown"
	}

	if browser == "" {
		browser = "Unknown"
	}

	if location == "" {
		location = "Unknown"
	}

	session := mongo.Session{
		ID:           bson.NewObjectID(),
		SID:          sid,
		Active:       true,
		Email:        email,
		CreatedAt:    time.Now().UTC(),
		LastSeen:     time.Now().UTC(),
		Ip:           ip,
		Os:           os,
		Browser:      browser,
		Location:     location,
		LogoutReason: "",
	}

	sessionsColl := mongo.Dns.Collection("sessions")

	_, err = sessionsColl.InsertOne(ctx, session)
	if err != nil {
		return "", err
	}

	return sid, nil
}
