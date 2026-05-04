package controllers

import (
	"hrms_go/dto"
	"hrms_go/models"
	"hrms_go/repositories"
	"hrms_go/utils"

	"github.com/gofiber/fiber/v3"
)

type SettingController struct {
	repo repositories.SettingRepository
}

func NewSettingController(repo repositories.SettingRepository) *SettingController {
	return &SettingController{repo}
}

// Create Setting godoc
// @Summary Create setting
// @Description Create new setting
// @Tags Setting
// @Accept json
// @Produce json
// @Param request body models.Setting true "Setting data"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/setting [post]
func (c *SettingController) Create(ctx fiber.Ctx) error {
	var setting models.Setting

	if err := ctx.Bind().Body(&setting); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	userId := ctx.Locals("user_id").(string)
	setting.CreatedBy = userId

	if err := c.repo.Create(&setting); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, setting)
}

// Get All Setting godoc
// @Summary Get all settings
// @Description Get list of settings with pagination
// @Tags Setting
// @Accept json
// @Produce json
// @Param sort_order query string false "Sort order (asc/desc)"
// @Param sort_by query string false "Sort by column"
// @Param search query string false "Search keyword"
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/setting [get]
func (c *SettingController) FindAll(ctx fiber.Ctx) error {
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

// Update Setting godoc
// @Summary Update setting
// @Description Update existing setting
// @Tags Setting
// @Accept json
// @Produce json
// @Param request body models.Setting true "Setting data"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/setting [put]
func (c *SettingController) Update(ctx fiber.Ctx) error {
	data := models.Setting{}

	if err := ctx.Bind().Body(&data); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}

	userId := ctx.Locals("user_id").(string)
	data.UpdatedBy = &userId

	if err := c.repo.Update(&data); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, data)
}

// Delete Setting godoc
// @Summary Delete setting
// @Description Delete setting by ID
// @Tags Setting
// @Accept json
// @Produce json
// @Param id path string true "Setting ID"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/setting/{id} [delete]
func (c *SettingController) Delete(ctx fiber.Ctx) error {
	id := ctx.Params("id")

	if err := c.repo.Delete(id); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}

	return utils.Success(ctx, nil)
}
