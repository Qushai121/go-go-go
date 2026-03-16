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
	attendanceRepo := repositories.NewAttendanceRepository(db)
	officeRepo := repositories.NewOfficeRepository(db)
	claimSubmissionRepo := repositories.NewClaimSubmissionRepository(db)
	leaveRepo := repositories.NewLeaveRepository(db)

	auth := controllers.NewAuthController(repo)
	user := controllers.NewUserController(repo)
	setup := controllers.NewSetupController(repo)
	customer := controllers.NewCustomerController(customerRepo)
	shift := controllers.NewShiftController(shiftRepo)
	companies := controllers.NewCompaniesController(companiesRepo)
	paramController := controllers.NewParamController(paramRepo)	
	paramGroupController := controllers.NewParamGroupController(paramGroupRepo)
	settingController := controllers.NewSettingController(settingRepo)
	attendanceController := controllers.NewAttendanceController(attendanceRepo)
	officeController := controllers.NewOfficeController(officeRepo)
	claimSubmissionController := controllers.NewClaimSubmissionController(claimSubmissionRepo)
	leaveController := controllers.NewLeaveController(leaveRepo)

	api := app.Group("/api")
	
	api.Post("/setup", setup.InitAdmin)
	api.Post("/login", auth.Login)

	users := api.Group("/users", middlewares.JWTProtected())
	users.Post("/", user.Create)
	users.Get("/", user.FindAll)
	users.Post("/update-shift", user.UpdateUserShift)
	users.Post("/update-profile-picture",user.UpdateUserPicture)

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
	
	param := api.Group("/param")
	param.Post("/", paramController.Create)
	param.Get("/", paramController.FindAll)
	param.Put("/", paramController.Update)
	param.Delete("/:id", paramController.Delete)

	setting := api.Group("/setting")

	setting.Post("/", settingController.Create)
	setting.Get("/", settingController.FindAll)
	setting.Put("/", settingController.Update)
	setting.Delete("/:id", settingController.Delete)
		
	attendance := api.Group("/attendance")
	attendance.Post("/", attendanceController.Create)
	attendance.Get("/", attendanceController.FindAll)
	attendance.Get("/user/:user_id", attendanceController.FindByUser)

	office := api.Group("/office")
	office.Post("/", officeController.Create)
	office.Get("/", officeController.FindAll)
	office.Put("/", officeController.Update)
	office.Delete("/:id", officeController.Delete)

	claimSubmission := api.Group("claim_submission")
	claimSubmission.Post("/",claimSubmissionController.Create)
	claimSubmission.Get("/",claimSubmissionController.FindAll)
	claimSubmission.Put("/",claimSubmissionController.Update)
	claimSubmission.Delete("/",claimSubmissionController.Delete)
	
	leaveRoute := api.Group("claim_submission")
	leaveRoute.Post("/",leaveController.Create)
	leaveRoute.Get("/",leaveController.FindAll)
	leaveRoute.Put("/",leaveController.Update)
	leaveRoute.Delete("/",leaveController.Delete)
	
}
