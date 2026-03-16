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
	customerRepo := repositories.NewCustomerRepository(db)
	shiftRepo := repositories.NewShiftRepository(db)
	companiesRepo := repositories.NewCompaniesRepository(db)
	paramRepo := repositories.NewParamRepository(db)
	paramGroupRepo := repositories.NewParamGroupRepository(db)
	settingRepo := repositories.NewSettingRepository(db)

	auth := controllers.NewAuthController(repo)
	user := controllers.NewUserController(repo)
	setup := controllers.NewSetupController(repo)
	customer := controllers.NewCustomerController(customerRepo)
	shift := controllers.NewShiftController(shiftRepo)
	companies := controllers.NewCompaniesController(companiesRepo)
	paramController := controllers.NewParamController(paramRepo)	
	paramGroupController := controllers.NewParamGroupController(paramGroupRepo)
	settingController := controllers.NewSettingController(settingRepo)



	api := app.Group("/api")
	
	api.Post("/setup", setup.InitAdmin)
	api.Post("/login", auth.Login)

	users := api.Group("/users", middlewares.JWTProtected())
	users.Post("/", user.Create)
	users.Get("/", user.FindAll)
	users.Post("/update-shift", user.UpdateUserShift)

	customerRoute := api.Group("/customer", middlewares.JWTProtected())
	customerRoute.Post("/", customer.Create)
	customerRoute.Get("/", customer.FindAll)
	customerRoute.Put("/",customer.Update);
	customerRoute.Delete("/:id",customer.Delete);

	shiftRoute := api.Group("/shift", middlewares.JWTProtected())
	shiftRoute.Post("/", shift.Create)
	shiftRoute.Get("/", shift.FindAll)
	shiftRoute.Put("/",shift.Update);
	shiftRoute.Delete("/:id",shift.Delete);

	companiesRoute := api.Group("/companies",middlewares.JWTProtected())
	companiesRoute.Post("/",companies.Create);
	companiesRoute.Get("/",companies.FindAll);
	companiesRoute.Put("/",companies.Update);
	companiesRoute.Delete("/:id",companies.Delete);
	
	paramGroup := api.Group("/param-group")
	paramGroup.Post("/", paramGroupController.Create)
	paramGroup.Get("/", paramGroupController.FindAll)
	paramGroup.Put("/", paramGroupController.Update)
	paramGroup.Delete("/:id", paramGroupController.Delete)
	
	param := app.Group("/param")
	param.Post("/", paramController.Create)
	param.Get("/", paramController.FindAll)
	param.Put("/", paramController.Update)
	param.Delete("/:id", paramController.Delete)

	setting := app.Group("/setting")

	setting.Post("/", settingController.Create)
	setting.Get("/", settingController.FindAll)
	setting.Put("/", settingController.Update)
	setting.Delete("/:id", settingController.Delete)
}
