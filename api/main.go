package main

import (
	"io"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rafinhacuri/SanchezDNS/db"
	"github.com/rafinhacuri/SanchezDNS/routes"
)

func init() {
	godotenv.Load("../.env")

	ssl := os.Getenv("MONGO_SSL") == "true"

	if err := db.InitDB(ssl, os.Getenv("MONGO_USERNAME"), os.Getenv("MONGO_URL"), os.Getenv("MONGO_PASSWORD"), os.Getenv("MONGO_DB_NAME")); err != nil {
		log.Fatal("Error to connect to database:", err)
	}
}

func main() {

	gin.DefaultWriter = io.Discard

	server := gin.Default()

	server.Use(gin.LoggerWithWriter(os.Stdout, "/healthcheck"))

	server.SetTrustedProxies([]string{"127.0.0.1", "::1"})

	routes.RegisterRoutes(server)
	server.Run(":8080")
}
