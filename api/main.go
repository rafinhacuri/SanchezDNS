package main

import (
	"io"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/rafinhacuri/SanchezDNS/api/env"
	"github.com/rafinhacuri/SanchezDNS/api/mongo"
	"github.com/rafinhacuri/SanchezDNS/api/redis"
	"github.com/rafinhacuri/SanchezDNS/api/routes"
	"github.com/rafinhacuri/SanchezDNS/api/s3"
)

func main() {
	_ = godotenv.Load("../.env")

	mongo.Connect()
	redis.Connect()
	s3.Connect()

	go mongo.Setup()

	gin.DefaultWriter = io.Discard

	server := gin.New()

	server.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		Output: os.Stdout,
		Skip: func(c *gin.Context) bool {
			if c.Request.URL.Path == "/healthcheck" {
				return true
			}

			ip := c.ClientIP()

			return ip == "152.84.253.18" || ip == "2804:1f10:8000:2::e"
		},
	}))

	server.Use(gin.Recovery())

	var err error

	err = server.SetTrustedProxies([]string{"127.0.0.1", "::1"})
	if err != nil {
		panic("Failed to set trusted proxies: " + err.Error())
	}

	routes.RegisterRoutes(server)

	if env.C.DevKey == "" || env.C.DevCert == "" {
		err = server.Run(":8080")
	} else {
		err = server.RunTLS(":8080", env.C.DevCert, env.C.DevKey)
	}

	if err != nil {
		panic("Failed to start gin: " + err.Error())
	}
}
