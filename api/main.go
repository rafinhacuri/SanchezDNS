package main

import (
	"io"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rafinhacuri/SanchezDNS/routes"
)

func main() {
	godotenv.Load("../.env")

	key := os.Getenv("DEV_KEY")
	cert := os.Getenv("DEV_CERT")

	gin.DefaultWriter = io.Discard

	server := gin.Default()

	server.Use(gin.LoggerWithWriter(os.Stdout, "/healthcheck"))

	server.SetTrustedProxies([]string{"127.0.0.1", "::1"})

	routes.RegisterRoutes(server)

	if key == "" || cert == "" {
		server.Run(":8080")
	} else {
		server.RunTLS(":8080", cert, key)
	}
}
