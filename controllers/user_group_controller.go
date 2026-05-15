package controllers

import (
	"hrms_go/dto"
	"hrms_go/models"
	"hrms_go/repositories"
	"hrms_go/utils"

	"github.com/gofiber/fiber/v3"
)

type UserGroupController struct {
	repo repositories.UserGroupRepository
}

func NewUserGroupController(repo repositories.UserGroupRepository) *UserGroupController {
	return &UserGroupController{repo}
}

// Create User Group godoc
// @Summary Create user group
// @Description Create new user group
// @Tags User Group
// @Accept json
// @Produce json
// @Param request body models.UserGroup true "User Group data"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/usergroup [post]
func (c *UserGroupController) Create(ctx fiber.Ctx) error {
	var data models.UserGroup
	if err := ctx.Bind().Body(&data); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}
	data.CreatedBy = ctx.Locals("user_id").(string)
	if err := c.repo.Create(&data); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}
	return utils.Success(ctx, data)
}

// Get All User Group godoc
// @Summary Get all user groups
// @Description Get list of user groups with pagination
// @Tags User Group
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
// @Router /api/usergroup [get]
func (c *UserGroupController) FindAll(ctx fiber.Ctx) error {
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

// Update User Group godoc
// @Summary Update user group
// @Description Update existing user group
// @Tags User Group
// @Accept json
// @Produce json
// @Param request body models.UserGroup true "User Group data"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/usergroup [put]
func (c *UserGroupController) Update(ctx fiber.Ctx) error {
	var data models.UserGroup
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

// Delete User Group godoc
// @Summary Delete user group
// @Description Delete user group by ID
// @Tags User Group
// @Accept json
// @Produce json
// @Param id path string true "User Group ID"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/usergroup/{id} [delete]
func (c *UserGroupController) Delete(ctx fiber.Ctx) error {
	if err := c.repo.Delete(ctx.Params("id")); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}
	return utils.Success(ctx, nil)
}
