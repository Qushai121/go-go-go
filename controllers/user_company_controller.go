package controllers

import (
	"hrms_go/dto"
	"hrms_go/models"
	"hrms_go/repositories"
	"hrms_go/utils"

	"github.com/gofiber/fiber/v3"
)

type UserCompanyController struct {
	repo repositories.UserCompanyRepository
}

func NewUserCompanyController(repo repositories.UserCompanyRepository) *UserCompanyController {
	return &UserCompanyController{repo}
}

// Create User Company godoc
// @Summary Create user company
// @Description Create new user company
// @Tags User Company
// @Accept json
// @Produce json
// @Param request body models.UserCompany true "User Company data"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/user-company [post]
func (c *UserCompanyController) Create(ctx fiber.Ctx) error {
	var data models.UserCompany
	if err := ctx.Bind().Body(&data); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}
	data.CreatedBy = ctx.Locals("user_id").(string)
	if err := c.repo.Create(&data); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}
	return utils.Success(ctx, data)
}

// Get All User Company godoc
// @Summary Get all user companies
// @Description Get list of user companies with pagination
// @Tags User Company
// @Accept json
// @Produce json
// @Param sort_order query string false "Sort order (asc/desc)"
// @Param sort_by query string false "Sort by column"
// @Param search query string false "Search keyword"
// @Param field_search query string false "Field name for dynamic search"
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/user-company [get]
func (c *UserCompanyController) FindAll(ctx fiber.Ctx) error {
	queryParams := dto.PaginateFieldDto{}
	if err := utils.BindPaginationParams(ctx, &queryParams); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}
	data, err := c.repo.FindAll(&queryParams)
	if err != nil {
		return utils.Error(ctx, 500, err.Error())
	}
	return utils.Success(ctx, data)
}

// Update User Company godoc
// @Summary Update user company
// @Description Update existing user company
// @Tags User Company
// @Accept json
// @Produce json
// @Param request body models.UserCompany true "User Company data"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/user-company [put]
func (c *UserCompanyController) Update(ctx fiber.Ctx) error {
	var data models.UserCompany
	if err := ctx.Bind().Body(&data); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}
	userID := ctx.Locals("user_id").(string)
	data.UpdatedBy = &userID
	if err := c.repo.Update(&data); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}
	return utils.Success(ctx, data)
}

// Delete User Company godoc
// @Summary Delete user company
// @Description Delete user company by ID
// @Tags User Company
// @Accept json
// @Produce json
// @Param id path string true "User Company ID"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/user-company/{id} [delete]
func (c *UserCompanyController) Delete(ctx fiber.Ctx) error {
	if err := c.repo.Delete(ctx.Params("id")); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}
	return utils.Success(ctx, nil)
}
