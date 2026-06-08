package controllers

import (
	"strings"

	"hrms_go/dto"
	"hrms_go/models"
	"hrms_go/repositories"
	"hrms_go/utils"

	"github.com/gofiber/fiber/v3"
)

type UserSiteController struct {
	repo repositories.UserSiteRepository
}

func NewUserSiteController(repo repositories.UserSiteRepository) *UserSiteController {
	return &UserSiteController{repo}
}

// Create User Site godoc
// @Summary Create user site mapping
// @Description Create mapping between employee and site
// @Tags User Site
// @Accept json
// @Produce json
// @Param request body models.UserSite true "User Site data"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/user-site [post]
func (c *UserSiteController) Create(ctx fiber.Ctx) error {
	var data models.UserSite
	if err := ctx.Bind().Body(&data); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	data.SiteType = strings.ToUpper(strings.TrimSpace(data.SiteType))
	if !isValidSiteType(data.SiteType) {
		return utils.Error(ctx, 400, "site_type must be BRANCH, OFFICE, or CUSTOMER")
	}
	data.CreatedBy = ctx.Locals("employee_nik").(string)

	if err := c.repo.Create(&data); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}
	return utils.Success(ctx, data)
}

// Get All User Site godoc
// @Summary Get all user site mappings
// @Description Get list of user site mappings with pagination
// @Tags User Site
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/user-site [get]
func (c *UserSiteController) FindAll(ctx fiber.Ctx) error {
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

// Update User Site godoc
// @Summary Update user site mapping
// @Description Update existing user site mapping
// @Tags User Site
// @Accept json
// @Produce json
// @Param request body models.UserSite true "User Site data"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/user-site [put]
func (c *UserSiteController) Update(ctx fiber.Ctx) error {
	var data models.UserSite
	if err := ctx.Bind().Body(&data); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	data.SiteType = strings.ToUpper(strings.TrimSpace(data.SiteType))
	if data.SiteType != "" && !isValidSiteType(data.SiteType) {
		return utils.Error(ctx, 400, "site_type must be BRANCH, OFFICE, or CUSTOMER")
	}
	userID := ctx.Locals("employee_nik").(string)
	data.UpdatedBy = &userID

	if err := c.repo.Update(&data); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}
	return utils.Success(ctx, data)
}

// Delete User Site godoc
// @Summary Delete user site mapping
// @Description Delete user site mapping by ID
// @Tags User Site
// @Accept json
// @Produce json
// @Param id path string true "User Site ID"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/user-site/{id} [delete]
func (c *UserSiteController) Delete(ctx fiber.Ctx) error {
	if err := c.repo.Delete(ctx.Params("id")); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}
	return utils.Success(ctx, nil)
}
