package controllers

import (
	"hrms_go/dto"
	"hrms_go/models"
	"hrms_go/repositories"
	"hrms_go/utils"

	"github.com/gofiber/fiber/v3"
)

type UserOfficeController struct {
	repo repositories.UserOfficeRepository
}

func NewUserOfficeController(repo repositories.UserOfficeRepository) *UserOfficeController {
	return &UserOfficeController{repo}
}

// Create User Office godoc
// @Summary Create user office
// @Description Create new user office
// @Tags User Office
// @Accept json
// @Produce json
// @Param request body models.UserOffice true "User Office data"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/user-office [post]
func (c *UserOfficeController) Create(ctx fiber.Ctx) error {
	var data models.UserOffice
	if err := ctx.Bind().Body(&data); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}
	data.CreatedBy = ctx.Locals("employee_nik").(string)
	if err := c.repo.Create(&data); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}
	return utils.Success(ctx, data)
}

// Get All User Office godoc
// @Summary Get all user offices
// @Description Get list of user offices with pagination
// @Tags User Office
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
// @Router /api/user-office [get]
func (c *UserOfficeController) FindAll(ctx fiber.Ctx) error {
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

// Update User Office godoc
// @Summary Update user office
// @Description Update existing user office
// @Tags User Office
// @Accept json
// @Produce json
// @Param request body models.UserOffice true "User Office data"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/user-office [put]
func (c *UserOfficeController) Update(ctx fiber.Ctx) error {
	var data models.UserOffice
	if err := ctx.Bind().Body(&data); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}
	userID := ctx.Locals("employee_nik").(string)
	data.UpdatedBy = &userID
	if err := c.repo.Update(&data); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}
	return utils.Success(ctx, data)
}

// Delete User Office godoc
// @Summary Delete user office
// @Description Delete user office by ID
// @Tags User Office
// @Accept json
// @Produce json
// @Param id path string true "User Office ID"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/user-office/{id} [delete]
func (c *UserOfficeController) Delete(ctx fiber.Ctx) error {
	if err := c.repo.Delete(ctx.Params("id")); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}
	return utils.Success(ctx, nil)
}
