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
	companiesRepo := repositories.NewCompaniesRepository(db)
	paramRepo := repositories.NewParamRepository(db)
	paramGroupRepo := repositories.NewParamGroupRepository(db)
	settingRepo := repositories.NewSettingRepository(db)
	attendanceRepo := repositories.NewAttendanceRepository(db)
	siteRepo := repositories.NewSiteRepository(db)
	userSiteRepo := repositories.NewUserSiteRepository(db)
	approvalRepo := repositories.NewApprovalRepository(db)
	approvalTemplateRepo := repositories.NewApprovalTemplateRepository(db)
	userCompanyRepo := repositories.NewUserCompanyRepository(db)
	divisionRepo := repositories.NewDivisionRepository(db)
	departmentRepo := repositories.NewDepartmentRepository(db)
	titleRepo := repositories.NewTitleRepository(db)
	userGroupRepo := repositories.NewUserGroupRepository(db)
	receiptRepo := repositories.NewReceiptRepository(db)
	officePermitRepo := repositories.NewOfficePermitRepository(db)

	// Tambahan Leave / Cuti
	leaveRepo := repositories.NewLeaveRepository(db)
	leaveTypeRepo := repositories.NewLeaveTypeRepository(db)

	auth := controllers.NewAuthController(repo)
	user := controllers.NewUserController(repo)
	setup := controllers.NewSetupController(repo)
	companies := controllers.NewCompaniesController(companiesRepo)
	paramController := controllers.NewParamController(paramRepo)
	paramGroupController := controllers.NewParamGroupController(paramGroupRepo)
	settingController := controllers.NewSettingController(settingRepo)
	attendanceController := controllers.NewAttendanceController(attendanceRepo)
	siteController := controllers.NewSiteController(siteRepo)
	userSiteController := controllers.NewUserSiteController(userSiteRepo)
	approvalController := controllers.NewApprovalController(approvalRepo)
	approvalTemplateController := controllers.NewApprovalTemplateController(approvalTemplateRepo)
	userCompanyController := controllers.NewUserCompanyController(userCompanyRepo)
	divisionController := controllers.NewDivisionController(divisionRepo)
	departmentController := controllers.NewDepartmentController(departmentRepo)
	titleController := controllers.NewTitleController(titleRepo)
	userGroupController := controllers.NewUserGroupController(userGroupRepo)
	receiptController := controllers.NewReceiptController(receiptRepo)
	officePermitController := controllers.NewOfficePermitController(officePermitRepo)

	// Tambahan Leave / Cuti
	leaveController := controllers.NewLeaveController(leaveRepo)
	leaveTypeController := controllers.NewLeaveTypeController(leaveTypeRepo)

	api := app.Group("/api")

	api.Post("/setup", setup.InitAdmin)
	api.Post("/login", auth.Login)

	users := api.Group("/users", middlewares.JWTProtected())
	users.Post("/", user.Create)
	users.Get("/", user.FindAll)
	users.Get("/me", user.Me)
	users.Post("/verify-face", user.VerifyFace)
	users.Post("/update-profile-picture", user.UpdateUserPicture)

	companiesRoute := api.Group("/companies", middlewares.JWTProtected())
	companiesRoute.Post("/", companies.Create)
	companiesRoute.Get("/", companies.FindAll)
	companiesRoute.Put("/", companies.Update)
	companiesRoute.Delete("/:id", companies.Delete)

	site := api.Group("/site", middlewares.JWTProtected())
	site.Post("/", siteController.Create)
	site.Get("/", siteController.FindAll)
	site.Get("/search", siteController.Search)
	site.Put("/", siteController.Update)
	site.Delete("/:id", siteController.Delete)

	paramGroup := api.Group("/param-group", middlewares.JWTProtected())
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

	userCompany := api.Group("/user-company", middlewares.JWTProtected())
	userCompany.Post("/", userCompanyController.Create)
	userCompany.Get("/", userCompanyController.FindAll)
	userCompany.Put("/", userCompanyController.Update)
	userCompany.Delete("/:id", userCompanyController.Delete)

	userSite := api.Group("/user-site", middlewares.JWTProtected())
	userSite.Post("/", userSiteController.Create)
	userSite.Get("/", userSiteController.FindAll)
	userSite.Put("/", userSiteController.Update)
	userSite.Delete("/:id", userSiteController.Delete)

	division := api.Group("/division", middlewares.JWTProtected())
	division.Post("/", divisionController.Create)
	division.Get("/", divisionController.FindAll)
	division.Put("/", divisionController.Update)
	division.Delete("/:id", divisionController.Delete)

	department := api.Group("/department", middlewares.JWTProtected())
	department.Post("/", departmentController.Create)
	department.Get("/", departmentController.FindAll)
	department.Put("/", departmentController.Update)
	department.Delete("/:id", departmentController.Delete)

	title := api.Group("/title", middlewares.JWTProtected())
	title.Post("/", titleController.Create)
	title.Get("/", titleController.FindAll)
	title.Put("/", titleController.Update)
	title.Delete("/:id", titleController.Delete)

	userGroup := api.Group("/usergroup", middlewares.JWTProtected())
	userGroup.Post("/", userGroupController.Create)
	userGroup.Get("/", userGroupController.FindAll)
	userGroup.Put("/", userGroupController.Update)
	userGroup.Delete("/:id", userGroupController.Delete)

	receipt := api.Group("/receipt", middlewares.JWTProtected())
	receipt.Post("/", receiptController.Create)
	receipt.Get("/", receiptController.FindAll)
	receipt.Post("/submit", receiptController.Submit)
	receipt.Post("/detail", receiptController.CreateDraftDetail)
	receipt.Get("/detail/draft", receiptController.FindDraftDetails)
	receipt.Post("/:id/detail", receiptController.CreateDetail)
	receipt.Put("/detail/:detail_id", receiptController.UpdateDetail)
	receipt.Delete("/detail/:detail_id", receiptController.DeleteDetail)
	receipt.Put("/:id/submission", receiptController.UpdateSubmission)
	receipt.Get("/:id", receiptController.Detail)
	receipt.Put("/:id", receiptController.UpdateHeader)
	receipt.Delete("/:id", receiptController.Delete)

	officePermit := api.Group("/office-permit", middlewares.JWTProtected())
	officePermit.Post("/", officePermitController.Create)
	officePermit.Get("/", officePermitController.FindAll)
	officePermit.Put("/", officePermitController.Update)
	officePermit.Delete("/:id", officePermitController.Delete)

	// Tambahan route Leave / Cuti
	leave := api.Group("/leave", middlewares.JWTProtected())
	leave.Post("/cuti", leaveController.AddCuti)
	leave.Get("/cuti", leaveController.FindCuti)
	leave.Get("/balance/:employee_nik", leaveController.Balance)
	leave.Get("/balance", leaveController.Balance)
	leave.Post("/transactions", leaveController.CreateTransaction)

	leaveType := api.Group("/leave-type", middlewares.JWTProtected())
	leaveType.Get("/", leaveTypeController.FindAll)
	leaveTypes := api.Group("/leave-types", middlewares.JWTProtected())
	leaveTypes.Get("/", leaveTypeController.FindAll)

	approval := api.Group("/approval", middlewares.JWTProtected())
	approval.Get("/", approvalController.FindAll)
	approval.Get("/:id", approvalController.Detail)
	approval.Delete("/:id", approvalController.Delete)
	approval.Post("/approve", approvalController.Approve)

	template := api.Group("/approval_template", middlewares.JWTProtected())
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
