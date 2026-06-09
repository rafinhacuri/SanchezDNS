package redis

import "time"

type SessionJSON struct {
	ID           string    `json:"_id"`
	SID          string    `json:"sid"`
	Active       bool      `json:"active"`
	Email        string    `json:"email"`
	CreatedAt    time.Time `json:"createdAt"`
	LastSeen     time.Time `json:"lastSeen"`
	Ip           string    `json:"ip"`
	Os           string    `json:"os"`
	Browser      string    `json:"browser"`
	Location     string    `json:"location"`
	Representar  string    `json:"representar,omitempty"`
	LogoutReason string    `json:"logoutReason"`
}
