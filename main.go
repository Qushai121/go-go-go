package main

import (
	"log"
	"os"

	"hrms_go/config"
	"hrms_go/routes"

	_ "hrms_go/docs"

	"github.com/gofiber/contrib/v3/swaggo"
	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
)

// @title HRMS API
// @version 1.0
// @description This is HRMS API documentation
// @host 127.0.0.1:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	godotenv.Load()

	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal("DB connection failed")
	}

	app := fiber.New()
	app.Get("/swagger/*", swaggo.HandlerDefault)
	// app.Static("/api/uploads", "./uploads")
	routes.Setup(app, db)

	log.Fatal(app.Listen(":" + os.Getenv("APP_PORT")))
}
