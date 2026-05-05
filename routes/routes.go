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
	companiesRepo := repositories.NewCompaniesRepository(db)
	paramRepo := repositories.NewParamRepository(db)
	paramGroupRepo := repositories.NewParamGroupRepository(db)
	settingRepo := repositories.NewSettingRepository(db)
	attendanceRepo := repositories.NewAttendanceRepository(db)
	officeRepo := repositories.NewOfficeRepository(db)
	approvalRepo := repositories.NewApprovalRepository(db)
	approvalTemplateRepo := repositories.NewApprovalTemplateRepository(db)
	branchRepo := repositories.NewBranchRepository(db)
	userCompanyRepo := repositories.NewUserCompanyRepository(db)
	userOfficeRepo := repositories.NewUserOfficeRepository(db)
	userCustomerRepo := repositories.NewUserCustomerRepository(db)
	divisionRepo := repositories.NewDivisionRepository(db)
	departmentRepo := repositories.NewDepartmentRepository(db)
	titleRepo := repositories.NewTitleRepository(db)
	userGroupRepo := repositories.NewUserGroupRepository(db)

	auth := controllers.NewAuthController(repo)
	user := controllers.NewUserController(repo)
	setup := controllers.NewSetupController(repo)
	customer := controllers.NewCustomerController(customerRepo)
	companies := controllers.NewCompaniesController(companiesRepo)
	paramController := controllers.NewParamController(paramRepo)
	paramGroupController := controllers.NewParamGroupController(paramGroupRepo)
	settingController := controllers.NewSettingController(settingRepo)
	attendanceController := controllers.NewAttendanceController(attendanceRepo)
	officeController := controllers.NewOfficeController(officeRepo)
	approvalController := controllers.NewApprovalController(approvalRepo)
	approvalTemplateController := controllers.NewApprovalTemplateController(approvalTemplateRepo)
	branchController := controllers.NewBranchController(branchRepo)
	userCompanyController := controllers.NewUserCompanyController(userCompanyRepo)
	userOfficeController := controllers.NewUserOfficeController(userOfficeRepo)
	userCustomerController := controllers.NewUserCustomerController(userCustomerRepo)
	divisionController := controllers.NewDivisionController(divisionRepo)
	departmentController := controllers.NewDepartmentController(departmentRepo)
	titleController := controllers.NewTitleController(titleRepo)
	userGroupController := controllers.NewUserGroupController(userGroupRepo)

	api := app.Group("/api")

	api.Post("/setup", setup.InitAdmin)
	api.Post("/login", auth.Login)

	users := api.Group("/users", middlewares.JWTProtected())
	users.Post("/", user.Create)
	users.Get("/", user.FindAll)
	users.Get("/me", user.Me)

	customerRoute := api.Group("/customer", middlewares.JWTProtected())
	customerRoute.Post("/", customer.Create)
	customerRoute.Get("/", customer.FindAll)
	customerRoute.Put("/", customer.Update)
	customerRoute.Delete("/:id", customer.Delete)

	companiesRoute := api.Group("/companies", middlewares.JWTProtected())
	companiesRoute.Post("/", companies.Create)
	companiesRoute.Get("/", companies.FindAll)
	companiesRoute.Put("/", companies.Update)
	companiesRoute.Delete("/:id", companies.Delete)

	branch := api.Group("/branch", middlewares.JWTProtected())
	branch.Post("/", branchController.Create)
	branch.Get("/", branchController.FindAll)
	branch.Put("/", branchController.Update)
	branch.Delete("/:id", branchController.Delete)

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

	office := api.Group("/office", middlewares.JWTProtected())
	office.Post("/", officeController.Create)
	office.Get("/", officeController.FindAll)
	office.Put("/", officeController.Update)
	office.Delete("/:id", officeController.Delete)

	userCompany := api.Group("/user-company", middlewares.JWTProtected())
	userCompany.Post("/", userCompanyController.Create)
	userCompany.Get("/", userCompanyController.FindAll)
	userCompany.Put("/", userCompanyController.Update)
	userCompany.Delete("/:id", userCompanyController.Delete)

	userOffice := api.Group("/user-office", middlewares.JWTProtected())
	userOffice.Post("/", userOfficeController.Create)
	userOffice.Get("/", userOfficeController.FindAll)
	userOffice.Put("/", userOfficeController.Update)
	userOffice.Delete("/:id", userOfficeController.Delete)

	userCustomer := api.Group("/user-customer", middlewares.JWTProtected())
	userCustomer.Post("/", userCustomerController.Create)
	userCustomer.Get("/", userCustomerController.FindAll)
	userCustomer.Put("/", userCustomerController.Update)
	userCustomer.Delete("/:id", userCustomerController.Delete)

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
