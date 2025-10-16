package models

import "time"

type Log struct {
	ID           string    `bson:"_id,omitempty" json:"id"`
	IdConnection string    `bson:"idConnection" json:"idConnection"`
	HostServer   string    `bson:"hostServer" json:"hostServer"`
	Zone         string    `bson:"zone" json:"zone"`
	Username     string    `bson:"username" json:"username"`
	Action       string    `bson:"action" json:"action"`
	Details      string    `bson:"details" json:"details"`
	CreatedAt    time.Time `bson:"createdAt" json:"createdAt"`
}
