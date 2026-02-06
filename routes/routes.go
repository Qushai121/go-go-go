package routes

import (
	"hrms_go/controllers"
	"hrms_go/middlewares"
	"hrms_go/repositories"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Setup(app *fiber.App, db *gorm.DB) {
	repo := repositories.NewUserRepository(db)

	auth := controllers.NewAuthController(repo)
	user := controllers.NewUserController(repo)
	setup := controllers.NewSetupController(repo)

	api := app.Group("/api")
	
	api.Post("/setup", setup.InitAdmin)
	api.Post("/login", auth.Login)

	users := api.Group("/users", middlewares.JWTProtected())
	users.Post("/", user.Create)
	users.Get("/", user.GetAll)
}
