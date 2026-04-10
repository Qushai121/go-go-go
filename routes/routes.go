package routes

import (
	"hrms_go/controllers"
	"hrms_go/middlewares"
	"hrms_go/repositories"

	"github.com/gofiber/fiber/v3"
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
	SubmissionRepo := repositories.NewSubmissionRepository(db)
	leaveRepo := repositories.NewLeaveRepository(db)
	wfhRepo := repositories.NewWFHRepository(db)
	approvalRepo := repositories.NewApprovalRepository(db)
	approvalTemplateRepo := repositories.NewApprovalTemplateRepository(db)

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
	submissionController := controllers.NewSubmissionController(SubmissionRepo)
	leaveController := controllers.NewLeaveController(leaveRepo)
	wfhControler := controllers.NewWFHController(wfhRepo)
	approvalController := controllers.NewApprovalController(approvalRepo)
	approvalTemplateController := controllers.NewApprovalTemplateController(approvalTemplateRepo)

	api := app.Group("/api")

	api.Post("/setup", setup.InitAdmin)
	api.Post("/login", auth.Login)

	users := api.Group("/users", middlewares.JWTProtected())
	users.Post("/", user.Create)
	users.Get("/", user.FindAll)
	users.Get("/me", user.Me)
	users.Post("/update-shift", user.UpdateUserShift)
	users.Post("/update-profile-picture", user.UpdateUserPicture)

	customerRoute := api.Group("/customer", middlewares.JWTProtected())
	customerRoute.Post("/", customer.Create)
	customerRoute.Get("/", customer.FindAll)
	customerRoute.Put("/", customer.Update)
	customerRoute.Delete("/:id", customer.Delete)

	shiftRoute := api.Group("/shift", middlewares.JWTProtected())
	shiftRoute.Post("/", shift.Create)
	shiftRoute.Get("/", shift.FindAll)
	shiftRoute.Put("/", shift.Update)
	shiftRoute.Delete("/:id", shift.Delete)

	companiesRoute := api.Group("/companies", middlewares.JWTProtected())
	companiesRoute.Post("/", companies.Create)
	companiesRoute.Get("/", companies.FindAll)
	companiesRoute.Put("/", companies.Update)
	companiesRoute.Delete("/:id", companies.Delete)

	paramGroup := api.Group("/param-group",middlewares.JWTProtected())
	paramGroup.Post("/", paramGroupController.Create)
	paramGroup.Get("/", paramGroupController.FindAll)
	paramGroup.Put("/", paramGroupController.Update)
	paramGroup.Delete("/:id", paramGroupController.Delete)

	param := api.Group("/param", middlewares.JWTProtected())
	param.Post("/", paramController.Create)
	param.Get("/", paramController.FindAll)
	param.Put("/", paramController.Update)
	param.Delete("/:id", paramController.Delete)

	setting := api.Group("/setting", middlewares.JWTProtected())
	setting.Post("/", settingController.Create)
	setting.Get("/", settingController.FindAll)
	setting.Put("/", settingController.Update)
	setting.Delete("/:id", settingController.Delete)

	attendance := api.Group("/attendance", middlewares.JWTProtected())
	attendance.Post("/", attendanceController.Create)
	attendance.Get("/", attendanceController.FindAll)
	attendance.Get("/me", attendanceController.FindByUser)

	office := api.Group("/office", middlewares.JWTProtected())
	office.Post("/", officeController.Create)
	office.Get("/", officeController.FindAll)
	office.Put("/", officeController.Update)
	office.Delete("/:id", officeController.Delete)

	Submission := api.Group("claim_submission", middlewares.JWTProtected())
	Submission.Post("/", submissionController.Create)
	Submission.Get("/", submissionController.FindAll)
	Submission.Get("/me", submissionController.FindByUser)
	Submission.Put("/", submissionController.Update)
	Submission.Delete("/:id", submissionController.Delete)

	leaveRoute := api.Group("leave", middlewares.JWTProtected())
	leaveRoute.Post("/", leaveController.Create)
	leaveRoute.Get("/", leaveController.FindAll)
	leaveRoute.Put("/", leaveController.Update)
	leaveRoute.Delete("/:id", leaveController.Delete)

	wfhRoute := api.Group("wfh", middlewares.JWTProtected())
	wfhRoute.Post("/", wfhControler.Create)
	wfhRoute.Get("/", wfhControler.FindAll)
	wfhRoute.Get("/me", wfhControler.FindAll)
	// wfhRoute.Put("/", wfhControler.Update)
	wfhRoute.Delete("/:id", wfhControler.Delete)

	approval := api.Group("/approval")
	approval.Get("/", approvalController.FindAll)
	approval.Get("/:id", approvalController.Detail)
	approval.Delete("/:id", approvalController.Delete)
	approval.Post("/approve", approvalController.Approve)

	template := api.Group("/approval_template")
	// header
	template.Post("/header", approvalTemplateController.CreateHeader)
	template.Get("/header", approvalTemplateController.FindAllHeader)
	template.Get("/header/:id", approvalTemplateController.DetailHeader)
	template.Put("/header/:id", approvalTemplateController.UpdateHeader)
	template.Delete("/header/:id", approvalTemplateController.DeleteHeader)

	// detail
	template.Post("/detail", approvalTemplateController.CreateDetail)
	template.Get("/detail/:header_id", approvalTemplateController.FindDetailByHeader)
	template.Put("/detail/:id", approvalTemplateController.UpdateDetail)
	template.Delete("/detail/:id", approvalTemplateController.DeleteDetail)

}
