package main

import (
	"log"
	"os"

	"hrms_go/config"
	"hrms_go/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal("DB connection failed")
	}

	app := fiber.New()

	app.Static("/api/uploads", "./uploads")
	routes.Setup(app, db)

	log.Fatal(app.Listen(":" + os.Getenv("APP_PORT")))
}
