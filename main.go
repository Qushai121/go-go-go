package main

import (
	"log"
	"os"

	"hrms_go/config"
	"hrms_go/routes"

	_ "hrms_go/docs"

	"github.com/gofiber/contrib/v3/swaggo"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/joho/godotenv"
)

// http://192.168.1.2:7701/swagger/index.html

// @title HRMS API
// @version 1.0
// @description This is HRMS API documentation
// @host 192.168.1.2:7701
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
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodPut,
			fiber.MethodDelete,
			fiber.MethodPatch,
			fiber.MethodOptions,
		},
		AllowHeaders: []string{
			fiber.HeaderOrigin,
			fiber.HeaderContentType,
			fiber.HeaderAccept,
			fiber.HeaderAuthorization,
		},
	}))
	app.Get("/swagger/*", swaggo.HandlerDefault)
	app.Use("/uploads", static.New("./uploads"))
	routes.Setup(app, db)

	log.Fatal(app.Listen(":" + os.Getenv("APP_PORT")))
}
