package controllers

import (
	"hrms_go/dto"
	"hrms_go/models"
	"hrms_go/repositories"
	"hrms_go/utils"

	"github.com/gofiber/fiber/v3"
)

type TitleController struct {
	repo repositories.TitleRepository
}

func NewTitleController(repo repositories.TitleRepository) *TitleController {
	return &TitleController{repo}
}

// Create Title godoc
// @Summary Create title
// @Description Create new title
// @Tags Title
// @Accept json
// @Produce json
// @Param request body models.Title true "Title data"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/title [post]
func (c *TitleController) Create(ctx fiber.Ctx) error {
	var data models.Title
	if err := ctx.Bind().Body(&data); err != nil {
		return utils.Error(ctx, 400, err.Error())
	}
	data.CreatedBy = ctx.Locals("user_id").(string)
	if err := c.repo.Create(&data); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}
	return utils.Success(ctx, data)
}

// Get All Title godoc
// @Summary Get all titles
// @Description Get list of titles with pagination
// @Tags Title
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
// @Router /api/title [get]
func (c *TitleController) FindAll(ctx fiber.Ctx) error {
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

// Update Title godoc
// @Summary Update title
// @Description Update existing title
// @Tags Title
// @Accept json
// @Produce json
// @Param request body models.Title true "Title data"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/title [put]
func (c *TitleController) Update(ctx fiber.Ctx) error {
	var data models.Title
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

// Delete Title godoc
// @Summary Delete title
// @Description Delete title by ID
// @Tags Title
// @Accept json
// @Produce json
// @Param id path string true "Title ID"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/title/{id} [delete]
func (c *TitleController) Delete(ctx fiber.Ctx) error {
	if err := c.repo.Delete(ctx.Params("id")); err != nil {
		return utils.Error(ctx, 500, err.Error())
	}
	return utils.Success(ctx, nil)
}
