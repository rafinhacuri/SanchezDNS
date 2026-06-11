package cron

import (
	"github.com/robfig/cron/v3"

	"github.com/rafinhacuri/SanchezDNS/api/auth"
)

func TimeoutSessions() {
	c := cron.New()

	_, err := c.AddFunc("@every 10m", func() {
		auth.TimeoutSessions()
	})
	if err != nil {
		panic("Failed to add cron job: " + err.Error())
	}

	c.Start()
}
