package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"hrms_go/config"
	"hrms_go/routes"
)

func main() {
	godotenv.Load()

	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal("DB connection failed")
	}

	app := fiber.New()
	routes.Setup(app, db)

	log.Fatal(app.Listen(":" + os.Getenv("APP_PORT")))
}
